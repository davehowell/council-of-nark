package harness

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type stringSet map[string]bool

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
func percentile(values []float64, q float64) any {
	if len(values) == 0 {
		return nil
	}
	ordered := append([]float64{}, values...)
	sort.Float64s(ordered)
	index := int(math.Ceil(q*float64(len(ordered)))) - 1
	if index < 0 {
		index = 0
	}
	return ordered[index]
}
func bootstrapMeanCI(values []float64, seed string, samples int) any {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return []float64{values[0], values[0]}
	}
	means := make([]float64, samples)
	for sample := 0; sample < samples; sample++ {
		sum := 0.0
		for draw := range values {
			digest := sha256.Sum256([]byte(fmt.Sprintf("%s\x1f%d\x1f%d", seed, sample, draw)))
			index := int(binary.BigEndian.Uint64(digest[:8]) % uint64(len(values)))
			sum += values[index]
		}
		means[sample] = sum / float64(len(values))
	}
	sort.Float64s(means)
	return []float64{means[int(.025*float64(samples))], means[int(.975*float64(samples))-1]}
}
func answerIDs(root, packet string) (stringSet, error) {
	text, err := readText(root, "experiment/scenarios/"+packet+"/answer-key.md")
	if err != nil {
		return nil, err
	}
	out := stringSet{}
	for _, line := range strings.Split(text, "\n") {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		parts := strings.Split(strings.Trim(line, "|"), "|")
		if len(parts) == 0 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		if len(id) == 5 && id[2] == '-' && id[3] >= '0' && id[3] <= '9' && id[4] >= '0' && id[4] <= '9' {
			out[id] = true
		}
	}
	return out, nil
}
func setUsage(run string, ids []string) (int, int, float64, float64) {
	input, output := 0, 0
	cost, latency := 0.0, 0.0
	for _, id := range ids {
		record, _ := readCallRecord(run, id)
		if record == nil {
			continue
		}
		in, out, c := tokenCounts(mapValue(record["usage"]))
		input += in
		output += out
		cost += c
		latency += floatValue(record["latency_seconds"])
	}
	return input, output, cost, latency
}
func groupName(metadata map[string]any) string {
	switch stringValue(metadata["design"]) {
	case "stage_a":
		return stringValue(metadata["arm"]) + ":" + stringValue(metadata["kind"])
	case "topology":
		return stringValue(metadata["topology"]) + ":" + stringValue(metadata["kind"])
	default:
		return fmt.Sprintf("provider-%d:%s:%s:%s", intValue(metadata["provider_index"]), stringValue(metadata["adapter"]), stringValue(metadata["model"]), stringValue(metadata["wrapper"]))
	}
}
func setDifference(a, b stringSet) []string {
	out := []string{}
	for value := range a {
		if !b[value] {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}
func setJaccard(a, b stringSet) float64 {
	union, intersection := stringSet{}, 0
	for value := range a {
		union[value] = true
		if b[value] {
			intersection++
		}
	}
	for value := range b {
		union[value] = true
	}
	if len(union) == 0 {
		return 1
	}
	return float64(intersection) / float64(len(union))
}
func readRatings(path string) (map[string]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	all, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("empty ratings CSV")
	}
	headers := all[0]
	out := map[string]map[string]string{}
	for _, values := range all[1:] {
		row := map[string]string{}
		for i, key := range headers {
			if i < len(values) {
				row[key] = values[i]
			}
		}
		id := row["item_id"]
		if out[id] != nil {
			return nil, fmt.Errorf("duplicate rating for %s; adjudicate to one row per item", id)
		}
		out[id] = row
	}
	return out, nil
}
func csvValue(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.15g", v)
	case int:
		return fmt.Sprint(v)
	case bool:
		return fmt.Sprint(v)
	default:
		data, _ := json.Marshal(v)
		return string(data)
	}
}
func writeRowsCSV(path string, rows []map[string]any) error {
	fieldsSet := map[string]bool{}
	for _, row := range rows {
		for key := range row {
			fieldsSet[key] = true
		}
	}
	fields := make([]string, 0, len(fieldsSet))
	for key := range fieldsSet {
		fields = append(fields, key)
	}
	sort.Strings(fields)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	_ = writer.Write(fields)
	for _, row := range rows {
		values := make([]string, len(fields))
		for i, key := range fields {
			values[i] = csvValue(row[key])
		}
		_ = writer.Write(values)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
func pairKey(parts ...any) string {
	values := make([]string, len(parts))
	for i, p := range parts {
		values[i] = fmt.Sprint(p)
	}
	return strings.Join(values, "\x1f")
}

func (h *Harness) Score(runArgument, ratingsArgument, label string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	ratingsPath := ratingsArgument
	if !filepath.IsAbs(ratingsPath) {
		ratingsPath = filepath.Join(run, ratingsPath)
	}
	var plan Plan
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	itemRows, err := readJSONL(filepath.Join(run, "blinded", "findings.jsonl"))
	if err != nil {
		return err
	}
	setRows, err := readJSONL(filepath.Join(run, "blinded", "sets.jsonl"))
	if err != nil {
		return err
	}
	items := map[string]map[string]any{}
	for _, row := range itemRows {
		items[stringValue(row["item_id"])] = row
	}
	blindedSets := map[string]map[string]any{}
	for _, row := range setRows {
		blindedSets[stringValue(row["set_id"])] = row
	}
	ratings, err := readRatings(ratingsPath)
	if err != nil {
		return err
	}
	missing, extra := []string{}, []string{}
	for id := range items {
		if ratings[id] == nil {
			missing = append(missing, id)
		}
	}
	for id := range ratings {
		if items[id] == nil {
			extra = append(extra, id)
		}
	}
	if len(missing) > 0 || len(extra) > 0 {
		return fmt.Errorf("rating coverage mismatch; missing=%v, extra=%v", missing, extra)
	}
	planSets := map[string]OutputSet{}
	for _, set := range plan.OutputSets {
		planSets[set.SetID] = set
	}
	rows := []map[string]any{}
	detectedBySet := map[string]stringSet{}
	for setID, blinded := range blindedSets {
		set := planSets[setID]
		valid, err := answerIDs(h.Root, set.Packet)
		if err != nil {
			return err
		}
		detected, falseClusters := stringSet{}, stringSet{}
		for _, id := range stringSlice(blinded["item_ids"]) {
			rating := ratings[id]
			defect := strings.ToUpper(strings.TrimSpace(rating["defect_id"]))
			material := strings.ToLower(strings.TrimSpace(rating["material"]))
			if material != "true" && material != "false" {
				return fmt.Errorf("invalid material value for %s", id)
			}
			if defect == "" || defect == "NONE" || defect == "NULL" || material == "false" {
				cluster := strings.Join(strings.Fields(strings.ToLower(rating["false_positive_cluster"])), " ")
				if cluster == "" {
					cluster = "legacy-item:" + id
				}
				falseClusters[cluster] = true
			} else if !valid[defect] {
				return fmt.Errorf("invalid defect ID %q for packet %s", defect, set.Packet)
			} else {
				detected[defect] = true
			}
		}
		tp, fp := len(detected), len(falseClusters)
		predicted := tp + fp
		precision, recall := 0.0, 0.0
		if predicted > 0 {
			precision = float64(tp) / float64(predicted)
		}
		if len(valid) > 0 {
			recall = float64(tp) / float64(len(valid))
		}
		f1 := 0.0
		if precision+recall > 0 {
			f1 = 2 * precision * recall / (precision + recall)
		}
		input, output, cost, latency := setUsage(run, set.CostCallIDs)
		row := map[string]any{"set_id": setID, "packet": set.Packet, "group": groupName(set.Metadata), "true_positives": tp, "detected_defect_ids": strings.Join(setDifference(detected, stringSet{}), "|"), "false_positives": fp, "false_positive_clusters": strings.Join(setDifference(falseClusters, stringSet{}), "|"), "possible_defects": len(valid), "precision": precision, "recall": recall, "f1": f1, "input_tokens": input, "output_tokens": output, "cost_usd": cost, "latency_seconds_serial_sum": latency}
		if output > 0 {
			row["true_findings_per_1k_output_tokens"] = 1000 * float64(tp) / float64(output)
		} else {
			row["true_findings_per_1k_output_tokens"] = nil
		}
		if cost > 0 {
			row["true_findings_per_dollar"] = float64(tp) / cost
		} else {
			row["true_findings_per_dollar"] = nil
		}
		for key, value := range set.Metadata {
			row[key] = value
		}
		rows = append(rows, row)
		detectedBySet[setID] = detected
	}
	sort.Slice(rows, func(i, j int) bool { return stringValue(rows[i]["set_id"]) < stringValue(rows[j]["set_id"]) })
	grouped := map[string][]map[string]any{}
	for _, row := range rows {
		name := stringValue(row["group"])
		grouped[name] = append(grouped[name], row)
	}
	summary := map[string]any{}
	for name, group := range grouped {
		f1s, precisions, recalls, outputs := []float64{}, []float64{}, []float64{}, []float64{}
		cost := 0.0
		for _, row := range group {
			f1s = append(f1s, floatValue(row["f1"]))
			precisions = append(precisions, floatValue(row["precision"]))
			recalls = append(recalls, floatValue(row["recall"]))
			outputs = append(outputs, floatValue(row["output_tokens"]))
			cost += floatValue(row["cost_usd"])
		}
		worst := 0.0
		if len(f1s) > 0 {
			worst = f1s[0]
			for _, v := range f1s {
				if v < worst {
					worst = v
				}
			}
		}
		summary[name] = map[string]any{"n_sets": len(group), "mean_f1": mean(f1s), "p10_f1": percentile(f1s, .10), "worst_f1": worst, "mean_precision": mean(precisions), "mean_recall": mean(recalls), "mean_output_tokens": mean(outputs), "total_cost_usd": cost}
	}
	fusion := []map[string]any{}
	fusionIndex := map[string]map[string]map[string]any{}
	for _, row := range rows {
		arm := stringValue(row["arm"])
		if stringValue(row["design"]) != "stage_a" || (arm != "M0" && arm != "M1" && arm != "M2") {
			continue
		}
		key := pairKey(row["packet"], arm, row["repeat"])
		if fusionIndex[key] == nil {
			fusionIndex[key] = map[string]map[string]any{}
		}
		fusionIndex[key][stringValue(row["kind"])] = row
	}
	for _, pair := range fusionIndex {
		raw, fused := pair["raw_union"], pair["fused"]
		if raw == nil || fused == nil {
			continue
		}
		retention := any(nil)
		if intValue(raw["true_positives"]) > 0 {
			retention = float64(intValue(fused["true_positives"])) / float64(intValue(raw["true_positives"]))
		}
		fusion = append(fusion, map[string]any{"packet": raw["packet"], "arm": raw["arm"], "repeat": raw["repeat"], "raw_true_positives": raw["true_positives"], "fused_true_positives": fused["true_positives"], "fusion_retention": retention, "raw_false_positives": raw["false_positives"], "fused_false_positives": fused["false_positives"]})
	}
	stageIndex := map[string]map[string]any{}
	for _, row := range rows {
		kind := stringValue(row["kind"])
		if stringValue(row["design"]) == "stage_a" && (kind == "final" || kind == "fused") {
			stageIndex[pairKey(row["packet"], row["repeat"], row["arm"])] = row
		}
	}
	overlaps := []map[string]any{}
	packets, repeats := stringSet{}, map[string]map[int]bool{}
	for _, row := range stageIndex {
		packet := stringValue(row["packet"])
		packets[packet] = true
		if repeats[packet] == nil {
			repeats[packet] = map[int]bool{}
		}
		repeats[packet][intValue(row["repeat"])] = true
	}
	for packet := range packets {
		for repeat := range repeats[packet] {
			for _, pair := range [][2]string{{"S1", "S2"}, {"M0", "M1"}, {"M1", "M2"}} {
				left, right := stageIndex[pairKey(packet, repeat, pair[0])], stageIndex[pairKey(packet, repeat, pair[1])]
				if left == nil || right == nil {
					continue
				}
				a, b := detectedBySet[stringValue(left["set_id"])], detectedBySet[stringValue(right["set_id"])]
				overlaps = append(overlaps, map[string]any{"packet": packet, "repeat": repeat, "left": pair[0], "right": pair[1], "jaccard": setJaccard(a, b), "left_only": setDifference(a, b), "right_only": setDifference(b, a)})
			}
		}
	}
	pairIndex := map[string]map[string]map[string]any{}
	for _, row := range rows {
		if stringValue(row["design"]) != "provider_pair" {
			continue
		}
		key := pairKey(row["packet"], row["repeat"], row["provider_index"])
		if pairIndex[key] == nil {
			pairIndex[key] = map[string]map[string]any{}
		}
		pairIndex[key][stringValue(row["wrapper"])] = row
	}
	pairedRows := []map[string]any{}
	for _, pair := range pairIndex {
		functional, fictional := pair["functional"], pair["fictional"]
		if functional == nil || fictional == nil {
			continue
		}
		a, b := detectedBySet[stringValue(functional["set_id"])], detectedBySet[stringValue(fictional["set_id"])]
		pairedRows = append(pairedRows, map[string]any{"packet": functional["packet"], "repeat": functional["repeat"], "provider_index": functional["provider_index"], "f1_delta_fictional_minus_functional": floatValue(fictional["f1"]) - floatValue(functional["f1"]), "precision_delta": floatValue(fictional["precision"]) - floatValue(functional["precision"]), "recall_delta": floatValue(fictional["recall"]) - floatValue(functional["recall"]), "jaccard": setJaccard(a, b), "functional_only": setDifference(a, b), "fictional_only": setDifference(b, a)})
	}
	pairedSummary := map[string]any{}
	if len(pairedRows) > 0 {
		packetNames := []string{"all"}
		seen := stringSet{}
		for _, row := range pairedRows {
			seen[stringValue(row["packet"])] = true
		}
		for packet := range seen {
			packetNames = append(packetNames, packet)
		}
		sort.Strings(packetNames[1:])
		for _, packet := range packetNames {
			selected := []map[string]any{}
			for _, row := range pairedRows {
				if packet == "all" || stringValue(row["packet"]) == packet {
					selected = append(selected, row)
				}
			}
			deltas, jaccards := []float64{}, []float64{}
			fictionalWins, ties, functionalWins := 0, 0, 0
			for _, row := range selected {
				delta := floatValue(row["f1_delta_fictional_minus_functional"])
				deltas = append(deltas, delta)
				jaccards = append(jaccards, floatValue(row["jaccard"]))
				if delta > 1e-12 {
					fictionalWins++
				} else if delta < -1e-12 {
					functionalWins++
				} else {
					ties++
				}
			}
			pairedSummary[packet] = map[string]any{"n_pairs": len(selected), "mean_f1_delta_fictional_minus_functional": mean(deltas), "bootstrap_95_ci_sampling_only": bootstrapMeanCI(deltas, plan.Seed+":"+packet+":paired-f1", 10000), "fictional_wins": fictionalWins, "ties": ties, "functional_wins": functionalWins, "mean_jaccard": mean(jaccards)}
		}
	}
	analysis := filepath.Join(run, "analysis", label)
	if err := os.MkdirAll(analysis, 0o755); err != nil {
		return err
	}
	report := map[string]any{"groups": summary, "fusion": fusion, "overlaps": overlaps, "paired_persona": map[string]any{"summary": pairedSummary, "pairs": pairedRows}}
	if err := writeJSON(filepath.Join(analysis, "summary.json"), report); err != nil {
		return err
	}
	if err := writeRowsCSV(filepath.Join(analysis, "sets.csv"), rows); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
	return nil
}
