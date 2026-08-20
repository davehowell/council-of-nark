# Review packet 01: revenue dashboard change

Review the proposed change before merge. Report only material findings that have a concrete consequence and fix.

## Service facts

- `raw.orders` is 20 TiB and ingestion-time partitioned by `_PARTITIONDATE`.
- `raw.shipment_events` is 6 TiB and can contain many events for one order.
- The dashboard needs the most recent 30 days and refreshes every hour.
- One output row must represent one stable `customer_id`.
- `customer_email` is restricted PII and must not be visible in Cube or Superset.
- Curated dashboards must query the Cube contract. They must not query raw tables.
- A disabled quality check must restore itself within two hours.

## `models/customer_revenue.sql`

```sql
-- L01
{{ config(materialized='table') }}

-- L04
select
  o.customer_email,
  sum(o.order_total) as total_revenue,
  count(s.event_id) as shipment_event_count
from `acme.raw.orders` o
left join `acme.raw.shipment_events` s
  on o.order_id = s.order_id
where date(o.created_at) >= date_sub(current_date(), interval 30 day)
group by 1
```

## `models/customer_revenue.yml`

```yaml
# L01
version: 2
models:
  - name: customer_revenue
    columns:
      - name: customer_email
        tests: [not_null]
```

## `cube/customer_revenue.yml`

```yaml
# L01
cubes:
  - name: customer_revenue
    sql_table: analytics.customer_revenue
    dimensions:
      - name: customer_email
        sql: customer_email
        type: string
        shown: true
    measures:
      - name: total_revenue
        sql: total_revenue
        type: sum
```

## Design note

`RevenueQueryStrategyFactory` will return `BigQueryRevenueStrategy`. The factory keeps room for Snowflake and Postgres strategies, although no migration or second database is planned. If Cube is delayed, Superset will temporarily query `raw.orders` through a saved SQL dataset.

## Recovery steps

1. Set `disable_revenue_checks=true`.
2. Run `dbt build --full-refresh --select customer_revenue`.
3. Restore the check after the dashboard looks right.

NOTE: Be careful because this may affect data.
