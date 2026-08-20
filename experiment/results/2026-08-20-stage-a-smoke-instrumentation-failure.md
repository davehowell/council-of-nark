# Stage A smoke: instrumentation failure

- Date: 2026-08-20
- Frozen source commit: `0e12ab56d00b2d0af3a2cf5d041d02440dc457ba`
- Config: `experiment/config/stage-a-smoke.json`
- Planned calls: 81
- Successful model responses: 0
- Local CLI errors: 72 reviewer calls
- Dependency-blocked calls: 9 fusers
- Inference status: none; this run contains no experimental observations

## Failure

Claude Code rejected the structured-output schema before sending a model request:

```text
--json-schema is not a valid JSON Schema: no schema with key or ref for the draft-2020 meta-schema URI
```

The source schema included a standards-documentation `$schema` annotation. The installed Claude Code validator accepts the schema subset but does not resolve that annotation. Every reviewer process exited locally. Each fuser was then correctly blocked by failed dependencies.

## Response

The adapter now keeps `$schema` in the committed source and frozen request record but removes that annotation from the CLI argument. A unit test fixes this behaviour. A one-call frozen adapter check must pass before the next 81-call smoke run.

The failed run remains sealed locally. It will not be reused, repaired, or counted as a sample.
