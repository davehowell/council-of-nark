from __future__ import annotations

from typing import Any

CONFIDENCE = {"high", "medium", "low"}
REQUIRED_FINDING = {"location", "claim", "consequence", "fix", "confidence"}


def findings_response(value: Any) -> list[str]:
    errors: list[str] = []
    if not isinstance(value, dict):
        return ["response is not an object"]
    if set(value) != {"findings"}:
        errors.append("root must contain only 'findings'")
    findings = value.get("findings")
    if not isinstance(findings, list):
        return errors + ["findings is not an array"]
    if len(findings) > 8:
        errors.append("findings contains more than 8 items")
    for index, finding in enumerate(findings):
        prefix = f"findings[{index}]"
        if not isinstance(finding, dict):
            errors.append(f"{prefix} is not an object")
            continue
        keys = set(finding)
        if not REQUIRED_FINDING.issubset(keys):
            errors.append(f"{prefix} lacks required fields")
        if keys - (REQUIRED_FINDING | {"raised_by"}):
            errors.append(f"{prefix} has unknown fields")
        for key in REQUIRED_FINDING - {"confidence"}:
            if not isinstance(finding.get(key), str) or not finding.get(key, "").strip():
                errors.append(f"{prefix}.{key} must be a non-empty string")
        if finding.get("confidence") not in CONFIDENCE:
            errors.append(f"{prefix}.confidence is invalid")
        raised_by = finding.get("raised_by")
        if raised_by is not None and (
            not isinstance(raised_by, list)
            or not all(isinstance(item, str) and item for item in raised_by)
            or len(set(raised_by)) != len(raised_by)
        ):
            errors.append(f"{prefix}.raised_by is invalid")
    return errors


def judgement_response(value: Any, expected_ids: set[str]) -> list[str]:
    if not isinstance(value, dict) or set(value) != {"judgements"}:
        return ["root must contain only 'judgements'"]
    rows = value["judgements"]
    if not isinstance(rows, list):
        return ["judgements is not an array"]
    errors: list[str] = []
    seen: set[str] = set()
    for index, row in enumerate(rows):
        if not isinstance(row, dict):
            errors.append(f"judgements[{index}] is not an object")
            continue
        required = {"item_id", "defect_id", "material", "confidence", "rationale"}
        if set(row) != required:
            errors.append(f"judgements[{index}] fields do not match schema")
            continue
        item_id = row["item_id"]
        if item_id not in expected_ids or item_id in seen:
            errors.append(f"judgements[{index}] has unexpected or duplicate item_id")
        seen.add(item_id)
        if row["defect_id"] is not None and not isinstance(row["defect_id"], str):
            errors.append(f"judgements[{index}].defect_id is invalid")
        if not isinstance(row["material"], bool):
            errors.append(f"judgements[{index}].material is invalid")
        if row["confidence"] not in CONFIDENCE:
            errors.append(f"judgements[{index}].confidence is invalid")
        if not isinstance(row["rationale"], str):
            errors.append(f"judgements[{index}].rationale is invalid")
    if seen != expected_ids:
        errors.append("judgements do not cover every expected item ID")
    return errors
