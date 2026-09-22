# Isolation threat model

## Objective

Hold local context approximately constant between conditions and prevent a respondent from discovering answer keys, sibling outputs, Git history, or an upstream solution. This is an experimental boundary, not a claim that isolated agents are always the best production workflow.

## Trusted computing base

- macOS kernel and Seatbelt implementation;
- Go toolchain used to build the committed harness;
- Git used by the controller before the model process starts;
- selected provider CLI, runtime, credential mechanism, and remote provider; OAuth clients that require shared keychain/home state are rejected;
- committed prompt, schema, config, rating, and scoring assets.

The freeze records provider CLI versions and SHA-256 entrypoint digests. It cannot hash or inspect remote provider code.

## Child-process filesystem boundary

The Go controller creates and verifies a detached worktree, assembles the exact prompt, and then starts the provider child elsewhere. The child receives:

- an empty current directory;
- a fresh ephemeral `HOME`, cache, config, and temporary directory;
- only adapter/runtime reads needed to start;
- only ephemeral scratch writes;
- no filesystem view of the worktree, repository, run directory, answer keys, or real home.

The generated Seatbelt policy denies by default. Executable probes must show an allowed scratch write, a denied repository read, and denial of an executable outside the explicit allowlist. Failure aborts doctor/freeze; there is no permissive fallback.

A 15 September 2026 launcher probe found that the earlier profiles used `(allow process*)`, which also permitted `process-exec` and made the later executable-path rule ineffective. The profiles now allow `process-fork` separately and list every executable path. Earlier model-call adapters exposed no model shell or discovered extension, so this finding is not evidence that a respondent read local files; it is evidence that the claimed executable restriction was not enforced. No ecological respondent had run. Future snapshots and calls must use the repaired profile and executable-denial probe.

## Network boundary

Outbound network is allowed to the provider CLI because inference requires it. Model tools are disabled, so the remote model has no local curl/browser/shell capability in current runs.

This does not prevent:

- a compromised provider CLI from misusing its own transport access;
- provider-side web search or retrieval hidden behind the API;
- remote caching, personalisation, moderation state, or benchmark memory;
- knowledge learned during model training.

Provider-side search must be explicitly disabled and recorded. If internet-assisted review is studied, make it a separate arm. A controlled offline documentation corpus or allowlisting proxy is preferred when the upstream patch must remain undiscoverable.

## Real open-source task preparation

For a pre-fix ecological task:

1. select the project/issue using frozen inclusion criteria;
2. export the target parent commit with `git archive` into a new directory;
3. remove upstream metadata, normalize the source tree, and leave `.git` absent;
4. exclude changelogs, issue caches, patch files, generated references, agent instructions, and evaluation assets that could reveal or steer a solution;
5. store provenance, evidence tests, dependency closure, and any edit baseline outside the source root;
6. give every condition the same local read/search/test tools through an allowlisted mediator;
7. keep unrestricted shell, Git remotes, and internet disabled unless they are declared factors;
8. retain the upstream patch, tests, and review discussion only for blinded evaluation;
9. accept supported novel findings rather than treating the eventual patch as the sole valid answer.

Deleting `.git` after a normal clone is insufficient if refs, worktrees, caches, or adjacent directories remain visible. Export into a new root. The Gortex pilot exporter also runs its controller-only focused test against a digested Go/module closure with `GOPROXY=off` and network-denied Seatbelt. This validates the snapshot; it does not itself grant model tools.

## Ecological mediator boundary

The ecological design does not mount source into Pi and does not override Pi's built-in read or shell tools. The intended process split is:

```text
provider transport + model
        ↕ Pi tool protocol
isolated Pi child (custom tools only; no source/controller paths)
        ↕ inherited LF-delimited JSON pipes
trusted Go mediator (policy, budgets, transcript)
        ├─ read-only sanitized source
        └─ separate network-denied focused-test sandbox
```

`experiment/ecological/mediator` validates every relative path beneath one canonical history-free source root. It rejects absolute/traversing paths, `.git`, symlinks, non-regular reads, binary/non-UTF-8 content, oversized files, unknown arguments, unlisted tests, and exhausted call/result budgets. Listing, reading, and RE2 searching happen in-process without a shell. Requests are serialized and paired with responses in a controller-owned JSONL transcript; transcript write failure closes the session.

`experiment/ecological/pi/ecological-tools.ts` knows only inherited request/response file descriptors. It receives no source, evidence, cache, or controller path. The launcher must use `--no-builtin-tools`, disable discovered extensions/skills/prompts/themes/context/session persistence, explicitly allow only the tracked tools, and capture Pi's JSON stream. The final submission is a terminating structured-output tool, not a writable source operation.

The no-model ecological doctor now verifies the snapshot, policy, extension and runtime digests; starts isolated Pi with only the custom tools; maps the two mediator pipes; proves source/controller/evidence reads are denied; exercises a mediated health request; explicitly verifies model/thinking state; and seals Pi events, profiles, probes, and the mediator transcript. The separate mediator check runs the allowlisted hidden test and its test-network probe. Both checks preserve attempts and make no provider call.

The single-stage claim runner now implements those requirements: exact prompt/input digests, JSON-mode Pi, pre-transport provider-payload capture, event and mediator transcripts, resource accounting, final-submission validation, ephemeral-state removal, complete sealing, and negative direct-read/Git/traversal/arbitrary-test/shell/test-network probes. Its tracked configuration is a deterministic mock and makes no provider call.

The boundary is still not authorized for a respondent call until the repaired profiles and runner pass from a clean committed controller and the compared arms, resource limits, retry rule, and human rating plan are preregistered. No tracked Pi-backed run configuration exists yet.

## OS accounts

The maintained harness refuses root. A standard non-admin dedicated macOS account is recommended for claim-bearing runs and can be asserted with `COUNCIL_EXPERIMENT_USER`.

The harness does not create one account per call. Doing so requires privileged Directory Services mutations, credential distribution, ownership changes, and cleanup, adding more state than it removes. Per-call Seatbelt profiles and ephemeral homes are the maintained call-level boundary.

agy is deliberately unsupported: when given an ephemeral home, its OAuth flow attempts to locate or create a login keychain. Direct Claude CLI login state is likewise unavailable without the real home/keychain context. Rather than relax isolation or permit UI, the harness fails before starting either client. Gemini and Anthropic experiments use pinned provider models through Pi's sterile credential copy.

## Unsupported platforms

Linux, Windows, containers, VMs, and alternative macOS sandbox schemes are out of scope. A port is a protocol adaptation, not an equivalent invocation; publish its threat model and isolation probe results.
