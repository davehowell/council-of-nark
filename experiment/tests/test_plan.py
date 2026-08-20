from __future__ import annotations

import unittest

from experiment.harness.common import ROOT, load_json
from experiment.harness.plan import build


class PlanTests(unittest.TestCase):
    def test_stage_a_smoke_counts(self) -> None:
        config = load_json(ROOT / "experiment/config/stage-a-smoke.json")
        plan = build(config)
        self.assertEqual(81, len(plan["calls"]))
        self.assertEqual(27, len(plan["output_sets"]))
        ids = [call["call_id"] for call in plan["calls"]]
        self.assertEqual(len(ids), len(set(ids)))

    def test_topology_smoke_counts_and_dependencies(self) -> None:
        config = load_json(ROOT / "experiment/config/topology-smoke.json")
        plan = build(config)
        self.assertEqual(144, len(plan["calls"]))
        call_ids = {call["call_id"] for call in plan["calls"]}
        for call in plan["calls"]:
            self.assertTrue(set(call["depends_on"]).issubset(call_ids))
            self.assertNotIn(call["call_id"], call["depends_on"])

    def test_provider_pair_is_byte_pair_design(self) -> None:
        config = load_json(ROOT / "experiment/config/provider-pair-smoke.json")
        plan = build(config)
        self.assertEqual(18, len(plan["calls"]))
        self.assertEqual(18, len(plan["output_sets"]))


if __name__ == "__main__":
    unittest.main()
