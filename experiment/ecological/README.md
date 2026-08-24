# Ecological open-source task curation

> Status: candidate set frozen for harness integration; no respondent run has occurred.

The core set uses public fixes merged from 1 June through 20 August 2026. It samples one Python graphics project, one Python data-loading project, one Go code-intelligence project, and one Rust/Python vector-index project. A separately gated Modular task tests an extreme-complexity stratum without weakening the core set's evidence rules.

## Selected set

| Stage | Task | Project | Evidence | Size | Difficulty |
|---|---|---|---|---:|---|
| Pilot | animation z-order regression | Manim | issue 4834 / PR 4918 | +33/−3 | medium |
| Pilot | paginator stop precedence | dlt | issue 4225 / PR 4227 | +36/−2 | medium |
| Pilot | Unicode tokenizer panic | Gortex | PR 569 | +37/−3 | medium |
| Pilot | long-name temporary-file leak | turbovec | issue 488 / PR 524 | +384/−55 | medium-hard |
| Reserve | memoizer identity collision | Manim | PR 4901 | +34/−1 | hard |
| Reserve | stale import after schema restore | dlt | issue 4242 / PR 4251 | +50/−15 | hard |
| Reserve | control-plane lock starvation | Gortex | issue 479 / PR 494 | +598/−18 | hard |
| Reserve | finite calibration poisoning | turbovec | issue 478 / PR 518 | +412/−15 | hard |

[`candidates.json`](candidates.json) records exact pre-fix parent commits, evidence commits, dates, licenses, patch sizes, and relevant regression tests. [`briefs/`](briefs/) contains symptom-oriented task text with root-cause and patch instructions removed.

## Extreme reserve: Modular

The candidate `eco-modular-cpu-splitk-arg-reduce` uses the newly open-sourced [Modular Platform](https://github.com/modular/modular). A large CPU `argmax`/`argmin` can return the wrong global index only after entering a multi-worker split-axis reduction. Diagnosing it requires navigating Mojo's generic monoid state, SIMD-lane collapse, scratch representation, parallel publication, and finalisation across a 225.9 MB, 10,538-blob source tree.

It is not yet eligible for a respondent run:

- Modular open-sourced Mojo on 18 August 2026, then synchronized this fix as a public linear commit. There is no public issue or PR review, so its evidence provenance is weaker and must remain a separate stratum.
- The parent pins Mojo `1.1.0.dev2026081813`, MAX `26.6.0.dev2026081813`, and the BuildBuddy Bazel wrapper `5.0.382`. Their complete artifacts and digests must be mirrored before offline execution.
- The focused CPU test must first fail at parent `f0f9e8f…` and pass at evidence commit `6df03ad2…` on the declared macOS host. The local screening host has sufficient CPU cores but lacks the required full Xcode installation, so this validation was not attempted.

Do not substitute it after looking at pilot outputs. Promote it only by a new preregistered extreme-task stage after the three gates above pass.

## Curation procedure

On 21 August 2026 the curator queried GitHub for merged pull requests in `2026-06-01..2026-08-20`. The search returned 75 Manim, 141 dlt, 434 Gortex, and 186 turbovec PRs. Modular returned only three merged PRs: its 561,408-line open-source import, notices, and documentation. The curator therefore screened post-import public source commits separately and retained one CPU-only, test-backed fix as blocked extreme reserve—not as PR-backed evidence. Titles/labels removed releases, dependency bumps, documentation-only changes, and obvious one-line chores. For finalists, the curator inspected PR/issue or commit text, changed-file lists, patch size, exact parent, regression-test evidence, and macOS feasibility. The curator necessarily saw the fixes. Respondents must not see this directory. Outcome raters may receive the patch/tests only after outputs are frozen, in a condition-blinded evidence bundle.

This is purposive stratified sampling, not a random sample of open-source defects. It favours well-documented merged bugs with tests. Freeze that limitation and the pilot/reserve split before running models.

## Inclusion criteria

A selected task must:

- have a fix merged no earlier than 1 June 2026;
- preferably have its issue or first public report after that date;
- reproduce a material correctness, reliability, performance, or safety consequence;
- have an exact pre-fix parent commit and a merged evidence commit;
- include discriminating tests or enough maintainer evidence to construct them;
- be reviewable on macOS without private services or credentials;
- be narrow enough for one task and one scoring key;
- require diagnosis rather than repeating a solution already stated in the supplied brief.

The four pilot tasks deliberately span languages and failure shapes. Reserve tasks are harder and should not replace a failed pilot task after outputs are inspected unless the exclusion rule was frozen beforehand. The Modular task does not yet satisfy the public-review and validated-offline-toolchain gates, so it is a watchlist item rather than a selected ninth task.

## Exclusions

Strong candidates were excluded when they bundled several independent fixes, had patches too broad for an interpretable first pilot, stated the one-line fix in the public report, or relied on an issue predating the cutoff. Exclusions are recorded in `candidates.json`, including turbovec PR 326, Gortex PR 629, dlt PR 4287, and Manim PR 4936.

## Training-cutoff caveat

Post-June publication makes direct training-set inclusion less plausible for a model whose training cutoff predates June 2026. It does **not** prove ignorance:

- providers may use later fine-tuning, retrieval, caches, or web search;
- a model may already know the surrounding project and infer the fix;
- public issue/PR text can be found if internet access leaks into the condition;
- this repository itself publishes the evidence metadata for reproducibility.

Respondents therefore receive only an exported pre-fix source tree and the sanitized brief. Provider-side search must be disabled and recorded; because it remains partly unobservable, report this threat rather than claiming perfect novelty.

## Snapshot protocol

For each task:

1. fetch the recorded `parent_commit` into a controller-only cache;
2. use `git archive`, not a working clone copy, to produce a new source directory;
3. remove project-level agent instructions, MCP configuration, workflow bot prompts, issue caches, patch files, generated changelogs, and evaluation assets that could instruct or reveal a solution;
4. create a neutral source root with normalized modes/timestamps and no `.git`; keep any later edit baseline outside the respondent-visible tree;
5. store provenance and file digests outside the source root visible to the model;
6. run a deny probe for council/controller metadata, sibling tasks, network access, writes, and evidence files;
7. supply the brief through the frozen request, not as an upstream issue URL;
8. expose the same allowlisted read/search/test tools to every arm;
9. prefetch and digest the exact toolchain/package closure controller-side, record dependency licenses, then run the focused test with network disabled.

### Implemented Gortex exporter

`config/eco-gortex-unicode-tokenizer.json` freezes the upstream commits, license-file digests, explicit removals, forbidden evidence markers, Go version, and focused test. Run:

```bash
just ecological-gortex-snapshot
```

Each invocation creates a new ignored `work/<attempt>/` and never reuses an attempt in place. `source/` is the read-only respondent candidate. `controller/` contains archives, the applied evidence test, logs, manifests, provenance, and a seal; it must never be mounted for a respondent. Immutable Go/tool module closures live under ignored `cache/` and are verified by tree digest before use. The evidence test is applied only in controller validation: it must reproduce the configured failure at the parent and pass at the evidence commit.

The first development validation preserved four failed attempts before completing this lifecycle. It made no model call. Its controller tree was dirty while the exporter was being written, which required the clean rerun recorded below.

On 24 August 2026, a clean rerun from committed controller `ef48b42…` completed as ignored attempt `20260824T012106Z-eco-gortex-unicode-tokenizer-62e6fe45`. It reproduced the same source tree digest, `41a9a14b…`, verified the seal, observed the configured parent failure, passed the evidence commit, and passed source/controller/evidence/network probes. This freezes the candidate snapshot as infrastructure evidence; it is not a respondent result.

### Narrow respondent mediator

The tracked policy `tool-policy/eco-gortex-unicode-tokenizer-v1.json` freezes read/list/search/test call and byte budgets. The Go package under `mediator/` enforces source-root confinement, rejects Git/traversal/symlinks/binary or oversized text, uses deterministic bounded listing and RE2 search, correlates requests, and writes a fail-closed JSONL transcript. `pi/ecological-tools.ts` exposes only those operations plus a terminating structured review tool.

The extension does not receive a source, evidence, cache, or controller path. The doctor connects it to the trusted Go mediator with inherited request/response pipes and starts Pi with all built-in/discovered tools and context disabled. A future claim runner must retain that boundary while capturing Pi's JSON event stream and the mediator transcript. The mediator—not the provider child—reads the source and runs the predeclared test. This keeps provider transport networking separate from the network-denied test sandbox.

A controller-side Gortex test runner now verifies the snapshot and closure, runs the hidden parent regression under a new network-denied Seatbelt profile, sanitizes controller paths from its returned output, and retains raw logs outside the respondent channel. `just ecological-gortex-mediator-check <snapshot>` exercises positive list/read/search/test requests plus traversal and arbitrary-target denials without a model call.

`just ecological-gortex-pi-doctor <snapshot>` now starts Pi in no-model RPC mode under the new provider-child profile, disables built-ins and discovered resources, exercises the inherited mediator pipes through a controller-only extension command, explicitly sets and verifies model/thinking state, checks source/controller/evidence denial, and seals its event/transcript/profile/runtime digests. Pi RPC is intentionally persistent, so the doctor uses controlled termination after all responses are captured; a claim runner should use JSON/print mode, which exits after the agent settles.

The mediator, extension, and doctor remain infrastructure only until a claim runner assembles the frozen brief/system prompt, captures the final structured submission and usage, and seals the complete lifecycle. Do not start a model by loading the extension manually. The first successful checks were run from a dirty development tree and must be repeated after commit before they can freeze runner inputs.

The answer evidence normally consists of the merged patch, its regression tests, issue/maintainer discussion, and independent human review. The Modular watchlist task has commit/test evidence but no public review discussion; report and analyse that stratum separately. An upstream correction is evidence, not the only acceptable answer: supported novel findings remain valid and must be adjudicated.

## Scoring

The detailed condition-blinded 1–5 anchors, claim-supportedness rules, efficiency measures, and adjudication procedure are frozen in [`SCORING.md`](SCORING.md). Task-specific keys under `evidence/` are rater/controller-only and must never enter a respondent mount or provider child.

Score separately:

- root-cause localisation;
- material consequence;
- proposed correction correctness and scope;
- regression-test discrimination;
- unsupported claims;
- remedy quality/actionability;
- token, latency, and tool-use cost.

Do not reduce ecological tasks to planted-defect F1 alone. Use blinded finding-level mappings plus paired qualitative ratings. Two humans rate independently; a rater may know the council project but must not see condition labels, prompt variants, PRs, or the private unblind map.
