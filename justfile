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
