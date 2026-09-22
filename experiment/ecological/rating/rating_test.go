package rating

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func assetDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate rating assets")
	}
	return filepath.Dir(file)
}

func TestRatingSchemasAreClosedJSONObjects(t *testing.T) {
	for _, name := range []string{"rating-bundle.schema.json", "rating-record.schema.json"} {
		data, err := os.ReadFile(filepath.Join(assetDir(t), name))
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if schema["type"] != "object" || schema["additionalProperties"] != false {
			t.Fatalf("%s must reject unlisted top-level fields", name)
		}
		required, ok := schema["required"].([]any)
		if !ok || len(required) == 0 {
			t.Fatalf("%s has no required fields", name)
		}
	}
}

func TestRatingPageHasNoRemoteDependenciesAndIncludesFrozenFields(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(assetDir(t), "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(data)
	for _, forbidden := range []string{
		`src="http://`, `src="https://`, `href="http://`, `href="https://`,
		"fetch(", "XMLHttpRequest", "WebSocket(", "EventSource(",
	} {
		if strings.Contains(page, forbidden) {
			t.Fatalf("rating page contains forbidden network capability %q", forbidden)
		}
	}
	for _, required := range []string{
		"default-src 'none'", "bundle_sha256", "reference_relation",
		"technical.root_cause", "utility.personal_usefulness",
		"treatment_guess", "wording_revealed",
	} {
		if !strings.Contains(page, required) {
			t.Fatalf("rating page is missing %q", required)
		}
	}
}
