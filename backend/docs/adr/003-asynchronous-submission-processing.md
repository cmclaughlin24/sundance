<div align="center">

# #003 Asynchronous Submission Processing

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Submission processing involves multiple steps — rule evaluation, field validation, data source lookup validation against the Tenants Service, and canonical tag mapping. The latency of these steps is variable and dependent on external infrastructure. Performing this work synchronously on the `POST /submissions` request would couple intake throughput to processing latency and expose clients to transient infrastructure failures.

## Decision

`POST /submissions` persists the submission with a status of `pending` and returns `202 Accepted` immediately with a reference ID. All validation and canonical mapping is performed asynchronously by a leader-elected background worker. Clients poll `GET /submissions/by-reference/{referenceId}/status` for the final outcome. Non-retryable errors (validation failures, draft version, unregistered field type) transition the submission to `rejected` immediately. All other errors are retried with exponential backoff up to a configurable limit before transitioning to `failed`. A `POST /submissions/{submissionId}/replay` endpoint resets a terminal submission to `pending` for reprocessing.

## Consequences

- Intake throughput is decoupled from processing latency; the service can accept submissions under infrastructure degradation.
- Clients must handle an asynchronous response model — a `202` on intake does not indicate acceptance; clients must poll for final status.
- Submission intake is idempotent via a client-supplied `Idempotency-Key` header; safe retries on network failure do not produce duplicate submissions.
- Processing failures are retried automatically without client involvement, with each attempt recorded on the submission as an audit trail.
