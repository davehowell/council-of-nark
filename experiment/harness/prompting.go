package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var defectIDPattern = regexp.MustCompile(`\b[A-Z]{2}-\d{2}\b`)
var markerPattern = regexp.MustCompile(`\{\{[A-Z_]+\}\}`)

func readText(root, relative string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
	return string(data), err
}
func packetText(root, packet string) (string, error) {
	return readText(root, "experiment/scenarios/"+packet+"/review-packet.md")
}
func answerText(root, packet string) (string, error) {
	return readText(root, "experiment/scenarios/"+packet+"/answer-key.md")
}

func loadSpecialists(root string) (map[string]map[string]string, error) {
	var rows []map[string]string
	if err := readJSON(filepath.Join(root, "experiment/prompts/specialists.json"), &rows); err != nil {
		return nil, err
	}
	out := map[string]map[string]string{}
	for _, row := range rows {
		out[row["id"]] = row
	}
	return out, nil
}
func replaceMarkers(template string, values map[string]string) (string, error) {
	out := template
	for key, value := range values {
		marker := "{{" + key + "}}"
		if !strings.Contains(out, marker) {
			return "", fmt.Errorf("template does not contain %s", marker)
		}
		out = strings.ReplaceAll(out, marker, value)
	}
	if leftovers := markerPattern.FindAllString(out, -1); len(leftovers) > 0 {
		return "", fmt.Errorf("unresolved template markers: %v", leftovers)
	}
	return out, nil
}
func assertAnswerKeyAbsent(prompt, answer string) error {
	if strings.Contains(prompt, "Private answer key") {
		return fmt.Errorf("answer-key marker leaked into reviewer prompt")
	}
	for _, id := range defectIDPattern.FindAllString(answer, -1) {
		if strings.Contains(prompt, id) {
			return fmt.Errorf("planted defect ID %s leaked into reviewer prompt", id)
		}
	}
	for _, line := range strings.Split(answer, "\n") {
		if !strings.HasPrefix(line, "|") || strings.HasPrefix(line, "|---") || strings.Contains(line, "Planted defect") {
			continue
		}
		parts := strings.Split(strings.Trim(line, "|"), "|")
		for i := 2; i < len(parts); i++ {
			cell := strings.TrimSpace(parts[i])
			if len(cell) >= 48 && strings.Contains(prompt, cell) {
				return fmt.Errorf("answer-key claim leaked verbatim into reviewer prompt")
			}
		}
	}
	return nil
}
func reviewContract(root, packet string) (string, error) {
	t, err := readText(root, "experiment/prompts/review-contract.txt")
	if err != nil {
		return "", err
	}
	p, err := packetText(root, packet)
	if err != nil {
		return "", err
	}
	return replaceMarkers(t, map[string]string{"REVIEW_PACKET": p})
}
func specialistPrompt(root, packet, role, wrapper string) (string, error) {
	rows, err := loadSpecialists(root)
	if err != nil {
		return "", err
	}
	row, ok := rows[role]
	if !ok {
		return "", fmt.Errorf("unknown specialist %s", role)
	}
	intro, err := readText(root, "experiment/prompts/specialist-intro.txt")
	if err != nil {
		return "", err
	}
	intro, err = replaceMarkers(intro, map[string]string{"LENS_KERNEL": row["kernel"]})
	if err != nil {
		return "", err
	}
	contract, err := reviewContract(root, packet)
	if err != nil {
		return "", err
	}
	prompt := row[wrapper+"_wrapper"] + "\n\n" + intro + "\n" + contract
	answer, err := answerText(root, packet)
	if err != nil {
		return "", err
	}
	if err := assertAnswerKeyAbsent(prompt, answer); err != nil {
		return "", err
	}
	return prompt, nil
}
func staticPrompt(root string, call Call) (string, error) {
	kind := stringValue(call.PromptSpec["kind"])
	var prompt string
	switch kind {
	case "generic":
		kernel, err := readText(root, "experiment/prompts/generic-kernel.txt")
		if err != nil {
			return "", err
		}
		contract, err := reviewContract(root, call.Packet)
		if err != nil {
			return "", err
		}
		prompt = kernel + "\n" + contract
	case "omnibus":
		var wrappers map[string]string
		if err := readJSON(filepath.Join(root, "experiment/prompts/omnibus-wrappers.json"), &wrappers); err != nil {
			return "", err
		}
		kernel, err := readText(root, "experiment/prompts/omnibus-kernel.txt")
		if err != nil {
			return "", err
		}
		contract, err := reviewContract(root, call.Packet)
		if err != nil {
			return "", err
		}
		prompt = wrappers[stringValue(call.PromptSpec["wrapper"])] + "\n\n" + kernel + "\n" + contract
	case "specialist":
		return specialistPrompt(root, call.Packet, stringValue(call.PromptSpec["role"]), stringValue(call.PromptSpec["wrapper"]))
	default:
		return "", fmt.Errorf("prompt kind %q requires dependency output", kind)
	}
	answer, err := answerText(root, call.Packet)
	if err != nil {
		return "", err
	}
	if err := assertAnswerKeyAbsent(prompt, answer); err != nil {
		return "", err
	}
	return prompt, nil
}
func withReviewerIDs(deps []dependencyPair) []map[string]any {
	payload := make([]map[string]any, 0, len(deps))
	for _, dep := range deps {
		findings := []any{}
		for _, raw := range sliceValue(dep.Response["findings"]) {
			finding := mapValue(raw)
			copyFinding := map[string]any{}
			for k, v := range finding {
				copyFinding[k] = v
			}
			if _, ok := copyFinding["raised_by"]; !ok {
				copyFinding["raised_by"] = []string{dep.Call.ReviewerID}
			}
			findings = append(findings, copyFinding)
		}
		payload = append(payload, map[string]any{"reviewer_id": dep.Call.ReviewerID, "findings": findings})
	}
	return payload
}

type dependencyPair struct {
	Call     Call
	Response map[string]any
}

func dynamicPrompt(root string, call Call, deps []dependencyPair) (string, error) {
	kind := stringValue(call.PromptSpec["kind"])
	packet, err := packetText(root, call.Packet)
	if err != nil {
		return "", err
	}
	ledger, _ := json.Marshal(withReviewerIDs(deps))
	var prompt string
	switch kind {
	case "fuser":
		t, err := readText(root, "experiment/prompts/fuser.txt")
		if err != nil {
			return "", err
		}
		prompt, err = replaceMarkers(t, map[string]string{"REVIEW_PACKET": packet, "REVIEW_FINDINGS": string(ledger)})
		if err != nil {
			return "", err
		}
	case "chain":
		rows, err := loadSpecialists(root)
		if err != nil {
			return "", err
		}
		row := rows[stringValue(call.PromptSpec["role"])]
		intro, err := readText(root, "experiment/prompts/specialist-intro.txt")
		if err != nil {
			return "", err
		}
		intro, err = replaceMarkers(intro, map[string]string{"LENS_KERNEL": row["kernel"]})
		if err != nil {
			return "", err
		}
		t, err := readText(root, "experiment/prompts/chain-contract.txt")
		if err != nil {
			return "", err
		}
		body, err := replaceMarkers(t, map[string]string{"REVIEW_PACKET": packet, "PRIOR_FINDINGS": string(ledger)})
		if err != nil {
			return "", err
		}
		prompt = row[stringValue(call.PromptSpec["wrapper"])+"_wrapper"] + "\n\n" + intro + "\n" + body
	default:
		return "", fmt.Errorf("prompt kind %q is not dynamic", kind)
	}
	answer, err := answerText(root, call.Packet)
	if err != nil {
		return "", err
	}
	if err := assertAnswerKeyAbsent(prompt, answer); err != nil {
		return "", err
	}
	return prompt, nil
}
