# Experiment checkpoint

> Checkpoint: 21 August 2026, after PRs #18 and #19. Start future sessions in the `council-of-nark` repository root and read this file plus [`LAB_NOTEBOOK.md`](LAB_NOTEBOOK.md) before changing or running the experiment.

## Repository state

- `main` was clean and synchronized with `origin/main` at checkpoint creation.
- Last substantive integration before this checkpoint: PR #19, baseline merge `9221554fdd3b1f1ba1ec17cfc00c04c364ce7c9a`.
- Strict Go/Seatbelt harness tag: `experiment-harness-go-v1.0`.
- Maintained platform: macOS, Go 1.22+, `/usr/bin/sandbox-exec`, isolated Pi adapters. Direct agy and direct Claude CLI remain fail-closed.
- Published synthetic results remain calibration, not confirmatory evidence.

Confirm the live state rather than assuming this hash is still HEAD:

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

No ecological source export, respondent call, or outcome rating has occurred.

## Highest-priority next work

### 1. Build the ecological snapshot pipeline

Do this before spending 480 calls on the synthetic persona factorial.

Start with the Gortex Unicode pilot because it has a focused Go test, no special hardware, and a relatively conventional toolchain. The exporter should:

1. fetch the exact parent into a controller-only cache;
2. verify commit identity and license;
3. export with `git archive` rather than copying a checkout;
4. remove `.git`, remotes, reflogs, patches, changelogs, issue caches, agent instructions (`AGENTS.md`, `CLAUDE.md`, `.cursor`, MCP/project AI configuration), and solution hints while retaining required licenses/notices;
5. create a neutral, history-free source root and store provenance/digests outside it;
6. prove the source cannot read the council repository, sibling tasks, evidence commit, or controller metadata;
7. prefetch and digest the dependency/test closure, then prove focused tests run with network denied;
8. apply the evidence regression test to the parent controller-side to demonstrate expected failure, and run it at the evidence commit to demonstrate pass;
9. preserve every failed export/build attempt rather than repairing artifacts in place.

Do not expose upstream evidence, `ecological/candidates.json`, or this repository to respondents.

### 2. Design explicit ecological tool access

The current one-shot harness disables tools and denies provider children all repository access. A real-source ecological task therefore needs a new, published isolation boundary; do not silently grant a shell or reuse the synthetic prompt path.

Preferred design questions to resolve:

- allow read-only access only to the sanitized task snapshot, never to the council/controller roots;
- expose narrow read/search operations plus predeclared test targets, not an unrestricted shell;
- execute tests in a separate network-denied sandbox with frozen dependencies while retaining network only for provider transport;
- log every tool request/result and include tool-policy/profile digests in `request.json` and the seal;
- keep tool budgets identical across compared arms;
- add executable probes for history, remotes, evidence, sibling paths, network, writes, and hidden solution search.

If Pi extensions/custom tools are used, first read the installed Pi extension and TUI/tool documentation and examples completely. Publish the changed threat model in [`ISOLATION.md`](ISOLATION.md).

### 3. Freeze ecological scoring before calls

Create a task-specific evidence key and anchored human rubric for:

- root-cause localisation;
- material consequence;
- correction correctness and scope;
- regression-test discrimination;
- unsupported claims;
- remedy quality/actionability;
- tool, token, latency, and cost efficiency.

The upstream patch is evidence, not the only acceptable solution. Preserve supported novel findings. Require two independent condition-blinded human raters and freeze exclusions/adjudication before output inspection.

### 4. Decide whether the persona factorial is still worth 480 calls

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
