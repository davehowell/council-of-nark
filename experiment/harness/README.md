# Go harness design

The maintained harness is implemented in Go's standard library and supports macOS only. Provider CLIs remain external so authentication stays in existing credential stores. Historical Python implementations remain reproducible from their tagged source commits; they are not active on `main`.

## Lifecycle

The single `council-exp` command provides these stages:

1. `doctor` validates macOS, Seatbelt, config, prompt assembly, wrapper length, model identifiers, call counts, and answer-key exclusion without model calls.
2. `freeze` requires a clean tree, reruns the public audit, records `HEAD`, hashes every tracked input, and records external CLI entrypoint digests/versions.
3. `plan` builds opaque call/output-set IDs and a deterministic hash-ordered schedule.
4. `run` schedules dependencies and starts every call in a fresh process from a detached worktree at the frozen commit.
5. `summarize` reports status, usage, cost, findings, and latency without scoring.
6. `seal` hashes and makes raw run files read-only; `verify` recomputes those digests.
7. `bundle` shuffles findings into a rating bundle that hides arm, wrapper, model, and provider.
8. `judge` can produce exploratory arm-blinded smoke triage. It requires a clean committed harness and records the derived-stage commit.
9. `score` consumes one adjudicated rating per finding and reports set/group metrics.

Use `just` rather than calling the binary directly. `go run` builds outside the worktree and no persistent harness binary is required.

## Mandatory Seatbelt boundary

Every non-mock provider child runs through `/usr/bin/sandbox-exec` with a generated deny-by-default profile:

- prompt assembly occurs in a fresh detached worktree;
- the child process starts in a different empty directory;
- the child cannot read the repository, worktree, answer keys, run records, or real home directory;
- `HOME`, cache, config, temporary, and current-working directories are fresh per attempt;
- only a small adapter-specific credential allowlist is copied into the ephemeral home; conversation/history stores are never copied;
- writes are limited to ephemeral scratch and `/dev/null`;
- executable/runtime paths are read-only;
- provider transport may use outbound networking;
- model tools, sessions, context files, extensions, skills, and project trust are disabled through adapter flags;
- ephemeral state is deleted after each attempt.

Run the executable probe explicitly:

```bash
just experiment-sandbox-check
```

It must prove that a scratch write succeeds and a repository read fails. `doctor` and `freeze` rerun the same probe. There is no unsandboxed fallback and the harness refuses non-macOS hosts or root.

Seatbelt is deprecated by Apple even though it remains present on supported macOS versions. If Apple removes it, the harness fails closed until a replacement threat model is implemented.

## Network and tools are separate capabilities

The provider CLI needs network access to submit the prompt. That does not grant the model a local shell or web tool: all current respondent and rater commands disable tools.

Seatbelt cannot observe or disable a provider-side search tool. A future internet-enabled study must therefore be a separate declared arm with provider-side search disabled/enabled explicitly. For real-project tasks, prefer frozen offline documentation/search corpora or a controlled proxy rather than unrestricted web search that can locate the matching upstream patch.

## Dedicated macOS account

For stronger UID/ACL separation, create one standard non-admin macOS account for experiment execution and run the entire harness there. Set:

```bash
export COUNCIL_EXPERIMENT_USER=experiment-account-name
```

The harness then refuses another account. It does not create users, alter Directory Services, copy credentials between accounts, or use `sudo`. A new account per call adds substantial privileged lifecycle state while the per-attempt empty HOME and Seatbelt profile already provide call-level separation; it is not part of the maintained protocol.

## Adapters

- **Claude:** safe mode, no tools, no session persistence, replacement system prompt, strict JSON Schema.
- **agy/Gemini:** fresh print process, pinned model/effort, CLI sandbox, slash commands disabled, strict JSON Schema.
- **Pi models:** no session, tools, extensions, skills, templates, themes, context files, or project trust; replacement system prompt. Pi is launched through the pinned Node executable and its assistant-role events are parsed without inspecting echoed user content.
- **mock:** deterministic empty response for plumbing tests; no child process.

External CLI entrypoints and versions are frozen because provider clients are part of the trusted computing base. Provider-side personalisation, hidden system prompts, caching, and server-side tooling remain unobservable validity limitations.

## Retry rule

Only non-zero process/transport failures are retried, with configured linear backoff. Pi assistant error events are promoted to retryable failures even if Pi exits zero. A zero-exit malformed substantive response is sealed and scores zero. Retries are infrastructure attempts, not additional samples.

## Data boundaries

`request.json` records the exact prompts, schema, model, source commit, isolation declaration, and digests. Prompt assembly rejects answer-key headings, planted IDs, and long key claims. Answer keys enter only the blinded rating stage. Raw and derived stages have separate integrity boundaries.
