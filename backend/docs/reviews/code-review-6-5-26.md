# Full Codebase Review: Forms and Tenants Services
**Date:** June 5, 2026

---

## Issues Resolved Since 5/29 Review

1. ~~`TagDocument.UpdatedAt` has `wson` typo instead of `bson`~~ (`adapters/persistence/mongodb/documents/tag.go:14`) -- Fixed in unstaged changes: `wson:"updated_at"` → `bson:"updated_at"`.

2. ~~`validateField` returns raw `fmt.Errorf` for missing required field~~ (`submission_jobs_service.go`) -- Fixed in unstaged changes: now wraps `strategies.ErrFieldRequired` with `%w`, making `shouldReject` correctly route missing-required-field submissions to `Reject` instead of `Fail`.

3. ~~`tagsService.Find` logs hardcoded empty string for `tenant_id`~~ (`tags_service.go:30`) -- Fixed in unstaged changes: now logs `query.TenantID`.

4. ~~`HydrateTagVersion` parameter name typo `deprectatedAt`~~ (`tag_version.go:65`) -- Fixed in unstaged changes: renamed to `deprecatedAt`.

5. ~~Command/query API boundary inconsistency~~ (`ports/primary.go`, `ports/commands.go`, `ports/query.go`) -- `FormsAPI`, `SubmissionsAPI`, and `SubmissionJobsAPI` interfaces migrated from pointer to value semantics for all flat command/query types. `CreateFormVersionCommand`, `UpdateFormVersionCommand`, and `CreateSubmissionCommand` intentionally retain pointer semantics (carry `[]*domain.Page` / `[]*domain.SubmissionFieldValue`). Constructors updated to return values; `Validate()` methods updated to use value receivers; service implementations, handlers, mocks, and tests all updated. Verified: clean build, all tests pass.

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

`RuleExpression.FieldKey` has no referential integrity check against version fields -- Rules may be created when a field does not yet have an ID and cannot be associated at creation time. Expression field keys are resolved at evaluation time against the `RuleEvaluationContext`; invalid keys evaluate to nil/zero in the `expr` environment, which is acceptable behavior for conditional visibility rules. *(Carried from 5/20.)*

**Duplicate error classification: `stratreg.ErrStrategyNotFound` produces `Failed` instead of `Rejected`** (`submission_jobs_service.go`, `adapters/workers/submissions_worker.go`) -- `isRetryableError` correctly returns `false` for `ErrStrategyNotFound` (stop retrying — the missing strategy will still be missing on the next attempt). `shouldReject` correctly excludes it — `Rejected` communicates "I processed your submission and it was invalid", which would be a lie; the submission itself is valid, the system was misconfigured. `Failed` is the correct semantic: an internal processing error that is not the submitter's fault. The two predicates are not misaligned; they are doing different jobs. *(Flagged 5/29 #4, closed as Will Not Fix.)*

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **No `SubmissionAttempt` ever created** (`submission_jobs_service.go:216-245`) -- `recordAttempt` transitions status (`Accept`/`Reject`/`Fail`) and persists the submission but never appends a `SubmissionAttempt` record. `NewSubmissionAttempt`, the `Submission.Attempts` field, the mongo document mapper, and the `SubmissionAttempt` domain type all exist as dead code. `Reject(err error)` and `Fail(err error)` accept an `err` parameter but never store it. With the retry loop in place, operators have no audit trail of per-attempt outcomes or specific errors. *(Carried from 5/29 #1, P2.)*

2. **`Tag` and `TagVersion` have no `validate` struct tags despite calling `validate.ValidateStruct`** (`tag.go:15-22`, `tag_version.go:29-39`) -- `Tag` fields (`TenantID`, `Key`, `DisplayName`) and `TagVersion` fields (`TagID`, `Version`, `Type`, `Status`) carry no `validate:"required"` tags. Both `NewTag` and `NewTagVersion` call `validate.ValidateStruct` but the call is a no-op on untagged structs. Empty `tenantID`, empty `key`, or empty `tagType` pass construction without error. Compare to `Form`, `FormVersion`, `Submission`, `DataSource` — all tag their required fields. *(New, P2.)*

3. **`TagType` has no constants and no `isValidTagType` predicate** (`tag_version.go:18`) -- `type TagType string` is declared with no constants. `NewTagVersion` accepts any `TagType` string, including empty string, without validation. Every other constrained domain type in the forms service — `FormVersionStatus`, `FieldType`, `SubmissionStatus`, `RuleType`, `ExprOperator`, `DataSourceType` — uses `validate.NewTypeValidator` with an explicit enum. `TagType` is the sole exception. *(New, P2.)*

#### Missing Functionality

4. **Field validator strategies: select and checkbox remain stubs** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`) -- Both return `nil` without performing any validation. Date has `checkValueRequired` but no date-specific validation (TODO at `date_field_validator.go:37`). *(Carried from 5/29 #3, P3.)*

5. **`tagsService.Delete` missing active-version guard** (`tags_service.go:138`) -- A `// FIXME` comment acknowledges the missing invariant: a tag with a historically active version should not be deletable in order to preserve audit history. The guard exists on `formsService.Delete` via `hasActiveVersion` but is absent here. Tags with published version history can currently be hard-deleted. *(New, P3.)*

---

### Tenants Service

#### Missing Functionality

6. **`Find()` has no pagination or filtering** (`tenants_service.go:30-40`) -- *(Carried from 5/29 #5, P3.)*

7. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/29 #6, P3.)*

---

### Cross-Service / pkg/

#### Bugs

8. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:188-194`) -- `onLeader` creates the pool, spawns worker goroutines, and defers `close(pool)`. On context cancellation, `close(pool)` fires before workers have observed the cancellation. A worker reaching `w.WorkerPool <- w.JobChannel` after close panics. The `sync.WaitGroup` in `Start` synchronizes `onLeader` itself but does not track the worker goroutines spawned inside it. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second WaitGroup inside `onLeader` to wait on worker exits before closing. *(Carried from 5/29 #7, P2.)*

#### Architectural

9. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service remains well-covered. *(Carried from 5/29 #8, P3.)*

10. **No domain events** for cross-service communication. *(Carried from 5/29 #9, P3.)*

11. **No real authentication** -- All services use `PlaceholderAuthenticator`. *(Carried from 5/29 #10, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | No `SubmissionAttempt` created; error params unused | Forms |
| **P2** | 2 | `Tag`/`TagVersion` missing `validate` struct tags — `ValidateStruct` is no-op | Forms |
| **P2** | 3 | `TagType` has no constants and no `isValidTagType` predicate | Forms |
| **P2** | 8 | `defer close(pool)` shutdown race | pkg/worker |
| **P3** | 4 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 5 | `tagsService.Delete` missing active-version guard | Forms |
| **P3** | 6 | Tenants `Find()` no pagination | Tenants |
| **P3** | 7 | `Lookup` no validation | Tenants |
| **P3** | 9 | Test coverage gaps (forms, pkg/) | Forms, pkg/ |
| **P3** | 10 | No domain events | All |
| **P3** | 11 | Placeholder authentication only | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8.5/10 — Production-Ready** | Tag slice is now live-wired with full service/handler coverage. `wson` BSON typo fix closes a silent data integrity gap on `UpdatedAt`. `validateField` fix correctly routes missing-required-field submissions to `Reject` rather than `Fail`, preventing incorrect retries. API boundary is now consistent — flat commands/queries pass by value; commands carrying `[]*domain.Page` or `[]*domain.SubmissionFieldValue` retain pointer semantics. Remaining P2 gaps: `Tag`/`TagVersion` domain validation is a no-op (missing struct tags and `TagType` constants), meaning the live tag API accepts invalid or empty inputs today; `SubmissionAttempt` audit trail remains absent despite the retry loop making it more necessary than ever. |
| **Tenants** | **9/10 — Production-Ready** | No changes this cycle. P3 gaps only (pagination, `Lookup` validation). |
| **pkg/** | **8/10 — Production-Ready** | No changes this cycle. `defer close(pool)` race remains the only material concern. |

---

## Summary

### Progress Since 5/29

- **4 bugs fixed in unstaged changes** — `wson` BSON typo on `Tag.UpdatedAt`, `validateField` incorrect error wrapping for missing required fields, `tags_service.Find` hardcoded empty `tenant_id` log attribute, `HydrateTagVersion` `deprectatedAt` parameter name typo.

- **API boundary value/pointer consistency** — `FormsAPI`, `SubmissionsAPI`, `SubmissionJobsAPI` interfaces and all implementing service methods, handlers, mocks, and tests migrated from pointer to value semantics for all flat types. The three commands carrying slices of domain aggregates intentionally retain pointer semantics per Go Code Review Comments guidance on structs whose elements are pointers to potentially mutating data. Verified: clean build, all tests pass.

- **Tag slice promoted from scaffold to live feature** — since 5/29, `Tag`/`TagVersion` are fully wired through service, repository, and REST handlers with all CRUD and version lifecycle (`Publish`/`Deprecate`/`Retire`) endpoints active.

- **`ErrStrategyNotFound` error classification closed as Will Not Fix** — `Failed` is the correct semantic for a developer misconfiguration; `Rejected` would misrepresent a valid submission as invalid. The two predicates (`shouldReject`, `isRetryableError`) are doing different jobs and their current behavior is intentional.

### Current State

**11 remaining issues** (4 P2, 7 P3). 0 P0, 0 P1.

The tag slice promotion from dead scaffold to live feature is the most significant development since 5/29. It closes the prior "canonical tag dead code" issue but surfaces two new P2 issues that were previously inert: `Tag`/`TagVersion` missing `validate` struct tags (construction validation is a no-op) and `TagType` having no constants or validator. These are now live production code paths accepting arbitrary or empty inputs.

**Hexagonal Architecture** -- The API boundary refactor correctly distinguishes between command/query types that are semantically immutable inputs (value types) and those that carry mutable domain aggregates (pointer types). The tag service adapter correctly isolates `isValidAccess` as a secondary port call before mutation, maintaining the pattern established by the forms service.

**DDD** -- `Tag.Update` correctly uses the copy-and-validate pattern (`cpy := *t; ... *t = cpy`) consistent with other domain mutations. `TagVersion` lifecycle methods (`Publish`, `Deprecate`, `Retire`) enforce valid state transitions with sentinel errors. The gap is that `NewTag` and `NewTagVersion` constructors do not enforce required-field invariants — the `validate.ValidateStruct` call is structurally present but functionally inert without struct tags. `SubmissionAttempt` remains the other outstanding DDD gap: the value object, its constructor, and its persistence path all exist; only the service-layer call to construct and append one is missing.

**Idiomatic Go** -- Tag commands use value types and value receivers consistently, matching the pattern established by this review cycle. `FindSubmissionByIDQuery[T]` remains a type alias (`=`) for `FindByIDQuery[T]` rather than a distinct type — retained intentionally.

### Highest-Impact Improvements

1. **Add `validate` struct tags to `Tag` and `TagVersion`** (P2) -- one-line fix per required field; immediately enforces construction-level invariants on the live tag API.
2. **Add `TagType` constants and `isValidTagType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTagVersion`; closes the unconstrained type string gap.
3. **Wire `SubmissionAttempt` into `recordAttempt()`** (P2) -- append `NewSubmissionAttempt` before persisting; store `err` from `Reject`/`Fail` in `ErrorDetails`.
4. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)` or track worker goroutines in a second WaitGroup that `onLeader` waits on before closing.
5. **Add active-version guard to `tagsService.Delete`** (P3) -- mirrors the existing `hasActiveVersion` check on `formsService.Delete`; the FIXME comment already identifies the invariant.
6. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data.
7. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3).
