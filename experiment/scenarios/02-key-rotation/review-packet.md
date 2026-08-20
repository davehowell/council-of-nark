# Review packet 02: payment key rotation at 03:00

Review this incident runbook before the on-call team adopts it. Report only material findings that have a concrete consequence and fix.

## Service facts

- `payments-api` runs in the `prod-payments` namespace. Deleting its only API-key Secret causes payment failures until pods receive a replacement.
- The deployment normally uses an immutable image digest.
- Pager alerts may be muted for at most 15 minutes and must restore automatically.
- The provider can keep the old and new key valid together for 30 minutes.
- The key value is secret and must never appear in a shell history, ticket, or log.

## `runbooks/rotate-payment-key.md`

```text
01  # Rotate the payment provider key
02
03  This process facilitates the seamless rotation of credentials in production.
04
05  1. Export the new key: export PAYMENT_KEY="example-do-not-use-key-material".
06  2. Silence the PaymentsDown alert in Alertmanager.
07  3. Delete the old secret: kubectl delete secret payment-api-key.
08  4. Create the new secret and restart payments-api.
09  5. Deploy payments-api:latest to make sure it picks everything up.
10  6. Check the logs and dashboard.
11  7. Re-enable the alert when things are stable.
12  8. Restart the worker if errors continue.
13
14  NOTE: Be very careful. This could cause unexpected issues.
```

## Deployment note

A new `RotationCoordinatorFactory` service will call separate validator, secret-writer, restart, and notifier services. The team expects one key rotation every 90 days. The runbook is the only rollback documentation.
