from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
import csv
import json
from pathlib import Path
import time
from typing import Any

from .adapters import extract, invoke
from .common import ROOT, load_json, opaque_id, resolve_run, slug_timestamp, source_commit, utc_now, write_json
from .prompting import replace
from .run import DetachedWorktree
from .validate import judgement_response

FIELDS = [
    "rater",
    "item_id",
    "defect_id",
    "false_positive_cluster",
    "material",
    "confidence",
    "notes",
]


def read_jsonl(path: Path) -> list[dict[str, Any]]:
    return [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines() if line]


def existing_ratings(path: Path) -> dict[str, dict[str, Any]]:
    if not path.exists():
        return {}
    ratings: dict[str, dict[str, Any]] = {}
    with path.open(newline="", encoding="utf-8") as handle:
        for row in csv.DictReader(handle):
            item_id = row["item_id"]
            if item_id in ratings:
                raise SystemExit(f"Duplicate existing LLM rating: {item_id}")
            ratings[item_id] = row
    return ratings


def recover_captured_response(
    destination: Path, expected_ids: set[str]
) -> dict[str, Any] | None:
    """Recover a valid structured output after a parser-only failure."""
    attempts = sorted(destination.glob("attempts/*/stdout.txt"), reverse=True)
    for path in attempts:
        parsed, _ = extract(path.read_text(encoding="utf-8"), "judgements")
        if parsed is not None and not judgement_response(parsed, expected_ids):
            return parsed
    return None


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
    previous = load_json(destination / "record.json") if (destination / "record.json").exists() else None
    with DetachedWorktree(run, judge_id, freeze["source_commit"]) as worktree:
        # Rating is a derived stage with its own recorded harness commit. Load its
        # prompt/schema from that commit while keeping the packet key at the frozen
        # review-source commit.
        template = (ROOT / "experiment/prompts/judge.txt").read_text(encoding="utf-8")
        answer = (
            worktree / "experiment" / "scenarios" / output_set["packet"] / "answer-key.md"
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
        schema = load_json(ROOT / "experiment/schema/judgements.schema.json")
        expected_ids = {item["item_id"] for item in items}
        attempts = []
        result = None
        batch = slug_timestamp()
        for attempt in range(1, config["max_attempts"] + 1):
            result = invoke(
                cwd=worktree,
                provider=config["provider"],
                prompt=prompt,
                system=system,
                schema=schema,
                schema_relative="experiment/schema/judgements.schema.json",
                timeout_seconds=config["request_timeout_seconds"],
                expected_root="judgements",
                expected_ids=expected_ids,
            )
            attempt_dir = destination / "attempts" / f"{batch}-{attempt}"
            attempt_dir.mkdir(parents=True, exist_ok=True)
            (attempt_dir / "stdout.txt").write_text(result.stdout, encoding="utf-8")
            (attempt_dir / "stderr.txt").write_text(result.stderr, encoding="utf-8")
            write_json(
                attempt_dir / "metadata.json",
                {
                    "attempt": attempt,
                    "returncode": result.returncode,
                    "parse_error": result.parse_error,
                    "latency_seconds": result.latency_seconds,
                    "usage": result.usage,
                },
            )
            attempts.append(
                {
                    "attempt": attempt,
                    "returncode": result.returncode,
                    "parse_error": result.parse_error,
                    "latency_seconds": result.latency_seconds,
                    "usage": result.usage,
                }
            )
            if result.returncode == 0:
                break
            if attempt < config["max_attempts"]:
                time.sleep(float(config.get("retry_backoff_seconds", 5)) * attempt)
        assert result is not None
        status = "success" if result.returncode == 0 and result.parsed is not None else "error"
        write_json(
            destination / "record.json",
            {
                "set_id": set_id,
                "packet": output_set["packet"],
                "provider": config["provider"],
                "status": status,
                "attempts": attempts,
                "parsed_response": result.parsed,
                "completed_at": utc_now(),
                "arm_blinded": True,
                "status_note": config["status"],
                "review_source_commit": freeze["source_commit"],
                "judge_harness_commit": source_commit(),
                "previous_record": previous,
            },
        )
        return set_id, result.parsed


def merge_ratings(
    ratings: dict[str, dict[str, Any]], response: dict[str, Any], model: str
) -> None:
    for judgement in response["judgements"]:
        ratings[judgement["item_id"]] = {
            "rater": f"llm:{model}",
            "item_id": judgement["item_id"],
            "defect_id": judgement["defect_id"],
            "false_positive_cluster": judgement["false_positive_cluster"],
            "material": judgement["material"],
            "confidence": judgement["confidence"],
            "notes": judgement["rationale"],
        }


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

    path = run / "blinded" / "ratings-llm.csv"
    ratings = existing_ratings(path)
    known_item_ids = {item["item_id"] for item in items}
    unknown = set(ratings) - known_item_ids
    if unknown:
        raise SystemExit(f"Existing ratings contain unknown item IDs: {sorted(unknown)}")

    pending = []
    for output_set in plan["output_sets"]:
        set_items = by_set.get(output_set["set_id"], [])
        expected = {item["item_id"] for item in set_items}
        if not expected or expected.issubset(ratings):
            continue
        # A partially rated set is rerated as one blinded unit.
        for item_id in expected:
            ratings.pop(item_id, None)
        destination = run / "blinded" / "llm-triage" / output_set["set_id"]
        recovered = recover_captured_response(destination, expected)
        if recovered is not None:
            merge_ratings(ratings, recovered, config["provider"]["model"])
            record_path = destination / "record.json"
            if record_path.exists():
                record = load_json(record_path)
                record["status"] = "success"
                record["parsed_response"] = recovered
                record["recovered_at"] = utc_now()
                record["recovery"] = "parser-only recovery from captured structured output"
                record["judge_harness_commit"] = source_commit()
                write_json(record_path, record)
            print(f"{output_set['set_id']} recovered from captured output", flush=True)
            continue
        pending.append(output_set)
    print(f"Resuming with {len(ratings)} existing/recovered ratings; {len(pending)} sets pending.")

    errors = []
    with ThreadPoolExecutor(max_workers=args.jobs) as pool:
        futures = {
            pool.submit(
                rate_set,
                run,
                output_set,
                by_set.get(output_set["set_id"], []),
                freeze,
                config,
            ): output_set
            for output_set in pending
        }
        for future in as_completed(futures):
            output_set = futures[future]
            set_id, response = future.result()
            if response is None:
                errors.append(set_id)
                print(f"{set_id} rating error", flush=True)
                continue
            merge_ratings(ratings, response, config["provider"]["model"])
            print(f"{set_id} rated", flush=True)

    with path.open("w", newline="", encoding="utf-8") as handle:
        writer = csv.DictWriter(handle, fieldnames=FIELDS, lineterminator="\n")
        writer.writeheader()
        writer.writerows(ratings[item_id] for item_id in sorted(ratings))
    missing = known_item_ids - set(ratings)
    print(f"Wrote {len(ratings)} exploratory ratings to {path.relative_to(run)}")
    if errors or missing:
        if errors:
            print(f"Unrated sets: {', '.join(errors)}")
        print(f"Unrated items remaining: {len(missing)}")
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
