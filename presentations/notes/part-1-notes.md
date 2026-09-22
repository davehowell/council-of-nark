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
