from __future__ import annotations

import json
import unittest

from experiment.harness.adapters import extract
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


if __name__ == "__main__":
    unittest.main()
