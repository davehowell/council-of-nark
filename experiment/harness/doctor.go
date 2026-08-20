package harness

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func wordCount(value string) int { return len(strings.Fields(value)) }
func ratio(a, b int) float64 {
	difference := a - b
	if difference < 0 {
		difference = -difference
	}
	maximum := a
	if b > a {
		maximum = b
	}
	if maximum == 0 {
		return 0
	}
	return float64(difference) / float64(maximum)
}
func (h *Harness) checkWrappers() []string {
	errors := []string{}
	rows, err := loadSpecialists(h.Root)
	if err != nil {
		return []string{err.Error()}
	}
	for _, role := range allRoles {
		row := rows[role]
		functional, fictional := wordCount(row["functional_wrapper"]), wordCount(row["fictional_wrapper"])
		difference := ratio(functional, fictional)
		fmt.Printf("wrapper %-13s functional=%2d fictional=%2d difference=%.1f%%\n", role, functional, fictional, difference*100)
		if difference > .10 {
			errors = append(errors, role+" wrappers differ by more than 10% in whitespace-token count")
		}
	}
	var wrappers map[string]string
	if err := readJSON(filepath.Join(h.Root, "experiment/prompts/omnibus-wrappers.json"), &wrappers); err != nil {
		return append(errors, err.Error())
	}
	functional, fictional := wordCount(wrappers["functional"]), wordCount(wrappers["fictional"])
	difference := ratio(functional, fictional)
	fmt.Printf("wrapper %-13s functional=%2d fictional=%2d difference=%.1f%%\n", "omnibus", functional, fictional, difference*100)
	if difference > .10 {
		errors = append(errors, "omnibus wrappers differ by more than 10%")
	}
	return errors
}
func (h *Harness) checkPrompts(config Config) ([]string, Plan) {
	errors := []string{}
	plan, err := BuildPlan(config)
	if err != nil {
		return []string{err.Error()}, Plan{}
	}
	callMap := map[string]Call{}
	for _, call := range plan.Calls {
		callMap[call.CallID] = call
	}
	for _, call := range plan.Calls {
		kind := stringValue(call.PromptSpec["kind"])
		var prompt string
		var err error
		if kind == "generic" || kind == "omnibus" || kind == "specialist" {
			prompt, err = staticPrompt(h.Root, call)
		} else if kind == "fuser" || kind == "chain" {
			deps := []dependencyPair{}
			for _, id := range call.DependsOn {
				deps = append(deps, dependencyPair{Call: callMap[id], Response: map[string]any{"findings": []any{}}})
			}
			prompt, err = dynamicPrompt(h.Root, call, deps)
		} else {
			err = fmt.Errorf("unknown prompt kind %s", kind)
		}
		if err == nil {
			var answer string
			answer, err = answerText(h.Root, call.Packet)
			if err == nil {
				err = assertAnswerKeyAbsent(prompt, answer)
			}
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("prompt %s: %v", call.CallID, err))
		}
	}
	fmt.Printf("assembled and contamination-checked %d planned prompts\n", len(plan.Calls))
	fmt.Printf("plan contains %d calls and %d scoreable output sets\n", len(plan.Calls), len(plan.OutputSets))
	return errors, plan
}
func runLookup(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
func (h *Harness) checkModels(config Config, skip bool) []string {
	if skip {
		return nil
	}
	providers := config.Providers
	if len(providers) == 0 {
		providers = []Provider{config.Provider}
	}
	errors := []string{}
	for _, p := range providers {
		if p.Adapter == "mock" {
			continue
		}
		if p.Adapter == "agy" {
			errors = append(errors, "agy is disabled: OAuth/keychain state conflicts with strict ephemeral-home isolation; use Pi with an explicit Google model")
			continue
		}
		if _, err := exec.LookPath(p.Adapter); err != nil {
			errors = append(errors, "missing executable: "+p.Adapter)
			continue
		}
		switch p.Adapter {
		case "pi", "claude":
			out, err := runLookup("pi", "--no-extensions", "--list-models", p.Model)
			registry := p.Model
			if parts := strings.SplitN(registry, "/", 2); len(parts) == 2 {
				registry = parts[1]
			}
			if err != nil || !strings.Contains(out, registry) {
				errors = append(errors, "model is not in local Pi registry: "+p.Model)
			}
		default:
			errors = append(errors, "unknown adapter: "+p.Adapter)
		}
	}
	return errors
}
func (h *Harness) Doctor(configPath string, skipLookup bool) error {
	if err := currentUserCheck(); err != nil {
		return err
	}
	if err := h.SandboxProbe(); err != nil {
		return err
	}
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(h.Root, configPath)
	}
	var config Config
	if err := readJSON(configPath, &config); err != nil {
		return err
	}
	errors := h.checkWrappers()
	promptErrors, _ := h.checkPrompts(config)
	errors = append(errors, promptErrors...)
	errors = append(errors, h.checkModels(config, skipLookup)...)
	if len(errors) > 0 {
		fmt.Println("Doctor failed:")
		for _, err := range errors {
			fmt.Println("- " + err)
		}
		return fmt.Errorf("doctor found %d error(s)", len(errors))
	}
	fmt.Println("Doctor passed. No model calls were made.")
	return nil
}
