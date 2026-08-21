# Ecological task: process-dependent cache identity collision

A serializer memoizes both hashable and unhashable objects. Under a crafted collision between an object's hash value and another object's runtime identity integer, an object can be treated as already processed even though it was never recorded. Because some hashes vary between processes, this can make otherwise identical scene cache keys diverge and cause partial movie files to re-render unexpectedly.

Review the frozen pre-fix source. Identify how separate identity domains are conflated, propose the smallest correction that keeps existing cycle handling, and specify regression tests for both collision directions.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
