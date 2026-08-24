package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/davehowell/council-of-nark/experiment/ecological/respondent"
	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

func main() {
	snapshotPath := flag.String("snapshot", "", "completed clean Gortex snapshot attempt")
	configPath := flag.String("config", "experiment/ecological/config/eco-gortex-unicode-tokenizer-tools.json", "tracked respondent-tool config")
	flag.Parse()
	if *snapshotPath == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: ecological-pi-doctor --snapshot <attempt> [--config <path>]")
		os.Exit(2)
	}
	root, err := repositoryRoot()
	if err != nil {
		fatal(err)
	}
	attempt, err := newAttempt(filepath.Join(root, "experiment", "ecological", "respondent-runs"))
	if err != nil {
		fatal(err)
	}
	started := time.Now().UTC()
	status := map[string]any{
		"schema_version": 1, "task_id": "eco-gortex-unicode-tokenizer",
		"status": "running", "started_at": started.Format(time.RFC3339Nano),
	}
	_ = writeJSON(filepath.Join(attempt, "status.json"), status)
	result, runErr := respondent.Doctor(root, *snapshotPath, *configPath, attempt)
	if runErr == nil {
		controller, err := controllerState(root)
		if err != nil {
			runErr = err
		} else {
			result.Artifacts["controller"] = controller
			runErr = writeJSON(filepath.Join(attempt, "doctor.json"), result)
		}
	}
	if runErr == nil {
		doctorDigest, err := snapshot.FileSHA256(filepath.Join(attempt, "doctor.json"))
		if err != nil {
			runErr = err
		} else {
			runErr = writeJSON(filepath.Join(attempt, "seal.json"), map[string]any{
				"schema_version": 1, "doctor_sha256": doctorDigest,
				"completed_at": time.Now().UTC().Format(time.RFC3339Nano),
			})
		}
	}
	status["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	if runErr != nil {
		status["status"] = "failed"
		status["error"] = runErr.Error()
	} else {
		status["status"] = "success"
	}
	_ = writeJSON(filepath.Join(attempt, "status.json"), status)
	printAttempt(root, attempt)
	if runErr != nil {
		fmt.Fprintln(os.Stderr, "error:", runErr)
		os.Exit(1)
	}
}

func controllerState(root string) (map[string]any, error) {
	commit, err := gitOutput(root, "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	status, err := gitOutput(root, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}
	manifest, err := snapshot.BuildTreeManifest(filepath.Join(root, "experiment", "ecological", "respondent"))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"commit": strings.TrimSpace(commit), "tree_dirty": strings.TrimSpace(status) != "",
		"respondent_tree_sha256": manifest.TreeSHA256,
	}, nil
}

func gitOutput(root string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func repositoryRoot() (string, error) {
	working, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(working, ".git")); err == nil {
			return working, nil
		}
		parent := filepath.Dir(working)
		if parent == working {
			return "", fmt.Errorf("not inside repository")
		}
		working = parent
	}
}

func newAttempt(parent string) (string, error) {
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return "", err
	}
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", err
	}
	path := filepath.Join(parent, time.Now().UTC().Format("20060102T150405Z")+"-gortex-pi-doctor-"+hex.EncodeToString(suffix[:]))
	if err := os.Mkdir(path, 0o700); err != nil {
		return "", err
	}
	return path, nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func printAttempt(root, attempt string) {
	if relative, err := filepath.Rel(root, attempt); err == nil {
		fmt.Println(filepath.ToSlash(relative))
		return
	}
	fmt.Println(attempt)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
