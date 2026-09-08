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
