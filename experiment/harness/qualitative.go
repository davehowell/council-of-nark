package harness

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func scoreOneToFive(row map[string]string, key string) (float64, error) {
	value, err := strconv.Atoi(strings.TrimSpace(row[key]))
	if err != nil || value < 1 || value > 5 {
		return 0, fmt.Errorf("%s must be an integer from 1 to 5", key)
	}
	return float64(value), nil
}

func (h *Harness) UnblindQualitative(runArgument, ratingsArgument, label string) error {
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	ratingsPath := ratingsArgument
	if !filepath.IsAbs(ratingsPath) {
		ratingsPath = filepath.Join(run, ratingsPath)
	}
	f, err := os.Open(ratingsPath)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	_ = f.Close()
	if err != nil {
		return err
	}
	if len(records) < 2 {
		return fmt.Errorf("qualitative ratings contain no rows")
	}
	header := records[0]
	ratings := []map[string]string{}
	seen := map[string]bool{}
	for _, values := range records[1:] {
		row := map[string]string{}
		for i, key := range header {
			if i < len(values) {
				row[key] = values[i]
			}
		}
		pairID, rater := row["pair_id"], strings.TrimSpace(row["rater"])
		if pairID == "" || rater == "" {
			return fmt.Errorf("every qualitative row needs rater and pair_id")
		}
		key := rater + "\x1f" + pairID
		if seen[key] {
			return fmt.Errorf("duplicate qualitative rating for %s by %s", pairID, rater)
		}
		seen[key] = true
		ratings = append(ratings, row)
	}
	var unblind map[string]any
	if err := readJSON(filepath.Join(run, "private", "unblind.json"), &unblind); err != nil {
		return err
	}
	pairMap := mapValue(unblind["pairs"])
	if len(pairMap) == 0 {
		return fmt.Errorf("run has no private paired-comparison map")
	}
	pairRows, err := readJSONL(filepath.Join(run, "blinded", "pairs.jsonl"))
	if err != nil {
		return err
	}
	publicPairs := map[string]map[string]any{}
	for _, row := range pairRows {
		publicPairs[stringValue(row["pair_id"])] = row
	}
	output := []map[string]any{}
	scoreSums := map[string]map[string]float64{"functional": {}, "fictional": {}}
	scoreCounts := map[string]map[string]int{"functional": {}, "fictional": {}}
	preferences := map[string]int{"functional": 0, "fictional": 0, "tie": 0}
	knownGuesses, correctGuesses := 0, 0
	wording := map[string]int{"true": 0, "false": 0, "unsure": 0}
	for _, rating := range ratings {
		pairID := rating["pair_id"]
		mapping := mapValue(pairMap[pairID])
		if len(mapping) == 0 {
			return fmt.Errorf("unknown pair_id %s", pairID)
		}
		leftCondition, rightCondition := stringValue(mapping["left_condition"]), stringValue(mapping["right_condition"])
		public := publicPairs[pairID]
		row := map[string]any{}
		for key, value := range rating {
			row[key] = value
		}
		row["packet"] = public["packet"]
		row["role"] = public["role"]
		row["kind"] = public["kind"]
		row["left_condition"] = leftCondition
		row["right_condition"] = rightCondition
		row["left_source_set_id"] = mapping["left_set_id"]
		row["right_source_set_id"] = mapping["right_set_id"]
		output = append(output, row)
		for _, metric := range []string{"supportedness", "actionability", "fix_quality"} {
			left, err := scoreOneToFive(rating, "left_"+metric+"_1_5")
			if err != nil {
				return fmt.Errorf("%s: %w", pairID, err)
			}
			right, err := scoreOneToFive(rating, "right_"+metric+"_1_5")
			if err != nil {
				return fmt.Errorf("%s: %w", pairID, err)
			}
			scoreSums[leftCondition][metric] += left
			scoreCounts[leftCondition][metric]++
			scoreSums[rightCondition][metric] += right
			scoreCounts[rightCondition][metric]++
		}
		switch rating["overall_preference_left_right_tie"] {
		case "left":
			preferences[leftCondition]++
		case "right":
			preferences[rightCondition]++
		case "tie":
			preferences["tie"]++
		default:
			return fmt.Errorf("%s has invalid overall preference", pairID)
		}
		guess := rating["condition_guess_left_functional_left_fictional_unsure"]
		if guess != "unsure" {
			knownGuesses++
			if (guess == "left-functional" && leftCondition == "functional") || (guess == "left-fictional" && leftCondition == "fictional") {
				correctGuesses++
			}
		}
		reveal := rating["wording_revealed_condition_true_false_unsure"]
		if _, ok := wording[reveal]; !ok {
			return fmt.Errorf("%s has invalid wording reveal value", pairID)
		}
		wording[reveal]++
	}
	means := map[string]any{}
	for _, condition := range []string{"functional", "fictional"} {
		metrics := map[string]any{}
		for _, metric := range []string{"supportedness", "actionability", "fix_quality"} {
			count := scoreCounts[condition][metric]
			if count > 0 {
				metrics["mean_"+metric] = scoreSums[condition][metric] / float64(count)
			}
		}
		means[condition] = metrics
	}
	guessAccuracy := any(nil)
	if knownGuesses > 0 {
		guessAccuracy = float64(correctGuesses) / float64(knownGuesses)
	}
	summary := map[string]any{"schema_version": 1, "n_ratings": len(ratings), "condition_scores": means, "preferences": preferences, "condition_guess": map[string]any{"known_guesses": knownGuesses, "correct_guesses": correctGuesses, "accuracy_when_guessed": guessAccuracy}, "wording_revealed_condition": wording, "warning": "condition labels were hidden, but output wording may reveal treatment; report guess/reveal results with quality ratings"}
	analysis := filepath.Join(run, "analysis", label)
	if err := os.MkdirAll(analysis, 0o755); err != nil {
		return err
	}
	if err := writeRowsCSV(filepath.Join(analysis, "qualitative-pairs-unblinded.csv"), output); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(analysis, "qualitative-summary.json"), summary); err != nil {
		return err
	}
	fmt.Printf("Unblinded %d locked qualitative ratings into %s while preserving opaque pair IDs.\n", len(ratings), h.relative(analysis))
	return nil
}
