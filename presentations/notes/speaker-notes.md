# The Council of Nark — speaker notes

Six slides · approximately five minutes · 616 spoken words.

Edit the closing comment in each slide in part-1/slides.md; these notes are generated from that source.

## 1. The Council of Nark

[Time: 0:00–0:45]
I started with a familiar feeling: a second pair of eyes catches things I miss. With agents, we can ask for several pairs almost instantly. That makes a council attractive, especially when each reviewer has a different job. But another review also costs context, tokens and time, and someone has to reconcile the answers. The question is whether the extra eyes earn that cost. I wanted a test that could tell me to use fewer agents just as comfortably as it could tell me to use more. This series is the story of trying to make that answer believable.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/README.md

## 2. One review can miss a whole lens

[Time: 0:45–1:35]
The useful part of the idea is the division of attention. A correctness reviewer follows edge cases and tests. An architecture reviewer follows interfaces and contracts. An operational reviewer asks what happens when the happy path has been running for six months and something breaks overnight. The full roster also includes simplicity, security, cost and technical language. I am showing three examples because the mechanism matters more than memorising seven robots. If independent attention helps, it should produce supported findings that a fair baseline misses. More comments alone would be a poor measure: a reviewer can manufacture work by inventing problems.

[Sources]
https://github.com/davehowell/council-of-nark/tree/main/agents
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md

## 3. The council has two jobs

[Time: 1:35–2:25]
The proposed workflow has two jobs. First, reviewers work independently on the same artifact. Second, a fuser combines the findings into something a person can use. Those jobs can fail separately. The panel might discover a real issue that the final answer drops. Or every reviewer might repeat the same unsupported claim and the fuser might make it sound convincing. That is why I need to inspect the raw collection as well as the final verdict. A polished answer is not enough. The council is review-only: these initial experiments measure what it says about an artifact, not whether it successfully implements a fix.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md

## 4. The characters are another hypothesis

[Time: 2:25–3:15]
The characters make the roles memorable, and they make the project fun. But memorability for me and better performance by the model are different questions. If I add K-2SO and also improve the correctness instructions, I cannot credit the character for any gain. The controlled comparison keeps the functional job byte-identical and changes only its wrapper. Then we can ask whether the prose helps, hinders, or makes a difference too small to matter. We are measuring observable behaviour. We cannot infer which internal features of the model caused it, and we should not confuse a distinctive voice with a useful review.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/prompts/README.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/PERSONA_FACTORIAL.md

## 5. The cheap alternative deserves a trial

[Time: 3:15–4:10]
The council should face alternatives that someone might actually choose. Start with one well-instructed reviewer. Give one continuing agent a chance to revisit its work. Then compare independent repeated reviews before adding specialist roles. The last two are not the same: an agent that sees its first answer may correct it or become anchored to it. Multiple calls also reuse context and may spend far more tokens than a single call. So I want both an ordinary cheap baseline and a comparison at an equal overall resource budget. A council could be better yet not worth its price, or worth using only for particular high-risk tasks.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 6. Let the council lose

[Time: 4:10–5:00]
That is the commitment behind the experiment. I do not need the council to win. If a simple functional prompt delivers the same useful outcome, that is a good result: I can save tokens and keep the names as convenient labels. If the council wins only under particular conditions, I want to identify those conditions. And if the experiment cannot distinguish the approaches, I need to say that rather than call it a tie. The next talk explains the controls that make those outcomes distinguishable. The interesting story then comes from discovering that the measurement system itself needed almost as much review as the agents did.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md

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

# The Experiment Fought Back — speaker notes

Six slides · approximately five minutes · 636 spoken words.

Edit the closing comment in each slide in part-3/slides.md; these notes are generated from that source.

## 1. The Experiment Fought Back

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
The course correction was to move towards real pre-fix open-source tasks. That required stronger isolation: an empty respondent environment, history-free source, a narrow mediator for source access, and tests with network access denied. Gortex's frozen regression has been shown to fail before its fix and pass afterwards. Several failed exporter and process-lifecycle attempts were preserved on the way. This is valuable infrastructure evidence. It is not evidence that a reviewer found the bug: no ecological respondent call has occurred. That distinction is the bridge to the final talk. We now know what the existing results support, and can define a finite closing test without claiming it has already happened.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CHECKPOINT.md

---

# What Would Settle It? — speaker notes

Six slides · approximately five minutes · 655 spoken words.

Edit the closing comment in each slide in part-4/slides.md; these notes are generated from that source.

## 1. What Would Settle It?

[Time: 0:00–0:45]
We can close the current evidence review now, but we cannot honestly close every empirical question with a definitive winner. The council has not demonstrated a reliable advantage over one functional review. Nor have we established equivalence. One persona underperformed in one controlled calibration. The question is what additional work would change that confidence enough to be useful. My proposal is a smaller final comparison focused on practical review quality and resource use. It deliberately leaves provider diversity, every chain order and the full eight-role factorial for another study. This talk separates that proposed work from the results we already have.

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

## 5. A small pilot cannot prove a small gap

[Time: 3:25–4:15]
Four tasks with two repeats across these arms require 136 planned review stages, before the separate initial feasibility block, tool-driven provider turns or retries. That can test whether the workflow completes and whether the rubric and budgets make sense. It cannot establish narrow equivalence. An illustrative precision calculation in the protocol can require around fifty-seven independent tasks under one assumed variability level, and repository clustering can require more. That is not a sample-size promise; it shows why adding repeats to four tasks is not a shortcut to generality. After excluded feasibility, freeze the confirmatory sample and spending cap. If the available resources cannot buy the required precision, stop and report the limit.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md

## 6. The council has not earned a default

[Time: 4:15–5:00]
My present operational choice is to start with one functional review and escalate for tasks whose risk justifies more attention. That is a decision under uncertainty, not a theorem that one agent is equivalent. The strongest persona result remains narrow: K-2SO reduced detection on this model and these three packets, with slightly fewer tokens. The council results remain mixed and exploratory. The closing protocol gives us a way to earn a stronger statement, including an equivalence result if the evidence is precise enough. Until the launcher, task environments and two human raters are ready, it remains a proposal. The useful outcome is a claim we can defend, even when that claim is smaller than the original idea.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/CONCLUSION.md
