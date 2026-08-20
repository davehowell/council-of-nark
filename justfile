set shell := ["bash", "-euc"]

# List recipes.
default:
    @just --list

# Scan public material for common secrets, private paths/hosts, and local forbidden terms.
audit:
    python3 scripts/public_audit.py

# Install pinned presentation dependencies.
slides-install:
    cd presentations && bun install --frozen-lockfile

# Serve the part-1 deck.
slides-talk:
    cd presentations && bunx --bun slidev part-1/slides.md

# Serve the part-2 protocol deck.
slides-experiment:
    cd presentations && bunx --bun slidev part-2/slides.md

# Export both PDFs.
slides-export:
    cd presentations && bunx --bun slidev export part-1/slides.md --output part-1/council-of-nark.pdf
    cd presentations && bunx --bun slidev export part-2/slides.md --output part-2/put-the-council-on-trial.pdf

# Render both decks to temporary PNG directories for visual inspection.
slides-verify:
    rm -rf /tmp/council-slides /tmp/council-experiment-slides
    cd presentations && bunx --bun slidev export part-1/slides.md --format png --output /tmp/council-slides
    cd presentations && bunx --bun slidev export part-2/slides.md --format png --output /tmp/council-experiment-slides
    @echo "Rendered /tmp/council-slides and /tmp/council-experiment-slides"

# Run Go harness unit and isolation tests without making model calls.
experiment-test:
    go test ./experiment/harness/...

# Prove Seatbelt permits scratch writes while denying repository reads.
experiment-sandbox-check:
    go run ./experiment/harness/cmd/council-exp sandbox-check

# Validate macOS, Seatbelt, config, prompts, model IDs, and key isolation without model calls.
experiment-doctor config="experiment/config/stage-a-smoke.json":
    go run ./experiment/harness/cmd/council-exp doctor "{{config}}"

# Freeze committed assets and create a deterministic plan; print the run path.
experiment-create config="experiment/config/stage-a-smoke.json":
    #!/usr/bin/env bash
    set -euo pipefail
    run="$(go run ./experiment/harness/cmd/council-exp freeze --config "{{config}}")"
    go run ./experiment/harness/cmd/council-exp plan "$run" >&2
    printf '%s\n' "$run"

# Execute or resume a planned run.
experiment-run run jobs="3":
    go run ./experiment/harness/cmd/council-exp run --jobs "{{jobs}}" "{{run}}"

# Summarize status, usage, cost, and latency without unblinding.
experiment-summary run:
    go run ./experiment/harness/cmd/council-exp summarize "{{run}}"

# Seal raw requests and responses by SHA-256 digest.
experiment-seal run:
    go run ./experiment/harness/cmd/council-exp seal "{{run}}"

# Verify a sealed run.
experiment-verify run:
    go run ./experiment/harness/cmd/council-exp verify "{{run}}"

# Create the arm-blinded human-rating bundle and runsheet.
experiment-bundle run:
    go run ./experiment/harness/cmd/council-exp bundle "{{run}}"

# Run exploratory arm-blinded LLM triage for a smoke test.
experiment-judge run jobs="2" config="experiment/config/judge-smoke.json":
    go run ./experiment/harness/cmd/council-exp judge --jobs "{{jobs}}" --config "{{config}}" "{{run}}"

# Score one adjudicated or exploratory blinded ratings CSV.
experiment-score run ratings label="adjudicated":
    go run ./experiment/harness/cmd/council-exp score --label "{{label}}" "{{run}}" "{{ratings}}"

# Freeze, plan, run, summarize, seal, verify, and blind any config.
experiment-complete config jobs="3":
    #!/usr/bin/env bash
    set -euo pipefail
    run="$(go run ./experiment/harness/cmd/council-exp freeze --config "{{config}}")"
    go run ./experiment/harness/cmd/council-exp plan "$run"
    go run ./experiment/harness/cmd/council-exp run --jobs "{{jobs}}" "$run"
    go run ./experiment/harness/cmd/council-exp summarize "$run"
    go run ./experiment/harness/cmd/council-exp seal "$run"
    go run ./experiment/harness/cmd/council-exp verify "$run"
    go run ./experiment/harness/cmd/council-exp bundle "$run"
    printf '\nRUN=%s\n' "$run"

# Make one frozen Anthropic/Haiku call through isolated Pi before a larger run.
experiment-adapter-check:
    @just experiment-complete experiment/config/adapter-check.json 1

# Verify explicit Pi/Gemma selection and prompt-enforced JSON with one frozen call.
experiment-adapter-check-gemma:
    @just experiment-complete experiment/config/adapter-check-gemma.json 1

# Run the original 81-call Haiku Stage A smoke test.
experiment-stage-a-smoke jobs="3":
    @just experiment-complete experiment/config/stage-a-smoke.json "{{jobs}}"

# Run the 81-call low-reasoning Gemma Stage A smoke test.
experiment-stage-a-smoke-gemma jobs="2":
    @just experiment-complete experiment/config/stage-a-smoke-gemma.json "{{jobs}}"

# Run 10 paired functional/fictional correctness samples per packet on Gemma.
experiment-persona-pair-gemma jobs="2":
    @just experiment-complete experiment/config/persona-pair-gemma-repeated.json "{{jobs}}"

# Exercise the complete harness with a deterministic local adapter and no model calls.
experiment-mock:
    @just experiment-complete experiment/config/mock-smoke.json 2
