# What this study can conclude

Status: 7 September 2026. Calibration is complete; comparative ecological evidence is not. This report closes the existing evidence review without presenting unrun work as a result.

## Supported statements

The council has not demonstrated a reliable advantage over a single functional review. It has also not been shown equivalent to one. In the Haiku smoke, S1 and M0 both averaged F1 0.860, with M0 costing 8.924 times as much in recorded usage. Equal scores did not always mean identical defects. In the Gemma smoke, functional specialists averaged F1 0.911 against repeated omnibus at 0.863 and single functional omnibus at 0.804. Each smoke has only three tasks and one draw per arm; the direction reversal cannot be attributed solely to the model because scoring and harness details changed too.

The K-2SO overlay underperformed its byte-identical correctness kernel on the three synthetic packets. Thirty pairs averaged fictional-minus-functional F1 −0.04245. This is one model, one role, three tasks, and exploratory LLM ratings. It does not establish how the other characters perform, or how any character affects a full council.

New sensitivity analysis: removing any one packet leaves a negative mean difference (−0.03364 to −0.04892). The observed direction is therefore not driven solely by one packet. However, flipping whole task blocks yields only eight sign assignments; the illustrative two-sided exact p-value is 0.25. This test assumes symmetric, exchangeable task effects and is post hoc. It illustrates the shortage of independent tasks; it neither reverses the sampling-only interval nor proves no effect.

The overlay used **fewer**, not more, recorded tokens: 40,182 versus 41,000 total input plus output tokens, about 2% less. It was worse on detection but marginally cheaper on this token measure. The data do not support saying this persona increased token spend. Recorded dollar cost was zero; that is not evidence that compute is free. Cross-model token accounting is not directly comparable.

## Reproduce the arithmetic

Run `python3 scripts/evidence_audit.py`. It verifies all published-file digests in the three result manifests, recalculates F1 from TP/FP/FN, reconstructs complete persona pairs from CSV, and emits stage comparisons, recorded costs/tokens, per-task effects, and leave-one-task-out sensitivity. It does not call a model, re-rate findings, or claim to reverify the unpublished raw seals. Changing a published input without updating its manifest fails verification. Historical artifacts stay unchanged.

TP is a unique keyed defect detected; FP is an unsupported claim under that run's clustering rules; FN is a keyed defect missed. Precision is TP/(TP+FP), recall is TP/(TP+FN), and F1 is 2TP/(2TP+FP+FN). Macro means give tasks equal weight, averaging repeated samples within task. The first smoke did not semantically cluster raw false positives, so its apparent fusion improvement is overstated. These metrics do not measure successful code changes or remedy quality.

## Important scope correction

These experiments compare **review outputs**, not an agent implementing a change against an implementation team. A conclusion about fixing software requires a separate endpoint: executable patches, hidden tests, regression checks, and equal implementation budgets. The proposed closing round stays with diagnosis and review, which is the council's stated purpose. Do not title it a coding-productivity benchmark.

Names may still be convenient labels. There is no demonstrated need to retain character prose. A reasonable operational default is one functional review, escalating when risk warrants it; that is a cost-conscious decision under uncertainty, not a proven equivalence result.

## What remains

[FINAL_ROUND.md](FINAL_ROUND.md) proposes a finite comparison of a single review, sequential self-review, repeated independent reviewers, functional specialists, and the same specialists with fictional overlays. It explicitly separates a small feasibility pilot from any precision-based confirmatory sample. Existing external-project snapshots supply candidate benchmark tasks, not benchmark outcomes. A completed launcher, eligible frozen tasks, and two independent human ratings are still required. No new respondent call was made for this report.

Sources: [Haiku smoke](results/2026-08-20-stage-a-smoke/README.md), [Gemma smoke](results/2026-08-20-stage-a-smoke-gemma/README.md), [persona pairs](results/2026-08-20-persona-pair-gemma/README.md), [checkpoint](CHECKPOINT.md), [scoring](ecological/SCORING.md).
