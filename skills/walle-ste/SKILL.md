---
name: walle-ste
description: Simplify and proofread any technical writing against rules adapted from ASD-STE100 — Confluence pages, Jira tickets, PR review comments, PR bodies, READMEs, runbooks, ADRs, dbt/Cube descriptions, error strings. Use when asked to "simplify this doc", "de-slop this", "run an STE pass", "make this readable", or "run WALL-E over it", or when drafting/reviewing Confluence, Jira, or PR-comment prose. Review mode reports findings. Rewrite mode edits only when the user explicitly asks for a rewrite.
---

# WALL-E STE kernel

This file is the always-loaded kernel: 23 checks plus the machinery. The remaining 47 checks live
in `references/`. Every check has a stable id, a tier, a detection step, and a fix. Cite the check
id and the citation on every finding. A finding with neither is an opinion.

## Precedence

1. **The house style guide beats this document.** Write the conflict down instead of resolving it
   silently per document.
2. **Meaning beats every rule.** A fix that changes what the sentence claims is discarded (`1.2`,
   `9.1`). `drop` is not `delete` is not `truncate`. `stop` is not `terminate`. `retry` is not
   `rerun`.
3. **Understood on the first read is the acceptance criterion and the tiebreaker** (`9.1`). The
   reader is a time-poor, competent engineer whose first language is not English, reading once, on
   a phone, at 03:00. Clumsy to a native skimmer is the price of unambiguous to everyone else.

## Scope

**Apply to:** READMEs, service and API docs, design docs, ADRs, runbooks, postmortems, migration
guides, release notes, PR and commit descriptions, tickets, review comments, error and log
strings, code comments, dbt and Cube descriptions, Confluence pages, durable Slack (pinned
summaries, handovers, decisions).

**Never apply to:** code, identifiers, command syntax, quoted text (reproduce byte for byte,
including British spelling), marketing copy and essays (they need a voice), ephemeral chat, and
the glossary and style guide themselves. A formatter owns heading levels, line width, and bullet
characters. A template owns which sections a document must contain. These rules own the sentences.

The standard's ~900-word approved dictionary is deliberately not applied. It cannot carry a README
or an ADR. What survives is its mechanism: one meaning per word, one name per thing, a small
controlled general vocabulary plus an open project glossary.

If the project has a glossary or exception lists (phrasal verbs, `-ing` nouns), load them in the
mask stage. Without a glossary, terminology findings (E1–E4) are guesses — cap them at minor.

## Register dial

Set per block, after classification:

- **strict** — procedural text, safety blocks, notes inside procedures, error and log strings.
  Every applicable check runs.
- **flavored** — descriptive prose: READMEs, ADR context, PR bodies, model descriptions. Suppress
  D3 (add `that`), D8 (articles), C21 (repeat the key noun), and C28 (comma readings). Everything
  else runs.

## Stage 1 — classify each block

Classify per block, not per document. A document that mixes types is normal. A block that mixes
types is the defect (`3.6`, `5#`, `6#`).

| Type | Reader is | Verb form | Sentence cap |
|---|---|---|---|
| Procedural | doing | Imperative, verb first (`5.3`) | 20 words (`5.1`) |
| Descriptive | reading | Active declarative, no imperative (`6#`) | 25 words (`6.3`) |
| Safety (WARNING/CAUTION) | about to break something | Imperative or negative imperative (`7.2`) | 20 words (`5.1`) |
| Note | reading context | Declarative, never imperative (`5.5`) | 25 words, unlimited sentences |

Tests: procedural means the reader types, clicks, or runs something here. Safety: if they get it
wrong, is the state recoverable by re-running something? No means WARNING. Yes, with a documented
recovery path, means CAUTION. No damage at all means NOTE (`7.1`).

The caps are alternatives selected by block class. Never report "over 20 words" against a design
doc.

## Stage 2 — mask, then lint

Do this once, before any check runs:

1. Mask code spans, fenced blocks, and inline SQL.
2. Mask quoted strings, UI labels, log lines, URLs, and file paths (`1.14` note, `8.6`).
3. Mask identifiers: `snake_case`, `camelCase`, `kebab-case`, dotted paths, ticket keys, SHAs,
   image tags, instance types (`8.6`).
4. Mask heading lines before any `-ing` check. `[EXT]` interpretation: `3.5` permits `-ing`
   technical nouns and says nothing about headings.
5. Strip list markers and step numbers before any word count (`8.6`).
6. Load the project glossary and exception lists if they exist.

Every `rg` command below runs on the masked text. ripgrep's default engine covers everything here
except lookaround (`-P`, used only by D3).

## Limits and counting

| Limit | Value | Scope | Source |
|---|---|---|---|
| Procedural or safety sentence | 20 words | Per sentence | `5.1` |
| Descriptive or note sentence | 25 words | Per sentence | `6.3`, `5.5` |
| Paragraph | 6 sentences | Descriptive text only. A lead-in plus its list counts as one sentence | `6.6` |
| Prose noun cluster | 3 words | A hyphenated compound counts as 1 | `2.1`, `8.7` |
| Coined term or introduced short form | 3 words | Never abbreviate a term of 3 words or fewer | `1.9`, `2.2` |
| List lead-in and each item | The cap in force | Counted separately. Items share no budget | `8.4` |
| Safety block | 2–3 sentences | `[EXT]` — the source states no count, its examples do | — |

**One word each:** a number (with its unit: `250 ms`, `4 TiB`), an abbreviation, an identifier or
flag, a quoted span or code span, a name, title, or label (**never paraphrase a name to satisfy a
limit**), a proper noun, a hyphenated compound, and a whole parenthetical (also checked separately
as its own sentence). Structural numbering is excluded entirely (`8.5`–`8.7`).

Worked: ``Run `kubectl -n prod rollout restart deploy/feature-writer` and wait 90 seconds.`` is
5 words, not 11.

## Tiers

| Marker | Meaning |
|---|---|
| `[M]` | Mechanical detection, deterministic fix. Apply it. |
| `[M+J]` | Mechanical detection, semantic fix. One decision per hit. Never bulk-apply. |
| `[J]` | Judgement. The entry states the question to ask. |
| `[EXT]` / `[GR-n]` | Provenance flags: extension / advisory in the source. Cite them, never as rule ids. |

Severity comes from consequence in context, never from tier or provenance (see Output contract).

---

# Kernel checks

## Actors and verbs

### A1 · Passive voice `[M+J]` `3.6`
- Catch: `rg -in '\b(is|are|was|were)\s+\w+(ed|en)\b'`. Agentful hits (`… by X`) are highest
  precision. On an agentless hit ask: by whom or by what? An answer means a violation.
- Discard: `by` naming a dimension or sequence (`partitioned by`, `followed by`), participial
  adjectives (`the failed shard`, `broken`, `hidden`, `frozen` — `3.3`), and false participles
  (`open`, `often`, `even`, `seven`, `token`, `red`, `golden`).
- Fix: promote the agent to subject. In a step, use the imperative. Or supply `you` / `we`.
- Exception, and it overrides the active-voice rule: **never invent a false actor.** Passive with
  a genuinely unknown agent is correct in descriptive text ("the offsets were corrupted" while
  root cause is open). A style rule must never manufacture a causal claim.

### A5 · Nominalization `[M+J]` `3.7`, `1.13`
- Catch: `rg -in '\b(do|perform|conduct|carry out|make|give|provide|execute)\s+(a|an|the)\s+\w+(tion|ment|ance|ence|sis|al)\b'`.
  Also bare `the ask`, `a deploy`, `an ingest`, `the spend`.
- Fix: use the verb directly. "Perform a migration of the table" → "Migrate the table".
  Nominalization also hides the actor, so it compounds with A1.

### A7 · Verb names no action, object, or direction `[M+J]` `1.12`
- Catch: `rg -in '\b(handles?|process(es)?|manages?|deal with|touch|address|support|update)\b'`.
  Also any transfer or motion verb with no source and no destination (`sync`, `promote`,
  `migrate`, `push`, `restore`, `cut over`): ask what object, from where, to where.
- Fix: name the exact operation, the object, and the endpoints. `update` is the worst umbrella
  verb: it hides `INSERT`, `UPSERT`, `MERGE`, `TRUNCATE`, and `DROP`.
- Exceptions: `update` naming the SQL operation, a dependency update, or a security update.
  `supports` with a named version or capability is precise.

## Claims

### B1 · Modal that hides a fact `[M+J]` Part 2 2-0-19
- Catch: `rg -n '\b(may|might|could|shall|ought|would)\b'` — lowercase only, which keeps the
  RFC 2119 exception mechanical.
- Fix: `may` → `can` (possibility) or `is permitted to`. `shall` → `must`. Conditional → state the
  condition: "If the lockfile is stale, the build breaks."
- Never convert a modal into a bare categorical claim. "Might break" and "breaks" are different
  claims (precedence rule 2).
- Exceptions: uppercase RFC 2119 keywords in a spec. Counterfactual `would` comparing a rejected
  option in an ADR. Past-ability `could` ("could not scale past 8 shards").

### B3 · Abstract claim, no mechanism and no magnitude `[J]` `4.1`
- Ask: can a reader act on this, or predict the system's behavior from it? An effect with no
  mechanism and no magnitude fails. "Improves maintainability" → name what moved where, with
  numbers.

### B8 · Unverifiable abstraction `[M+J]` `1.5`, `[EXT]`
- Catch adjectives `[EXT]`: `rg -in '\b(seamless|robust|powerful|cutting-edge|effortless|world-class|next-generation|revolutionary|blazing|elegant|intuitive|first-class|battle-tested|enterprise-grade)\b'`
- Catch verbs: `rg -in '\b(enables?|empowers?|streamlines?|operationali[sz]e[sd]?|fosters?|unblocks?|drives? (adoption|alignment|value)|delivers? value)\b'`
- Catch nouns: `rg -in '\b(estate|cadence|fabric|posture|surface area|north star|alignment|learnings|synerg\w+|ways of working)\b'`
- Test: point at the command, artifact, config, or measurable change it names. None means
  decoration. Fix: delete, or state the measured property.

### B9 · Filler word `[M]` Part 2 2-0-13
- Catch: `rg -in '\b(just|simply|basically|actually|really|very|quite|essentially|of course|obviously|it is important to note that|it should be noted that)\b'`
- Fix: delete. `just` and `simply` also smuggle in a claim that the task is easy.
- `already` is `[M+J]`: keep it only where it names a prior stage or state.
- Owns these adverbs. D14 owns the adjective `actual`. `just-in-time` is a term — the hyphen is
  the boundary.

### B10 · Reference point not stated `[M+J]` `[EXT]`
- Catch: `rg -in '\b(currently|at present|recently|soon|in the near future|as of today|going forward|at this time)\b'`
  plus bare `old`, `new`, `stale`, `legacy`, `outdated` as adjectives.
- Fix: state the version, date, threshold, or condition. Documentation outlives its writing. A
  vague relative adjective is a data-loss risk in retention and TTL docs.
- Exception: `now` in a log line or a step ("the queue is now empty"). Judge each hit.

### B13 · Mixed or missing unit and time convention `[M+J]` `[EXT]`
- Catch: `rg -in '\b\d{1,2}\s?(am|pm)\b'`, bare `HH:MM` with no zone, `GB`/`GiB` mixed in one
  document, durations written two ways.
- Fix: ISO 8601 with an explicit zone, one byte-size convention, one duration format, stated once
  per project. Nothing else owns this check.

### B14 · Rhetorical padding `[M+J]` `[EXT]`
- Catch: `rg -in '\b(in summary|in conclusion|to summari[sz]e|overall,|as mentioned above|as previously stated|it is worth noting|that said,|at the end of the day)\b'`
  and `rg -in 'not only .{0,60} but also|is(n.t| not) just .{0,40}(—|--|,) it'`
- Fix: delete the closing-summary paragraph — topic sentences already form the outline (C15). A
  needed summary goes first, as a `TL;DR`. Rewrite reversal constructions as plain statements. A
  list of exactly three parallel items where two are enough, or four exist, is padding.

## Sentences and structure

### C1 · Over-length sentence, and the shape of the split `[M+J]` `5.1`, `6.3`, `4.4`, `6.2`
Absorbs retired C2 (guard rails) and C16 (reconnect).
- Catch: tokenize with the counting convention, compare against the block's cap: 20 procedural
  and safety, 25 descriptive and note. One cap per sentence, selected by block class.
- Fix: split at a clause boundary into complete sentences. Then reconnect: mark a causal,
  contrastive, or sequential relation explicitly (`then`, `so`, `but`, `as a result`). Unmarked
  causation is fatal in incident write-ups.
- Guard rails: keep every article, auxiliary, and subject — telegraphic prose is a failed fix.
  Never let a word count invent a step boundary: a step is a stopping point.

### C8 · Noun pile-up `[M+J]` `2.1`, `2.2`, `8.7`
- Catch: a run of four or more consecutive nouns and adjectives ending in a head noun. A
  hyphenated compound counts as one. Also any token with three or more internal hyphens.
- Fix, in this order:
  1. Unpack with a preposition, rewriting from the head noun backwards. If no preposition fits,
     the relationship is unknown — ask the owner.
  2. Rewrite a hidden action (`-tion`, `-ment`, `-ing`) as a verb.
  3. Hyphenate the sub-groups that genuinely bind.
  4. For an established name, write it in full once and introduce a short form of 3 words or
     fewer.
- When unsure which words bind, method 1 wins. Never guess with punctuation (`2.2`).

### C9 · Buried conditional `[M+J]` `5.4`, `[GR-2]`
- Catch: `if`, `when`, `unless`, `after`, `once` appearing after the main verb of a step. Also a
  leading condition with no closing comma.
- Fix: condition first, closed with a comma. A skimming on-call reader executes the verb without
  the condition.

### C15 · Paragraph and section structure `[M+J]` `6.4`–`6.6`, `9.4`
- Catch: more than 6 sentences in a descriptive paragraph (lead-in plus list = one sentence). Six
  is a ceiling, not a target.
- Outline test `[J]`: extract every heading and the first sentence of every paragraph, in order.
  That list must summarize the document. Failures: gaps, repeats, a heading doing the topic
  sentence's job, misordered sections. Fix section order before paragraphs.
- Opening sentence `[J]`: it must name the thing, say what it is for, and add something the name
  does not. `X is a service that Xs` is padding.

### D1 · Unresolvable reference `[M+J]` `[GR-3]`, `[GR-4]`, `4.4`
- Catch: `rg -in '\bthis\s+(is|are|can|will|causes|means|makes|results|happens|breaks|fixes|allows|helps|should)\b'`,
  sentence-initial `This`/`These` with no following noun, and any pronoun with more than one
  number-matching candidate in reach.
- Fix: write `this <noun>`, or restate the cause outright. An ambiguous `it` in a runbook sends
  someone to restart the wrong service.
- Not a violation: a document-scoped `This` in an opening or closing sentence.

## Words

### E1 · Synonym rotation `[M+J]` `1.11`, `6.2`
- Catch: more than one member of an alias cluster in one document: {customer, client, account,
  user} · {job, task, DAG, pipeline, workflow} · {service, app, component} · {table, dataset,
  model, view} · {feature, attribute, field, column} · {deploy, ship, release, roll out} ·
  {delete, remove, drop, purge}.
- Fix: one name for one thing, repeated every time. Variation reads as significance and sends
  readers hunting for a difference that does not exist.

### E5 · Abbreviation failures `[M+J]` `2.2` M1, `8.3`, `6.6`
- Catch: all-caps tokens of 2–6 characters with no `Full Term (ABC)` expansion earlier in the
  same document. Four defects: never expanded · the long form reappears after the definition ·
  an abbreviation for a term of 3 words or fewer · abbreviation soup inside a procedure (write
  terms out in steps even when defined — this overrides defect 2 there).
- The introduced short form must itself be 3 words or fewer.
- Exempt: terms universal in the field — `SQL`, `API`, `HTTP`, `JSON`, `CSV`, `UTC`, `ETL`,
  `GPU`, `CI`.

### E8 · Inflated word `[M+J]` `1.12`, Part 2 recurring errors
- Replace: utilize/leverage→use · commence/initiate→start · facilitate→help ·
  ensure/assure→make sure that · prior to→before · subsequent to→after · subsequently→then ·
  in order to→to · in the event that→if · a number of→some · regarding/concerning→about ·
  obtain/acquire→get · demonstrate→show · transmit→send · additionally/furthermore/moreover→also ·
  however/nevertheless→but · therefore/hence/consequently→as a result ·
  simultaneously/concurrently→at the same time · portion→part · instantiate→create ·
  orchestrate→run, schedule · surface (v)→show · action (v)→do · aforementioned→name the thing.
- Fix: apply the replacement, then reread. A different part of speech means restructure, not
  force the swap (`9.1`).
- Exceptions: `INSERT`, `EXECUTE`, `MATERIALIZE`, `ROTATE` as commands. `provision` in
  infrastructure. **`terminate` is not a synonym for `stop`:** terminating an AWS instance
  destroys it. Never rewrite the technical sense.
- Owns `ensure`, `in order to`, `surface (v)`. Do not also report them under D6 or B9.

## Procedures

### F1 · Step that is not a clean imperative `[M+J]` `5.3`, `3.6`
- Catch: `rg -in '^\s*[0-9]+[.)]\s+(The |You must|You should|You need to|You will need to|You have to|It is necessary to|We |.*\b(is|are) to be\b)'`
- Fix: start every step with an imperative verb. Delete `you must` from steps — a step is already
  an order. Reserve `must` for safety blocks and stated limits: once every step says must, the one
  that truly must is invisible.
- A step can open with its condition (C9). F1 rejects a step whose command is not imperative.

### F4 · Load-bearing note `[M+J]` `5.5`
- Catch: `NOTE:` / `[!NOTE]` / `[!TIP]` blocks containing an imperative (including `Make sure`,
  `Ensure`, `Be aware`, `Remember to`) or `you must` / `you need to`.
- The delete-the-notes test: strip every note and tip, read the remaining steps, and ask whether a
  competent stranger can still finish. Any gap is a step misfiled as commentary. Run it as a
  ritual, not a judgement call.
- Fix: rewrite the note as an imperative step, insert it at the right point (usually before the
  step it was attached to), re-run the test.

## Safety

Attach a safety block to: `DROP` / `TRUNCATE` / `DELETE` with no `WHERE` · `git push --force` ·
`terraform apply|destroy` in prod · `kubectl delete ns|pvc` · `dbt run --full-refresh` · an
unbounded backfill or replay · key or credential rotation · a prod schema migration · Kafka topic
deletion · an index rebuild.

### G1 · Destructive operation with no warning `[M+J]` `7#`
- Catch: scan for the trigger list. For each hit, check that a safety block sits immediately
  before that step.
- Fix: attach one. Danger must be a labeled element, not a tone of voice. The block must pass G2.

### G2 · Safety block shape `[J]` `7.1`–`7.3` — runs on every safety block, never skipped
- Shape: `<LABEL>: <command or condition>. <consequence>.` Three slots, fixed order, none
  optional.
  1. Label first token: `WARNING:` (irreversible or externally visible) · `CAUTION:` (recoverable
     with a documented path) · `NOTE:` (no damage).
  2. First sentence: an imperative, negative imperative, or condition clause. Never a noun
     phrase, never `There is`, never a passive.
  3. Consequence: what happens if the reader ignores it. Omit only if genuinely uncharacterizable,
     and then say what is uncertain (`7.3`).
- Case is a house choice (`7.1` note). The label word carries the severity, not the formatting.

### G3 · Vague warning `[M+J]` `7.1`, `7.3`
- Catch: `rg -in 'be careful|take care|use caution|proceed with caution|is critical|is imperative|is essential|at your own risk|may have unintended consequences'`
  plus softened consequences inside a safety block: `rg -in '\b(impacted|affected|inconsistencies|unexpected behaviou?r|issues|problems|degraded experience)\b'`
- Fix: name the action and the outcome with blunt words: `deleted`, `lost`, `unrecoverable`,
  `outage`, `leaked`, `corrupted`, `overwritten`, `stale`. Softening the harm is what makes a
  warning get skimmed.

---

# The catalog

**Hard gate: never report a finding for a catalog check without first reading its entry in the
reference file.** The entries carry the carve-outs that stop this tool from mangling correct
prose. Kernel checks above are exempt — their carve-outs are co-located here.

| File | Checks | Open it when you see |
|---|---|---|
| `references/checks-words.md` | E2 E3 E4 E6 E7 E13 A6 D10 | one action worded three ways, paraphrased asset names, phrasal verbs, slang, capitalized common nouns, verbed nouns, ambiguous modifier binding |
| `references/checks-grammar.md` | C4 A3 A10 C7 B2 A9 C25 B4 B5 B12 C12 C13 C27 C14 D8 | tense drift, modal passive, `-ing` chains, dropped subjects, contractions, unquantified degrees, broken comparatives, prohibitions with no action, inline enumerations, list mechanics, mixed modes, articles |
| `references/checks-structure.md` | F2 F3 F5 F7 C11 C21 C24 | two actions per step, unnumbered steps, limits hidden in notes, imperatives in overviews, tangled or disconnected paragraphs, prose duplicating a table |
| `references/checks-safety.md` | G5 G6 G7 G9 G10 G11 | wrong severity, mixed label vocabulary, warning after the command, over-long warnings, warnings in reference docs, prohibition only in a lead-in |
| `references/checks-punctuation.md` | C10 C28 C18 E15 D3 D4 D6 D11 D12 D13 D14 | semicolons, ambiguous commas, heavy parentheticals, hyphen questions, dropped `that`, `with`/`using`, `any`/`since`/`over`, `e.g.`, gendered terms, apostrophes, false friends |

`references/registry.md` is the authoritative index: every id, tier, citation, home, and the
retired list. If any file and the registry disagree, the registry wins.

## Rewriting method (`9.1`)

- Rewrite a defective sentence whole, not one flagged word at a time. Rules interact.
- A substitute must keep the meaning and the part of speech. Awkward output means the
  construction is the problem, not the word.
- Ask two questions before rewriting an unclear sentence: what does the vague word mean here, and
  what must the reader do?
- When a sentence resists rewriting: split it, delete information that is not necessary, or get
  the missing fact from the engineer who owns the system. Never simplify by deleting the
  information a sentence carries.
- After edits, grep for what the edit broke: renamed terms, "see step 4", "as described above".

## Stop rule

Stop and report what you stopped on when any of these is true:

1. The sentence is understood on the first read. Further edits are taste (`9.1`).
2. The next change alters the technical claim. Discard the fix (`1.2`).
3. The remaining defect needs a fact you do not have. Name the fact.
4. The check's own entry carves the hit out.
5. A formatter, a template, or the house style guide owns the finding.
6. You rewrote the same sentence twice and it is still unclear. Split it, cut it, or ask the
   owner.
7. The findings exceed what a reviewer will read. Rank by severity and cut the tail.

## Self-lint order

1. Classify every block, set the register dial, run the mask, load the glossary.
2. Generative checks first, so the text they create gets linted: G1 (attach missing safety
   blocks), G2 (shape them), F4 (promote misfiled steps), C12 (build vertical lists — catalog).
3. Run `[M]`, then `[M+J]`, then `[J]`. Within a tier: kernel checks first, then any catalog file
   the symptoms made you load.
4. Re-run step 3 once over any text you created or restructured. One re-run only. If the second
   pass flags your own output, report the finding and stop.
5. One span, one finding. Deduplicate to the most specific check. Entries name the owner of
   contested words.
6. Apply the stop rule.

## Output contract

Severity is set by the consequence in context, never by tier or provenance:

- **blocker** — the defect can cause a wrong action on a destructive or irreversible path, or the
  only available fix would change a technical claim (report it, never apply it).
- **major** — actor, referent, or condition ambiguity in procedural, safety, or incident text.
- **minor** — clarity and consistency defects in descriptive text.
- **nit** — polish, and `[EXT]` / `[GR-n]` advisory items.

**Review mode (default).** Findings only, most severe first. Each finding: severity, check id,
citation (rule id, `[EXT]`, or `[GR-n]`), the offending span, the rewrite, one line of rationale.
Never report a carved-out hit or something a formatter owns.

**Rewrite mode.** Only when the user explicitly asks for a rewrite. Return only the requested
text — no preamble, no summary. If a sentence needs a fact you do not have, append one line
naming the fact. Council and subagent runs never rewrite: they report findings only.

## What this cannot do

The checks remove the form of slop. They cannot make a hollow paragraph true, and they cannot
tell whether a runbook states the rollback procedure at all. That needs a reader who knows the
system.

## Attribution

Adapted from ASD-STE100 Simplified Technical English, Issue 9 (2025-01-15), copyright © ASD
(AeroSpace, Security and Defence Industries Association of Europe), <https://asd-ste100.org>.
This adaptation is unofficial and is not endorsed by ASD. It adapts the writing rules for
software, data, and ML engineering prose and deliberately omits the approved-word dictionary.
Items marked `[EXT]` are additions and are not part of the standard. Get the real standard from
ASD.

walle-ste v0.2 · 2026-07-31
