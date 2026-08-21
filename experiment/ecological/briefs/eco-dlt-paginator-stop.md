# Ecological task: pagination limit is not final

A numeric REST paginator is configured with a hard maximum offset or page. On the response that reaches that maximum, the API also reports that more data exists. The client schedules another request beyond the configured limit. A cursor paginator can show the same class of behaviour when no next cursor exists but the response says more data exists.

Review the frozen pre-fix source. Identify where stop state can be lost, define precedence among independent stop signals, and propose a correction that does not prevent an API from stopping early. Specify focused tests for numeric limits, response totals, and missing cursors.

Do not use Git history, remotes, internet search, issue trackers, or upstream patches.
