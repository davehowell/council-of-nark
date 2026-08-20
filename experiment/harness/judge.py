from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
import csv
import json
from pathlib import Path
from typing import Any

from .adapters import invoke
from .common import load_json, opaque_id, resolve_run, utc_now, write_json
from .prompting import replace
from .run import DetachedWorktree


def read_jsonl(path: Path) -> list[dict[str, Any]]:
    return [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines() if line]


def rate_set(
    run: Path,
    output_set: dict[str, Any],
    items: list[dict[str, Any]],
    freeze: dict[str, Any],
    config: dict[str, Any],
) -> tuple[str, dict[str, Any] | None]:
    set_id = output_set["set_id"]
    if not items:
        return set_id, {"judgements": []}
    judge_id = "j-" + opaque_id(config["label"], set_id)
    destination = run / "blinded" / "llm-triage" / set_id
    destination.mkdir(parents=True, exist_ok=True)
    with DetachedWorktree(run, judge_id, freeze["source_commit"]) as worktree:
        template = (worktree / "experiment/prompts/judge.txt").read_text(encoding="utf-8")
        answer = (
            worktree
            / "experiment"
            / "scenarios"
            / output_set["packet"]
            / "answer-key.md"
        ).read_text(encoding="utf-8")
        public_items = [
            {"item_id": item["item_id"], "finding": item["finding"]} for item in items
        ]
        prompt = replace(
            template,
            ANSWER_KEY=answer,
            FINDINGS=json.dumps(public_items, ensure_ascii=False),
        )
        system = (
            "You are an isolated blinded evaluation rater. Use only the supplied answer key and findings. "
            "Do not use tools, memories, project context, or other sessions. Return only schema-valid JSON."
        )
        schema = load_json(worktree / "experiment/schema/judgements.schema.json")
        result = invoke(
            cwd=worktree,
            provider=config["provider"],
            prompt=prompt,
            system=system,
            schema=schema,
            schema_relative="experiment/schema/judgements.schema.json",
            timeout_seconds=config["request_timeout_seconds"],
            expected_root="judgements",
            expected_ids={item["item_id"] for item in items},
        )
        (destination / "stdout.txt").write_text(result.stdout, encoding="utf-8")
        (destination / "stderr.txt").write_text(result.stderr, encoding="utf-8")
        write_json(
            destination / "record.json",
            {
                "set_id": set_id,
                "packet": output_set["packet"],
                "provider": config["provider"],
                "status": "success" if result.returncode == 0 and result.parsed else "error",
                "returncode": result.returncode,
                "parse_error": result.parse_error,
                "latency_seconds": result.latency_seconds,
                "usage": result.usage,
                "completed_at": utc_now(),
                "arm_blinded": True,
                "status_note": config["status"],
            },
        )
        return set_id, result.parsed


def main() -> int:
    parser = argparse.ArgumentParser(description="Create exploratory arm-blinded LLM smoke ratings")
    parser.add_argument("run")
    parser.add_argument("--config", default="experiment/config/judge-smoke.json")
    parser.add_argument("--jobs", type=int, default=2)
    args = parser.parse_args()
    run = resolve_run(args.run)
    if not (run / "seal.json").exists():
        raise SystemExit("Seal raw model output before rating")
    root = Path(__file__).resolve().parents[2]
    config_path = Path(args.config)
    if not config_path.is_absolute():
        config_path = root / config_path
    config = load_json(config_path)
    freeze = load_json(run / "freeze.json")
    plan = load_json(run / "plan.json")
    items = read_jsonl(run / "blinded" / "findings.jsonl")
    by_set: dict[str, list[dict[str, Any]]] = {}
    for item in items:
        by_set.setdefault(item["set_id"], []).append(item)

    ratings: list[dict[str, Any]] = []
    errors = []
    with ThreadPoolExecutor(max_workers=args.jobs) as pool:
        futures = {
            pool.submit(rate_set, run, output_set, by_set.get(output_set["set_id"], []), freeze, config): output_set
            for output_set in plan["output_sets"]
        }
        for future in as_completed(futures):
            output_set = futures[future]
            set_id, response = future.result()
            if response is None:
                errors.append(set_id)
                print(f"{set_id} rating error", flush=True)
                continue
            for judgement in response["judgements"]:
                ratings.append(
                    {
                        "rater": f"llm:{config['provider']['model']}",
                        "item_id": judgement["item_id"],
                        "defect_id": judgement["defect_id"],
                        "material": judgement["material"],
                        "confidence": judgement["confidence"],
                        "notes": judgement["rationale"],
                    }
                )
            print(f"{set_id} rated", flush=True)

    path = run / "blinded" / "ratings-llm.csv"
    with path.open("w", newline="", encoding="utf-8") as handle:
        fields = ["rater", "item_id", "defect_id", "material", "confidence", "notes"]
        writer = csv.DictWriter(handle, fieldnames=fields)
        writer.writeheader()
        writer.writerows(ratings)
    print(f"Wrote {len(ratings)} exploratory ratings to {path.relative_to(run)}")
    if errors:
        print(f"Unrated sets: {', '.join(errors)}")
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
