 Full Codebase Review: Forms and Tenants Services

**Date:** July 9, 2026

---

## Issues Resolved Since 7/5 Review

1. ~~`Submission.Reject` emits a nil-payload event~~ (`submission.go`) -- `submissionRejectedPayload` struct introduced with `referenceId`, `tenantId`, `formId`, `versionId`, and `reason` (the stringified error passed to `Reject`). `submissionAcceptedPayload` struct introduced in the same commit, replacing the bare `ToFactMap()` marshal with a structured envelope carrying the same identity fields alongside `facts`. Both are marshaled into `json.RawMessage` and passed to `addEvent`. The `// FIXME` comment is removed. _(Carried from prior cycle, P2.)_

2. ~~`EventTypeSubmissionRejected` payload is nil~~ -- Resolved by the above. _(Carried from prior cycle.)_

---

## Will Not Fix

See [6/14 review](code-review-6-14-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`TagPrimitiveType` has no constants and no `isValidTagPrimitiveType` predicate** (`tag_version.go:22`) -- `type TagPrimitiveType string` is declared with no constants and is accepted as a pointer field on `Tag` without any enum validation. Every other constrained type in the domain uses `validate.NewTypeValidator` with an explicit enum. _(Carried from 7/5 #1, P2.)_

2. **`FieldTagMapping` missing `validate` struct tags -- `NewFieldTagMapping` validates nothing** (`field_tag_mapping.go:23-38`) -- `NewFieldTagMapping` calls `validate.ValidateStruct(ftm)` but neither `FieldTagMapping` nor its embedded `FieldTagMappingConfig` carry any `validate` tags. `FieldID` and `TagVersionID` can be empty strings at construction time. _(Carried from 7/5 #2, P2.)_

3. **`submissionFailedPayload` is dead code** (`submission.go`) -- `submissionFailedPayload` struct is declared with `referenceId`, `tenantId`, `formId`, `versionId`, and `reason` fields but `Submission.Fail` never calls `addEvent` -- it only calls `addAttempt`. No `EventTypeSubmissionFailed` constant exists. The struct is unreachable. `Submission.Accept` and `Submission.Reject` both emit domain events; `Fail` is the sole terminal outcome that does not. Either `Submission.Fail` should emit a `failed` event using this payload -- consistent with the other two terminal transitions -- or the struct should be removed. _(New, P3.)_

#### Missing Functionality

4. **Field validator strategies: select and checkbox remain stubs; date partial** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:37`) -- Both select and checkbox return `nil` without performing any validation. Date has `checkValueRequired` but no date-range validation (TODO comment present). Submissions with these field types pass validation unconditionally. _(Carried from 7/5 #4, P3.)_

---

### Tenants Service

#### Missing Functionality

5. **`Find()` has no pagination or filtering** (`tenants_service.go:31-41`) -- _(Carried from 7/5 #5, P3.)_

6. **`Lookup` value object has no validation** (`lookup.go`) -- _(Carried from 7/5 #6, P3.)_

---

### Cross-Service / pkg/

#### Bugs

7. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:189`) -- `onLeader` creates the pool, spawns worker goroutines, and defers `close(pool)`. On context cancellation, `close(pool)` fires before the worker goroutines have observed the cancellation. A worker that reaches `w.WorkerPool <- w.JobChannel` after the channel is closed panics. The existing `sync.WaitGroup` in `Start` synchronizes `onLeader` itself but does not track the goroutines spawned inside it. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second `sync.WaitGroup` inside `onLeader` that waits on worker exits before closing. _(Carried from 7/5 #7, P2.)_

8. **`BigQueryDataLakeClient` is a live-registered stub** (`adapters/clients/big_query_data_lake_client.go:24`) -- `DataSourceTypeDataLake` is fully accepted by the domain, persistence, and REST API layers. Any `data-lake` data source created today will be persisted and processed by the worker, but `Query` unconditionally returns `ErrBigQueryDataLakeNotConfigured`. A `data-lake` data source will silently cycle through worker retries and exhaust the retry limit with no meaningful error message. Options: guard the type at the command/API layer until BigQuery is implemented, or surface `ErrBigQueryDataLakeNotConfigured` as a non-retryable error to avoid retry churn. _(Carried from 7/5 #8, P2.)_

9. **Submission and outbox workers share a single `RetryLimit`** (`workers.go:26-35`) -- Both `newSubmissionsBackgroundWorker` and `newOutboxRelayBackgroundWorker` receive `options.RetryLimit` from the same `WorkerOptions`. For submissions, `RetryLimit` governs how many times the processing pipeline is retried per job. For the outbox relay, it governs how many publish attempts an event can accumulate before being permanently abandoned. These semantics are distinct and operators cannot tune them independently without a code change. _(Carried from 7/5 #9, P2.)_

#### Architectural

10. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service remains well-covered. _(Carried from 7/5 #10, P3.)_

11. **`FindJobs` fetches full submission documents** (`submissions_repository.go:89`) -- Acknowledged by an existing `// TODO` comment. `FindJobs` fetches full `SubmissionDocument` objects and discards everything except the `_id`. A MongoDB projection (`{"_id": 1}`) would avoid deserializing potentially large `Values`/`Facts`/`Attempts` arrays for a query that returns only IDs. _(Carried from 7/5 #11, P3.)_

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `TagPrimitiveType` unconstrained string type; no constants or predicate | Forms |
| **P2** | 2 | `FieldTagMapping` missing `validate` struct tags -- `ValidateStruct` is no-op | Forms |
| **P2** | 7 | `defer close(pool)` shutdown race | pkg/worker |
| **P2** | 8 | `BigQueryDataLakeClient` stub accepts API requests, always fails at processing time | Tenants |
| **P2** | 9 | Submission and outbox workers share `RetryLimit`; semantics are distinct | pkg/worker, Forms |
| **P3** | 3 | `submissionFailedPayload` is dead code; `Submission.Fail` emits no event | Forms |
| **P3** | 4 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 5 | Tenants `Find()` no pagination | Tenants |
| **P3** | 6 | `Lookup` no validation | Tenants |
| **P3** | 10 | Test coverage gaps (forms, pkg/) | Forms, pkg/ |
| **P3** | 11 | `FindJobs` fetches full submission documents; projection missing | Forms |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **9/10 -- Production-Ready** | Both `Submission.Accept` and `Submission.Reject` now emit fully structured domain event payloads carrying identity fields (`referenceId`, `tenantId`, `formId`, `versionId`) alongside outcome-specific data (`facts` or `reason`). The Kafka topic strategy has been simplified: topic is now `AggregateType` alone (`"submission"`) with `eventType` moved to a message header, making the topic stable as new event types are added. `FormVersion` metadata is fully wired end-to-end. No remaining live processing path bugs. The two P2 gaps are construction-time only validation no-ops that do not affect the processing pipeline. |
| **Tenants** | **9/10 -- Production-Ready** | No functional changes since 7/5. `BigQueryDataLakeClient` stub remains the only visible production gap. |
| **pkg/** | **8/10 -- Production-Ready** | `defer close(pool)` race in the background worker unchanged. Zero test coverage across `pkg/`. |

---

## Summary

### Progress Since 7/5

- **`Submission.Reject` event payload completed** -- `submissionRejectedPayload` struct carries `referenceId`, `tenantId`, `formId`, `versionId`, and `reason` (the stringified rejection error). `submissionAcceptedPayload` struct replaces the bare `ToFactMap()` marshal on the accept path with the same identity envelope plus `facts`. Both payloads are marshaled to `json.RawMessage` before being passed to `addEvent`. The `// FIXME` comment is removed and the outbox relay now publishes meaningful structured payloads for both terminal acceptance and rejection outcomes.

- **Kafka topic strategy simplified** -- `KafkaPublisher.Publish` now uses `string(event.AggregateType)` as the topic (e.g. `"submission"`) and moves `event.Type` to a Kafka message header (`"eventType"`). The `topic()` helper function and `strings` import are removed. This is a **breaking change** for any existing consumers subscribed to the previous `"submission.accepted"` / `"submission.rejected"` topic names; they must resubscribe to `"submission"` and filter on the `eventType` header instead.

- **`FormVersion` metadata** -- `Metadata map[string]string` added to the `FormVersion` domain aggregate, `NewFormVersion`, `HydrateFormVersion`, `FormVersionDocument`, `UpsertFormVersionRequest`, `FormVersionResponse`, `CreateFormVersionCommand`, and `UpdateFormVersionCommand`. `FormVersion.Update(metadata) error` added and called from `formsService.UpdateVersion`. The method returns `error` for consistency with other domain mutation methods and to accommodate future validation without breaking callers. Fully wired end-to-end from REST handler through service, domain, and MongoDB document mapper.

### Progress Since 7/5 (inherited context)

See [7/5 review](code-review-7-5-26.md) for full context on the outbox pattern, `SubmissionAttempt` audit trail, `Tag`/`TagVersion` validation enforcement, `tagsService.Delete` version guard, `GetSubmissionFacts` endpoint, and `keypath` regex validator.

### Current State

**11 remaining issues** (5 P2, 6 P3). 0 P0, 0 P1.

All live processing path bugs are resolved. The domain event model is now complete for accepted and rejected submissions. The Kafka topic restructuring makes the event bus more stable -- a single topic per aggregate type with event type discrimination via headers is a better long-term design than per-event-type topics.

**Hexagonal Architecture** -- No structural changes this cycle. The Kafka topic change is entirely contained within the `KafkaPublisher` adapter; the `ports.Publisher` interface and `domain.Event` are unchanged. `FormVersion.Update` is correctly placed on the domain aggregate; the service layer calls it and the handler layer passes `cmd.Metadata` through the command, maintaining the inversion of control.

**DDD** -- `submissionAcceptedPayload` and `submissionRejectedPayload` are package-private structs within the domain, correctly encapsulating the event serialization concern alongside the aggregate that produces it. `Submission.Fail` is the only terminal state transition that does not emit a domain event; `submissionFailedPayload` existing as dead code is the outstanding gap. `TagPrimitiveType` remains the sole unconstrained domain string type. `FieldTagMapping` construction remains the outstanding validation no-op.

**Idiomatic Go** -- `submissionAcceptedPayload` and `submissionRejectedPayload` are unexported types, correctly scoped to the package. The `json.Marshal` error is silently discarded (`p, _ := json.Marshal(...)`) on both payload paths -- this is acceptable given the structs contain only basic string and map types that cannot produce marshal errors, but is worth noting as a pattern to be conscious of when payload types grow more complex.

### Highest-Impact Improvements

1. **Add `validate` struct tags to `FieldTagMapping`** (P2) -- `validate:"required"` on `FieldID` and `TagVersionID`; closes the last construction validation no-op.
2. **Add `TagPrimitiveType` constants and `isValidTagPrimitiveType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTag`; closes the remaining unconstrained domain type gap.
3. **Emit a `failed` event from `Submission.Fail`** (P3) -- use `submissionFailedPayload` (already declared) to emit `EventTypeSubmissionFailed`; makes all three terminal outcomes observable via the event bus and removes the dead code.
4. **Split `RetryLimit` into separate submission and outbox config values** (P2) -- add `OutboxRetryLimit int` to `WorkerOptions` (or a dedicated `OutboxWorkerOptions`); pass each independently to their respective workers.
5. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)` or track worker goroutines in a second `sync.WaitGroup` that `onLeader` waits on before closing.
6. **Guard `DataSourceTypeDataLake` at the API or make `ErrBigQueryDataLakeNotConfigured` non-retryable** (P2) -- prevents `data-lake` data sources from silently consuming all worker retries before failing.
7. **Add projection to `FindJobs`** (P3) -- replace full-document fetch with `{"_id": 1}` projection; reduces unnecessary deserialization of large submission documents in the worker polling path.
8. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data on the live processing pipeline.
9. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3).
