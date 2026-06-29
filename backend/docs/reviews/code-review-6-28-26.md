# Full Codebase Review: Forms and Tenants Services

---

## Issues Resolved Since 6/14 Review

1. ~~`SetRequestDateContext` stored the date value under `correlationIDKey` instead of `requestDateKey`~~ (`pkg/common/httputil/request_date_middleware.go`) -- Unstaged: parameter renamed from `correlationID` to `requestDate`, context key corrected from `correlationIDKey` to `requestDateKey`. Both services register `NewRequestDateMiddleware("X-Request-Date")`; the header value now correctly populates `RequestDateFromContext` downstream.

2. ~~`rankCandiates` typo~~ (`submission_jobs_service.go`) -- Unstaged: function renamed to `rankCandidates` at declaration and at both call sites in `evaluateCollectionCandidates` and `evaluateScalarCandidates`.

3. ~~Tenants `main.go` embeds `\n` in `fmt.Sprintf` log format strings~~ (`tenants/cmd/server/main.go`) -- Unstaged: trailing `\n` removed from all three format strings (`server error`, `application received shutdown signal`, `application shutdown failed`). JSON structured log `msg` fields are now clean. Forms `main.go` did not have this issue.

---

## Will Not Fix

See [6/14 review](code-review-6-14-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **No `SubmissionAttempt` ever created -- `Reject`/`Fail` error params unused** (`submission_jobs_service.go:290-319`) -- `recordAttempt` transitions status and persists the submission but never calls `NewSubmissionAttempt` or appends to `submission.Attempts`. `Submission.Fail(err error)` and `Submission.Reject(err error)` accept an `err` parameter but discard it entirely (`submission.go:116-123`). The domain type, constructor, `Submission.Attempts` field, and mongo document mapper all exist as dead code. With the retry loop in place, operators have no audit trail of per-attempt outcomes or the errors that caused them. _(Carried from 6/14 #1, P2.)_

2. **`Tag` has no `validate` struct tags despite calling `validate.ValidateStruct`** (`tag.go:27-36`) -- `Tag` fields (`TenantID`, `KeyPath`, `DisplayName`) carry no `validate:"required"` tags. `NewTag` calls `validate.ValidateStruct` but the call is a no-op on an untagged struct. An empty `tenantID`, empty `keyPath`, or empty `displayName` passes construction without error. Note: `NodeType` is now validated via the explicit `isTagValueKind` predicate before construction, closing that gap. `TagVersion` fields are similarly untagged. Compare to `Form`, `FormVersion`, `Submission`, and `DataSource`, which all tag their required fields. _(Carried from 6/14 #2, P2.)_

3. **`TagPrimitiveType` has no constants and no `isValidTagPrimitiveType` predicate** (`tag_version.go:22`) -- `type TagPrimitiveType string` is declared with no constants and is accepted as a pointer field on `Tag` without any enum validation. The prior `TagType` gap is closed -- `TagNodeType` now has `isTagValueKind` -- but `TagPrimitiveType` remains the sole unconstrained domain string type. Any string, including empty, is accepted. Every other constrained type in the domain uses `validate.NewTypeValidator` with an explicit enum. _(Updated from 6/14 #3, P2.)_

4. **`FieldTagMapping` missing `validate` struct tags -- `NewFieldTagMapping` validates nothing** (`field_tag_mapping.go:23-38`) -- `NewFieldTagMapping` calls `validate.ValidateStruct(ftm)` but neither `FieldTagMapping` nor its embedded `FieldTagMappingConfig` carry any `validate` tags. `FieldID` and `TagVersionID` can be empty strings at construction time. _(Carried from 6/14 #4, P2.)_

5. **`evaluateCollectionCandidates` dereferences `collectionIndex` without a nil guard** (`submission_jobs_service.go:350`) -- `byCollectionIdx[*fc.collectionIndex]` panics at runtime if any `factCandidate` for a collection-ancestor tag has a nil `CollectionIndex`. The scalar path in `evaluateScalarCandidates` correctly passes `nil`; only the collection path unconditionally dereferences. A submission against a form with collection-tagged fields that omits `CollectionIndex` on any field value will panic the worker goroutine. The recover in `Worker.process` catches the panic but the submission is left perpetually pending with no attempt recorded and no error surfaced. _(New, P2.)_

6. **`ErrNoEligibleTagVersion` and `ErrMultipleActiveTagVersion` declared with empty error messages** (`tag_version.go:14-15`) -- Both are `errors.New("")`. `ErrMultipleActiveTagVersion` propagates from `ResolveTagVersion` into submission processing via `normalize`; `ErrNoEligibleTagVersion` surfaces when a submission references a tag with no active or deprecated version. Empty strings produce completely uninformative log entries and error responses on both paths. _(New, P2.)_

#### Missing Functionality

7. **Field validator strategies: select and checkbox remain stubs; date partial** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:37`) -- Both select and checkbox return `nil` without performing any validation. Date has `checkValueRequired` but no date-range validation (TODO comment present). Submissions with these field types pass validation unconditionally. _(Carried from 6/14 #5, P3.)_

8. **`tagsService.Delete` missing active-version guard** (`tags_service.go:139`) -- A `// FIXME` comment acknowledges the missing invariant. Deleting a tag with a historically active version should be prevented to preserve audit history, consistent with the guard on `formsService.Delete` via `hasActiveVersion`. Without it, a tag backing live `FieldTagMapping` records can be hard-deleted. _(Carried from 6/14 #6, P3.)_

---

### Tenants Service

#### Missing Functionality

9. **`Find()` has no pagination or filtering** (`tenants_service.go:31-41`) -- _(Carried from 6/14 #7, P3.)_

10. **`Lookup` value object has no validation** (`lookup.go`) -- _(Carried from 6/14 #8, P3.)_

---

### Cross-Service / pkg/

#### Bugs

11. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:188-194`) -- `onLeader` creates the pool, spawns worker goroutines, and defers `close(pool)`. On context cancellation, `close(pool)` fires before the worker goroutines have observed the cancellation. A worker that reaches `w.WorkerPool <- w.JobChannel` after the channel is closed panics. The existing `sync.WaitGroup` in `Start` synchronizes `onLeader` itself but does not track the goroutines spawned inside it. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second `sync.WaitGroup` inside `onLeader` that waits on worker exits before closing. _(Carried from 6/14 #9, P2.)_

12. **`BigQueryDataLakeClient` is a live-registered stub** (`adapters/clients/big_query_data_lake_client.go:24`) -- `DataSourceTypeDataLake` is fully accepted by the domain, persistence, and REST API layers. Any `data-lake` data source created today will be persisted and processed by the worker, but `Query` unconditionally returns `ErrBigQueryDataLakeNotConfigured`. A `data-lake` data source will silently cycle through worker retries and exhaust the retry limit with no meaningful error message. Options: guard the type at the command/API layer until BigQuery is implemented, or surface `ErrBigQueryDataLakeNotConfigured` as a non-retryable error to avoid retry churn. _(Carried from 6/14 #10, P2.)_

#### Architectural

13. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service remains well-covered. _(Carried from 6/14 #11, P3.)_

14. **No domain events** for cross-service communication. _(Carried from 6/14 #12, P3.)_

---

## Priority Summary

| Priority | #   | Issue                                                                                | Service(s)  |
| -------- | --- | ------------------------------------------------------------------------------------ | ----------- |
| **P2**   | 1   | No `SubmissionAttempt` created; error params unused                                  | Forms       |
| **P2**   | 2   | `Tag`/`TagVersion` missing `validate` struct tags -- `ValidateStruct` is no-op       | Forms       |
| **P2**   | 3   | `TagPrimitiveType` unconstrained string type; no constants or predicate              | Forms       |
| **P2**   | 4   | `FieldTagMapping` missing `validate` struct tags -- `ValidateStruct` is no-op        | Forms       |
| **P2**   | 5   | `evaluateCollectionCandidates` nil `collectionIndex` deref -- worker goroutine panic | Forms       |
| **P2**   | 6   | `ErrNoEligibleTagVersion` / `ErrMultipleActiveTagVersion` empty error messages       | Forms       |
| **P2**   | 11  | `defer close(pool)` shutdown race                                                    | pkg/worker  |
| **P2**   | 12  | `BigQueryDataLakeClient` stub accepts API requests, always fails at processing time  | Tenants     |
| **P3**   | 7   | Select/checkbox validator stubs, date partial                                        | Forms       |
| **P3**   | 8   | `tagsService.Delete` missing active-version guard                                    | Forms       |
| **P3**   | 9   | Tenants `Find()` no pagination                                                       | Tenants     |
| **P3**   | 10  | `Lookup` no validation                                                               | Tenants     |
| **P3**   | 13  | Test coverage gaps (forms, pkg/)                                                     | Forms, pkg/ |
| **P3**   | 14  | No domain events                                                                     | All         |

---

## Production Readiness

| Service     | Rating                       | Assessment                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ----------- | ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Forms**   | **8/10 -- Production-Ready** | Submission processing pipeline now fully implemented: visibility and required-field rules are evaluated per-field, field validators run against each visible field, fact candidates are collected and normalized into `CanonicalFact` records via priority-ranked tag resolution. Unstaged fixes `SetRequestDateContext` wrong context key and `rankCandidates` typo. Two new P2 findings: nil deref in `evaluateCollectionCandidates` can panic worker goroutines on collection-tagged submissions; `ErrNoEligibleTagVersion` and `ErrMultipleActiveTagVersion` have empty messages making failure paths completely uninformative. Four carried P2 validation no-ops unchanged. |
| **Tenants** | **9/10 -- Production-Ready** | No functional changes since 6/14. `BigQueryDataLakeClient` stub remains the only visible production gap. `\n` in log format strings fixed in unstaged.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| **pkg/**    | **8/10 -- Production-Ready** | `SetRequestDateContext` wrong context key fixed in unstaged -- a live correctness bug in both services' request date middleware. `defer close(pool)` race in the background worker unchanged. Zero test coverage across `pkg/`.                                                                                                                                                                                                                                                                                                                                                                                                                                                  |

---

## Summary

### Progress Since 6/14

- **Submission processing pipeline implemented** (committed) -- `submissionJobsService` now has a complete `Process` implementation. `sanitize` drives two phases: `extractFactCandidates` walks the form version page/section/field tree, evaluating `RuleTypeVisible` rules via `shouldValidate` to skip hidden elements, evaluating `RuleTypeRequired` rules via `isRequired` to dynamically override field required-ness, and running `validateField` against each visible field before collecting `factCandidate` records for every `FieldTagMapping` on that field. `normalize` then groups candidates by `TagVersionID`, resolves the active/deprecated tag version via `domain.ResolveTagVersion`, and dispatches to either `evaluateCollectionCandidates` (for tags with a collection ancestor, grouping by `CollectionIndex`) or `evaluateScalarCandidates`, both using `rankCandidates` to select the highest-priority candidate per group. The resulting `[]*domain.CanonicalFact` represents the canonical output of the submission. `recordAttempt` routes the outcome to `Accept`, `Reject` (for field validation and version status errors), or `Fail` (for infrastructure errors), then persists the submission in a transaction. The `shouldReject` and `isRetryableError` predicates in the worker correctly classify the error types produced by this pipeline.

- **`SetRequestDateContext` context key bug** (unstaged) -- `context.WithValue` was called with `correlationIDKey` instead of `requestDateKey`, silently discarding the validated `X-Request-Date` header on every request in both services. Fixed: parameter renamed to `requestDate`, key corrected to `requestDateKey`.

- **`rankCandidates` typo** (unstaged) -- `rankCandiates` renamed to `rankCandidates` at the function declaration and both call sites in `submission_jobs_service.go`.

- **Tenants `main.go` log format strings** (unstaged) -- Trailing `\n` removed from three `fmt.Sprintf` format strings passed to `slog`. JSON log `msg` field values are now clean for structured log consumers.

### Current State

**14 remaining issues** (8 P2, 6 P3). 0 P0, 0 P1.

The submission processing pipeline is the most significant development since 6/14. The full sanitize → extract → normalize → rank → record path is now in place and correctly wired through the worker retry loop. The two new P2 findings are direct consequences of this pipeline being live: the nil deref in `evaluateCollectionCandidates` (`submission_jobs_service.go:350`) is a worker goroutine panic risk on collection-tagged submissions where any field value omits `CollectionIndex`; the empty error messages on `ErrNoEligibleTagVersion` and `ErrMultipleActiveTagVersion` make it impossible to diagnose failures when `ResolveTagVersion` fails during `normalize`.

**Hexagonal Architecture** -- No structural changes this cycle. The submission pipeline correctly remains in the core services layer: `extractFactCandidates` and `normalize` operate entirely on domain types, with repository access through the defined secondary ports. The rule evaluator is injected as a `ports.RuleEvaluator` and the field validators as a `ports.FieldValidatorRegistry`, maintaining the inversion of control.

**DDD** -- `TagNodeType` is now properly constrained via `isTagValueKind`, closing the prior `TagType` gap from 6/14. `TagPrimitiveType` is the remaining unconstrained domain string type. `SubmissionAttempt` remains fully modeled but entirely unused at the service layer -- the audit trail gap grows more significant now that the processing pipeline is live and producing real accept/reject/fail outcomes with no per-attempt record.

**Idiomatic Go** -- `rankCandidates` rename is the only idiomatic improvement this cycle. The `evaluateCollectionCandidates` nil deref is the clearest new idiomatic gap: collection candidates must guard `fc.collectionIndex != nil` before dereferencing, consistent with the nil-safe handling already present in `evaluateScalarCandidates`.

### Highest-Impact Improvements

1. **Add nil guard before `*fc.collectionIndex` in `evaluateCollectionCandidates`** (P2) -- check `fc.collectionIndex != nil` before the map key dereference; return an error or skip the candidate if nil to prevent worker goroutine panics on collection-tagged submissions.
2. **Give `ErrNoEligibleTagVersion` and `ErrMultipleActiveTagVersion` meaningful messages** (P2) -- one-line fix per sentinel; immediately makes submission processing failure logs actionable.
3. **Wire `SubmissionAttempt` into `recordAttempt()`** (P2) -- the pipeline now produces real outcomes; append `NewSubmissionAttempt` before persisting and store the `err` that `Reject`/`Fail` currently discard.
4. **Add `validate` struct tags to `Tag`, `TagVersion`, and `FieldTagMapping`** (P2) -- one-line fix per required field; immediately enforces construction-level invariants on the live tag and field-mapping API paths.
5. **Add `TagPrimitiveType` constants and `isValidTagPrimitiveType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTag`; closes the remaining unconstrained domain type gap.
6. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)` or track worker goroutines in a second `sync.WaitGroup` that `onLeader` waits on before closing.
7. **Guard `DataSourceTypeDataLake` at the API or make `ErrBigQueryDataLakeNotConfigured` non-retryable** (P2) -- prevents `data-lake` data sources from silently consuming all worker retries before failing.
8. **Add active-version guard to `tagsService.Delete`** (P3) -- mirrors the existing `hasActiveVersion` check on `formsService.Delete`; FIXME comment already identifies the invariant.
9. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data on the now-live processing pipeline.
10. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3).
