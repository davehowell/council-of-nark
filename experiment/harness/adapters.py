from __future__ import annotations

from dataclasses import dataclass
import json
import os
from pathlib import Path
import re
import subprocess
import time
from typing import Any

from .validate import findings_response, judgement_response

AUTH_ENV = {
    "ANTHROPIC_API_KEY",
    "ANTHROPIC_OAUTH_TOKEN",
    "OPENAI_API_KEY",
    "GEMINI_API_KEY",
    "GOOGLE_API_KEY",
    "GOOGLE_APPLICATION_CREDENTIALS",
    "GOOGLE_CLOUD_PROJECT",
    "GOOGLE_CLOUD_LOCATION",
    "AWS_PROFILE",
    "AWS_REGION",
}
BASE_ENV = {"HOME", "PATH", "TMPDIR", "LANG", "LC_ALL", "USER", "SHELL", "TERM"}


@dataclass
class Invocation:
    command: list[str]
    command_display: list[str]
    returncode: int
    stdout: str
    stderr: str
    latency_seconds: float
    parsed: Any | None
    parse_error: str | None
    usage: dict[str, Any]


def clean_environment() -> dict[str, str]:
    env = {key: value for key, value in os.environ.items() if key in BASE_ENV or key in AUTH_ENV}
    env.update(
        {
            "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
            "DISABLE_AUTOUPDATER": "1",
            "PI_SKIP_VERSION_CHECK": "1",
            "PI_TELEMETRY": "0",
            "NO_COLOR": "1",
        }
    )
    return env


def command_for(
    provider: dict[str, Any],
    prompt: str,
    system: str,
    schema: dict[str, Any],
    schema_relative: str,
) -> tuple[list[str], list[str]]:
    adapter = provider["adapter"]
    model = provider["model"]
    effort = provider.get("effort", "low")
    # Claude Code's structured-output validator accepts the supported schema subset
    # but rejects the draft-2020 meta-schema URI. Keep `$schema` in the frozen file
    # and omit only that annotation from the CLI argument.
    cli_schema = {key: value for key, value in schema.items() if key != "$schema"}
    schema_text = json.dumps(cli_schema, separators=(",", ":"))
    if adapter == "claude":
        command = [
            "claude",
            "--print",
            "--safe-mode",
            "--disable-slash-commands",
            "--no-chrome",
            "--no-session-persistence",
            "--permission-mode",
            "dontAsk",
            "--tools",
            "",
            "--model",
            model,
            "--effort",
            effort,
            "--system-prompt",
            system,
            "--json-schema",
            schema_text,
            "--output-format",
            "json",
            prompt,
        ]
    elif adapter == "agy":
        command = [
            "agy",
            "--print",
            prompt,
            "--model",
            model,
            "--effort",
            effort,
            "--output-format",
            "json",
            "--json-schema",
            schema_text,
            "--disable-slash-commands",
            "--sandbox",
        ]
    elif adapter == "pi":
        command = [
            "pi",
            "--print",
            "--no-session",
            "--no-tools",
            "--no-extensions",
            "--no-skills",
            "--no-prompt-templates",
            "--no-themes",
            "--no-context-files",
            "--no-approve",
            "--model",
            model,
            "--thinking",
            effort,
            "--system-prompt",
            system,
            "--mode",
            "json",
            prompt,
        ]
    elif adapter == "mock":
        command = []
    else:
        raise ValueError(f"Unknown adapter: {adapter}")

    display = []
    for item in command:
        if item == prompt:
            display.append("<PROMPT>")
        elif item == system:
            display.append("<SYSTEM_PROMPT>")
        elif item == schema_text:
            display.append("<JSON_SCHEMA>")
        else:
            display.append(item)
    return command, display


def raw_json_values(text: str) -> list[Any]:
    values: list[Any] = []
    stripped = text.strip()
    if stripped.startswith("```"):
        stripped = re.sub(r"^```(?:json)?\s*|\s*```$", "", stripped, flags=re.I | re.S)
    for candidate in [stripped, *text.splitlines()]:
        try:
            values.append(json.loads(candidate))
        except (json.JSONDecodeError, TypeError):
            pass
    decoder = json.JSONDecoder()
    for match in re.finditer(r"[\[{]", text):
        try:
            value, _ = decoder.raw_decode(text[match.start() :])
            values.append(value)
        except json.JSONDecodeError:
            continue
    return values


def nested_candidates(value: Any) -> list[Any]:
    values = [value]
    if isinstance(value, dict):
        for key in ("result", "structured_output", "output", "response", "message"):
            nested = value.get(key)
            if isinstance(nested, str):
                values.extend(raw_json_values(nested))
            elif isinstance(nested, (dict, list)):
                values.extend(nested_candidates(nested))
        content = value.get("content")
        if isinstance(content, list):
            for block in content:
                if isinstance(block, dict) and isinstance(block.get("text"), str):
                    values.extend(raw_json_values(block["text"]))
    elif isinstance(value, list):
        for item in value:
            values.extend(nested_candidates(item))
    return values


def extract(text: str, expected_root: str) -> tuple[Any | None, str | None]:
    candidates: list[Any] = []
    for value in raw_json_values(text):
        candidates.extend(nested_candidates(value))
    for value in reversed(candidates):
        if isinstance(value, dict) and expected_root in value:
            return value, None

    # Pi JSON mode can stream text deltas without a final text block.
    deltas: list[str] = []
    for value in raw_json_values(text):
        if isinstance(value, dict):
            event = value.get("assistantMessageEvent") or value.get("event")
            if isinstance(event, dict) and event.get("type") == "text_delta":
                delta = event.get("delta") or event.get("text")
                if isinstance(delta, str):
                    deltas.append(delta)
    if deltas:
        for value in raw_json_values("".join(deltas)):
            if isinstance(value, dict) and expected_root in value:
                return value, None
    return None, f"No JSON object with root key {expected_root!r} found"


def ndjson_events(text: str) -> list[dict[str, Any]]:
    events: list[dict[str, Any]] = []
    for line in text.splitlines():
        try:
            value = json.loads(line)
        except json.JSONDecodeError:
            continue
        if isinstance(value, dict):
            events.append(value)
    return events


def message_text(message: Any) -> list[str]:
    if not isinstance(message, dict) or message.get("role") != "assistant":
        return []
    content = message.get("content")
    if not isinstance(content, list):
        return []
    return [
        block["text"]
        for block in content
        if isinstance(block, dict) and block.get("type") == "text" and isinstance(block.get("text"), str)
    ]


def extract_pi_assistant(text: str, expected_root: str) -> tuple[Any | None, str | None]:
    """Parse only Pi assistant events; never inspect echoed user prompt content."""
    completed: list[str] = []
    deltas: list[str] = []
    for value in ndjson_events(text):
        if value.get("type") in {"message_end", "turn_end"}:
            completed.extend(message_text(value.get("message")))
        if value.get("type") == "agent_end" and isinstance(value.get("messages"), list):
            for message in value["messages"]:
                completed.extend(message_text(message))
        event = value.get("assistantMessageEvent") or value.get("event")
        if isinstance(event, dict):
            if event.get("type") == "text_end" and isinstance(event.get("content"), str):
                completed.append(event["content"])
            elif event.get("type") == "text_delta":
                delta = event.get("delta") or event.get("text")
                if isinstance(delta, str):
                    deltas.append(delta)
    if deltas:
        completed.append("".join(deltas))
    for candidate in reversed(completed):
        for value in reversed(raw_json_values(candidate)):
            if isinstance(value, dict) and expected_root in value:
                return value, None
    return None, f"No assistant JSON object with root key {expected_root!r} found"


def pi_provider_error(text: str) -> str | None:
    for value in ndjson_events(text):
        message = value.get("message")
        if isinstance(message, dict) and message.get("role") == "assistant":
            if message.get("stopReason") == "error" or message.get("errorMessage"):
                return str(message.get("errorMessage") or "Pi provider error")
        if value.get("type") == "error":
            return str(value.get("error") or value.get("message") or "Pi provider error")
    return None


def usage_from(text: str) -> dict[str, Any]:
    result: dict[str, Any] = {}
    values = raw_json_values(text)
    for value in values:
        if not isinstance(value, dict):
            continue
        for key in ("usage", "modelUsage"):
            if isinstance(value.get(key), dict):
                result[key] = value[key]
        for key in ("total_cost_usd", "cost_usd", "duration_ms", "duration_api_ms"):
            if isinstance(value.get(key), (int, float)):
                result[key] = value[key]
        message = value.get("message")
        if isinstance(message, dict) and isinstance(message.get("usage"), dict):
            result["usage"] = message["usage"]
    return result


def invoke(
    *,
    cwd: Path,
    provider: dict[str, Any],
    prompt: str,
    system: str,
    schema: dict[str, Any],
    schema_relative: str,
    timeout_seconds: int,
    expected_root: str = "findings",
    expected_ids: set[str] | None = None,
) -> Invocation:
    command, display = command_for(provider, prompt, system, schema, schema_relative)
    started = time.monotonic()
    if provider["adapter"] == "mock":
        stdout, stderr, returncode = '{"findings":[]}\n', "", 0
    else:
        try:
            process = subprocess.run(
                command,
                cwd=cwd,
                env=clean_environment(),
                stdin=subprocess.DEVNULL,
                capture_output=True,
                text=True,
                timeout=timeout_seconds,
            )
            stdout, stderr, returncode = process.stdout, process.stderr, process.returncode
        except subprocess.TimeoutExpired as error:
            stdout = error.stdout.decode() if isinstance(error.stdout, bytes) else (error.stdout or "")
            stderr = error.stderr.decode() if isinstance(error.stderr, bytes) else (error.stderr or "")
            stderr += f"\nTimed out after {timeout_seconds} seconds."
            returncode = 124
    latency = time.monotonic() - started
    if provider["adapter"] == "pi":
        provider_error = pi_provider_error(stdout)
        if provider_error:
            # Pi can report a provider failure in JSON while its process exits zero.
            # Promote it to a retryable infrastructure status.
            if returncode == 0:
                returncode = 75
            stderr += "\nPi provider error: " + provider_error
            parsed, parse_error = None, provider_error
        else:
            parsed, parse_error = extract_pi_assistant(stdout, expected_root)
    else:
        parsed, parse_error = extract(stdout, expected_root)
    if parsed is not None:
        errors = (
            findings_response(parsed)
            if expected_root == "findings"
            else judgement_response(parsed, expected_ids or set())
        )
        if errors:
            parse_error = "; ".join(errors)
            parsed = None
    return Invocation(
        command=command,
        command_display=display,
        returncode=returncode,
        stdout=stdout,
        stderr=stderr,
        latency_seconds=latency,
        parsed=parsed,
        parse_error=parse_error,
        usage=usage_from(stdout),
    )
