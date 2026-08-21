# Experiment lab notebook

This is the chronological engineering record. It preserves failed runs and the rationale available at each pivot instead of rewriting the project as a clean success story after the fact. Times are UTC. Published snapshots and Git history remain the authoritative artifacts.

## 2026-08-19 — turn the story into falsifiable claims

- Split the council story into role specialisation, fictional overlay, fusion, topology, and provider questions.
- Created synthetic packets with planted keys, controlled prompt pairs, strict response contracts, and a preregistration template.
- Decision: character names are useful mnemonic labels even if fictional prose fails. The overlay must earn its tokens empirically.
- Decision: use fresh processes and detached worktrees, seal raw output before rating, and blind raters to arm/provider.

## 2026-08-20 04:23 — instrumentation failed before inference

- Planned the first 81-call Stage A smoke.
- Claude rejected the draft-2020 `$schema` annotation locally. No model request succeeded; nine dependent fusers were blocked.
- Preserved the failed run rather than deleting or repairing it.
- Repair: retain `$schema` in source/provenance but omit the unsupported annotation from the CLI argument. Added an adapter check before a larger run.
- Published record: [`results/2026-08-20-stage-a-smoke-instrumentation-failure.md`](results/2026-08-20-stage-a-smoke-instrumentation-failure.md).

## 2026-08-20 04:27 — first complete Stage A plumbing smoke

- Completed 81/81 pinned Haiku calls and sealed the raw run.
- Exploratory arm-blinded rating required a small Pi/OpenAI fallback after the primary judge exhausted a session limit.
- Observed high baseline performance and near-equal fused scores despite conditions finding different defect IDs.
- Observed specialist raw recall above repeated omnibus recall, followed by fusion loss.
- Initial mistake avoided: do not interpret one draw on three easy packets as a verdict on council value.
- Published record: [`results/2026-08-20-stage-a-smoke/`](results/2026-08-20-stage-a-smoke/).

## 2026-08-20 07:00–08:15 — post-smoke contamination and scoring audit

- Found that seeded plan order had been discarded by set iteration in the scheduler.
- Found that raw panels counted repeated unsupported mechanisms as multiple false positives.
- Identified fuser flattening, one-sample decoding noise, and LLM-rater limitations as design threats.
- Repair: preserve seeded order, semantically cluster false claims, publish issue overlap, and make raw-union coverage primary for specialisation.
- Decision: rerun unchanged material on an explicit lower-capability model before making tasks harder.

## 2026-08-20 08:15 — first Gemma attempt discarded

- Selected `gemma-4-31b-it` with thinking off rather than inheriting Pi's configured default.
- Pi echoed the user prompt in JSON events and represented provider quota failures while exiting zero. The generic parser could mistake prompt examples for responses.
- Sealed and discarded the partial 71/81 run before rating.
- Repair: parse Pi assistant events only, promote provider errors to retryable failures, add backoff, and lower concurrency.
- Published record: [`results/2026-08-20-gemma-smoke-incomplete.md`](results/2026-08-20-gemma-smoke-incomplete.md).

## 2026-08-20 08:31 — clean low-reasoning Stage A run

- Completed 81/81 Gemma respondent calls with zero malformed responses.
- The optional Gemini-low rating stage exposed two more derived-stage defects: enums needed explicit types for Vertex, and the extractor selected the embedded schema's `judgements` property instead of the returned array.
- Repair: type every enum and require requested structured-output roots to be arrays.
- Recovered 419 already captured, schema-valid judgements mechanically. No replacement rating call or changed judgement was introduced.
- Observed M1 raw recall `0.917` versus M0 `0.792`; M1 fused F1 was highest in this draw. S0 still scored `0.819`, so the packets remained easy.
- Published record: [`results/2026-08-20-stage-a-smoke-gemma/`](results/2026-08-20-stage-a-smoke-gemma/).

## 2026-08-20 09:00 — repeated correctness overlay calibration

- Ran 30 packet-blocked functional/K-2SO pairs: three packets × ten repeats, 60/60 calls successful.
- Functional mean F1: `0.790`; fictional mean F1: `0.748`.
- Paired fictional-minus-functional mean: `−0.0425`; sampling-only interval `[−0.0808, −0.0047]`.
- Outcomes: 2 fictional wins, 16 ties, 12 functional wins. Defect identity still changed in 14 pairs.
- Interpretation: wording changed behaviour, but this fictional overlay did not help this model/kernel/task set. This does not establish a universal persona effect.
- Published record: [`results/2026-08-20-persona-pair-gemma/`](results/2026-08-20-persona-pair-gemma/).

## 2026-08-20 — current pivot

- Do not rescue fictional overlays by changing outcomes or inventing post-hoc metrics. Names can remain mnemonic labels if the prose treatment loses.
- Do not escalate toward increasingly contrived synthetic defects merely to manufacture spread.
- Before ecological real-project work, migrate the active harness to Go and require macOS Seatbelt isolation with executable probes.
- Keep provider API access separate from model tool access. A local sandbox cannot prevent provider-side search tools; those must be disabled and recorded.
- Prefer an exported pre-fix source tree with Git history removed for future open-source tasks. Treat the eventual patch and tests as evidence, not the only valid answer.
- Evaluate the remaining fictional/functional role pairs as one preregistered family rather than selectively running only promising personas.

## 2026-08-20 12:23–13:00 — Go and Seatbelt migration

- Replaced the active Python harness with a standard-library Go command. Historical Python runs remain reproducible from their source tags.
- Added a mandatory deny-by-default macOS Seatbelt profile, empty child cwd, per-attempt ephemeral home/cache/temp, repository-read denial probe, external CLI digest freeze, and fail-closed macOS/root checks.
- Completed a 9/9 local mock lifecycle and reproduced the published Gemma score groups/overlaps with the Go scorer (apart from the deliberately new deterministic bootstrap resampling sequence).
- Completed live isolated checks through Pi for Gemma, Google/Gemini rating, and Anthropic/Haiku.
- Direct agy testing exposed an interactive keychain dependency under an ephemeral home; the prompts were cancelled without creating/resetting a keychain. Direct Claude CLI likewise could not use shared login state without its real home.
- Decision: do not weaken isolation to accommodate those clients. Reject direct agy/Claude before launch and route explicitly pinned Anthropic, Google, OpenAI, and Gemma models through Pi's sterile auth/model-registry copy. Do not copy Pi settings, skills, history, or sessions.
- Network remains available only to the trusted provider client transport. Current model tools remain disabled; provider-side search is still unobservable and must be a separately declared arm.

## 2026-08-20 13:08 — preserve the engineering narrative

- Added Part 3, *The Experiment Fought Back*, as a living presentation of instrumentation failures, discarded runs, scoring repairs, negative findings, and isolation hardening.
- Decision: published talks should show the failed paths and why the methodology changed, not reconstruct a falsely linear success story.

## 2026-08-21 13:20 — human blinding and ecological curation

- Clarified that one council author may be a human rater if the prior is disclosed and an independent second rater is retained.
- Found that earlier opaque IDs were hashes of a public seed: procedurally blinded, but reproducible by a motivated rater.
- Repair: private-key HMAC IDs for findings, sets, and A/B pairs; finding mapping precedes pair comparison; qualitative ratings record condition guesses and whether wording revealed treatment.
- Decision: do not rewrite fictional output into neutral prose. Wording is an outcome of the treatment; instead report imperfect blinding directly.
- Curated four pilot and four reserve pre-fix tasks merged after 1 June 2026 across Manim, dlt, Gortex, and turbovec. Froze parent/evidence commits, sanitized symptom briefs, exclusions, and task-specific evidence tests.
- Caveat: post-cutoff dates reduce training-set risk but cannot rule out provider retrieval, later fine-tuning, or hidden search. Real-task respondents must receive history-free archives with no internet/evidence access.

## 2026-08-21 14:45 — Modular extreme-task screening

- Screened `modular/modular` after it was suggested as an extreme ecological option.
- Provenance finding: Mojo was open-sourced on 18 August in a 561,408-line import. The June–August GitHub search contains no public bug-fix PR suitable for the existing PR-backed set; subsequent public history includes synchronized internal commits.
- Retained one CPU-only watchlist candidate: commit `6df03ad2…`, whose parent exhibits wrong `argmax`/`argmin` indices when a large row enters a multi-worker split-axis reduction. The evidence adds focused Mojo tests and changes 134 lines across three files.
- Kept it out of the selected set because it has no public issue/PR review and because the exact Bazel/prebuilt-Mojo dependency closure has not been frozen or run offline.
- Environment check: the screening host is macOS arm64 with 14 logical cores, but only Command Line Tools are active while Modular's development guide specifies Xcode 16+. With 40 GiB free, no speculative 225.9 MB source plus toolchain build/cache was started.
- Promotion rule: independently demonstrate the focused test failing at the exact parent and passing at the evidence commit, mirror and digest all build artifacts, then preregister an extreme-task stage. Do not use it as an outcome-dependent pilot replacement.

## 2026-08-21 15:13 — repository handoff checkpoint

- Recorded the complete session handoff in `experiment/CHECKPOINT.md` before moving future work from the `prezo` workspace to the council repository itself.
- No respondent or rating calls were made for this checkpoint.
- Priority is now ecological infrastructure: begin with a history-free Gortex Unicode snapshot, verify the regression controller-side, then design narrowly allowlisted read/search/test access under a changed and probed isolation boundary.
- The 480-call persona factorial remains unrun and is not the default next action.

## Notebook rule

Append material decisions before or immediately after their run. Record source commit/tag, config, exclusions, failures, repairs, interpretation, and next decision. Correct factual errors explicitly; do not silently rewrite earlier reasoning.
