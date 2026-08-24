# Experiment checkpoint

> Checkpoint: 24 August 2026, after the clean Gortex snapshot rerun and first ecological mediator/rubric implementation. Start future sessions in the `council-of-nark` repository root and read this file plus [`LAB_NOTEBOOK.md`](LAB_NOTEBOOK.md) before changing or running the experiment.

## Repository state

- The snapshot work started from clean `main` at `69638aea2c429b64a43886098fde1f3cd8bee3f9`, synchronized with `origin/main`.
- Controller `ef48b42…` contains the Gortex exporter; `0ba767a…` adds the mediator/rubric/doctor boundary and `a5972ab…` completes raw-artifact sealing. Confirm the live commit and clean state before any further snapshot or respondent work.
- Strict Go/Seatbelt harness tag: `experiment-harness-go-v1.0`.
- Maintained platform: macOS, Go 1.22+, `/usr/bin/sandbox-exec`, isolated Pi adapters. Direct agy and direct Claude CLI remain fail-closed.
- Published synthetic results remain calibration, not confirmatory evidence.
- No persona-factorial, ecological respondent, or ecological rating call has occurred.

Confirm the live state rather than assuming a recorded hash is still HEAD:

```bash
git status --short --branch
git log -5 --oneline --decorate
just audit
just experiment-test
just experiment-sandbox-check
```

## What is complete

### Synthetic calibration

- Published the 81-call Haiku Stage A smoke, clean 81-call Gemma Stage A smoke, and 30-pair Gemma correctness wrapper calibration.
- Preserved instrumentation failures, the discarded partial Gemma run, parser recovery, metric definitions, and chronology.
- The repeated correctness result favoured the functional prompt for that model/task family. It is sampling calibration over three synthetic packets, not task generalisation.
- A balanced eight-role fictional-overlay factorial is designed in [`PERSONA_FACTORIAL.md`](PERSONA_FACTORIAL.md) and `config/persona-factorial-gemma.json`, but its 480 calls have **not** been run.

### Isolation and reproducibility

- The active harness is standard-library Go and mandatory macOS Seatbelt.
- Provider children use empty working directories and ephemeral homes and cannot read the council repository, prompt-assembly worktree, answer keys, or real home.
- External CLI versions/digests, retries, sealing, verification, health, blinding, judging, and scoring are implemented.
- Current provider/model tools are disabled. Provider transport networking remains necessary; provider-side hidden search is an unobservable limitation.

### Human rating and blinding

- PR #18 replaced public-seed opaque-looking IDs with private random-key HMAC-SHA-256 IDs for sets, findings, and A/B pairs.
- Phase 1 maps shuffled findings independently. Phase 2 compares matched outputs with randomized left/right placement and rates supportedness, actionability, fix quality, preference, condition guess, confidence, and wording leakage.
- Raw wording is preserved because it is part of the treatment. This is label-blinded, not guaranteed treatment-blinded.
- `qualitative` requires two complete independent raters before unblinding and preserves opaque IDs in the derived report.
- The project author volunteered as one disclosed-prior rater. A second independent rater is still mandatory. No claim-bearing human ratings have occurred.
- A detailed anchored 1–5 remedy-quality rubric and disagreement/adjudication procedure still need to be frozen before a claim-bearing run.

### Ecological curation

The core set has four pilot and four reserve tasks across Manim, dlt, Gortex, and turbovec. Exact parents, evidence commits, tests, dates, licenses, patch sizes, exclusions, and symptom-only briefs are in [`ecological/candidates.json`](ecological/candidates.json).

Pilot:

1. Manim animation z-order regression.
2. dlt paginator stop precedence.
3. Gortex Unicode tokenizer panic.
4. turbovec long-name stale-temporary-file leak.

Reserve:

1. Manim memoizer identity collision.
2. dlt stale imported-schema restoration.
3. Gortex control-plane lock starvation.
4. turbovec finite calibration poisoning.

PR #19 added a separately gated Modular extreme reserve: CPU split-axis `argmax`/`argmin` index corruption. It is **not eligible** until its exact external build closure is mirrored/digested and the focused test is shown failing at the parent and passing at the evidence commit. It has public commit/test evidence but no public issue or PR review, so never pool it silently with the PR-backed set.

The original local development export completed with source tree digest `41a9a14b…`; the exact evidence test failed at the parent with the expected panic and passed at the evidence commit using the frozen offline Go/module closure. A clean rerun from committed controller `ef48b42…` completed on 24 August as ignored attempt `20260824T012106Z-eco-gortex-unicode-tokenizer-62e6fe45`, reproduced the same tree digest, verified its seal, and passed all configured isolation probes. It is frozen infrastructure evidence, not a respondent result. No ecological respondent call or outcome rating has occurred.

## Highest-priority next work

### 1. Complete and probe the ecological launcher

Do this before spending 480 calls on the synthetic persona factorial. The clean snapshot, tool policy, Go mediator/test runner, Pi extension, and no-model doctor now exist. No claim runner or ecological model call exists yet.

The controller-side test runner and `just ecological-gortex-mediator-check <snapshot>` now verify the sealed snapshot/closure, exercise bounded list/read/search, reproduce the hidden regression under network-denied Seatbelt, sanitize returned paths, and deny traversal/arbitrary targets. The first successful check (`20260824T014109Z-gortex-mediator-check-2324d548`) came from a dirty development tree and is engineering evidence only.

`just ecological-gortex-pi-doctor <snapshot>` now probes the other half of the boundary without a provider call. It starts isolated Pi with no built-ins/discovered resources/session, denies source/controller/evidence reads, exercises inherited mediator pipes, explicitly verifies model/thinking state, and seals events/profiles/runtime/transcript digests. Four failed attempts exposed persistent RPC shutdown behavior and Pi clamping Gemma `off` to `minimal`; the repair pins explicit `minimal` and uses controlled post-check termination.

After commits `0ba767a…` and `a5972ab…`, clean mediator attempt `20260824T020111Z-gortex-mediator-check-662f680b` and clean Pi doctor `20260824T020121Z-gortex-pi-doctor-a93ff02e` both completed from `a5972ab…`. Their summaries/raw artifacts and seals verified, the returned mediator responses contained no absolute controller path, and the doctor made zero provider turns. These are frozen infrastructure checks, not respondent results.

The claim runner must:

1. use the verified snapshot/policy/mediator/extension/runtime inputs and JSON/print mode, which exits after the agent settles;
2. assemble and digest the exact system prompt, sanitized brief, tool schemas, model, thinking level, and decoding state;
3. capture Pi events and mediator transcript, validate exactly one terminating structured submission, and seal responses, usage, cost, latency, profiles, and probes;
4. rerun negative direct-read/Git/write/arbitrary-target/shell/pipe/test-network probes from a clean committed controller;
5. rerun both no-model checks after committing this infrastructure, then run a deterministic mock lifecycle before any ecological model call.

Do not load `pi/ecological-tools.ts` manually for a respondent.

### 2. Finish ecological scoring/run design before calls

[`ecological/SCORING.md`](ecological/SCORING.md) now freezes claim-supportedness rules, five anchored dimensions, efficiency measures, missing-output handling, and disagreement/adjudication. `ecological/evidence/eco-gortex-unicode-tokenizer.md` freezes the task-specific evidence key and acceptable equivalent remedies. Both raters must remain condition-blinded; the evidence directory is controller/rater-only.

Still freeze:

- compared arms and the exact byte-level prompt factor;
- call count, randomisation/blocking, model/thinking level, retries, and smallest effect of interest;
- the cross-rater aggregation and any primary endpoint derived from the five separate dimensions;
- two available raters and the rating CSV/bundle format;
- exclusions and pilot-to-reserve substitution rules.

The upstream patch remains evidence, not the only acceptable solution. Preserve supported novel findings.

### 3. Decide whether the persona factorial is still worth 480 calls

Do not run `just experiment-persona-factorial-gemma 2` by default. First confirm:

- the anchored remedy rubric is committed;
- two raters are available;
- all eight roles will be reported, including null/negative outcomes;
- the ±0.02 family margin and Holm-corrected role claims remain preregistered;
- synthetic mechanism evidence is worth the cost relative to ecological integration.

If run, rerun correctness within the frozen family; do not pool the observed pilot selectively.

## Non-negotiable rules

- Never edit, patch, combine, or count excluded sealed samples.
- Every result report defines TP, FP, FN, precision, recall, F1, and macro mean inline, including not-applicable explanations.
- LLM ratings are exploratory only; claim-bearing work requires two independent humans.
- Preserve failures, negative results, deviations, costs, and methodological pivots in [`LAB_NOTEBOOK.md`](LAB_NOTEBOOK.md).
- Keep functional kernels byte-identical when fictional prose is the factor.
- Keep answer keys/evidence out of respondent and fusion contexts.
- Separate raw union coverage from fused practical output.
- Unrestricted internet search is not a base ecological capability. A future controlled-search arm needs a separate preregistration and threat model.
- Do not weaken ephemeral-home isolation to support shared-login provider clients.

## Useful paths and commands

- Ecological protocol: [`ecological/README.md`](ecological/README.md)
- Candidate manifest: [`ecological/candidates.json`](ecological/candidates.json)
- Human/operator procedure: [`RUNSHEET.md`](RUNSHEET.md)
- Harness design: [`harness/README.md`](harness/README.md)
- Isolation model: [`ISOLATION.md`](ISOLATION.md)
- Metrics: [`METRICS.md`](METRICS.md)
- Preregistration: [`PREREGISTRATION.md`](PREREGISTRATION.md)

```bash
just experiment-doctor experiment/config/persona-factorial-gemma.json
just experiment-adapter-check-gemma
just experiment-verify "$RUN"
just experiment-bundle "$RUN"
just experiment-score "$RUN" blinded/ratings-adjudicated.csv adjudicated
just experiment-qualitative "$RUN" blinded/pairwise-ratings-both.csv qualitative
```

## Ignored local artifact warning

`experiment/runs/20260821T135400Z-mock-pair-fcb1a5efea` is a two-call mock plumbing run created while testing HMAC pairing. It contains no respondent findings and any locally fabricated qualitative rows/derived analyses are plumbing fixtures, not human evidence. Do not report or pool it. Historical real runs remain identified in the published result manifests and notebook.

The ignored `experiment/ecological/work/` and `experiment/ecological/cache/` trees contain the Gortex controller cache, four preserved failed exporter attempts, and the successful dirty-tree development attempt `20260821T160711Z-eco-gortex-unicode-tokenizer-f3a5b8c9`. They contain upstream evidence and controller metadata. Never mount or publish them as respondent source.
