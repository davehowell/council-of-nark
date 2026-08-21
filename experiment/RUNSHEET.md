# Experiment runsheet

Use this runsheet when one person must operate the study. The `just` recipes do not require that person to inspect prompts or outputs.

## 1. Prepare once

1. On macOS, install `git`, Go 1.22+, `just`, and Pi. The maintained harness does not support other operating systems.
2. Authenticate Pi outside this repository.
3. Check out the frozen commit on a clean working tree.
4. Run `just audit`.
5. Run `just experiment-test` and `just experiment-sandbox-check`.
6. Run `just experiment-doctor experiment/config/stage-a-smoke-gemma.json`.

The doctor checks Seatbelt denial, executables, model identifiers, prompt assembly, wrapper length, call counts, and answer-key exclusion. It makes no model calls.

## 2. Verify the adapter, then run Stage A smoke

Make one frozen call before the larger run:

```bash
just experiment-adapter-check-gemma
```

Verify that its run reports one successful call from `gemma-4-31b-it` with thinking disabled. Then run:

```bash
just experiment-stage-a-smoke-gemma 2
```

The command prints a run path when it finishes. Save that path as `RUN`. The recipe:

- freezes committed assets and SHA-256 digests;
- builds a seeded, randomised 81-call plan;
- assembles each prompt from a fresh detached worktree;
- starts the provider child under Seatbelt in an empty cwd and ephemeral home without repository/worktree access;
- disables tools, project context, skills, extensions, and session persistence;
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

Two raters independently follow `$RUN/blinded/RUNSHEET.md`. Bundle item, set, and pair IDs use HMAC-SHA-256 with a private random key; unlike the earlier public-seed IDs, conditions cannot be derived from the committed config alone.

A project author or council advocate may be one rater if they disclose that prior. They must not be the only rater. Complete the phases in order:

1. copy `rating-template.csv` and map shuffled findings individually before seeing pairs;
2. lock that CSV;
3. copy `pairwise-rating-template.csv` and score randomised A/B outputs for supportedness, actionability, and fix quality;
4. only after each A/B score, guess the condition and record whether wording revealed it.

Output language is part of the treatment and is not rewritten. Therefore condition labels and IDs are blinded, but treatment blinding can fail. Report guess accuracy and self-reported reveal rate rather than claiming perfect blindness.

For smoke-test triage only, an arm-blinded LLM rater is available:

```bash
just experiment-judge "$RUN" 2 experiment/config/judge-smoke-gemini-low.json
```

This output is exploratory. It does not replace two blinded human raters for a claim-bearing run. If the default rater reaches a provider limit, resume only the missing sets with the documented cross-provider fallback:

```bash
just experiment-judge "$RUN" 1 experiment/config/judge-smoke-openai.json
```

The rating CSV records the model for each item. Report mixed-rater use as a smoke-stage limitation.

## 5. Adjudicate and score

Resolve finding-map disagreements while arm and provider identities remain hidden. Preserve both original files. Save one adjudicated row per finding in `ratings-adjudicated.csv`, then run:

```bash
just experiment-score "$RUN" blinded/ratings-adjudicated.csv adjudicated
```

For LLM smoke triage:

```bash
just experiment-score "$RUN" blinded/ratings-llm.csv llm-smoke
```

After both qualitative A/B files are locked, concatenate them with one header and unblind only the derived report:

```bash
just experiment-qualitative "$RUN" blinded/pairwise-ratings-both.csv qualitative
```

Opaque IDs remain in the output for auditability; the derived CSV adds actual left/right conditions. Never replace the original blinded ratings.

## 6. Decide whether to pilot

Before the pilot:

1. repair only failures declared as smoke-test calibration targets;
2. add frozen task variants and clean controls;
3. preregister the smallest effect of interest and exclusion rules;
4. commit and merge those changes;
5. freeze a new commit;
6. run the pilot config with `just experiment-complete experiment/config/stage-a-pilot.json 3`.

Never alter or overwrite a sealed run. Never select thresholds after unblinding confirmatory output.
