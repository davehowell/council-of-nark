# Experiment runsheet

Use this runsheet when one person must operate the study. The `just` recipes do not require that person to inspect prompts or outputs.

## 1. Prepare once

1. Install `git`, Python 3.11 or later, `just`, and the selected model CLI.
2. Authenticate the CLI outside this repository.
3. Check out the frozen commit on a clean working tree.
4. Run `just audit`.
5. Run `just experiment-test`.
6. Run `just experiment-doctor experiment/config/stage-a-smoke.json`.

The doctor checks executables, model identifiers, prompt assembly, wrapper length, call counts, and answer-key exclusion. It makes no model calls.

## 2. Verify the adapter, then run Stage A smoke

Make one frozen call before the larger run:

```bash
just experiment-adapter-check
```

Verify that its run reports one successful call. Then run:

```bash
just experiment-stage-a-smoke 3
```

The command prints a run path when it finishes. Save that path as `RUN`. The recipe:

- freezes committed assets and SHA-256 digests;
- builds a seeded, randomised 81-call plan;
- gives each call a fresh detached worktree and process;
- disables tools, project context, skills, extensions, and session persistence where the CLI exposes those controls;
- captures requests, responses, usage, errors, retries, and latency;
- seals raw files by digest;
- creates an arm-blinded rating bundle.

Do not interpret the smoke run as evidence. Use it to find malformed output, obvious ceiling/floor effects, rate limits, and ambiguous keys.

## 3. Verify and inspect run health

```bash
just experiment-verify "$RUN"
just experiment-summary "$RUN"
```

If calls failed for infrastructure reasons, keep the raw attempt records. Create a new run rather than editing a sealed run.

## 4. Rate blind

Two raters independently follow `$RUN/blinded/RUNSHEET.md`. Each rater copies `rating-template.csv`, fills every row, and avoids `plan.json`, `calls/`, and `private/`.

For smoke-test triage only, an arm-blinded LLM rater is available:

```bash
just experiment-judge "$RUN" 2
```

This output is exploratory. It does not replace two blinded human raters for a claim-bearing run.

## 5. Adjudicate and score

Resolve disagreements while arm and provider identities remain hidden. Save one row per finding in `ratings-adjudicated.csv`, then run:

```bash
just experiment-score "$RUN" blinded/ratings-adjudicated.csv adjudicated
```

For LLM smoke triage:

```bash
just experiment-score "$RUN" blinded/ratings-llm.csv llm-smoke
```

## 6. Decide whether to pilot

Before the pilot:

1. repair only failures declared as smoke-test calibration targets;
2. add frozen task variants and clean controls;
3. preregister the smallest effect of interest and exclusion rules;
4. commit and merge those changes;
5. freeze a new commit;
6. run the pilot config with `just experiment-complete experiment/config/stage-a-pilot.json 3`.

Never alter or overwrite a sealed run. Never select thresholds after unblinding confirmatory output.
