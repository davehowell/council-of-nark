# Controlled prompt drafts

These prompts are ready for the smoke test. They are not preregistered or frozen until the smoke test confirms that the schemas and difficulty work.

- `s0-generic.txt` — weak generic single reviewer.
- `s1-omnibus.txt` — functional omnibus reviewer and the repeated-review prompt for M0.
- `s2-omnibus-glados.txt` — same omnibus checklist with a compact fictional wrapper.
- `specialist-template.txt` — assembly template for M1/M2 and Stage B.
- `specialist-kernels.md` — the seven controlled kernels and paired wrappers.
- `fuser.txt` — common persona-free fuser for M0/M1/M2 and the topology test.

Replace `{{REVIEW_PACKET}}`, `{{LENS_KERNEL}}`, `{{STYLE_WRAPPER}}`, and `{{REVIEW_FINDINGS}}` mechanically. Do not include any `answer-key.md` file in model context.

Before the confirmatory run:

1. Pin model snapshot IDs and decoding parameters.
2. Measure the actual token counts with each provider tokenizer.
3. Shorten paired wrappers until functional and fictional specialist prompts are within the preregistered tolerance.
4. Hash the final prompt files and record the hashes with the run manifest.
