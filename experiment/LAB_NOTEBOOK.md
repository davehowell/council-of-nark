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

## 2026-08-21 15:57–16:08 — Gortex ecological snapshot development validation

- Added the first ecological exporter and task config for `eco-gortex-unicode-tokenizer`. It fetches the exact parent/evidence objects, verifies commit ancestry and Apache license/notice digests, uses `git archive`, removes configured agent/evaluation/benchmark assets, normalizes the source, and leaves `.git` absent. Provenance and evidence remain outside the respondent-visible source.
- Froze Go `1.26.5` plus the focused package's six-module closure and dependency license files. Controller validation applies the exact evidence test to the parent and runs `go test -count=1 ./internal/search/rerank -run '^TestTokenize$'` with module download and VCS stamping disabled.
- Preserved four failed attempts: validation extraction initially omitted an intermediate directory; macOS refused to move a read-only module-cache root; cleanup of a duplicate read-only toolchain staging tree failed; then Go attempted VCS stamping from the enclosing council checkout. Repairs were limited to new code/new attempts: create extraction parents, atomically stage the module closure with a temporary writable root, make duplicate staging directories removable, and set `-buildvcs=false`.
- The fifth attempt, `20260821T160711Z-eco-gortex-unicode-tokenizer-f3a5b8c9`, completed. The sanitized source has 4,093 entries / 37,789,193 bytes, tree digest `41a9a14b…`; the parent reproduced `index out of range [1] with length 1`, and the evidence commit passed. Source/controller/evidence/network probes passed. No model or rating call occurred.
- This was a development run from a dirty controller tree (`69638ae…`), not a claim-bearing frozen snapshot. Rerun after committing the exporter. Next, define the narrow respondent read/search/test mediator and ecological scoring rubric; do not run the persona factorial by default.

## 2026-08-24 01:21–01:54 — clean Gortex freeze and mediator boundary

- Ran the required audit, Go tests, and synthetic Seatbelt probe from clean committed controller `ef48b42…`; all passed. Clean snapshot attempt `20260824T012106Z-eco-gortex-unicode-tokenizer-62e6fe45` completed in 39.6 seconds with `tree_dirty: false`.
- Verified `status.json`, the provenance and source-manifest seal digests, and source tree digest `41a9a14bd447791f254c79a7441d81451dd69ce6e03367dddbc20c4679c2b523` (4,093 entries / 37,789,193 bytes). The parent produced the configured panic, evidence passed, the frozen closure matched, and source/controller/evidence/write/network probes passed. No model or rating call occurred.
- Read the installed Pi README, extension documentation, TUI/tool documentation, JSON-event documentation, and relevant custom-tool/sandbox/structured-output examples before choosing an integration boundary.
- Decision: do not mount source into the network-enabled Pi provider child and do not override a broad built-in read or shell tool. The Pi extension receives only two inherited JSON pipes. A trusted Go mediator owns the sanitized source, strict path checks, deterministic bounded list/read/RE2 search, tool/result budgets, and a fail-closed transcript. Focused tests remain a separate controller operation under network-denied Seatbelt.
- Added the versioned Gortex tool policy, mediator implementation/tests, Pi custom tools, terminating structured review contract, task evidence key, and anchored ecological scoring/adjudication rubric. The upstream patch remains evidence rather than the only acceptable answer; supported novel findings survive.
- Added snapshot seal/tree verification and a controller-side Gortex test runner. It verifies the closure, runs the hidden parent test under network-denied Seatbelt, sanitizes controller paths from returned output, and retains raw logs controller-side. Development check `20260824T014109Z-gortex-mediator-check-2324d548` passed bounded list/read/search/test operations plus traversal and arbitrary-target denial; respondent-visible responses contained no absolute controller path. The 13.2-second check made no model call.
- Added a no-model Pi RPC doctor to connect the inherited pipes and probe the provider-child profile. Preserved four failed attempts: stdin EOF did not terminate persistent RPC; explicit `off` was accepted but clamped back to Gemma's `minimal`; extension `ctx.shutdown()` did not end the idle RPC process; and SIGTERM followed the same persistent path. Repair: pin the observed supported `minimal` level, verify it via RPC before any prompt, and use controlled SIGKILL only after doctor responses are complete. A future claim runner will use JSON/print mode, which exits after the agent settles.
- Dirty-tree doctor attempt `20260824T015407Z-gortex-pi-doctor-ceb9c4b8` completed in 3.2 seconds with zero provider turns and extension errors. It selected `gemma-4-31b-it`/`minimal`, exercised one mediated pipe request, denied source/controller/council/evidence reads, allowed only scratch writes, recorded Node `v22.19.0` and Pi `0.84.2` entrypoint digests, and verified its doctor seal. No model or rating call occurred.
- Both successful checks ran from the current dirty development tree, so they are engineering evidence only and must be repeated after commit. Next: implement the JSON-mode claim runner, frozen prompt/request assembly, final-submission validation, full sealing, and a deterministic mock lifecycle. Do not make an ecological model call before that runner and the compared-arm design are frozen.

## 2026-08-24 02:01 — committed ecological boundary checks

- Committed the mediator/rubric/doctor boundary as `0ba767a…`, then added complete summary/raw-artifact sealing as `a5972ab…`. The repository was clean for both checks below; ignored prior attempts were retained.
- Clean mediator attempt `20260824T020111Z-gortex-mediator-check-662f680b` completed from `a5972ab…`. It reverified snapshot `41a9a14b…`, closure and policy/extension/mediator digests; passed bounded list/read/search and the hidden network-denied regression; denied traversal and an arbitrary test target; exposed no absolute controller path in returned responses; and sealed the summary, transcript, and test-artifact tree.
- Clean Pi doctor `20260824T020121Z-gortex-pi-doctor-a93ff02e` completed from `a5972ab…`. It passed provider-child source/controller/council/evidence denial, inherited-pipe health, explicit `gemma-4-31b-it`/`minimal` state, zero-provider-turn and zero-extension-error checks, and sealed the doctor plus raw event, stderr, probe, profile, config, runtime, and transcript digests.
- No model or rating call occurred. The next gate is still a JSON-mode claim runner with frozen prompt assembly, exactly-one-final-submission validation, mock lifecycle, and complete usage/cost/latency sealing. Compared arms and human-rating aggregation must be preregistered before the first ecological respondent call.

## 2026-09-07 — close the evidence review and shorten the talks

- Recomputed all three published calibration CSVs with `scripts/evidence_audit.py`, verifying their published-file manifest digests and F1 arithmetic. This checks published artifacts, not the unpublished raw seals; no new respondent or rating call occurred.
- New descriptive persona sensitivity: every leave-one-task-out mean stays negative, from −0.03364 to −0.04892. A post-hoc whole-task sign-flip check has only eight assignments and two-sided p=0.25 under its symmetry assumption. It illustrates that three tasks do not support broad inference; it does not invalidate the historical sampling-only interval.
- Corrected the cost narrative: fictional total tokens were 40,182 versus functional 41,000, around 2% fewer. Detection underperformed, but token waste was not demonstrated. The Haiku repeated-omnibus/S1 recorded cost ratio recomputes to 8.924 at equal mean F1, without establishing equivalence.
- Added `CONCLUSION.md` and a proposed finite `FINAL_ROUND.md`: single review, continuing self-review, repeated independent omnibus, three functional specialists, and the same specialist kernels with character overlays. Proposed resource matching, task-level analysis, explicit equivalence margins, human aggregation, precision planning and a fixed stopping rule. No claim-bearing launch is authorized by that document alone: exact prompts, launcher, eligibility checks, human raters and final preregistration are unfinished.
- Rewrote the talks into four six-slide presentations with one shared theme, complete timed scripts and per-slide sources. Moved engineering detail into the retained notebook/timeline. Part 4 explicitly labels the final round as proposed and unrun.
- Corrected the Pages Haiku tile: 0.752 was fused F1, not fused recall (0.708). Historical result files remain unchanged.
- Pages now builds on PRs and regenerates PDFs and readable notes from slide sources before publishing on main. Prior presentations remain recoverable in Git history at `a76d9ad` and earlier.

## Notebook rule

Append material decisions before or immediately after their run. Record source commit/tag, config, exclusions, failures, repairs, interpretation, and next decision. Correct factual errors explicitly; do not silently rewrite earlier reasoning.
