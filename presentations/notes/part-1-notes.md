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
