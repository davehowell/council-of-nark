# Experiment checkpoint

> 15 September 2026 update: launch gate 1 in [FINAL_ROUND.md](FINAL_ROUND.md) is complete. A committed JSON-mode respondent launcher now captures exact provider payloads, enforces declared resource limits, validates one terminating submission, removes ephemeral credentials, and seals the whole attempt. Clean snapshot, mediator, no-provider Pi doctor, and deterministic mock checks passed. The work found and repaired an overbroad Seatbelt process permission before any ecological respondent call. Gates 2–4, two named raters, and spending approval remain unresolved; the closing round is still not authorized.

> Checkpoint: 15 September 2026, after the ecological launcher and repaired boundary passed from committed code. Start future sessions in the `council-of-nark` repository root and read this file plus [`LAB_NOTEBOOK.md`](LAB_NOTEBOOK.md) before changing or running the experiment.

## Repository state

- The snapshot work started from clean `main` at `69638aea2c429b64a43886098fde1f3cd8bee3f9`, synchronized with `origin/main`.
- Controller `a000459…` adds the respondent launcher and repaired executable boundary; `fc9b6d4…` pins its clean Gortex snapshot, and `d22e3c8…` extends the Pi doctor with executable-denial probes. Confirm the live commit and clean state before any further snapshot or respondent work.
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
- Provider children use empty working directories and ephemeral homes and cannot read the council repository, prompt-assembly worktree, answer keys, or real home. Seatbelt permits process creation separately from explicitly allowlisted executables; probes require an unlisted executable to fail.
- External CLI versions/digests, retries, sealing, verification, health, blinding, judging, and scoring are implemented.
- Current provider/model tools are disabled. Provider transport networking remains necessary; provider-side hidden search is an unobservable limitation.

### Human rating and blinding

- PR #18 replaced public-seed opaque-looking IDs with private random-key HMAC-SHA-256 IDs for sets, findings, and A/B pairs.
- Phase 1 maps shuffled findings independently. Phase 2 compares matched outputs with randomized left/right placement and rates supportedness, actionability, fix quality, preference, condition guess, confidence, and wording leakage.
- Raw wording is preserved because it is part of the treatment. This is label-blinded, not guaranteed treatment-blinded.
- `qualitative` requires two complete independent raters before unblinding and preserves opaque IDs in the derived report.
- Dave is confirmed as the disclosed-prior rater. Dan is the prospective independent rater but is not confirmed yet. A second independent rater remains mandatory. No claim-bearing human ratings have occurred.
- The Gortex task has an anchored 1–5 quality rubric and disagreement/adjudication procedure. Equivalent anchors must be frozen separately for every later task.

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

The original local development export completed with source tree digest `41a9a14b…`; the exact evidence test failed at the parent with the expected panic and passed at the evidence commit using the frozen offline Go/module closure. After the executable-policy repair, clean attempt `20260915T032052Z-eco-gortex-unicode-tokenizer-f123f261` from controller `a000459…` reproduced that tree and closure, denied network and an unlisted executable during validation, and verified its seal. It is frozen infrastructure evidence, not a respondent result. No ecological respondent call or outcome rating has occurred.

## Highest-priority next work

### 1. Freeze the comparative Gortex block

Launch gate 1 is complete. Clean mediator attempt `20260915T032205Z-gortex-mediator-check-8ec5ecdb`, Pi doctor `20260915T032313Z-gortex-pi-doctor-2c85ae4d`, and deterministic mock `20260915T032315Z-eco-gortex-unicode-tokenizer-respondent-mock-8006b0b4` passed from committed code. The mock sealed and independently verified 12 files with digest `ce298226…`; it used zero provider calls. See [`ecological/respondent/README.md`](ecological/respondent/README.md).

Before a model call, freeze the exact prompts and deterministic schedule for S/I/R/F/P, define how continuing-session and fuser histories are represented, and add a whole-block controller rather than invoking the single-stage launcher by hand. Dave is confirmed as the disclosed-author rater; confirm Dan or another independent rater, set the deadline and model-specific spending ceiling, and document the excluded feasibility record. The current no-model doctor selects `google/gemma-4-31b-it`/`minimal`, but the live comparison model is not frozen. Do not load `pi/ecological-tools.ts` manually for a respondent.

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
