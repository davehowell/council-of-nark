---
name: nark-matrix
description: Convene a task-specific council of opinionated, review-only specialists; run them independently, then fuse their structured findings into one ranked verdict. Use when asked to convene the council, run the Nark Matrix, obtain adversarial review, or compare several review lenses.
---

# Nark Matrix

The Nark Matrix is a fan-out/fuse review pattern. The controller selects relevant reviewers, gives each the same source artifact in an independent session, and fuses their findings. Character voices make the lenses memorable. The functional checklists do the work.

This is an experimental technique. Do not claim that multiple agents, providers, or fictional wrappers improve quality without measuring the task at hand.

## Roster

| Reviewer | Lens | Definition |
|---|---|---|
| HK-47 | simplicity and unnecessary machinery | `agents/hk47-reviewer.md` |
| K-2SO | correctness, tests, observability, edge cases | `agents/k2so-observability.md` |
| GLaDOS | architecture, contracts, migrations, optional synthesis | `agents/glados-architect.md` |
| C-3PO | security, personal data, permissions, explicit controls | `agents/c3po-compliance.md` |
| WALL-E | durable technical language and post-fusion translation | `agents/walle-simplifier.md` plus `skills/walle-ste/` |
| Bender | cost and compute | `skills/nark-matrix/personas/bender-finops.md` |
| Holly | long-horizon operational entropy | `skills/nark-matrix/personas/holly-entropy.md` |

All reviewers are review-only. They report findings and never edit files.

## Select the panel

Choose the smallest relevant panel, with at least two reviewers:

- **Plan or design:** HK-47 + GLaDOS. Add K-2SO for state, retry, or correctness risk. Add Holly for manual steps or temporary state.
- **Code or data-flow change:** K-2SO + HK-47. Add GLaDOS for contracts, Bender for material cost, and C-3PO for sensitive data or access boundaries.
- **Infrastructure or operational change:** C-3PO + Holly. Add K-2SO for rollback and observability, Bender for capacity/cost, and GLaDOS for coupling.
- **Runbook, ADR, README, or other durable prose:** add WALL-E.

Do not add an irrelevant reviewer to meet a provider quota. Provider diversity is a variable to test, not a substitute for a relevant lens.

## Run independent reviews

1. Freeze the exact artifact or diff.
2. Give every reviewer that same frozen material and a narrow question.
3. Start each review in a fresh session with no prior council output, answer key, or unrelated memory.
4. Disable editing tools. Prefer no tools when the artifact is inlined.
5. Run reviewers in parallel when the harness supports it.
6. Preserve each raw response before fusion.

Claude Code can load the files in `agents/` as custom subagents. Other harnesses can inline an agent definition into a fresh one-shot call. The Bender and Holly files are standalone prompt templates for cross-provider calls.

Every reviewer must return:

```text
## Findings (<reviewer>)
- severity: blocker | major | minor | nit
  location: <file:line or design section>
  claim: <one supported claim>
  consequence: <concrete effect>
  fix: <smallest effective change>
  confidence: high | medium | low
```

An empty block is valid. Never reward a reviewer for filling a quota.

## Fuse once

The controller is the default arbiter. Use GLaDOS in synthesis mode only when volume or genuine conflict warrants another call.

During fusion:

- merge findings with the same mechanism and consequence;
- retain a supported minority finding;
- discard unsupported claims and style-only preferences;
- resolve conflicts explicitly;
- keep the most specific location and smallest effective fix;
- rank by expected harm, not by personality;
- retain which reviewers raised each finding;
- bank material positives that should not regress.

Do not let the fuser inspect hidden evaluation keys. Save both the raw union and fused verdict when evaluating the method, because fusion can delete valid minority findings.

## Translate

After fusion, run WALL-E in translator mode when non-specialists need the result. Translation is separate from review and does not count toward the minimum panel size.

## Report

Return:

1. **The council speaks:** at most one short quote from each consulted character.
2. **Ranked verdict:** severity, location, claim, consequence, fix, raised-by, confidence.
3. **Resolved conflicts:** only when reviewers disagreed.
4. **Banked positives.**
5. **Bottom line:** `ship`, `fix-first`, or `redesign`.
6. **For the mortals:** optional WALL-E translation.

Keep the character quotes short. The structured verdict is the deliverable.
