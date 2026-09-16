# A finite closing round

Status: proposed 7 September 2026; launch gate 1 completed 15 September 2026; **not preregistered or executable yet**. This supersedes neither historical protocols nor frozen results. It recommends deferring the 480-call synthetic factorial, provider diversity, all chain permutations, and Modular integration to keep the closing question tractable.

## Question and estimand

For the declared model, task set, and resource policy, does a small council produce better supported, actionable reviews than a single reviewer or sequential self-review? Does fictional wording improve, harm, or leave practical quality unchanged when the functional kernels are identical?

This is diagnosis/review, not end-to-end implementation. A working patch benchmark would require a separate design. The target of generalisation must be declared before running: a fixed benchmark average can be estimated by repeats, but software engineering generally requires new tasks and repositories.

## Five arms

| Arm | Procedure | Review/fusion sessions per task-repeat | Purpose |
|---|---|---:|---|
| S | One omnibus functional review | 1 | Cheapest practical baseline |
| I | Initial omnibus review, two informed self-revisions, final consolidation in one continuing session | 1 session / 4 stages | Fewer-agent alternative with an equal total budget |
| R | Three independent omnibus reviews, one neutral fuser | 4 | Controls repeated sampling and fusion |
| F | Correctness, architecture, and operational-risk specialists, one neutral fuser | 4 | Tests functional specialisation |
| P | Exactly the F kernels plus K-2SO, GLaDOS, and Holly overlays, same neutral fuser | 4 | Tests this three-persona package |

Use the same pinned model/provider and observed thinking state throughout. All arms use the same schema, source access, non-leading brief, and capped final length. The three lenses are selected for relevance before observing responses. Results apply to this reduced council, not every original council member. P−F tests the package; individual role results are descriptive and cannot identify eight individual persona effects.

The live comparison will use `google/gemma-4-31b-it` with `minimal` thinking, matching the earlier low-reasoning calibration and the clean boundary doctor. Freeze the exact Pi/runtime digests in the live configuration before the first call. Changing the model, provider or thinking level later changes the target of the experiment and requires a new preregistration rather than an in-place substitution.

Count **provider turns**, sessions, and tool calls separately: a tool-using stage may make multiple provider turns. The 17 stages per block are not a promise of 17 API requests. I's final stage sees its own evolving history; R/F/P fusers see all three reviews in a deterministically shuffled order, with role labels removed but wording retained. This history difference is the treatment. Fusers receive identical source capability and their own allocated budget; evidence keys remain unavailable.

Budget proposal for feasibility only: S receives 8,000 total charged input/output tokens; I/R/F/P each receive 32,000, allocated 8,000 per stage. The resulting 136,000-token ceiling for the excluded Gortex block is the agreed spend limit; provider cost is approved within it. Set equal aggregate source bytes, tool calls, timeout and final-length caps for I/R/F/P. Count repeated context, cached input, reasoning tokens when exposed, retries and fuser spend. Enforce aggregate limits in the controller, recording unavoidable last-response overrun; do not imply an unenforced cap is exact. Freeze actual numerical tool/byte/time limits after the excluded feasibility run. Report both actual resource use and caps. Call-count matching alone is not token matching.

There is no calendar deadline because this is a personal project. Calendar time is not a stopping rule. Each provider stage still has a frozen wall-clock timeout to stop a hung process and keep treatment budgets comparable. Once a block starts, follow its finite allocation schedule without inspecting outcomes or adding stages; preserve and report any infrastructure pause.

## Tasks and feasibility

Use the four curated pilot tasks (one each from Manim, dlt, Gortex and turbovec) only after each history-free parent and offline dependency closure passes the exact before/after regression check. Gortex currently has infrastructure evidence; the other three are candidates. Symptom briefs must not give away the unsafe operation. New task selection must be performed without inspecting arm outcomes.

The isolation probes and deterministic single-stage mock lifecycle passed from committed code on 15 September 2026 with zero provider calls. Next, one Gortex block checks completeness, time, tool consumption, rubric fit and budget feasibility; exclude it from confirmation. A four-task × two-repeat feasibility pilot uses 8 × 17 = **136 stages**, excluding the initial 17-stage block, tool-driven turns, retries and human-rating work. It can reveal a ceiling or broken arm. It cannot establish narrow equivalence. Stop after feasibility if the available time or money cannot support an adequate task sample; publish that limitation instead of multiplying repeats on easy packets.

The four reserve tasks are not automatic outcome-driven replacements. Substitute only for a documented pre-response eligibility failure within repository, preserving the exclusion log. Never replace a task because the baseline scored too well. Keep Modular excluded until its distinct evidence and build requirements are independently satisfied.

## Primary outcome and human aggregation

Before the first comparative run, freeze task-specific anchors equivalent to [SCORING.md](ecological/SCORING.md). Do not reuse Unicode-specific anchors literally for unrelated repositories. Two humans score independently, with one permitted disclosed author prior and a second independent rater. Preserve raw wording and collect treatment guesses after scoring.

Define review utility Q on [0,1]: mean of the five adjudicated dimension scores, rescaled as (mean−1)/4. Set Q=0 for missing/malformed final output, an unsupported primary mechanism, or a critical unsupported remedy claim. Publish every original dimension and both raters' scores; the composite is a declared decision aid, not a substitute for correctness. Adjudicate mechanism disagreements before aggregation. Material unsupported-claim rate and critical-failure rate are safety endpoints; supported novel findings count even if absent from the eventual patch.

For every comparison, average paired differences within each task before averaging across tasks. Report fixed-set results explicitly. For broader inference account for repository clustering; four repositories are insufficient for precise population inference. More repeated calls cannot repair that limitation. Report raw panel coverage and fusion retention descriptively using blinded semantic claim mapping; do not invent planted-defect F1 for ecological tasks.

## Decision rules and uncertainty

Freeze a practical Q margin of ±0.05 (equivalent to 0.2 points on the 1–5 mean rubric), subject to a documented pre-run maintainer judgement that such a difference would not change adoption. Margin choice is a value judgement, not something the existing F1 pilot validates.

Primary contrasts are F−S (practical benefit), F−I (council versus self-review), F−R (specialisation), and P−F (personality package). Use simultaneous intervals controlling the four-comparison family at 5%; a conservative option is 98.75% two-sided intervals for each contrast. Freeze an inference method suitable for the final sampling frame, including task/repository clustering, before confirmation.

- Meaningful benefit: entire interval above +0.05, with no disqualifying increase in critical errors.
- Meaningful harm: entire interval below −0.05.
- Practical equivalence: entire interval inside [−0.05,+0.05]. Merely including zero does not qualify.
- Otherwise: inconclusive. A confidently nonzero but sub-margin effect is distinguishable from practical importance; publish the interval.

Claim a cheaper substitute only when quality meets the preregistered noninferiority condition and a paired resource interval supports the saving. Require at least 20% fewer actual total tokens for a material token-saving claim; report dollar and latency outcomes separately. A zero reported price cannot establish cost efficiency. Do not assert equal safety from observing zero rare critical errors; publish uncertainty for those rates.

## Precision and stopping

Use the excluded feasibility run to estimate variability and ceiling risk, then freeze sample size and a hard resource cap **before** confirmation. A rough planning approximation is n ≈ (z × s / h)^2 independent task differences for desired interval half-width h; simulation must incorporate bounded outcomes, repository clustering, four contrasts and equivalence power before selecting n. For illustration only, s=0.15, h=0.05 and z≈2.50 requires about 57 independent tasks, not four tasks repeated 57 times. An interval narrow enough is necessary but not sufficient to pass equivalence.

One confirmatory analysis at the frozen sample size. No repeated looks until significance, no changing margins, no adding favourable roles. If the budget cannot fund that precision, conclude the study with calibrated limits. Replicators rerun frozen manifests and analysis, expecting a distribution of outputs rather than identical prose. Preserve all sealed attempts, exclusions, rating locks and analysis versions.

## Remaining launch gates

1. **Complete (15 September 2026):** the JSON/print single-stage launcher captures exact requests, validates one final submission, enforces declared aggregate limits, removes ephemeral credentials, and seals the attempt. Clean snapshot, mediator, Pi doctor, deterministic mock, and negative boundary probes passed. The mock used zero provider calls and is not evidence.
2. Validate eligible snapshots and freeze non-Gortex evidence keys and anchors.
3. **Partly complete:** Dave is confirmed as the disclosed-author rater; Gemma/minimal and the 136,000-token Gortex-block ceiling are agreed, with no calendar deadline. Confirm the independent rater (Dan is the current candidate), then freeze exact arm prompts, per-stage wall timeouts and the allocation schedule.
4. Run the excluded feasibility stage, then commit a completed preregistration with model/runtime digests, numerical budgets, sample size, seed, failure/retry rules, intervals, safety rule and stopping rule.

These gates require real work and human input. This proposal does not convert infrastructure success into a council result.
