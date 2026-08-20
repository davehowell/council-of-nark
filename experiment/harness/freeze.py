from __future__ import annotations

import argparse
import json
from pathlib import Path
import subprocess

from .common import (
    EXPERIMENT,
    ROOT,
    git_blob,
    relative,
    require_clean_tree,
    sha256_bytes,
    slug_timestamp,
    source_commit,
    tracked_files_at,
    utc_now,
    write_json,
)


def main() -> int:
    parser = argparse.ArgumentParser(description="Freeze committed experiment assets into a run manifest")
    parser.add_argument("--config", required=True, help="tracked config path, relative to repository root")
    args = parser.parse_args()

    require_clean_tree()
    audit = subprocess.run(
        ["python3", "scripts/public_audit.py"], cwd=ROOT, text=True, capture_output=True
    )
    if audit.returncode:
        raise SystemExit(audit.stdout + audit.stderr)

    commit = source_commit()
    config_path = (ROOT / args.config).resolve()
    try:
        config_relative = config_path.relative_to(ROOT).as_posix()
    except ValueError as error:
        raise SystemExit("Config must be inside the repository") from error
    config_bytes = git_blob(commit, config_relative)
    config = json.loads(config_bytes)

    files = tracked_files_at(commit, ("experiment", "justfile", "scripts/public_audit.py"))
    manifest_files = [
        {"path": path, "sha256": sha256_bytes(git_blob(commit, path))}
        for path in files
        if not path.startswith("experiment/results/")
    ]
    run_name = f"{slug_timestamp()}-{config['label']}-{commit[:10]}"
    run = EXPERIMENT / "runs" / run_name
    if run.exists():
        raise SystemExit(f"Run already exists: {relative(run)}")
    run.mkdir(parents=True)

    manifest = {
        "schema_version": 1,
        "created_at": utc_now(),
        "source_commit": commit,
        "config_source": config_relative,
        "config_sha256": sha256_bytes(config_bytes),
        "assets": manifest_files,
        "isolation": {
            "git_worktree": "one clean detached worktree per model call",
            "process": "one fresh process per model call",
            "session_persistence": false,
            "tools": false,
            "project_context": false,
            "provider_side_state": "unobservable; recorded as a validity limitation"
        },
        "public_audit": {
            "command": "python3 scripts/public_audit.py",
            "passed": true,
            "stdout": audit.stdout.strip(),
            "stderr": audit.stderr.strip(),
        },
    }
    write_json(run / "freeze.json", manifest)
    (run / "config.json").write_bytes(config_bytes)
    print(relative(run))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
