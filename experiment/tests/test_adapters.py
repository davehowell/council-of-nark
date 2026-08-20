from __future__ import annotations

import json
import unittest

from experiment.harness.adapters import (
    command_for,
    extract,
    extract_pi_assistant,
    pi_provider_error,
)
from experiment.harness.summarize import token_counts
from experiment.harness.validate import findings_response, judgement_response


VALID = {
    "findings": [
        {
            "location": "design item 1",
            "claim": "The path blocks.",
            "consequence": "Acknowledgement exceeds the limit.",
            "fix": "Queue before acknowledgement.",
            "confidence": "high",
        }
    ]
}


class AdapterTests(unittest.TestCase):
    def test_extracts_plain_json(self) -> None:
        parsed, error = extract(json.dumps(VALID), "findings")
        self.assertIsNone(error)
        self.assertEqual(VALID, parsed)

    def test_extracts_claude_result_wrapper(self) -> None:
        wrapper = {"type": "result", "result": json.dumps(VALID), "usage": {"input_tokens": 12}}
        parsed, error = extract(json.dumps(wrapper), "findings")
        self.assertIsNone(error)
        self.assertEqual(VALID, parsed)

    def test_rejects_unknown_fields(self) -> None:
        invalid = {"findings": [{**VALID["findings"][0], "severity": "major"}]}
        self.assertTrue(findings_response(invalid))

    def test_pi_parser_never_uses_echoed_user_prompt(self) -> None:
        stream = "\n".join(
            [
                json.dumps(
                    {
                        "type": "message_end",
                        "message": {
                            "role": "user",
                            "content": [{"type": "text", "text": json.dumps(VALID)}],
                        },
                    }
                ),
                json.dumps(
                    {
                        "type": "message_end",
                        "message": {
                            "role": "assistant",
                            "content": [],
                            "stopReason": "error",
                            "errorMessage": "quota exceeded",
                        },
                    }
                ),
            ]
        )
        parsed, _ = extract_pi_assistant(stream, "findings")
        self.assertIsNone(parsed)
        self.assertEqual("quota exceeded", pi_provider_error(stream))

    def test_pi_parser_reads_assistant_text_only(self) -> None:
        stream = "\n".join(
            [
                json.dumps(
                    {
                        "type": "message_end",
                        "message": {
                            "role": "user",
                            "content": [{"type": "text", "text": '{"findings":[]}'}],
                        },
                    }
                ),
                json.dumps(
                    {
                        "type": "message_end",
                        "message": {
                            "role": "assistant",
                            "content": [{"type": "text", "text": json.dumps(VALID)}],
                            "stopReason": "stop",
                        },
                    }
                ),
            ]
        )
        parsed, error = extract_pi_assistant(stream, "findings")
        self.assertIsNone(error)
        self.assertEqual(VALID, parsed)

    def test_cli_schema_omits_unsupported_meta_schema_annotation(self) -> None:
        schema = {"$schema": "https://json-schema.org/draft/2020-12/schema", "type": "object"}
        command, _ = command_for(
            {"adapter": "claude", "model": "example-model", "effort": "low"},
            "prompt",
            "system",
            schema,
            "schema.json",
        )
        schema_argument = command[command.index("--json-schema") + 1]
        self.assertNotIn("$schema", schema_argument)
        self.assertIn('"type":"object"', schema_argument)

    def test_pi_usage_includes_reasoning_and_nested_cost(self) -> None:
        usage = {
            "usage": {
                "input": 100,
                "output": 20,
                "reasoning": 5,
                "cost": {"total": 0.0125},
            }
        }
        self.assertEqual((100, 25, 0.0125), token_counts(usage))

    def test_false_findings_require_a_semantic_cluster(self) -> None:
        response = {
            "judgements": [
                {
                    "item_id": "i-one",
                    "defect_id": None,
                    "false_positive_cluster": "invented cache outage",
                    "material": False,
                    "confidence": "high",
                    "rationale": "No packet support.",
                }
            ]
        }
        self.assertEqual([], judgement_response(response, {"i-one"}))
        response["judgements"][0]["false_positive_cluster"] = None
        self.assertTrue(judgement_response(response, {"i-one"}))


if __name__ == "__main__":
    unittest.main()
