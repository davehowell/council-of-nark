from __future__ import annotations

import argparse
import os
from pathlib import Path

from .common import (
    TERMINAL_STATUSES,
    load_json,
    read_call_record,
    relative,
    resolve_run,
    sha256_file,
    utc_now,
    write_json,
)

RAW_NAMES = {"freeze.json", "config.json", "plan.json", "run-summary.json"}


def raw_files(run: Path) -> list[Path]:
    files = [run / name for name in RAW_NAMES if (run / name).exists()]
    calls = run / "calls"
    if calls.exists():
        files.extend(path for path in calls.rglob("*") if path.is_file())
    return sorted(files)


def main() -> int:
    parser = argparse.ArgumentParser(description="Seal immutable raw requests and responses by digest")
    parser.add_argument("run")
    args = parser.parse_args()
    run = resolve_run(args.run)
    if (run / "seal.json").exists():
        raise SystemExit("Run is already sealed")
    plan = load_json(run / "plan.json")
    incomplete = []
    for call in plan["calls"]:
        record = read_call_record(run, call["call_id"])
        if not record or record.get("status") not in TERMINAL_STATUSES:
            incomplete.append(call["call_id"])
    if incomplete:
        raise SystemExit(f"Cannot seal; incomplete calls: {', '.join(incomplete)}")
    manifest = [
        {"path": path.relative_to(run).as_posix(), "sha256": sha256_file(path)}
        for path in raw_files(run)
    ]
    write_json(
        run / "seal.json",
        {
            "schema_version": 1,
            "sealed_at": utc_now(),
            "scope": "freeze/config/plan/run-summary and calls/**; derived blinded and analysis files are excluded",
            "files": manifest,
        },
    )
    for path in raw_files(run):
        path.chmod(path.stat().st_mode & ~0o222)
    print(f"Sealed {len(manifest)} raw files in {relative(run)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
