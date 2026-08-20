# Gemma Stage A smoke: incomplete infrastructure run

- Date: 2026-08-20
- Frozen source commit: `8682c842ec19567301e77eeac1f31873cfb4b7b9`
- Source tag: `experiment-harness-v0.2`
- Model: `gemma-4-31b-it`, thinking off
- Planned calls: 81
- Successful responses: 71
- Provider quota failures misclassified as malformed: 7
- Dependency-blocked fusers: 3
- Inference status: none; the complete run is discarded before rating or scoring

Pi's JSON mode echoed the user prompt, including its example response object. When Google returned a quota error inside an assistant event while Pi exited zero, the generic stream parser could find the prompt's example JSON and label the call malformed instead of retrying an infrastructure failure.

The correction parses only assistant-role text, promotes assistant error events to non-zero retryable status, adds backoff, lowers concurrency, and freezes a new seeded run. The partial run remains sealed locally and will not be patched, combined with another run, or counted as an experimental sample.
