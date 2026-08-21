# Ecological task: restored schema permanently skips a newer import

A pipeline uses an import-schema file and restores local state from its destination on ephemeral storage. The destination contains a schema older than the current import file. After restore, the pipeline behaves as if the newer import has already been applied, skips it, and fails on the new column. Repeating the run restores the same stale state and repeats the failure indefinitely.

Review the frozen pre-fix source. Trace how import provenance is recorded during normal saves versus destination restoration, identify the invalid state transition, and propose a narrow correction that preserves ordinary persistence behaviour. Specify storage-level and pipeline-level regression tests.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
