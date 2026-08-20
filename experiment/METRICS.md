# Metrics and comparability

## F1

For one output set:

- **True positive (TP):** a unique planted defect ID found by at least one supported finding.
- **False positive (FP):** a unique unsupported or unkeyed defect claim after semantic de-duplication.
- **False negative (FN):** a planted defect ID the output set missed.
- **Precision:** `TP / (TP + FP)`.
- **Recall:** `TP / (TP + FN)`.
- **F1:** the harmonic mean, `2 × precision × recall / (precision + recall)`, equivalently `2TP / (2TP + FP + FN)`.

Example: finding 6 of 8 planted defects with 2 unique false claims gives precision `6/8 = 0.75`, recall `6/8 = 0.75`, and F1 `0.75`.

F1 is zero when the output finds nothing. A duplicate wording of one true or false claim does not create another prediction.

## Macro F1

The Stage A summary calculates F1 separately for each packet/output set, then gives each packet equal weight in the mean. A long packet cannot dominate a short packet merely because it contains more findings.

With only eight defects per packet, recall moves in increments of `0.125`. With three packets and one repeat, mean F1 is coarse. A difference such as `0.03` can be less than one changed detection across the study and has no useful uncertainty estimate.

## What equal F1 does not mean

Equal F1 means equal balance of TP, FP, and FN counts. It does not mean two conditions found the same defects, proposed the same fixes, or contributed the same categories.

In the first smoke, S1 and fused M0 had identical packet-level F1 values. On the revenue packet, S1 found `RD-04` but missed `RD-06`; M0 found `RD-06` but missed `RD-04`. F1 correctly called their counts equal while hiding that qualitative difference.

Publish these alongside F1:

- detected defect IDs and per-defect frequency;
- pairwise Jaccard overlap;
- unique true contributions by reviewer/lens;
- false-positive clusters;
- fix quality or developer acceptance for real tasks;
- p10, variance, and confidence intervals over independent repeats and task variants;
- token, latency, and cost-normalised outcomes.

## Raw union versus fusion

A raw union must de-duplicate both true and false claims semantically. Counting seven reviewers' duplicate unsupported claim as seven false positives biases raw panels downward and exaggerates fusion gain.

Fusion F1 measures the practical final verdict. Raw-union recall and precision measure what the panel made available before arbitration. Use both. A fuser sees the packet and can flatten differences between panels, so role specialisation should not be judged from fused F1 alone.

## Limits

F1 treats every planted defect as equally important and gives no credit for a better remedy to the same defect. Severity-weighted recall and blinded fix-quality ratings can supplement it, but their weights and rubrics must be frozen before a claim-bearing run.
