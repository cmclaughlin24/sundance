# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/20 Review

1. ~~Submission processing has no retry on transient failure~~ (Forms) -- New `submissionJob.Process` retry loop in `adapters/workers/submissions_worker.go:40-66` with exponential backoff (`backoff *= 2`) and `retryLimit` from settings. `isRetryableError` classifies `ErrFieldValidation`, `ErrFieldRequired`, `ErrFieldTypeValue`, `ErrVersionStatus`, and `stratreg.ErrStrategyNotFound` as non-retryable; everything else (infrastructure/transient) is retried. *(Committed `7cf1589b`.)*

2. ~~Worker options scattered across constructor calls~~ (All) -- New `WorkerOptions` struct in `adapters/workers/workers.go` with `Interval`, `PoolSize`, `RetryLimit` consolidated. Driven from `settings.json` (`worker.interval`, `worker.pool_size`, `worker.retry_limit`). Both services now bootstrap workers through a uniform `Bootstrap(app, settings)` entry point. *(Committed `d86be47e`.)*

3. ~~Worker pool dispatch was unbuffered~~ (`pkg/worker/background_worker.go:188`) -- `pool := make(chan chan J, bw.size)` is now a buffered channel sized to the pool, eliminating dispatcher blocking when workers are mid-job. *(Committed `d9c6888e`.)*

4. ~~Tenants data sources retry indefinitely on failure~~ (Tenants) -- New `Attempts int` field and `RecordAttempt()` method on `ScheduledDataSourceAttributes`. `FindJobs` filter now includes `RetryLimit` and the mongo query uses `attributes.attempts: $lt RetryLimit` (corrected from `$lte` which was off-by-one). `RefreshData` resets `Attempts` to 0 on successful refresh; `Process` increments on failure. New tests `data_source_test.go`, `data_source_attributes_test.go` cover the retry semantics. *(Committed `32739c48`.)*

5. ~~Domain type `Version` ambiguous with `int` version number~~ (Forms) -- `Version` renamed to `FormVersion` across domain, ports, repositories, DTOs, mongo documents, and handlers. Aligns the type name with its collection (`form_versions`) and removes overloading with `FormVersion.Version int`. *(Committed `8f3a2508`.)*

6. ~~`Process` re-validates already-terminal submissions~~ (Forms) -- New skip guard at `submission_jobs_service.go:78-81` short-circuits when `submission.Status` is `Accepted` or `Rejected`. `Failed` submissions still flow through the retry loop (intentional — `Failed` indicates a transient/infra outcome). *(Committed `7cf1589b`.)*

7. ~~Monolithic `handlers.go` in forms service~~ (Forms) -- Split into `form_handlers.go` (form/version CRUD), `submission_handlers.go` (submission endpoints), and `handlers.go` (constructor, shared error mapping). Handler tests split accordingly. *(Committed `70b16b41`.)*

8. ~~Cache type string `"inmemory"` non-conventional~~ (`pkg/cache/cache.go:17`) -- Renamed to `"in-memory"` for consistency with documented naming. *(Committed `de55cf23`.)*

9. ~~Architecture diagrams missing~~ (Docs) -- New architecture diagram and supporting images added under `backend/docs/`. *(Committed `42d559f1`, `85df5059`.)*

10. ~~Swagger docs stale after `FormVersion` rename~~ (Forms) -- Regenerated `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go`. *(Committed `a940cc43`.)*

11. ~~`ExprRuleEvaluator.statement` overwrites accumulator on every iteration~~ (`adapters/evaluators/expr_rule_evaluator.go:72`) -- Was `stmt = join + statementFn(re)` which silently dropped all but the last expression in a compound rule. Now `stmt = stmt + join + statementFn(re)`. Compound visibility/required rules with AND/OR across ≥2 expressions evaluate correctly. *(Committed `9a2f985e`.)*

12. ~~Commented-out canonical tag routes in `routes.go`~~ (`adapters/rest/routes.go:69-81`) -- Removed. *(Committed `9a2f985e`.)*

13. ~~Worker duration log used `fmt.Sprintf("%d", ...)`~~ (`pkg/worker/worker.go:88`) -- Replaced with native int attribute `"duration_ms", time.Since(start).Milliseconds()`. `fmt` import removed. *(Committed `9a2f985e`.)*

14. ~~`ScheduledDataSourceAttributes.Attempts` misaligned~~ (`data_source_attributes.go:33`) -- gofmt indentation corrected. *(Committed `9a2f985e`.)*

15. ~~`getReferenceIDPathValue` colocated in `handlers.go`~~ -- Moved to `submission_handlers.go` next to `replaySubmission`, matching the colocation of `getFormIDPathValue`/`getVersionIDPathValue` in `form_handlers.go`. *(Committed `9a2f985e`.)*

16. ~~Canonical tag scaffold: `HydrateCaonicalTagVersion` typo~~ (`canonical_tag_version.go:36`) -- Renamed to `HydrateCanonicalTagVersion`. *(Committed `9a2f985e`.)*

17. ~~Canonical tag scaffold: `CanonicalTagID` dropped in constructors~~ (`canonical_tag_version.go`) -- Both `NewCanonicalTagVersion` and `HydrateCanonicalTagVersion` now assign `CanonicalTagID` from their parameter. *(Committed `9a2f985e`.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

`RuleExpression.FieldKey` has no referential integrity check against version fields -- Rules may be created when a field does not yet have an ID and therefore cannot be associated at creation time. Expression field keys are resolved at evaluation time against the `RuleEvaluationContext`; invalid keys evaluate to nil/zero in the `expr` environment, which is acceptable behavior for conditional visibility rules. *(Carried from 5/20.)*

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **No `SubmissionAttempt` ever created** (`submission_jobs_service.go:216-245`) -- `recordAttempt` transitions status (`Accept`/`Reject`/`Fail`) and persists the submission but never appends a `SubmissionAttempt` record. The domain type (`submission_attempt.go`), the `Submission.Attempts` field, the mongo document mapper, and `NewSubmissionAttempt` constructor all exist as dead code. `Reject(err error)` and `Fail(err error)` accept an `err` parameter but never store it — once attempts are created, the error should land in `SubmissionAttempt.ErrorDetails`. With the retry loop now in place, operators have no way to see how many attempts a `Failed` submission consumed or what specific error caused each one. *(Carried from 5/20 #1, P2.)*

2. **Canonical tag scaffold: unconstrained status/type and missing `validate` tags** (`canonical_tag.go`, `canonical_tag_version.go`) -- Two issues remain in the in-progress canonical tag domain after the recent fixes:
   - `CanonicalTagStatus` and `CanonicalTagType` are unconstrained string aliases. No constants, no `isValid*` predicate, no `validate.NewTypeValidator`. Diverges from the established pattern (`FormVersionStatus`, `RuleType`, `DataSourceType`, `SubmissionStatus`). Any string is accepted.
   - `CanonicalTag` and `CanonicalTagVersion` carry no `validate:"required"` struct tags despite both `NewCanonicalTag` and `NewCanonicalTagVersion` calling `validate.ValidateStruct(...)`. The validation calls are no-ops. Compare to `Submission`, `DataSource`, `Form`, `FormVersion`, etc., which tag their required fields.

   The slice is still dead code (no ports/services/repos/DTOs/routes), so this isn't blocking production today, but the moment the slice is wired the domain will accept arbitrary status/type strings and skip validation on empty fields. *(New, P2.)*

#### Missing Functionality

3. **Field validator strategies: select and checkbox remain stubs** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`) -- Both return `nil` without performing any validation. Date has `checkValueRequired` but no date-specific validation (TODO at `date_field_validator.go:37`). *(Carried from 5/20 #2, P3.)*

4. **`FindJobs` `Take` limit not applied at repository level** (`inmemory/submissions_repository.go:87-106`, `mongodb/submissions_repository.go:76-91`) -- `FindSubmissionsFilter.Take` is passed by the service layer, but neither the in-memory nor MongoDB `FindJobs` implementations apply the limit. The mongo implementation also has a leftover TODO about projecting only the `_id` field. *(Carried from 5/20 #3, P3.)*

#### Architectural

5. **Duplicate error classification across service and worker** (`services/submission_jobs_service.go:247-252`, `adapters/workers/submissions_worker.go:116-122`) -- `shouldReject` (service: terminal-vs-retryable for status transition) and `isRetryableError` (worker: continue-vs-stop the retry loop) reference overlapping error sets via mirrored predicates with inverted polarity. The two are aligned by accident today and drift would be silent. Specifically: `isRetryableError` classifies `stratreg.ErrStrategyNotFound` as non-retryable but `shouldReject` does not include it, so an unknown field type today produces a `Failed` submission (suggesting retryability) for what is semantically a permanent error. Recommend a single `strategies.Classify(err) Disposition` (or similar) consumed by both call sites. *(New, P3.)*

6. **Inconsistent command/query `Validate()` methods** (`ports/commands.go:155-177`, `ports/query.go:75-97`) -- `CreateSubmissionCommand`, `FindSubmissionsQuery`, and `FindSubmissionByIDQuery` are missing the `Validate()` method that the rest of the commands/queries expose. `submissionsService` validates these via direct `validate.ValidateStruct(command)` calls, but the API surface is inconsistent. *(New, P3.)*

7. **`formsService.hasActiveVersion` returns `true, err` on match path** (`forms_service.go:384`) -- Inside the loop, the success-case return is `return true, err` where `err` is nil at that point. Reads as a bug and should be `return true, nil`. *(New, P3.)*

---

### Tenants Service

#### Missing Functionality

8. **`Find()` has no pagination or filtering** (`tenants_service.go:30-40`) -- *(Carried from 5/20 #4, P3.)*

9. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/20 #5, P3.)*

---

### Cross-Service / pkg/

#### Bugs

10. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:188-194`) -- `onLeader` creates the pool, spawns worker goroutines inside it (`w.Start(ctx)` at line 193), and defers `close(pool)`. On context cancellation or failover, `onLeader` returns and the deferred `close(pool)` runs. The worker goroutines still send to `pool` via `w.WorkerPool <- w.JobChannel:` and observe `wctx.Done()` independently; if a worker reaches the send case after the channel is closed, it panics on a closed channel. The existing `sync.WaitGroup` in `Start` synchronizes `onLeader` (registered via `wg.Go` at line 144) but does not track the worker goroutines, so `close(pool)` runs without waiting for workers to observe the cancellation. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second WaitGroup inside `onLeader` to wait on worker exits. *(New, P2.)*

#### Architectural

11. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. Tenants service is now well-covered (22 test files). *(Carried from 5/20 #6, P3.)*

12. **No domain events** for cross-service communication. *(Carried from 5/20 #7, P3.)*

13. **No real authentication** -- All services use `PlaceholderAuthenticator`. *(Carried from 5/20 #8, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | No `SubmissionAttempt` created; error params unused | Forms |
| **P2** | 2 | Canonical tag scaffold: unconstrained status/type, no validate tags | Forms |
| **P2** | 10 | `defer close(pool)` shutdown race | pkg/worker |
| **P3** | 3 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 4 | `FindJobs` `Take` not honored | Forms |
| **P3** | 5 | Duplicate error classification (service vs worker) | Forms |
| **P3** | 6 | Inconsistent `Validate()` on submission commands/queries | Forms |
| **P3** | 7 | `hasActiveVersion` returns `true, err` instead of `true, nil` | Forms |
| **P3** | 8 | Tenants `Find()` no pagination | Tenants |
| **P3** | 9 | `Lookup` no validation | Tenants |
| **P3** | 11 | Test coverage gaps (forms, pkg/) | Forms, pkg/ |
| **P3** | 12 | No domain events | All |
| **P3** | 13 | Placeholder authentication only | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8.5/10 -- Production-Ready** | Submission processing pipeline now includes a bounded retry loop with exponential backoff and disciplined error classification (`isRetryableError` separates infra from validation/permanent errors). `Process` short-circuits already-terminal submissions, preventing rework. `FormVersion` rename clarifies the domain language. Handler refactor improves maintainability. `ExprRuleEvaluator` accumulator bug — which previously caused compound rules to evaluate only their last expression — is fixed, restoring compound-rule correctness. Remaining P2 work: `SubmissionAttempt` is still never created despite the model and persistence path existing, so the system has no audit trail of per-attempt outcomes — a gap that becomes more visible now that retries exist. Canonical tag scaffold introduces correctness debt (unconstrained status/type, missing validate tags) that is harmless today (dead code) but would corrupt persisted state the moment the slice is wired. No service/domain/repository tests yet. |
| **Tenants** | **9/10 -- Production-Ready** | Data-source retry limit closes a real correctness gap: `Attempts` tracked on `ScheduledDataSourceAttributes`, mongo query uses `$lt` (the previous `$lte` was off-by-one), and `RefreshData` resets the counter on success. Comprehensive test coverage carried over from 5/20 remains intact, including new `data_source_test.go` cases for `Attempts`/`RecordAttempt`. Only P3 gaps remain (pagination, `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | Worker pool channel is now buffered, eliminating dispatcher blocking. `WorkerOptions` consolidation gives both services a uniform settings surface. Worker duration log is now a native int attribute. Only material concern is the `defer close(pool)` shutdown race — the existing `sync.WaitGroup` synchronizes `onLeader` but not the worker goroutines spawned inside it. Zero test coverage in the package. |

---

## Summary

### Progress Since 5/20

- **Submission retry pipeline** (committed) -- `submissionJob.Process` now wraps `service.Process` in a bounded retry loop with exponential backoff. `isRetryableError` distinguishes permanent failures (`ErrFieldValidation`, `ErrFieldRequired`, `ErrFieldTypeValue`, `ErrVersionStatus`, `ErrStrategyNotFound`) from transient ones; transient errors retry up to `retryLimit` (5 by default) with `backoff *= 2`. The skip guard in `Process` ensures already-terminal submissions are not re-validated.

- **Worker options consolidated into settings** (committed) -- Both services expose `worker.interval`, `worker.pool_size`, `worker.retry_limit` in `settings.json` and wire through a uniform `Bootstrap(app, settings)` entry point. The `WorkerOptions` struct centralizes defaults (1m interval, pool size 5, retry limit 5).

- **Buffered worker pool channel** (committed) -- `pool := make(chan chan J, bw.size)` eliminates dispatcher blocking when workers are mid-job. The dispatch path can now enqueue without contention.

- **Tenants data-source retry limit** (committed) -- `Attempts` counter added to `ScheduledDataSourceAttributes` with `RecordAttempt()` increment on failure and reset-to-zero in `RefreshData`. `FindDataSourceJobsFilter.RetryLimit` flows from settings through query → filter → repository. MongoDB query uses `$lt` (previously `$lte`, an off-by-one bug). New domain tests cover the semantics.

- **`Version` → `FormVersion` rename** (committed) -- Domain type, repository, ports, DTOs, mongo documents, and handlers all updated. The collection name (`form_versions`) and the type name now agree. Swagger docs regenerated.

- **Handler refactor** (committed) -- `handlers.go` split into `form_handlers.go` and `submission_handlers.go`. Tests split correspondingly.

- **Bug fixes discovered during review** (committed) -- `ExprRuleEvaluator.statement` accumulator (compound rules now correctly accumulate the expression string across iterations); commented-out canonical tag routes removed; worker duration log now uses native int with `duration_ms` key; gofmt fix on `Attempts` indentation; `getReferenceIDPathValue` colocated with the submission handler that uses it; `HydrateCaonicalTagVersion` typo corrected; `CanonicalTagID` now assigned in both canonical tag version constructors.

### Current State

**13 remaining issues** (3 P2, 10 P3). 0 P0, 0 P1. Multiple issues resolved this cycle across the committed work and the same-cycle bug-fix commit (`9a2f985e`).

**Forms Service at 8.5/10** (down 0.5 from 9/10). Net-positive cycle: the retry pipeline is a real production-readiness improvement and the `FormVersion` rename clarifies the model. The evaluator accumulator fix restores compound-rule correctness. However, the canonical tag scaffold introduces correctness debt — harmless today but loaded for the slice's eventual wiring — and the `SubmissionAttempt` audit trail remains absent precisely when the new retry loop makes it most useful.

**Tenants Service at 9/10** (unchanged). The retry-limit work is a clean addition with domain tests; the `$lte` → `$lt` fix corrects an off-by-one that would have under-counted final attempts. No regressions.

**pkg/ at 8/10** (unchanged). Buffered pool and structured duration logging are positives; the `defer close(pool)` race is a real but bounded concern that needs addressing before any high-frequency leader churn scenario.

**Hexagonal Architecture** -- The retry loop correctly lives in the worker adapter, not the service. `isRetryableError` is an adapter-level concern (orchestration policy) while `shouldReject` is a service-level concern (domain status decision); both reference shared sentinels from `strategies` and `services`, maintaining proper dependency direction. The `WorkerOptions` consolidation pushes infrastructure configuration to the composition root via `settings.json`, keeping the workers package agnostic of how its parameters arrive. The tenants `Attempts`/`RecordAttempt`/`RefreshData` semantics are correctly modeled as domain methods on `ScheduledDataSourceAttributes` rather than service-layer bookkeeping.

**DDD** -- `Attempts` and `RecordAttempt()` belong on the domain value object that owns the lifecycle (`ScheduledDataSourceAttributes`), and the symmetric reset in `RefreshData` keeps the invariant local. The `FormVersion` rename brings ubiquitous language into alignment with persistence. The remaining DDD gap is the same one carried from 5/20: `SubmissionAttempt` is a value object modeled on the submission aggregate but never appended; the model exists, the persistence path exists, only the service-layer call to construct one is missing. Canonical tags are scaffolded but the domain invariants (status/type enumeration, required fields) are not yet enforced — fixing those before the slice is wired keeps the aggregate honest.

**Idiomatic Go** -- Retry loop uses `errors.Is` for sentinel checking with multiple `||`-chained predicates, which scales acceptably for the current set but would benefit from a `Classify(err)` style helper as more sentinels accumulate. `defer close(pool)` is a common Go pattern that fits when the channel sender is the same goroutine; here it's not — the sender goroutines are spawned inside `onLeader` but outlive it via context cancellation. The switch from `fmt.Sprintf("%d", ms)` to native `int64` slog attribute is correct structured-logging form.

### Highest-Impact Improvements

1. **Wire `SubmissionAttempt` into `recordAttempt()`** (P2) -- append a `NewSubmissionAttempt(attempt, result, errorDetails)` to `submission.Attempts` before persisting; store the `err` that `Reject`/`Fail` currently discard. With the retry loop in place, this gives operators a complete history of how each submission was processed.
2. **Finish the canonical tag domain before wiring the slice** (P2) -- add status/type constants + `isValid*` predicate + `NewTypeValidator`, add `validate:"required"` struct tags on `CanonicalTag` and `CanonicalTagVersion`. Cheap now, expensive after the first persisted record.
3. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)`, or track worker goroutines in a WaitGroup that `onLeader` waits on before closing.
4. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid data.
5. **Apply `FindJobs` `Take` limit in both repository implementations** (P3) -- prevents unbounded result sets and removes the leftover mongo TODO.
6. **Consolidate error classification** (P3) -- single `Classify(err) Disposition` consumed by both `recordAttempt` and `isRetryableError` to prevent drift.
7. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3) -- narrows the largest remaining test gap.
