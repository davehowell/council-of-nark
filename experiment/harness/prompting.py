from __future__ import annotations

import json
from pathlib import Path
import re
from typing import Any

from .common import load_json, sha256_text

DEFECT_ID = re.compile(r"\b[A-Z]{2}-\d{2}\b")


def text(root: Path, relative: str) -> str:
    return (root / relative).read_text(encoding="utf-8")


def packet_text(root: Path, packet: str) -> str:
    return text(root, f"experiment/scenarios/{packet}/review-packet.md")


def answer_text(root: Path, packet: str) -> str:
    return text(root, f"experiment/scenarios/{packet}/answer-key.md")


def specialists(root: Path) -> dict[str, dict[str, str]]:
    rows = load_json(root / "experiment/prompts/specialists.json")
    return {row["id"]: row for row in rows}


def replace(template: str, **values: str) -> str:
    rendered = template
    for key, value in values.items():
        marker = "{{" + key + "}}"
        if marker not in rendered:
            raise ValueError(f"Template does not contain {marker}")
        rendered = rendered.replace(marker, value)
    leftovers = re.findall(r"{{[A-Z_]+}}", rendered)
    if leftovers:
        raise ValueError(f"Unresolved template markers: {leftovers}")
    return rendered


def assert_answer_key_absent(prompt: str, answer_key: str) -> None:
    planted_ids = set(DEFECT_ID.findall(answer_key))
    if "Private answer key" in prompt or any(defect_id in prompt for defect_id in planted_ids):
        raise ValueError("Answer-key marker or planted defect ID leaked into reviewer prompt")
    cells: list[str] = []
    for line in answer_key.splitlines():
        if not line.startswith("|") or line.startswith("|---") or "Planted defect" in line:
            continue
        parts = [part.strip() for part in line.strip("|").split("|")]
        cells.extend(part for part in parts[2:] if len(part) >= 48)
    for cell in cells:
        if cell in prompt:
            raise ValueError("Answer-key claim leaked verbatim into reviewer prompt")


def review_contract(root: Path, packet: str) -> str:
    contract = text(root, "experiment/prompts/review-contract.txt")
    return replace(contract, REVIEW_PACKET=packet_text(root, packet))


def specialist_prompt(root: Path, packet: str, role: str, wrapper: str) -> str:
    row = specialists(root)[role]
    style = row[f"{wrapper}_wrapper"]
    intro = replace(
        text(root, "experiment/prompts/specialist-intro.txt"),
        LENS_KERNEL=row["kernel"],
    )
    prompt = f"{style}\n\n{intro}\n{review_contract(root, packet)}"
    assert_answer_key_absent(prompt, answer_text(root, packet))
    return prompt


def static_prompt(root: Path, call: dict[str, Any]) -> str:
    packet = call["packet"]
    kind = call["prompt_spec"]["kind"]
    if kind == "generic":
        prompt = text(root, "experiment/prompts/generic-kernel.txt") + "\n" + review_contract(root, packet)
    elif kind == "omnibus":
        wrappers = load_json(root / "experiment/prompts/omnibus-wrappers.json")
        prompt = (
            wrappers[call["prompt_spec"]["wrapper"]]
            + "\n\n"
            + text(root, "experiment/prompts/omnibus-kernel.txt")
            + "\n"
            + review_contract(root, packet)
        )
    elif kind == "specialist":
        prompt = specialist_prompt(
            root, packet, call["prompt_spec"]["role"], call["prompt_spec"]["wrapper"]
        )
    else:
        raise ValueError(f"Prompt kind {kind!r} requires dependency output")
    assert_answer_key_absent(prompt, answer_text(root, packet))
    return prompt


def with_reviewer_ids(dependencies: list[tuple[dict[str, Any], dict[str, Any]]]) -> list[dict[str, Any]]:
    payload = []
    for call, response in dependencies:
        findings = []
        for finding in response["findings"]:
            finding = dict(finding)
            finding.setdefault("raised_by", [call["reviewer_id"]])
            findings.append(finding)
        payload.append({"reviewer_id": call["reviewer_id"], "findings": findings})
    return payload


def dynamic_prompt(
    root: Path,
    call: dict[str, Any],
    dependencies: list[tuple[dict[str, Any], dict[str, Any]]],
) -> str:
    packet = call["packet"]
    kind = call["prompt_spec"]["kind"]
    if kind == "fuser":
        template = text(root, "experiment/prompts/fuser.txt")
        prompt = replace(
            template,
            REVIEW_PACKET=packet_text(root, packet),
            REVIEW_FINDINGS=json.dumps(with_reviewer_ids(dependencies), ensure_ascii=False),
        )
    elif kind == "chain":
        row = specialists(root)[call["prompt_spec"]["role"]]
        style = row[f"{call['prompt_spec']['wrapper']}_wrapper"]
        intro = replace(
            text(root, "experiment/prompts/specialist-intro.txt"),
            LENS_KERNEL=row["kernel"],
        )
        ledger = with_reviewer_ids(dependencies)
        template = text(root, "experiment/prompts/chain-contract.txt")
        prompt = f"{style}\n\n{intro}\n" + replace(
            template,
            REVIEW_PACKET=packet_text(root, packet),
            PRIOR_FINDINGS=json.dumps(ledger, ensure_ascii=False),
        )
    else:
        raise ValueError(f"Prompt kind {kind!r} is not dynamic")
    assert_answer_key_absent(prompt, answer_text(root, packet))
    return prompt


def prompt_digest(prompt: str) -> str:
    return sha256_text(prompt)
