# Private answer key 03

Do not send this file to reviewers.

| ID | Location | Planted defect | Minimum accepted claim |
|---|---|---|---|
| WH-01 | Design items 2–3 | Check-then-insert is race-prone and there is no idempotent merge/unique enforcement for 24-hour retries. | Concurrent retries can create duplicates; ingest by `event_id` with an atomic idempotent write or deduplication step. |
| WH-02 | Item 1 | The HTTP acknowledgement blocks on BigQuery, dbt, and Cube refreshes. | The path cannot reliably meet two seconds and couples receipt to analytics availability; durably enqueue then acknowledge. |
| WH-03 | Item 4 | Complete headers/payloads log bearer tokens and restricted email. | Secrets and PII enter logs; redact/allow-list fields before logging. |
| WH-04 | Item 5 | Six services and six queues are disproportionate to 20 events/minute. | The topology adds operations and failure modes without scale need; start with one receiver and one durable queue/worker boundary. |
| WH-05 | Item 6 | An unversioned shared JSON contract can change after only a Slack message. | Consumers can break on rename/type drift; validate `schema_version` and use a compatible version/deprecation contract. |
| WH-06 | Item 7 | Failed messages have no dead-letter/replay state, age metric, or alert; detection waits for a partner report. | Data loss/backlog is silent; add bounded retries, dead-letter storage, age/depth metrics, and an actionable alert. |
| WH-07 | Item 8 | A manual pause has no expiry and returns success while discarding events. | Forgotten state or maintenance traffic causes acknowledged data loss; make storage the safe default and use an expiring pause/backpressure response. |
| WH-08 | Item 9 | Replay performs an unpartitioned full-table read and exposes an operational replay action directly to the dashboard. | Every replay can scan the full table and an end user can trigger duplicate/expensive work; require partition bounds and put an authorised, idempotent control API between the dashboard and replay. |

Do not award generic “add tests/docs” findings unless they identify one of these mechanisms and consequences.
