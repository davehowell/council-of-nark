from __future__ import annotations

import argparse
from concurrent.futures import Future, ThreadPoolExecutor, wait, FIRST_COMPLETED
import json
from pathlib import Path
import shutil
import subprocess
import threading
import time
from typing import Any

from .adapters import invoke
from .common import (
    ROOT,
    TERMINAL_STATUSES,
    git_blob,
    load_json,
    read_call_record,
    relative,
    resolve_run,
    sha256_bytes,
    sha256_text,
    source_commit,
    successful_response,
    utc_now,
    write_json,
)
from .prompting import dynamic_prompt, packet_text, static_prompt

WORKTREE_LOCK = threading.Lock()


class DetachedWorktree:
    def __init__(self, run: Path, call_id: str, commit: str):
        self.path = ROOT / "experiment" / "worktrees" / run.name / call_id
        self.commit = commit

    def __enter__(self) -> Path:
        self.path.parent.mkdir(parents=True, exist_ok=True)
        with WORKTREE_LOCK:
            if self.path.exists():
                subprocess.run(["git", "worktree", "remove", "--force", str(self.path)], cwd=ROOT, capture_output=True)
                shutil.rmtree(self.path, ignore_errors=True)
            subprocess.run(
                ["git", "worktree", "add", "--detach", "--quiet", str(self.path), self.commit],
                cwd=ROOT,
                check=True,
                capture_output=True,
            )
        actual = subprocess.run(
            ["git", "rev-parse", "HEAD"], cwd=self.path, check=True, text=True, capture_output=True
        ).stdout.strip()
        if actual != self.commit:
            raise RuntimeError("Detached worktree is not at the frozen commit")
        return self.path

    def __exit__(self, exc_type, exc, traceback) -> None:
        with WORKTREE_LOCK:
            subprocess.run(
                ["git", "worktree", "remove", "--force", str(self.path)],
                cwd=ROOT,
                capture_output=True,
                text=True,
            )
            shutil.rmtree(self.path, ignore_errors=True)
            subprocess.run(["git", "worktree", "prune"], cwd=ROOT, capture_output=True)


def verify_manifest(freeze: dict[str, Any]) -> None:
    commit = freeze["source_commit"]
    if source_commit() != commit:
        raise SystemExit(
            "HEAD differs from the frozen source commit. Check out the recorded commit before running."
        )
    for asset in freeze["assets"]:
        actual = sha256_bytes(git_blob(commit, asset["path"]))
        if actual != asset["sha256"]:
            raise SystemExit(f"Frozen asset digest mismatch: {asset['path']}")


def dependency_pairs(
    run: Path, call: dict[str, Any], call_map: dict[str, dict[str, Any]]
) -> list[tuple[dict[str, Any], dict[str, Any]]]:
    return [
        (call_map[call_id], successful_response(run, call_id))
        for call_id in call["depends_on"]
    ]


def execute_call(
    run: Path,
    call: dict[str, Any],
    call_map: dict[str, dict[str, Any]],
    config: dict[str, Any],
    freeze: dict[str, Any],
) -> dict[str, Any]:
    call_id = call["call_id"]
    call_dir = run / "calls" / call_id
    call_dir.mkdir(parents=True, exist_ok=True)
    started_at = utc_now()
    try:
        with DetachedWorktree(run, call_id, freeze["source_commit"]) as worktree:
            kind = call["prompt_spec"]["kind"]
            if kind in {"fuser", "chain"}:
                prompt = dynamic_prompt(worktree, call, dependency_pairs(run, call, call_map))
            else:
                prompt = static_prompt(worktree, call)
            system = (worktree / "experiment/prompts/system.txt").read_text(encoding="utf-8")
            schema_path = worktree / "experiment/schema/findings.schema.json"
            schema = load_json(schema_path)
            request = {
                "schema_version": 1,
                "call_id": call_id,
                "created_at": started_at,
                "source_commit": freeze["source_commit"],
                "provider": call["provider"],
                "phase": call["phase"],
                "packet_sha256": sha256_text(packet_text(worktree, call["packet"])),
                "system_prompt": system,
                "system_prompt_sha256": sha256_text(system),
                "prompt": prompt,
                "prompt_sha256": sha256_text(prompt),
                "schema": schema,
                "schema_sha256": sha256_text(json.dumps(schema, sort_keys=True)),
                "answer_key_in_context": False,
                "isolation": freeze["isolation"],
            }
            write_json(call_dir / "request.json", request)

            max_attempts = config["max_attempts"]
            final = None
            for attempt in range(1, max_attempts + 1):
                attempt_dir = call_dir / "attempts" / str(attempt)
                attempt_dir.mkdir(parents=True, exist_ok=True)
                result = invoke(
                    cwd=worktree,
                    provider=call["provider"],
                    prompt=prompt,
                    system=system,
                    schema=schema,
                    schema_relative="experiment/schema/findings.schema.json",
                    timeout_seconds=config["request_timeout_seconds"],
                )
                (attempt_dir / "stdout.txt").write_text(result.stdout, encoding="utf-8")
                (attempt_dir / "stderr.txt").write_text(result.stderr, encoding="utf-8")
                write_json(
                    attempt_dir / "metadata.json",
                    {
                        "attempt": attempt,
                        "command": result.command_display,
                        "returncode": result.returncode,
                        "latency_seconds": result.latency_seconds,
                        "parse_error": result.parse_error,
                        "usage": result.usage,
                    },
                )
                final = (attempt, result)
                if result.returncode == 0:
                    break
                # Only process/transport failures are retried. A malformed substantive
                # response has return code zero and is preserved as a zero-scoring sample.
                if attempt < max_attempts:
                    backoff = float(config.get("retry_backoff_seconds", 5)) * attempt
                    time.sleep(backoff)

            assert final is not None
            attempt, result = final
            if result.returncode != 0:
                status = "error"
            elif result.parsed is None:
                status = "malformed"
            else:
                status = "success"
            record = {
                "schema_version": 1,
                "call_id": call_id,
                "status": status,
                "started_at": started_at,
                "completed_at": utc_now(),
                "attempts": attempt,
                "returncode": result.returncode,
                "latency_seconds": result.latency_seconds,
                "usage": result.usage,
                "parse_error": result.parse_error,
                "parsed_response": result.parsed,
                "request_sha256": request["prompt_sha256"],
                "source_commit": freeze["source_commit"],
                "worktree": {"fresh": True, "detached": True, "removed_after_call": True},
            }
            write_json(call_dir / "record.json", record)
            return record
    except Exception as error:
        record = {
            "schema_version": 1,
            "call_id": call_id,
            "status": "error",
            "started_at": started_at,
            "completed_at": utc_now(),
            "attempts": 0,
            "error_type": type(error).__name__,
            "error": str(error),
            "source_commit": freeze["source_commit"],
        }
        write_json(call_dir / "record.json", record)
        return record


def blocked_record(run: Path, call: dict[str, Any], failed: list[str]) -> None:
    write_json(
        run / "calls" / call["call_id"] / "record.json",
        {
            "schema_version": 1,
            "call_id": call["call_id"],
            "status": "blocked",
            "completed_at": utc_now(),
            "failed_dependencies": failed,
        },
    )


def ordered_pending_call_ids(run: Path, plan: dict[str, Any]) -> list[str]:
    """Return unfinished calls in the seeded order recorded by plan.json."""
    return [
        call["call_id"]
        for call in plan["calls"]
        if (read_call_record(run, call["call_id"]) or {}).get("status") not in TERMINAL_STATUSES
    ]


def main() -> int:
    parser = argparse.ArgumentParser(description="Execute a frozen experiment plan")
    parser.add_argument("run")
    parser.add_argument("--jobs", type=int)
    args = parser.parse_args()
    run = resolve_run(args.run)
    if (run / "seal.json").exists():
        raise SystemExit("Run is sealed and cannot be modified")

    freeze = load_json(run / "freeze.json")
    config = load_json(run / "config.json")
    plan = load_json(run / "plan.json")
    verify_manifest(freeze)
    call_map = {call["call_id"]: call for call in plan["calls"]}
    # Preserve the seeded order stored in plan.json. A set would make execution
    # order depend on Python hash iteration and reintroduce an avoidable time confound.
    pending = ordered_pending_call_ids(run, plan)
    jobs = args.jobs or config["max_concurrency"]
    running: dict[Future[dict[str, Any]], str] = {}

    with ThreadPoolExecutor(max_workers=jobs) as pool:
        while pending or running:
            made_progress = False
            for call_id in pending.copy():
                if len(running) >= jobs:
                    break
                call = call_map[call_id]
                dependency_records = {
                    dep: read_call_record(run, dep) for dep in call["depends_on"]
                }
                if any(record is None for record in dependency_records.values()):
                    continue
                failed = [
                    dep for dep, record in dependency_records.items() if record["status"] != "success"
                ]
                pending.remove(call_id)
                made_progress = True
                if failed:
                    blocked_record(run, call, failed)
                    print(f"{call_id} blocked by {', '.join(failed)}", flush=True)
                    continue
                future = pool.submit(execute_call, run, call, call_map, config, freeze)
                running[future] = call_id

            if running:
                done, _ = wait(running, return_when=FIRST_COMPLETED)
                for future in done:
                    call_id = running.pop(future)
                    record = future.result()
                    print(f"{call_id} {record['status']}", flush=True)
                    made_progress = True
            if not made_progress and pending:
                raise SystemExit("Plan contains a dependency cycle or missing call")

    statuses: dict[str, int] = {}
    for call_id in call_map:
        status = (read_call_record(run, call_id) or {}).get("status", "missing")
        statuses[status] = statuses.get(status, 0) + 1
    write_json(run / "run-summary.json", {"completed_at": utc_now(), "statuses": statuses})
    print(json.dumps(statuses, sort_keys=True))
    return 0 if statuses.get("missing", 0) == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
