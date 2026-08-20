# Council of Nark — empirical experiment design

> Status: protocol draft. Run a smoke test before preregistration, freeze the prompts and answer keys, then run the confirmatory trial. The aim is to falsify the claims, not illustrate them.

## 1. Claims and limits

The deck currently makes three separable claims:

1. **Role diversity:** one model prompted as independent specialists finds more real defects than the same model making repeated general reviews.
2. **Persona wrapper:** pop-culture characterisation adds value beyond the specialist's explicit checklist, rather than spending tokens on style.
3. **Topology:** independent fan-out followed by one fusion retains more signal than a serial review chain.

The experiment measures repeatable changes in the conditional output: defect coverage, errors, variance, tokens, and cost. Role cues may change hidden-state and attention trajectories, and which learned features contribute through fixed weights; this changes logits and token probabilities; decoding then adds sample variance. Black-box APIs cannot isolate those stages or show that a prompt “activates a different model weight space.”

### Falsifiable predictions

- At matched calls, a functional specialist panel beats repeated omnibus reviewers on planted-defect F1.
- Adding the character wrapper to an otherwise identical specialist prompt either improves F1 or does not. More jokes alone do not count.
- Fan-out has higher final F1, lower run-to-run variance, and less loss of previously found defects than a chain.
- If fusion helps, it raises precision without discarding enough true findings to lower F1.
- Any gain should appear in the p10/floor and variance as well as the mean. A mean-only claim is weaker than the current deck's "narrower cone" claim.

## 2. Toy review packets

Use three small, frozen packets. Each has a private answer key with planted defect IDs. The packets cover different mixes so a single lucky checklist cannot win the whole study.

| Packet | Material | Planted lenses | Files |
|---|---|---|---|
| 01 — Revenue dashboard | dbt/SQL, semantic layer, design note, recovery steps | correctness, cost, exposure, architecture, simplicity, entropy, prose | `experiment/scenarios/01-revenue-dashboard/` |
| 02 — Key rotation at 03:00 | incident runbook and deployment note | safety, compliance, observability, entropy, prose, needless machinery | `experiment/scenarios/02-key-rotation/` |
| 03 — Webhook redesign | short architecture proposal | idempotency, contracts, PII, cost, observability, scale, manual-state decay | `experiment/scenarios/03-webhook-redesign/` |

The model sees only `review-packet.md`. The scorer sees `answer-key.md`. Do not put persona names, expected lens labels, or the experiment hypothesis in a packet.

Three packets are enough for a mechanism pilot, not a general claim about software review. A confirmatory study should create 10 semantics-preserving mutations of each packet: rename entities, vary defect locations, swap equivalent syntax, and include clean controls. Freeze those variants before the confirmatory run.

## 3. Stage A — same model, different prompting

Pin one cheap model snapshot, for example the available Claude Haiku snapshot. Use the same model for every reviewer and the fuser. Start each call in a fresh session with no tools or memory.

Seven panel lenses are used. The fuser is a separate, common persona-free call so GLaDOS can be tested as the architecture specialist:

- simplicity / over-engineering
- correctness / testing / observability
- architecture / contracts / coupling
- security / compliance
- FinOps / compute waste
- long-horizon operational entropy
- technical-language clarity

### Arms

| ID | Calls | Review condition | Purpose |
|---|---:|---|---|
| S0 | 1 | Plain generic adversarial review | Weak practical baseline |
| S1 | 1 | Functional omnibus prompt that explicitly lists all seven lenses, with no fictional persona | Fair single-reviewer baseline |
| S2 | 1 | The S1 prompt plus a length-matched GLaDOS character wrapper | Single reviewer: persona vs no persona |
| M0 | 7 + 1 | Seven independent copies of the functional omnibus S1, then common fuser | Matched-call sampling control |
| M1 | 7 + 1 | Seven functional specialists, no fictional characterisation, then common fuser | Value of role specialisation |
| M2 | 7 + 1 | The same seven specialist kernels plus character wrappers, then common fuser | Incremental value or cost of pop culture |

The critical comparisons are:

- **specialisation:** M1 − M0
- **pop-culture persona:** M2 − M1
- **single-call persona:** S2 − S1
- **more samples:** M0 − S1
- **practical council vs functional omnibus reviewer:** M2 − S1, reported with its extra cost

S0 is useful context but is not the fair control. Do not claim victory from M2 beating S0.

### Prompt control

Each specialist pair must have an identical lens kernel, task text, finding schema, output cap, and severity rules. M2 adds only a short character/style wrapper. Keep the functional and fictional wrappers within ±10% input tokens. Do not use the full production persona files in the controlled trial; their unequal length and richer checklists would confound persona with instructions. Run those later as an ecological-validity replication.

Use the same finding contract in all arms:

```text
Return at most 8 findings. Each finding must contain:
- location
- claim
- consequence
- fix
- confidence: high | medium | low
Do not rewrite the artifact. Do not invent a finding to fill the quota.
```

The omnibus prompt must name all seven perspectives in comparable detail. This gives the single reviewer a genuine chance to organise its attention.

### Raw panel versus fusion

Score both:

1. the union of the seven raw finding sets; and
2. the fuser's final set.

This reveals fusion gain and **fusion loss**. The fuser helps when it removes false positives and duplicates while retaining true positives. If it only drops valid minority findings, the arbiter claim fails even when the raw ensemble is strong.

## 4. Stage B — does the effect cross providers?

After Stage A is frozen, send the exact same UTF-8 prompt text and packets to one pinned cheap model from each provider:

- Claude Haiku
- Gemini Flash or Flash-Lite
- the available low-cost GPT model

First replicate one specialist pair: the persona-free correctness reviewer and its K-2SO wrapper. This estimates `provider × persona` interaction cheaply. If an effect survives, replicate S1, M1, and M2 across providers.

"Same prompt" means byte-identical text, not equal tokenisation; provider tokenisers and hidden system prompts differ. Pin model IDs and dates, record API parameters, and treat provider as a factor rather than proof that one provider is intrinsically better.

## 5. Stage C — fan-out versus chain

Use three roles so every chain order is tractable:

- HK-47 / simplicity
- K-2SO / correctness and observability
- C-3PO / security and compliance

There are `3! = 6` chain orders. Test all six.

### Primary topology comparison

Both topologies use three reviewer calls plus the same final fuser:

**Fan-out**

```text
controller ─┬─ reviewer A(original packet) ─┐
            ├─ reviewer B(original packet) ─┼─ same fuser ─ verdict
            └─ reviewer C(original packet) ─┘
```

**Informed chain**

```text
reviewer A(original packet)
  → reviewer B(original packet + A's finding ledger)
  → reviewer C(original packet + A/B ledger)
  → same fuser(A, B, C outputs)
```

The informed chain is the fair primary test. Every reviewer can still inspect the source. A relay where later agents see only the prior summary is a useful stress test, but it is too lossy to be the main control.

Use the same role prompts and output caps. Cap the inherited ledger at a fixed size and record all input tokens because the chain consumes more context. For each chain order, give the fan-out fuser the independent outputs in that same order; this controls for fuser recency and input ordering.

### Chain metrics

At each hop record:

- true findings added
- true findings retained from the prior ledger
- true findings dropped or materially weakened (**regression**)
- false findings added and corrected
- final F1

Report the final score for every order and the range across all six. If one favourable order wins and the others fail, the chain is order-sensitive rather than generally superior.

## 6. Blinding and run procedure

1. Give each packet, arm, repeat, and chain order an opaque run ID.
2. Generate all reviewer calls in fresh sessions. The model never sees arm names, other conditions, answer keys, or hypotheses.
3. Shuffle finding order before scoring. Remove character quotes only from the scorer copy; retain the raw output for token/style analysis.
4. Give the scorer the answer key and finding, but no arm, model, provider, prompt, or run metadata.
5. Use a deterministic matcher for exact locations and defect IDs where possible. Send ambiguous semantic matches to two blinded human raters. A strong LLM judge may triage, but human agreement decides the final label.
6. Freeze prompts, caps, decoding settings, model snapshot IDs, mutation variants, exclusion rules, and the analysis script before the confirmatory run.
7. Save raw requests, raw responses, latency, input/output tokens, errors, retries, and cost.

Retries are infrastructure events, not new samples. Predefine whether a refused or malformed response scores zero; the recommended rule is zero after one format-repair attempt that cannot alter substantive content.

## 7. Scoring

Map each finding to zero or one planted defect ID. A finding that maps to no defect enters a blinded semantic false-positive cluster. Two phrasings of the same true or false claim count once.

### Primary metric

- **Macro F1 over planted defects**, calculated per packet/run and averaged with each packet weighted equally.
- For role specialisation (`M1 − M0`), compare the semantically de-duplicated raw unions first. This measures what the panel found before a capable fuser can flatten differences.
- For practical shipped output, compare fused verdicts and report fusion retention separately.

### Secondary metrics

- recall and precision
- severity-weighted recall
- p10, worst case, variance, and interquartile range across repeats
- true findings per 1,000 output tokens and per dollar
- fusion retention: fused true positives / raw-union true positives
- false-positive removal by fusion
- malformed-output rate and latency

### Diversity and overlap

Compare semantic defect IDs, not wording:

- pairwise Jaccard overlap between reviewers
- unique true contribution from each lens
- duplicate true findings (useful consensus)
- unique false findings (noise, not diversity)
- coverage entropy across defect categories

A persona condition that produces more varied prose but no additional true defect IDs has not produced useful diversity.

## 8. Sample sizes

### Smoke test

Run each arm once on each packet. Its purpose is to catch broken prompts, ceiling/floor effects, malformed schemas, and answer-key ambiguity. Do not report inferential statistics.

### Toy pilot

Run 10 independent samples per `(packet × arm)` in Stage A. For topology, run 3 samples per `(packet × order × topology)`. Show bootstrap confidence intervals, but label the study a pilot because there are only three base packets.

### Confirmatory extension

Use 10 frozen variants per packet and at least 5 independent samples per variant/arm. Bootstrap by packet variant, not by individual finding. A mixed-effects logistic model can analyse defect detection with arm and scenario as fixed effects and defect/variant as random effects.

Run a power simulation from the pilot before choosing the final repeat count. Thirty to fifty repeats of only the same three packets would estimate sampling variance but would not fix the lack of task diversity.

## 9. Fairness and threats to validity

- **Compute:** M arms spend eight calls. M0 is the matched-call control; also report tokens, latency, and cost-normalised scores.
- **Prompt content:** explicit lens kernels, fictional wrappers, and provider identity are separate factors. Do not compare the full production Bender prompt with a two-line generic prompt and call the difference "persona".
- **Context cost:** character wrappers and chain ledgers consume tokens. Measure useful findings per token and truncation rate.
- **Ceiling/floor:** calibrate with a smoke test. Add clean packets so reviewers can demonstrate restraint.
- **Grader bias:** prefer planted defects and locations. Blind human raters to condition and report Cohen's kappa.
- **Correlated samples:** repeated calls on one packet are not independent tasks. Treat packet variants as the unit for generalisation.
- **Model drift:** pin snapshots where possible and finish a randomised block close in time.
- **Contamination:** use synthetic private packets, not common benchmark tasks.
- **Style leakage:** the scorer must not infer the arm from catchphrases; strip character-only prose from scorer copies.
- **Ecological validity:** toy defects test mechanism. A later trial on real PRs should use developer acceptance and defects found before merge.

## 10. Preregistered interpretation

Use effect sizes and 95% bootstrap confidence intervals. Replace the thresholds below after the smoke test, before the confirmatory run.

- **Role diversity supported:** M1's de-duplicated raw union beats M0 on macro F1, recall, and p10 with no material precision collapse. Fused output is the practical secondary check.
- **Pop-culture wrapper supported:** M2 beats M1 by the preregistered smallest effect of interest and improves true findings per token. Characterful wording alone is failure.
- **Pop-culture wrapper wasteful:** M2 is equivalent or worse than M1 and uses more tokens or produces more malformed/noisy output.
- **Fusion supported:** fused F1 exceeds raw-union F1 through enough false-positive removal to offset fusion loss.
- **Fusion weakened:** fused verdict approximately matches the raw union; the value came from sampling/specialisation.
- **Fan-out supported:** fan-out beats the average chain order and its p10, while chains show higher regression and order sensitivity.
- **Topology claim refuted:** informed chains match or beat fan-out across orders at comparable cost.
- **Weight-space mechanism remains unproven:** any behavioural win warrants mechanistic follow-up, not an internal-weights claim.

## 11. Deliverables

- frozen packets and private answer keys under `experiment/scenarios/`
- versioned prompt templates and request logs
- one tidy row per finding plus one summary row per run
- plots of F1, p10, variance, overlap, unique true contribution, fusion loss, token cost, and chain-order sensitivity
- the follow-up deck: `presentations/part-2/slides.md`
- an honest write-up that reports null and negative results
