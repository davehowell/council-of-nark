# Presentations

The presentable GitHub Pages hub is [davehowell.github.io/council-of-nark](https://davehowell.github.io/council-of-nark/). It hosts all three Slidev decks, downloadable PDFs, published-result tiles, and the living experiment timeline.

## Part 1: The Council of Nark

[`part-1/slides.md`](part-1/slides.md) introduces the review roster, explains the proposed fan-out/fuse topology, and frames improved coverage and lower variance as hypotheses. The exported deck is [`part-1/council-of-nark.pdf`](part-1/council-of-nark.pdf).

## Part 2: Put the Council on Trial

[`part-2/slides.md`](part-2/slides.md) turns the idea into separate tests of role specialisation, character wrappers, fusion, providers, and fan-out versus informed chains. The exported deck is [`part-2/put-the-council-on-trial.pdf`](part-2/put-the-council-on-trial.pdf).

## Part 3: The Experiment Fought Back

[`part-3/slides.md`](part-3/slides.md) preserves the engineering story: instrumentation failure, ceiling effects, contamination audit, discarded runs, parser repairs, negative persona evidence, the Go/Seatbelt migration, ecological curation, repeated snapshot failures, and the sealed mediator/Pi boundary. The exported deck is [`part-3/the-experiment-fought-back.pdf`](part-3/the-experiment-fought-back.pdf).

## Run and export

The repository pins Slidev dependencies in this directory.

```bash
just slides-install
just slides-talk
just slides-experiment
just slides-trials
just slides-site-build
just slides-site-preview
just slides-export
just slides-verify
```

`site/` contains the static Pages shell. `scripts/build-site.mjs` builds each deck beneath `dist/decks/`, copies the PDFs, and applies the repository base path in CI. Append timeline entries to [`site/timeline/events.js`](site/timeline/events.js), keeping each claim aligned with [`experiment/LAB_NOTEBOOK.md`](../experiment/LAB_NOTEBOOK.md). Story mode makes the timeline keyboard-presentable.

Slidev serves [`public/`](public/) at the site root. The PNG headshots are crops from the main council artwork. GitHub Actions deploys `presentations/dist/` from `main`; generated output remains ignored locally.
