# The Council of Nark — speaker notes

Ten slides · approximately eight minutes · 981 spoken words.

Edit the closing comment in each slide in part-1/slides.md; these notes are generated from that source.

## 1. The Council of Nark

[Time: 0:00–0:50]
The Council of Nark began with a workflow I was already using. Fable 5 could ask several agents to inspect the same task and combine their replies. The run was slower and used more tokens than one call. While it ran, I often stepped away from the screen and planned across projects on A3 paper. I also used Gemini for second opinions and an HK-47 reviewer for unnecessary complexity. Those tools did not prove that a panel was better. They gave me a practical question: will reviewers with different, stable responsibilities find useful problems that one general review misses? The diagrams in this talk show that original hypothesis, not measured results.

[Sources]
https://davehowell.github.io/experimental/
https://github.com/davehowell/council-of-nark/blob/main/README.md

## 2. The workflow already mixed reviewers

[Time: 0:50–1:40]
The workflow already mixed models and instructions. Gemini gave me an independent answer from another provider. HK-47 used the same underlying model as other agents but had one narrow job: challenge over-engineering. Fable showed a convenient way to run several reviews and collect the output. The council made those parts explicit. Each reviewer would have a stable checklist. The controller would choose only the reviewers relevant to the task. They would inspect the same source independently, so one review would not steer the next. Repeated findings might strengthen confidence, while a supported finding from one reviewer might expose a blind spot. Better review quality was still only a hypothesis.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/tree/main/agents

## 3. Meet the council

[Time: 1:40–2:05]
The characters made the responsibilities easier to remember and made the project enjoyable to use. That was their known purpose. Whether character wording improved a model's review was a separate question. The functional instructions had to stand on their own. A reviewer needed a defined scope, evidence for each finding and a common output format. The names could remain useful labels even if an experiment found no performance benefit from the fictional wording.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/README.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/PERSONA_FACTORIAL.md

## 4. Who reviews what?

[Time: 2:05–3:05]
The roster covers seven kinds of risk. HK-47 asks whether the solution contains machinery that the problem does not require. K-2SO checks behaviour, tests and failure signals. GLaDOS follows boundaries between components. C-3PO checks sensitive data, permissions and destructive operations. Bender examines compute and cost. Holly takes the long view: manual steps, forgotten toggles, state drift and knowledge that can disappear when people leave. WALL-E checks whether technical writing will remain clear to its next reader. The controller does not run all seven by default. It chooses the smallest group that matches the code, plan or document, with at least two independent reviewers.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/README.md

## 5. Holly and WALL-E work at different points

[Time: 3:05–3:55]
Holly and WALL-E widened the design in different ways. Holly sits on the review panel when a change depends on memory, manual work or long-term operational care. WALL-E can review prose as a specialist, but also has a separate job after the arbiter finishes. The arbiter first merges duplicate findings, resolves conflicts and ranks the result. WALL-E then translates that settled review for people who do not need file-and-line details. Translation does not create another vote and does not count as a panel review. Keeping those jobs separate prevents a simpler explanation from quietly changing the technical decision.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/skills/walle-ste/SKILL.md

## 6. One review can produce different outcomes

[Time: 3:55–4:45]
This cone was a way to draw variation, not a statistical model fitted to observations. Start with the same code and prompt. A model can produce a strong review, an average one or a weak one. More reasoning does not guarantee a better answer, and a later run can miss something that an earlier run found. Better context and clearer instructions can improve the starting conditions, but they do not make every output identical. The diagram gave me language for the next question: what happens when several model calls pass work from one to another?

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md

## 7. A serial chain can carry mistakes forward

[Time: 4:45–5:35]
In a serial chain, the second reviewer does not begin independently. It inherits the first answer. A strong first answer can still lose a valid finding during the next rewrite. A weak first answer can lead the next reviewer towards the wrong concern. The second reviewer can also correct the first, so the diagram does not claim that chains always get worse. It shows why the order matters and why a chain needs its own test. The protocol later compared every ordering of a three-reviewer chain instead of choosing one favourable sequence.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md

## 8. One arbiter combines independent reviews

[Time: 5:35–6:25]
The alternative used parallel reviews and one combining step. Each specialist receives the same source and works independently. The arbiter sees all of the findings together. It can merge duplicate reports, reject unsupported claims and preserve a useful finding that only one specialist noticed. My original expectation was that this process would raise the typical quality and reduce the weak end of the range. That expectation could be wrong. The combining step can also discard a correct finding or make repeated errors sound authoritative. The experiment therefore had to score the individual findings and the combined review separately.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md

## 9. Flower, not a chain

[Time: 6:25–7:20]
These diagrams show the structural difference. On the left, reviewers work from the same source and return findings to one arbiter. The arbiter can compare disagreements because no reviewer has rewritten another review first. On the right, output passes through a sequence. Later reviewers can see and change earlier work, so the path through the reviewers becomes part of the result. I called the first shape a flower because all work returns to one centre. The useful claim was not that flowers are inherently better. It was that independent review followed by one combining step should be compared with serial review under controlled conditions.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md

## 10. A testable proposal, not a conclusion

[Time: 7:20–8:00]
The original idea produced three claims that could be tested separately. First, specialist roles might find different real defects from repeated general reviews. Second, fictional character wording might help, hinder or make no useful difference when the functional instructions stay fixed. Third, independent fan-out and one fusion step might produce a better final review than a chain. Cost, tokens and time belong beside quality in every comparison. The desired outcome was never the largest council. It was the smallest review process that produced a dependable result. The next talk explains how the protocol tried to separate those claims.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

---

# Put the Council on Trial — speaker notes

Six slides · approximately five minutes · 639 spoken words.

Edit the closing comment in each slide in part-2/slides.md; these notes are generated from that source.

## 1. Put the Council on Trial

[Time: 0:00–0:45]
Suppose seven reviewers find more defects than one. That sounds like a council victory, but perhaps seven ordinary reviews would do just as well. Or perhaps one reviewer with the same budget would outperform both. This is the core design problem: extra computation, specialist instructions, character prose and fusion all change the output. If we change them together, we learn whether a package happened to work, but not which part earned its place. The study therefore began by separating the claims. A fair test needs controls that can make the exciting explanation unnecessary.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md

## 2. Change one thing at a time

[Time: 0:45–1:35]
Here are four comparisons we can understand without learning every arm code. First, does repeating an omnibus review add value? Second, at the same number of calls, do specialists add more than that repetition? Third, with the specialist instructions unchanged, does fictional wording help? Fourth, what does the fuser retain or discard? The original council arms each used seven reviewers and one fuser. Matching those calls is useful for isolating roles, but it does not guarantee equal token spend: different prompts can produce different lengths. Provider diversity and informed chains were additional planned questions. They were not completed results and need not all be solved to close this study.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md

## 3. A score needs an answer key

[Time: 1:35–2:25]
We started with three synthetic packets: a revenue dashboard, key rotation and a webhook redesign. Each had planted defects and clean facts, so we could check whether a review found a real issue or complained about something sound. That gives us a controlled calibration set. It also creates a limitation: doing ten runs on each packet still leaves only three different tasks. We can estimate how much outputs fluctuate on those packets, but cannot pretend we have sampled thirty independent software problems. The answer keys are public for reproducibility and excluded from respondent prompts. Public availability means this is a reproducible demonstration, not a permanently secret benchmark.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/README.md

## 4. Equal scores can hide different reviews

[Time: 2:25–3:20]
A true positive is a unique planted defect correctly identified. A false positive is an unsupported claim, and a false negative is a missed planted defect. Precision is true positives divided by all positive claims; recall is true positives divided by all planted defects. F1 balances those two, and the macro average gives each task equal weight. Those definitions matter because identical F1 can conceal different issue identities. A review can also diagnose a problem correctly and recommend a bad fix. We therefore need overlap and fusion retention, and the real-code stage needs an anchored remedy rubric. Tokens, cost and latency stay beside quality rather than disappearing into a flattering composite score.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/SCORING.md

## 5. The measurement needs its own controls

[Time: 3:20–4:10]
The controls are practical. Freeze the prompts, packets, model settings and analysis before a claim-bearing run. Save raw requests and responses, then seal them by digest so later analysis cannot quietly change the observations. Keep source evidence and answer keys out of the respondent environment. Have two humans rate independently before the condition labels are revealed. Character language can still reveal its treatment, so the raters also record their guesses; we should not claim perfect blinding. The published runs used LLM triage, which is useful for finding plumbing problems but falls short of that confirmatory standard. Our confidence has to reflect how the scores were produced.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/RUNSHEET.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/ISOLATION.md

## 6. “No clear winner” is not equivalence

[Time: 4:10–5:00]
The final distinction is the one most likely to get lost in a short presentation. An uncertain difference is not evidence that two approaches are equivalent. To claim a cheaper substitute, we need to decide beforehand how much quality loss would matter, then collect enough evidence to rule that loss out. To call approaches practically equivalent, the whole interval must fit inside our agreed margin. A wide interval that crosses zero answers neither question. This means a small pilot can honestly conclude that the setup works, or that a direction deserves another test, without claiming the underlying question is settled. That is exactly the discipline the first runs turned out to need.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

---

# The Experiments Failed. A Lot. — speaker notes

Six slides · approximately five minutes · 635 spoken words.

Edit the closing comment in each slide in part-3/slides.md; these notes are generated from that source.

## 1. The Experiments Failed. A Lot.

[Time: 0:00–0:40]
The most useful part of this project was not a leaderboard. It was finding the ways a leaderboard could mislead us. Some runs failed before a model answered. One partial run was discarded before scoring. Completed runs exposed both easy tasks and a fuser that could lose good findings. Then a focused repetition went against the persona I was testing. The story is a series of course corrections: each failure changed the measurement or narrowed the claim. I want to show three of those lessons, and the results that survived them, without turning this into a tour of every implementation detail.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/LAB_NOTEBOOK.md

## 2. First, make sure an answer is an answer

[Time: 0:40–1:30]
The first run produced no experimental observations because the client rejected a schema annotation. A later stream parser could mistake example JSON echoed from the prompt for a model answer, while provider errors could arrive with a successful process exit code. We preserved and excluded the partial run, repaired the parser and started again. Scoring also needed repair: repeated unsupported claims were being counted as separate false positives in a raw panel. That made fusion look more beneficial than it was. These were not evidence against the council. They were evidence that the experiment needed controls of its own. Historical scores keep their original limitations attached.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONTAMINATION_REVIEW.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-gemma-smoke-incomplete.md

## 3. Eight calls matched one review's score

[Time: 1:30–2:25]
Here is the first result that makes a cheaper alternative interesting. One functional omnibus review and seven independent omnibus reviews followed by fusion both averaged F1 0.860. The repeated-review package cost about 8.9 times as much in recorded usage. But this is not an equivalence finding: there were only three tasks and one sample per condition, and sometimes the same score came from different defects. Specialists actually exposed higher raw recall, but fusion dropped valid findings. The baseline was already strong, leaving limited headroom. This justified testing the same material on a smaller model before spending more calls or making the synthetic tasks increasingly contrived.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-stage-a-smoke/README.md
https://github.com/davehowell/council-of-nark/blob/main/scripts/evidence_audit.py

## 4. The next smoke favoured specialists

[Time: 2:25–3:20]
The clean Gemma smoke gave a different ordering. Functional specialists had the best final F1 in this draw. Their raw recall was 0.917 against 0.792 for repeated omnibus reviews, so there was evidence of additional coverage as well as additional noise. Fictional specialists did worse than functional specialists. The omnibus character comparison went the other way, which is another warning against interpreting one sample as a stable persona effect. We also cannot call the difference from Haiku a clean causal model effect: harness and scoring details changed between historical runs. The lesson was to repeat a tightly controlled pair before claiming that wording helped or hurt.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-stage-a-smoke-gemma/README.md

## 5. K-2SO lost the repeated comparison

[Time: 3:20–4:15]
This comparison held the correctness kernel constant and repeated the functional and fictional versions ten times on each packet. The fictional-minus-functional mean F1 difference was minus 0.0425. Its historical sampling-only interval was minus 0.0808 to minus 0.0047. That interval concerns repeated calls on these packets, not a population of software tasks. A new sensitivity check leaves the mean negative whichever packet we remove, so one packet is not solely driving the direction. There is also a useful correction to the cost story: the fictional outputs used about two percent fewer total tokens. The overlay was worse on detection here, but it was not more expensive by that token measure.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/results/2026-08-20-persona-pair-gemma/README.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md

## 6. Better isolation enabled the next test

[Time: 4:15–5:00]
The course correction was to move towards real pre-fix open-source tasks. That required stronger isolation: an empty respondent environment, history-free source, a narrow mediator for source access, and tests with network access denied. Gortex's frozen regression has been shown to fail before its fix and pass afterwards. Several failed exporter and process-lifecycle attempts were preserved on the way. This is valuable infrastructure evidence. It is not evidence that a reviewer found the bug: no ecological respondent call has occurred. That distinction is the bridge to the pivot talk. We now know what the existing results support and can define the next test without claiming it has already happened.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CHECKPOINT.md

---

# What Would Settle It? — speaker notes

Six slides · approximately five minutes · 609 spoken words.

Edit the closing comment in each slide in part-4/slides.md; these notes are generated from that source.

## 1. What Would Settle It?

[Time: 0:00–0:45]
Part 3 ends without a verdict because calibration and infrastructure checks do not answer the original question. This talk is the pivot, not the conclusion. It defines the next comparison before any new respondent call: practical review tasks, matched resource limits, blinded human ratings, and explicit decision rules. The council has not demonstrated a reliable advantage over one functional review, and equivalence is also untested. The next stage must collect the evidence needed for a later conclusion rather than make the current uncertainty sound final.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 2. Give one agent a fair chance

[Time: 0:45–1:40]
The key addition is a continuing agent that reviews its own answer twice and then consolidates it. Compare that with three independent omnibus reviews, three specialists and those same specialists with fictional overlays. All four extended arms get the same total resource cap; the cheap single review remains a practical reference. The reduced council uses correctness, architecture and operational risk, selected before outcomes are observed. This tests the three-persona package, not every character individually. Count provider turns as well as sessions because tool use can require several turns. Record actual tokens including repeated context and fusion: equal call counts alone do not make the comparison fair.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 3. Real projects make the test relevant

[Time: 1:40–2:30]
The repository already curates tasks from these four projects. They are useful because there is a pre-fix revision and external evidence of what eventually corrected the symptom. Respondents receive the earlier source and a non-leading brief, with history and upstream answers removed. Only Gortex currently has the validated infrastructure described in the previous talk. The others still need their frozen dependency closures and before-and-after test checks. We must allow supported alternative remedies rather than reward copying the upstream patch. These are review and diagnosis tasks; passing infrastructure tests does not turn this into a benchmark of agents implementing software. Four projects improve realism but remain a very small population sample.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/candidates.json
https://github.com/davehowell/council-of-nark/blob/main/experiment/CHECKPOINT.md

## 4. Define what would change the decision

[Time: 2:30–3:25]
The real-code rubric measures localisation, consequence, correction, regression tests and actionability. The proposed primary utility averages those anchored dimensions and gives zero to a missing answer, an unsupported primary mechanism or a critically harmful remedy. We preserve the individual scores and unsupported claims alongside it. A proposed practical margin is five percentage points on the rescaled utility, which needs agreement before the run. Four primary comparisons need simultaneous uncertainty control, not four uncorrected chances to find a win. A cheaper substitute must satisfy the quality rule and demonstrate resource savings. No critical failures in a tiny sample is not proof of equal safety, so those rates need uncertainty too.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/ecological/SCORING.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 5. Run feasibility, then confirmation

[Time: 3:25–4:15]
Start with the four curated tasks and two repeats to test the workflow, rubric, and resource accounting. Exclude those feasibility observations from the result. Then freeze a larger set of independent eligible tasks and run the precision-sized confirmation. The current calculation gives an illustrative target of about fifty-seven tasks under one variability assumption; repository clustering can increase it. Repeated calls on four tasks cannot replace task diversity. The final count, spending cap, exclusion rules, and analysis must be fixed before the confirmatory responses. Part 5 should report that run, including failures and uncertainty, rather than promote the feasibility pilot into an answer.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 6. Next: run the ecological comparison

[Time: 4:15–5:00]
This talk hands the project back to execution. Finish the claim runner and deterministic mock lifecycle. Validate enough history-free task environments for the declared precision target. Freeze the compared arms, budgets, exclusions, and analysis. Then collect condition-blinded ratings from two independent humans. The existing K-2SO result and mixed council calibrations remain part of the record, but they are not the ending. Part 5 can be the conclusion after the ecological comparison produces results. It must report the outcome even if the council loses, the cheaper alternative holds up, or the interval still crosses the declared decision boundary.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
