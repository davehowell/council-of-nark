from __future__ import annotations

import argparse
from collections import defaultdict
import csv
import json
import math
from pathlib import Path
from statistics import mean
from typing import Any

from .common import load_json, read_call_record, resolve_run, write_json
from .summarize import token_counts


def read_jsonl(path: Path) -> list[dict[str, Any]]:
    return [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines() if line]


def answer_ids(root: Path, packet: str) -> set[str]:
    text = (root / "experiment" / "scenarios" / packet / "answer-key.md").read_text(encoding="utf-8")
    return {
        part.strip()
        for line in text.splitlines()
        if line.startswith("|")
        for part in [line.strip("|").split("|")[0]]
        if len(part.strip()) == 5 and part.strip()[2] == "-" and part.strip()[3:].isdigit()
    }


def percentile(values: list[float], q: float) -> float | None:
    if not values:
        return None
    ordered = sorted(values)
    index = max(0, math.ceil(q * len(ordered)) - 1)
    return ordered[index]


def set_usage(run: Path, call_ids: list[str]) -> tuple[int, int, float, float]:
    input_tokens = output_tokens = 0
    cost = latency = 0.0
    for call_id in call_ids:
        record = read_call_record(run, call_id)
        if not record:
            continue
        inp, out, call_cost = token_counts(record.get("usage", {}))
        input_tokens += inp
        output_tokens += out
        cost += call_cost
        latency += float(record.get("latency_seconds", 0) or 0)
    return input_tokens, output_tokens, cost, latency


def group_name(metadata: dict[str, Any]) -> str:
    design = metadata["design"]
    if design == "stage_a":
        return f"{metadata['arm']}:{metadata['kind']}"
    if design == "topology":
        return f"{metadata['topology']}:{metadata['kind']}"
    return f"{metadata['adapter']}:{metadata['wrapper']}"


def main() -> int:
    parser = argparse.ArgumentParser(description="Score blinded ratings and unblind aggregate conditions")
    parser.add_argument("run")
    parser.add_argument("ratings")
    parser.add_argument("--label", default="adjudicated")
    args = parser.parse_args()
    run = resolve_run(args.run)
    root = Path(__file__).resolve().parents[2]
    ratings_path = Path(args.ratings)
    if not ratings_path.is_absolute():
        ratings_path = run / ratings_path
    plan = load_json(run / "plan.json")
    items = {row["item_id"]: row for row in read_jsonl(run / "blinded" / "findings.jsonl")}
    sets = {row["set_id"]: row for row in read_jsonl(run / "blinded" / "sets.jsonl")}
    ratings: dict[str, dict[str, str]] = {}
    with ratings_path.open(newline="", encoding="utf-8") as handle:
        for row in csv.DictReader(handle):
            item_id = row["item_id"]
            if item_id in ratings:
                raise SystemExit(f"Duplicate rating for {item_id}; adjudicate to one row per item")
            ratings[item_id] = row
    missing = sorted(set(items) - set(ratings))
    extra = sorted(set(ratings) - set(items))
    if missing or extra:
        raise SystemExit(f"Rating coverage mismatch; missing={missing}, extra={extra}")

    plan_sets = {row["set_id"]: row for row in plan["output_sets"]}
    rows: list[dict[str, Any]] = []
    detected_by_set: dict[str, set[str]] = {}
    for set_id, blinded_set in sets.items():
        output_set = plan_sets[set_id]
        valid = answer_ids(root, output_set["packet"])
        detected: set[str] = set()
        false_positives = 0
        for item_id in blinded_set["item_ids"]:
            rating = ratings[item_id]
            defect = rating["defect_id"].strip().upper()
            material = rating["material"].strip().lower()
            if material not in {"true", "false"}:
                raise SystemExit(f"Invalid material value for {item_id}")
            if defect in {"", "NONE", "NULL"} or material == "false":
                false_positives += 1
            elif defect not in valid:
                raise SystemExit(f"Invalid defect ID {defect!r} for packet {output_set['packet']}")
            else:
                detected.add(defect)
        tp = len(detected)
        predicted = tp + false_positives
        precision = tp / predicted if predicted else 0.0
        recall = tp / len(valid) if valid else 0.0
        f1 = 2 * precision * recall / (precision + recall) if precision + recall else 0.0
        inp, out, cost, latency = set_usage(run, output_set["cost_call_ids"])
        metadata = output_set["metadata"]
        row = {
            "set_id": set_id,
            "packet": output_set["packet"],
            **metadata,
            "group": group_name(metadata),
            "true_positives": tp,
            "false_positives": false_positives,
            "possible_defects": len(valid),
            "precision": precision,
            "recall": recall,
            "f1": f1,
            "input_tokens": inp,
            "output_tokens": out,
            "cost_usd": cost,
            "latency_seconds_serial_sum": latency,
            "true_findings_per_1k_output_tokens": (1000 * tp / out) if out else None,
            "true_findings_per_dollar": (tp / cost) if cost else None,
        }
        rows.append(row)
        detected_by_set[set_id] = detected

    grouped: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for row in rows:
        grouped[row["group"]].append(row)
    summary = {}
    for name, group in grouped.items():
        f1s = [row["f1"] for row in group]
        summary[name] = {
            "n_sets": len(group),
            "mean_f1": mean(f1s),
            "p10_f1": percentile(f1s, 0.10),
            "worst_f1": min(f1s),
            "mean_precision": mean(row["precision"] for row in group),
            "mean_recall": mean(row["recall"] for row in group),
            "mean_output_tokens": mean(row["output_tokens"] for row in group),
            "total_cost_usd": sum(row["cost_usd"] for row in group),
        }

    fusion = []
    index: dict[tuple[Any, ...], dict[str, dict[str, Any]]] = defaultdict(dict)
    for row in rows:
        if row.get("design") != "stage_a" or row.get("arm") not in {"M0", "M1", "M2"}:
            continue
        key = (row["packet"], row["arm"], row["repeat"])
        index[key][row["kind"]] = row
    for key, pair in index.items():
        if "raw_union" in pair and "fused" in pair:
            raw, fused = pair["raw_union"], pair["fused"]
            fusion.append(
                {
                    "packet": key[0],
                    "arm": key[1],
                    "repeat": key[2],
                    "raw_true_positives": raw["true_positives"],
                    "fused_true_positives": fused["true_positives"],
                    "fusion_retention": (
                        fused["true_positives"] / raw["true_positives"]
                        if raw["true_positives"]
                        else None
                    ),
                    "raw_false_positives": raw["false_positives"],
                    "fused_false_positives": fused["false_positives"],
                }
            )

    analysis = run / "analysis" / args.label
    analysis.mkdir(parents=True, exist_ok=True)
    write_json(analysis / "summary.json", {"groups": summary, "fusion": fusion})
    fields = sorted({key for row in rows for key in row})
    with (analysis / "sets.csv").open("w", newline="", encoding="utf-8") as handle:
        writer = csv.DictWriter(handle, fieldnames=fields)
        writer.writeheader()
        writer.writerows(rows)
    print(json.dumps({"groups": summary, "fusion": fusion}, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
