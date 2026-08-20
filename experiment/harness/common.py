from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
EXPERIMENT = ROOT / "experiment"
TERMINAL_STATUSES = {"success", "malformed", "error", "blocked"}


def utc_now() -> str:
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def slug_timestamp() -> str:
    return datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def write_json(path: Path, value: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    payload = json.dumps(value, indent=2, sort_keys=True, ensure_ascii=False) + "\n"
    atomic_write(path, payload)


def atomic_write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    fd, temporary = tempfile.mkstemp(prefix=path.name + ".", dir=path.parent)
    try:
        with os.fdopen(fd, "w", encoding="utf-8") as handle:
            handle.write(content)
            handle.flush()
            os.fsync(handle.fileno())
        os.replace(temporary, path)
    finally:
        if os.path.exists(temporary):
            os.unlink(temporary)


def sha256_bytes(content: bytes) -> str:
    return hashlib.sha256(content).hexdigest()


def sha256_text(content: str) -> str:
    return sha256_bytes(content.encode("utf-8"))


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def opaque_id(seed: str, *parts: object, length: int = 16) -> str:
    material = "\x1f".join([seed, *(str(part) for part in parts)])
    return hashlib.sha256(material.encode("utf-8")).hexdigest()[:length]


def git(*args: str, cwd: Path = ROOT, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["git", *args], cwd=cwd, check=check, text=True, capture_output=True
    )


def source_commit() -> str:
    return git("rev-parse", "HEAD").stdout.strip()


def require_clean_tree() -> None:
    status = git("status", "--porcelain=v1").stdout.strip()
    if status:
        raise SystemExit("Refusing to freeze a dirty tree. Commit or discard every source change first.")


def tracked_files_at(commit: str, prefixes: tuple[str, ...]) -> list[str]:
    result = git("ls-tree", "-r", "--name-only", commit, "--", *prefixes)
    return [line for line in result.stdout.splitlines() if line]


def git_blob(commit: str, path: str) -> bytes:
    result = subprocess.run(
        ["git", "show", f"{commit}:{path}"], cwd=ROOT, check=True, capture_output=True
    )
    return result.stdout


def resolve_run(value: str) -> Path:
    path = Path(value)
    if not path.is_absolute():
        path = ROOT / path
    path = path.resolve()
    runs = (EXPERIMENT / "runs").resolve()
    if runs not in path.parents:
        raise SystemExit(f"Run path must be below {runs.relative_to(ROOT)}")
    if not path.exists():
        raise SystemExit(f"Run does not exist: {path}")
    return path


def relative(path: Path) -> str:
    try:
        return path.resolve().relative_to(ROOT.resolve()).as_posix()
    except ValueError:
        return path.as_posix()


def read_call_record(run: Path, call_id: str) -> dict[str, Any] | None:
    path = run / "calls" / call_id / "record.json"
    return load_json(path) if path.exists() else None


def successful_response(run: Path, call_id: str) -> dict[str, Any]:
    record = read_call_record(run, call_id)
    if not record or record.get("status") != "success":
        raise RuntimeError(f"Dependency {call_id} has no successful response")
    return record["parsed_response"]
