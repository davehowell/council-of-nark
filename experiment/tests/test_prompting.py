from __future__ import annotations

import unittest

from experiment.harness.common import ROOT, load_json
from experiment.harness.plan import build
from experiment.harness.prompting import answer_text, assert_answer_key_absent, static_prompt


class PromptTests(unittest.TestCase):
    def test_every_static_stage_a_prompt_excludes_answer_key(self) -> None:
        config = load_json(ROOT / "experiment/config/stage-a-smoke.json")
        plan = build(config)
        for call in plan["calls"]:
            if call["prompt_spec"]["kind"] not in {"generic", "omnibus", "specialist"}:
                continue
            prompt = static_prompt(ROOT, call)
            assert_answer_key_absent(prompt, answer_text(ROOT, call["packet"]))
            self.assertIn('"findings"', prompt)

    def test_specialist_wrapper_word_counts_are_within_ten_percent(self) -> None:
        for row in load_json(ROOT / "experiment/prompts/specialists.json"):
            functional = len(row["functional_wrapper"].split())
            fictional = len(row["fictional_wrapper"].split())
            self.assertLessEqual(abs(functional - fictional) / max(functional, fictional), 0.10, row["id"])


if __name__ == "__main__":
    unittest.main()
