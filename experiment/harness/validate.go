package harness

import (
	"fmt"
	"strings"
)

var confidenceValues = map[string]bool{"high": true, "medium": true, "low": true}
var requiredFinding = []string{"location", "claim", "consequence", "fix", "confidence"}

func findingsErrors(value map[string]any) []string {
	errors := []string{}
	if len(value) != 1 {
		errors = append(errors, "root must contain only 'findings'")
	}
	raw, ok := value["findings"]
	if !ok {
		return append(errors, "findings is not an array")
	}
	findings, ok := raw.([]any)
	if !ok {
		return append(errors, "findings is not an array")
	}
	if len(findings) > 8 {
		errors = append(errors, "findings contains more than 8 items")
	}
	allowed := map[string]bool{"raised_by": true}
	for _, k := range requiredFinding {
		allowed[k] = true
	}
	for i, rawFinding := range findings {
		prefix := fmt.Sprintf("findings[%d]", i)
		finding, ok := rawFinding.(map[string]any)
		if !ok {
			errors = append(errors, prefix+" is not an object")
			continue
		}
		for _, key := range requiredFinding {
			if _, ok := finding[key]; !ok {
				errors = append(errors, prefix+" lacks required field "+key)
			}
		}
		for key := range finding {
			if !allowed[key] {
				errors = append(errors, prefix+" has unknown field "+key)
			}
		}
		for _, key := range []string{"location", "claim", "consequence", "fix"} {
			if strings.TrimSpace(stringValue(finding[key])) == "" {
				errors = append(errors, prefix+"."+key+" must be a non-empty string")
			}
		}
		if !confidenceValues[stringValue(finding["confidence"])] {
			errors = append(errors, prefix+".confidence is invalid")
		}
		if raised, ok := finding["raised_by"]; ok {
			seen := map[string]bool{}
			values, valid := raised.([]any)
			if !valid {
				errors = append(errors, prefix+".raised_by is invalid")
			} else {
				for _, v := range values {
					s := stringValue(v)
					if s == "" || seen[s] {
						valid = false
					}
					seen[s] = true
				}
				if !valid {
					errors = append(errors, prefix+".raised_by is invalid")
				}
			}
		}
	}
	return errors
}
func judgementErrors(value map[string]any, expected map[string]bool) []string {
	if len(value) != 1 {
		return []string{"root must contain only 'judgements'"}
	}
	raw, ok := value["judgements"]
	if !ok {
		return []string{"root must contain only 'judgements'"}
	}
	rows, ok := raw.([]any)
	if !ok {
		return []string{"judgements is not an array"}
	}
	errors := []string{}
	seen := map[string]bool{}
	required := map[string]bool{"item_id": true, "defect_id": true, "false_positive_cluster": true, "material": true, "confidence": true, "rationale": true}
	for i, rawRow := range rows {
		prefix := fmt.Sprintf("judgements[%d]", i)
		row, ok := rawRow.(map[string]any)
		if !ok {
			errors = append(errors, prefix+" is not an object")
			continue
		}
		if len(row) != len(required) {
			errors = append(errors, prefix+" fields do not match schema")
			continue
		}
		validFields := true
		for k := range row {
			if !required[k] {
				validFields = false
			}
		}
		if !validFields {
			errors = append(errors, prefix+" fields do not match schema")
			continue
		}
		id := stringValue(row["item_id"])
		if !expected[id] || seen[id] {
			errors = append(errors, prefix+" has unexpected or duplicate item_id")
		}
		seen[id] = true
		defect := row["defect_id"]
		cluster := row["false_positive_cluster"]
		if defect != nil {
			if _, ok := defect.(string); !ok {
				errors = append(errors, prefix+".defect_id is invalid")
			}
		}
		if cluster != nil {
			if s, ok := cluster.(string); !ok || strings.TrimSpace(s) == "" {
				errors = append(errors, prefix+".false_positive_cluster is invalid")
			}
		}
		if defect == nil && cluster == nil {
			errors = append(errors, prefix+" needs a false-positive cluster")
		}
		if defect != nil && cluster != nil {
			errors = append(errors, prefix+" cannot have both defect ID and false-positive cluster")
		}
		if _, ok := row["material"].(bool); !ok {
			errors = append(errors, prefix+".material is invalid")
		}
		if !confidenceValues[stringValue(row["confidence"])] {
			errors = append(errors, prefix+".confidence is invalid")
		}
		if _, ok := row["rationale"].(string); !ok {
			errors = append(errors, prefix+".rationale is invalid")
		}
	}
	if len(seen) != len(expected) {
		errors = append(errors, "judgements do not cover every expected item ID")
	} else {
		for id := range expected {
			if !seen[id] {
				errors = append(errors, "judgements do not cover every expected item ID")
				break
			}
		}
	}
	return errors
}
