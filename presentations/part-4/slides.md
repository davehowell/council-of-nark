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

<p class="eyebrow">04 / The conclusion</p>

# What Would Settle It?

<p class="lead">A smaller final test. An honest stopping rule.</p>

<div class="foot"><span>Council of Nark</span><span>1 / 6</span></div>

<!--
[Time: 0:00–0:45]
We can close the current evidence review now, but we cannot honestly close every empirical question with a definitive winner. The council has not demonstrated a reliable advantage over one functional review. Nor have we established equivalence. One persona underperformed in one controlled calibration. The question is what additional work would change that confidence enough to be useful. My proposal is a smaller final comparison focused on practical review quality and resource use. It deliberately leaves provider diversity, every chain order and the full eight-role factorial for another study. This talk separates that proposed work from the results we already have.

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

<p class="verdict">Otherwise, publish “inconclusive”.</p>

<div class="foot"><span>Two independent humans · locked rubric · declared comparisons</span><span>4 / 6</span></div>

<!--
[Time: 2:30–3:25]
The real-code rubric measures localisation, consequence, correction, regression tests and actionability. The proposed primary utility averages those anchored dimensions and gives zero to a missing answer, an unsupported primary mechanism or a critically harmful remedy. We preserve the individual scores and unsupported claims alongside it. A proposed practical margin is five percentage points on the rescaled utility, which needs agreement before the run. Four primary comparisons need simultaneous uncertainty control, not four uncorrected chances to find a win. A cheaper substitute must satisfy the quality rule and demonstrate resource savings. No critical failures in a tiny sample is not proof of equal safety, so those rates need uncertainty too.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/SCORING.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# A small pilot cannot prove a small gap

<div class="pair">
<div><span class="number">4 × 2</span><span class="label">Tasks × repeats for feasibility</span></div>
<div><span class="number warm">136</span><span class="label">Planned stages, not API calls</span></div>
</div>

<p class="caution">Confirmation needs a sample sized for precision, then one final analysis.</p>

<div class="foot"><span>Proposed · not run · extra feasibility, tools and retries excluded</span><span>5 / 6</span></div>

<!--
[Time: 3:25–4:15]
Four tasks with two repeats across these arms require 136 planned review stages, before the separate initial feasibility block, tool-driven provider turns or retries. That can test whether the workflow completes and whether the rubric and budgets make sense. It cannot establish narrow equivalence. An illustrative precision calculation in the protocol can require around fifty-seven independent tasks under one assumed variability level, and repository clustering can require more. That is not a sample-size promise; it shows why adding repeats to four tasks is not a shortcut to generality. After excluded feasibility, freeze the confirmatory sample and spending cap. If the available resources cannot buy the required precision, stop and report the limit.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->

---

# The council has not earned a default

<p class="lead">Start with one functional review.<br>Escalate when the task warrants it.</p>

<p class="verdict">A practical choice under uncertainty.<br>Not a claim of proven equivalence.</p>

<div class="foot"><span>Keep the negative results. Keep the route to replication.</span><span>6 / 6</span></div>

<!--
[Time: 4:15–5:00]
My present operational choice is to start with one functional review and escalate for tasks whose risk justifies more attention. That is a decision under uncertainty, not a theorem that one agent is equivalent. The strongest persona result remains narrow: K-2SO reduced detection on this model and these three packets, with slightly fewer tokens. The council results remain mixed and exploratory. The closing protocol gives us a way to earn a stronger statement, including an equivalence result if the evidence is precise enough. Until the launcher, task environments and two human raters are ready, it remains a proposal. The useful outcome is a claim we can defend, even when that claim is smaller than the original idea.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
-->
