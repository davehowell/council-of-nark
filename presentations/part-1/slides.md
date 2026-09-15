---
theme: default
title: The Council of Nark
colorSchema: dark
fonts:
  sans: Arial
  local: Arial
transition: none
layout: default
class: cover origin-cover
---

<div class="cover-grid">
<div>

<p class="eyebrow">01 / The original hypothesis</p>

# The Council<br>of Nark

<p class="lead">Why try several specialist reviewers instead of one general review?</p>

<p class="origin-note">The idea grew out of Fable 5's multi-agent review mode.<br>It ran while I planned across projects on paper.</p>

</div>
<figure class="fable-card">
  <img src="./fable-5.jpg" alt="Fable 5 interface used as the starting point for the council idea">
  <figcaption>Starting point: several reviews combined into one response</figcaption>
</figure>
</div>

<div class="foot"><span>Council of Nark</span><span>1 / 10</span></div>

<!--
[Time: 0:00–0:50]
The Council of Nark began with a workflow I was already using. Fable 5 could ask several agents to inspect the same task and combine their replies. The run was slower and used more tokens than one call. While it ran, I often stepped away from the screen and planned across projects on A3 paper. I also used Gemini for second opinions and an HK-47 reviewer for unnecessary complexity. Those tools did not prove that a panel was better. They gave me a practical question: will reviewers with different, stable responsibilities find useful problems that one general review misses? The diagrams in this talk show that original hypothesis, not measured results.

[Sources]
https://davehowell.github.io/experimental/
https://github.com/davehowell/council-of-nark/blob/main/README.md
-->

---
class: origin-slide
---

# The workflow already mixed reviewers

<p class="slide-intro">The proposal made an informal practice explicit.</p>

<div class="origin-cards">
  <div class="origin-card">
    <span class="card-index">01</span>
    <strong>Fable 5</strong>
    <p>Several agent reviews, followed by one combined response.</p>
  </div>
  <div class="origin-card">
    <span class="card-index">02</span>
    <strong>Gemini</strong>
    <p>A separate model for research, challenge and a second opinion.</p>
  </div>
  <div class="origin-card">
    <span class="card-index">03</span>
    <strong>HK-47</strong>
    <p>A fixed role that looks for code and process we do not need.</p>
  </div>
</div>

<div class="proposal-strip"><strong>Proposal</strong><span>Give each reviewer one clear responsibility, run relevant reviewers independently, then combine their findings once.</span></div>

<div class="foot"><span>Where the idea came from</span><span>2 / 10</span></div>

<!--
[Time: 0:50–1:40]
The workflow already mixed models and instructions. Gemini gave me an independent answer from another provider. HK-47 used the same underlying model as other agents but had one narrow job: challenge over-engineering. Fable showed a convenient way to run several reviews and collect the output. The council made those parts explicit. Each reviewer would have a stable checklist. The controller would choose only the reviewers relevant to the task. They would inspect the same source independently, so one review would not steer the next. Repeated findings might strengthen confidence, while a supported finding from one reviewer might expose a blind spot. Better review quality was still only a hypothesis.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/tree/main/agents
-->

---
class: council-reveal
---

# Meet the council

<img class="council-art" src="./Nark-council.png" alt="Seven robot reviewers gathered around the Council of Nark console">

<p class="art-caption">Seven memorable characters. Seven fixed review responsibilities.</p>

<div class="foot"><span>The roles make the review lenses easy to recall</span><span>3 / 10</span></div>

<!--
[Time: 1:40–2:05]
The characters made the responsibilities easier to remember and made the project enjoyable to use. That was their known purpose. Whether character wording improved a model's review was a separate question. The functional instructions had to stand on their own. A reviewer needed a defined scope, evidence for each finding and a common output format. The names could remain useful labels even if an experiment found no performance benefit from the fictional wording.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/README.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/PERSONA_FACTORIAL.md
-->

---
class: roster-slide
---

# Who reviews what?

<p class="slide-intro">Each reviewer has a narrow lens. None of them edits the code, plan or document.</p>

<table class="council-table">
<thead><tr><th></th><th>Reviewer</th><th>Responsibility</th></tr></thead>
<tbody>
<tr><td><img src="./hk47.png" alt=""></td><td><strong>HK-47</strong></td><td>Remove unnecessary complexity and speculative machinery.</td></tr>
<tr><td><img src="./K-2SO.png" alt=""></td><td><strong>K-2SO</strong></td><td>Check correctness, tests, observability and data edge cases.</td></tr>
<tr><td><img src="./glados.png" alt=""></td><td><strong>GLaDOS</strong></td><td>Check architecture, contracts, coupling and migrations.</td></tr>
<tr><td><img src="./c3po.png" alt=""></td><td><strong>C-3PO</strong></td><td>Check secrets, personal data, permissions and destructive actions.</td></tr>
<tr><td><img src="./bender.png" alt=""></td><td><strong>Bender</strong></td><td>Check compute use and avoidable cost.</td></tr>
<tr><td><img src="./holly.png" alt=""></td><td><strong>Holly</strong></td><td>Check what breaks after months of neglect or staff change.</td></tr>
<tr><td><img src="./walle.png" alt=""></td><td><strong>WALL-E</strong></td><td>Check durable technical language and translate the final review.</td></tr>
</tbody>
</table>

<div class="foot"><span>The controller selects the smallest relevant group</span><span>4 / 10</span></div>

<!--
[Time: 2:05–3:05]
The roster covers seven kinds of risk. HK-47 asks whether the solution contains machinery that the problem does not require. K-2SO checks behaviour, tests and failure signals. GLaDOS follows boundaries between components. C-3PO checks sensitive data, permissions and destructive operations. Bender examines compute and cost. Holly takes the long view: manual steps, forgotten toggles, state drift and knowledge that can disappear when people leave. WALL-E checks whether technical writing will remain clear to its next reader. The controller does not run all seven by default. It chooses the smallest group that matches the code, plan or document, with at least two independent reviewers.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/README.md
-->

---
class: handoff-slide
---

# Holly and WALL-E work at different points

<div class="seat-grid">
  <div class="seat-card">
    <div class="seat-head"><img src="./holly.png" alt="Holly"><div><strong>Holly</strong><small>Independent panel review</small></div></div>
    <p class="seat-question">What happens after six months of ordinary neglect?</p>
    <ul>
      <li>Forgotten toggles and manual steps</li>
      <li>State drift and missing ownership</li>
      <li>Noisy alerts and difficult handovers</li>
    </ul>
  </div>
  <div class="seat-card">
    <div class="seat-head"><img src="./walle.png" alt="WALL-E"><div><strong>WALL-E</strong><small>Writing review and final translation</small></div></div>
    <p class="seat-question">Can a tired reader understand the decision once?</p>
    <ul>
      <li>Unclear conditions and vague warnings</li>
      <li>Dense jargon and inconsistent terms</li>
      <li>A plain-language version of the combined review</li>
    </ul>
  </div>
</div>

<div class="translation-flow"><span>independent findings</span><b>→</b><span class="accent">arbiter merges and ranks</span><b>→</b><span class="warm">WALL-E translates</span><b>→</b><span>human decision</span></div>

<div class="foot"><span>Review and translation are separate jobs</span><span>5 / 10</span></div>

<!--
[Time: 3:05–3:55]
Holly and WALL-E widened the design in different ways. Holly sits on the review panel when a change depends on memory, manual work or long-term operational care. WALL-E can review prose as a specialist, but also has a separate job after the arbiter finishes. The arbiter first merges duplicate findings, resolves conflicts and ranks the result. WALL-E then translates that settled review for people who do not need file-and-line details. Translation does not create another vote and does not count as a panel review. Keeping those jobs separate prevents a simpler explanation from quietly changing the technical decision.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/skills/walle-ste/SKILL.md
-->

---
class: concept-slide
---

# One review can produce different outcomes

<div class="concept-label">Visual model of the original hypothesis</div>

<div class="concept-layout">
<div>
<svg viewBox="0 0 520 300" class="concept-chart" role="img" aria-label="A conceptual cone showing that one review can improve or worsen the starting analysis">
  <defs><marker id="one-tip" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#91d9cb"/></marker></defs>
  <line x1="58" y1="28" x2="58" y2="260" class="axis"/>
  <line x1="58" y1="260" x2="500" y2="260" class="axis"/>
  <text x="64" y="22" class="axis-label">review quality ↑</text>
  <text x="306" y="286" class="axis-label">review process →</text>
  <polygon points="64,186 492,62 492,244" class="range neutral"/>
  <line x1="64" y1="186" x2="486" y2="150" class="mean-line" marker-end="url(#one-tip)"/>
  <line x1="64" y1="186" x2="492" y2="186" class="baseline"/>
  <circle cx="64" cy="186" r="5" class="start-dot"/>
  <text x="72" y="177" class="chart-label">same task and prompt</text>
  <text x="426" y="79" class="chart-label good-text">strong</text>
  <text x="430" y="235" class="chart-label warm-text">weak</text>
  <text x="404" y="139" class="chart-label">expected middle</text>
</svg>
</div>
<div class="concept-copy">
  <p>The same model and prompt can give a stronger or weaker review on different runs.</p>
  <p>The original shorthand was a <strong>cone of possible outcomes</strong>.</p>
  <p>Better instructions can move the whole range, but they do not remove variation.</p>
</div>
</div>

<div class="foot"><span>The cone is a picture of the hypothesis</span><span>6 / 10</span></div>

<!--
[Time: 3:55–4:45]
This cone was a way to draw variation, not a statistical model fitted to observations. Start with the same code and prompt. A model can produce a strong review, an average one or a weak one. More reasoning does not guarantee a better answer, and a later run can miss something that an earlier run found. Better context and clearer instructions can improve the starting conditions, but they do not make every output identical. The diagram gave me language for the next question: what happens when several model calls pass work from one to another?

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md
-->

---
class: concept-slide
---

# A serial chain can carry mistakes forward

<div class="concept-label">Original hypothesis: serial review</div>

<div class="concept-layout">
<div>
<svg viewBox="0 0 520 300" class="concept-chart" role="img" aria-label="A conceptual diagram showing a second review branching from strong and weak first reviews">
  <line x1="58" y1="28" x2="58" y2="260" class="axis"/>
  <line x1="58" y1="260" x2="500" y2="260" class="axis"/>
  <text x="64" y="22" class="axis-label">review quality ↑</text>
  <text x="306" y="286" class="axis-label">review sequence →</text>
  <polygon points="64,186 270,122 270,224" class="range neutral"/>
  <polygon points="270,122 492,54 492,178" class="range good"/>
  <polygon points="270,224 492,188 492,254" class="range weak"/>
  <line x1="270" y1="38" x2="270" y2="260" class="divider"/>
  <circle cx="270" cy="122" r="5" class="good-dot"/>
  <circle cx="270" cy="224" r="5" class="weak-dot"/>
  <text x="116" y="173" class="chart-label">first review</text>
  <text x="314" y="42" class="chart-label good-text">strong first answer</text>
  <text x="316" y="250" class="chart-label warm-text">weak first answer</text>
</svg>
</div>
<div class="concept-copy">
  <p>In a chain, each reviewer sees the previous answer as part of its starting point.</p>
  <p>A later reviewer can correct an error, repeat it or remove a valid finding.</p>
  <p>The order of reviewers can therefore change the final review.</p>
</div>
</div>

<div class="foot"><span>Passing output forward creates order effects</span><span>7 / 10</span></div>

<!--
[Time: 4:45–5:35]
In a serial chain, the second reviewer does not begin independently. It inherits the first answer. A strong first answer can still lose a valid finding during the next rewrite. A weak first answer can lead the next reviewer towards the wrong concern. The second reviewer can also correct the first, so the diagram does not claim that chains always get worse. It shows why the order matters and why a chain needs its own test. The protocol later compared every ordering of a three-reviewer chain instead of choosing one favourable sequence.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md
-->

---
class: concept-slide
---

# One arbiter combines independent reviews

<div class="concept-label">Original hypothesis: independent review and fusion</div>

<div class="concept-layout">
<div>
<svg viewBox="0 0 520 300" class="concept-chart" role="img" aria-label="A conceptual diagram comparing the range before fusion with the expected range after fusion">
  <defs><marker id="fuse-tip" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#91d9cb"/></marker></defs>
  <line x1="58" y1="28" x2="58" y2="260" class="axis"/>
  <line x1="58" y1="260" x2="500" y2="260" class="axis"/>
  <text x="64" y="22" class="axis-label">review quality ↑</text>
  <text x="306" y="286" class="axis-label">review process →</text>
  <polygon points="64,190 492,64 492,244" class="range before"/>
  <text x="418" y="235" class="chart-label muted-text">one review</text>
  <polygon points="64,160 492,78 492,188" class="range fused"/>
  <line x1="64" y1="160" x2="486" y2="132" class="mean-line" marker-end="url(#fuse-tip)"/>
  <circle cx="64" cy="190" r="4" class="muted-dot"/>
  <circle cx="64" cy="160" r="5" class="good-dot"/>
  <line x1="64" y1="184" x2="64" y2="167" class="lift-line" marker-end="url(#fuse-tip)"/>
  <text x="224" y="153" class="chart-label good-text">expected after combining reviews</text>
</svg>
</div>
<div class="concept-copy">
  <p>Every specialist receives the same source, without seeing another review.</p>
  <p>The arbiter removes duplicates and keeps supported findings raised by only one reviewer.</p>
  <p><strong>Hypothesis:</strong> this raises useful coverage and reduces weak final reviews.</p>
</div>
</div>

<div class="foot"><span>Combining reviews is a claim to test, not a quality guarantee</span><span>8 / 10</span></div>

<!--
[Time: 5:35–6:25]
The alternative used parallel reviews and one combining step. Each specialist receives the same source and works independently. The arbiter sees all of the findings together. It can merge duplicate reports, reject unsupported claims and preserve a useful finding that only one specialist noticed. My original expectation was that this process would raise the typical quality and reduce the weak end of the range. That expectation could be wrong. The combining step can also discard a correct finding or make repeated errors sound authoritative. The experiment therefore had to score the individual findings and the combined review separately.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/METRICS.md
-->

---
class: topology-slide
---

# Flower, not a chain

<p class="slide-intro">The same reviewers can be connected in two different ways.</p>

<div class="topology-grid">
  <div class="topology-panel good-panel">
    <div class="topology-title">FAN OUT, THEN FUSE</div>
    <div class="flower">
      <svg viewBox="0 0 340 250" aria-hidden="true">
        <g class="flower-links">
          <line x1="170" y1="121" x2="170" y2="34"/>
          <line x1="156" y1="128" x2="64" y2="74"/>
          <line x1="184" y1="128" x2="276" y2="74"/>
          <line x1="156" y1="146" x2="70" y2="204"/>
          <line x1="184" y1="146" x2="270" y2="204"/>
        </g>
      </svg>
      <div class="arbiter"><b>ARBITER</b><small>compare · merge · rank</small></div>
      <img class="avatar" src="./hk47.png" alt="HK-47" style="left:145px;top:5px">
      <img class="avatar" src="./K-2SO.png" alt="K-2SO" style="left:38px;top:49px">
      <img class="avatar" src="./glados.png" alt="GLaDOS" style="left:252px;top:49px">
      <img class="avatar" src="./c3po.png" alt="C-3PO" style="left:45px;top:178px">
      <img class="avatar" src="./bender.png" alt="Bender" style="left:245px;top:178px">
    </div>
    <p>Each reviewer sees the source. The arbiter compares independent findings once.</p>
  </div>

  <div class="topology-panel chain-panel">
    <div class="topology-title">PASS THE ANSWER ALONG</div>
    <div class="chain">
      <svg viewBox="0 0 390 250" aria-hidden="true">
        <g class="chain-links">
          <line x1="54" y1="125" x2="118" y2="65"/>
          <line x1="54" y1="125" x2="118" y2="185"/>
          <line x1="151" y1="65" x2="215" y2="125"/>
          <line x1="151" y1="185" x2="215" y2="125"/>
          <line x1="248" y1="125" x2="312" y2="65"/>
          <line x1="248" y1="125" x2="312" y2="185"/>
        </g>
      </svg>
      <img class="avatar" src="./hk47.png" alt="HK-47" style="left:26px;top:101px">
      <img class="avatar" src="./K-2SO.png" alt="K-2SO" style="left:119px;top:40px">
      <img class="avatar" src="./glados.png" alt="GLaDOS" style="left:119px;top:160px">
      <img class="avatar" src="./c3po.png" alt="C-3PO" style="left:216px;top:101px">
      <img class="avatar" src="./bender.png" alt="Bender" style="left:313px;top:40px">
      <img class="avatar" src="./holly.png" alt="Holly" style="left:313px;top:160px">
    </div>
    <p>Each reviewer inherits earlier output. Sequence and omissions can affect the result.</p>
  </div>
</div>

<div class="foot"><span>How reviewers are connected is part of the experiment</span><span>9 / 10</span></div>

<!--
[Time: 6:25–7:20]
These diagrams show the structural difference. On the left, reviewers work from the same source and return findings to one arbiter. The arbiter can compare disagreements because no reviewer has rewritten another review first. On the right, output passes through a sequence. Later reviewers can see and change earlier work, so the path through the reviewers becomes part of the result. I called the first shape a flower because all work returns to one centre. The useful claim was not that flowers are inherently better. It was that independent review followed by one combining step should be compared with serial review under controlled conditions.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/skills/nark-matrix/SKILL.md
-->

---
class: closing-slide
---

# A testable proposal, not a conclusion

<div class="test-grid">
  <div><span>01</span><strong>Specialist roles</strong><p>Compare them with the same number of repeated general reviews.</p></div>
  <div><span>02</span><strong>Character wording</strong><p>Keep the functional instructions fixed and change only the character text.</p></div>
  <div><span>03</span><strong>Review structure</strong><p>Compare fan-out and fusion with every order of a short chain.</p></div>
</div>

<p class="closing-rule">Keep extra reviewers only when they improve useful findings enough to justify their time and cost.</p>

<div class="faces" aria-label="The seven council reviewers">
  <img src="./hk47.png" alt="HK-47"><img src="./K-2SO.png" alt="K-2SO"><img src="./glados.png" alt="GLaDOS"><img src="./c3po.png" alt="C-3PO"><img src="./bender.png" alt="Bender"><img src="./holly.png" alt="Holly"><img src="./walle.png" alt="WALL-E">
</div>

<div class="foot"><span>Next: design a fair test</span><span>10 / 10</span></div>

<!--
[Time: 7:20–8:00]
The original idea produced three claims that could be tested separately. First, specialist roles might find different real defects from repeated general reviews. Second, fictional character wording might help, hinder or make no useful difference when the functional instructions stay fixed. Third, independent fan-out and one fusion step might produce a better final review than a chain. Cost, tokens and time belong beside quality in every comparison. The desired outcome was never the largest council. It was the smallest review process that produced a dependable result. The next talk explains how the protocol tried to separate those claims.

[Sources]
https://github.com/davehowell/council-of-nark/blob/main/experiment/protocol.md
https://github.com/davehowell/council-of-nark/blob/main/experiment/FINAL_ROUND.md
-->
