# Offline ecological rating page

`index.html` is a network-free viewer and form for absolute human ratings. Open it locally in a modern browser, enter the rater ID and disclosed role, then import a bundle that conforms to `rating-bundle.schema.json`.

The page shows:

- the problem statement;
- the original repository's reference-resolution title and summary;
- the applied patch and changed tests;
- one opaque experimental review at a time;
- the five technical 1–5 scores and four personal-utility 1–7 scores;
- reference relation, qualitative comments, and a post-score treatment guess.

It rejects extra bundle fields so arm, provider, token, latency, and run metadata cannot be included accidentally. Bundle text is inserted with `textContent`, not interpreted as HTML. The page makes no network request and stores progress only in the current page. Export a checkpoint after each review if the browser session may be interrupted. The exported JSON remains blinded and conforms to `rating-record.schema.json`.

The whole-block controller must eventually generate and seal the bundle, private condition map, order, and digest. Do not assemble a claim-bearing bundle by hand. Pairwise presentation happens only after both absolute-rating files are locked; its generator and page are not implemented yet.

See [`../RATING.md`](../RATING.md) for the procedure and interpretation rules.
