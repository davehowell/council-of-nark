from __future__ import annotations

import csv
from pathlib import Path
import tempfile
import unittest

from experiment.harness.common import ROOT
from experiment.harness.judge import existing_ratings


class JudgeResumeTests(unittest.TestCase):
    def test_existing_ratings_are_loaded_by_item_id(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "ratings.csv"
            with path.open("w", newline="", encoding="utf-8") as handle:
                writer = csv.DictWriter(
                    handle,
                    fieldnames=["rater", "item_id", "defect_id", "material", "confidence", "notes"],
                )
                writer.writeheader()
                writer.writerow(
                    {
                        "rater": "test",
                        "item_id": "i-one",
                        "defect_id": "RD-01",
                        "material": "True",
                        "confidence": "high",
                        "notes": "",
                    }
                )
            rows = existing_ratings(path)
            self.assertEqual({"i-one"}, set(rows))
            self.assertEqual("RD-01", rows["i-one"]["defect_id"])

    def test_judge_prompt_contains_contract_for_adapters_without_schema_flags(self) -> None:
        prompt = (ROOT / "experiment/prompts/judge.txt").read_text(encoding="utf-8")
        self.assertIn('{"judgements"', prompt)
        self.assertIn('"confidence":"high|medium|low"', prompt)


if __name__ == "__main__":
    unittest.main()
