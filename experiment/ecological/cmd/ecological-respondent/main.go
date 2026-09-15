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
)

func main() {
	snapshotPath := flag.String("snapshot", "", "completed clean ecological snapshot attempt")
	configPath := flag.String("config", "experiment/ecological/config/eco-gortex-unicode-tokenizer-respondent-mock.json", "tracked ecological respondent config")
	allowLive := flag.Bool("allow-live", false, "permit a provider-backed run from an explicitly live config")
	verifyPath := flag.String("verify", "", "verify a sealed respondent attempt and exit")
	flag.Parse()
	if flag.NArg() != 0 || (*verifyPath == "" && *snapshotPath == "") || (*verifyPath != "" && *snapshotPath != "") {
		fmt.Fprintln(os.Stderr, "usage: ecological-respondent (--snapshot <snapshot> [--config <config>] [--allow-live] | --verify <attempt>)")
		os.Exit(2)
	}
	root, err := repositoryRoot()
	if err != nil {
		fatal(err)
	}
	if *verifyPath != "" {
		path := *verifyPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		seal, err := respondent.VerifyAttemptSeal(path)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("verified %d files: %s\n", len(seal.Files), seal.TreeSHA256)
		return
	}
	config, _, err := respondent.LoadRunConfig(root, *configPath)
	if err != nil {
		fatal(err)
	}
	if config.Adapter == "pi" && !*allowLive {
		fatal(fmt.Errorf("live provider execution requires --allow-live after every launch gate is frozen"))
	}
	attempt, err := newAttempt(filepath.Join(root, "experiment", "ecological", "respondent-runs"), config.TaskID, config.Adapter)
	if err != nil {
		fatal(err)
	}
	started := time.Now().UTC()
	status := map[string]any{
		"schema_version": 1, "task_id": config.TaskID, "adapter": config.Adapter,
		"status": "running", "started_at": started.Format(time.RFC3339Nano),
	}
	if err := writeJSON(filepath.Join(attempt, "status.json"), status); err != nil {
		fatal(err)
	}
	result, runErr := respondent.ExecuteRun(root, *snapshotPath, *configPath, attempt)
	if err := respondent.RemoveEphemeralRunState(attempt); err != nil && runErr == nil {
		runErr = err
	}
	status["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	status["duration_seconds"] = time.Since(started).Seconds()
	if runErr != nil {
		status["status"] = "failed"
		status["error"] = runErr.Error()
	} else {
		status["status"] = "success"
		status["outcome"] = result.Outcome
	}
	if err := writeJSON(filepath.Join(attempt, "status.json"), status); err != nil && runErr == nil {
		runErr = err
	}
	if _, err := respondent.SealAttempt(attempt); err != nil && runErr == nil {
		runErr = err
	}
	if _, err := respondent.VerifyAttemptSeal(attempt); err != nil && runErr == nil {
		runErr = err
	}
	printAttempt(root, attempt)
	if runErr != nil {
		fmt.Fprintln(os.Stderr, "error:", runErr)
		os.Exit(1)
	}
}

func repositoryRoot() (string, error) {
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("locate repository root: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func newAttempt(parent, task, adapter string) (string, error) {
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return "", err
	}
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s-%s-respondent-%s-%s", time.Now().UTC().Format("20060102T150405Z"), task, adapter, hex.EncodeToString(suffix[:]))
	path := filepath.Join(parent, name)
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
