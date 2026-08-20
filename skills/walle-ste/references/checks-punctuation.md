# Punctuation and general recommendations

Catalog entries. Read the whole entry before reporting — the carve-outs travel with the check.
All commands run on masked text (SKILL.md Stage 2). Tiers and citations: `registry.md`.
`GR-1` to `GR-8` are advisory in the source: cite them as `[GR-n]`, never as rule ids.

### C10 · Semicolon `[M]` `8.1`
- Catch: `rg -n ';'` The mask already excludes SQL, code, CSS, HTML entities, and URLs — the ban
  applies to prose only.
- Fix: two sentences, or a vertical list.
- **Only the semicolon is banned.** Commas, colons, parentheses, question marks, hyphens, and
  dashes are all permitted. The em dash is not banned by this standard — if the house wants it
  gone, that is a house rule, labeled `[EXT]`, never attributed to ASD-STE100.

### C28 · A comma that permits two readings `[M+J]` `8#`, `5.4`
Suppressed in the flavored register.
- Restrictive versus non-restrictive: `, which` describes ALL of the preceding noun. Bare `that`
  restricts to SOME. "Drop the tables, which are partitioned by date" drops every table. Ask
  which reading is meant, then punctuate for it.
- Serial comma: in a series of three or more, put a comma before the final `and`/`or`.
- Adverb at a comma boundary `[J]`: "If the replica does not sync, cleanly restart the pod" is
  not "If the replica does not sync cleanly, restart the pod". Move the adverb away from the
  boundary, or delete it.

### C18 · Parenthetical carrying a whole sentence `[J]` `8.3`
Parentheses are for: a cross-reference (`see ADR-014`) · a labeled diagram or listing element,
named exactly as the artifact labels it — never invent callout numbers in prose · a step number ·
an abbreviation gloss at first use (`Change Data Capture (CDC)`) · a short factual clarification
(`no more than 512 rows per step`) · one strictly parallel alternative pair, at most.
- Anything else is commentary the writer did not place. Promote it to a sentence or delete it.
- `[EXT]` inversion: write the plural instead of `test(s)` — engineering prose has no
  legal-precision constraint, and `(s)` reads as refusal to commit. Never cite this as `8.3`.

### E15 · Hyphenation `[M+J]` `8.2`
- Hyphenate: a multi-word modifier before its noun (`read-only replica`, `end-to-end test` —
  but "the pipeline is up to date" after the noun takes none) `[J]` · spelled-out compound
  numbers (`forty-seven`) `[M]` · a letter or numeral bound to a noun (`3-node cluster`,
  `B-tree index`) `[M]` · noun-plus-verb compound verbs (`force-push`, `fine-tune`,
  `cross-validate`) `[M+J]` · a prefix vowel before a root vowel (`re-enable`, `pre-aggregate`,
  `de-identify`) `[M]`.
- Carve-outs: established closed forms take no hyphen — `backfill`, `rollback`, `checkout`,
  `logout`, `upsert`, `autoscale`, `preprocess`, `reindex`, `coordinate`. Keep a project word
  list.
- A hyphen is not a dash (`8.2` note): the hyphen joins, the en dash marks a range, the em dash
  sets off an interjection. In ASCII-only contexts (log lines, commit messages), a spaced hyphen
  is the correct rendering of a dash, not a violation.

### D3 · Dropped `that` `[M+J]` `[GR-1]`
Suppressed in the flavored register.
- Catch: `rg -inP '\b(make sure|verify|confirm|shows?|recommends?|assumes?|notes?)\s+(?!that\b)(the|a|an|you|it|we|this)\b'`
  The only check needing `-P` (lookaround).
- Fix: keep the conjunction — it marks where the main clause ends. Keeping `that` is the
  opposite of most editing advice, and it is deliberate: the reader's first language is not
  English.

### D4 · Buried or ambiguous instrument: `with`, `using`, `Use X to` `[M+J]` `[GR-2]`, `3.5`
- Catch: `with` carrying a condition (`rg -in '\bwith\s+(the\s+|a\s+|an\s+)?\w+(\s+\w+)?\s+(enabled|disabled|running|paused|stopped|set|configured|applied|held|open)\b'`),
  dangling `using` clauses, and step-initial `Use X to do Y`.
- Ask, on a flagged `with`: which reading is meant — `that has`, `together with`, or
  `by means of`?
- Fix: lead with the real action verb and name the tool after it ("Use `kubectl` to scale the
  deployment" → "Scale the deployment with `kubectl`"). Convert a condition-carrying `with` into
  a `When …` clause. Split dangling participles into finite sentences.
- Carve-outs: instrumental `with` naming an unambiguous tool or flag is the mandated
  construction — never flag it. `using` as an ordinary participle modifying a noun ("the pods
  using the old image") is fine — the finding is the dangling participle only.

### D6 · Multi-sense ordinary words `[M+J]` `9.2`, Part 2
Flag the word only where the sense is genuinely ambiguous, then write the replacement for the
sense meant:
- `any` → `all`, `each`, `one`, `some`. Dangerous only where the sentence deletes, filters, or
  thresholds ("Delete any duplicate rows" — one copy or all copies?). `any` in a universal
  statement ("flag any run of four or more") is fine.
- `over` / `under` / `above` / `below` → `more than` / `less than` for quantities, `during` for
  time. State whether a bound is inclusive. Document-position uses ("see the table below") are
  fine. Never flag SQL window syntax (`OVER (PARTITION BY …)`) or prose discussing it.
- `since` → `because` for cause, `after` or a date for time. `while` / `as` / `once` → name the
  relation.
- `complete` → `completed` for finished ("the load is complete" versus "the dataset is
  complete" are different claims).
- `check` → keep it where the object and the expected result are both named. Flag bare "check
  the pipeline". The periphrastic "do a check of" stays dropped (dictionary artifact).
- `repeat` → `do steps 3 to 5 again`, plus a count or stop condition.
- `people` / `users` / `someone` / `the team` → the specific role.
- No figurative physical verbs: `reach`, `hit`, `land`, `unlock`, `bubble up`,
  `move the needle`, `turns green` → the literal claim.
- Owned elsewhere: `ensure`, `in order to`, `surface` (E8) · `with`, `using` (D4) · `old`,
  `new` (B10) · `actually` (B9). Do not double-report.

### D11 · Latin abbreviation `[M]` `[GR-6]`
- Catch: `rg -n '\b(e\.g\.|i\.e\.|etc\.|viz\.|cf\.|N\.B\.|et al\.)'`
- Fix: `for example`, `that is`, `and so on` — or delete. `etc.` at the end of a list usually
  means the writer stopped thinking.

### D12 · Gendered pronoun and exclusionary term `[M+J]` `[GR-3]`, `[GR-7]`
- Catch: `rg -in '\b(he|she|his|her|hers|him|manned|manpower|man-hours)\b'` plus the house word
  list (`[EXT]` — the source defers the list to external guidance): `master`/`slave`,
  `whitelist`/`blacklist`, `sanity check`, `dummy value`.
- Fix: flag a gendered pronoun only where the referent is generic — address the reader as `you`
  or name the role. **A named individual keeps their pronoun**: "Priya merged the change. She
  then paged the on-call." is not a finding (`GR-3` scopes the prohibition to a generic person).

### D13 · Possessive apostrophe `[M+J]` `[GR-8]`
- Catch: `rg -n "'s\b|s'\b"` Check singular versus plural placement, and `it's`/`its`.
- Fix: rewrite with `of` or a noun modifier when in any doubt. Names ending in `s` (`users`,
  `events`, `orders`) make the apostrophe genuinely ambiguous.

### D14 · False friends `[J]` `[GR-5]`
Reviewer awareness for a multilingual team, not a self-check.
- `actual` for `current` gets its own lint `[M+J]`: `rg -in '\bactual\b'` → "the current schema
  version". Owns the adjective. B9 owns the adverb `actually`.
- Watch: `eventual` for `possible`, `sensible` for `sensitive`.
