package harness

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func tokenCounts(usage map[string]any) (int, int, float64) {
	source := usage
	if nested := mapValue(usage["usage"]); len(nested) > 0 {
		source = nested
	}
	first := func(keys ...string) any {
		for _, key := range keys {
			if source[key] != nil {
				return source[key]
			}
		}
		return nil
	}
	input := intValue(first("input_tokens", "inputTokens", "input"))
	output := intValue(first("output_tokens", "outputTokens", "output")) + intValue(source["reasoning"])
	cost := floatValue(usage["total_cost_usd"])
	if cost == 0 {
		cost = floatValue(usage["cost_usd"])
	}
	if cost == 0 {
		if costs := mapValue(source["cost"]); len(costs) > 0 {
			cost = floatValue(costs["total"])
		}
	}
	if input == 0 && output == 0 {
		for _, raw := range mapValue(usage["modelUsage"]) {
			model := mapValue(raw)
			input += intValue(model["inputTokens"])
			output += intValue(model["outputTokens"])
			cost += floatValue(model["costUSD"])
		}
	}
	return input, output, cost
}
func (h *Harness) Summarize(runArgument string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	var plan Plan
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	statuses := map[string]int{}
	providers := map[string]map[string]any{}
	totalFindings := 0
	for _, call := range plan.Calls {
		record, err := readCallRecord(run, call.CallID)
		if err != nil {
			return err
		}
		status := "missing"
		if record != nil {
			status = stringValue(record["status"])
		}
		statuses[status]++
		key := call.Provider.Adapter + ":" + call.Provider.Model
		bucket := providers[key]
		if bucket == nil {
			bucket = map[string]any{"calls": 0, "latencies": []float64{}, "input_tokens": 0, "output_tokens": 0, "cost_usd": 0.0}
			providers[key] = bucket
		}
		bucket["calls"] = intValue(bucket["calls"]) + 1
		if record != nil {
			if record["latency_seconds"] != nil {
				bucket["latencies"] = append(bucket["latencies"].([]float64), floatValue(record["latency_seconds"]))
			}
			in, out, cost := tokenCounts(mapValue(record["usage"]))
			bucket["input_tokens"] = intValue(bucket["input_tokens"]) + in
			bucket["output_tokens"] = intValue(bucket["output_tokens"]) + out
			bucket["cost_usd"] = floatValue(bucket["cost_usd"]) + cost
			if status == "success" {
				totalFindings += len(sliceValue(mapValue(record["parsed_response"])["findings"]))
			}
		}
	}
	outputProviders := map[string]any{}
	for key, bucket := range providers {
		latencies := bucket["latencies"].([]float64)
		sum := 0.0
		for _, v := range latencies {
			sum += v
		}
		meanValue := any(nil)
		if len(latencies) > 0 {
			meanValue = sum / float64(len(latencies))
		}
		delete(bucket, "latencies")
		bucket["mean_latency_seconds"] = meanValue
		bucket["total_latency_seconds"] = sum
		outputProviders[key] = bucket
	}
	summary := map[string]any{"statuses": statuses, "total_findings_before_set_union": totalFindings, "providers": outputProviders}
	analysis := filepath.Join(run, "analysis")
	if err := os.MkdirAll(analysis, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(analysis, "run-health.json"), summary); err != nil {
		return err
	}
	var lines []string
	lines = append(lines, "# Run health", "", fmt.Sprintf("- Calls: %d", len(plan.Calls)), fmt.Sprintf("- Raw findings: %d", totalFindings))
	keys := make([]string, 0, len(statuses))
	for key := range statuses {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("- %s: %d", key, statuses[key]))
	}
	for key, raw := range outputProviders {
		bucket := mapValue(raw)
		lines = append(lines, "", "## "+key, fmt.Sprintf("- Input tokens: %d", intValue(bucket["input_tokens"])), fmt.Sprintf("- Output tokens: %d", intValue(bucket["output_tokens"])), fmt.Sprintf("- Recorded cost: $%.4f", floatValue(bucket["cost_usd"])), fmt.Sprintf("- Mean latency: %v", bucket["mean_latency_seconds"]))
	}
	if err := writeText(filepath.Join(analysis, "run-health.md"), strings.Join(lines, "\n")+"\n"); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(data))
	return nil
}

var rawNames = map[string]bool{"freeze.json": true, "config.json": true, "plan.json": true, "run-summary.json": true}

func rawFiles(run string) ([]string, error) {
	files := []string{}
	for name := range rawNames {
		path := filepath.Join(run, name)
		if stat, err := os.Stat(path); err == nil && stat.Mode().IsRegular() {
			files = append(files, path)
		}
	}
	calls := filepath.Join(run, "calls")
	_ = filepath.WalkDir(calls, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files, nil
}
func (h *Harness) Seal(runArgument string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(run, "seal.json")); err == nil {
		return fmt.Errorf("run is already sealed")
	}
	var plan Plan
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	incomplete := []string{}
	for _, call := range plan.Calls {
		record, err := readCallRecord(run, call.CallID)
		if err != nil {
			return err
		}
		if record == nil || !terminalStatuses[stringValue(record["status"])] {
			incomplete = append(incomplete, call.CallID)
		}
	}
	if len(incomplete) > 0 {
		return fmt.Errorf("cannot seal; incomplete calls: %s", strings.Join(incomplete, ", "))
	}
	files, _ := rawFiles(run)
	manifest := []map[string]any{}
	for _, path := range files {
		digest, err := shaFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(run, path)
		manifest = append(manifest, map[string]any{"path": filepath.ToSlash(rel), "sha256": digest})
	}
	if err := writeJSON(filepath.Join(run, "seal.json"), map[string]any{"schema_version": 2, "sealed_at": utcNow(), "scope": "freeze/config/plan/run-summary and calls/**; derived blinded and analysis files are excluded", "files": manifest}); err != nil {
		return err
	}
	for _, path := range files {
		stat, _ := os.Stat(path)
		_ = os.Chmod(path, stat.Mode()&^0o222)
	}
	fmt.Printf("Sealed %d raw files in %s\n", len(manifest), h.relative(run))
	return nil
}
func (h *Harness) Verify(runArgument string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	sealPath := filepath.Join(run, "seal.json")
	if _, err := os.Stat(sealPath); err != nil {
		return fmt.Errorf("run is not sealed")
	}
	var seal struct {
		Files []Asset `json:"files"`
	}
	if err := readJSON(sealPath, &seal); err != nil {
		return err
	}
	expected := map[string]bool{}
	errorsFound := []string{}
	for _, row := range seal.Files {
		expected[row.Path] = true
		path := filepath.Join(run, filepath.FromSlash(row.Path))
		digest, err := shaFile(path)
		if err != nil {
			errorsFound = append(errorsFound, "missing: "+row.Path)
		} else if digest != row.SHA256 {
			errorsFound = append(errorsFound, "digest mismatch: "+row.Path)
		}
	}
	calls := filepath.Join(run, "calls")
	_ = filepath.WalkDir(calls, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(run, path)
			rel = filepath.ToSlash(rel)
			if !expected[rel] {
				errorsFound = append(errorsFound, "unsealed raw call file: "+rel)
			}
		}
		return nil
	})
	if len(errorsFound) > 0 {
		fmt.Println("Integrity verification failed:")
		for _, value := range errorsFound {
			fmt.Println("- " + value)
		}
		return fmt.Errorf("integrity verification failed")
	}
	fmt.Printf("Integrity verification passed (%d raw files).\n", len(expected))
	return nil
}

func findingsForSet(run string, set OutputSet) ([]map[string]any, error) {
	out := []map[string]any{}
	for _, id := range set.CallIDs {
		record, err := readCallRecord(run, id)
		if err != nil {
			return nil, err
		}
		if record == nil || stringValue(record["status"]) != "success" {
			continue
		}
		for _, raw := range sliceValue(mapValue(record["parsed_response"])["findings"]) {
			finding := mapValue(raw)
			copyFinding := map[string]any{}
			for k, v := range finding {
				copyFinding[k] = v
			}
			out = append(out, copyFinding)
		}
	}
	return out, nil
}
func writeJSONL(path string, rows []map[string]any) error {
	var b strings.Builder
	for _, row := range rows {
		data, err := json.Marshal(row)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteByte('\n')
	}
	return writeText(path, b.String())
}
func blindedID(key []byte, prefix string, parts ...any) string {
	values := []string{prefix}
	for _, part := range parts {
		values = append(values, fmt.Sprint(part))
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(strings.Join(values, "\x1f")))
	return prefix + "-" + hex.EncodeToString(mac.Sum(nil))[:16]
}

func loadOrCreateBlindingKey(run string) ([]byte, error) {
	private := filepath.Join(run, "private")
	if err := os.MkdirAll(private, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(private, "blinding-key")
	if value, err := os.ReadFile(path); err == nil {
		if len(value) != 32 {
			return nil, fmt.Errorf("invalid private blinding key length")
		}
		return value, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return nil, err
	}
	if err := atomicWrite(path, value, 0o600); err != nil {
		return nil, err
	}
	return value, nil
}

func blindedOrder(rows []map[string]any, key []byte, field string) {
	sort.SliceStable(rows, func(i, j int) bool {
		return blindedID(key, "order", stringValue(rows[i][field])) < blindedID(key, "order", stringValue(rows[j][field]))
	})
}

type comparisonMember struct {
	condition, role, kind, group string
	set                          OutputSet
}

func comparisonFor(set OutputSet) (comparisonMember, bool) {
	metadata := set.Metadata
	design := stringValue(metadata["design"])
	if design == "provider_pair" || design == "persona_factorial" {
		condition := stringValue(metadata["wrapper"])
		if condition != "functional" && condition != "fictional" {
			return comparisonMember{}, false
		}
		role, kind := stringValue(metadata["role"]), stringValue(metadata["kind"])
		group := pairKey(design, set.Packet, metadata["repeat"], metadata["provider_index"], metadata["model"], role, kind)
		return comparisonMember{condition: condition, role: role, kind: kind, group: group, set: set}, true
	}
	if design != "stage_a" {
		return comparisonMember{}, false
	}
	arm, kind := stringValue(metadata["arm"]), stringValue(metadata["kind"])
	condition, role := "", ""
	switch arm {
	case "S1":
		condition, role = "functional", "omnibus"
	case "S2":
		condition, role = "fictional", "omnibus"
	case "M1":
		condition, role = "functional", "specialist-panel"
	case "M2":
		condition, role = "fictional", "specialist-panel"
	}
	if condition == "" {
		return comparisonMember{}, false
	}
	group := pairKey(design, set.Packet, metadata["repeat"], role, kind)
	return comparisonMember{condition: condition, role: role, kind: kind, group: group, set: set}, true
}

func writeTemplate(path string, header []string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	_ = writer.Write(header)
	for _, row := range rows {
		_ = writer.Write(row)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (h *Harness) Bundle(runArgument string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	var plan Plan
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	key, err := loadOrCreateBlindingKey(run)
	if err != nil {
		return err
	}
	blinded := filepath.Join(run, "blinded")
	_ = os.RemoveAll(blinded)
	if err := os.MkdirAll(blinded, 0o755); err != nil {
		return err
	}
	items, sets, pairs := []map[string]any{}, []map[string]any{}, []map[string]any{}
	unblindSets, unblindItems, unblindPairs := map[string]any{}, map[string]any{}, map[string]any{}
	packets := map[string]bool{}
	blindedSetByActual := map[string]string{}
	itemRowsBySet := map[string][]map[string]any{}
	for _, set := range plan.OutputSets {
		blindSetID := blindedID(key, "s", set.SetID)
		blindedSetByActual[set.SetID] = blindSetID
		findings, err := findingsForSet(run, set)
		if err != nil {
			return err
		}
		ids := []string{}
		for index, finding := range findings {
			id := blindedID(key, "i", set.SetID, index)
			ids = append(ids, id)
			item := map[string]any{"item_id": id, "set_id": blindSetID, "packet": set.Packet, "finding": finding}
			items = append(items, item)
			itemRowsBySet[blindSetID] = append(itemRowsBySet[blindSetID], map[string]any{"item_id": id, "finding": finding})
			unblindItems[id] = map[string]any{"set_id": blindSetID, "source_set_id": set.SetID, "source_index": index}
		}
		sets = append(sets, map[string]any{"set_id": blindSetID, "packet": set.Packet, "item_ids": ids})
		unblindSets[blindSetID] = set
		packets[set.Packet] = true
	}
	comparisonGroups := map[string]map[string]comparisonMember{}
	for _, set := range plan.OutputSets {
		member, ok := comparisonFor(set)
		if !ok {
			continue
		}
		if comparisonGroups[member.group] == nil {
			comparisonGroups[member.group] = map[string]comparisonMember{}
		}
		comparisonGroups[member.group][member.condition] = member
	}
	for group, members := range comparisonGroups {
		functional, hasFunctional := members["functional"]
		fictional, hasFictional := members["fictional"]
		if !hasFunctional || !hasFictional {
			continue
		}
		pairID := blindedID(key, "p", group)
		left, right := functional, fictional
		if h := hmac.New(sha256.New, key); true {
			_, _ = h.Write([]byte("side\x1f" + group))
			if h.Sum(nil)[0]&1 == 1 {
				left, right = fictional, functional
			}
		}
		leftSet, rightSet := blindedSetByActual[left.set.SetID], blindedSetByActual[right.set.SetID]
		pairs = append(pairs, map[string]any{"pair_id": pairID, "packet": left.set.Packet, "role": left.role, "kind": left.kind, "left": map[string]any{"output_id": leftSet, "findings": itemRowsBySet[leftSet]}, "right": map[string]any{"output_id": rightSet, "findings": itemRowsBySet[rightSet]}})
		unblindPairs[pairID] = map[string]any{"left_condition": left.condition, "right_condition": right.condition, "left_set_id": left.set.SetID, "right_set_id": right.set.SetID}
	}
	blindedOrder(items, key, "item_id")
	blindedOrder(sets, key, "set_id")
	blindedOrder(pairs, key, "pair_id")
	if err := writeJSONL(filepath.Join(blinded, "findings.jsonl"), items); err != nil {
		return err
	}
	if err := writeJSONL(filepath.Join(blinded, "sets.jsonl"), sets); err != nil {
		return err
	}
	if err := writeJSONL(filepath.Join(blinded, "pairs.jsonl"), pairs); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(run, "private", "unblind.json"), map[string]any{"sets": unblindSets, "items": unblindItems, "pairs": unblindPairs}); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(blinded, "manifest.json"), map[string]any{"schema_version": 2, "id_scheme": "HMAC-SHA-256 with private random key", "blinding_key_sha256": shaBytes(key), "condition_labels_hidden": true, "wording_may_reveal_treatment": true, "item_count": len(items), "set_count": len(sets), "pair_count": len(pairs)}); err != nil {
		return err
	}
	answerDir, packetDir := filepath.Join(blinded, "answer-keys"), filepath.Join(blinded, "review-packets")
	if err := os.MkdirAll(answerDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(packetDir, 0o755); err != nil {
		return err
	}
	for packet := range packets {
		if err := copyFile(filepath.Join(h.Root, "experiment/scenarios", packet, "answer-key.md"), filepath.Join(answerDir, packet+".md"), 0o644); err != nil {
			return err
		}
		if err := copyFile(filepath.Join(h.Root, "experiment/scenarios", packet, "review-packet.md"), filepath.Join(packetDir, packet+".md"), 0o644); err != nil {
			return err
		}
	}
	itemTemplate := [][]string{}
	for _, row := range items {
		itemTemplate = append(itemTemplate, []string{"", stringValue(row["item_id"]), "", "", "", "", ""})
	}
	if err := writeTemplate(filepath.Join(blinded, "rating-template.csv"), []string{"rater", "item_id", "defect_id", "false_positive_cluster", "material", "confidence", "notes"}, itemTemplate); err != nil {
		return err
	}
	pairTemplate := [][]string{}
	for _, row := range pairs {
		pairTemplate = append(pairTemplate, []string{"", stringValue(row["pair_id"]), "", "", "", "", "", "", "", "", "", "", ""})
	}
	if err := writeTemplate(filepath.Join(blinded, "pairwise-rating-template.csv"), []string{"rater", "pair_id", "left_supportedness_1_5", "right_supportedness_1_5", "left_actionability_1_5", "right_actionability_1_5", "left_fix_quality_1_5", "right_fix_quality_1_5", "overall_preference_left_right_tie", "condition_guess_left_functional_left_fictional_unsure", "guess_confidence_1_5", "wording_revealed_condition_true_false_unsure", "notes"}, pairTemplate); err != nil {
		return err
	}
	runsheet := "# Blinded human rating runsheet\n\n## Phase 1 — independent defect mapping\n\n1. Do not open `plan.json`, `calls/`, or `private/`. Do not derive IDs from public configs.\n2. Rate `findings.jsonl` one item at a time in its shuffled order; do not compare paired outputs yet.\n3. Use the packet and answer key for the item's `packet`. Map to one defect ID or `NONE`.\n4. Cluster semantically duplicate false claims within the same opaque set. Record materiality, confidence, and notes.\n5. Finish and lock your Phase 1 CSV before opening `pairs.jsonl`.\n\n## Phase 2 — paired qualitative judgement\n\n6. `pairs.jsonl` randomises functional/fictional output between left and right. Score supportedness, actionability, and fix quality from 1 (poor) to 5 (excellent).\n7. Choose left, right, or tie overall. Only after scoring, guess which side is functional, or choose unsure; record confidence and whether wording revealed the condition.\n8. Wording is not transformed because style and token use are part of the treatment. Therefore this is label-blinded, not guaranteed treatment-blinded.\n\n## Independence and unblinding\n\n9. Work independently. Disclose relevant priors; one rater's preference for the fictional robots is acceptable but must not be the only rating.\n10. Use two raters. Reconcile Phase 1 disagreements before unblinding. Preserve both original rating files.\n11. The controller joins opaque IDs to conditions only after both ratings are locked. Do not delete the hashes; publish an unblinded derived table so the audit trail survives.\n"
	if err := writeText(filepath.Join(blinded, "RUNSHEET.md"), runsheet); err != nil {
		return err
	}
	fmt.Printf("Created %d blinded items across %d output sets and %d paired comparisons.\n", len(items), len(sets), len(pairs))
	return nil
}
