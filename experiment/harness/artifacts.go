package harness

import (
	"encoding/csv"
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
func deterministicOrder(rows []map[string]any, seed, key string) {
	sort.SliceStable(rows, func(i, j int) bool {
		return shaText(seed+"\x1f"+stringValue(rows[i][key])) < shaText(seed+"\x1f"+stringValue(rows[j][key]))
	})
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
	blinded := filepath.Join(run, "blinded")
	_ = os.RemoveAll(blinded)
	if err := os.MkdirAll(blinded, 0o755); err != nil {
		return err
	}
	seed := plan.Seed + ":blinding"
	items, sets := []map[string]any{}, []map[string]any{}
	unblindSets, unblindItems := map[string]any{}, map[string]any{}
	packets := map[string]bool{}
	for _, set := range plan.OutputSets {
		findings, err := findingsForSet(run, set)
		if err != nil {
			return err
		}
		ids := []string{}
		for index, finding := range findings {
			id := "i-" + opaqueID(seed, 16, set.SetID, index)
			ids = append(ids, id)
			items = append(items, map[string]any{"item_id": id, "set_id": set.SetID, "packet": set.Packet, "finding": finding})
			unblindItems[id] = map[string]any{"set_id": set.SetID, "source_index": index}
		}
		sets = append(sets, map[string]any{"set_id": set.SetID, "packet": set.Packet, "item_ids": ids})
		unblindSets[set.SetID] = set
		packets[set.Packet] = true
	}
	deterministicOrder(items, seed, "item_id")
	deterministicOrder(sets, seed+":sets", "set_id")
	if err := writeJSONL(filepath.Join(blinded, "findings.jsonl"), items); err != nil {
		return err
	}
	if err := writeJSONL(filepath.Join(blinded, "sets.jsonl"), sets); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(run, "private", "unblind.json"), map[string]any{"sets": unblindSets, "items": unblindItems}); err != nil {
		return err
	}
	answerDir := filepath.Join(blinded, "answer-keys")
	if err := os.MkdirAll(answerDir, 0o755); err != nil {
		return err
	}
	for packet := range packets {
		if err := copyFile(filepath.Join(h.Root, "experiment/scenarios", packet, "answer-key.md"), filepath.Join(answerDir, packet+".md"), 0o644); err != nil {
			return err
		}
	}
	f, err := os.Create(filepath.Join(blinded, "rating-template.csv"))
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	_ = writer.Write([]string{"rater", "item_id", "defect_id", "false_positive_cluster", "material", "confidence", "notes"})
	for _, row := range items {
		_ = writer.Write([]string{"", stringValue(row["item_id"]), "", "", "", "", ""})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	runsheet := "# Blinded rating runsheet\n\n1. Do not open `plan.json`, `calls/`, or `private/unblind.json`.\n2. Open `findings.jsonl` and the answer key for each item's `packet`.\n3. Enter one row per item in a copy of `rating-template.csv`.\n4. Set `defect_id` to one matching planted ID, or `NONE`.\n5. For `NONE`, assign a short `false_positive_cluster`; reuse it for the same claimed mechanism in this set.\n6. Set `material` to `true` only when the packet supports the claim and its consequence.\n7. Use `confidence` = `high`, `medium`, or `low`; explain ambiguous matches in `notes`.\n8. Work independently. Do not compare ratings until both raters finish.\n9. Resolve disagreements without revealing arm, model, or provider metadata.\n10. Save the adjudicated file and run the scoring recipe.\n\nDifferent wording can map to the same defect. An unmatched material claim is a false positive. A style preference with no planted mechanism maps to `NONE`.\n"
	if err := writeText(filepath.Join(blinded, "RUNSHEET.md"), runsheet); err != nil {
		return err
	}
	fmt.Printf("Created %d blinded items across %d output sets.\n", len(items), len(sets))
	return nil
}
