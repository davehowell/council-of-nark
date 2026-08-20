---
theme: seriph
title: The Experiment Fought Back
colorSchema: light
info: |
  Part 3 of the Council of Nark series: instrumentation failures, ceiling effects,
  parser bugs, contamination audits, negative persona results, and a stricter Go/Seatbelt harness.
transition: slide-left
mdc: true
layout: image-right
image: ./Nark-council.png
backgroundSize: cover
---

# The Experiment Fought Back

### Trials, failures and useful negative results

<div class="mt-8 text-sm opacity-70 max-w-md">
The engineering story between “let's test it” and “we can trust the test.”
</div>

<img src="./Nark-council.png" class="absolute right-0 top-0 w-1/2 h-full object-cover" />

<!-- Presenter cues
- Part 1 proposed the council. Part 2 designed a falsifiable trial.
- Part 3 is the missing middle: the experiment itself became the system under review.
- This is not a victory lap. The failures changed both the harness and our interpretation.
-->

---
layout: default
---

# The plan looked linear

<div class="linear mt-12">
  <span>freeze</span><b>→</b><span>run</span><b>→</b><span>seal</span><b>→</b><span>blind</span><b>→</b><span>score</span><b>→</b><span>conclude</span>
</div>

<div class="reality mt-12">
  <b>Reality</b>
  <code>fail → inspect → preserve → repair → rerun → doubt → audit → narrow the claim</code>
</div>

<style>
.linear { display:flex; justify-content:center; align-items:center; gap:14px; }.linear span { border:1px solid #cbd5e1; background:#f8fafc; border-radius:12px; padding:16px 20px; font-weight:700; }.linear b { color:#7c3aed; }
.reality { border:2px solid #f59e0b; background:#fffbeb; border-radius:16px; padding:24px; text-align:center; }.reality b { display:block; color:#b45309; margin-bottom:10px; }.reality code { font-size:16px; }
</style>

<!-- Presenter cues
- Reproducibility is not just replaying the final happy path.
- We kept failed and incomplete runs because they explain why later controls exist.
- The discipline: never patch an observed sample into validity. Seal it, exclude it, start again.
-->

---
layout: default
---

# Five runs and a harness rewrite

<div class="timeline mt-6">
  <div><time>04:23</time><b>Schema failure</b><span>0 observations</span></div>
  <div><time>04:27</time><b>Haiku smoke</b><span>81 / 81</span></div>
  <div><time>08:15</time><b>Gemma partial</b><span>discarded</span></div>
  <div><time>08:31</time><b>Gemma clean</b><span>81 / 81</span></div>
  <div><time>09:00</time><b>Persona pairs</b><span>60 / 60</span></div>
  <div><time>12:23</time><b>Go + Seatbelt</b><span>fail closed</span></div>
</div>

<div class="mt-7 text-center text-sm opacity-70">Every arrow represents a changed threat model, not just another command.</div>

<style>
.timeline { display:grid; grid-template-columns:repeat(3,1fr); gap:14px; }.timeline div { border:1px solid #cbd5e1; border-radius:13px; padding:16px; background:#f8fafc; }.timeline time { color:#7c3aed; font-family:monospace; font-weight:800; }.timeline b,.timeline span { display:block; }.timeline b { margin-top:7px; }.timeline span { font-size:12px; opacity:.6; margin-top:4px; }
</style>

<!-- Presenter cues
- All times are UTC from the append-only lab notebook.
- We moved from “does the council work?” to a sequence of smaller questions: did inference happen, was output parsed, was scoring fair, was context isolated?
- The timeline is itself an experimental artifact.
-->

---
layout: default
---

# Failure 1: no model was called

<div class="grid grid-cols-2 gap-8 mt-10">
  <div class="failure"><b>Expected</b><p>81 structured Claude responses</p></div>
  <div class="failure red"><b>Observed</b><p>72 local schema errors<br>9 blocked fusers<br>0 inference responses</p></div>
</div>

<div class="fix mt-9"><code>$schema</code> remained in source and provenance, but was omitted from the CLI argument.</div>

<div class="mt-5 text-sm opacity-70 text-center">The failed run stayed sealed. It never became data.</div>

<style>
.failure { border:1px solid #cbd5e1; border-radius:14px; padding:24px; background:#f8fafc; }.failure.red { background:#fef2f2; border-color:#fca5a5; }.failure b { font-size:20px; }.failure p { margin-top:14px; line-height:1.7; }.fix { border-left:5px solid #10b981; background:#ecfdf5; border-radius:9px; padding:17px 22px; text-align:center; }
</style>

<!-- Presenter cues
- Standards-valid JSON Schema was not accepted by that CLI validator.
- The important distinction: instrumentation failure versus a malformed model answer.
- We added a frozen one-call adapter check before every expensive batch.
-->

---
layout: default
---

# Success created a new problem

<div class="big-number mt-6">0.81–0.86</div>
<div class="text-center text-lg">Mean F1 for most final Haiku conditions</div>

<div class="grid grid-cols-3 gap-5 mt-10">
  <div class="lesson"><b>Plumbing</b><span>worked</span></div>
  <div class="lesson"><b>Headroom</b><span>did not</span></div>
  <div class="lesson"><b>Verdict</b><span>not available</span></div>
</div>

<div class="mt-8 px-5 py-3 rounded-xl bg-amber-50 border border-amber-300 text-sm text-center">
A frontier-ish model on easy synthetic packets mostly measured ceiling effects and one decoding draw.
</div>

<style>
.big-number { text-align:center; font-size:70px; font-weight:900; color:#7c3aed; }.lesson { text-align:center; border:1px solid #cbd5e1; border-radius:13px; padding:20px; }.lesson b,.lesson span { display:block; }.lesson span { margin-top:7px; opacity:.6; }
</style>

<!-- Presenter cues
- The first complete smoke proved scheduling, sealing, blinding and scoring.
- It did not prove that the council was good or useless.
- The correct pivot was a cheaper explicit model on unchanged scenarios—not post-hoc harder traps designed to rescue the favourite idea.
-->

---
layout: default
---

# Metric key — because nobody remembers F1

| Key | Meaning |
|---|---|
| **TP** | one unique planted defect correctly found |
| **FP** | one unique unsupported claim after semantic de-duplication |
| **FN** | one planted defect missed |
| **Precision** | `TP / (TP + FP)` — how much of the review was supported |
| **Recall** | `TP / (TP + FN)` — how much of the defect set was found |
| **F1** | `2TP / (2TP + FP + FN)` — balance of precision and recall |

<div class="mt-7 text-center text-sm opacity-70">Equal F1 means equal counts. It does not mean equal defects, remedies or severity.</div>

<!-- Presenter cues
- Eight defects per packet means recall moves in coarse steps of 0.125.
- Macro F1 scores each output set and then weights sets equally.
- F1 cannot judge whether a proposed fix is sensible.
-->

---
layout: default
---

# Equal scores hid different reviews

<div class="compare mt-10">
  <div><code>S1</code><b>found RD-04</b><span>missed RD-06</span></div>
  <em>same F1</em>
  <div><code>M0</code><b>found RD-06</b><span>missed RD-04</span></div>
</div>

<div class="mt-12 grid grid-cols-2 gap-6 text-sm">
  <div class="note"><b>Count metric</b><br>Same TP/FP/FN balance</div>
  <div class="note"><b>Coverage metric</b><br>Different defect identity and contribution</div>
</div>

<style>
.compare { display:grid; grid-template-columns:1fr 130px 1fr; align-items:center; gap:20px; }.compare>div { border:1px solid #cbd5e1; border-radius:15px; padding:28px; text-align:center; background:#f8fafc; }.compare code { color:#7c3aed; font-size:24px; font-weight:800; }.compare b,.compare span { display:block; margin-top:10px; }.compare span { opacity:.6; }.compare em { text-align:center; color:#b45309; font-style:normal; font-weight:800; }.note { background:#f5f3ff; border-radius:10px; padding:17px; }
</style>

<!-- Presenter cues
- This forced overlap reporting: detected IDs, Jaccard similarity and one-sided contributions alongside F1.
- It also changed how we interpret a council: broader raw coverage can exist even when a fuser makes final scores converge.
-->

---
layout: default
---

# The audit found four quiet biases

<div class="grid grid-cols-2 gap-5 mt-7">
  <div class="bug"><b>Order</b><p>A set discarded the seeded schedule.</p></div>
  <div class="bug"><b>False positives</b><p>Seven repeats of one bad claim counted seven times.</p></div>
  <div class="bug"><b>Fusion</b><p>A capable arbiter flattened panel differences.</p></div>
  <div class="bug"><b>Sampling</b><p>One output could not identify a prompt effect.</p></div>
</div>

<div class="mt-7 text-center text-sm"><b>Repairs:</b> preserve order · cluster false claims · publish raw unions · run paired repeats</div>

<style>
.bug { border-left:5px solid #ef4444; background:#fef2f2; border-radius:9px; padding:17px 20px; }.bug b { color:#991b1b; }.bug p { font-size:14px; margin-top:7px; }
</style>

<!-- Presenter cues
- None of these leaked the answer key, but each could distort the measured effect.
- A post-smoke contamination audit should be expected, not embarrassing.
- The raw union became primary evidence for role diversity; fused F1 became the practical final-verdict measure.
-->

---
layout: default
---

# Failure 2: the prompt looked like an answer

<div class="flow mt-8">
  <span>Pi JSON stream</span><b>→</b><span>echoed user prompt</span><b>→</b><span>provider 429</span><b>→</b><span class="bad">example JSON parsed</span>
</div>

<div class="grid grid-cols-2 gap-7 mt-10">
  <div class="failure red"><b>Partial run</b><p>71 success · 7 quota failures · 3 blocked fusers</p></div>
  <div class="failure"><b>Repair</b><p>assistant events only · provider errors promoted · retries + backoff</p></div>
</div>

<div class="mt-6 text-center text-sm opacity-70">Sealed. Discarded before rating. Replaced from a new seed.</div>

<style>
.flow { display:flex; align-items:center; justify-content:center; gap:10px; }.flow span { border:1px solid #cbd5e1; border-radius:10px; padding:13px; font-size:13px; }.flow .bad { border-color:#ef4444; background:#fef2f2; }.flow b { color:#7c3aed; }
.failure { border:1px solid #cbd5e1; border-radius:14px; padding:24px; background:#f8fafc; }.failure.red { background:#fef2f2; border-color:#fca5a5; }.failure b { font-size:20px; }.failure p { margin-top:14px; line-height:1.55; }
</style>

<!-- Presenter cues
- The local process exited zero even though the provider failed.
- Searching the whole stream found the schema example inside the echoed prompt.
- The parser now considers only assistant-role text and treats assistant error events as infrastructure failures.
-->

---
layout: default
---

# The rater broke too

<div class="grid grid-cols-2 gap-7 mt-8">
  <div class="bug"><b>Schema dialect</b><p>Vertex required explicit types beside enums.</p></div>
  <div class="bug"><b>Wrong JSON object</b><p>The extractor selected the embedded schema's <code>judgements</code> property.</p></div>
</div>

<div class="fix mt-9">Accept the requested root only when its value is an array.</div>

<div class="mt-6 text-center text-sm opacity-70">
419 captured, schema-valid ratings were recovered mechanically.<br>No replacement model call. No changed judgement.
</div>

<style>
.fix { border-left:5px solid #10b981; background:#ecfdf5; border-radius:9px; padding:17px 22px; text-align:center; }
</style>

<!-- Presenter cues
- Derived analysis stages need provenance and tests just as much as respondent stages.
- Recovery was acceptable because the model output already existed and validated; only parser selection had failed.
- Later, rating stages require a clean committed harness so uncommitted parser code cannot become invisible provenance.
-->

---
layout: default
---

# Cheap model, useful spread

| Arm | Final F1 | Recall |
|---|---:|---:|
| S0 generic | 0.819 | 0.708 |
| M0 repeated omnibus | 0.863 | 0.792 |
| **M1 functional specialists** | **0.911** | **0.875** |
| M2 fictional specialists | 0.857 | 0.792 |

<div class="mt-8 grid grid-cols-2 gap-6 text-sm">
  <div class="note"><b>Before fusion</b><br>M1 recall 0.917 vs M0 0.792</div>
  <div class="note"><b>Still easy</b><br>S0 remained above 0.81</div>
</div>

<div class="mt-6 text-center text-sm opacity-70">The council exposed more signal—and more arbitration work.</div>

<style>
.note { background:#f5f3ff; border:1px solid #ddd6fe; border-radius:10px; padding:17px; }
</style>

<!-- Presenter cues
- This is the behaviour we expected to test: specialists increased raw coverage.
- But the three packets remained easy even for Gemma, so this is calibration rather than general proof.
- M2 reversed the first smoke's small persona advantage. One draw was never enough.
-->

---
layout: default
---

# The fictional overlay lost this round

<div class="grid grid-cols-2 gap-8 mt-5">
  <div class="score good"><small>Functional correctness</small><b>0.790</b><span>mean F1</span></div>
  <div class="score bad"><small>K-2SO wrapper</small><b>0.748</b><span>mean F1</span></div>
</div>

<div class="mt-8 text-center text-[17px]"><code>fictional − functional = −0.0425</code></div>
<div class="mt-3 text-center text-sm opacity-70">2 fictional wins · 16 ties · 12 functional wins</div>

<div class="mt-7 px-5 py-3 rounded-xl bg-blue-50 border border-blue-200 text-sm text-center">
Names can remain mnemonic labels. The prose still has to earn its tokens.
</div>

<style>
.score { border-radius:16px; padding:24px; text-align:center; }.score.good { background:#ecfdf5; border:1px solid #6ee7b7; }.score.bad { background:#fef2f2; border:1px solid #fca5a5; }.score small,.score span { display:block; }.score b { display:block; font-size:50px; margin:8px; }
</style>

<!-- Presenter cues
- Thirty packet-blocked pairs used the same model, same kernel and ten repeats.
- The overlay changed defect identity in 14 pairs, so wording affected behaviour; the shift was not beneficial on average.
- This is one persona on toy tasks, not a universal anti-persona claim.
-->

---
layout: default
---

# Isolation had to become executable

<div class="grid grid-cols-2 gap-7 mt-7">
<div>
  <h3>Before</h3>
  <ul class="mt-4 space-y-3 text-sm">
    <li>Python controller</li>
    <li>provider child inside clean worktree</li>
    <li>correct flags were the main boundary</li>
  </ul>
</div>
<div>
  <h3 class="text-emerald-700">Now</h3>
  <ul class="mt-4 space-y-3 text-sm">
    <li>standard-library Go harness</li>
    <li>mandatory macOS Seatbelt</li>
    <li>empty cwd + ephemeral HOME</li>
    <li>repository/worktree read must fail</li>
    <li>external CLI entrypoints frozen by digest</li>
  </ul>
</div>
</div>

<div class="mt-8 font-mono text-center text-sm bg-slate-900 text-green-300 rounded-xl p-4">scratch write: allowed · repository read: denied</div>

<!-- Presenter cues
- Prompt assembly still needs the frozen worktree. The provider child does not.
- The harness fails on non-macOS, root, a failed Seatbelt probe, a dirty freeze, or changed runtime digest.
- There is no silent unsandboxed fallback.
-->

---
layout: default
---

# Authentication fought the sandbox

<div class="grid grid-cols-2 gap-7 mt-8">
  <div class="failure red"><b>agy</b><p>Ephemeral HOME triggered an interactive keychain dependency.</p></div>
  <div class="failure red"><b>Direct Claude CLI</b><p>Shared login state disappeared with the real HOME.</p></div>
</div>

<div class="fix mt-8"><b>Do not weaken isolation.</b> Reject both clients before launch; pin Anthropic, Google, OpenAI and Gemma models through sterile Pi.</div>

<div class="mt-5 text-center text-sm opacity-70">Copy auth + model registry only. Never settings, skills, history, sessions or project context.</div>

<style>
.failure { border:1px solid #fca5a5; border-radius:14px; padding:24px; background:#fef2f2; }.failure b { font-size:20px; }.failure p { margin-top:14px; line-height:1.55; }
.fix { border-left:5px solid #10b981; background:#ecfdf5; border-radius:9px; padding:17px 22px; text-align:center; }
</style>

<!-- Presenter cues
- The keychain prompt was cancelled. We did not reset or create a keychain.
- Compatibility is subordinate to the threat model.
- Using one sterile client also reduces client-specific differences in the cross-provider arm.
-->

---
layout: default
---

# Internet is useful — and contaminating

<div class="grid grid-cols-2 gap-8 mt-8">
  <div class="failure"><b>Production goal</b><p>Use every legal tool that gets the right fix quickly.</p></div>
  <div class="failure red"><b>Experimental goal</b><p>Do not find the matching upstream PR and copy its answer.</p></div>
</div>

<div class="mt-9 grid grid-cols-3 gap-4 text-sm text-center">
  <div class="note"><b>Base arm</b><br>frozen source + local search/tests</div>
  <div class="note"><b>Docs arm</b><br>frozen offline corpus</div>
  <div class="note"><b>Internet arm</b><br>declared factor + controlled proxy</div>
</div>

<div class="mt-6 text-center text-xs opacity-60">Provider-side search remains unobservable unless the provider exposes and disables it.</div>

<style>
.failure { border:1px solid #cbd5e1; border-radius:14px; padding:24px; background:#f8fafc; }.failure.red { background:#fef2f2; border-color:#fca5a5; }.failure b { font-size:20px; }.failure p { margin-top:14px; line-height:1.55; }
.note { background:#f5f3ff; border:1px solid #ddd6fe; border-radius:10px; padding:17px; }
</style>

<!-- Presenter cues
- Seatbelt can allow provider transport while model tools stay disabled.
- Future real-code tasks need allowlisted local tools. Unrestricted web access is a separate treatment, not an invisible convenience.
- Export the pre-fix commit with git archive into a new root; deleting .git in place is not enough.
-->

---
layout: default
---

# Next tests — without moving the goalposts

<div class="grid grid-cols-2 gap-7 mt-6">
  <div class="next"><b>Overlay family</b><span>8 role pairs</span><strong>480 calls</strong><p>Fresh correctness replication. ±0.02 practical margin. Report every role.</p></div>
  <div class="next"><b>Ecological tasks</b><span>pre-fix open source</span><strong>history removed</strong><p>Actual patch + tests as evidence, not the only valid answer.</p></div>
</div>

<div class="mt-8 text-sm text-center opacity-70">Two blinded humans before a claim. LLM ratings remain triage.</div>

<style>
.next { border:1px solid #cbd5e1; border-radius:15px; padding:24px; background:#f8fafc; }.next b,.next span,.next strong { display:block; }.next b { font-size:21px; }.next span { margin-top:7px; opacity:.6; }.next strong { color:#7c3aed; font-size:25px; margin-top:16px; }.next p { font-size:13px; margin-top:14px; line-height:1.5; }
</style>

<!-- Presenter cues
- The eight-role factorial is worth doing if we want a general statement about fictional overlays.
- The earlier correctness run is pilot evidence and is not pooled.
- Real-project tasks are the answer to synthetic ceiling effects—not increasingly absurd planted defects.
-->

---
layout: center
class: text-center
---

# Keep the receipts.

<div class="mx-auto max-w-3xl mt-8 grid grid-cols-2 gap-4 text-sm text-left">
  <div class="note"><b>Failures</b><br>are methodology</div>
  <div class="note"><b>Nulls</b><br>are product decisions</div>
  <div class="note"><b>Negative results</b><br>remove theatre</div>
  <div class="note"><b>Reproducibility</b><br>includes the pivots</div>
</div>

<div class="mt-10 text-[17px]">The council reviews the artifact.<br><b>The experiment must review itself.</b></div>

<div class="mt-8 text-xs opacity-55"><code>experiment/LAB_NOTEBOOK.md</code> · <code>experiment/CONTAMINATION_REVIEW.md</code> · <code>experiment/results/</code></div>

<style>
.note { background:#f5f3ff; border:1px solid #ddd6fe; border-radius:10px; padding:17px; }
</style>

<!-- Presenter cues
- The interesting engineering story is not that every run worked. It is that failed assumptions became explicit controls.
- Closing line: “If we only publish the final score, we hide the part that made the score worth believing.”
-->
