---
theme: seriph
title: The Council of Nark
colorSchema: light
info: |
  AI check-in 001 — heterogeneous multi-agent review, fused by a single arbiter.
  Holly audits what rots over time; WALL-E translates the fused verdict for humans.
transition: slide-left
mdc: true
layout: image-right
image: /fable-5.jpg
backgroundSize: contain
---

# The Council of Nark

#### From Fable 5's swarm to a panel of narky bots

<v-clicks>

- **Fable 5, ultra / max effort** — a *swarm* of subagents whose outputs are fused into one stronger result.
- Like a model's internal reasoning tokens — but **external, parallel, inspectable**.
- Slow and token-hungry. On the Claude **Max** plan I don't hit thresholds, and I hold high standards — so the trade is worth it (for me).
- While it deliberates, I step away to **A3 paper** — longer-horizon, cross-project planning by hand.

</v-clicks>

<!--
Open on the trade-off: max effort buys quality at the cost of time + tokens. The
"step away to paper" point is the human-in-the-loop framing — the model deliberates,
I plan across projects. This sets up why a heterogeneous panel is worth the wait.
-->

---
layout: default
---

# Already heterogeneous

I'd been mixing models for a while — Fable's swarm just nudged me to make it deliberate.

<div class="grid grid-cols-2 gap-8 mt-6">
<div>

#### What I was already running

<v-clicks>

- **Gemini** (via MCP / `agy`) — second opinion, web research, cross-model sanity checks.
- **HK-47** — a narky reviewer subagent that hunts over-engineering and complexity demons.

</v-clicks>

</div>
<div>

#### The idea

<v-clicks>

- Different models — and the *same* model in a different **role** — have different strengths.
- Round out a **staple of narky bots**, each holding a fixed angle on *plan · build · review*.
- Gain from where they **overlap** (consensus) and where they **don't** (edge cases).

</v-clicks>

</div>
</div>

<!--
The pitch: heterogeneity isn't only across providers — a persona is a *role* that
shifts the same model's attention. HK-47 was the proof of concept; the council
generalises it.
-->

---
layout: image
image: /Nark-council.png
backgroundSize: cover
class: flex items-start justify-center
---

<h1 class="!text-5xl !text-white !font-normal tracking-wide mt-10 drop-shadow-[0_2px_10px_rgba(0,0,0,0.95)]">
Introducing the Council of Nark
</h1>

<!--
Let the image breathe. One line, no box — just a shadow so it reads over the art.
Beat, then move to the roster.
-->

---
layout: default
---

# The Council

Each bot holds one angle. The arbiter convenes the relevant subset, then fuses.

<table class="mt-3 council-table">
<thead><tr><th></th><th>Bot</th><th>Angle</th><th>Provider</th></tr></thead>
<tbody>
<tr><td><img src="./hk47.png" /></td><td><b>HK-47</b></td><td>simplicity · anti-over-engineering</td><td>Claude</td></tr>
<tr><td><img src="./K-2SO.png" /></td><td><b>K-2SO</b></td><td>data quality · tests · observability</td><td>Claude</td></tr>
<tr><td><img src="./glados.png" /></td><td><b>GLaDOS</b></td><td>architecture · contracts · synthesis</td><td>Claude</td></tr>
<tr><td><img src="./c3po.png" /></td><td><b>C-3PO</b></td><td>PII · secrets · IAM · compliance</td><td>Claude</td></tr>
<tr><td><img src="./bender.png" /></td><td><b>Bender</b></td><td>FinOps · cost and compute waste</td><td><b>Gemini</b></td></tr>
<tr><td><img src="./holly.png" /></td><td><b>Holly</b></td><td>six-month entropy · forgotten state</td><td><b>OpenAI</b></td></tr>
<tr><td><img src="./walle.png" /></td><td><b>WALL-E</b></td><td>plain technical prose · translation</td><td>Claude</td></tr>
</tbody>
</table>

<div class="text-xs opacity-70 mt-2">3 providers · structured findings · review-only. WALL-E also translates every fused verdict.</div>

<style>
.council-table td, .council-table th { vertical-align: middle; padding: 2px 14px; border: none; font-size: 14px; }
.council-table img { height:2.25rem; width:2.25rem; border-radius:9999px; object-fit:cover; }
</style>

<!-- Public assets are served from root paths. Holly and WALL-E are cropped from Nark-council.png. -->

<!--
Bender and Holly make provider diversity real rather than theatrical. Every panel
reviewer emits the same findings shape. WALL-E has a second, separate job after
fusion: make the verdict usable by people who do not read file:line diagnostics.
-->

---
layout: default
---

# Two new seats, at different moments

<div class="grid grid-cols-2 gap-10 mt-5">
  <div class="new-seat">
    <div class="seat-head"><img src="./holly.png" /><div><b>Holly</b><small>on the panel · OpenAI</small></div></div>
    <p class="seat-q">“What does this look like after six months of neglect?”</p>
    <ul>
      <li>Forgotten toggles and manual steps</li>
      <li>State drift and knowledge evaporation</li>
      <li>Signal fatigue and the 03:00 hand-off</li>
    </ul>
  </div>
  <div class="new-seat">
    <div class="seat-head"><img src="./walle.png" /><div><b>WALL-E</b><small>after fusion · Claude</small></div></div>
    <p class="seat-q">“Can a tired human understand this once?”</p>
    <ul>
      <li>Reviews technical writing with STE-derived rules</li>
      <li>Compacts jargon, vague warnings and buried conditions</li>
      <li>Translates every fused verdict for non-experts</li>
    </ul>
  </div>
</div>

<div class="translation-flow mt-7">
  <span>specialist findings</span><b>→</b><span class="fuse">GLaDOS / controller<br><small>dedupe · rank · arbitrate</small></span><b>→</b><span class="walle">WALL-E<br><small>plain language</small></span><b>→</b><span>humans</span>
</div>

<style>
.new-seat { border:1px solid #cbd5e1; border-radius:14px; padding:18px 20px; background:#f8fafc; }
.seat-head { display:flex; align-items:center; gap:14px; font-size:22px; }
.seat-head img { width:64px; height:64px; border-radius:9999px; object-fit:cover; }
.seat-head small { display:block; font-size:11px; font-weight:400; opacity:.65; }
.seat-q { color:#475569; font-size:14px; font-style:italic; margin:10px 0 8px; }
.new-seat ul { font-size:14px; line-height:1.55; padding-left:18px; }
.translation-flow { display:flex; justify-content:center; align-items:center; gap:12px; font-size:13px; text-align:center; }
.translation-flow span { border:1px solid #94a3b8; border-radius:9px; padding:8px 13px; background:white; }
.translation-flow .fuse { border-color:#8b5cf6; background:#f5f3ff; }
.translation-flow .walle { border-color:#d97706; background:#fffbeb; }
.translation-flow small { opacity:.65; font-size:10px; }
</style>

<!-- Holly widens the time horizon and adds a third provider. WALL-E can review prose
on the panel, but his universal job is downstream: translate the already-fused verdict.
That translator pass does not count toward the minimum panel size. -->

---
layout: default
---

# A single call is a cone

<div class="grid grid-cols-5 gap-6 items-center mt-2">
<div class="col-span-3">

<svg viewBox="0 0 520 320" class="w-full">
<defs><marker id="ar1" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#b45309"/></marker></defs>
<line x1="60" y1="30" x2="60" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<line x1="60" y1="280" x2="500" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<text x="64" y="24" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">↑ order / quality</text>
<text x="300" y="304" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">time · reasoning steps →</text>
<polygon points="60,210 500,80 500,268" fill="#f59e0b" fill-opacity="0.14" stroke="#f59e0b" stroke-opacity="0.5" stroke-dasharray="5 3"/>
<polyline points="60,210 110,204 160,196 210,186 260,176 310,162 360,150 410,138 460,124 500,112" fill="none" stroke="#9ca3af" stroke-width="1" stroke-opacity="0.45"/>
<polyline points="60,210 110,216 160,212 210,224 260,232 310,240 360,246 410,254 460,258 500,260" fill="none" stroke="#9ca3af" stroke-width="1" stroke-opacity="0.45"/>
<line x1="60" y1="210" x2="494" y2="170" stroke="#b45309" stroke-width="2.5" marker-end="url(#ar1)"/>
<text x="398" y="158" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#b45309">average</text>
<text x="330" y="96" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#b45309">one model · one pass</text>
<circle cx="60" cy="210" r="4.5" fill="#b45309"/>
<text x="66" y="228" style="font-size:11px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">start · prompt</text>
</svg>

</div>
<div class="col-span-2">

<v-clicks>

- A single model call = a **cone of outcomes** fanning from the prompt — it can land **above** the start (better) or **below** it (worse: hallucination, regressions). A smarter model lifts the cone; better context lifts the **start**.
- More steps ≠ more order. Chaining LLM→LLM→LLM **reseeds** the cone at each endpoint — entropy flares, best *and* worst amplify.

</v-clicks>

</div>
</div>

<!--
Frame the axes: X is time / a model's step from prompt to response; Y is order /
anti-chaos / a proxy for quality. One call is a probability cone. The faint jagged
lines are sample random walks. Set up the trap: serial chaining widens, not narrows.
-->

---
layout: default
---

# Hypothesis: chaining flares the cone

<div class="grid grid-cols-5 gap-6 items-center mt-2">
<div class="col-span-3">

<svg viewBox="0 0 520 320" class="w-full">
<line x1="60" y1="30" x2="60" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<line x1="60" y1="280" x2="500" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<text x="64" y="24" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">↑ order / quality</text>
<text x="300" y="304" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">time · reasoning steps →</text>
<line x1="280" y1="66" x2="280" y2="280" stroke="#94a3b8" stroke-width="1" stroke-opacity="0.35" stroke-dasharray="3 4"/>
<polygon points="60,210 280,145 280,239" fill="#f59e0b" fill-opacity="0.14" stroke="#f59e0b" stroke-opacity="0.5" stroke-dasharray="5 3"/>
<polygon points="280,145 500,72 500,200" fill="#f59e0b" fill-opacity="0.13" stroke="#f59e0b" stroke-opacity="0.55"/>
<polygon points="280,239 500,206 500,278" fill="#dc2626" fill-opacity="0.16" stroke="#dc2626" stroke-opacity="0.6"/>
<circle cx="60" cy="210" r="4" fill="#94a3b8"/>
<circle cx="280" cy="145" r="4.5" fill="#b45309"/>
<circle cx="280" cy="239" r="4.5" fill="#dc2626"/>
<text x="108" y="192" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#b45309">1st call</text>
<text x="346" y="104" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#b45309">from best case</text>
<text x="344" y="262" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#dc2626">from worst case</text>
<text x="286" y="60" style="font-size:11px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">2nd call →</text>
</svg>

</div>
<div class="col-span-2">

<v-clicks>

- Chain LLM→LLM and each endpoint becomes a **fresh start** — a new cone.
- From the best case (**orange**) it can still **regress**; from the worst case (**red**) it tends to **stay bad**, and good work gets discarded.
- The union of outcomes is **wider** than one call — variance compounds (Markov × Monte-Carlo).

</v-clicks>

</div>
</div>

<!--
The trap: a DAG of LLMs reseeds the cone at every hop. Each intermediate endpoint
becomes a new apex, so the best case fans up-and-down again and the worst case mostly
stays bad. The reachable spread after two calls is wider than after one.
-->

---
layout: default
---

# One arbiter, not a DAG

<div class="grid grid-cols-5 gap-6 items-center mt-2">
<div class="col-span-3">

<svg viewBox="0 0 520 320" class="w-full">
<defs><marker id="ar2" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#047857"/></marker></defs>
<line x1="60" y1="30" x2="60" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<line x1="60" y1="280" x2="500" y2="280" stroke="#94a3b8" stroke-width="1.5"/>
<text x="64" y="24" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">↑ order / quality</text>
<text x="300" y="304" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#64748b">time · reasoning steps →</text>
<polygon points="60,210 500,80 500,268" fill="#9ca3af" fill-opacity="0.08" stroke="#9ca3af" stroke-opacity="0.3" stroke-dasharray="4 4"/>
<line x1="60" y1="210" x2="500" y2="172" stroke="#9ca3af" stroke-width="1.5" stroke-opacity="0.5" stroke-dasharray="4 4"/>
<text x="430" y="184" style="font-size:11px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#9ca3af">before</text>
<polygon points="60,180 500,92 500,202" fill="#10b981" fill-opacity="0.28" stroke="#10b981" stroke-opacity="0.75"/>
<line x1="60" y1="180" x2="494" y2="147" stroke="#047857" stroke-width="2.5" marker-end="url(#ar2)"/>
<text x="410" y="136" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#047857">average</text>
<text x="238" y="172" style="font-size:12px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#047857">consult + fuse</text>
<circle cx="60" cy="210" r="4" fill="#9ca3af"/>
<circle cx="60" cy="180" r="4.5" fill="#047857"/>
<line x1="60" y1="206" x2="60" y2="184" stroke="#047857" stroke-width="1.5" marker-end="url(#ar2)"/>
<text x="66" y="202" style="font-size:11px;font-family:ui-sans-serif,system-ui,sans-serif" fill="#047857">start lifted</text>
</svg>

</div>
<div class="col-span-2">

The **main window is the arbiter** — it consults the council, then folds the replies back in.

<v-clicks>

- Working hypothesis: fusion **lifts the start** and **narrows the cone** by trimming the low tail.
- Overlap can signal consensus; divergence can expose edge cases.
- It is not a hallucination cure. The follow-up experiment tests whether the cone actually changes.

</v-clicks>

</div>
</div>

<!--
The payoff: folding a subagent's reply back into the arbiter ≈ lifting the start
point and narrowing the cone — trimming the low tail rather than reseeding it (which
serial chaining does). The green average bisects the green cone; the lower edge still
slopes down (fusion reduces, never eliminates, the downside).
-->

---
layout: default
---

# Flower, not a chain

Same reviewers, two topologies — the experiment will test whether shape decides the spread.

<div class="topo-wrap">

  <div class="panel">
    <div class="topo-head good">ONE ARBITER · HUB &amp; FUSE</div>
    <div class="flower">
      <svg viewBox="0 0 340 340" width="340" height="340">
        <defs>
          <marker id="gtip" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#059669"/></marker>
        </defs>
        <g stroke="#10b981" stroke-opacity="0.15" stroke-width="24" stroke-linecap="round">
          <line x1="170" y1="134" x2="170" y2="48"/>
          <line x1="204" y1="159" x2="286" y2="132"/>
          <line x1="191" y1="199" x2="242" y2="269"/>
          <line x1="149" y1="199" x2="98" y2="269"/>
          <line x1="136" y1="159" x2="54" y2="132"/>
        </g>
        <g stroke="#059669" stroke-width="1.6" marker-start="url(#gtip)" marker-end="url(#gtip)">
          <line x1="170" y1="126" x2="170" y2="74"/>
          <line x1="212" y1="156" x2="261" y2="140"/>
          <line x1="196" y1="206" x2="226" y2="248"/>
          <line x1="144" y1="206" x2="114" y2="248"/>
          <line x1="128" y1="156" x2="79" y2="140"/>
        </g>
      </svg>
      <div class="arbiter">&gt;_<b>Arbiter</b><small>main session</small></div>
      <img class="avatar" src="./hk47.png"   style="left:145px; top:20px" />
      <img class="avatar" src="./K-2SO.png"  style="left:264px; top:106px" />
      <img class="avatar" src="./glados.png" style="left:219px; top:246px" />
      <img class="avatar" src="./c3po.png"   style="left:72px;  top:246px" />
      <img class="avatar" src="./bender.png" style="left:26px;  top:106px" />
    </div>
    <div class="cap good">Parallel <b>consult</b>, one <b>fuse</b> back into the arbiter's context. Expected to trim the low tail.</div>
  </div>

  <div class="panel">
    <div class="topo-head bad">A DAG OF BOTS · SERIAL + FAN-OUT</div>
    <div class="dag">
      <svg viewBox="0 0 440 300" width="440" height="300">
        <defs>
          <marker id="rtip" markerWidth="7" markerHeight="7" refX="5" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="#dc2626"/></marker>
        </defs>
        <g stroke="#dc2626" stroke-width="2" marker-end="url(#rtip)" fill="none">
          <line x1="62" y1="138" x2="120" y2="95"/>
          <line x1="62" y1="162" x2="120" y2="205"/>
          <line x1="157" y1="92"  x2="215" y2="135"/>
          <line x1="157" y1="208" x2="215" y2="165"/>
          <line x1="252" y1="138" x2="310" y2="95"/>
          <line x1="252" y1="162" x2="310" y2="205"/>
          <line x1="346" y1="93"  x2="396" y2="134"/>
          <line x1="346" y1="207" x2="396" y2="166"/>
        </g>
      </svg>
      <img class="avatar" src="./hk47.png"   style="left:24px;  top:129px" />
      <img class="avatar" src="./K-2SO.png"  style="left:119px; top:59px" />
      <img class="avatar" src="./glados.png" style="left:119px; top:199px" />
      <img class="avatar" src="./c3po.png"   style="left:214px; top:129px" />
      <img class="avatar" src="./bender.png" style="left:309px; top:59px" />
      <img class="avatar" src="./hk47.png"   style="left:309px; top:199px" />
      <img class="avatar" src="./K-2SO.png"  style="left:394px; top:129px" />
    </div>
    <div class="cap bad">Each hop sees prior output. The test measures <b>regressions</b>, order effects, and dropped findings.</div>
  </div>

</div>

<style>
.topo-wrap { display:flex; gap:30px; justify-content:center; align-items:flex-start; margin-top:6px; }
.panel { display:flex; flex-direction:column; align-items:center; }
.topo-head { font-size:12px; letter-spacing:.12em; font-weight:700; margin-bottom:2px; }
.topo-head.good { color:#047857; } .topo-head.bad { color:#b91c1c; }
.flower { position:relative; width:340px; height:340px; }
.dag { position:relative; width:440px; height:300px; }
.avatar { position:absolute; border-radius:9999px; object-fit:cover; background:#fff; }
.flower .avatar { width:50px; height:50px; border:2px solid #10b981; box-shadow:0 2px 6px rgba(0,0,0,.2); }
.dag .avatar { width:42px; height:42px; border:2px solid #dc2626; box-shadow:0 2px 6px rgba(0,0,0,.2); }
.arbiter { position:absolute; left:131px; top:131px; width:78px; height:78px; border-radius:50%;
  background:#06281f; border:2px solid #10b981; color:#a7f3d0; box-shadow:0 0 18px rgba(16,185,129,.45);
  display:flex; flex-direction:column; align-items:center; justify-content:center; font-family:ui-monospace,monospace; }
.arbiter b { font-size:12px; line-height:1.1; margin-top:1px; } .arbiter small { font-size:8px; opacity:.75; }
.cap { font-size:12px; line-height:1.4; text-align:center; max-width:340px; margin-top:8px; color:#475569; }
.cap.good b { color:#047857; } .cap.bad b { color:#b91c1c; }
</style>

<!--
The topology IS the argument. Left: the main session is the hub — it fans out to the
council in parallel and folds every reply back into ONE context (lift the start,
narrow the cone). Right: a serial DAG hands output bot-to-bot, so every hop reseeds
the cone — the failure mode from the previous two slides. Same bots, opposite variance.
-->

---
layout: center
class: text-center
---

# Fuse, don't chain.

<div class="mx-auto max-w-3xl text-left mt-1 space-y-2 text-[15px] leading-relaxed">

- **Hypothesis:** specialist roles add true findings beyond repeated generic reviews.
- **Question:** do pop-culture wrappers help, or only spend context and style tokens?
- **Topology test:** fan-out + fusion versus every ordering of a three-reviewer chain.

</div>

<div class="faces">
  <img src="./hk47.png" /><img src="./K-2SO.png" /><img src="./glados.png" /><img src="./c3po.png" /><img src="./bender.png" /><img src="./holly.png" /><img src="./walle.png" />
</div>

<div class="text-sm opacity-70 mt-3">Next — the blinded, planted-defect protocol in <code>../part-2/slides.md</code>.</div>
<div class="text-xs opacity-50 mt-5">Convene the council. Questions?</div>

<style>
.faces { display:flex; gap:10px; justify-content:center; margin-top:22px; }
.faces img { width:52px; height:52px; border-radius:9999px; object-fit:cover; border:2px solid #cbd5e1; box-shadow:0 2px 6px rgba(0,0,0,.15); }
</style>
