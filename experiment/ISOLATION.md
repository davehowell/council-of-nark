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

The generated Seatbelt policy denies by default. Executable probes must show an allowed scratch write and denied repository read. Failure aborts doctor/freeze; there is no permissive fallback.

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
3. remove upstream metadata and initialise one neutral root commit;
4. exclude changelogs, issue caches, patch files, and generated references that reveal the solution;
5. give every condition the same local read/search/test tools through an allowlisted mediator;
6. keep unrestricted shell, Git remotes, and internet disabled unless they are declared factors;
7. retain the upstream patch, tests, and review discussion only for blinded evaluation;
8. accept supported novel findings rather than treating the eventual patch as the sole valid answer.

Deleting `.git` after a normal clone is insufficient if refs, worktrees, caches, or adjacent directories remain visible. Export into a new sandbox root.

## OS accounts

The maintained harness refuses root. A standard non-admin dedicated macOS account is recommended for claim-bearing runs and can be asserted with `COUNCIL_EXPERIMENT_USER`.

The harness does not create one account per call. Doing so requires privileged Directory Services mutations, credential distribution, ownership changes, and cleanup, adding more state than it removes. Per-call Seatbelt profiles and ephemeral homes are the maintained call-level boundary.

agy is deliberately unsupported: when given an ephemeral home, its OAuth flow attempts to locate or create a login keychain. Direct Claude CLI login state is likewise unavailable without the real home/keychain context. Rather than relax isolation or permit UI, the harness fails before starting either client. Gemini and Anthropic experiments use pinned provider models through Pi's sterile credential copy.

## Unsupported platforms

Linux, Windows, containers, VMs, and alternative macOS sandbox schemes are out of scope. A port is a protocol adaptation, not an equivalent invocation; publish its threat model and isolation probe results.
