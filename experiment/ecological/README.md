# Ecological open-source task curation

> Status: candidate set frozen for harness integration; no respondent run has occurred.

This set uses public fixes merged from 1 June through 20 August 2026. It samples one Python graphics project, one Python data-loading project, one Go code-intelligence project, and one Rust/Python vector-index project.

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

## Curation procedure

On 21 August 2026 the curator queried GitHub for merged pull requests in `2026-06-01..2026-08-20`. The search returned 75 Manim, 141 dlt, 434 Gortex, and 186 turbovec PRs. Titles/labels removed releases, dependency bumps, documentation-only changes, and obvious one-line chores. For finalists, the curator inspected PR/issue text, changed-file lists, patch size, exact merge parent, regression-test evidence, and macOS feasibility. The curator necessarily saw the fixes. Respondents must not see this directory. Outcome raters may receive the patch/tests only after outputs are frozen, in a condition-blinded evidence bundle.

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

The four pilot tasks deliberately span languages and failure shapes. Reserve tasks are harder and should not replace a failed pilot task after outputs are inspected unless the exclusion rule was frozen beforehand.

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
3. remove project-level agent instructions, MCP configuration, workflow bot prompts, issue caches, patch files, and generated changelogs that could instruct or reveal a solution;
4. initialise one neutral root commit with no remote, reflog, alternate object store, tags, branches, or upstream refs;
5. store provenance and file digests outside the source root visible to the model;
6. run a deny probe for `.git` history, sibling directories, the council repository, network tools, and evidence files;
7. supply the brief through the frozen request, not as an upstream issue URL;
8. expose the same allowlisted read/search/test tools to every arm.

The answer evidence consists of the merged patch, its regression tests, issue/maintainer discussion, and independent human review. The upstream patch is evidence, not the only acceptable answer: supported novel findings remain valid and must be adjudicated.

## Scoring

Score separately:

- root-cause localisation;
- material consequence;
- proposed correction correctness and scope;
- regression-test discrimination;
- unsupported claims;
- remedy quality/actionability;
- token, latency, and tool-use cost.

Do not reduce ecological tasks to planted-defect F1 alone. Use blinded finding-level mappings plus paired qualitative ratings. Two humans rate independently; a rater may know the council project but must not see condition labels, prompt variants, PRs, or the private unblind map.
