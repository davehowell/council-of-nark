from __future__ import annotations

import json
import unittest

from experiment.harness.adapters import command_for, extract
from experiment.harness.validate import findings_response


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


if __name__ == "__main__":
    unittest.main()
