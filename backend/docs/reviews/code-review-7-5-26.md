# Full Codebase Review: Forms and Tenants Services

**Date:** July 5, 2026

---

## Issues Resolved Since 6/28 Review

1. ~~No `SubmissionAttempt` ever created -- `Reject`/`Fail` error params unused~~ (`submission.go`) -- `addAttempt` private method introduced; `Accept`, `Fail`, and `Reject` all call `s.addAttempt(len(s.Attempts)+1, string(status), err)`. The persistence mapper (`documents/submission_attempt.go`) and `Submission.Attempts` field were already complete and are now reachable code. _(Carried from 6/28 #1, P2.)_

2. ~~`ErrNoEligibleTagVersion` and `ErrMultipleActiveTagVersion` declared with empty error messages~~ (`tag_version.go:14-15`) -- Both now wrap `ErrInvalidTagVersion` with descriptive messages: `"no eligible tag versions"` and `"multiple active tag versions"`. _(Carried from 6/28 #6, P2.)_

3. ~~No domain events for cross-service communication~~ -- Outbox pattern fully implemented: `domain.Event` aggregate with `withEvents` mixin (`HasEvents` interface, `PeekEvents()`/`DrainEvents()`), `OutboxRepository` port with MongoDB and in-memory adapters, `KafkaPublisher` and `InMemoryPublisher`, and `outboxRelayBackgroundWorker`. `Submission.Accept` emits `EventTypeSubmissionAccepted` with a `ToFactMap()` payload serialized into `json.RawMessage`. _(Carried from 6/28 #14, P3 -- partially resolved, see #13 below.)_

4. ~~`mongoDBTagsRespository` typo~~ (`tags_repository.go:28`) -- Struct renamed to `mongoDBTagsRepository`; all receiver method declarations updated. Stale log messages `"upsert canonical tag"` / `"canonical_tag_id"` corrected to `"upsert tag"` / `"tag_id"` in the same change.

5. ~~`evaluateCollectionCandidates` dereferences `collectionIndex` without a nil guard~~ (`submission_jobs_service.go`) -- `ErrMissingCollectionIndex` sentinel introduced at the service level. The nil guard now returns the error rather than silently skipping the candidate. `shouldReject` updated to classify `ErrMissingCollectionIndex` as a rejection: a collection-tagged field value that omits `CollectionIndex` is structurally invalid input from the submitter, not an infrastructure failure. `evaluateCollectionCandidates` and `evaluateScalarCandidates` both updated to return `([]*domain.CanonicalFact, error)` for signature consistency. _(Carried from 6/28 #5, P2.)_

---

## Will Not Fix

See [6/14 review](code-review-6-14-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`Tag` has no `validate` struct tags despite calling `validate.ValidateStruct`** (`tag.go:29-38`) -- `Tag` fields (`TenantID`, `KeyPath`, `DisplayName`) carry no `validate:"required"` tags. `NewTag` calls `validate.ValidateStruct` but the call is a no-op on an untagged struct. An empty `tenantID`, empty `keyPath`, or empty `displayName` passes construction without error. `TagVersion` fields are similarly untagged. Compare to `Form`, `FormVersion`, `Submission`, and `DataSource`, which all tag their required fields. _(Carried from 6/28 #2, P2.)_

2. **`TagPrimitiveType` has no constants and no `isValidTagPrimitiveType` predicate** (`tag_version.go:22`) -- `type TagPrimitiveType string` is declared with no constants and is accepted as a pointer field on `Tag` without any enum validation. Every other constrained type in the domain uses `validate.NewTypeValidator` with an explicit enum. _(Carried from 6/28 #3, P2.)_

3. **`FieldTagMapping` missing `validate` struct tags -- `NewFieldTagMapping` validates nothing** (`field_tag_mapping.go:23-38`) -- `NewFieldTagMapping` calls `validate.ValidateStruct(ftm)` but neither `FieldTagMapping` nor its embedded `FieldTagMappingConfig` carry any `validate` tags. `FieldID` and `TagVersionID` can be empty strings at construction time. _(Carried from 6/28 #4, P2.)_

4. **`Submission.Reject` emits a nil-payload event** (`submission.go:142-143`) -- `s.addEvent(EventTypeSubmissionRejected, nil)` is called with a `// FIXME: Create event payload.` comment. The outbox relay worker publishes this event to Kafka with `message.Value = event.Payload`, which will be a zero-length byte slice. Consumers expecting a rejection payload receive nothing actionable. `EventTypeSubmissionAccepted` correctly serializes `s.ToFactMap()` as its payload; the `Reject` path is half-implemented. _(New, P2.)_

#### Missing Functionality

5. **Field validator strategies: select and checkbox remain stubs; date partial** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:37`) -- Both select and checkbox return `nil` without performing any validation. Date has `checkValueRequired` but no date-range validation (TODO comment present). Submissions with these field types pass validation unconditionally. _(Carried from 6/28 #7, P3.)_

6. **`tagsService.Delete` missing active-version guard** (`tags_service.go:139`) -- A `// FIXME` comment acknowledges the missing invariant. Deleting a tag with a historically active version should be prevented to preserve audit history, consistent with the guard on `formsService.Delete` via `hasActiveVersion`. Without it, a tag backing live `FieldTagMapping` records can be hard-deleted. _(Carried from 6/28 #8, P3.)_

---

### Tenants Service

#### Missing Functionality

7. **`Find()` has no pagination or filtering** (`tenants_service.go:31-41`) -- _(Carried from 6/28 #9, P3.)_

8. **`Lookup` value object has no validation** (`lookup.go`) -- _(Carried from 6/28 #10, P3.)_

---

### Cross-Service / pkg/

#### Bugs

9. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:189`) -- `onLeader` creates the pool, spawns worker goroutines, and defers `close(pool)`. On context cancellation, `close(pool)` fires before the worker goroutines have observed the cancellation. A worker that reaches `w.WorkerPool <- w.JobChannel` after the channel is closed panics. The existing `sync.WaitGroup` in `Start` synchronizes `onLeader` itself but does not track the goroutines spawned inside it. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second `sync.WaitGroup` inside `onLeader` that waits on worker exits before closing. _(Carried from 6/28 #11, P2.)_

10. **`BigQueryDataLakeClient` is a live-registered stub** (`adapters/clients/big_query_data_lake_client.go:24`) -- `DataSourceTypeDataLake` is fully accepted by the domain, persistence, and REST API layers. Any `data-lake` data source created today will be persisted and processed by the worker, but `Query` unconditionally returns `ErrBigQueryDataLakeNotConfigured`. A `data-lake` data source will silently cycle through worker retries and exhaust the retry limit with no meaningful error message. Options: guard the type at the command/API layer until BigQuery is implemented, or surface `ErrBigQueryDataLakeNotConfigured` as a non-retryable error to avoid retry churn. _(Carried from 6/28 #12, P2.)_

11. **Submission and outbox workers share a single `RetryLimit`** (`workers.go:26-35`) -- Both `newSubmissionsBackgroundWorker` and `newOutboxRelayBackgroundWorker` receive `options.RetryLimit` from the same `WorkerOptions`. For submissions, `RetryLimit` governs how many times the processing pipeline is retried per job. For the outbox relay, it governs how many publish attempts an event can accumulate before being permanently abandoned. These semantics are distinct and operators cannot tune them independently without a code change. _(New, P2.)_

#### Architectural

12. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service remains well-covered. _(Carried from 6/28 #13, P3.)_

13. **`FindJobs` fetches full submission documents** (`submissions_repository.go:89`) -- Acknowledged by an existing `// TODO` comment. `FindJobs` fetches full `SubmissionDocument` objects and discards everything except the `_id`. A MongoDB projection (`{"_id": 1}`) would avoid deserializing potentially large `Values`/`Facts`/`Attempts` arrays for a query that returns only IDs. _(New, P3.)_

14. **`EventTypeSubmissionRejected` payload is nil** -- Covered under bug #4 above; the domain event infrastructure is otherwise complete. _(New, partially resolves 6/28 #14.)_

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `Tag`/`TagVersion` missing `validate` struct tags -- `ValidateStruct` is no-op | Forms |
| **P2** | 2 | `TagPrimitiveType` unconstrained string type; no constants or predicate | Forms |
| **P2** | 3 | `FieldTagMapping` missing `validate` struct tags -- `ValidateStruct` is no-op | Forms |
| **P2** | 4 | `Submission.Reject` emits nil-payload event | Forms |
| **P2** | 9 | `defer close(pool)` shutdown race | pkg/worker |
| **P2** | 10 | `BigQueryDataLakeClient` stub accepts API requests, always fails at processing time | Tenants |
| **P2** | 11 | Submission and outbox workers share `RetryLimit`; semantics are distinct | pkg/worker, Forms |
| **P3** | 5 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 6 | `tagsService.Delete` missing active-version guard | Forms |
| **P3** | 7 | Tenants `Find()` no pagination | Tenants |
| **P3** | 8 | `Lookup` no validation | Tenants |
| **P3** | 12 | Test coverage gaps (forms, pkg/) | Forms, pkg/ |
| **P3** | 13 | `FindJobs` fetches full submission documents; projection missing | Forms |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8/10 -- Production-Ready** | Outbox pattern and domain events are the most significant addition since 6/28. `SubmissionAttempt` audit trail is now live end-to-end: `Accept`, `Fail`, and `Reject` all call `addAttempt`, the document mapper serializes it, and the submission is persisted in a transaction that also writes outbox events atomically. The `EventTypeSubmissionAccepted` path is complete. `ErrMissingCollectionIndex` closes the nil deref panic risk in the collection candidate path and correctly routes malformed submissions to `Rejected`. One new P2: `Submission.Reject` emits a nil-payload event, leaving consumers with no rejection context. One new P2: submission and outbox workers share a single `RetryLimit` that serves two semantically different purposes. Three carried P2 validation no-ops remain. |
| **Tenants** | **9/10 -- Production-Ready** | No functional changes since 6/28. `BigQueryDataLakeClient` stub remains the only visible production gap. |
| **pkg/** | **8/10 -- Production-Ready** | `defer close(pool)` race in the background worker unchanged. Zero test coverage across `pkg/`. |

---

## Summary

### Progress Since 6/28

- **Outbox / domain event pattern implemented** -- `domain.Event` aggregate with `withEvents` mixin added to `Submission`. `Submission.Accept` serializes `ToFactMap()` into a `json.RawMessage` payload and emits `EventTypeSubmissionAccepted`. `Submission.Reject` emits `EventTypeSubmissionRejected` with a nil payload (see bug #4). `mongoDBSubmissionsRepository.Upsert` calls `WriteEvents` inside the same MongoDB session context as the submission write, providing atomic outbox delivery. The `outboxRelayBackgroundWorker` polls the outbox on the same interval/pool/elector pattern as the submissions worker, dispatches each pending or errored event to `publisher.Publish`, then upserts the result back to the outbox. `KafkaPublisher` and `InMemoryPublisher` both satisfy the `ports.Publisher` secondary port. The composition root in `main.go` bootstraps the publisher from settings and both workers are started from the same `workers.Bootstrap` call.

- **`SubmissionAttempt` wired** -- `addAttempt` private method on `Submission` appends a `NewSubmissionAttempt` record on every terminal outcome. All three paths -- `Accept`, `Fail`, `Reject` -- now record the attempt number, result status string, and error details. The persistence mapper round-trips `Attempts` through `submissionAttemptDocument`. This closes the longest-standing P2 gap in the forms service.

- **`evaluateCollectionCandidates` nil deref fixed** -- `ErrMissingCollectionIndex` sentinel declared at the service level. Both `evaluateCollectionCandidates` and `evaluateScalarCandidates` updated to return `([]*domain.CanonicalFact, error)`. The nil guard returns `ErrMissingCollectionIndex` immediately rather than silently skipping the candidate, preventing an incomplete fact set from being accepted. `shouldReject` updated to classify it as a rejection: missing `CollectionIndex` on a collection-tagged field value is structurally invalid submitter input.

- **`mongoDBTagsRespository` typo corrected** -- Struct and all receiver methods renamed to `mongoDBTagsRepository`. Stale log attributes `"upsert canonical tag"` / `"canonical_tag_id"` updated to `"upsert tag"` / `"tag_id"` for consistency with the rest of the repository layer.

### Current State

**13 remaining issues** (7 P2, 6 P3). 0 P0, 0 P1.

The outbox implementation is the most consequential change since 6/28. The transactional write pattern -- submission upsert and outbox event write sharing the same MongoDB session context -- is architecturally sound and consistent with the Debezium outbox pattern referenced in the code comments. The relay worker correctly re-drives failed publishes using the `EventStatusError` filter and respects the `RetryLimit` and `CreatedAfter` window.

**Hexagonal Architecture** -- `ports.Publisher` is a clean secondary port in the core package. Both `KafkaPublisher` and `InMemoryPublisher` are in the adapters layer. The `OutboxRepository` follows the same pattern as all other repositories. The `withEvents` mixin and `HasEvents` interface are domain-layer constructs with no adapter dependencies. The composition root correctly wires the publisher into `core.Application` via `WithPublisher`. `WriteEvents` is a method on the concrete repository structs rather than on the `SubmissionsRepository` port itself -- this is acceptable since it is called from within the adapter, not from the service layer.

**DDD** -- `Submission` now correctly aggregates its own audit trail (`Attempts`) and its own domain events (`withEvents`). The `addAttempt` and `addEvent` helpers are private, enforcing that status transitions are the only path to producing these records. `ErrMissingCollectionIndex` is correctly placed at the service layer: the nil check cannot be enforced earlier since `CollectionIndex` validity is only knowable once the tag mapping is resolved during processing. `TagPrimitiveType` remains the sole unconstrained domain string type. The three `validate` struct tag no-ops (`Tag`, `TagVersion`, `FieldTagMapping`) are the outstanding DDD construction invariant gaps.

**Idiomatic Go** -- The `withEvents` embedding pattern using `iter.Seq[Event]` for `PeekEvents` is idiomatic Go 1.23+. `DrainEvents` using `we.events = nil` correctly releases the backing array. The `HasEvents` interface is minimal and testable. Updating both `evaluateCollectionCandidates` and `evaluateScalarCandidates` to return `error` keeps the `evalFn` function variable type uniform -- a clean approach that avoids a type switch on the dispatch path.

### Highest-Impact Improvements

1. **Complete `Submission.Reject` event payload** (P2) -- serialize a rejection context struct (at minimum `submissionID`, `status`, `errorDetails`) into `json.RawMessage` and pass it to `addEvent`; removes the `// FIXME` and makes rejection events actionable for consumers.
2. **Add `validate` struct tags to `Tag`, `TagVersion`, and `FieldTagMapping`** (P2) -- one-line fix per required field; immediately enforces construction-level invariants on the live tag and field-mapping API paths.
3. **Add `TagPrimitiveType` constants and `isValidTagPrimitiveType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTag`; closes the remaining unconstrained domain type gap.
4. **Split `RetryLimit` into separate submission and outbox config values** (P2) -- add `OutboxRetryLimit int` to `WorkerOptions` (or a dedicated `OutboxWorkerOptions`); pass each independently to their respective workers.
5. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)` or track worker goroutines in a second `sync.WaitGroup` that `onLeader` waits on before closing.
6. **Guard `DataSourceTypeDataLake` at the API or make `ErrBigQueryDataLakeNotConfigured` non-retryable** (P2) -- prevents `data-lake` data sources from silently consuming all worker retries before failing.
7. **Add active-version guard to `tagsService.Delete`** (P3) -- mirrors the existing `hasActiveVersion` check on `formsService.Delete`; FIXME comment already identifies the invariant.
8. **Add projection to `FindJobs`** (P3) -- replace full-document fetch with `{"_id": 1}` projection; reduces unnecessary deserialization of large submission documents in the worker polling path.
9. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data on the live processing pipeline.
10. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3).
