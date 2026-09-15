<p align="center">
  <img src="assets/nark-council.png" alt="The Council of Nark: a group of robot reviewers around a shared console" width="100%">
</p>

# Council of Nark

An open, experimental toolkit for adversarial software review with a small panel of independent specialist agents and one fusion step.

The council is designed to widen attention, not manufacture consensus. Each reviewer uses a narrow functional lens, reports evidence in the same schema, and cannot edit the artifact. An arbiter then de-duplicates, ranks, and reconciles the findings. The robot personas make each lens memorable; whether the fictional wrappers improve review quality is a testable question, not an assumption.

## What is here

| Path | Purpose |
|---|---|
| [`agents/`](agents/) | Review-only Claude Code subagent definitions for simplicity, correctness, architecture, security, and technical language. |
| [`skills/nark-matrix/`](skills/nark-matrix/) | The orchestration skill plus standalone Bender and Holly prompt templates. |
| [`skills/walle-ste/`](skills/walle-ste/) | An STE-derived technical-writing review kernel used by WALL-E. |
| [`experiment/`](experiment/) | Synthetic packets, controlled prompts, answer keys, and the empirical protocol. |
| [`presentations/`](presentations/) | Four talks: the hypothesis, protocol, failed experiments and evidence, and the next experimental stage. |
| [`scripts/public_audit.py`](scripts/public_audit.py) | A pre-publication scan for common secrets, private paths, private hosts, and local forbidden terms. |

## The review lenses

- **HK-47:** unnecessary complexity and speculative machinery.
- **K-2SO:** correctness, tests, observability, idempotency, and silent failure.
- **GLaDOS:** architecture, contracts, coupling, migrations, and optional synthesis.
- **C-3PO:** secrets, personal data, permissions, destructive actions, and explicit controls.
- **Bender:** cost and compute waste.
- **Holly:** long-horizon operational entropy and the 03:00 hand-off.
- **WALL-E:** durable technical language and plain-language translation after fusion.

The controller selects the smallest relevant subset. A council has at least two reviewers. All reviewers are review-only.

## Install the skills

For Claude Code, copy the agents and skills into your user configuration:

```bash
mkdir -p ~/.claude/agents ~/.claude/skills
cp agents/*.md ~/.claude/agents/
cp -R skills/nark-matrix ~/.claude/skills/
cp -R skills/walle-ste ~/.claude/skills/
```

Agent Skills-compatible harnesses can load the skill directly from `skills/nark-matrix/SKILL.md`. Pi can load it for a session with:

```bash
pi --skill skills/nark-matrix
```

Pi deliberately has no built-in subagent policy. The skill therefore describes the orchestration pattern, while the caller decides how to start independent sessions.

## Presentations

The [GitHub Pages presentation hub](https://davehowell.github.io/council-of-nark/) hosts all decks, published-result tiles, downloadable PDFs, and the living experiment timeline.

- **Part 1: The Council of Nark** presents the original hypothesis, the reviewer roles, and the proposed fan-out-and-fuse structure.
- **Part 2: Put the Council on Trial** explains fair controls, scoring and the difference between uncertainty and equivalence.
- **Part 3: The Experiments Failed. A Lot.** follows measurement failures, mixed council results and the negative K-2SO calibration.
- **Part 4: What Would Settle It?** is the pivot into an ecological comparison with one agent reviewing its own work. It is not the conclusion.

Part 1 has ten slides and an approximately eight-minute script. Parts 2–4 each have six slides and an approximately five-minute script. Read [all speaker notes](presentations/notes/speaker-notes.md), or run `just slides-talk`, `just slides-experiment`, `just slides-trials`, or `just slides-pivot`. See [`presentations/README.md`](presentations/README.md) for editing and exports. A conclusion talk belongs after the ecological comparison has results.

## Experiments

The [7 September evidence review](experiment/CONCLUSION.md) concludes that council advantage and single-review equivalence are both unestablished. The repeated K-2SO result is negative within its narrow calibration scope, with slightly fewer total tokens. Reproduce the published arithmetic with `just experiment-evidence-audit`. The [proposed closing round](experiment/FINAL_ROUND.md) compares practical alternatives on external-project review tasks; it has not been run.

The supported experiment runner targets **macOS only** and requires Go 1.22+ plus `/usr/bin/sandbox-exec`. Linux, Windows, containers, and other sandbox implementations are intentionally out of scope; adapting the protocol is left to replicators, who must document any isolation differences.

The study is an empirical pilot, not a benchmark of software engineering in general. It uses synthetic, frozen review packets with planted defects and clean facts. The core controls are:

1. match call counts before crediting specialist roles;
2. keep functional kernels identical when testing character wrappers;
3. score raw panel unions and fused verdicts separately;
4. test all orders of a three-reviewer chain;
5. blind scoring to arm and provider;
6. report false positives, variance, tokens, cost, and latency alongside recall;
7. state that black-box output cannot prove an internal model mechanism.

The answer keys are public for reproducibility, but the run harness keeps them out of reviewer prompts. This makes the repository suitable for replication, not a permanent uncontaminated benchmark. See [`experiment/README.md`](experiment/README.md), [`experiment/protocol.md`](experiment/protocol.md), and the operator [`experiment/RUNSHEET.md`](experiment/RUNSHEET.md).

```bash
just experiment-test                                      # no model calls
just experiment-sandbox-check                              # no model calls
just experiment-doctor experiment/config/stage-a-smoke.json  # no model calls
just experiment-stage-a-smoke 3                           # frozen 81-call smoke
```

Each prompt is assembled in a fresh detached worktree. The provider child then starts in an empty directory and ephemeral home under a deny-by-default macOS Seatbelt profile, without repository/worktree access, tools, project context, optional memory, or session persistence. Raw requests and responses are sealed by digest before arm-blinded scoring.

The first 81-call Stage A plumbing smoke and its mixed calibration results are published in [`experiment/results/2026-08-20-stage-a-smoke/`](experiment/results/2026-08-20-stage-a-smoke/). The low-reasoning Gemma rerun and 30-pair persona calibration are also under [`experiment/results/`](experiment/results/). They are explicitly non-confirmatory and use LLM triage; the repository does not present those scores as proof that the council works or is worthless.

The next ecological stage curates eight PR-backed post–June 2026 pre-fix tasks across Manim, dlt, Gortex, and turbovec under [`experiment/ecological/`](experiment/ecological/), plus one separately gated Modular extreme-reserve candidate. The Gortex pilot now has a clean history-free snapshot, bounded controller-side source/test mediator, anchored scoring key, and no-model isolated-Pi doctor. No ecological respondent call has occurred; the JSON-mode claim runner and compared-arm design are still gated. See [`experiment/CHECKPOINT.md`](experiment/CHECKPOINT.md) for the current handoff.

## Public-release audit

Run this before every public commit:

```bash
just audit
```

For organisation-specific terms, put one case-insensitive term per line in the ignored file `.public-audit-forbidden`, or set `PUBLIC_AUDIT_FORBIDDEN` to a comma-separated list. The scanner supplements manual review; it cannot determine ownership or identify every trade secret.

## Caveats

- More agents can add correlated noise and cost. Fusion can also discard valid minority findings.
- Character imitation and provider diversity must earn their context and latency through measured outcomes.
- The initial packets test a mechanism on a small synthetic sample. Generalisation requires frozen variants, clean controls, blinded human ratings, and replication on representative work.
- Character names and visual references belong to their respective rights holders. See [`NOTICE.md`](NOTICE.md).

## License

MIT for the original code and text in this repository, subject to the third-party notices in [`NOTICE.md`](NOTICE.md).
