from __future__ import annotations

import argparse
import json
from pathlib import Path
import shutil
import subprocess

from .common import ROOT, load_json
from .plan import build
from .prompting import answer_text, assert_answer_key_absent, dynamic_prompt, static_prompt


def words(value: str) -> int:
    return len(value.split())


def ratio(a: int, b: int) -> float:
    return abs(a - b) / max(a, b)


def check_wrappers() -> list[str]:
    errors: list[str] = []
    specialists = load_json(ROOT / "experiment/prompts/specialists.json")
    for row in specialists:
        functional = words(row["functional_wrapper"])
        fictional = words(row["fictional_wrapper"])
        difference = ratio(functional, fictional)
        print(
            f"wrapper {row['id']:<13} functional={functional:2d} fictional={fictional:2d} difference={difference:.1%}"
        )
        if difference > 0.10:
            errors.append(f"{row['id']} wrappers differ by more than 10% in whitespace-token count")
    wrappers = load_json(ROOT / "experiment/prompts/omnibus-wrappers.json")
    functional, fictional = words(wrappers["functional"]), words(wrappers["fictional"])
    difference = ratio(functional, fictional)
    print(f"wrapper {'omnibus':<13} functional={functional:2d} fictional={fictional:2d} difference={difference:.1%}")
    if difference > 0.10:
        errors.append("omnibus wrappers differ by more than 10% in whitespace-token count")
    return errors


def lookup_models(config: dict, skip: bool) -> list[str]:
    errors: list[str] = []
    providers = config.get("providers", [config.get("provider")])
    for provider in [value for value in providers if value]:
        adapter, model = provider["adapter"], provider["model"]
        executable = {"claude": "claude", "agy": "agy", "pi": "pi", "mock": None}[adapter]
        if executable and not shutil.which(executable):
            errors.append(f"missing executable: {executable}")
            continue
        if skip or adapter == "mock":
            continue
        if adapter == "agy":
            result = subprocess.run(["agy", "models"], text=True, capture_output=True, timeout=60)
            if result.returncode or model not in result.stdout:
                errors.append(f"agy model is not listed: {model}")
        elif adapter in {"pi", "claude"}:
            result = subprocess.run(
                ["pi", "--no-extensions", "--list-models", model],
                text=True,
                capture_output=True,
                timeout=60,
            )
            registry_id = model.split("/", 1)[-1]
            if result.returncode or registry_id not in result.stdout:
                errors.append(f"model is not in the local pi registry: {model}")
    return errors


def check_prompts(config: dict) -> list[str]:
    errors: list[str] = []
    plan = build(config)
    call_map = {call["call_id"]: call for call in plan["calls"]}
    for call in plan["calls"]:
        kind = call["prompt_spec"]["kind"]
        try:
            if kind in {"generic", "omnibus", "specialist"}:
                prompt = static_prompt(ROOT, call)
            elif kind == "fuser":
                dependencies = [
                    (call_map[call_id], {"findings": []}) for call_id in call["depends_on"]
                ]
                prompt = dynamic_prompt(ROOT, call, dependencies)
            elif kind == "chain":
                dependency = call_map[call["depends_on"][0]]
                prompt = dynamic_prompt(ROOT, call, [(dependency, {"findings": []})])
            else:
                errors.append(f"unknown prompt kind: {kind}")
                continue
            assert_answer_key_absent(prompt, answer_text(ROOT, call["packet"]))
        except Exception as error:  # surfaced with call identity for diagnosis
            errors.append(f"prompt {call['call_id']}: {error}")
    print(f"assembled and contamination-checked {len(plan['calls'])} planned prompts")
    print(f"plan contains {len(plan['calls'])} calls and {len(plan['output_sets'])} scoreable output sets")
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate an experiment config without model calls")
    parser.add_argument("config")
    parser.add_argument("--skip-model-lookup", action="store_true")
    args = parser.parse_args()
    config_path = Path(args.config)
    if not config_path.is_absolute():
        config_path = ROOT / config_path
    config = load_json(config_path)

    errors = []
    errors.extend(check_wrappers())
    errors.extend(check_prompts(config))
    errors.extend(lookup_models(config, args.skip_model_lookup))
    if errors:
        print("Doctor failed:")
        for error in errors:
            print(f"- {error}")
        return 1
    print("Doctor passed. No model calls were made.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
