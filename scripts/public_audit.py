#!/usr/bin/env python3
"""Fail when public repository material contains common secret or privacy indicators.

Add project-specific forbidden terms to the ignored `.public-audit-forbidden` file
(one term per line) or the comma-separated PUBLIC_AUDIT_FORBIDDEN environment
variable. This scanner supplements manual review; it is not a proof of safety.
"""

from __future__ import annotations

import argparse
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
TEXT_SUFFIXES = {
    "", ".css", ".csv", ".html", ".ini", ".js", ".json", ".md", ".mjs",
    ".py", ".sh", ".sql", ".toml", ".ts", ".txt", ".yaml", ".yml",
}
SKIP_PARTS = {".git", "node_modules", "experiment/runs", "experiment/worktrees"}
PLACEHOLDER_WORDS = {
    "example", "fake", "placeholder", "redacted", "dummy", "test", "changeme",
    "do-not-use", "do_not_use", "<", "{{", "${",
}

PATTERNS: list[tuple[str, re.Pattern[str]]] = [
    ("private key block", re.compile(r"BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY")),
    ("AWS access key", re.compile(r"\b(?:AKIA|ASIA)[A-Z0-9]{16}\b")),
    ("GitHub token", re.compile(r"\bgh[pousr]_[A-Za-z0-9_]{30,}\b")),
    ("Google API key", re.compile(r"\bAIza[0-9A-Za-z_-]{30,}\b")),
    ("OpenAI-style key", re.compile(r"\bsk-[A-Za-z0-9_-]{20,}\b")),
    ("credential in URL", re.compile(r"https?://[^\s/:]+:[^\s/@]+@", re.I)),
    ("macOS home path", re.compile(r"/Users/[A-Za-z0-9._-]+/")),
    ("Linux home path", re.compile(r"/home/[A-Za-z0-9._-]+/")),
    ("private hostname", re.compile(r"\b[A-Za-z0-9.-]+\.(?:corp|internal|private)\b", re.I)),
    ("private IPv4 address", re.compile(r"\b(?:10\.\d{1,3}\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3}|172\.(?:1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3})\b")),
]

ASSIGNMENT = re.compile(
    r"(?i)\b(password|passwd|secret|api[_-]?key|access[_-]?token|client[_-]?secret)"
    r"\s*[:=]\s*['\"]([^'\"\n]{8,})['\"]"
)


def candidate_paths() -> list[Path]:
    command = ["git", "ls-files", "--cached", "--others", "--exclude-standard", "-z"]
    result = subprocess.run(command, cwd=ROOT, check=True, capture_output=True)
    paths = []
    for raw in result.stdout.split(b"\0"):
        if not raw:
            continue
        path = ROOT / os.fsdecode(raw)
        relative = path.relative_to(ROOT).as_posix()
        if any(relative == part or relative.startswith(part + "/") for part in SKIP_PARTS):
            continue
        if path.is_file():
            paths.append(path)
    return sorted(paths)


def extra_forbidden() -> list[str]:
    terms: list[str] = []
    local = ROOT / ".public-audit-forbidden"
    if local.exists():
        terms.extend(
            line.strip() for line in local.read_text(encoding="utf-8").splitlines()
            if line.strip() and not line.lstrip().startswith("#")
        )
    terms.extend(term.strip() for term in os.getenv("PUBLIC_AUDIT_FORBIDDEN", "").split(",") if term.strip())
    return terms


def read_public_text(path: Path) -> str | None:
    if path.suffix.lower() == ".pdf":
        tool = shutil.which("pdftotext")
        if tool is None:
            return None
        result = subprocess.run([tool, str(path), "-"], check=True, capture_output=True)
        return result.stdout.decode("utf-8", errors="replace")
    if path.suffix.lower() not in TEXT_SUFFIXES:
        return None
    try:
        return path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        return None


def line_number(text: str, offset: int) -> int:
    return text.count("\n", 0, offset) + 1


def looks_like_placeholder(value: str) -> bool:
    lowered = value.lower()
    return any(word in lowered for word in PLACEHOLDER_WORDS)


def scan(path: Path, text: str, forbidden: list[str]) -> list[str]:
    relative = path.relative_to(ROOT).as_posix()
    findings: list[str] = []
    for label, pattern in PATTERNS:
        for match in pattern.finditer(text):
            findings.append(f"{relative}:{line_number(text, match.start())}: {label}")
    for match in ASSIGNMENT.finditer(text):
        if not looks_like_placeholder(match.group(2)):
            findings.append(
                f"{relative}:{line_number(text, match.start())}: possible credential assignment"
            )
    lowered = text.casefold()
    lowered_path = relative.casefold()
    for term in forbidden:
        needle = term.casefold()
        if needle in lowered_path:
            findings.append(f"{relative}: forbidden local term in path")
        start = 0
        while True:
            offset = lowered.find(needle, start)
            if offset < 0:
                break
            findings.append(
                f"{relative}:{line_number(text, offset)}: forbidden local term"
            )
            start = offset + max(1, len(needle))
    return findings


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--require-pdf-text",
        action="store_true",
        help="fail if pdftotext is unavailable while PDFs are present",
    )
    args = parser.parse_args()

    forbidden = extra_forbidden()
    findings: list[str] = []
    skipped_pdfs: list[str] = []
    scanned = 0

    for path in candidate_paths():
        text = read_public_text(path)
        if text is None:
            if path.suffix.lower() == ".pdf" and shutil.which("pdftotext") is None:
                skipped_pdfs.append(path.relative_to(ROOT).as_posix())
            continue
        scanned += 1
        findings.extend(scan(path, text, forbidden))

    if skipped_pdfs:
        message = "pdftotext unavailable; did not inspect PDF text: " + ", ".join(skipped_pdfs)
        if args.require_pdf_text:
            findings.append(message)
        else:
            print("warning: " + message, file=sys.stderr)

    if findings:
        print("Public-release audit failed:", file=sys.stderr)
        for finding in findings:
            print(f"- {finding}", file=sys.stderr)
        return 1

    print(f"Public-release audit passed ({scanned} text/PDF files, {len(forbidden)} local terms).")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
