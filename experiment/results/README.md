# Published results

Selected result snapshots belong here only after raw material and derived reports pass the public-release audit.

A published snapshot should include:

- frozen source commit and config digest;
- run seal digest manifest;
- status, model, usage, latency, and malformed-response summary;
- blinded rating protocol and rater type;
- aggregate set-level scores and uncertainty;
- nulls, negative results, exclusions, and deviations;
- no credentials, private paths, session IDs, raw environment details, or organisation-specific material.

Raw local runs remain under the ignored `experiment/runs/` directory.

## Published snapshots

- [`2026-08-20-stage-a-smoke/`](2026-08-20-stage-a-smoke/): successful 81-call Stage A calibration with exploratory arm-blinded LLM triage.
- [`2026-08-20-stage-a-smoke-instrumentation-failure.md`](2026-08-20-stage-a-smoke-instrumentation-failure.md): the preceding local schema-validation failure, which made no model requests and produced no observations.
- [`2026-08-20-gemma-smoke-incomplete.md`](2026-08-20-gemma-smoke-incomplete.md): sealed and discarded partial Gemma run that exposed Pi stream/error classification defects before scoring.
- [`2026-08-20-stage-a-smoke-gemma/`](2026-08-20-stage-a-smoke-gemma/): clean 81-call low-reasoning Gemma smoke; specialists increased raw recall and conditions found different defects.
- [`2026-08-20-persona-pair-gemma/`](2026-08-20-persona-pair-gemma/): 30 repeated functional/K-2SO correctness pairs; the fictional wrapper reduced mean F1 in this narrow test.
