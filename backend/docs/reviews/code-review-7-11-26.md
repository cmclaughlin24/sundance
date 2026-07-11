# Full Codebase Review: Forms and Tenants Services

**Date:** July 11, 2026

---

## Issues Resolved Since 7/9 Review

1. ~~`defer close(pool)` race on shutdown~~ (`pkg/worker/background_worker.go`) -- `background_worker.go` (256 lines) deleted and replaced with three focused files: `distributed_worker.go`, `periodic_worker.go`, and `type.go`. `DistributedWorker.Start` now maintains a `sync.WaitGroup` tracking the `onLeader` goroutine and performs a 30-second drain wait on `ctx.Done()` before returning. See issue #7 below for the residual race that remains in both new types. _(Carried from 7/9 #7, P2 -- partially resolved.)_

2. ~~Outbox relay runs on leader-elected worker~~ -- Outbox pattern changed from leader-elected to claim-based (`830e9574`). `OutboxRepository.Claim()` now uses an atomic `FindOneAndUpdate` with a `locked_until` field and `LeaseDuration`. Events in `processing` status with an expired `locked_until` are re-eligible, preventing double-processing across instances without requiring Redis leader election. `$unset locked_until` on `Upsert` correctly releases the lease after processing. The outbox relay now runs on `PeriodicWorker` (no election) rather than `DistributedWorker`. _(New resolution, architectural improvement.)_

---

## Will Not Fix

See [6/14 review](code-review-6-14-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`TagPrimitiveType` has no constants and no `isValidTagPrimitiveType` predicate** (`tag_version.go:22`) -- `type TagPrimitiveType string` is declared with no constants and is accepted as a pointer field on `Tag` without any enum validation. Every other constrained type in the domain uses `validate.NewTypeValidator` with an explicit enum. _(Carried from 7/9 #1, P2.)_

2. **`FieldTagMapping` missing `validate` struct tags -- `NewFieldTagMapping` validates nothing** (`field_tag_mapping.go:23-38`) -- `NewFieldTagMapping` calls `validate.ValidateStruct(ftm)` but neither `FieldTagMapping` nor its embedded `FieldTagMappingConfig` carry any `validate` tags. `FieldID` and `TagVersionID` can be empty strings at construction time. _(Carried from 7/9 #2, P2.)_

3. **`submissionFailedPayload` is dead code** (`submission.go:233-239`) -- `submissionFailedPayload` struct is declared with `referenceId`, `tenantId`, `formId`, `versionId`, and `reason` fields but `Submission.Fail` never calls `addEvent` -- it only calls `addAttempt`. No `EventTypeSubmissionFailed` constant exists. The struct is unreachable. `Submission.Accept` and `Submission.Reject` both emit domain events; `Fail` is the sole terminal outcome that does not. Either `Submission.Fail` should emit a `failed` event using this payload -- consistent with the other two terminal transitions -- or the struct should be removed. _(Carried from 7/9 #3, P3.)_

#### Missing Functionality

4. **Field validator strategies: select and checkbox remain stubs; date partial** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:37`) -- Both select and checkbox return `nil` without performing any validation. Date has `checkValueRequired` but no date-range validation (TODO comment present). Submissions with these field types pass validation unconditionally. _(Carried from 7/9 #4, P3.)_

---

### Tenants Service

#### Missing Functionality

5. **`Find()` has no pagination or filtering** (`tenants_service.go:31-41`) -- _(Carried from 7/9 #5, P3.)_

6. **`Lookup` value object has no validation** (`lookup.go`) -- _(Carried from 7/9 #6, P3.)_

---

### Cross-Service / pkg/

#### Bugs

7. **`defer close(pool)` shutdown race persists in both new worker types** (`pkg/worker/periodic_worker.go:72`, `pkg/worker/distributed_worker.go:186`) -- The refactor resolved the `sync.WaitGroup` gap at the `DistributedWorker.Start` level (the `onLeader` goroutine is now tracked and drained with a 30-second timeout on shutdown). However, the core race is still present inside both `PeriodicWorker.Start` and `DistributedWorker.onLeader`: each creates a `pool` channel, spawns `Worker[J]` goroutines via `w.Start(ctx)`, and defers `close(pool)`. The spawned `Worker[J]` goroutines are fire-and-forget with no `WaitGroup` tracking them. On context cancellation, `close(pool)` fires before the inner worker goroutines have exited. A goroutine racing to send `w.WorkerPool <- w.JobChannel` after `close(pool)` will panic. Suggest either dropping `defer close(pool)` (let GC reclaim the channel) or tracking inner worker goroutines in a `WaitGroup` that `onLeader`/`Start` waits on before closing. _(Carried from 7/9 #7, P2.)_

8. **`BigQueryDataLakeClient` is a live-registered stub** (`adapters/clients/big_query_data_lake_client.go:24`) -- `DataSourceTypeDataLake` is fully accepted by the domain, persistence, and REST API layers. Any `data-lake` data source created today will be persisted and processed by the worker, but `Query` unconditionally returns `ErrBigQueryDataLakeNotConfigured`. A `data-lake` data source will silently cycle through worker retries and exhaust the retry limit with no meaningful error message. Options: guard the type at the command/API layer until BigQuery is implemented, or surface `ErrBigQueryDataLakeNotConfigured` as a non-retryable error to avoid retry churn. _(Carried from 7/9 #8, P2.)_

9. **Submission and outbox workers share a single `RetryLimit`** (`workers.go:26-35`) -- Both `newSubmissionsDistributedWorker` and `newOutboxRelayPeriodicWorker` receive `options.RetryLimit` from the same `WorkerOptions`. For submissions, `RetryLimit` governs how many times the processing pipeline is retried per job. For the outbox relay, it governs how many publish attempts an event can accumulate before being permanently abandoned. These semantics are distinct and operators cannot tune them independently without a code change. _(Carried from 7/9 #9, P2.)_

#### Architectural

10. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service remains well-covered. _(Carried from 7/9 #10, P3.)_

11. **`FindJobs` fetches full submission documents** (`submissions_repository.go:89`) -- Acknowledged by an existing `// TODO` comment. `FindJobs` fetches full `SubmissionDocument` objects and discards everything except the `_id`. A MongoDB projection (`{"_id": 1}`) would avoid deserializing potentially large `Values`/`Facts`/`Attempts` arrays for a query that returns only IDs. _(Carried from 7/9 #11, P3.)_

#### Observations

12. **`PeriodicWorker` has no backoff under sustained failure** (`pkg/worker/periodic_worker.go`) -- `PeriodicWorker.work()` errors are logged and silently continued at the same fixed interval. For the outbox relay, sustained Kafka unavailability will result in repeated claim-attempt-fail cycles at the configured cadence. The claim-based lease model prevents event loss (events return to eligibility after `locked_until` expires), so this is not a correctness issue. A simple fixed cooldown or exponential backoff after N consecutive failures would reduce unnecessary load on the broker during outage periods. _(New, P4.)_

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `TagPrimitiveType` unconstrained string type; no constants or predicate | Forms |
| **P2** | 2 | `FieldTagMapping` missing `validate` struct tags -- `ValidateStruct` is no-op | Forms |
| **P2** | 7 | `defer close(pool)` shutdown race persists in `PeriodicWorker` and `DistributedWorker.onLeader` | pkg/worker |
| **P2** | 8 | `BigQueryDataLakeClient` stub accepts API requests, always fails at processing time | Tenants |
| **P2** | 9 | Submission and outbox workers share `RetryLimit`; semantics are distinct | pkg/worker, Forms |
| **P3** | 3 | `submissionFailedPayload` is dead code; `Submission.Fail` emits no event | Forms |
| **P3** | 4 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 5 | Tenants `Find()` no pagination | Tenants |
| **P3** | 6 | `Lookup` no validation | Tenants |
| **P3** | 10 | Test coverage gaps (forms, pkg/) | Forms, pkg/ |
| **P3** | 11 | `FindJobs` fetches full submission documents; projection missing | Forms |
| **P4** | 12 | `PeriodicWorker` no backoff under sustained failure | pkg/worker |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **9/10 -- Production-Ready** | No regressions. Claim-based outbox is a meaningful improvement in resilience over the leader-elected approach: the outbox relay no longer requires Redis coordination, and events cannot be double-processed even when multiple instances are running. The `locked_until` / `$unset` lifecycle is correctly implemented end-to-end. No new live-path bugs introduced this cycle. The two P2 gaps remain construction-time validation no-ops that do not affect the processing pipeline. |
| **Tenants** | **9/10 -- Production-Ready** | No functional changes since 7/9. `BigQueryDataLakeClient` stub remains the only visible production gap. |
| **pkg/** | **8/10 -- Production-Ready** | Worker package refactored into well-scoped `DistributedWorker` and `PeriodicWorker` types with clear separation of concerns. `DistributedWorker.Start` now correctly drains `onLeader` on shutdown with a 30-second timeout. The `defer close(pool)` panic risk remains inside both `PeriodicWorker.Start` and `DistributedWorker.onLeader` for the inner `Worker[J]` goroutines. Zero test coverage across `pkg/` unchanged. |

---

## Summary

### Progress Since 7/9

- **Worker package refactored** -- `background_worker.go` deleted and replaced with `distributed_worker.go`, `periodic_worker.go`, and `type.go`. The split correctly separates the leader-election concern (`DistributedWorker`) from simple interval polling (`PeriodicWorker`). `FetchJobsFn[J]` and the sentinel errors are extracted to `type.go`. `DistributedWorker.Start` now manages `onLeader` via a `sync.WaitGroup` and performs a 30-second drain wait on shutdown, closing the previous gap where `onLeader` had no lifecycle tracking at the `Start` level.

- **Outbox relay migrated to claim-based processing** -- `OutboxRepository.Claim()` uses `FindOneAndUpdate` to atomically set `status: processing` and `locked_until: now + LeaseDuration (5m)` on each claimed event. Events with an expired `locked_until` are re-eligible, providing at-least-once delivery without Redis leader election. `Upsert` uses `$unset locked_until` to release the lease on completion or error. The outbox relay worker now runs on `PeriodicWorker` -- appropriate since claim-based mutual exclusion at the data layer replaces the need for distributed leader election at the worker layer.

### Current State

**12 remaining issues** (5 P2, 6 P3, 1 P4). 0 P0, 0 P1.

The claim-based outbox is the most architecturally significant change since 7/9. Removing the Redis dependency from the outbox relay path simplifies the operational topology: the relay can safely run on every instance without coordination overhead, and the `locked_until` lease provides the mutual exclusion guarantee that previously required leader election.

**Hexagonal Architecture** -- The `PeriodicWorker` and `DistributedWorker` split is a clean internal refactor of `pkg/worker` with no changes to any port or domain interface. The `OutboxRepository.Claim()` signature change (`ClaimEventsOptions` now includes `LeaseDuration` and `CreatedAfter`) is a port-level change that is correctly absorbed at the adapter layer; the service layer passes `ClaimEventsOptions` through the `newOutboxWorkFn` closure without coupling to the MongoDB implementation.

**DDD** -- No domain changes this cycle. `submissionFailedPayload` remains the only dead domain struct. `TagPrimitiveType` and `FieldTagMapping` validation gaps are unchanged.

**Idiomatic Go** -- `DistributedWorker` and `PeriodicWorker` use the functional options pattern (`opts ...func(*T)`) consistently with the existing `Worker[J]` type. The empty `var ( )` block in `distributed_worker.go` introduced during the refactor has been removed. The `wg.Go` call (from `sync.WaitGroup` extended API, Go 1.22+) in `DistributedWorker.Start` is idiomatic for tracking a single goroutine with panic propagation.

### Highest-Impact Improvements

1. **Fix the inner worker goroutine shutdown race** (P2) -- drop `defer close(pool)` in `PeriodicWorker.Start` and `DistributedWorker.onLeader`, or track `Worker[J]` goroutines in a second `WaitGroup` and wait on them before closing; closes the panic risk under concurrent shutdown.
2. **Add `validate` struct tags to `FieldTagMapping`** (P2) -- `validate:"required"` on `FieldID` and `TagVersionID`; closes the last construction validation no-op.
3. **Add `TagPrimitiveType` constants and `isValidTagPrimitiveType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTag`; closes the remaining unconstrained domain type gap.
4. **Split `RetryLimit` into separate submission and outbox config values** (P2) -- add `OutboxRetryLimit int` to `WorkerOptions` (or a dedicated `OutboxWorkerOptions`); pass each independently to their respective workers.
5. **Guard `DataSourceTypeDataLake` at the API or make `ErrBigQueryDataLakeNotConfigured` non-retryable** (P2) -- prevents `data-lake` data sources from silently consuming all worker retries before failing.
6. **Emit a `failed` event from `Submission.Fail`** (P3) -- use `submissionFailedPayload` (already declared) to emit `EventTypeSubmissionFailed`; makes all three terminal outcomes observable via the event bus and removes the dead code.
7. **Add projection to `FindJobs`** (P3) -- replace full-document fetch with `{"_id": 1}` projection; reduces unnecessary deserialization of large submission documents in the worker polling path.
8. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data on the live processing pipeline.
9. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3).
