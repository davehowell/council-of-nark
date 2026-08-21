# Extreme ecological task: parallel CPU arg-reduction index corruption

On a CPU host with at least four logical workers, `argmax` and `argmin` can return the wrong global index for a single row with a reduction axis in the hundreds of thousands or millions of elements. Small rows are correct. A representative failure places the unique winning value well beyond the first worker's portion of a `[1, 2097152]` input. The incorrect behaviour appears only after the implementation selects its split-axis CPU tier.

Review the frozen pre-fix Modular source. Trace the generic reduction state from each worker through publication, cross-worker combination, and final output. Identify the violated state or representation invariant and any scratch-layout consequence. Propose a correction that remains valid for other reduction monoids and worker counts, and state why its ordering is safe.

Specify focused regressions for both `argmax` and `argmin`, winners in different worker portions, threshold-adjacent axis sizes, NaN behaviour, hosts that cannot activate the split tier, and compile-time protection against an undersized partial-state slot.

Do not use Git history, changelogs, remotes, internet search, issue trackers, upstream commits, or synchronized internal revision IDs.
