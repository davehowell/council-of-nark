# Stage A smoke results: exploratory only

This run calibrated the harness and packets. It is one sample per `(packet × arm)`, has three synthetic packets, and uses LLM triage rather than human adjudication. It cannot support inferential claims.

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

- Review source: `d421a4d3aba1a73dcf4706d23a22e69df8d6e1a8` (`experiment-harness-v0.1.1`)
- Model: `claude-haiku-4-5-20251001`, low effort, provider-default temperature
- Isolation: 81 fresh processes in 81 clean detached worktrees; tools, project context, and session persistence disabled
- Status: 81/81 successful structured responses; 0 malformed; 0 retries that became samples
- Wall time: about 48 minutes at concurrency 3
- Provider-reported usage: 143,658 input tokens; 849,024 output/reasoning tokens; $4.8011
- Raw findings before output-set union: 550

The raw run is sealed locally. [`manifest.json`](manifest.json) records its source/config digests and raw-seal digest.

## Fused/final output sets

| Arm | Mean F1 | Precision | Recall | p10/worst F1 | Total recorded cost for 3 packets |
|---|---:|---:|---:|---:|---:|
| S0 generic | 0.812 | 1.000 | 0.708 | 0.667 | $0.1295 |
| S1 functional omnibus | **0.860** | 0.944 | **0.792** | **0.714** | $0.1611 |
| S2 omnibus + GLaDOS | 0.802 | 0.933 | 0.708 | 0.615 | $0.1641 |
| M0 repeated omnibus + fuse | **0.860** | **0.944** | **0.792** | **0.714** | $1.4380 |
| M1 functional specialists + fuse | 0.752 | 0.808 | 0.708 | 0.667 | $1.3784 |
| M2 fictional specialists + fuse | 0.783 | 0.821 | 0.750 | 0.667 | $1.5300 |

Predefined smoke contrasts on mean F1:

- Specialisation, `M1 − M0`: **−0.108**.
- Fictional specialist wrapper, `M2 − M1`: **+0.031**.
- Single-call fictional wrapper, `S2 − S1`: **−0.058**.
- Repeated sampling, `M0 − S1`: **0.000**, at about 8.9× the recorded cost.
- Practical fictional council, `M2 − S1`: **−0.077**.

These are descriptive differences over three sets. No confidence interval or significance claim is valid.

## Raw unions and fusion

| Arm | Raw-union F1 | Raw precision | Raw recall | Fused F1 | Mean fusion retention |
|---|---:|---:|---:|---:|---:|
| M0 | 0.603 | 0.492 | 0.792 | 0.860 | 1.000 |
| M1 | 0.554 | 0.390 | **0.958** | 0.752 | 0.738 |
| M2 | **0.632** | **0.497** | 0.917 | 0.783 | 0.819 |

Fusion removed enough false positives to improve F1 in every multi-reviewer arm. It retained every detected true defect in M0, but dropped valid findings in M1 and M2. This supports testing fusion as a separate component rather than crediting the whole panel.

This historical smoke counted every unmatched raw finding as a false positive. It did not cluster semantically duplicate false claims, so raw-union precision is biased downward and fusion gain is overstated. The corrected harness clusters false claims in later runs.

## Calibration interpretation

This was a plumbing smoke, not a go/no-go decision about the council:

- S0 already achieved 0.812 mean F1 and perfect recall on the webhook packet, showing limited headroom for this model/task pairing.
- One sample per packet cannot distinguish a prompt effect from decoding variation.
- M1 raw recall reached 0.958, above M0's 0.792, while its precision was lower. That is consistent with specialists exposing more different issues and more noise; it is not captured by fused F1 alone.
- The functional specialist panel did not beat matched repeated omnibus sampling after fusion in this draw.
- Fictional wrappers were inconsistent: S2 was below S1, while M2 was slightly above M1 and below M0.

S1 and M0 had identical packet-level F1 values but did not always find the same defects. On the revenue packet, S1 found `RD-04` and missed `RD-06`; M0 did the reverse. Equal F1 describes equal counts, not equivalent reviews.

Next: repeat the unchanged smoke with explicit low-reasoning Gemma to restore prompt headroom. If that run has useful spread, measure within-arm variance with repeated randomised blocks. Harder or real scenarios come only after the cheaper-model test.

## Rating limitation

The arm-blinded triage contains 550 item ratings:

- 544 from `claude-sonnet-4-6` at high effort;
- 6 from `openai-codex/gpt-5.4-mini` at low effort after the default rater hit a session limit;
- 0 human ratings.

The mixed fallback set was still arm-blinded, but model-judge bias and correlated model families remain serious threats. Treat every score here as smoke-stage instrumentation. [`ratings-tidy.csv`](ratings-tidy.csv) publishes labels and provenance without raw finding prose.

## Files

- [`summary.json`](summary.json): grouped scores and fusion retention.
- [`sets.csv`](sets.csv): one scored row per output set.
- [`ratings-tidy.csv`](ratings-tidy.csv): one unblinded label row per finding, without finding text.
- [`run-health.json`](run-health.json): status, usage, cost, latency, and raw finding count.
- [`manifest.json`](manifest.json): source, config, seal, rater, and published-file digests.
