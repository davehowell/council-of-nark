---
theme: seriph
title: Put the Council on Trial
colorSchema: light
info: |
  Follow-up to The Council of Nark: a blinded, planted-defect experiment for role diversity,
  pop-culture personas, fusion, provider effects, and fan-out versus serial chains.
transition: slide-left
mdc: true
layout: image-right
image: /Nark-council.png
backgroundSize: cover
---

# Put the Council on Trial

### A falsifiable test of roles, robots and topology

<div class="mt-8 text-sm opacity-75 max-w-md">
Same model. Frozen artifacts. Hidden defects. Blinded scoring.
</div>

<!-- Presenter cues
- Open with the tension: the first presentation introduced the Council of Nark, and then made an appealing claim; In part 2 we put give the council a way to fail.
- “Same model” isolates prompting. “Hidden defects” gives us a known answer. “Blinded scoring” keeps the grader honest.
- The question is simple: which setup finds more real problems for its cost?
- Transition: split the Council idea into three claims so one result cannot hide another.
-->

---
layout: default
---

# Three claims — not one

<div class="grid grid-cols-3 gap-6 mt-7">
  <div class="claim"><b>1 · Roles</b><p>Independent specialists find more real defects than repeated general reviews.</p></div>
  <div class="claim"><b>2 · Personas</b><p>Character wrappers add value beyond the explicit specialist checklist.</p></div>
  <div class="claim"><b>3 · Topology</b><p>Fan-out + one fusion retains more signal than serial hand-offs.</p></div>
</div>

<div class="mt-7 px-5 py-3 rounded-xl bg-amber-50 border border-amber-300 text-[13px] leading-relaxed">
  <b>Black-box limit:</b> we measure repeatable shifts in defect coverage, errors, variance and cost.
  Role cues may change hidden-state and attention trajectories, and which learned features contribute through fixed weights; this changes logits and the conditional token distribution; decoding then adds sample variance.
  API outputs cannot isolate those stages or prove that different weights were “activated.”
</div>

<style>
.claim { border:1px solid #cbd5e1; border-radius:14px; padding:18px 20px; background:#f8fafc; min-height:145px; }
.claim b { color:#7c3aed; font-size:19px; } .claim p { font-size:15px; line-height:1.45; margin-top:10px; }
</style>

<!-- Presenter cues
- Separate the claims: functional roles, fictional wrappers, and orchestration topology can each win or lose independently.
- Technical caveat: the model parameters stay fixed. “Latent space” here is shorthand for different hidden-state and attention trajectories, not an observed mechanism.
- The pipeline is: prompt changes internal activations → logits/token probabilities change → decoding selects one output. Sampling adds noise at the final step.
- We only observe the endpoint. The useful claim is repeatable improvement in coverage, reliability or cost, not model transparency.
- Transition: to score that endpoint, give every condition the same traps with a private answer key.
-->

---
layout: default
---

# Three synthetic traps

<div class="grid grid-cols-3 gap-6 mt-7">
  <div class="packet green"><span>01</span><b>Revenue dashboard</b><small>SQL · dbt · Cube · recovery</small><p>Fan-out, full scans, PII, contract bypass, speculative factory, forgotten flag, vague safety.</p></div>
  <div class="packet amber"><span>02</span><b>Key rotation at 03:00</b><small>runbook · deployment note</small><p>Leaked key, outage window, wrong context, <code>latest</code>, muted alerts, no rollback, four-service factory.</p></div>
  <div class="packet blue"><span>03</span><b>Webhook redesign</b><small>architecture proposal</small><p>Race-prone dedupe, slow ACK, payload logging, six queues, schema drift, silent backlog, lossy pause.</p></div>
</div>

<div class="text-center mt-7 text-sm opacity-70"><b>Packet</b> = one frozen review scenario and its context · 8 planted defects · private answer key</div>

<style>
.packet { position:relative; border:1px solid #cbd5e1; border-radius:14px; padding:22px 18px 18px; min-height:245px; }
.packet span { position:absolute; right:14px; top:10px; font-size:38px; font-weight:800; opacity:.12; }
.packet b { display:block; font-size:18px; margin-top:12px; }.packet small { display:block; opacity:.6; margin-top:5px; }
.packet p { font-size:14px; line-height:1.5; margin-top:22px; }.packet.green { background:#ecfdf5; }.packet.amber { background:#fffbeb; }.packet.blue { background:#eff6ff; }
</style>

<!-- Presenter cues
- Define “packet” now: the frozen artifact plus the facts needed to review it. Every arm sees the same packet; only the review setup changes.
- Each packet plants eight material defects across several lenses and also includes correct details, so restraint is measurable.
- The answer key never enters model context. It maps semantically equivalent findings to defect IDs.
- Three packets are a mechanism pilot, not a universal benchmark. Confirmation needs frozen variants and clean controls.
- Transition: first hold the call count at one and vary only the prompt wrapper.
-->

---
layout: default
---

# First: one reviewer

<div class="arm-list mt-8">
  <div><code>S0</code><b>Plain generic review</b><span>Weak baseline</span></div>
  <div><code>S1</code><b>Functional omnibus review</b><span>All seven lenses listed fairly · no persona</span></div>
  <div><code>S2</code><b>Functional omnibus + GLaDOS</b><span>Length-matched character layer</span></div>
</div>

<div class="grid grid-cols-2 gap-6 mt-9 text-[15px]">
  <div class="test"><b>Single-call persona effect</b><code>S2 − S1</code></div>
  <div class="test"><b>Fair baseline</b><span>S1, never S0 alone</span></div>
</div>

<div class="text-center mt-5 text-xs opacity-60"><b>Arm</b> = one prompt and orchestration condition tested against every packet</div>

<style>
.arm-list { display:flex; flex-direction:column; gap:12px; }.arm-list>div { display:grid; grid-template-columns:62px 310px 1fr; align-items:center; padding:14px 18px; border:1px solid #cbd5e1; border-radius:12px; background:#f8fafc; }
.arm-list code { color:#7c3aed; font-size:18px; font-weight:700; }.arm-list b { font-size:16px; }.arm-list span { font-size:14px; opacity:.7; }
.test { border-left:4px solid #8b5cf6; background:#f5f3ff; padding:14px 18px; border-radius:8px; display:flex; justify-content:space-between; align-items:center; }.test code { font-size:18px; }
</style>

<!-- Presenter cues
- Define “arm”: one experimental condition. Each arm processes all three packets.
- S0 is the ordinary “review this” prompt. It is context, not the control we want to beat.
- S1 is the fair control: one reviewer explicitly receives every functional lens.
- S2 keeps that checklist and adds only the GLaDOS wrapper. S2 minus S1 isolates the single-call persona increment.
- Say this plainly: if the council only beats S0, we learned almost nothing.
- Transition: panels spend more calls, so match that compute before crediting specialisation.
-->

---
layout: default
---

# Then match the calls

<div class="grid grid-cols-3 gap-5 mt-7">
  <div class="multi"><code>M0</code><b>Repeated functional omnibus</b><p>7 independent copies of S1</p><small>Controls for more samples</small></div>
  <div class="multi focus"><code>M1</code><b>Functional specialists</b><p>7 explicit lenses, no fiction</p><small>Tests role specialisation</small></div>
  <div class="multi"><code>M2</code><b>Robot specialists</b><p>Same 7 kernels + character wrappers</p><small>Tests pop-culture increment</small></div>
</div>

<div class="fuse-row mt-8"><span>seven fresh review calls</span><b>→</b><span class="arb">same model · same fuser</span><b>→</b><span>ranked verdict</span></div>

<div class="mt-7 text-center text-sm opacity-70">Every arm uses one pinned cheap model first. Same artifact, schema, caps and fresh sessions.</div>

<style>
.multi { border:1px solid #cbd5e1; border-radius:14px; padding:22px 18px; text-align:center; min-height:185px; background:#f8fafc; }.multi.focus { border-color:#8b5cf6; background:#f5f3ff; }.multi code { display:block; font-size:22px; color:#7c3aed; font-weight:800; }.multi b { display:block; font-size:18px; margin:8px 0; }.multi p { font-size:14px; }.multi small { display:block; opacity:.6; margin-top:15px; }
.fuse-row { display:flex; justify-content:center; align-items:center; gap:14px; font-size:14px; }.fuse-row span { border:1px solid #94a3b8; border-radius:9px; padding:10px 16px; }.fuse-row .arb { border-color:#8b5cf6; background:#f5f3ff; }
</style>

<!-- Presenter cues
- M0 buys seven independent lottery tickets with the same functional omnibus prompt.
- M1 spends the same seven calls on separate functional specialists. M2 adds the robot wrappers to those same kernels.
- All three use a separate, persona-free fuser, so the fuser is not the treatment.
- The full production persona files are too unequal for this controlled comparison; use distilled, length-matched kernels first.
- Transition: these matched arms let each subtraction answer one narrow question.
-->

---
layout: default
---

# Each subtraction answers one question

<div class="contrast-layout mt-7">
  <div class="equations">
    <div><code>M1 − M0</code><span><b>Specialisation</b> beyond repeated general sampling</span></div>
    <div><code>M2 − M1</code><span><b>Pop-culture wrapper</b> beyond the same role instructions</span></div>
    <div><code>M0 − S1</code><span><b>More samples</b> without specialist diversity</span></div>
    <div><code>M2 − S1</code><span><b>Practical council uplift</b> — with cost reported</span></div>
  </div>
  <aside class="arm-key">
    <b>Key</b>
    <div><code>S1</code><span>Functional omnibus<br><small>1 call</small></span></div>
    <div><code>M0</code><span>Repeated functional omnibus<br><small>7 reviewers + fuser</small></span></div>
    <div><code>M1</code><span>Functional specialists<br><small>7 reviewers + fuser</small></span></div>
    <div><code>M2</code><span>Robot specialists<br><small>7 reviewers + fuser</small></span></div>
  </aside>
</div>

<div class="mt-6 px-5 py-3 rounded-xl bg-red-50 border border-red-200 text-sm">
If M2 only produces more colourful prose, it loses. Count additional <b>true defect IDs</b>, not personality.
</div>

<style>
.contrast-layout { display:grid; grid-template-columns:1fr 260px; gap:30px; align-items:start; }
.equations { display:flex; flex-direction:column; gap:8px; }.equations>div { display:grid; grid-template-columns:175px 1fr; align-items:center; border-bottom:1px solid #e2e8f0; padding:10px 10px; }.equations code { color:#7c3aed; font-weight:800; font-size:21px; }.equations span { font-size:15px; }
.arm-key { border:1px solid #cbd5e1; border-radius:12px; background:#f8fafc; padding:13px 15px; }.arm-key>b { display:block; font-size:13px; letter-spacing:.12em; text-transform:uppercase; opacity:.55; margin-bottom:6px; }.arm-key>div { display:grid; grid-template-columns:38px 1fr; gap:7px; padding:6px 0; border-top:1px solid #e2e8f0; }.arm-key code { color:#7c3aed; font-weight:800; }.arm-key span { font-size:12px; line-height:1.25; }.arm-key small { opacity:.6; font-size:10px; }
</style>

<!-- Presenter cues
- Read these as controlled contrasts, not literal arithmetic on one output.
- M1 minus M0 asks whether dividing attention by function beats seven general samples.
- M2 minus M1 is the clean test of the pop-culture layer; the specialist instructions are otherwise identical.
- M0 minus S1 prices the benefit of sampling alone. M2 minus S1 is the practical product comparison, with its extra calls exposed.
- Without M0, “the council won” may mean only more lottery tickets. Without M1, “personas won” may mean only a better checklist.
- Transition: define “win” in terms of signal, noise, reliability and cost.
-->

---
layout: default
---

# Score signal, noise and cost

<div class="grid grid-cols-2 gap-8 mt-8">
<div>
  <h3 class="text-emerald-700">Primary</h3>
  <div class="metric big">Macro F1 on planted defects</div>
  <h3 class="mt-6 text-slate-600">Cone claim</h3>
  <div class="metric">p10 · worst case · variance · IQR</div>
</div>
<div>
  <h3 class="text-blue-700">Diversity</h3>
  <div class="metric">Jaccard overlap · unique true findings · category coverage</div>
  <h3 class="mt-6 text-slate-600">Efficiency</h3>
  <div class="metric">true findings / 1k tokens · $ · latency</div>
</div>
</div>

<div class="mt-9 text-sm text-center opacity-75">A unique false positive is noise. Duplicate true findings are consensus. Different wording is neither.</div>

<style>
.metric { border:1px solid #cbd5e1; border-radius:12px; background:#f8fafc; padding:18px; margin-top:8px; font-size:16px; }.metric.big { font-size:22px; font-weight:700; border-color:#34d399; background:#ecfdf5; }
</style>

<!-- Presenter cues
- Primary score is macro F1: recall rewards finding planted defects; precision penalises confident nonsense.
- The cone claim predicts a better floor and narrower spread, so show p10 and variance, not only the mean.
- Jaccard overlap is computed on defect IDs, not wording. A new true ID is useful diversity; a novel hallucination is noise.
- Normalise by tokens, dollars and latency because a seven-reviewer win may still be a poor product choice.
- Scoring is blind to arm, provider and catchphrases; ambiguous semantic matches go to two human raters.
- Transition: score the panel before and after fusion to see whether the arbiter preserves its advantage.
-->

---
layout: default
---

# Did the arbiter help — or drop the minority report?

<div class="grid grid-cols-2 gap-10 mt-9">
  <div class="fusion raw"><b>Raw union</b><p>All true and false findings from seven reviewers</p><span>maximum coverage<br>maximum noise</span></div>
  <div class="fusion final"><b>Fused verdict</b><p>Same fuser de-duplicates, ranks and rules</p><span>precision gain?<br>fusion loss?</span></div>
</div>

<div class="mt-10 text-center text-[17px]">
<b>Fusion retention</b> = fused true positives ÷ raw-union true positives
</div>

<div class="mt-5 text-center text-sm opacity-70">Score both outputs. A neat summary that drops valid edge cases weakens the arbiter claim.</div>

<style>
.fusion { border-radius:16px; padding:28px; min-height:210px; text-align:center; }.fusion.raw { background:#f1f5f9; border:1px dashed #64748b; }.fusion.final { background:#f5f3ff; border:2px solid #8b5cf6; }.fusion b { font-size:24px; }.fusion p { font-size:15px; margin:18px 0; }.fusion span { font-size:14px; opacity:.65; }
</style>

<!-- Presenter cues
- Treat the raw union as the panel’s maximum available evidence, including its noise.
- The fuser should merge duplicates and reject unsupported claims without deleting a valid minority report.
- Fusion retention makes that loss visible. A polished short verdict can score worse than the untidy panel.
- This separates “more reviewers helped” from “the arbiter helped.”
- Transition: now keep the reviewers and fuser fixed, and change only how information flows between them.
-->

---
layout: default
---

# Fan-out versus every chain order

<div class="grid grid-cols-2 gap-8 mt-4">
<div>
  <div class="top-label good">FAN-OUT</div>
  <div class="diagram"><span>A</span><span>B</span><span>C</span><b>↓ &nbsp; ↓ &nbsp; ↓</b><em>same fuser</em></div>
  <p class="caption">Each reviewer sees the original packet independently.</p>
</div>
<div>
  <div class="top-label bad">INFORMED CHAIN</div>
  <div class="diagram chain"><span>A</span><b>→</b><span>B + ledger</span><b>→</b><span>C + ledger</span><em>↓ same fuser</em></div>
  <p class="caption">Every reviewer still sees the original. Prior findings can anchor, help, or vanish.</p>
</div>
</div>

<div class="orders mt-7"><b>3 roles → 6 orders</b><code>ABC · ACB · BAC · BCA · CAB · CBA</code></div>
<div class="mt-5 text-center text-sm opacity-70">Measure true findings added, retained and dropped at each hop — plus final F1 and order sensitivity.</div>

<style>
.top-label { text-align:center; font-size:13px; font-weight:800; letter-spacing:.16em; }.top-label.good { color:#047857; }.top-label.bad { color:#b91c1c; }.diagram { margin-top:9px; height:150px; border:1px solid #a7f3d0; background:#ecfdf5; border-radius:14px; display:flex; align-items:center; justify-content:center; gap:16px; position:relative; }.diagram.chain { border-color:#fecaca; background:#fef2f2; gap:9px; }.diagram span,.diagram em { border:1px solid #94a3b8; border-radius:9999px; padding:9px 13px; background:white; font-style:normal; font-size:13px; }.diagram em { position:absolute; bottom:12px; }.caption { font-size:13px; opacity:.7; text-align:center; margin-top:9px; }.orders { display:flex; justify-content:center; gap:30px; align-items:center; }.orders code { font-size:17px; color:#7c3aed; }
</style>

<!-- Presenter cues
- This is a fair informed chain: B and C still see the original packet. The test is not a straw-man summary relay.
- The chain also sees prior findings, which may help, anchor later reviewers, or cause valid findings to disappear.
- Three roles give six possible orders; run all six because a chain that wins only in one order is operationally fragile.
- Pair each chain order with the same fuser input order in fan-out to control recency effects.
- Track findings added, retained and dropped at each hop, and record the chain’s extra input tokens.
- Transition: after isolating prompting and topology within one model, ask whether the effect survives another provider.
-->

---
layout: default
---

# Then cross the provider boundary

<div class="providers mt-12">
  <div><b>Claude</b><span>Haiku</span></div><i>×</i><div><b>Gemini</b><span>Flash / Flash-Lite</span></div><i>×</i><div><b>OpenAI</b><span>low-cost GPT</span></div>
</div>

<div class="mt-10 mx-auto max-w-3xl text-[16px] leading-relaxed">
Send the <b>byte-identical K-2SO prompt</b> and its persona-free twin to each pinned model. Test <code>provider × persona</code> before paying to replicate the whole council.
</div>

<div class="mt-8 text-sm opacity-70 text-center">Same text ≠ same tokens or hidden system prompt. Provider is a factor, not a verdict.</div>

<style>
.providers { display:flex; justify-content:center; align-items:center; gap:24px; }.providers div { width:210px; border:1px solid #cbd5e1; border-radius:14px; padding:22px; text-align:center; background:#f8fafc; }.providers b { display:block; font-size:22px; }.providers span { display:block; font-size:13px; opacity:.6; margin-top:7px; }.providers i { font-size:24px; color:#7c3aed; font-style:normal; }
</style>

<!-- Presenter cues
- Start cheaply with one paired prompt: K-2SO and the identical correctness kernel without the character wrapper.
- Send byte-identical text, but do not call it a perfectly identical treatment: tokenisers and hidden system prompts differ.
- Provider × persona tells us whether the wrapper effect is robust or model-family-specific.
- If there is no Stage A effect, stop; do not pay to replicate a null across the whole council.
- Transition: close with the staged run plan and the decisions each outcome permits.
-->

---
layout: center
class: text-center
---

# Let the robots lose.

<div class="terms mx-auto max-w-4xl mt-5">
  <span><b>Packet</b> one frozen review scenario</span>
  <span><b>Arm</b> one prompt/orchestration condition: S0–M2</span>
</div>

<div class="mx-auto max-w-3xl text-left mt-6 space-y-3 text-[15px]">
  <div><b>Smoke:</b> run every arm once on every packet; repair prompts and ceiling/floor effects.</div>
  <div><b>Pilot:</b> 10 samples per packet × arm; 3 samples per topology order.</div>
  <div><b>Confirm:</b> freeze packet variants and thresholds; score blind against private answer keys.</div>
</div>

<div class="mt-8 grid grid-cols-3 gap-4 text-sm">
  <div class="outcome yes"><b>Roles win</b><br>M1 &gt; M0</div>
  <div class="outcome maybe"><b>Characters earn rent</b><br>M2 &gt; M1 per token</div>
  <div class="outcome no"><b>Or we delete the bit</b><br>nulls count</div>
</div>

<div class="mt-8 text-xs opacity-55"><code>experiment/protocol.md</code> · frozen scenarios under <code>experiment/scenarios/</code></div>

<style>
.terms { display:flex; justify-content:center; gap:12px; }.terms span { border:1px solid #cbd5e1; border-radius:9999px; padding:8px 15px; background:#f8fafc; font-size:12px; }.terms b { color:#7c3aed; margin-right:5px; }
.outcome { border-radius:12px; padding:15px; }.outcome.yes { background:#ecfdf5; }.outcome.maybe { background:#fffbeb; }.outcome.no { background:#fef2f2; }
</style>

<!-- Presenter cues
- Re-state the terms: a packet is one frozen scenario; an arm is one competing way to review it.
- Smoke is plumbing and calibration, not evidence. Pilot estimates effects and variance. Confirm is the first preregistered claim-bearing run.
- Explain the title: falsifiability means giving the favourite idea permission to lose before seeing results.
- Decision rules: keep functional roles if M1 beats M0; keep character wrappers only if M2 beats M1 per token; prefer fan-out only if it beats the chain orders reliably.
- A null is useful product evidence. It tells us to remove context, latency and theatre that do not improve the endpoint.
- Closing line: “The goal is not to prove the Council is clever. It is to learn which parts earn their place.”
-->
