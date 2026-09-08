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
