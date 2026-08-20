from __future__ import annotations

import argparse
import csv
import json
from pathlib import Path
import random
import shutil
from typing import Any

from .common import load_json, opaque_id, read_call_record, resolve_run, write_json


def findings_for_set(run: Path, output_set: dict[str, Any]) -> list[dict[str, Any]]:
    findings: list[dict[str, Any]] = []
    for call_id in output_set["call_ids"]:
        record = read_call_record(run, call_id)
        if not record or record.get("status") != "success":
            continue
        for finding in record["parsed_response"]["findings"]:
            findings.append(dict(finding))
    return findings


def write_jsonl(path: Path, rows: list[dict[str, Any]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(
        "".join(json.dumps(row, sort_keys=True, ensure_ascii=False) + "\n" for row in rows),
        encoding="utf-8",
    )


def main() -> int:
    parser = argparse.ArgumentParser(description="Create an arm-blinded human-rating bundle")
    parser.add_argument("run")
    args = parser.parse_args()
    run = resolve_run(args.run)
    plan = load_json(run / "plan.json")
    blinded = run / "blinded"
    if blinded.exists():
        shutil.rmtree(blinded)
    blinded.mkdir()

    seed = plan["seed"] + ":blinding"
    item_rows: list[dict[str, Any]] = []
    set_rows: list[dict[str, Any]] = []
    unblind: dict[str, Any] = {"sets": {}, "items": {}}
    for output_set in plan["output_sets"]:
        item_ids = []
        for index, finding in enumerate(findings_for_set(run, output_set)):
            item_id = "i-" + opaque_id(seed, output_set["set_id"], index)
            item_ids.append(item_id)
            item_rows.append(
                {
                    "item_id": item_id,
                    "set_id": output_set["set_id"],
                    "packet": output_set["packet"],
                    "finding": finding,
                }
            )
            unblind["items"][item_id] = {
                "set_id": output_set["set_id"],
                "source_index": index,
            }
        set_rows.append(
            {
                "set_id": output_set["set_id"],
                "packet": output_set["packet"],
                "item_ids": item_ids,
            }
        )
        unblind["sets"][output_set["set_id"]] = output_set

    random.Random(seed).shuffle(item_rows)
    random.Random(seed + ":sets").shuffle(set_rows)
    write_jsonl(blinded / "findings.jsonl", item_rows)
    write_jsonl(blinded / "sets.jsonl", set_rows)
    write_json(run / "private" / "unblind.json", unblind)

    answer_dir = blinded / "answer-keys"
    answer_dir.mkdir()
    packets = sorted({row["packet"] for row in set_rows})
    for packet in packets:
        shutil.copyfile(
            Path(__file__).resolve().parents[2] / "experiment" / "scenarios" / packet / "answer-key.md",
            answer_dir / f"{packet}.md",
        )

    with (blinded / "rating-template.csv").open("w", newline="", encoding="utf-8") as handle:
        writer = csv.writer(handle)
        writer.writerow(
            [
                "rater",
                "item_id",
                "defect_id",
                "false_positive_cluster",
                "material",
                "confidence",
                "notes",
            ]
        )
        for row in item_rows:
            writer.writerow(["", row["item_id"], "", "", "", "", ""])

    (blinded / "RUNSHEET.md").write_text(
        "# Blinded rating runsheet\n\n"
        "1. Do not open `plan.json`, `calls/`, or `private/unblind.json`.\n"
        "2. Open `findings.jsonl` and the answer key for each item's `packet`.\n"
        "3. Enter one row per item in a copy of `rating-template.csv`.\n"
        "4. Set `defect_id` to one matching planted ID, or `NONE`.\n"
        "5. For `NONE`, assign a short `false_positive_cluster`; reuse it for the same claimed mechanism in this set.\n"
        "6. Set `material` to `true` only when the packet supports the claim and its consequence.\n"
        "7. Use `confidence` = `high`, `medium`, or `low`; explain ambiguous matches in `notes`.\n"
        "8. Work independently. Do not compare ratings until both raters finish.\n"
        "9. Resolve disagreements without revealing arm, model, or provider metadata.\n"
        "10. Save the adjudicated file and run the scoring recipe.\n\n"
        "Different wording can map to the same defect. An unmatched material claim is a false positive. "
        "A style preference with no planted mechanism maps to `NONE`.\n",
        encoding="utf-8",
    )
    print(f"Created {len(item_rows)} blinded items across {len(set_rows)} output sets.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
