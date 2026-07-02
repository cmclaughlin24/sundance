<div align="center">

# #009 Outbox Pattern: Polling vs. Streaming

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

After a submission is accepted, the Forms Service must reliably publish a canonical submission event to the message broker for downstream consumption. To guarantee delivery without coupling event publishing to the acceptance transaction, the system uses a transactional outbox: the accepted submission record itself serves as the outbox entry, and a relay component is responsible for reading undelivered records and publishing them.

Two relay approaches were considered:

- **Polling** — a background worker ticks on a configurable interval, queries for accepted submissions that have not yet been published, and publishes each one. This is the same model used by the data source refresh worker in the Tenants Service (ADR-004).
- **Streaming (Change Data Capture)** — the relay subscribes to MongoDB change streams, which expose the database's internal oplog as a real-time event feed. The moment a submission document transitions to `accepted`, the change stream fires and the relay publishes the event immediately, with no polling interval.

The project already has a generic `BackgroundWorker[J Job]` with Redis-backed leader election (ADR-004), making a polling relay a natural extension of existing infrastructure. A streaming relay would require subscribing to and managing a MongoDB change stream cursor — a distinct operational pattern not currently used anywhere in the system.

## Decision

The outbox relay is implemented as a polling-based background worker using the existing `BackgroundWorker[J Job]` from `pkg/worker`. On each tick the worker queries for accepted submissions that have not yet been published and dispatches each as a job. This is the chosen approach for the initial release, prioritising consistency with the existing worker pattern and operational simplicity over real-time event delivery.

## Consequences

- The relay is consistent with the existing background worker pattern; no new infrastructure, driver configuration, or operational concepts are introduced.
- The relay inherits all worker guarantees from ADR-004: Redis-backed leader election ensures a single active relay across replicas, graceful 30-second shutdown drains in-flight publish jobs, and the `InMemoryElector` removes the Redis dependency for local development and testing.
- Events are not published in real time — delivery latency is bounded by the worker's tick interval rather than by the time of acceptance.
- At high submission volume, each tick issues a database query regardless of whether new records exist, adding polling overhead that a change stream would avoid.
- If event delivery latency or database polling overhead becomes a concern, a streaming relay backed by MongoDB change streams should be evaluated as a replacement. Change streams would reduce latency to near-zero and eliminate the per-tick query cost, at the expense of cursor lifecycle management and the requirement for a MongoDB replica set oplog (already satisfied in production per ADR-006).
