package respondent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
)

type EcologicalReview struct {
	ReviewSummary string              `json:"review_summary"`
	Findings      []EcologicalFinding `json:"findings"`
	Uncertainties []string            `json:"uncertainties"`
}

type EcologicalFinding struct {
	Title              string           `json:"title"`
	Locations          []SourceLocation `json:"locations"`
	Claim              string           `json:"claim"`
	Mechanism          string           `json:"mechanism"`
	Consequence        string           `json:"consequence"`
	Correction         string           `json:"correction"`
	AllocationAndScope string           `json:"allocation_and_scope"`
	RegressionTests    []RegressionCase `json:"regression_tests"`
	Evidence           []string         `json:"evidence"`
	Confidence         string           `json:"confidence"`
}

type SourceLocation struct {
	Path      string `json:"path"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
}

type RegressionCase struct {
	Case              string `json:"case"`
	Expected          string `json:"expected"`
	WhyDiscriminating string `json:"why_discriminating"`
}

func DecodeReview(data []byte, maxBytes int) (EcologicalReview, error) {
	if len(data) > maxBytes {
		return EcologicalReview{}, fmt.Errorf("final submission is %d bytes; cap is %d", len(data), maxBytes)
	}
	var review EcologicalReview
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&review); err != nil {
		return EcologicalReview{}, fmt.Errorf("decode final submission: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return EcologicalReview{}, fmt.Errorf("final submission contains trailing JSON")
	}
	if err := review.Validate(); err != nil {
		return EcologicalReview{}, err
	}
	return review, nil
}

func (r EcologicalReview) Validate() error {
	if strings.TrimSpace(r.ReviewSummary) == "" {
		return fmt.Errorf("review_summary is required")
	}
	if len(r.Findings) < 1 || len(r.Findings) > 8 || len(r.Uncertainties) > 8 {
		return fmt.Errorf("final submission has an invalid finding or uncertainty count")
	}
	for index, finding := range r.Findings {
		if err := finding.validate(); err != nil {
			return fmt.Errorf("finding %d: %w", index+1, err)
		}
	}
	for _, uncertainty := range r.Uncertainties {
		if strings.TrimSpace(uncertainty) == "" {
			return fmt.Errorf("uncertainties must not contain empty values")
		}
	}
	return nil
}

func (f EcologicalFinding) validate() error {
	for label, value := range map[string]string{
		"title": f.Title, "claim": f.Claim, "mechanism": f.Mechanism,
		"consequence": f.Consequence, "correction": f.Correction,
		"allocation_and_scope": f.AllocationAndScope,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", label)
		}
	}
	if len(f.Locations) < 1 || len(f.Locations) > 4 || len(f.RegressionTests) > 12 || len(f.Evidence) < 1 || len(f.Evidence) > 12 {
		return fmt.Errorf("location, test, or evidence count is invalid")
	}
	for _, location := range f.Locations {
		clean := path.Clean(location.Path)
		segments := strings.Split(location.Path, "/")
		if location.Path == "" || clean == "." || clean != location.Path || strings.Contains(location.Path, "\\") || path.IsAbs(clean) || containsSegment(segments, "..") || containsSegment(segments, ".git") {
			return fmt.Errorf("location path must stay inside the source snapshot")
		}
		if location.LineStart < 1 || location.LineEnd < location.LineStart {
			return fmt.Errorf("location line range is invalid")
		}
	}
	for _, test := range f.RegressionTests {
		if strings.TrimSpace(test.Case) == "" || strings.TrimSpace(test.Expected) == "" || strings.TrimSpace(test.WhyDiscriminating) == "" {
			return fmt.Errorf("regression test fields must not be empty")
		}
	}
	for _, evidence := range f.Evidence {
		if strings.TrimSpace(evidence) == "" {
			return fmt.Errorf("evidence must not contain empty values")
		}
	}
	if f.Confidence != "high" && f.Confidence != "medium" && f.Confidence != "low" {
		return fmt.Errorf("confidence must be high, medium, or low")
	}
	return nil
}

func containsSegment(segments []string, target string) bool {
	for _, segment := range segments {
		if segment == target {
			return true
		}
	}
	return false
}
