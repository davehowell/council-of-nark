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

# Run harness unit tests without making model calls.
experiment-test:
    python3 -m unittest discover -s experiment/tests -v

# Validate config, prompt assembly, model IDs, and answer-key isolation without model calls.
experiment-doctor config="experiment/config/stage-a-smoke.json":
    python3 -m experiment.harness.doctor "{{config}}"

# Freeze committed assets and create a deterministic plan; print the run path.
experiment-create config="experiment/config/stage-a-smoke.json":
    #!/usr/bin/env bash
    set -euo pipefail
    run="$(python3 -m experiment.harness.freeze --config "{{config}}")"
    python3 -m experiment.harness.plan "$run" >&2
    printf '%s\n' "$run"

# Execute or resume a planned run.
experiment-run run jobs="3":
    python3 -m experiment.harness.run "{{run}}" --jobs "{{jobs}}"

# Summarize status, usage, cost, and latency without unblinding.
experiment-summary run:
    python3 -m experiment.harness.summarize "{{run}}"

# Seal raw requests and responses by SHA-256 digest.
experiment-seal run:
    python3 -m experiment.harness.seal "{{run}}"

# Verify a sealed run.
experiment-verify run:
    python3 -m experiment.harness.verify "{{run}}"

# Create the arm-blinded human-rating bundle and runsheet.
experiment-bundle run:
    python3 -m experiment.harness.bundle "{{run}}"

# Run exploratory arm-blinded LLM triage for a smoke test.
experiment-judge run jobs="2":
    python3 -m experiment.harness.judge "{{run}}" --jobs "{{jobs}}"

# Score one adjudicated or exploratory blinded ratings CSV.
experiment-score run ratings label="adjudicated":
    python3 -m experiment.harness.score "{{run}}" "{{ratings}}" --label "{{label}}"

# Freeze, plan, run, summarize, seal, verify, and blind any config.
experiment-complete config jobs="3":
    #!/usr/bin/env bash
    set -euo pipefail
    run="$(python3 -m experiment.harness.freeze --config "{{config}}")"
    python3 -m experiment.harness.plan "$run"
    python3 -m experiment.harness.run "$run" --jobs "{{jobs}}"
    python3 -m experiment.harness.summarize "$run"
    python3 -m experiment.harness.seal "$run"
    python3 -m experiment.harness.verify "$run"
    python3 -m experiment.harness.bundle "$run"
    printf '\nRUN=%s\n' "$run"

# Make one frozen model call to verify the Claude adapter before a larger run.
experiment-adapter-check:
    @just experiment-complete experiment/config/adapter-check.json 1

# Run the 81-call Stage A smoke test.
experiment-stage-a-smoke jobs="3":
    @just experiment-complete experiment/config/stage-a-smoke.json "{{jobs}}"

# Exercise the complete harness with a deterministic local adapter and no model calls.
experiment-mock:
    @just experiment-complete experiment/config/mock-smoke.json 2
