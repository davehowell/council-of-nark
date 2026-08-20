# Fictional-overlay factorial

> Status: designed, not run. Freeze and commit the final rating plan before execution.

## Why run the remaining pairs?

Yes—if the research question is whether fictional overlays are generally useful. Repeating only K-2SO cannot distinguish a general overlay effect from one character/lens interaction. Running whichever persona looks promising would create selection bias.

The appropriate next synthetic study is one balanced family, not seven disconnected follow-ups.

## Conditions

Pair the functional and fictional wrappers for:

1. simplicity / HK-47;
2. correctness / K-2SO;
3. architecture / GLaDOS;
4. security/compliance / C-3PO;
5. cost / Bender;
6. entropy / Holly;
7. technical language / WALL-E;
8. omnibus / GLaDOS overlay.

The lens kernel, packet, output schema, cap, model, and thinking level remain identical within each pair. The existing correctness result is pilot evidence and is **not pooled** into this fresh family.

## Frozen design

- Model: `gemma-4-31b-it` through isolated Pi, thinking off.
- Packets: the three unchanged synthetic packets.
- Repeats: ten per `(role × packet × wrapper)`.
- Calls: `8 roles × 3 packets × 10 repeats × 2 wrappers = 480`.
- Pairing blocks on role, packet, and repeat index. Calls remain independent; providers do not expose a common decoding seed.
- Primary sign: fictional minus functional, so a negative value favours the functional prompt.

Run only after the source and rating plan are clean and committed:

```bash
just experiment-doctor experiment/config/persona-factorial-gemma.json
just experiment-persona-factorial-gemma 2
```

## Estimands

Primary descriptive estimand:

- mean paired `F1_fictional − F1_functional` across all 240 pairs.

Mandatory supporting outcomes:

- paired precision and recall deltas;
- win/tie/loss counts;
- per-role and per-packet deltas;
- keyed-defect Jaccard overlap and one-sided contributions;
- input/output tokens and latency;
- false-positive clusters;
- malformed/error rates.

The harness's repeated-call bootstrap is labelled sampling-only. Reusing three tasks does not estimate generalisation to other software work.

## Decision rule to preregister

Use a practical F1 margin of `0.02` for the family mean:

- interval entirely below `−0.02`: fictional overlay is practically harmful here;
- interval within `[−0.02, +0.02]`: practically equivalent; retain names only if useful to humans;
- interval entirely above `+0.02`: fictional overlay is practically beneficial here;
- otherwise: inconclusive.

Do not infer that an individual persona helps from an uncorrected per-role interval. If individual roles are claim-bearing, preregister Holm correction across the eight role tests. Report every role regardless of direction.

## Rating requirement

LLM triage may calibrate the run but cannot make the confirmatory claim. Two independent humans must rate arm-blinded findings and reconcile disagreements without provider, role, or wrapper identity. Add a separate blinded remedy-quality rubric before rating; F1 alone cannot show that a fix is good.

## Scope

This factorial answers a narrow mechanism question on synthetic packets: does the prose overlay change supported defect detection for this model? It does not replace real-project ecological replication. If the family result is null or negative, mnemonic names can remain while functional prompts do the work.
