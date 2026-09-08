"""Recompute published calibration claims without model calls or new ratings.

Run: python3 scripts/evidence_audit.py
All input files are checked against their published manifest digests. Output is
JSON on stdout; historical results are never rewritten. Standard library only.
"""
import csv
import hashlib
import io
import itertools
import json
import math
from pathlib import Path
from statistics import mean

ROOT = Path(__file__).resolve().parents[1]
RESULTS = ROOT / "experiment/results"


def read_run(name):
    folder = RESULTS / name
    manifest = json.loads((folder / "manifest.json").read_text())
    for filename, digest in manifest["published_files"].items():
        if hashlib.sha256((folder / filename).read_bytes()).hexdigest() != digest:
            raise ValueError(f"Published digest mismatch: {name}/{filename}")
    rows = list(csv.DictReader(io.StringIO((folder / "sets.csv").read_text())))
    for row in rows:
        tp, fp, possible = (int(row[k]) for k in
                            ("true_positives", "false_positives", "possible_defects"))
        expected = 2 * tp / (tp + fp + possible)
        if not math.isclose(float(row["f1"]), expected, abs_tol=1e-12):
            raise ValueError(f"F1 mismatch: {row['set_id']}")
    return rows


def stage(rows):
    groups = {}
    for arm in sorted({r["arm"] for r in rows}):
        final = [r for r in rows if r["arm"] == arm and r["kind"] in ("final", "fused")]
        groups[arm] = {
            "n_tasks": len(final),
            "mean_f1": mean(float(r["f1"]) for r in final),
            "mean_recall": mean(float(r["recall"]) for r in final),
            "recorded_cost_usd": sum(float(r["cost_usd"]) for r in final),
            "recorded_total_tokens": sum(int(r["input_tokens"]) + int(r["output_tokens"]) for r in final),
        }
    return groups


def audit():
    haiku = stage(read_run("2026-08-20-stage-a-smoke"))
    gemma = stage(read_run("2026-08-20-stage-a-smoke-gemma"))
    rows = read_run("2026-08-20-persona-pair-gemma")
    paired = {}
    for row in rows:
        key = (row["packet"], int(row["repeat"]), row["provider_index"])
        pair = paired.setdefault(key, {})
        if row["wrapper"] in pair:
            raise ValueError(f"Duplicate pair member: {key}")
        pair[row["wrapper"]] = row
    if len(paired) != 30 or any(set(p) != {"functional", "fictional"} for p in paired.values()):
        raise ValueError("Expected exactly 30 complete pairs")
    deltas = [(k[0], float(p["fictional"]["f1"]) - float(p["functional"]["f1"]))
              for k, p in sorted(paired.items())]
    task_means = {t: mean(d for task, d in deltas if task == t) for t in sorted({t for t, _ in deltas})}
    # Descriptive sensitivity: entire task blocks change sign together. With
    # three blocks the two-sided exact test cannot have p below 2/8 = .25.
    observed = abs(mean(task_means.values()))
    values = list(task_means.values())
    permuted = [abs(mean(s * d for s, d in zip(signs, values)))
                for signs in itertools.product((-1, 1), repeat=len(values))]
    groups = {}
    for wrapper in ("functional", "fictional"):
        selected = [r for r in rows if r["wrapper"] == wrapper]
        groups[wrapper] = {
            "mean_f1": mean(float(r["f1"]) for r in selected),
            "input_tokens": sum(int(r["input_tokens"]) for r in selected),
            "output_tokens": sum(int(r["output_tokens"]) for r in selected),
        }
        groups[wrapper]["total_tokens"] = groups[wrapper]["input_tokens"] + groups[wrapper]["output_tokens"]
    return {
        "status": "post-hoc descriptive sensitivity; no new respondents or ratings",
        "published_file_digests_verified": True,
        "raw_seals_reverified": False,
        "haiku_final": haiku,
        "gemma_final": gemma,
        "haiku_M0_to_S1_recorded_cost_ratio": haiku["M0"]["recorded_cost_usd"] / haiku["S1"]["recorded_cost_usd"],
        "persona": {
            "groups": groups,
            "pairs": len(paired),
            "independent_task_count": len(task_means),
            "fictional_minus_functional_f1": mean(d for _, d in deltas),
            "per_task_delta": task_means,
            "leave_one_task_out_delta": {t: mean(d for task, d in deltas if task != t) for t in task_means},
            "fictional_wins": sum(d > 1e-12 for _, d in deltas),
            "ties": sum(abs(d) <= 1e-12 for _, d in deltas),
            "functional_wins": sum(d < -1e-12 for _, d in deltas),
            "task_block_sign_flip_two_sided_p": sum(v >= observed - 1e-12 for v in permuted) / len(permuted),
            "task_block_test_caution": "Illustrative symmetry/exchangeability assumption, only three selected tasks; not a preregistered population test.",
        },
        "metric_key": {
            "TP": "unique keyed defect detected", "FP": "unsupported claim (historical clustering rules differ)",
            "FN": "keyed defect missed", "precision": "TP/(TP+FP)", "recall": "TP/(TP+FN)",
            "F1": "2TP/(2TP+FP+FN)", "macro": "equal weight per task; persona repeats averaged within task",
            "cost": "provider-recorded usage, not independently verified billing; zero is not proof of free compute",
        },
    }


if __name__ == "__main__":
    print(json.dumps(audit(), indent=2, sort_keys=True))
