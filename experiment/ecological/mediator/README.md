# Ecological tool mediator

This package is the trusted, controller-side boundary for ecological source access. It is not an agent and does not call a model.

## Protocol

The Pi extension writes one LF-delimited JSON request per tool call to inherited file descriptor 3 and reads correlated responses from descriptor 4. Frames use UTF-8 JSON:

```json
{"request_id":"provider-tool-id","tool":"source_read","arguments":{"path":"internal/x.go","start_line":1,"line_count":80}}
```

```json
{"request_id":"provider-tool-id","sequence":2,"ok":true,"output":"1:package x","truncated":false,"metadata":{"result_bytes":11}}
```

The launcher, not the extension, owns the other pipe ends. The extension receives no source or controller path. The mediator serializes requests even if Pi presents parallel tool calls, applies overall and per-tool budgets, and records each valid request/response in controller-owned JSONL. A transcript write failure closes the session.

## Operations

- `source_list`: deterministic bounded directory listing, depth at most 3.
- `source_read`: numbered lines from one regular UTF-8 source file.
- `source_search`: bounded Go RE2 line search with an optional basename glob.
- `run_focused_test`: exact target allowlist; implementation remains controller-side.

There is no write, edit, shell, Git, URL, arbitrary command, arbitrary test, or internet operation. Absolute/traversing paths, `.git`, symlinks, binary/invalid UTF-8 content, oversized files, and unknown JSON fields are rejected.

The final `submit_ecological_review` tool lives in [`../pi/ecological-tools.ts`](../pi/ecological-tools.ts). It terminates Pi with structured output and does not enter the mediator protocol.

## Gortex infrastructure check

From a committed clean tree, after creating the clean snapshot:

```bash
just ecological-gortex-mediator-check \
  experiment/ecological/work/<clean-snapshot-attempt>
```

The check verifies the snapshot seal and exact source tree, hashes policy/extension/mediator inputs, runs positive list/read/search/test probes, and requires traversal and arbitrary-test denial. The hidden parent regression runs under a separate deny-by-default Seatbelt profile using the frozen Go/module closure. Each test invocation proves that network access and an unlisted shell executable are denied. Raw test logs remain under ignored `experiment/ecological/mediator-runs/`; only sanitized output crosses the mediator response.

This command makes no model call and does not exercise Pi. Follow it with:

```bash
just ecological-gortex-pi-doctor \
  experiment/ecological/work/<clean-snapshot-attempt>
```

The doctor starts isolated Pi in persistent RPC mode without sending an agent prompt, exercises the inherited pipes through a controller-only extension command, verifies model/thinking state, and seals events and isolation metadata. It uses controlled termination after checks because RPC mode is a long-lived integration server. Retain failed attempts; do not repair them in place.
