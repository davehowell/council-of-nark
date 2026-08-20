from __future__ import annotations

import unittest

from experiment.harness.score import bootstrap_mean_ci


class ScoreTests(unittest.TestCase):
    def test_bootstrap_interval_is_deterministic(self) -> None:
        first = bootstrap_mean_ci([-0.1, 0.0, 0.1, 0.2], "fixed", samples=1_000)
        second = bootstrap_mean_ci([-0.1, 0.0, 0.1, 0.2], "fixed", samples=1_000)
        self.assertEqual(first, second)
        assert first is not None
        self.assertLessEqual(first[0], 0.05)
        self.assertGreaterEqual(first[1], 0.05)


if __name__ == "__main__":
    unittest.main()
