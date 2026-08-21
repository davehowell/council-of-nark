# Ecological task: stale temporary files survive forever

The persistence layer creates sibling temporary files and later sweeps old ones. For destination basenames near the filesystem component-length limit, a killed writer can leave a full-size temporary file that subsequent sweeps never remove. Repeated crashes can fill the destination volume. Short destination names are reclaimed correctly.

Review the frozen pre-fix source. Reconcile temporary-name creation with stale-file recognition, identify collision or over-deletion risks in any correction, and propose portable tests for long names, names containing the temporary marker, age thresholds, and unrelated files sharing a prefix.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
