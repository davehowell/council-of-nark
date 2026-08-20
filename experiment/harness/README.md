# Harness design

The harness uses Python's standard library. Model CLIs remain external so authentication stays in each user's existing credential store.

## Lifecycle

1. `doctor.py` validates a config, prompt assembly, wrapper length, model identifiers, call counts, and answer-key exclusion without model calls.
2. `freeze.py` requires a clean tree, records `HEAD`, and hashes every tracked experiment input.
3. `plan.py` builds opaque call and output-set IDs from the config seed.
4. `run.py` schedules dependencies and starts every call in a fresh process and detached worktree at the frozen commit.
5. `summarize.py` reports status, usage, cost, findings, and latency without scoring.
6. `seal.py` hashes and makes raw run files read-only. `verify.py` recomputes those digests.
7. `bundle.py` shuffles findings into a bundle that hides arm, wrapper, model, and provider.
8. `judge.py` can produce exploratory arm-blinded smoke triage.
9. `score.py` consumes one adjudicated rating per finding and reports set/group metrics.

## Adapters

- **Claude:** `--safe-mode`, no tools, no session persistence, replacement system prompt, strict JSON Schema.
- **agy/Gemini:** fresh print process, pinned model and effort, sandbox, slash commands disabled, strict JSON Schema.
- **Pi/OpenAI:** no session, tools, extensions, skills, templates, themes, context files, or project trust; replacement system prompt.
- **mock:** deterministic empty response for plumbing tests.

The runner passes prompts as process arguments and closes stdin. It sends no file tools to the model. Provider CLIs can still have unobservable server-side behaviour, and the agy CLI does not expose every isolation control that Claude and Pi expose.

## Retry rule

The runner retries only a non-zero process or transport failure, up to the config's `max_attempts`. A zero-exit malformed substantive response is not repaired or rerolled; it is sealed and scores zero. Retries are infrastructure attempts, not additional samples.

## Data boundaries

`request.json` records the exact system prompt, user prompt, schema, model, source commit, and digests. It explicitly records `answer_key_in_context: false`. Prompt assembly also rejects answer-key headings, planted IDs, and long answer-key claims.

Answer keys enter only the blinded rating stage. The raw run and the derived rating bundle have separate integrity boundaries. Derived files can be regenerated from the sealed raw run.
