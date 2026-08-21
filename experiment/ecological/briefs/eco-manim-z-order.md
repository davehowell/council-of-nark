# Ecological task: animation z-order regression

A Cairo scene contains two persistent mobjects above a passing-flash animation. The first flash renders behind both, as expected. After one persistent mobject participates in a `Succession`, an equivalent flash can render in front of the other persistent mobject even though their z-index values did not change.

Review the frozen pre-fix source. Identify the most plausible mechanism, affected code, consequences beyond the minimal example, and a minimal durable correction. Specify a regression test that distinguishes the fix from a superficial redraw workaround.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
