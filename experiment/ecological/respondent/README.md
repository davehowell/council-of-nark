# Ecological respondent launcher

Status: infrastructure only. The tracked configuration uses a deterministic mock and makes no model call. No live five-arm configuration is committed.

The launcher runs one isolated respondent stage. A later arm controller must assemble the continuing-session and fusion procedures in the finite closing round; this command does not silently approximate those arms.

## Mock lifecycle

After a clean Gortex snapshot exists, run:

```bash
just ecological-gortex-respondent-mock experiment/ecological/work/<snapshot-attempt>
```

The deterministic mock exercises:

- snapshot seal and source-tree verification;
- bounded list, read, search, and focused-test requests;
- Git, traversal, and arbitrary-test denials;
- exact prompt, config, policy, extension, model, and budget capture;
- one terminating structured submission;
- event, usage, provider-request, and mediator accounting;
- removal of copied credentials and reproducible build caches;
- an immutable whole-attempt SHA-256 seal.

It emits a fixed fake finding for validation only. Its `adapter.provider_calls` value is zero, and it must never enter an outcome analysis.

Verify a completed attempt independently:

```bash
just ecological-respondent-verify experiment/ecological/respondent-runs/<attempt>
```

## Live boundary

The Pi path uses JSON mode from an empty working directory and ephemeral home. Built-in tools, sessions, skills, prompt templates, themes, context files, project trust, and extension discovery are disabled. The provider child receives only the tracked extension plus inherited mediator and audit pipes. Source, controller, evidence, and council files remain unreadable.

The extension records each serialized provider payload before transport, including the effective system content, conversation, and tool schemas. It records response metadata and assistant usage without changing the payload. Aggregate provider-request and token caps are checked before each later request. The response that crosses a token cap cannot be recalled, so any last-response overrun is recorded explicitly. The mediator separately enforces tool counts and returned-source bytes. A wall timeout and final-submission byte cap also apply.

The launcher accepts a Pi-backed config only with `--allow-live`, from a clean committed tree. This mechanical switch is not research authorization. Before using it, freeze the five arm prompts and schedule, rerun every boundary check from the committed controller, name two human raters, set a deadline and spend ceiling, and complete the excluded feasibility preregistration required by [`../../FINAL_ROUND.md`](../../FINAL_ROUND.md).

## Attempt contents

A completed attempt contains:

- `request.json`: exact tracked inputs, text, digests, provider state, and budgets;
- `provider-audit.jsonl`: serialized provider payloads and response/usage lifecycle;
- `pi-events.jsonl`: Pi's raw JSON event stream;
- `mediator-transcript.jsonl`: every source/test request and bounded response;
- `focused-tests/`: network, executable, profile, and exact test logs;
- `final-submission.json`: present only for one schema-valid final action;
- `result.json`: health, resource totals, and valid/malformed outcome;
- `status.json`: attempt completion state;
- `seal.json`: digest, path, and byte count for every other retained file.

A malformed substantive response remains sealed and receives the protocol's declared zero utility. Infrastructure errors remain failed attempts. Never edit or replace either kind in place.
