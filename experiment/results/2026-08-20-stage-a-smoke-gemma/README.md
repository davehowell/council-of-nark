# Stage A smoke with low-reasoning Gemma

This repeats the unchanged six-arm plumbing smoke with an explicitly selected lower-capability model. It is one sample per `(packet × arm)` and uses LLM triage, so it remains non-confirmatory.

## Run

- Review source: `27d80c652cf208649e637266bd4c5dc925214016` (`experiment-harness-v0.2.1`)
- Respondent: `gemma-4-31b-it` through Pi, thinking off
- Rater: `gemini-3.5-flash-low` through agy, arm-blinded
- Status: 81/81 respondent calls successful; 0 malformed
- Isolation: 81 fresh processes and detached worktrees
- Wall time: about 12 minutes at concurrency 2
- Usage: 103,149 input tokens; 43,476 output/reasoning tokens; recorded cost $0
- Raw findings before output-set union: 419

## Fused/final sets

| Arm | Mean F1 | Precision | Recall | p10/worst F1 |
|---|---:|---:|---:|---:|
| S0 generic | 0.819 | **1.000** | 0.708 | 0.667 |
| S1 functional omnibus | 0.804 | 0.958 | 0.708 | 0.769 |
| S2 omnibus + GLaDOS | 0.834 | 0.958 | 0.750 | 0.769 |
| M0 repeated omnibus + fuse | 0.863 | 0.958 | 0.792 | 0.857 |
| M1 functional specialists + fuse | **0.911** | 0.958 | **0.875** | 0.857 |
| M2 fictional specialists + fuse | 0.857 | 0.944 | 0.792 | 0.714 |

Descriptive one-draw contrasts:

- `M1 − M0` fused F1: **+0.048**.
- `M2 − M1` fused F1: **−0.054**.
- `S2 − S1`: **+0.029**.
- `M0 − S1`: **+0.059**.

The fictional-wrapper directions disagree between the single omnibus and specialist panel. One draw cannot identify a persona effect.

## Raw semantic unions

| Arm | Raw F1 | Precision | Recall |
|---|---:|---:|---:|
| M0 | 0.863 | 0.958 | 0.792 |
| M1 | 0.861 | 0.824 | **0.917** |
| M2 | **0.877** | 0.849 | **0.917** |

M1 exposed more planted defects than M0 before fusion but also more unique false-positive clusters. The fuser retained enough M1 signal to produce the highest final F1 in this draw.

F1 again hides issue identity:

- on the revenue packet, M1 added `RD-06` beyond M0;
- on the webhook packet, M1 added `WH-01` beyond M0;
- M1 and M0 had the same keyed set on key rotation;
- M1 found `RD-06` and `KR-07` where M2 did not; M2 added no keyed defect beyond M1 in this draw.

## Interpretation

Using a smaller model changed the ordering and exposed the behaviour the first smoke could not resolve: specialist raw recall increased and conditions found different defects. It did not create much lower absolute performance; S0 still scored 0.819, so the three packets remain easy for modern models.

This does not prove specialist or fictional persona effects. It justifies measuring sampling variation before changing scenarios. The next focused run repeats the functional and K-2SO correctness prompts ten times per packet.

## Limitations

- Three synthetic packets, one respondent sample per arm.
- No common temperature control.
- Arm-blinded LLM ratings; no human adjudication.
- Respondent and rater are different models but from the same broad provider ecosystem.
- F1 gives no credit for remedy quality.

See [`summary.json`](summary.json), [`sets.csv`](sets.csv), [`ratings-tidy.csv`](ratings-tidy.csv), [`run-health.json`](run-health.json), and [`manifest.json`](manifest.json).
