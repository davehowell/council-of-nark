---
theme: default
title: What Would Settle It?
colorSchema: dark
fonts:
  sans: Arial
  local: Arial
transition: none
layout: default
class: cover
---

<p class="eyebrow">04 / The pivot</p>

# What Would Settle It?

<p class="lead">The next comparison, defined before the next call.</p>

<div class="foot"><span>Council of Nark</span><span>1 / 6</span></div>

<!--
[Time: 0:00–0:45]
Part 3 ends without a verdict because calibration and infrastructure checks do not answer the original question. This talk is the pivot, not the conclusion. It defines the next comparison before any new respondent call: practical review tasks, matched resource limits, blinded human ratings, and explicit decision rules. The council has not demonstrated a reliable advantage over one functional review, and equivalence is also untested. The next stage must collect the evidence needed for a later conclusion rather than make the current uncertainty sound final.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# Give one agent a fair chance

| Proposed arm | What it isolates |
|---|---|
| One functional review | Cheap baseline |
| One agent, self-review | Fewer-agent alternative |
| Three omnibus + fusion | Independent sampling |
| Three specialists + fusion | Functional roles |
| Same specialists + personas | Personality package |

<div class="foot"><span>Same model · matched total budgets for the four extended arms</span><span>2 / 6</span></div>

<!--
[Time: 0:45–1:40]
The key addition is a continuing agent that reviews its own answer twice and then consolidates it. Compare that with three independent omnibus reviews, three specialists and those same specialists with fictional overlays. All four extended arms get the same total resource cap; the cheap single review remains a practical reference. The reduced council uses correctness, architecture and operational risk, selected before outcomes are observed. This tests the three-persona package, not every character individually. Count provider turns as well as sessions because tool use can require several turns. Record actual tokens including repeated context and fusion: equal call counts alone do not make the comparison fair.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# Real projects make the test relevant

| Candidate project | Frozen task |
|---|---|
| Manim | Animation ordering regression |
| dlt | Paginator stopping precedence |
| Gortex | Unicode tokenizer panic |
| turbovec | Stale temporary-file cleanup |

<p class="caution">Candidates are not outcomes. Each needs a validated test environment.</p>

<div class="foot"><span>Open-source fixes are evidence, not the only acceptable remedies</span><span>3 / 6</span></div>

<!--
[Time: 1:40–2:30]
The repository already curates tasks from these four projects. They are useful because there is a pre-fix revision and external evidence of what eventually corrected the symptom. Respondents receive the earlier source and a non-leading brief, with history and upstream answers removed. Only Gortex currently has the validated infrastructure described in the previous talk. The others still need their frozen dependency closures and before-and-after test checks. We must allow supported alternative remedies rather than reward copying the upstream patch. These are review and diagnosis tasks; passing infrastructure tests does not turn this into a benchmark of agents implementing software. Four projects improve realism but remain a very small population sample.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/candidates.json
https://github.com/davehowell/council-of-nark/blob/main/experiment/CHECKPOINT.md
-->

---

# Define what would change the decision

- Better: a useful quality gain with acceptable failure risk.
- Equivalent: uncertainty fits inside a declared small margin.
- Cheaper substitute: quality holds up and resources fall.

<p class="verdict">Apply the rule after the confirmatory run.</p>

<div class="foot"><span>Two independent humans · locked rubric · declared comparisons</span><span>4 / 6</span></div>

<!--
[Time: 2:30–3:25]
The real-code rubric measures localisation, consequence, correction, regression tests and actionability. The proposed primary utility averages those anchored dimensions and gives zero to a missing answer, an unsupported primary mechanism or a critically harmful remedy. We preserve the individual scores and unsupported claims alongside it. A proposed practical margin is five percentage points on the rescaled utility, which needs agreement before the run. Four primary comparisons need simultaneous uncertainty control, not four uncorrected chances to find a win. A cheaper substitute must satisfy the quality rule and demonstrate resource savings. No critical failures in a tiny sample is not proof of equal safety, so those rates need uncertainty too.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/SCORING.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# Run feasibility, then confirmation

<div class="pair">
<div><span class="number">4 × 2</span><span class="label">Tasks × repeats for feasibility</span></div>
<div><span class="number warm">57*</span><span class="label">Illustrative independent-task target</span></div>
</div>

<p class="caution">*Freeze the final count from the declared precision rule.</p>

<div class="foot"><span>Exclude feasibility data from the confirmatory result</span><span>5 / 6</span></div>

<!--
[Time: 3:25–4:15]
Start with the four curated tasks and two repeats to test the workflow, rubric, and resource accounting. Exclude those feasibility observations from the result. Then freeze a larger set of independent eligible tasks and run the precision-sized confirmation. The current calculation gives an illustrative target of about fifty-seven tasks under one variability assumption; repository clustering can increase it. Repeated calls on four tasks cannot replace task diversity. The final count, spending cap, exclusion rules, and analysis must be fixed before the confirmatory responses. Part 5 should report that run, including failures and uncertainty, rather than promote the feasibility pilot into an answer.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# Next: run the ecological comparison

<p class="lead">Finish the runner.<br>Freeze the tasks.<br>Collect blinded human ratings.</p>

<p class="verdict">Then publish Part 5: the result.</p>

<div class="foot"><span>Part 4 is the pivot, not the conclusion</span><span>6 / 6</span></div>

<!--
[Time: 4:15–5:00]
This talk hands the project back to execution. Finish the claim runner and deterministic mock lifecycle. Validate enough history-free task environments for the declared precision target. Freeze the compared arms, budgets, exclusions, and analysis. Then collect condition-blinded ratings from two independent humans. The existing K-2SO result and mixed council calibrations remain part of the record, but they are not the ending. Part 5 can be the conclusion after the ecological comparison produces results. It must report the outcome even if the council loses, the cheaper alternative holds up, or the interval still crosses the declared decision boundary.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
-->
