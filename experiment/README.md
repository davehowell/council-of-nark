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

Each reviewer receives only `review-packet.md`. The harness must reject any request containing answer-key content. The keys are public after publication, so future users must treat this as a reproducible demonstration rather than a permanently hidden benchmark.

## Stages

- **Smoke:** one sample per packet and arm. Repair broken schemas and ceiling/floor effects. Do not make inferential claims.
- **Pilot:** ten samples per packet and Stage A arm; three samples per packet, topology, and chain order. Estimate effects and variance.
- **Confirmatory extension:** freeze semantics-preserving variants and clean controls before choosing thresholds or inspecting outputs. Use blinded human ratings for ambiguous matches.

The three base packets are too few for broad claims. Repeated sampling estimates model variance, not task diversity.

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
