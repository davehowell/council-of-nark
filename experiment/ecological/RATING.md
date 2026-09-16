# Ecological review rating protocol

> Status: proposed for the excluded Gortex feasibility block. Freeze the bundle generator, scales, output order, two raters, and automated-assessor configuration before rating a live response.

## Purpose

This study evaluates review and research reports, not generated patches. The practical question is whether a report helps a maintainer understand a problem, judge risk, choose a response, or notice something the upstream resolution missed.

The merged upstream change is the **reference resolution**. It is a demonstrated baseline because it was accepted by the project and passes the frozen regression evidence. It is not a gold standard. An experimental review may recommend an equally sound alternative, identify a supported concern that the merged change omitted, or make a different idiomatic trade-off.

Compilation and test success validate the task fixture and any upstream patch. They are not quality scores for a review report. The respondent may run the one allowed diagnostic test as evidence, but it cannot submit code.

## What each rater sees

The offline rating page shows the same task panel for every output:

1. the ecological problem statement given to respondents;
2. the original repository's pull-request title and succinct resolution description;
3. the applied patch;
4. the added or changed regression tests;
5. a notice that the reference is strong evidence, not the only acceptable answer.

Below that panel, the page shows one complete experimental review under an opaque ID. It preserves the review's wording, cited locations, mechanism, consequences, proposed correction, tests, evidence, confidence, and uncertainties.

The page does **not** include the arm, persona/kernel label, provider metadata, run ID, token use, latency, tool transcript, output position before shuffling, or other reviews. Those fields are rejected from the rating-bundle format.

## Blinding

- Give every output a private-key HMAC ID and shuffle output order before either rater starts.
- Show one output at a time. Lock its absolute ratings before showing another output.
- Use the same reference panel and scales for every output and every arm.
- Preserve raw output wording. Rewriting it would remove part of the treatment. After scoring, ask the rater which treatment they think produced it and whether wording disclosed the treatment.
- Dave is a disclosed-author rater. He has repository access and may remember study details; his procedural blinding is therefore based on honest non-use of the private condition map, not claimed ignorance of the project.
- The second rater must be independent and must not receive the private condition map.
- Do not unblind either rater's file until both files and any adjudication rows are locked.

The reference patch is common evidence, so showing it does not reveal which arm produced an output. It can anchor raters too closely to the maintainer's solution. The instructions must therefore credit supported alternatives and novel findings explicitly.

## Phase 1: absolute rating

Rate each valid review independently before any side-by-side comparison. A missing or malformed final submission is not converted into a displayable review: the controller applies the declared Q=0 rule, records the failure, and omits it from human scoring. Do not fabricate placeholder prose for a failed arm.

### Evidence-grounded dimensions

Use the five task-specific 1–5 dimensions in [`SCORING.md`](SCORING.md):

1. root-cause localisation;
2. material consequence;
3. correction correctness and scope;
4. regression-test discrimination;
5. remedy quality and actionability.

These form technical quality **Q** as specified in [`../FINAL_ROUND.md`](../FINAL_ROUND.md). Claim-supportedness labels and critical unsupported-remedy rules still apply. Agreement and adjudication use these scores, not personal taste.

### Personal review utility

Record these 1–7 scales separately. Do not average them into Q.

| Field | 1 | 4 | 7 |
|---|---|---|---|
| **Personal usefulness** | Actively misleading or wastes review time. | Mixed or neutral; some usable material but no clear decision improvement. | Materially improves what I would investigate, say, or decide. |
| **Judgment and idiomatic fit** | Recommendations conflict with sound language/project practice. | Reasonable judgment, though I would make different trade-offs. | Matches or improves on the judgment I want from a trusted reviewer. |
| **Useful value beyond the reference** | Adds unsupported distraction. | Correctly restates the reference but adds no material value. | Adds an important supported risk, alternative, test, or framing. |
| **Signal to noise** | Important points are obscured by noise or poor priorities. | Usable but could be shorter or better ordered. | Concise, well-prioritised, and easy to act on. |

Scores 2, 3, 5, and 6 represent positions between the printed anchors. **Dave's personal-usefulness score is a declared user-centred outcome**, not a proxy for universal quality. Report the independent rater's utility scores separately; do not average away a real taste difference.

Also record:

- relation to the reference: `worse`, `roughly_equivalent`, `better_in_useful_ways`, or `different_not_comparable`;
- the most useful point;
- the main weakness or misleading claim;
- what, if anything, the rater would act on;
- optional qualitative commentary;
- treatment guess and confidence;
- whether wording appeared to reveal the treatment.

## Phase 2: paired preference

After all absolute ratings are locked, present only the four preregistered contrasts: F–S, F–I, F–R, and P–F. Randomise left/right placement independently. Keep the common task/reference panel visible.

For each pair, ask:

- Which review would you rather receive for a real review or research decision?
- Preference strength from strong-left (−3) through tie (0) to strong-right (+3).
- Which review is more technically trustworthy?
- Which review has better judgment and prioritisation?
- What specific difference drove the choice?
- Treatment guess, confidence, and wording disclosure.

Pairwise preference is secondary to the independently locked scores. It tests practical choice directly without allowing one unusually polished comparison to rewrite earlier absolute ratings.

## Claim mapping and adjudication

After whole-review ratings, split findings into checkable claims as described in [`SCORING.md`](SCORING.md). The reference resolution, source, focused test, maintainer explanation, and language/runtime rules are evidence. The upstream patch is not a whitelist.

A supported alternative remedy remains supported even when it differs from the merged patch. A novel finding remains eligible when it is relevant, evidenced, and material. Mark a claim unsupported only because evidence contradicts it, it overstates certainty or impact, or it proposes a harmful remedy—not because maintainers chose another style.

The two raters lock their claim labels independently. Adjudicate material disagreements while arm labels remain hidden. Preserve both original ratings.

## Automated assessment

Automated assessment is secondary and uses the same opaque output IDs.

### Deterministic validation

Record, but do not turn into a quality score:

- schema validity and one terminating submission;
- cited path and line-range validity;
- whether cited excerpts exist in the frozen source;
- duplicate findings and empty evidence fields;
- malformed output, resource use, denied requests, and budget overruns;
- validity of the frozen parent/evidence regression fixture.

Do not award points because code compiles or tests pass. Respondents submit reviews, not patches. Test invocation is an evidence-gathering choice and a resource measure, not proof that the final reasoning is good.

### Blinded model assessor

One or more separately declared assessor models may receive the same problem statement, reference panel, frozen source excerpts, and opaque review shown to humans. They return the same five technical dimensions, claim labels, four utility scales, reference relation, and concise reasons.

Automated ratings must:

- use a frozen prompt/schema/model configuration;
- omit arm labels and resource metadata;
- run only after respondent attempts are sealed;
- remain separate from human ratings;
- report agreement and systematic differences rather than replacing human judgment;
- never adjudicate Dave's personal usefulness or taste on his behalf.

## Rating records

The offline page in [`rating/`](rating/) imports a strict label-blinded JSON bundle and exports one JSON file per rater. It makes no network request. The exported file contains opaque output IDs, scores, comments, condition guesses, and timestamps, but no condition map.

The future whole-block controller must create:

- one sealed absolute-rating bundle with deterministic HMAC IDs and shuffled order;
- one private output-to-condition map;
- paired bundles only after absolute ratings are locked;
- a digest of each input bundle in each rating record;
- a mechanical check that forbidden condition/resource fields are absent.

Do not hand-edit a bundle after rating starts. If presentation is wrong, preserve the failed bundle, repair the generator, create new opaque IDs, and restart rating before unblinding.
