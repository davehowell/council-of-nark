package mediator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

type AttemptRequirements struct {
	TaskID           string
	ControllerCommit string
	SourceTreeSHA256 string
}

type VerifiedAttempt struct {
	Attempt          string
	Source           string
	Controller       string
	TaskID           string
	ControllerCommit string
	SourceTreeSHA256 string
	SourceEntries    int
	SourceBytes      int64
}

func VerifySnapshotAttempt(attempt string, requirements AttemptRequirements) (VerifiedAttempt, error) {
	attempt, err := filepath.Abs(attempt)
	if err != nil {
		return VerifiedAttempt{}, err
	}
	var status struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
		Phase  string `json:"phase"`
	}
	if err := readJSONFile(filepath.Join(attempt, "status.json"), &status); err != nil {
		return VerifiedAttempt{}, fmt.Errorf("read snapshot status: %w", err)
	}
	if status.Status != "success" || status.Phase != "complete" {
		return VerifiedAttempt{}, fmt.Errorf("snapshot attempt is not complete and successful")
	}
	if requirements.TaskID != "" && status.TaskID != requirements.TaskID {
		return VerifiedAttempt{}, fmt.Errorf("snapshot task mismatch: got %q", status.TaskID)
	}
	controller := filepath.Join(attempt, "controller")
	source := filepath.Join(attempt, "source")
	var seal struct {
		ProvenanceSHA256     string `json:"provenance_sha256"`
		SourceManifestSHA256 string `json:"source_manifest_sha256"`
		SourceTreeSHA256     string `json:"source_tree_sha256"`
	}
	if err := readJSONFile(filepath.Join(controller, "seal.json"), &seal); err != nil {
		return VerifiedAttempt{}, fmt.Errorf("read snapshot seal: %w", err)
	}
	for path, expected := range map[string]string{
		filepath.Join(controller, "provenance.json"):      seal.ProvenanceSHA256,
		filepath.Join(controller, "source-manifest.json"): seal.SourceManifestSHA256,
	} {
		actual, err := snapshot.FileSHA256(path)
		if err != nil || actual != expected {
			return VerifiedAttempt{}, fmt.Errorf("snapshot seal digest mismatch: %s", filepath.Base(path))
		}
	}
	var provenance struct {
		TaskID     string `json:"task_id"`
		Controller struct {
			Commit    string `json:"commit"`
			TreeDirty bool   `json:"tree_dirty"`
		} `json:"controller"`
		Validation map[string]any `json:"validation"`
		Isolation  map[string]any `json:"isolation"`
	}
	if err := readJSONFile(filepath.Join(controller, "provenance.json"), &provenance); err != nil {
		return VerifiedAttempt{}, err
	}
	if provenance.TaskID != status.TaskID || provenance.Controller.TreeDirty {
		return VerifiedAttempt{}, fmt.Errorf("snapshot provenance task differs or controller tree was dirty")
	}
	if requirements.ControllerCommit != "" && provenance.Controller.Commit != requirements.ControllerCommit {
		return VerifiedAttempt{}, fmt.Errorf("snapshot controller commit mismatch")
	}
	for label, values := range map[string]map[string]any{"validation": provenance.Validation, "isolation": provenance.Isolation} {
		if len(values) == 0 {
			return VerifiedAttempt{}, fmt.Errorf("snapshot %s evidence is absent", label)
		}
	}
	for _, key := range []string{"closure_verified_offline", "network_probe_denied", "unlisted_executable_denied"} {
		if value, ok := provenance.Validation[key].(bool); !ok || !value {
			return VerifiedAttempt{}, fmt.Errorf("snapshot validation check %q did not pass", key)
		}
	}
	for _, key := range []string{"source_read", "source_write_denied", "council_read_denied", "controller_metadata_read_denied", "sibling_task_read_denied", "evidence_read_denied", "git_history_absent", "network_denied"} {
		if value, ok := provenance.Isolation[key].(bool); !ok || !value {
			return VerifiedAttempt{}, fmt.Errorf("snapshot isolation check %q did not pass", key)
		}
	}
	if _, err := os.Lstat(filepath.Join(source, ".git")); err == nil || !os.IsNotExist(err) {
		return VerifiedAttempt{}, fmt.Errorf("snapshot source contains Git metadata")
	}
	manifest, err := snapshot.BuildTreeManifest(source)
	if err != nil {
		return VerifiedAttempt{}, fmt.Errorf("recompute source tree: %w", err)
	}
	if manifest.TreeSHA256 != seal.SourceTreeSHA256 {
		return VerifiedAttempt{}, fmt.Errorf("snapshot source tree differs from seal")
	}
	if requirements.SourceTreeSHA256 != "" && manifest.TreeSHA256 != requirements.SourceTreeSHA256 {
		return VerifiedAttempt{}, fmt.Errorf("snapshot source tree differs from frozen requirement")
	}
	return VerifiedAttempt{
		Attempt: attempt, Source: source, Controller: controller, TaskID: status.TaskID,
		ControllerCommit: provenance.Controller.Commit, SourceTreeSHA256: manifest.TreeSHA256,
		SourceEntries: manifest.FileCount, SourceBytes: manifest.ByteCount,
	}, nil
}
