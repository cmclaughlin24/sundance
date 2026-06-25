<div align="center">

# #004 Generic Background Worker with Redis Leader Election

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Both services require a background worker for periodic job processing — the Tenants Service refreshes expired scheduled data source caches, and the Forms Service processes pending submissions. Both services are horizontally scalable, meaning multiple replicas may run simultaneously. Only one replica should process jobs at a time to avoid duplicate processing. A reusable worker implementation was preferred over duplicating the polling, concurrency, and leader election logic in each service.

## Decision

A single generic `BackgroundWorker[J Job]` is implemented in `pkg/worker` and reused by both services. The worker manages a ticker-driven loop, leader election state, job fetching, and a goroutine pool. The consuming service provides only a `FetchJobsFn` and a `Job` implementation. Leader election is handled by a `CacheElector` backed by Redis `SetNX` + Lua scripts with a configurable TTL and renewal interval, ensuring only one replica runs the worker at a time across horizontal scale. An `InMemoryElector` is provided as a no-op for local development and testing, removing the Redis dependency outside of production.

This is a deliberate tradeoff against a fully distributed approach — such as a dedicated message queue or job scheduler — which would allow multiple replicas to process jobs concurrently. The leader election model is simpler to operate and sufficient for current throughput requirements, but may need to be revisited if processing latency or job queue depth becomes a bottleneck at scale.

## Consequences

- Both services share a single, consistently behaved worker implementation; changes to retry logic, shutdown behaviour, or leader election apply uniformly.
- Redis is a hard runtime dependency in production for both services. Loss of Redis connectivity will cause all replicas to fail leader election and stop processing jobs.
- The worker provides graceful shutdown with a 30-second drain window, allowing in-flight jobs to complete before the process exits.
- New job types can be introduced by implementing the `Job` interface and providing a `FetchJobsFn`; no changes to the worker itself are required.
- If submission volume grows significantly, the single-leader model may become a throughput bottleneck and a distributed queue approach should be evaluated.
