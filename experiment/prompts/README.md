# Controlled prompts

The harness assembles prompts mechanically from small frozen components.

- `system.txt`: common isolation and JSON-only system instruction.
- `review-contract.txt`: common evidence, output-cap, schema, and packet contract.
- `generic-kernel.txt`: S0 generic baseline.
- `omnibus-kernel.txt`: byte-identical seven-lens kernel for S1, S2, and M0.
- `omnibus-wrappers.json`: paired functional and GLaDOS wrappers for S1/S2.
- `specialists.json`: seven byte-identical functional kernels with paired functional/fictional wrappers for M1/M2 and provider replication.
- `specialist-intro.txt`: common narrow-lens instruction.
- `chain-contract.txt`: cumulative inherited-ledger contract for informed chains.
- `fuser.txt`: common persona-free fusion contract.
- `judge.txt`: arm-blinded smoke-triage prompt. Human ratings remain definitive for claim-bearing runs.

`experiment/harness/prompting.py` is the only prompt assembler. It fails if a reviewer prompt contains a planted defect ID, answer-key heading, or a long answer-key claim verbatim.

## Fairness checks

`just experiment-doctor <config>` assembles every planned prompt and checks that:

1. no answer-key content enters reviewer or fuser context;
2. each specialist pair differs only in its wrapper;
3. paired wrappers are within 10% by whitespace-token count;
4. the model identifiers are available locally;
5. the deterministic call count matches the design.

Whitespace-token count is a preflight proxy. Record provider-reported input tokens during smoke and revise wrapper wording before preregistration if provider tokenisation breaches the tolerance.

Every response must match `experiment/schema/findings.schema.json`. A malformed response is preserved and scores zero; the harness does not repair substantive output.
