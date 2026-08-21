# Experiment

This directory contains the protocol and frozen source material for a staged empirical study of the Council of Nark.

## Platform support

The maintained runner supports **macOS only**. Its isolation design uses macOS Seatbelt and is not silently weakened on another operating system. Linux, Windows, container, and VM ports are left as an exercise for replicators; publish the changed threat model and rerun the isolation probes when adapting it.

## Research questions

1. Do independent functional specialists find more planted defects than the same number of repeated omnibus reviews?
2. Does a short fictional wrapper add useful findings beyond a byte-identical functional lens kernel?
3. Does one fusion call improve precision without discarding enough true findings to reduce F1?
4. Does independent fan-out outperform an informed serial chain across all reviewer orders?
5. Do any effects replicate across model providers?

The study measures model output. It cannot infer which internal features, activations, or weights caused an output difference.

## Material

- [`CHECKPOINT.md`](CHECKPOINT.md): current handoff, completed work, blockers, non-negotiable rules, and prioritized next steps.
- [`protocol.md`](protocol.md): hypotheses, arms, topology, scoring, sampling, threats, and interpretation rules.
- [`prompts/`](prompts/): controlled single-reviewer, specialist, and fuser prompts.
- [`scenarios/`](scenarios/): three synthetic packets and their planted-defect keys.
- [`config/`](config/): pinned smoke, pilot, topology, provider-pair, and judge configurations.
- [`schema/`](schema/): strict response contracts.
- [`harness/`](harness/): standard-library Go runner, macOS Seatbelt profiles, adapters, blinding, sealing, rating, and scoring code.
- [`HARDNESS.md`](HARDNESS.md): scenario difficulty dimensions, smoke calibration triggers, and confirmatory mutation rules.
- [`METRICS.md`](METRICS.md): F1, semantic union, overlap, and comparability definitions.
- [`PERSONA_FACTORIAL.md`](PERSONA_FACTORIAL.md): balanced eight-role functional/fictional follow-up, estimands, multiplicity, and decision rule.
- [`ecological/`](ecological/): eight PR-backed post–June 2026 tasks across Manim, dlt, Gortex, and turbovec, plus one blocked Modular extreme-reserve candidate.
- [`CONTAMINATION_REVIEW.md`](CONTAMINATION_REVIEW.md): post-smoke review of context boundaries, scoring, scheduling, and remaining threats.
- [`ISOLATION.md`](ISOLATION.md): macOS Seatbelt, network, dedicated-account, and real-project threat model.
- [`LAB_NOTEBOOK.md`](LAB_NOTEBOOK.md): chronological engineering decisions, failures, repairs, results, and pivots.
- [`PREREGISTRATION.md`](PREREGISTRATION.md): commit-before-run hypothesis, threshold, sample, rating, and exclusion template.
- [`RUNSHEET.md`](RUNSHEET.md): the operator and human-rating procedure.

Each reviewer receives only `review-packet.md`. The harness rejects a request containing answer-key markers or claims. The keys are public in this repository, so future users must treat this as a reproducible demonstration rather than a permanently hidden benchmark.

## Stages

- **Smoke:** one sample per packet and arm. Repair broken schemas and ceiling/floor effects. Do not make inferential claims.
- **Pilot:** ten samples per packet and Stage A arm; three samples per packet, topology, and chain order. Estimate effects and variance.
- **Confirmatory extension:** freeze semantics-preserving variants and clean controls before choosing thresholds or inspecting outputs. Use blinded human ratings for ambiguous matches.

The three base packets are too few for broad claims. Repeated sampling estimates model variance, not task diversity.

The first successful Stage A calibration is published under [`results/2026-08-20-stage-a-smoke/`](results/2026-08-20-stage-a-smoke/). It proved the plumbing and exposed limited headroom with a strong model. Its one-sample descriptive scores are not a decision about council value.

The unchanged low-reasoning Gemma rerun is under [`results/2026-08-20-stage-a-smoke-gemma/`](results/2026-08-20-stage-a-smoke-gemma/). Specialists increased raw recall and conditions found different defect IDs, although S0 remained high. A focused 30-pair correctness test under [`results/2026-08-20-persona-pair-gemma/`](results/2026-08-20-persona-pair-gemma/) found that the K-2SO wrapper reduced mean F1 for this model and task set. These remain calibration results, not a general verdict on personas or councils.

## Freeze rules

Before a claim-bearing run:

1. commit every prompt, packet, key, schema, model identifier, decoding setting, exclusion rule, and analysis script;
2. run the public-release audit;
3. generate a manifest with the Git commit and SHA-256 digest of every input;
4. randomise run order in blocks;
5. assemble each prompt in a clean detached worktree, then start the provider in a fresh Seatbelt-confined process with an empty cwd and ephemeral home; deny the child any repository/worktree access;
6. capture raw request, raw response, model metadata, external CLI digests, Seatbelt profile digest, usage, latency, errors, and retries;
7. make the completed run immutable by digest and verify it before analysis.

A clean process prevents local context leakage. It cannot prove that a provider has no unobserved server-side state. Record that limitation.

## Scoring

A finding maps to zero or one planted defect. Semantic duplicates of the same true or false claim count once. Unmatched claims form blinded false-positive clusters before scoring.

The primary metric is macro F1 over planted defects, weighted equally by packet. Also report precision, recall, p10, worst case, variance, token/cost efficiency, raw-union coverage, fusion retention, malformed responses, and latency.

Use deterministic location checks where possible. Two blinded humans decide ambiguous semantic mappings; an LLM judge may triage but must not be the final confirmatory rater. One human may be the council's author if that prior is disclosed and an independent second rater is retained. HMAC-keyed IDs hide labels; raw wording is not rewritten, so measure and report treatment-guess accuracy instead of claiming perfect blinding.

## Run the harness

Preflight does not call a model:

```bash
just experiment-test
just experiment-sandbox-check
just experiment-doctor experiment/config/stage-a-smoke-gemma.json
```

After the harness and assets are committed on a clean tree, verify explicit model selection and run the low-reasoning Stage A smoke:

```bash
just experiment-adapter-check-gemma
just experiment-stage-a-smoke-gemma 2
```

The Stage A smoke makes 81 calls: 3 packets × (S0/S1/S2 at one call each + M0/M1/M2 at seven reviewers and one fuser each). The focused Gemma correctness pair makes 60 calls. The preregistration-ready eight-role persona factorial makes 480 calls: 8 roles × 3 packets × 10 repeats × 2 wrappers. The topology smoke makes 144 calls across all six role orders. The cross-provider pair smoke makes 18 calls.

If the cheap smoke restores useful spread, estimate one wrapper effect against sampling noise before redesigning scenarios:

```bash
just experiment-persona-pair-gemma 2
```

Use the generic recipe for another frozen config:

```bash
just experiment-complete experiment/config/provider-pair-smoke.json 3
```

Every prompt is assembled from a fresh detached worktree at the recorded commit. The model process then runs from an empty directory and ephemeral home under a deny-by-default Seatbelt profile; it cannot read that worktree or repository. Adapters disable tools, sessions, extensions, skills, and project context. Provider transport networking remains available, while model web/shell tools remain disabled. The runner stores raw attempts under ignored `experiment/runs/`, then seals the raw file set by SHA-256 digest.

See [`RUNSHEET.md`](RUNSHEET.md) for blinded ratings and scoring. A smoke-only arm-blinded LLM triage recipe is available, but confirmatory scoring requires two independent human raters.

## Model-state limitation

Fresh local sessions prevent this harness, earlier council calls, local memory files, and the current worktree from entering model context. They do not prove that a provider has no server-side personalisation, caching, abuse classifiers, or other hidden state. The first Stage A smoke used a pinned Anthropic snapshot; the headroom calibration explicitly selects Gemma through Pi rather than inheriting Pi's powerful default. Stage B treats provider as an experimental factor and reports hidden-system differences as a limitation.
