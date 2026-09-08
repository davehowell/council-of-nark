---
theme: default
title: The Experiments Failed. A Lot.
colorSchema: dark
fonts:
  sans: Arial
  local: Arial
transition: none
layout: default
class: cover
---

<p class="eyebrow">03 / The evidence</p>

# The Experiments<br>Failed. A Lot.

<p class="lead">What broke, what survived, and what changed.</p>

<div class="foot"><span>Council of Nark</span><span>1 / 6</span></div>

<!--
[Time: 0:00–0:40]
The most useful part of this project was not a leaderboard. It was finding the ways a leaderboard could mislead us. Some runs failed before a model answered. One partial run was discarded before scoring. Completed runs exposed both easy tasks and a fuser that could lose good findings. Then a focused repetition went against the persona I was testing. The story is a series of course corrections: each failure changed the measurement or narrowed the claim. I want to show three of those lessons, and the results that survived them, without turning this into a tour of every implementation detail.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/LAB_NOTEBOOK.md
-->

---

# First, make sure an answer is an answer

| Failure | Course correction |
|---|---|
| Schema rejected before inference | Check the adapter first |
| Prompt echo looked like output | Parse assistant events only |
| Duplicate complaints inflated noise | Cluster unsupported claims |

<p class="caution">The partial 71/81 run was sealed and excluded.</p>

<div class="foot"><span>Engineering failures, not model outcomes</span><span>2 / 6</span></div>

<!--
[Time: 0:40–1:30]
The first run produced no experimental observations because the client rejected a schema annotation. A later stream parser could mistake example JSON echoed from the prompt for a model answer, while provider errors could arrive with a successful process exit code. We preserved and excluded the partial run, repaired the parser and started again. Scoring also needed repair: repeated unsupported claims were being counted as separate false positives in a raw panel. That made fusion look more beneficial than it was. These were not evidence against the council. They were evidence that the experiment needed controls of its own. Historical scores keep their original limitations attached.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONTAMINATION_REVIEW.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-gemma-smoke-incomplete.md
-->

---

# Eight calls matched one review's score

<div class="pair">
<div><span class="number">0.860</span><span class="label">One functional review</span></div>
<div><span class="number">0.860</span><span class="label">Seven reviews + fusion</span></div>
</div>

<p class="verdict">8.9× the recorded cost in the Haiku smoke.</p>
<p class="caution">Three tasks, one draw. Equal F1 did not mean identical findings.</p>

<div class="foot"><span>Mean F1 · exploratory LLM scoring</span><span>3 / 6</span></div>

<!--
[Time: 1:30–2:25]
Here is the first result that makes a cheaper alternative interesting. One functional omnibus review and seven independent omnibus reviews followed by fusion both averaged F1 0.860. The repeated-review package cost about 8.9 times as much in recorded usage. But this is not an equivalence finding: there were only three tasks and one sample per condition, and sometimes the same score came from different defects. Specialists actually exposed higher raw recall, but fusion dropped valid findings. The baseline was already strong, leaving limited headroom. This justified testing the same material on a smaller model before spending more calls or making the synthetic tasks increasingly contrived.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-stage-a-smoke/README.md
https://github.com/davehowell/council-of-nark/blob/main/scripts/evidence_audit.py
-->

---

# The next smoke favoured specialists

| Gemma final review | Mean F1 |
|---|---:|
| One functional omnibus | 0.804 |
| Repeated omnibus + fusion | 0.863 |
| Functional specialists + fusion | **0.911** |
| Fictional specialists + fusion | 0.857 |

<p class="caution">One draw on three tasks. The ordering changed; proof did not arrive.</p>

<div class="foot"><span>81 successful calls · exploratory LLM scoring</span><span>4 / 6</span></div>

<!--
[Time: 2:25–3:20]
The clean Gemma smoke gave a different ordering. Functional specialists had the best final F1 in this draw. Their raw recall was 0.917 against 0.792 for repeated omnibus reviews, so there was evidence of additional coverage as well as additional noise. Fictional specialists did worse than functional specialists. The omnibus character comparison went the other way, which is another warning against interpreting one sample as a stable persona effect. We also cannot call the difference from Haiku a clean causal model effect: harness and scoring details changed between historical runs. The lesson was to repeat a tightly controlled pair before claiming that wording helped or hurt.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-stage-a-smoke-gemma/README.md
-->

---

# K-2SO lost the repeated comparison

<div class="pair">
<div><span class="number">0.790</span><span class="label">Functional correctness</span></div>
<div><span class="number warm">0.748</span><span class="label">Same kernel + K-2SO</span></div>
</div>

<p class="verdict">2 fictional wins · 16 ties · 12 functional wins</p>
<p class="caution">30 pairs, but only three tasks. One model and one persona.</p>

<div class="foot"><span>Mean F1 · exploratory LLM scoring</span><span>5 / 6</span></div>

<!--
[Time: 3:20–4:15]
This comparison held the correctness kernel constant and repeated the functional and fictional versions ten times on each packet. The fictional-minus-functional mean F1 difference was minus 0.0425. Its historical sampling-only interval was minus 0.0808 to minus 0.0047. That interval concerns repeated calls on these packets, not a population of software tasks. A new sensitivity check leaves the mean negative whichever packet we remove, so one packet is not solely driving the direction. There is also a useful correction to the cost story: the fictional outputs used about two percent fewer total tokens. The overlay was worse on detection here, but it was not more expensive by that token measure.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-persona-pair-gemma/README.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
-->

---

# Better isolation enabled the next test

- History-free source, with upstream answers removed.
- Bounded source access and offline regression tests.
- Gortex's parent fails; its evidence revision passes.

<p class="verdict">Zero ecological respondent calls so far.</p>

<div class="foot"><span>A valid test environment is not a council result</span><span>6 / 6</span></div>

<!--
[Time: 4:15–5:00]
The course correction was to move towards real pre-fix open-source tasks. That required stronger isolation: an empty respondent environment, history-free source, a narrow mediator for source access, and tests with network access denied. Gortex's frozen regression has been shown to fail before its fix and pass afterwards. Several failed exporter and process-lifecycle attempts were preserved on the way. This is valuable infrastructure evidence. It is not evidence that a reviewer found the bug: no ecological respondent call has occurred. That distinction is the bridge to the pivot talk. We now know what the existing results support and can define the next test without claiming it has already happened.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CHECKPOINT.md
-->
