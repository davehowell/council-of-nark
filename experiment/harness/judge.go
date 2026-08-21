package harness

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ratingFields = []string{"rater", "item_id", "defect_id", "false_positive_cluster", "material", "confidence", "notes"}

func readJSONL(path string) ([]map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rows := []map[string]any{}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row map[string]any
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
func readCSVMaps(path string) (map[string]map[string]string, error) {
	out := map[string]map[string]string{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return out, nil
	}
	header := records[0]
	for _, values := range records[1:] {
		row := map[string]string{}
		for i, key := range header {
			if i < len(values) {
				row[key] = values[i]
			}
		}
		id := row["item_id"]
		if out[id] != nil {
			return nil, fmt.Errorf("duplicate existing LLM rating: %s", id)
		}
		out[id] = row
	}
	return out, nil
}
func recoverCaptured(destination string, expected map[string]bool) (map[string]any, error) {
	paths, err := filepath.Glob(filepath.Join(destination, "attempts", "*", "stdout.txt"))
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		parsed, _ := extract(string(data), "judgements")
		if parsed != nil && len(judgementErrors(parsed, expected)) == 0 {
			return parsed, nil
		}
	}
	return nil, nil
}
func mergeRatings(ratings map[string]map[string]string, response map[string]any, model string) {
	for _, raw := range sliceValue(response["judgements"]) {
		row := mapValue(raw)
		defect, cluster := "", ""
		if row["defect_id"] != nil {
			defect = stringValue(row["defect_id"])
		}
		if row["false_positive_cluster"] != nil {
			cluster = stringValue(row["false_positive_cluster"])
		}
		id := stringValue(row["item_id"])
		ratings[id] = map[string]string{"rater": "llm:" + model, "item_id": id, "defect_id": defect, "false_positive_cluster": cluster, "material": strconv.FormatBool(boolValue(row["material"])), "confidence": stringValue(row["confidence"]), "notes": stringValue(row["rationale"])}
	}
}
func (h *Harness) rateSet(run string, set OutputSet, items []map[string]any, freeze Freeze, config Config, judgeCommit string) (map[string]any, error) {
	if len(items) == 0 {
		return map[string]any{"judgements": []any{}}, nil
	}
	judgeID := "j-" + opaqueID(config.Label, 16, set.SetID)
	destination := filepath.Join(run, "blinded", "llm-triage", set.SetID)
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return nil, err
	}
	var previous any
	if _, err := os.Stat(filepath.Join(destination, "record.json")); err == nil {
		var p map[string]any
		if err := readJSON(filepath.Join(destination, "record.json"), &p); err == nil {
			previous = p
		}
	}
	var response map[string]any
	err := h.withWorktree(run, judgeID, freeze.SourceCommit, func(worktree string) error {
		template, err := readText(h.Root, "experiment/prompts/judge.txt")
		if err != nil {
			return err
		}
		answer, err := readText(worktree, "experiment/scenarios/"+set.Packet+"/answer-key.md")
		if err != nil {
			return err
		}
		publicItems := []map[string]any{}
		expected := map[string]bool{}
		for _, item := range items {
			id := stringValue(item["item_id"])
			expected[id] = true
			publicItems = append(publicItems, map[string]any{"item_id": id, "finding": item["finding"]})
		}
		itemJSON, _ := json.Marshal(publicItems)
		prompt, err := replaceMarkers(template, map[string]string{"ANSWER_KEY": answer, "FINDINGS": string(itemJSON)})
		if err != nil {
			return err
		}
		system := "You are an isolated blinded evaluation rater. Use only the supplied answer key and findings. Do not use tools, memories, project context, or other sessions. Return only schema-valid JSON."
		var schema map[string]any
		if err := readJSON(filepath.Join(h.Root, "experiment/schema/judgements.schema.json"), &schema); err != nil {
			return err
		}
		maxAttempts := config.MaxAttempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}
		attempts := []map[string]any{}
		batch := slugTimestamp()
		var result Invocation
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			privateBase := filepath.Join(h.Experiment, "private", "sandboxes", filepath.Base(run), judgeID, fmt.Sprintf("%s-%d", batch, attempt))
			_ = os.RemoveAll(privateBase)
			result, err = invoke(privateBase, config.Provider, prompt, system, schema, config.RequestTimeoutSeconds, "judgements", expected)
			_ = os.RemoveAll(privateBase)
			if err != nil {
				return err
			}
			attemptDir := filepath.Join(destination, "attempts", fmt.Sprintf("%s-%d", batch, attempt))
			if err := os.MkdirAll(attemptDir, 0o755); err != nil {
				return err
			}
			if err := writeText(filepath.Join(attemptDir, "stdout.txt"), result.Stdout); err != nil {
				return err
			}
			if err := writeText(filepath.Join(attemptDir, "stderr.txt"), result.Stderr); err != nil {
				return err
			}
			meta := map[string]any{"attempt": attempt, "returncode": result.ReturnCode, "parse_error": nil, "latency_seconds": result.LatencySeconds, "usage": result.Usage, "sandbox": result.Sandbox}
			if result.ParseError != "" {
				meta["parse_error"] = result.ParseError
			}
			if err := writeJSON(filepath.Join(attemptDir, "metadata.json"), meta); err != nil {
				return err
			}
			attempts = append(attempts, meta)
			if result.ReturnCode == 0 {
				break
			}
			if attempt < maxAttempts {
				backoff := config.RetryBackoffSeconds
				if backoff == 0 {
					backoff = 5
				}
				time.Sleep(time.Duration(backoff * float64(attempt) * float64(time.Second)))
			}
		}
		status := "error"
		if result.ReturnCode == 0 && result.Parsed != nil {
			status = "success"
			response = result.Parsed
		}
		record := map[string]any{"set_id": set.SetID, "packet": set.Packet, "provider": providerMap(config.Provider), "status": status, "attempts": attempts, "parsed_response": response, "completed_at": utcNow(), "arm_blinded": true, "status_note": config.Status, "review_source_commit": freeze.SourceCommit, "judge_harness_commit": judgeCommit, "previous_record": previous}
		return writeJSON(filepath.Join(destination, "record.json"), record)
	})
	return response, err
}

func writeRatings(path string, ratings map[string]map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(f)
	if err := writer.Write(ratingFields); err != nil {
		return err
	}
	ids := make([]string, 0, len(ratings))
	for id := range ratings {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		row := ratings[id]
		values := make([]string, len(ratingFields))
		for i, key := range ratingFields {
			values[i] = row[key]
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return f.Close()
}
func blindedSetIDs(run string, plan Plan) map[string]string {
	mapping := map[string]string{}
	var unblind map[string]any
	if err := readJSON(filepath.Join(run, "private", "unblind.json"), &unblind); err == nil {
		for blindID, raw := range mapValue(unblind["sets"]) {
			actual := stringValue(mapValue(raw)["set_id"])
			if actual != "" {
				mapping[actual] = blindID
			}
		}
	}
	for _, set := range plan.OutputSets {
		if mapping[set.SetID] == "" {
			mapping[set.SetID] = set.SetID
		}
	}
	return mapping
}

func (h *Harness) Judge(runArgument, configArgument string, jobs int) error {
	if err := h.requireCleanTree(); err != nil {
		return fmt.Errorf("derived rating stages require a clean committed harness: %w", err)
	}
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(run, "seal.json")); err != nil {
		return fmt.Errorf("seal raw model output before rating")
	}
	configPath := configArgument
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(h.Root, configPath)
	}
	var config Config
	var freeze Freeze
	var plan Plan
	if err := readJSON(configPath, &config); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(run, "freeze.json"), &freeze); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	items, err := readJSONL(filepath.Join(run, "blinded", "findings.jsonl"))
	if err != nil {
		return err
	}
	bySet := map[string][]map[string]any{}
	known := map[string]bool{}
	for _, item := range items {
		id := stringValue(item["item_id"])
		known[id] = true
		setID := stringValue(item["set_id"])
		bySet[setID] = append(bySet[setID], item)
	}
	ratingsPath := filepath.Join(run, "blinded", "ratings-llm.csv")
	ratings, err := readCSVMaps(ratingsPath)
	if err != nil {
		return err
	}
	for id := range ratings {
		if !known[id] {
			return fmt.Errorf("existing ratings contain unknown item ID %s", id)
		}
	}
	pending := []OutputSet{}
	blindSetByActual := blindedSetIDs(run, plan)
	judgeCommit, err := h.sourceCommit()
	if err != nil {
		return err
	}
	for _, sourceSet := range plan.OutputSets {
		set := sourceSet
		set.SetID = blindSetByActual[sourceSet.SetID]
		expected := map[string]bool{}
		complete := true
		for _, item := range bySet[set.SetID] {
			id := stringValue(item["item_id"])
			expected[id] = true
			if ratings[id] == nil {
				complete = false
			}
		}
		if len(expected) == 0 || complete {
			continue
		}
		for id := range expected {
			delete(ratings, id)
		}
		destination := filepath.Join(run, "blinded", "llm-triage", set.SetID)
		recovered, err := recoverCaptured(destination, expected)
		if err != nil {
			return err
		}
		if recovered != nil {
			mergeRatings(ratings, recovered, config.Provider.Model)
			recordPath := filepath.Join(destination, "record.json")
			var record map[string]any
			if err := readJSON(recordPath, &record); err == nil {
				record["status"] = "success"
				record["parsed_response"] = recovered
				record["recovered_at"] = utcNow()
				record["recovery"] = "parser-only recovery from captured structured output"
				record["judge_harness_commit"] = judgeCommit
				_ = writeJSON(recordPath, record)
			}
			fmt.Printf("%s recovered from captured output\n", set.SetID)
		} else {
			pending = append(pending, set)
		}
	}
	fmt.Printf("Resuming with %d existing/recovered ratings; %d sets pending.\n", len(ratings), len(pending))
	if jobs < 1 {
		jobs = 2
	}
	sem := make(chan struct{}, jobs)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errorsFound := []string{}
	for _, set := range pending {
		set := set
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			response, err := h.rateSet(run, set, bySet[set.SetID], freeze, config, judgeCommit)
			<-sem
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil || response == nil {
				errorsFound = append(errorsFound, set.SetID)
				fmt.Printf("%s rating error\n", set.SetID)
				return
			}
			mergeRatings(ratings, response, config.Provider.Model)
			fmt.Printf("%s rated\n", set.SetID)
		}()
	}
	wg.Wait()
	if err := writeRatings(ratingsPath, ratings); err != nil {
		return err
	}
	missing := 0
	for id := range known {
		if ratings[id] == nil {
			missing++
		}
	}
	rel, _ := filepath.Rel(run, ratingsPath)
	fmt.Printf("Wrote %d exploratory ratings to %s\n", len(ratings), filepath.ToSlash(rel))
	if len(errorsFound) > 0 || missing > 0 {
		return fmt.Errorf("unrated sets=%v; unrated items=%d", errorsFound, missing)
	}
	return nil
}
