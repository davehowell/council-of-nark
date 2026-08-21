# Ecological task: healthy daemon reports as unavailable under indexing

While a large repository is being tracked or enriched, serial status requests can vary from milliseconds to more than ten seconds even when they return the same payload. A short-budget session-start probe interprets the delay as an unhealthy daemon and injects that false warning into later agent context. The control plane should remain responsive while long data-plane work proceeds.

Review the frozen pre-fix source. Identify lock scope and snapshot-consistency constraints, propose a concurrency-safe separation, and discuss stale-but-coherent versus blocked status. Specify tests that prove status remains bounded while search/index mutation is active and that shutdown remains safe.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
