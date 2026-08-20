# Safety instructions

Catalog entries. Read the whole entry before reporting — the carve-outs travel with the check.
All commands run on masked text (SKILL.md Stage 2). Tiers and citations: `registry.md`.

G1 (missing block), G2 (block shape — mandatory on every safety block), and G3 (vague warning)
live in the kernel. Any block these checks create or rewrite must still pass G2.

### G5 · Wrong severity `[M+J]` `7.1`
- Ask: what is the worst realistic outcome if the reader gets this wrong? Compare with the
  label. A mismatch in either direction is a violation — a WARNING on a 40-minute backfill is a
  NOTE, and a CAUTION on a key rotation that breaks every running job is a WARNING.
- Every label must be defensible by naming the outcome that justified it.

### G6 · Label vocabulary and rendering `[M+J]` `7#`, `7.1`
- Catch: `rg -in '^\s*>?\s*\**(WARNING|CAUTION|DANGER|NOTE|ATTENTION|NOTICE|IMPORTANT|HEADS UP)\b'`
  across the corpus. More than one vocabulary is a violation. Also a block whose severity is
  conveyed only by color, bold, or an emoji.
- Fix: one severity vocabulary for the whole documentation set, defined in one place. Reserve
  `WARNING` so it is not spent on trivia. The label word is the payload — formatting does not
  survive GitHub, mkdocs, a terminal, a log aggregator, and a Slack unfurl.

### G7 · Warning placement and split severity `[M+J]` `7.1`, `7.2`
- Catch: a destructive command whose nearest preceding element is not its safety block. A
  warning after the command, or only in the document preamble. Also two adjacent blocks of
  different severity on one step, or a CAUTION whose body names a top-severity consequence.
- Fix: move the block immediately before the step it protects. Merge split labels and keep the
  higher severity. Both moves change what the reader sees before acting — check that the merged
  block still names every consequence, and never relocate the recovery command out of reach.

### G9 · Over-long warning `[M+J]` `5.1` + `[EXT]`
- Scope: WARNING and CAUTION only. A NOTE follows `5.5` — 25 words per sentence, unlimited
  sentences.
- Catch: more than three sentences, any sentence over 20 words, or a code fence, table, or
  nested list inside the block. (`[EXT]`: the source states no sentence count. Its worked
  examples are all two or three short sentences.)
- Fix: one command or condition plus the consequence. Move rationale and the recovery procedure
  to the body or a link — but confirm the recovery path is still reachable from the block. A
  warning that deletes its own recovery command is a worse warning.

### G10 · Warning in a descriptive document `[J]` `7#`
- Catch: a WARNING in a model description, architecture overview, API reference entry, or ADR
  context section, with no command or step nearby.
- Fix: move it to the runbook and link to it. Warnings scattered through reference docs go stale
  and get skimmed.

### G11 · Prohibition only in the lead-in `[M]` `4.3`
- Catch: a warning list whose lead-in ends `do not:`, `never:`, or `avoid:`.
- Fix: repeat the prohibition at the start of each item, so an item copied out on its own still
  reads as a prohibition.
