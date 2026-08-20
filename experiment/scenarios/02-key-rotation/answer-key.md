# Private answer key 02

Do not send this file to reviewers.

| ID | Location | Planted defect | Minimum accepted claim |
|---|---|---|---|
| KR-01 | Runbook line 05 | Example key material is embedded in the runbook and exported into shell history. | The runbook teaches a pattern that can leak real key material through source control or shell history; use an approved secret-input mechanism. |
| KR-02 | Lines 07–08 | The only Secret is deleted before its replacement exists, despite a 30-minute overlap window. | Payments can fail between deletion and rollout; create/validate the new version first, roll pods, then revoke the old key. |
| KR-03 | Lines 07–09 | Commands omit the `prod-payments` namespace/context guard. | The on-call operator can mutate the wrong cluster or namespace; verify context and pass the namespace explicitly. |
| KR-04 | Line 09 | The deployment changes from an immutable digest to `latest`. | A restart can pull unreviewed or different code and cannot be reproduced; retain the approved digest. |
| KR-05 | Lines 06 and 11 | Alert restoration depends on a human notion of “stable” and has no 15-minute expiry. | The alert can stay muted after the incident; use a timed silence and explicit owner/success check. |
| KR-06 | Lines 10–14 | Validation, the conditional restart, and the safety note are vague; the condition is buried after the action and no concrete rollback exists. | A 03:00 operator cannot tell what to inspect, when to restart, what damage can occur, or how to reverse it; give exact checks, condition-first steps, consequence, and rollback. |
| KR-07 | Deployment note | Four services plus a factory are proposed for one rotation every 90 days. | The machinery creates deployment and failure surface without demonstrated scale; use one small script/job or the provider workflow. |
| KR-08 | Whole runbook | No metric/alert verifies payment success during and after the rotation. | A technically successful rollout can silently reject payments; specify a payment success/error-rate gate and rollback threshold. |

Line 03 is deliberately poor prose but is not a separate planted defect; it may support KR-06. Do not award style-only findings that state no operational consequence.
