# Ecological task: corpus-dependent Unicode panic

An `explore` request over a repository containing non-ASCII prose can fail with an internal `index out of range [1] with length 1` panic. A lowercase Cyrillic query can trigger it even though the index is healthy. ASCII-only repositories have not reproduced the failure.

Review the frozen pre-fix source. Locate the unsafe assumption, explain why the query text alone is not sufficient to predict the crash, and propose a Unicode-correct correction without turning tokenization into an allocation-heavy hot path. Specify regression cases across relevant script/case boundaries.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
