# Experiment contamination review

This review was performed after the first plumbing smoke and before the cheaper-model rerun.

## Context isolation

The implemented boundaries are sound for local context contamination:

- the Go controller assembles each prompt in a fresh detached worktree at the frozen commit; the provider child starts elsewhere under Seatbelt and cannot read the worktree;
- reviewer prompts contain the packet but not the answer-key file, planted IDs, arm name, hypothesis, or other outputs;
- prompt assembly rejects answer-key headings, IDs, and long verbatim key claims;
- Claude runs in safe mode with tools and session persistence disabled;
- Pi runs without tools, sessions, extensions, skills, templates, themes, project context, or project trust;
- model stdin is closed;
- raw requests preserve the exact system/user prompts and digests;
- fusers see opaque reviewer IDs;
- rating bundles hide arm, wrapper, model, and provider metadata.

These controls cannot observe provider-side personalisation, hidden system prompts, caching, safety classifiers, or future benchmark contamination. Public answer keys make this a reproducible demonstration, not a durable secret benchmark.

## Issues found

### 1. Smoke model/task mismatch

The pinned Haiku snapshot solved much of the toy material in one call. S0 reached perfect recall on one packet. This is a ceiling/headroom problem, not evidence that prompt roles have no effect.

**Correction:** run the same frozen Stage A design with explicit `pi --model gemma-4-31b-it --thinking off` before changing scenarios. Do not use Pi's configured default model.

### 2. Seeded order was not honoured by the scheduler

`plan.json` contained a seeded shuffle, but `run.py` converted pending call IDs to a Python set. Set iteration discarded the planned ordering. Arms remained hidden and calls remained independent, but time/order blocking was no longer reproducible.

**Correction:** preserve pending calls as the ordered plan list. Dependency readiness can skip an item without reordering the remaining ready calls.

### 3. Raw false positives were not semantic unions

True duplicates collapsed to one planted ID. Unmatched findings were each counted as a false positive, even when several reviewers repeated the same unsupported claim. This biased raw panels downward and overstated fusion's precision gain.

**Correction:** blinded raters assign a false-positive cluster within each output set. Scoring counts unique clusters, just as it counts unique true defect IDs. Existing published smoke raw-union scores retain the old rule and are labelled accordingly.

### 4. Fused F1 can flatten reviewer differences

The common fuser sees the packet so it can validate support. A capable fuser can independently recognise mechanisms and make M0/M1/M2 final verdicts converge, despite the instruction not to originate findings.

This is a design trade-off rather than answer leakage. All multi-reviewer arms receive the same extra call, so M1−M0 remains call-matched. It does mean fused F1 alone is a poor test of role diversity.

**Correction:** make semantically de-duplicated raw-union coverage/precision primary for specialist diversity. Keep fused F1 as the practical arbiter outcome and report fusion retention.

### 5. One sample cannot separate treatment from decoding noise

The CLI does not expose a common temperature control. One output per packet/arm estimates neither within-condition variance nor the size of a prompt effect.

**Correction:** the next run remains a smoke. If Gemma restores headroom, run randomised repeated blocks before interpreting M2−M1. A difference smaller than within-arm variation is not a persona effect.

### 6. LLM ratings are calibration only

The first smoke used an arm-blinded LLM rater, mostly from the same provider family as the respondents. Six fallback labels used another provider. Character style can also leak through wording even when catchphrases are schema-excluded.

**Correction:** use LLM ratings only to calibrate plumbing and hardness. Two independent humans, blinded to condition, decide claim-bearing semantic matches and fix quality.

### 7. Pi JSON streams echo the user prompt

Pi JSON mode emits user and assistant events. The first Pi parser searched the complete stream, so after a provider error it could find the JSON example inside the echoed user prompt and misclassify infrastructure failure as malformed model output. Pi also reports provider quota errors while the local process exits zero.

**Correction:** parse only assistant-role text events for Pi. Promote assistant `stopReason: error` events to retryable infrastructure failures, preserve each attempt, add backoff, and test prompt/response separation explicitly. The first Gemma run was sealed and discarded before scoring.

### 8. Structured-output wrappers embed their own schema

agy returns both `structured_output` and `json_schema`. A generic recursive extractor accepted any object containing a `judgements` key, so it could select the schema's `properties.judgements` object instead of the model's judgement array.

**Correction:** a response root is eligible only when the requested root value is an array. Parser-only failures can be recovered mechanically from captured structured output after schema validation, without another model call or changed judgement.

### 9. A clean worktree was still visible to the provider process

The earlier controller disabled model tools and context loading, but its provider child used the detached worktree as its current directory. Correct adapter flags made accidental loading unlikely, not impossible. A changed or compromised CLI could still inspect public answer keys, Git metadata, neighbouring project files, or user configuration.

**Correction:** the active harness is now Go/macOS-only and fails closed unless a Seatbelt probe passes. Prompt assembly still uses the verified worktree, but the child receives an empty cwd and ephemeral home under a deny-by-default profile. Repository/worktree reads fail, writes are scratch-only, and external CLI entrypoint digests are frozen. Provider transport networking remains permitted; server-side tools/state remain unobservable.

## Checks that did not reveal contamination

- M1 and M2 use byte-identical lens kernels; only paired wrappers differ.
- S1 and S2 share the byte-identical omnibus kernel and output contract.
- M0, M1, and M2 use seven reviewer calls plus the same fuser contract.
- The answer key does not appear in sampled raw requests.
- No model receives Git history or another arm's response.
- The current controller's conversation is not passed to subprocess respondents.

## Interpretation of the first smoke

The run proved that 81 isolated structured calls, dependency scheduling, sealing, blinding, and scoring worked after adapter repair. Its scores are descriptive calibration data. They do not establish that the council works or is worthless.
