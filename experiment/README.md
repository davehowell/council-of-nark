# Experiment

This directory contains the protocol and frozen source material for a staged empirical study of the Council of Nark.

## Research questions

1. Do independent functional specialists find more planted defects than the same number of repeated omnibus reviews?
2. Does a short fictional wrapper add useful findings beyond a byte-identical functional lens kernel?
3. Does one fusion call improve precision without discarding enough true findings to reduce F1?
4. Does independent fan-out outperform an informed serial chain across all reviewer orders?
5. Do any effects replicate across model providers?

The study measures model output. It cannot infer which internal features, activations, or weights caused an output difference.

## Material

- [`protocol.md`](protocol.md): hypotheses, arms, topology, scoring, sampling, threats, and interpretation rules.
- [`prompts/`](prompts/): controlled single-reviewer, specialist, and fuser prompts.
- [`scenarios/`](scenarios/): three synthetic packets and their planted-defect keys.
- [`config/`](config/): pinned smoke, pilot, topology, provider-pair, and judge configurations.
- [`schema/`](schema/): strict response contracts.
- [`harness/`](harness/): standard-library Python runner, adapters, blinding, sealing, rating, and scoring code.
- [`HARDNESS.md`](HARDNESS.md): scenario difficulty dimensions, smoke calibration triggers, and confirmatory mutation rules.
- [`PREREGISTRATION.md`](PREREGISTRATION.md): commit-before-run hypothesis, threshold, sample, rating, and exclusion template.
- [`RUNSHEET.md`](RUNSHEET.md): the operator and human-rating procedure.

Each reviewer receives only `review-packet.md`. The harness rejects a request containing answer-key markers or claims. The keys are public in this repository, so future users must treat this as a reproducible demonstration rather than a permanently hidden benchmark.

## Stages

- **Smoke:** one sample per packet and arm. Repair broken schemas and ceiling/floor effects. Do not make inferential claims.
- **Pilot:** ten samples per packet and Stage A arm; three samples per packet, topology, and chain order. Estimate effects and variance.
- **Confirmatory extension:** freeze semantics-preserving variants and clean controls before choosing thresholds or inspecting outputs. Use blinded human ratings for ambiguous matches.

The three base packets are too few for broad claims. Repeated sampling estimates model variance, not task diversity.

The first successful Stage A calibration is published under [`results/2026-08-20-stage-a-smoke/`](results/2026-08-20-stage-a-smoke/). Its descriptive results do not support advancing directly to the pilot: matched repeated omnibus sampling beat the functional specialists after fusion, fictional-wrapper effects were inconsistent, and the weak baseline had limited headroom. Harder frozen variants, clean controls, and human ratings come next.

## Freeze rules

Before a claim-bearing run:

1. commit every prompt, packet, key, schema, model identifier, decoding setting, exclusion rule, and analysis script;
2. run the public-release audit;
3. generate a manifest with the Git commit and SHA-256 digest of every input;
4. randomise run order in blocks;
5. start each call in a fresh process and clean detached worktree with sessions, tools, project context, and optional memory disabled;
6. capture raw request, raw response, model metadata, usage, latency, errors, and retries;
7. make the completed run immutable by digest and verify it before analysis.

A clean process prevents local context leakage. It cannot prove that a provider has no unobserved server-side state. Record that limitation.

## Scoring

A finding maps to zero or one planted defect. Duplicate phrasings of the same defect count once. Unmatched material claims are false positives.

The primary metric is macro F1 over planted defects, weighted equally by packet. Also report precision, recall, p10, worst case, variance, token/cost efficiency, raw-union coverage, fusion retention, malformed responses, and latency.

Use deterministic location checks where possible. Two blinded humans decide ambiguous semantic mappings; an LLM judge may triage but must not be the final confirmatory rater.

## Run the harness

Preflight does not call a model:

```bash
just experiment-test
just experiment-doctor experiment/config/stage-a-smoke.json
```

After the harness and assets are committed on a clean tree, run the complete Stage A smoke:

```bash
just experiment-stage-a-smoke 3
```

The Stage A smoke makes 81 calls: 3 packets × (S0/S1/S2 at one call each + M0/M1/M2 at seven reviewers and one fuser each). The topology smoke makes 144 calls across all six role orders. The provider-pair smoke makes 18 calls.

Use the generic recipe for another frozen config:

```bash
just experiment-complete experiment/config/provider-pair-smoke.json 3
```

Every model call runs in a fresh process from a fresh detached worktree at the recorded commit. The adapters disable tools, sessions, extensions, skills, and project context where the CLI exposes those switches. The runner keeps answer keys out of review and fusion prompts. It stores raw attempts under ignored `experiment/runs/`, then seals the raw file set by SHA-256 digest.

See [`RUNSHEET.md`](RUNSHEET.md) for blinded ratings and scoring. A smoke-only arm-blinded LLM triage recipe is available, but confirmatory scoring requires two independent human raters.

## Model-state limitation

Fresh local sessions prevent this harness, earlier council calls, local memory files, and the current worktree from entering model context. They do not prove that a provider has no server-side personalisation, caching, abuse classifiers, or other hidden state. Stage A uses a pinned Anthropic snapshot so an OpenAI controller does not also serve as the primary respondent. Stage B treats provider as an experimental factor and reports hidden-system differences as a limitation.
