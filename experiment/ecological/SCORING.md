# Ecological scoring protocol

> Status: frozen for Gortex infrastructure calibration only. Do not make a comparative claim until the arm design, sample count, two raters, aggregation rule, and smallest effect of interest are preregistered.

Ecological outputs are review and research reports, not generated patches, and are not scored as planted-defect F1. The upstream correction and tests are a demonstrated reference resolution, not the only acceptable answer. Raters first score each complete review independently, then map its claims, and only afterward see preregistered blinded pairs. Supported alternative remedies and novel findings remain eligible.

## Rater boundary

Two humans rate independently while condition, provider, prompt wrapper, and output order remain hidden. One may disclose a prior preference for the council, but the other must be independent. Raters receive:

- the ecological problem statement;
- the original repository's pull-request title and succinct resolution description;
- the applied patch and changed regression tests;
- the frozen source locations cited by the output;
- the task-specific evidence key and regression evidence;
- one complete review under an opaque ID.

The reference resolution is common to every condition. It is a baseline for technical adequacy, not a whitelist or gold standard. Raters do not receive the condition map, prompt variants, resource use, run IDs, or other respondents' output during absolute rating. Raw wording is preserved. After scoring, each rater records a condition guess, confidence, and whether wording revealed the likely treatment. This is label blinding, not guaranteed treatment blinding. The complete presentation and personal-utility scales are in [`RATING.md`](RATING.md).

## Atomic claim mapping

Before applying the 1–5 rubric, split a finding into independently checkable claims. Mark each claim:

- **supported** — source, focused test, maintainer evidence, or sound language/runtime reasoning supports it;
- **unsupported** — contradicted, materially overstated, or presented as fact without adequate support;
- **indeterminate** — available evidence cannot decide it.

Semantically duplicate claims count once within an output. Record unsupported claims by severity:

- **minor:** does not change diagnosis or remedy;
- **material:** could misdirect implementation or testing;
- **critical:** recommends a harmful change or invents the primary mechanism.

Conditional proposals and clearly labelled uncertainty are not unsupported claims merely because the upstream patch chose another design.

## Anchored dimensions

Use integer scores. Scores 2 and 4 mean the response falls between the adjacent anchors; write one sentence explaining the choice. A polished style cannot compensate for a technically wrong mechanism.

### 1. Root-cause localisation

| Score | Anchor |
|---:|---|
| 1 | No relevant location, or the identified code cannot cause the symptom. |
| 2 | Points to a broadly relevant subsystem but not the unsafe operation or state transition. |
| 3 | Identifies the unsafe operation and a plausible trigger, but omits an important boundary condition or causal step. |
| 4 | Correct file/region, unsafe assumption, and trigger with only a minor omission. |
| 5 | Exact location and complete causal chain, including why the supplied symptom can depend on repository/corpus state rather than only the query. |

### 2. Material consequence

| Score | Anchor |
|---:|---|
| 1 | Consequence is absent, contradicted, or unrelated. |
| 2 | Mentions failure vaguely without connecting it to the request path or affected users. |
| 3 | Correctly describes the principal failure but overstates prevalence or misses an important precondition. |
| 4 | Correct failure mode, affected path, and materiality with a minor precision gap. |
| 5 | Precisely connects trigger, request/corpus interaction, process or request impact, and scope without unsupported severity inflation. |

### 3. Correction correctness and scope

| Score | Anchor |
|---:|---|
| 1 | Does not prevent the defect or introduces a clear correctness regression. |
| 2 | Could avoid the observed crash but breaks required tokenisation behavior or relies on a fragile special case. |
| 3 | Correct in the reported case, but materially over-broad, allocation-heavy, or incomplete at adjacent Unicode/boundary cases. |
| 4 | Unicode-correct and scoped, with only a minor compatibility, invalid-input, or hot-path concern. |
| 5 | Uses a boundary-safe lookahead or equivalent state machine, preserves existing ASCII/camel/acronym semantics, avoids whole-string or repeated suffix rune allocation, and changes no unrelated behavior. |

### 4. Regression-test discrimination

| Score | Anchor |
|---:|---|
| 1 | No executable case, or cases would pass before the correction. |
| 2 | Repeats the symptom without an expected result or misses the crashing boundary. |
| 3 | Includes at least one pre-fix-failing Unicode case with a correct expected result. |
| 4 | Covers the crash plus a neighboring script/case boundary and an existing behavior guard. |
| 5 | A compact matrix distinguishes the defect and protects lowercase/separator behavior, all-uppercase non-ASCII endings, mixed acronym-to-Camel transitions, and established ASCII behavior without encoding the implementation. |

### 5. Remedy quality and actionability

| Score | Anchor |
|---:|---|
| 1 | Cannot be implemented safely from the description. |
| 2 | Names a direction but leaves substantial mechanism, location, or verification work unspecified. |
| 3 | Implementable by a maintainer, but important scope, performance, or test details are missing. |
| 4 | Clear location, algorithm, edge handling, and tests; only a minor implementation choice remains. |
| 5 | Concise, directly implementable, evidence-linked, explicit about trade-offs and uncertainty, and avoids unrelated refactoring. |

## Personal review utility

Technical quality Q remains the mean of the five dimensions above after the declared rescaling and failure rules. Record the four 1–7 personal-utility scales from [`RATING.md`](RATING.md) separately: personal usefulness, judgment and idiomatic fit, useful value beyond the reference, and signal to noise. Do not average them into Q.

Dave's personal-usefulness score is a declared user-centred outcome because the council is intended for his workflow. The independent rater's score tests whether that preference travels. Publish them separately instead of treating taste disagreement as measurement error.

## Tool and efficiency measures

The controller records these mechanically; raters do not estimate them:

- tool calls by operation and denied calls;
- result bytes and truncation;
- focused-test count and duration;
- input, output, and cached tokens;
- wall latency and provider latency where available;
- reported cost;
- malformed or missing final submission.

Compare arms with identical tool policy and budgets. Report each quality dimension beside cost and latency. Do not create a post-hoc quality-per-dollar composite after seeing outcomes. Compilation and test success validate the frozen task fixture; they do not score a review report. Whether the respondent called the focused test is an evidence-gathering and resource-use observation, not a quality point.

## Disagreement and adjudication

1. Each rater locks claim labels and five dimension scores independently.
2. Calculate exact agreement for claim labels and weighted Cohen's kappa for ordinal dimensions when sample size permits; otherwise publish the raw disagreement table.
3. Discuss only items that differ by a claim label or by at least two rubric points. Keep both originals unchanged.
4. An adjudicator resolves evidence interpretation while labels remain blinded. The adjudicated row cites the evidence and records whether either original changed.
5. Unblind only after all finding mappings, scores, exclusions, and paired preferences are locked.

Never average away a mechanism disagreement before adjudication. Publish individual-rater and adjudicated results.

## Missing and malformed output

A response with no valid final submission receives no finding scores and is reported as malformed. A valid submission with no supported primary diagnosis receives root-cause, consequence, correction, regression, and remedy scores of 1; do not exclude it. Infrastructure failures follow the preregistered retry rule and remain in the run health record.
