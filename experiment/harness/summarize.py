from __future__ import annotations

import argparse
from collections import defaultdict
import json
from statistics import mean
from typing import Any

from .common import load_json, read_call_record, resolve_run, write_json


def token_counts(usage: dict[str, Any]) -> tuple[int, int, float]:
    source = usage.get("usage") if isinstance(usage.get("usage"), dict) else usage
    input_tokens = int(
        source.get("input_tokens", source.get("inputTokens", source.get("input", 0))) or 0
    )
    output_tokens = int(
        source.get("output_tokens", source.get("outputTokens", source.get("output", 0))) or 0
    )
    # Pi reports reasoning separately from visible output. Include it in the
    # billable/generated token total used for cross-adapter efficiency.
    output_tokens += int(source.get("reasoning", 0) or 0)
    cost = float(usage.get("total_cost_usd", usage.get("cost_usd", 0)) or 0)
    if not cost and isinstance(source.get("cost"), dict):
        cost = float(source["cost"].get("total", 0) or 0)
    if not (input_tokens or output_tokens) and isinstance(usage.get("modelUsage"), dict):
        for model in usage["modelUsage"].values():
            if isinstance(model, dict):
                input_tokens += int(model.get("inputTokens", 0) or 0)
                output_tokens += int(model.get("outputTokens", 0) or 0)
                cost += float(model.get("costUSD", 0) or 0)
    return input_tokens, output_tokens, cost


def main() -> int:
    parser = argparse.ArgumentParser(description="Summarize run health, usage, and latency without scoring")
    parser.add_argument("run")
    args = parser.parse_args()
    run = resolve_run(args.run)
    plan = load_json(run / "plan.json")
    statuses: dict[str, int] = defaultdict(int)
    providers: dict[str, dict[str, Any]] = defaultdict(lambda: {
        "calls": 0, "latencies": [], "input_tokens": 0, "output_tokens": 0, "cost_usd": 0.0
    })
    total_findings = 0
    for call in plan["calls"]:
        record = read_call_record(run, call["call_id"])
        status = record.get("status", "missing") if record else "missing"
        statuses[status] += 1
        key = f"{call['provider']['adapter']}:{call['provider']['model']}"
        bucket = providers[key]
        bucket["calls"] += 1
        if record:
            if isinstance(record.get("latency_seconds"), (int, float)):
                bucket["latencies"].append(record["latency_seconds"])
            inp, out, cost = token_counts(record.get("usage", {}))
            bucket["input_tokens"] += inp
            bucket["output_tokens"] += out
            bucket["cost_usd"] += cost
            if status == "success":
                total_findings += len(record["parsed_response"]["findings"])
    output_providers = {}
    for key, bucket in providers.items():
        latencies = bucket.pop("latencies")
        output_providers[key] = {
            **bucket,
            "mean_latency_seconds": mean(latencies) if latencies else None,
            "total_latency_seconds": sum(latencies),
        }
    summary = {
        "statuses": dict(statuses),
        "total_findings_before_set_union": total_findings,
        "providers": output_providers,
    }
    analysis = run / "analysis"
    analysis.mkdir(exist_ok=True)
    write_json(analysis / "run-health.json", summary)
    lines = ["# Run health", "", f"- Calls: {len(plan['calls'])}", f"- Raw findings: {total_findings}"]
    lines.extend(f"- {status}: {count}" for status, count in sorted(statuses.items()))
    for key, bucket in output_providers.items():
        lines.extend(
            [
                "",
                f"## {key}",
                f"- Input tokens: {bucket['input_tokens']}",
                f"- Output tokens: {bucket['output_tokens']}",
                f"- Recorded cost: ${bucket['cost_usd']:.4f}",
                f"- Mean latency: {bucket['mean_latency_seconds']}",
            ]
        )
    (analysis / "run-health.md").write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(json.dumps(summary, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
