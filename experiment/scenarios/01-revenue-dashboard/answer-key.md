# Private answer key 01

Do not send this file to reviewers. A finding maps to at most one ID. Equivalent wording is accepted only when it states the consequence.

| ID | Location | Planted defect | Minimum accepted claim |
|---|---|---|---|
| RD-01 | `customer_revenue.sql:L04-L13` | The one-to-many shipment join duplicates each order before `sum(order_total)`. | Revenue is over-counted for orders with multiple shipment events; pre-aggregate/deduplicate shipment data or aggregate orders before the join. |
| RD-02 | SQL config/filter and recovery step 2 | An hourly full table rebuild has no `_PARTITIONDATE` filters and the recovery command forces full refresh. | The query scans the 20 TiB and 6 TiB sources instead of bounded partitions; add partition predicates and use incremental processing. |
| RD-03 | Cube dimensions | Restricted `customer_email` is exposed with `shown: true`. | PII becomes visible to dashboard users; remove, hide, or mask it and use `customer_id`. |
| RD-04 | SQL select/group and schema YAML | The required stable `customer_id` grain is absent; mutable email is the grouping key and there is no uniqueness test for the required grain. | Customers can split/merge when email changes; group by `customer_id` and test `unique`/`not_null`. |
| RD-05 | Design note, second sentence | The fallback lets Superset bypass the required Cube contract and query raw tables. | Business logic/security can diverge across layers; keep the fallback behind the curated Cube contract. |
| RD-06 | Design note, first sentence | A strategy factory and unused database strategies are speculative machinery for one fixed backend. | The abstraction has no current consumer; call the BigQuery implementation directly until a second backend exists. |
| RD-07 | Recovery steps 1 and 3 | The quality-check flag relies on a human to restore it and has no two-hour TTL. | Checks can remain disabled after the incident; use an expiring override or automatic restoration. |
| RD-08 | Recovery step 2 and NOTE | A destructive/expensive full refresh has only a vague warning, and “looks right” gives no validation criterion. | The operator cannot predict cost/data impact or know when to restore checks; name the consequence, bounded command, checks, owner, and success criteria. |

Not planted: `shipment_event_count` itself is allowed to count events; `total_revenue` as a Cube sum is valid if the upstream grain is fixed. Do not award claims based only on taste or catchphrases.
