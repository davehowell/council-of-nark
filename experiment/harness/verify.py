from __future__ import annotations

import argparse

from .common import load_json, resolve_run, sha256_file


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify a sealed run's raw-file integrity")
    parser.add_argument("run")
    args = parser.parse_args()
    run = resolve_run(args.run)
    seal_path = run / "seal.json"
    if not seal_path.exists():
        raise SystemExit("Run is not sealed")
    seal = load_json(seal_path)
    errors = []
    expected = {row["path"] for row in seal["files"]}
    for row in seal["files"]:
        path = run / row["path"]
        if not path.exists():
            errors.append(f"missing: {row['path']}")
        elif sha256_file(path) != row["sha256"]:
            errors.append(f"digest mismatch: {row['path']}")
    actual = {
        path.relative_to(run).as_posix()
        for path in (run / "calls").rglob("*")
        if path.is_file()
    }
    unsealed = sorted(actual - expected)
    errors.extend(f"unsealed raw call file: {path}" for path in unsealed)
    if errors:
        print("Integrity verification failed:")
        for error in errors:
            print(f"- {error}")
        return 1
    print(f"Integrity verification passed ({len(expected)} raw files).")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
