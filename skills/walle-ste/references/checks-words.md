# Words and multi-word nouns

Catalog entries. Read the whole entry before reporting — the carve-outs travel with the check.
All commands run on masked text (SKILL.md Stage 2). Tiers and citations: `registry.md`.

### E2 · One wording per kind of step, commit, or ticket `[M+J]` `9.4`
Stronger than E1: the same sentence shape for the same action, not merely the same noun.
- Catch: extract the leading verb of every step, commit subject, ticket title, and error-message
  family. Near-duplicates naming one action ("Trigger the DAG" / "Kick off the job" / "Run the
  pipeline") are violations.
- Fix: fix the wording once and reuse it. Variation reads as significance.
- `[EXT]` A published ~30-verb list for step openings and commit subjects (`add`, `remove`,
  `fix`, `rename`, `move`, `update`, `revert`, `document`, `test`, `start`, `stop`, `deploy`,
  `roll back`) makes this enforceable in a commit hook.

### E3 · One word, two meanings `[M+J]` `1.3`, `9.2`
- Catch: inspect high-risk domain words for sense drift within one document: `model`, `load`,
  `run`, `job`, `key`, `partition`, `batch`, `stage`, `source`, `target`, `refresh`, `sync`,
  `client`, `feature`, `service`.
- Fix: one meaning per word per document. Two genuine domain senses ("the dbt model" and "the
  churn model") get separated contexts, each defined at first use, never both unqualified in one
  paragraph (`1.6`).

### E4 · Canonical name and canonical spelling `[M+J]` `1.6`, `1.8`, `2.2` M2
- Catch: grep referenced asset names against the schema, manifest, service registry, or glossary.
  Flag paraphrases and near-misses. A paraphrased table, dashboard, or endpoint name is
  unsearchable.
- Fix: use the exact name. Never swap a synonym into an established term (`primary key`, not
  `main key`). Never strip a hyphen from a canonical spelling (`dead-letter queue`,
  `write-ahead log`).
- Reverse error: do not hyphenate a conventionally open term (`the daily aggregate table`). Keep
  conventional attributive hyphens (`real-time pipeline`, `read-only replica`). Do not strip
  hyphens that C8 inserted to group a sub-cluster.

### E6 · Phrasal verb `[M+J]` `9.3`
- Catch: `rg -in '\b(spin up|tear down|kick off|stand up|bring up|back out|fall over|blow up|take down|hook up|wire up|look into|sort out|figure out|give off)\b'`
- Fix: one precise verb: `start`, `delete`, `revert`, `fail`, `investigate`, `connect`, `emit`.
  Then reread — a single-word replacement often needs a changed object or preposition (`9.3`).
- Exception: keep a written exception list of phrasal verbs that are genuine technical terms:
  `roll back`, `back up`, `roll out`, `fall back`, `log in`, `time out`. The list is what makes
  the rule enforceable without absurd rewrites.

### E7 · Slang and in-group jargon `[M+J]` `1.10`
- Catch: `rg -in '\b(nuke|blow away|yeet|wedged|hosed|borked|yak-shave|bikeshed|footgun|fire drill|over the wall|sideways|hairy)\b'`
- Fix: the plain, widely understood word plus the concrete referent: "nuke the release" → "roll
  back to `v1.42.3`".
- Carve-outs: `magic number`, `magic bytes`, `brick the device`, `flaky test`,
  `thundering herd` are established terms with precise meanings. Flag only the loose sense.

### E13 · Capitalize only when necessary `[M+J]` `1.5` note
- Catch: mid-sentence capitalized common nouns (`Data`, `Pipeline`, `Model`, `Service`,
  `Cluster`) checked against a known proper-noun list.
- Fix: capitals are for proper nouns, product names, acronyms, and exact identifiers, never for
  emphasis. The reverse error also counts: `LabelEncoder` is not `labelencoder`.

### A6 · Wrong part of speech `[M+J]` `1.2`, `1.4`, `9.2`
- Catch, verbed noun: a product or component name followed by a determiner —
  `rg -in '^\s*[-*0-9.)]*\s*(Slack|Jira|Docker|Terraform|S3|Kafka|diff|spec|ticket|action|architect|sunset|onboard|productioni[sz]e)\s+(the|a|an|this|these|all|your|it)\b'`
  (`Docker the service` matches, `Docker builds the image` does not).
- Catch, coined inflection: `Slacked`, `Kafka'd`, `dedeploy`, `unbackfill` (`[EXT]` — adapted
  from Part 2's no-antonyms note).
- Fix: a plain verb plus the name: "We `Kafka'd` the events" → "We published the events to
  Kafka".
- Carve-out: words that legitimately take both roles when position makes the role obvious
  (`9.2`): `cache`, `index`, `partition`, `filter`, `join`, `batch`, `log`, `query`, `train`,
  `fit`, `shard`, `deploy`, `merge`, `vacuum`. Keep the two roles in separate sentences.

### D10 · Modifier attachment ambiguity `[J]` `2.1`, `2.2` M2, `8.2`
Ambiguity is the defect. Length is only the proxy — a three-word cluster can still fail.
- Ask: bracket the cluster two ways. If both readings are plausible and differ, rewrite.
  `(model training) data` versus `model (training data)`.
- Fix: insert a preposition, or hyphenate the pair that binds. A hyphen asserts a relationship —
  a wrongly placed one states something false. Test: does the hyphenated compound name a thing
  that exists in the system?
