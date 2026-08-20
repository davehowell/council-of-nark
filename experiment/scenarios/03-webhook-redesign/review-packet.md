# Review packet 03: partner webhook redesign

Review the proposed design before implementation. Report only material findings that have a concrete consequence and fix.

## Requirements and scale

- Partners send 20 events per minute on average and can retry the same event for 24 hours.
- Every event contains `event_id`, `schema_version`, `customer_email`, and a bearer token used only to authenticate receipt.
- Acknowledgement must complete within two seconds.
- Analytics may lag by five minutes.
- Raw bearer tokens and customer email must not enter application logs.
- BigQuery is partitioned by `received_date`; on-demand query cost is charged to this team.
- The service must support replay after a downstream outage.

## Proposed design

1. The HTTP handler writes each payload to BigQuery, then waits for dbt and Cube refresh calls before returning HTTP 200.
2. Before every insert, the handler runs `SELECT COUNT(*) FROM analytics.webhook_events WHERE event_id = ?` to check whether the event exists.
3. The insert table has no unique constraint or merge step because partner retries are expected to be rare.
4. For debugging, the service logs the complete request headers and payload.
5. Six new services will handle receipt, validation, routing, persistence, notification, and replay. Each service has its own queue.
6. All services consume the same unversioned JSON shape. Producers can add or rename fields after notifying the team in Slack.
7. Failed queue messages remain on the queue. The team will inspect queue depth when a partner reports missing data.
8. During downstream maintenance, an operator sets `pause_ingestion=true` and clears it after the maintenance window. Paused requests return HTTP 200 without storing the payload.
9. The replay service reads the complete BigQuery table and republishes every matching partner record. The initial dashboard will call this service directly for a “replay” button.
