# Repeated Gemma correctness-persona calibration

This focused run compares the same correctness/observability kernel with either a functional wrapper or a K-2SO wrapper. It uses ten paired repeats on each of the three base packets.

## Metric key

| Metric | Meaning |
|---|---|
| TP — true positive | One unique planted defect correctly identified. |
| FP — false positive | One unique unsupported claim after semantic de-duplication. |
| FN — false negative | One planted defect the output missed. |
| Precision | `TP / (TP + FP)`: how much of the review was supported. |
| Recall | `TP / (TP + FN)`: how much of the planted defect set was found. |
| F1 | `2TP / (2TP + FP + FN)`: harmonic balance of precision and recall. |
| Macro mean | Calculate the metric per output set, then weight each set equally. |

F1 does not measure severity, remedy quality, or whether two conditions found the same defects.

## Run

- Source: `cb57a8c0c48f1fe841905b8f0c3a1c61f38e5d92` (`experiment-harness-v0.3`)
- Respondent: `gemma-4-31b-it` through Pi, thinking off
- Rater: `gemini-3.5-flash-low` through agy, arm-blinded
- Calls: 60/60 successful; 0 malformed
- Pairs: 30 = 3 packets × 10 repeats
- Wall time: about 8 minutes at concurrency 2
- Usage: 51,314 input tokens; 29,868 output/reasoning tokens; recorded cost $0

## Aggregate results

| Wrapper | Mean F1 | Precision | Recall | p10 F1 | Worst F1 | Mean output tokens |
|---|---:|---:|---:|---:|---:|---:|
| Functional correctness | **0.790** | **0.980** | **0.675** | **0.667** | 0.462 | 514 |
| K-2SO fictional | 0.748 | 0.954 | 0.629 | 0.615 | 0.462 | 481 |

Paired fictional-minus-functional F1:

- Mean: **−0.0425**.
- Sampling-only bootstrap 95% interval: **[−0.0808, −0.0047]**.
- Fictional wins: **2**.
- Ties: **16**.
- Functional wins: **12**.
- Mean keyed-defect Jaccard overlap: **0.896**.

Per packet:

| Packet | Mean delta | Sampling-only 95% interval | Fictional / tie / functional wins |
|---|---:|---:|---:|
| Revenue dashboard | −0.060 | [−0.155, +0.046] | 1 / 4 / 5 |
| Key rotation | −0.038 | [−0.072, −0.005] | 1 / 4 / 5 |
| Webhook redesign | −0.030 | [−0.073, 0.000] | 0 / 8 / 2 |

## Interpretation

For this model, kernel, and toy task set, the K-2SO wrapper did not improve detection. It reduced mean F1 and recall. The repeated result also shows why the first smoke's small M2-over-M1 difference was not persuasive: wrapper effects can reverse across samples and orchestration levels.

The result does **not** show that all personas are harmful:

- it tests one persona and one functional lens;
- it does not test provider diversity;
- it does not test full council fusion;
- the bootstrap estimates sampling variation over repeated calls to the same three packets, not generalisation to other tasks;
- ratings are LLM-based rather than human;
- remedy quality is unscored.

The wrappers still changed which defects appeared in 14 of 30 pairs, despite a high mean Jaccard overlap. That supports the behavioural premise that wording can shift attention. In this controlled test, the shift was not beneficial on average.

## Next decision

Do not make the synthetic packets more contrived yet. A real-task ecological replication is now more valuable: freeze a pre-fix open-source revision, remove history, provide a non-leading issue brief, and compare findings against the eventual fix and maintainer evidence. Keep the repeated functional control and blinded scoring.

See [`summary.json`](summary.json), [`sets.csv`](sets.csv), [`ratings-tidy.csv`](ratings-tidy.csv), [`run-health.json`](run-health.json), and [`manifest.json`](manifest.json).
