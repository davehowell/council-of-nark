# Verbs, sentences, and lists

Catalog entries. Read the whole entry before reporting — the carve-outs travel with the check.
All commands run on masked text (SKILL.md Stage 2). Tiers and citations: `registry.md`.

The six permitted verb forms (`3.2`): infinitive · imperative · simple present · simple past ·
simple future · past participle as an adjective. A technical term does not license the passive or
an odd tense (`1.12`).

### C4 · Verb form outside the six permitted `[M+J]` `3.2`, Part 2 2-0-7
- Catch: `rg -in '\b(have|has|had)\s+\w+(ed|en)\b'` · `rg -in '\b(am|is|are|was|were)\s+\w+ing\b'` ·
  `rg -in '\b(will|would)\s+have\b'` · `rg -in '\b(being|been)\b'`
- Fix: rewrite into a permitted form. Part 2 bans the tokens `being` and `been` outright, so
  rewrite around them. Perfect and progressive tenses scramble incident timelines — simple past
  in order is what makes a timeline recoverable.
- Carve-outs: `is running` as a live-state report in a status line or log message (do not extend
  to prose). `will be available` and `to be deprecated` are permitted forms.

### A3 · Modal passive, requirement with no actor `[M+J]` `3.4`
- Catch: `rg -in '\b(can|must|should|will|may|is to|are to|needs to|has to)\s+be\s+\w+(ed|en)\b'`
  and `rg -in '\b(is|are) (required|necessary|needed|mandatory)\b'`
- Apply A1's false-participle discard list.
- Fix: decide whether the sentence instructs or describes. Procedural → imperative ("Grant the
  role"). Descriptive → name the actor, simple present or future. "The secret must be rotated"
  never says whose job it is — the single most useful lint for runbooks.
- Carve-out: uppercase `REQUIRED` in an RFC 2119 spec.

### A10 · The document avoids `you` and `we` `[J]` `3.6` Method 4
- Catch: a document full of agentless passives with no first or second person anywhere. The
  absence is the finding.
- Fix: `you` for the reader, `we` for your own team. Most house styles ban both. For runbooks
  and READMEs the standard is right and the house style is wrong — cite this entry.

### C7 · `-ing` forms: verbal, chained, and dangling `[M+J]` `3.5`
- Catch: an `-ing` token after a form of *be*, or opening a clause after a comma. Two or more
  verbal `-ing` forms in one sentence is a chained-participle finding.
- Discard: `, including` · `, depending on` · `, according to` · `, regarding` — prepositions,
  not participles. Also `-ing` technical nouns and modifiers: `staging`, `streaming`,
  `training`, `logging`, `caching`, `polling`, `sharding`. Keep a project allowlist and diff
  against it (`[EXT]` — that is how you implement `3.5` in CI).
- Fix: rewrite to a finite verb. Split a chained sentence into a lead-in, a numbered list, and a
  separate consequence.

### B2 · `should` used to state a requirement `[M+J]` Part 2
- Catch: `rg -in '\bshould\b'` Ask: mandatory, or genuinely optional?
- Fix: `should` → `must` for a requirement. If genuinely optional, say so and state what happens
  when it is skipped ("CI fails without it").
- Carve-out: RFC 2119 uppercase `SHOULD`. A blanket ban destroys a load-bearing distinction.

### A9 · Dropped noun, verb, or subject `[M+J]` `4.2`
- Catch: sentences with no noun phrase before the main verb (outside imperatives). Steps with no
  verb: `X to Y`, `X = Y`. Also `rg -in '^\s*If (installed|enabled|present|set|configured|missing|null|absent)\b'`
- Fix: restore the missing part, naming the real subject — "If enabled, drop the tables" → "If
  the `cleanup_staging` flag is enabled, drop the staging tables". Never shorten by deleting
  grammar.

### C25 · Contraction `[M]` `4.2`
- Catch: `rg -in "\b\w+['’](t|re|ve|ll|d|m)\b"` Possessive `'s` belongs to D13.
- Fix: expand in full. Matters most in error and log strings, which get read under pressure and
  grepped.

### B4 · Unquantified degree word `[M+J]` `4.1`, Part 2 2-0-8
- Catch: `rg -in '\b(significantly|substantially|dramatically|considerably|frequently|regularly|a lot of|plenty of|large|small|slow|fast|heavy)\b'`
  Suppress hits with a figure, unit, or version in the same sentence.
- `[EXT]` extension: also `possibly`, `potentially`, `generally`, `typically`, `tend to`. Not
  modals — never cite them under B1.
- Fix: the number, threshold, unit, or version.
- Carve-out: never flag inside an established or hyphenated term: `large language model`,
  `low-latency service`, `high-throughput scoring`, `fast-forward merge`, `slow query log`
  (`8.2`, `1.5`).

### B5 · Comparative defects `[M+J]` `4.1`, `1.5`, Part 2 2-0-8
- Catch, double comparative: `rg -in '\bmore \w+er\b|\bmost \w+est\b'` — delete the `more`/`most`.
- Catch, comparative on an absolute: `more|most|less|least` + `unique`, `complete`, `correct`,
  `deterministic`, `idempotent`, `final`, `optimal`, `null`. The property is true or false.
- Ask, no baseline `[J]`: faster than what? An unbaselined comparative is meaningless — give
  both numbers.
- Carve-outs: `more accurate` (continuous metric, quantifiable). `more than 500 ms` /
  `less than 1000 rows` are the mandated D6 replacements for `above`/`below` — never flag them.

### B12 · Prohibition stated instead of an action `[M+J]` `4.1`
- Catch: `rg -in '\b(is|are) not (permitted|allowed)\b|must not be present|should be avoided|\bavoid \w+ing\b'`
- Fix: write the check or command the reader performs, and the consequence: "Avoid concurrent
  migrations" → "Do not run two migrations at once. The second fails on the advisory lock."

### C12 · Inline enumeration `[M+J]` `4.3`, `6.4`
- Catch: three or more comma-separated noun phrases before a final `and`/`or`, or two or more
  actions joined by `then`.
- Fix: end the lead-in with a colon, build a vertical list.
- Carve-out: a short coordinated series of bare noun phrases (3 words or fewer each) inside the
  sentence cap stays inline. Convert when an item carries its own clause, when items are actions,
  or when the series pushes the sentence over its cap. D8 owns articles in a series that stays
  inline.

### C13 · Broken vertical list `[M+J]` `4.3`
Each mechanic separately checkable:
- The lead-in ends with a colon.
- Every item carries a marker: numbers when order matters, bullets when it does not. A bulleted
  list of sequential steps is a finding.
- Items start uppercase, unless the item opens with a case-sensitive identifier or flag.
- No item ends with a comma or semicolon (`rg -n '[,;]\s*$'` on list lines) — this also kills
  the trailing `and`/`or`.
- One indentation level. `[EXT]` Two tolerable for a genuine taxonomy. Three is a rewrite.
- Deferred to the house guide (precedence rule 1): terminal-period policy. `4.3` says period on
  full sentences only. Google and Microsoft say all-or-nothing. Report as suggestion only, and a
  list of bare identifiers takes none — a trailing period gets mis-copied.

### C27 · List items do not fit the lead-in `[J]` `4.3`, `1.5` note
- Catch: concatenate the lead-in with each item and read the pair as one sentence.
- Ask: does every pairing parse? Does every item answer the question the lead-in asked? Does the
  lead-in claim a completeness the list lacks ("The supported types are:" over a partial list —
  write "include", and point at the full list).
- Fix: rewrite the lead-in or the item. No item in a different grammatical shape from its
  siblings.

### C14 · Modes mixed in one block or one list `[M+J]` `4.3`, `5#`, `6#`
- Catch: one list holding both imperative-initial and declarative items. A section headed
  *Overview* or *Architecture* containing bare-verb sentences.
- Fix: facts in one list under a descriptive heading, steps in another under a procedural one.

### D8 · Articles: presence, absence, and adjective scope `[M+J]` `4.5`, `4.2`
`4.5` is conditional — half the rule is about when *not* to add an article. Suppressed in the
flavored register.
- Add an article or demonstrative before a countable noun in subject or object position:
  "Restart scheduler pod" → "Restart the scheduler pod".
- No article before a general class, mass concept, or abstract quality: "Vector embeddings can
  improve retrieval quality" is correct — over-correction is a real failure mode.
- No definite article before a noun directly followed by an identifier, version, or ticket key
  `[M]`: "Restart pod `web-7f3a`", "reopen ticket DATA-4821". The pair is a proper name.
- Two coordinated nouns in a short sentence: repeat the article ("Remove the partition and the
  index" — also kills a noun/verb ambiguity). A series of three or more with no adjective at
  stake: one article before the first noun is enough.
- Adjective scope `[J]`: repeat the article before every item a leading adjective does NOT
  modify. "The new partitioned tables, the views, and the scheduled queries" (only the tables
  are new). Getting this wrong in a migration runbook recreates the wrong objects.
- Disputes resolve against a named grammar reference or the house guide, not by preference.
