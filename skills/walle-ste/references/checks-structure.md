# Procedural and descriptive structure

Catalog entries. Read the whole entry before reporting — the carve-outs travel with the check.
All commands run on masked text (SKILL.md Stage 2). Tiers and citations: `registry.md`.

### F2 · Two actions in one step `[M+J]` `5.2`
- Catch: step sentences with ` and ` or ` then ` and an imperative verb on both sides. Ask: both
  at once, or one after the other?
- Carve-out 1: actions performed as a single operation stay together — "Drop and recreate the
  partition", "Stop and restart the scheduler", "Tag and push the image". Test: can the reader
  pause between them and leave the system in a sane state? If not, they are simultaneous.
- Carve-out 2: a second sentence stating the immediate result or acceptance limit is permitted —
  "Run the smoke suite. Every test must exit 0." That is the correct home for an acceptance
  criterion, not a note.

### F3 · Unnumbered steps `[M+J]` `5.2`
- Catch: procedural sections whose steps are unnumbered bullets or run-on prose. Also a step
  that packs several actions to keep the list short.
- Fix: number every work step, and use as many steps as the task needs. During an incident,
  responders report progress by step number. Where the flow is genuinely unordered, say so
  instead of implying an order.

### F5 · Requirement or danger in a note `[M+J]` `5.5`
- Catch, requirement: note blocks containing `must`, `at least`, `no more than`, `maximum`,
  `minimum`, or a number with a unit.
- Catch, danger: a note containing `lose`, `delete`, `overwrite`, `irreversible`,
  `cannot be recovered`, `downtime`, `corrupt`, `production`.
- Fix: a note gives information only. Move a threshold or expected result into the step it
  qualifies. Promote a danger to a safety block (G2 shapes it). Severity is decided by
  consequence, not by the grammar the author happened to use.

### F7 · Wrong construct in descriptive text `[M+J]` `5.5`, `6#`
- Catch: a sentence-initial bare verb inside an overview, architecture, or model-description
  block. Also a callout in a document with no steps and no adjacent diagram or table.
- Fix: restate the imperative as fact ("Register the feature group before training" → "A feature
  group must be registered before a training job can reference it"), or move it under a
  procedural heading. Fold stray asides into the prose.

### C11 · One topic and one new idea per sentence `[J]` `6.1`, `4.1`
- Catch: more than one finite clause with a distinct subject (`X does A, and Y does B`,
  `, which` chains that shift topic). Proxy: more than one new proper noun or system name
  introduced per sentence.
- Ask: did you re-read any sentence to work out what refers to what?
- Fix: one subject per sentence. Restate the subject noun rather than leaning on a bare `it`.

### C21 · Carry key words forward `[M+J]` `6.2`, `6.5`
Suppressed in the flavored register.
- Catch: adjacent sentences sharing no content word and no referring pronoun. A paragraph
  opening that shares nothing with the previous paragraph and has no connector.
- Fix: repeat the key noun. This is the opposite of the instinct to vary wording — repetition
  binds technical prose and makes it retrievable by search and by a model. Say this out loud to
  writers. They will resist it.

### C24 · Prose duplicating a table `[J]` `9.1`
- Catch: an introductory sentence that enumerates the same values or thresholds as the table it
  introduces.
- Fix: cut the prose to a pointer. Only one copy ever gets updated — duplication is the top
  source of drift in runbooks and SLA docs.
