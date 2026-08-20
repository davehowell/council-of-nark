# Scenario hardness and calibration

Hardness means the probability that a review condition finds a planted mechanism without producing unsupported claims. It is not the length of the packet or the obscurity of a technology name. Hardness is relative to model capability: a packet that separates small models can be a ceiling test for a frontier model.

Use the least-capable model that can reliably follow the response contract when calibrating prompt effects. If every arm is near the ceiling, lower model capability before adding contrived defects. Preserve a spread of easy, medium, and hard components once task variants are introduced.

## Hardness dimensions

Each packet mixes these dimensions:

| Dimension | Easier | Harder |
|---|---|---|
| Locality | mechanism and consequence appear together | facts must be joined across sections |
| Explicitness | packet states the violated limit beside the action | limit and action are separated |
| Lens overlap | several lenses can identify the same defect | one narrow lens is needed |
| Distractors | every suspicious detail is defective | correct details resemble common defects |
| Remedy specificity | one obvious fix | several fixes exist but only one preserves stated constraints |
| Operational horizon | immediate failure | failure appears after retry, neglect, drift, or hand-off |

The base packets deliberately contain local and cross-section defects. They also name facts that are not defects. A reviewer must show a concrete consequence and fix, which reduces credit for keyword spotting.

## Smoke calibration rules

The smoke stage checks instruments; it does not test hypotheses. Inspect scores only after the run is sealed and ratings are arm-blinded.

Flag a packet for revision when any condition holds:

- S1 mean recall is below 0.15: likely floor, missing context, or an ambiguous key.
- S1 mean recall is above 0.85 with precision above 0.90: likely ceiling.
- every M arm's raw union reaches all defects: raw coverage has no room to distinguish specialisation.
- more than 10% of responses are malformed: schema or adapter failure.
- two independent raters disagree on more than 20% of finding-to-defect mappings: answer key lacks operational boundaries.
- a clean control attracts more than two material findings in one output: the prompt rewards quota filling.

These are calibration triggers, not exclusion rules. Report every trigger and every resulting edit. Do not choose a threshold because it favours the council.

## Permitted post-smoke edits

Before preregistration, the study may:

- clarify packet facts that raters could not interpret consistently;
- split one compound planted defect when it has two independent consequences;
- merge defects that raters cannot distinguish semantically;
- shorten wrappers to meet provider token tolerance;
- adjust output caps when truncation, rather than review ability, determines scores;
- add clean controls and frozen variants.

Do not move or reword a defect merely because one arm missed it. Record changed file digests and start a new frozen run.

## Confirmatory variants

Create variants before the confirmatory run. For each base packet, include:

- semantic renames of organisations, services, tables, and fields;
- equivalent syntax and reordered sections;
- defect locations moved without changing evidence;
- scale values changed while preserving the same violated ratio;
- at least two clean controls where risky-looking details satisfy every stated constraint;
- ablations that remove one planted defect at a time.

A second person should verify that every mutation preserves its answer key. Freeze variants before selecting the confirmatory sample count.

## Hardness reporting

For each packet and variant, publish:

- S0 and S1 recall and precision;
- per-defect detection frequency;
- clean-control false positives;
- rater disagreement;
- response truncation and malformed rates;
- input length and provider-reported tokens.

A packet that differentiates prompts but depends on one obscure product fact is narrow, not scientifically hard. The packet must supply every fact needed for its answer.
