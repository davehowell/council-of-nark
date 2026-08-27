---
theme: seriph
title: The Experiment Fought Back
colorSchema: light
info: |
  Part 3 of the Council of Nark series: instrumentation failures, ceiling effects,
  negative persona evidence, isolation hardening, and the pivot to sealed ecological tasks.
transition: slide-left
mdc: true
layout: image-right
image: /Nark-council.png
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

# The next test changed shape

<div class="fork mt-8">
  <div class="road paused"><small>planned</small><b>Persona factorial</b><strong>480 calls</strong><p>Generalise one negative overlay result across eight roles.</p></div>
  <div class="junction">?</div>
  <div class="road chosen"><small>prioritised</small><b>Ecological pilot</b><strong>real code</strong><p>First prove the mechanism survives outside planted packets.</p></div>
</div>

<div class="decision mt-8">Ceiling effects made “more synthetic calls” less urgent than “better evidence.”</div>

<style>
.fork { display:grid; grid-template-columns:1fr 70px 1fr; align-items:center; gap:18px; }.road { border:1px solid #cbd5e1; border-radius:16px; padding:22px; }.road small,.road b,.road strong { display:block; }.road small { text-transform:uppercase; letter-spacing:.12em; opacity:.55; }.road b { font-size:21px; margin-top:8px; }.road strong { font-size:28px; margin-top:10px; }.road p { font-size:13px; line-height:1.45; margin-top:10px; }.paused { background:#f8fafc; }.paused strong { color:#64748b; }.chosen { background:#ecfdf5; border-color:#6ee7b7; }.chosen strong { color:#047857; }.junction { font-size:46px; text-align:center; color:#7c3aed; font-weight:900; }.decision { background:#f5f3ff; border-radius:11px; padding:15px; text-align:center; font-weight:700; }
</style>

<!-- Presenter cues
- The 480-call factorial remains a valid preregistered question; it was not cancelled because one persona lost.
- But another large synthetic batch would not solve the observed lack of headroom.
- We prioritised ecological validity and left the factorial visibly unrun rather than selectively testing promising characters.
-->

---
layout: default
---

# Real-code curation was a funnel

<div class="funnel mt-5">
  <div><b>8</b><span>PR-backed tasks</span><small>Manim · dlt · Gortex · turbovec</small></div>
  <i>+</i>
  <div class="watch"><b>1</b><span>extreme watchlist</span><small>Modular CPU reduction</small></div>
  <i>→</i>
  <div class="pilot"><b>1</b><span>first pilot</span><small>Gortex Unicode tokenizer</small></div>
</div>

<div class="grid grid-cols-3 gap-4 mt-8 text-sm">
  <div class="rule"><b>After cutoff</b><br>reduces—not removes—training risk</div>
  <div class="rule"><b>Pre-fix parent</b><br>history and upstream evidence removed</div>
  <div class="rule"><b>Known test</b><br>must fail before and pass after</div>
</div>

<div class="mt-6 text-center text-xs opacity-60">Modular stayed gated: no PR review and no frozen external build closure.</div>

<style>
.funnel { display:grid; grid-template-columns:1fr 35px 1fr 45px 1fr; align-items:center; gap:10px; }.funnel>div { padding:18px; border:1px solid #cbd5e1; border-radius:14px; text-align:center; background:#f8fafc; }.funnel b { display:block; font-size:38px; color:#7c3aed; }.funnel span,.funnel small { display:block; }.funnel span { font-weight:800; }.funnel small { margin-top:7px; opacity:.6; }.funnel i { font-size:27px; color:#94a3b8; text-align:center; font-style:normal; }.funnel .watch { border-style:dashed; }.funnel .pilot { background:#ecfdf5; border-color:#6ee7b7; }.rule { border-left:4px solid #7c3aed; padding:12px 14px; background:#f5f3ff; border-radius:8px; }
</style>

<!-- Presenter cues
- The selected set has four pilot and four reserve tasks, all backed by public fixes merged after 1 June 2026.
- The eventual patch is evidence, not the only acceptable answer. Supported novel findings remain valid.
- Modular looked attractive because it was extreme, but methodological eligibility mattered more than spectacle.
-->

---
layout: default
---

# The snapshot failed four different ways

<div class="stair mt-5">
  <div><em>1</em><b>Extract</b><span>missing intermediate directory</span></div>
  <div><em>2</em><b>Move</b><span>read-only module cache</span></div>
  <div><em>3</em><b>Clean</b><span>duplicate toolchain staging</span></div>
  <div><em>4</em><b>Test</b><span>Go VCS stamping escaped upward</span></div>
  <div class="pass"><em>5</em><b>Validate</b><span>parent failed · evidence passed</span></div>
</div>

<div class="mt-8 text-center text-sm"><b>Rule:</b> preserve each attempt; repair code; start a new attempt.</div>

<style>
.stair { display:grid; grid-template-columns:repeat(5,1fr); gap:9px; align-items:end; height:245px; }.stair div { border:1px solid #fca5a5; background:#fef2f2; border-radius:11px 11px 0 0; padding:12px; min-height:105px; }.stair div:nth-child(2){min-height:135px}.stair div:nth-child(3){min-height:165px}.stair div:nth-child(4){min-height:195px}.stair div:nth-child(5){min-height:225px}.stair em { display:block; font-size:26px; color:#b91c1c; font-weight:900; font-style:normal; }.stair b,.stair span { display:block; }.stair span { font-size:11px; line-height:1.35; margin-top:8px; opacity:.7; }.stair .pass { background:#ecfdf5; border-color:#6ee7b7; }.stair .pass em { color:#047857; }
</style>

<!-- Presenter cues
- These were exporter and environment failures, not model failures.
- The fourth failure mattered most: the test process discovered the enclosing council checkout through build metadata.
- Setting buildvcs=false closed that path. The fifth development attempt reproduced the panic and passed at the evidence commit.
-->

---
layout: default
---

# Clean source is infrastructure—not a result

<div class="grid grid-cols-2 gap-7 mt-5">
  <div class="snapshot"><small>history-free Gortex parent</small><b>4,093</b><span>entries</span><b>37.8 MB</b><span>sanitised source</span></div>
  <div class="checks">
    <div>✓ <span>tree digest <code>41a9a14b…</code></span></div>
    <div>✓ <span><code>.git</code> absent</span></div>
    <div>✓ <span>parent panic reproduced</span></div>
    <div>✓ <span>evidence commit passed</span></div>
    <div>✓ <span>source / controller / network probes</span></div>
  </div>
</div>

<div class="zero mt-7"><b>0</b><span>ecological respondent calls</span></div>

<style>
.snapshot { border:1px solid #cbd5e1; border-radius:15px; padding:20px; text-align:center; background:#f8fafc; }.snapshot small,.snapshot b,.snapshot span { display:block; }.snapshot small { opacity:.6; }.snapshot b { color:#7c3aed; font-size:32px; margin-top:8px; }.snapshot span { font-size:12px; }.checks { display:grid; gap:8px; }.checks div { background:#ecfdf5; border:1px solid #a7f3d0; border-radius:9px; padding:11px 14px; color:#047857; }.checks span { color:#1f2937; margin-left:7px; }.zero { display:flex; justify-content:center; align-items:center; gap:13px; border:2px solid #f59e0b; background:#fffbeb; border-radius:12px; padding:13px; }.zero b { font-size:32px; color:#b45309; }.zero span { font-weight:800; }
</style>

<!-- Presenter cues
- The clean rerun came from a committed controller and reproduced the exact source digest from development.
- Snapshot validity means the task and closure are reproducible. It says nothing yet about reviewer performance.
- We keep “infrastructure evidence” visually distinct from “respondent result” to prevent accidental claim inflation.
-->

---
layout: default
---

# Source and network could not share a process

<div class="boundary mt-4">
  <div class="box provider"><b>Provider + model</b><span>transport network</span><small>no source mount</small></div>
  <div class="pipe">Pi tool protocol<br>↕</div>
  <div class="box pi"><b>Isolated Pi</b><span>custom tools only</span><small>two inherited JSON pipes</small></div>
  <div class="pipe">bounded requests<br>↕</div>
  <div class="box mediator"><b>Go mediator</b><span>list · read · RE2 search</span><small>allowlisted focused test</small></div>
</div>

<div class="grid grid-cols-2 gap-5 mt-7 text-sm">
  <div class="zone green"><b>Source operations</b><br>in-process, confined, byte/call budgets</div>
  <div class="zone red"><b>Test process</b><br>separate Seatbelt, frozen closure, network denied</div>
</div>

<style>
.boundary { display:grid; grid-template-columns:1fr 120px 1fr 120px 1fr; align-items:center; gap:8px; }.box { text-align:center; border-radius:14px; padding:20px 12px; border:1px solid #cbd5e1; }.box b,.box span,.box small { display:block; }.box span { margin-top:8px; }.box small { margin-top:7px; opacity:.6; }.provider { background:#eff6ff; border-color:#93c5fd; }.pi { background:#f5f3ff; border-color:#c4b5fd; }.mediator { background:#ecfdf5; border-color:#6ee7b7; }.pipe { text-align:center; font:12px monospace; color:#64748b; }.zone { border-radius:10px; padding:14px 18px; }.green { background:#ecfdf5; }.red { background:#fef2f2; }
</style>

<!-- Presenter cues
- Mounting source into a network-enabled model process would collapse the contamination boundary.
- The extension knows file descriptors, not paths. The mediator—not Pi—owns source and test access.
- Every tool request and result is correlated in a controller-owned, fail-closed transcript.
-->

---
layout: default
---

# The doctor found assumptions, not bugs

<div class="doctor-flow mt-5">
  <div><b>EOF</b><span>RPC stayed alive</span></div><i>→</i>
  <div><b>thinking off</b><span>clamped to minimal</span></div><i>→</i>
  <div><b>shutdown()</b><span>idle process remained</span></div><i>→</i>
  <div><b>SIGTERM</b><span>same persistent path</span></div>
</div>

<div class="repair mt-8">
  <div><small>doctor repair</small><b>verify <code>minimal</code> explicitly</b></div>
  <div><small>lifecycle repair</small><b>controlled termination after checks</b></div>
  <div><small>claim runner</small><b>use JSON/print mode</b></div>
</div>

<div class="mt-7 text-center text-sm opacity-70">All failed attempts preserved. No prompt sent. No provider turn.</div>

<style>
.doctor-flow { display:grid; grid-template-columns:1fr 28px 1fr 28px 1fr 28px 1fr; align-items:center; gap:5px; }.doctor-flow div { padding:17px 10px; border:1px solid #fca5a5; background:#fef2f2; border-radius:11px; text-align:center; }.doctor-flow b,.doctor-flow span { display:block; }.doctor-flow b { font-family:monospace; }.doctor-flow span { font-size:11px; margin-top:7px; opacity:.65; }.doctor-flow i { font-style:normal; color:#7c3aed; font-size:22px; text-align:center; }.repair { display:grid; grid-template-columns:repeat(3,1fr); gap:12px; }.repair div { background:#ecfdf5; border:1px solid #a7f3d0; border-radius:11px; padding:15px; text-align:center; }.repair small,.repair b { display:block; }.repair small { text-transform:uppercase; opacity:.55; }.repair b { margin-top:7px; }
</style>

<!-- Presenter cues
- We initially treated persistent RPC shutdown as a defect. It is the wrong lifecycle for a one-shot claim runner.
- Pi also accepted “off” syntactically but selected the model's supported minimum. State must be observed, not inferred from an accepted command.
- The final clean doctor sealed runtime, profiles, events, stderr, boundary probe and mediator transcript with zero provider turns.
-->

---
layout: default
---

# The evidence ladder now has a hard gate

<div class="ladder mt-4">
  <div class="done"><b>1</b><span>Synthetic plumbing</span><small>complete</small></div>
  <div class="done"><b>2</b><span>Negative overlay calibration</span><small>complete</small></div>
  <div class="done"><b>3</b><span>History-free snapshot</span><small>sealed</small></div>
  <div class="done"><b>4</b><span>Mediator + Pi doctor</span><small>sealed · zero calls</small></div>
  <div class="gate"><b>5</b><span>Claim runner + mock</span><small>current gate</small></div>
  <div><b>6</b><span>Compared arms</span><small>not started</small></div>
</div>

<div class="mt-7 grid grid-cols-2 gap-5 text-sm">
  <div class="next"><b>Before first call</b><br>freeze byte-level prompt factor, model, retries, endpoint and human aggregation</div>
  <div class="next"><b>Then measure</b><br>supported findings, remedy quality, tool/token/cost/latency efficiency</div>
</div>

<style>
.ladder { display:grid; grid-template-columns:repeat(6,1fr); gap:8px; align-items:end; }.ladder div { min-height:145px; border:1px solid #cbd5e1; border-radius:11px; padding:12px; background:#f8fafc; }.ladder b,.ladder span,.ladder small { display:block; }.ladder b { font-size:25px; color:#94a3b8; }.ladder span { font-size:12px; font-weight:800; margin-top:9px; }.ladder small { font-size:10px; margin-top:8px; opacity:.6; }.ladder .done { background:#ecfdf5; border-color:#6ee7b7; }.ladder .done b { color:#047857; }.ladder .gate { background:#fffbeb; border:2px solid #f59e0b; transform:translateY(-8px); }.ladder .gate b { color:#b45309; }.next { background:#f5f3ff; border-radius:10px; padding:15px 18px; }
</style>

<!-- Presenter cues
- The clean mediator and doctor are frozen infrastructure checks, not ecological outcomes.
- This ladder prevents a successful sandbox check from being narrated as evidence that the council found a real defect.
- The next model call remains gated by a deterministic mock lifecycle and a preregistered comparison.
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
