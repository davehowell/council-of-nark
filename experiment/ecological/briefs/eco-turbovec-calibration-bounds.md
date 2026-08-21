# Ecological task: finite calibration values silently poison search

A persisted index accepts finite, positive calibration scales with extreme magnitudes. Certain accepted values make later search arithmetic overflow or produce NaN for every score; top-k then degenerates to input order without reporting corruption. One poisoned coordinate is sufficient and survives a save/load round trip.

Review the frozen pre-fix source. Derive validation from the quantities used by search rather than choosing an arbitrary epsilon, identify every construction/load boundary that must enforce it, and propose tests for silent NaN/Inf, round trips, boundary values, and valid low-magnitude data.

Scope this task to calibration-value validation; ignore unrelated public helpers. Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
