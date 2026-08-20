# Presentations

## Part 1: The Council of Nark

[`part-1/slides.md`](part-1/slides.md) introduces the review roster, explains the proposed fan-out/fuse topology, and frames improved coverage and lower variance as hypotheses. The exported deck is [`part-1/council-of-nark.pdf`](part-1/council-of-nark.pdf).

## Part 2: Put the Council on Trial

[`part-2/slides.md`](part-2/slides.md) turns the idea into separate tests of role specialisation, character wrappers, fusion, providers, and fan-out versus informed chains. The exported deck is [`part-2/put-the-council-on-trial.pdf`](part-2/put-the-council-on-trial.pdf).

## Run and export

The repository pins Slidev dependencies in this directory.

```bash
just slides-install
just slides-talk
just slides-experiment
just slides-export
just slides-verify
```

Slidev serves [`public/`](public/) at the site root. The PNG headshots are crops from the main council artwork.
