# Four short talks

The [GitHub Pages hub](https://davehowell.github.io/council-of-nark/) hosts four talks in one dark, restrained theme. The visual hypothesis talk has ten slides and runs for approximately eight minutes. The other talks have six slides and run for approximately five minutes each. Every slide has a complete spoken script, timing cue and source references in its final Markdown comment.

| Talk | Editable source | Speaker script | Offline slides |
|---|---|---|---|
| 1. The Council of Nark | [slides](part-1/slides.md) | [notes](notes/part-1-notes.md) | [PDF](part-1/council-of-nark.pdf) |
| 2. Put the Council on Trial | [slides](part-2/slides.md) | [notes](notes/part-2-notes.md) | [PDF](part-2/put-the-council-on-trial.pdf) |
| 3. The Experiments Failed. A Lot. | [slides](part-3/slides.md) | [notes](notes/part-3-notes.md) | [PDF](part-3/the-experiments-failed-a-lot.pdf) |
| 4. What Would Settle It? | [slides](part-4/slides.md) | [notes](notes/part-4-notes.md) | [PDF](part-4/what-would-settle-it.pdf) |

Read [all scripts together](notes/speaker-notes.md). Part 1 contains the longer visual introduction; Parts 2–4 contain approximately 620–660 spoken words each. Rehearse at your own pace and trim if you speak slowly. Source blocks are reference material, not part of the spoken script. On the site, each talk also has a printable notes page. Slidev presenter mode displays the same notes; press P to open it.

## Edit and reproduce

Edit slide text and its final comment in `part-N/slides.md`. Edit `theme/styles.css` for all four talks. Keep notes in the slide source as the single authority; regenerate the standalone copies with `just slides-notes` after edits. [Slidev's notes syntax](https://sli.dev/guide/syntax#notes) documents the closing-comment convention.

```bash
just slides-install
just slides-talk            # part 1
just slides-experiment      # part 2
just slides-trials          # part 3
just slides-pivot           # part 4
just slides-notes           # readable scripts from slide comments
just slides-site-build      # four web decks + fresh PDFs + notes
just slides-export          # refresh tracked offline PDFs from that build
just slides-verify          # inspect all 28 production slides
```

Dependencies remain pinned in `bun.lock`. If the browser is missing, run `bunx playwright install chromium` from this directory. CI installs Chromium and its Linux dependencies before building.

The Pages workflow builds on pull requests and deploys only after changes reach `main`. Every build generates PDFs from the same Markdown as the web decks, preventing stale downloads. Notes are also generated during the build. The site retains published-result pages and the complete historical timeline; the short talks no longer need to carry the whole changelog.

## Evidence status

Part 3 reports exploratory calibration, including scoring limitations. Part 4 is a pivot into the **proposed, unrun** ecological comparison. It does not claim a conclusion, ecological outcomes, general persona harm, or single-agent equivalence. A conclusion talk belongs after that comparison has results. See [the evidence review](../experiment/CONCLUSION.md) and [next-stage protocol](../experiment/FINAL_ROUND.md).

Prior decks and PDFs remain available in Git history. Part 1 restores the visual narrative from `eca311a` while keeping the current shared theme and the later experimental caveats. The historical experiment results are not rewritten.
