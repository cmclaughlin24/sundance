# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/10 Review

1. ~~`errors.Is` not used for sentinel error comparison~~ (Cross-Service, P2) -- All services now use `errors.Is` for sentinel error comparisons. `tenants_service.go:152`, `data_sources_service.go:216`, `submissions_service.go:149,157`, `forms_service.go:392`, and `forms_service.go:400` all use `errors.Is` correctly. *(Resolves 5/10 carried cross-service concern.)*

2. ~~`InMemoryDataSourceRepository.FindJobs` does not filter~~ (Tenants, P2) -- Filtering by type, limit, and expiration is fully implemented. *(Resolves 5/10 #1.)*

3. ~~`DataSourcesJobService.Find` ignores `query.Limit`~~ (Tenants, P2) -- The `Limit` field is now passed through to the `FindDataSourceJobsFilter`. *(Found and fixed same cycle.)*

4. ~~`TenantsService.Update` logs before persistence~~ (Tenants, P2) -- Log ordering corrected. *(Found and fixed same cycle.)*

5. ~~`scheduledDataSourceAttributesResponse` DTO omits fields~~ (Tenants, P3) -- `IntervalHours` and `ExpirationDate` now included in the response DTO. *(Found and fixed same cycle.)*

6. ~~`mockDataSourcesRepository` missing `FindJobs`~~ (Tenants, P3) -- Mock now satisfies `ports.DataSourcesRepository`. *(Found and fixed same cycle.)*

7. ~~`createSubmission` handler uses unbuffered channel~~ (Submissions, P2) -- All handler channels now consistently buffered. *(Found and fixed same cycle.)*

8. ~~`Elector` interface has no distributed implementation~~ (pkg/worker, P3) -- `CacheElector` wrapping `RedisCacheManager` provides distributed leader election via `AcquireLock`/`RenewLock`/`ReleaseLock` backed by Redis `SetNX` and Lua scripts. *(Resolves 5/10 #4.)*

9. ~~`logFindFormByIDError` uses `==` instead of `errors.Is`~~ (Forms, P2) -- Now uses `errors.Is` correctly. *(Found and fixed same cycle.)*

10. ~~Validation rules inconsistent across layers~~ (Forms, P2) -- DTO and command validation rules are now aligned. *(Found and fixed same cycle.)*

11. ~~Worker pool dispatch blocks indefinitely if all workers are busy~~ (pkg/worker, P3) -- Per-job timeouts via `SetTimeout` on the `BackgroundWorkerBuilder` ensure that hung jobs are cancelled, workers return to the pool, and the dispatch loop unblocks. *(Resolves 5/10 #6.)*

12. ~~`InMemorySubmissionsRepository.Find` ignores `Statuses` filter~~ (Submissions, P2) -- The in-memory implementation now filters by `filter.Statuses` using `slices.Contains`, matching the MongoDB adapter behavior. *(Unstaged.)*

13. ~~`submissionAttemptDocument` has unused `IdempotencyID` field~~ (Submissions, P3) -- The `IdempotencyID` field has been removed from the `submissionAttemptDocument` BSON struct. *(Unstaged.)*

14. ~~`InMemoryVersionsRepository` lacks duplicate `(form_id, version)` enforcement~~ (Forms, P3) -- `Upsert` now checks for existing versions with the same number before inserting and returns `domain.ErrDuplicateVersion`, matching the MongoDB unique index behavior. *(Unstaged.)*

15. ~~Service structs exported but returned as interfaces~~ (All, P3) -- `FormsService` renamed to `formsService`, `TenantsService` to `tenantsService`, `DataSourcesService` to `dataSourcesService`, `DataSourcesJobService` to `dataSourcesJobService`, `SubmissionsService` to `submissionsService`, `SubmissionJobsService` to `submissionJobsService`. All receiver methods updated. Test files updated to use unexported types. *(Unstaged.)*

16. ~~Typo: `"placholder"` in authenticator init~~ (Forms, P3) -- Corrected to `"placeholder"`. *(Unstaged.)*

17. ~~`mockDatabase.Close` signature mismatch~~ (Tenants, P3) -- `Close()` updated to `Close(context.Context)`, matching the `database.Database` interface. *(Unstaged.)*

18. ~~`ClaimsContextKey` is exported~~ (pkg/auth, P3) -- `ClaimsContextKey` renamed to `claimsContextKey` (unexported), matching the unexported context key pattern used in `httputil`. *(Unstaged.)*

19. ~~`MongoDBDatabase.BeginTx` leaks session on `StartTransaction` failure~~ (pkg/database, P2) -- `session.EndSession(ctx)` is now called before returning on the `StartTransaction` error path, preventing session leaks. *(Unstaged.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`Create` blocked -- `FindByIdempotencyID` returns `ErrNotFound` on miss, treated as failure** (`submissions_service.go:96-105`) -- The `Create` method calls `FindByIdempotencyID` and treats any non-nil error as a failure. The in-memory repository returns `common.ErrNotFound` when no match exists, meaning every genuinely new submission creation fails. The service assumes `nil, nil` means "not found" but repositories return `nil, ErrNotFound`. *(P0 -- blocks all submission creation.)*

2. **`Replay` does not check tenant authorization** (`submissions_service.go:130-145`) -- Unlike `FindByID` and `FindByReferenceID`, the `Replay` method does NOT verify `submission.TenantID != command.TenantID`. Any tenant can replay another tenant's submission. *(P1 -- authorization bypass.)*

3. **`submissionJob.Process` does not pass the held `*domain.Submission` to the service** (`workers/submissions_worker.go:26-28`) -- The `submissionJob` holds a `*domain.Submission` in field `s`, but `Process` calls `j.service.Process(ctx)` without passing it. The `SubmissionJobsService.Process` method has no way to know which submission to process. Worker data flow is fundamentally broken. *(P1 -- worker non-functional.)*

4. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go:29-32`) -- `TenantID` and `ID` fields have no `validate:"required"` tags, so `validate.ValidateStruct(command)` always passes even with empty values. *(P2.)*

#### Missing Functionality

5. **`sendErrorResponse` wrapper adds no value** (`handlers.go:194-199`) -- *(Carried from 5/10 #7, P3.)*

6. **`Replay` service method is a stub** (`submissions_service.go:130-145`) -- *(Carried from 5/10 #8, P3.)*

7. **`SubmissionJobsService.Process` is a stub** (`submission_jobs_service.go:35-37`) -- Returns `nil` without doing anything. Separate from the `Replay` stub. *(P3.)*

#### Code Quality

8. **`Payload` typed as `any`** (`submission.go:36`) -- *(Carried from 5/10 #9, P3.)*

---

### Tenants Service

#### Missing Functionality

9. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

10. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Forms Service

#### Bugs

11. **`Form.Update` mutates fields before validation** (`domain/form.go:53-66`) -- The `Update` method sets fields, then calls `validate.ValidateStruct(f)`. If validation fails, the entity is left in a dirty state with the invalid data already applied. Should validate first or operate on a copy. *(P2.)*

---

### pkg/

#### Bugs

12. **`CacheManager.Get` cannot distinguish cache miss from cached zero-value** (`redis_cache_manager.go:45-69`, `inmemory_cache_manager.go:28-43`) -- Returns `nil` error for both cache miss and cached zero-value. Should return a `bool` hit/miss indicator or a sentinel `ErrCacheMiss`. *(P2.)*

13. **`ErrMissingIdempotencyHeader` and `ErrMissingTenantID` map to 500** (`httputil/http.go:104-109`) -- Neither error matches any `errors.Is` check in `SendErrorResponse`, falling through to the default 500 case. Should be 400 Bad Request. *(P2.)*

#### Architectural

14. **`BackgroundWorker.Start` releases leadership on shutdown but not on `workFn` errors** (`background_worker.go:98`) -- *(Carried from 5/10 #5, P3.)*

15. **`CacheManager` interface conflates caching with distributed locking** (`cache/cache.go:29-36`) -- `Get/Set/Del` (caching) bundled with `AcquireLock/RenewLock/ReleaseLock` (distributed locking). These are separate concerns; violates ISP. *(P2.)*

#### Missing Functionality

16. **No `Close()` on `CacheManager`** (`cache/cache.go`) -- Redis connection never cleanly shut down. Compare with `Database` which has `Close`. *(P3.)*

17. **`RedisCacheManager.Set` hardcodes TTL=0** (`redis_cache_manager.go:79`) -- All cache entries stored with no expiry. Cache grows unboundedly. *(P3.)*

---

### Cross-Service

#### Architectural

18. **Test coverage gaps** -- Submissions has no handler or service tests. No domain-layer or repository-layer tests exist across services. Zero test files in entire `pkg/` directory. *(Carried from 5/10 #10, P3.)*

19. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

20. **No real authentication** -- Placeholder only. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 1 | `Create` blocked by `FindByIdempotencyID` error handling | Submissions |
| **P1** | 2 | `Replay` tenant authorization bypass | Submissions |
| **P1** | 3 | Worker data flow broken -- submission not passed to `Process` | Submissions |
| **P2** | 4 | `ReplaySubmissionCommand` no validation tags | Submissions |
| **P2** | 11 | `Form.Update` mutates before validation | Forms |
| **P2** | 12 | Cache miss indistinguishable from zero-value | pkg/cache |
| **P2** | 13 | Middleware errors map to 500 instead of 400 | pkg/common |
| **P2** | 15 | `CacheManager` conflates caching with locking (ISP) | pkg/cache |
| **P3** | 5 | `sendErrorResponse` wrapper is no-op | Submissions |
| **P3** | 6 | `Replay` is a stub | Submissions |
| **P3** | 7 | `SubmissionJobsService.Process` is a stub | Submissions |
| **P3** | 8 | `Payload` typed as `any` | Submissions |
| **P3** | 9 | `Find()` no pagination | Tenants |
| **P3** | 10 | `Lookup` no validation | Tenants |
| **P3** | 14 | Leadership held despite failures | pkg/worker |
| **P3** | 16 | No `Close()` on `CacheManager` | pkg/cache |
| **P3** | 17 | Redis cache has no TTL | pkg/cache |
| **P3** | 18 | Test coverage gaps | All |
| **P3** | 19 | No domain events | All |
| **P3** | 20 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8/10 -- Beta** | Most mature service. Complete CRUD with versioning lifecycle (create, publish, retire). Handler tests provide good HTTP layer coverage. In-memory adapter now enforces version uniqueness matching MongoDB. Service structs properly unexported. Blocked from production by: no service/domain layer tests, `Form.Update` dirty-state bug, no real authentication. |
| **Tenants** | **7/10 -- Beta** | Fully functional including background job processing pipeline and leader election infrastructure. Clean hexagonal structure with strategies pattern. Service structs properly unexported. Multiple issues resolved this cycle. Remaining gaps are pagination and `Lookup` validation -- both P3. |
| **Submissions** | **3/10 -- Alpha** | **Critical P0 bug blocks all submission creation.** Authorization bypass in `Replay`. Worker data flow fundamentally broken. `Process` is a stub. Nearly zero test coverage (only route registration tested). Improved from 2/10: in-memory `Statuses` filtering fixed, unused BSON field removed, service structs properly unexported. Not deployable in any environment. |
| **pkg/** | **5/10 -- Beta** | Cache has no TTL, no miss detection, no shutdown. Middleware errors incorrectly map to 500. `CacheManager` conflates caching with distributed locking (ISP violation). Zero test coverage. The abstractions are well-designed but the implementations have gaps. Session leak in `BeginTx` fixed (unstaged). |

---

## Summary

### Progress Since 5/10

- **`errors.Is` adopted across all services** -- All sentinel error comparisons now correctly use `errors.Is`. The cross-service concern from prior reviews is fully resolved.

- **`InMemoryDataSourceRepository.FindJobs` fully implemented** -- Type, limit, and expiration filtering all working correctly.

- **`DataSourcesJobService.Find` passes `query.Limit`** -- The limit field is now forwarded to the filter.

- **`TenantsService.Update` log ordering fixed** -- Success log now appears after persistence succeeds.

- **Scheduled DTO fields added** -- `IntervalHours` and `ExpirationDate` now included in `scheduledDataSourceAttributesResponse`.

- **Test mock completed** -- `mockDataSourcesRepository` now implements `FindJobs`, satisfying the full interface.

- **Handler channels buffered consistently** -- Submissions handlers now use buffered channels, preventing goroutine leaks on context cancellation.

- **Distributed elector resolved** -- `CacheElector` wrapping `RedisCacheManager` provides distributed leader election via Redis-backed locking.

- **`logFindFormByIDError` fixed** -- Now uses `errors.Is` correctly.

- **Validation rules aligned** -- DTO and command validation rules for forms are now consistent.

- **Worker pool dispatch timeout** -- Per-job timeouts via `SetTimeout` ensure hung jobs are cancelled and workers return to the pool.

- **In-memory `Statuses` filter implemented** (unstaged) -- `InMemorySubmissionsRepository.Find` now filters by `Statuses` using `slices.Contains`, matching the MongoDB adapter.

- **Unused BSON field removed** (unstaged) -- `IdempotencyID` removed from `submissionAttemptDocument`.

- **In-memory version uniqueness enforced** (unstaged) -- `InMemoryVersionsRepository.Upsert` now checks for duplicate `(form_id, version)` and returns `ErrDuplicateVersion`, matching the MongoDB unique index.

- **Service structs unexported across all services** (unstaged) -- `FormsService`, `TenantsService`, `DataSourcesService`, `DataSourcesJobService`, `SubmissionsService`, and `SubmissionJobsService` all renamed to unexported types. Follows idiomatic Go: consumers interact via port interfaces, not concrete types.

- **Authenticator typo fixed** (unstaged) -- `"placholder"` corrected to `"placeholder"`.

- **`mockDatabase.Close` signature fixed** (unstaged) -- Now matches the `database.Database` interface with `Close(context.Context) error`.

- **`ClaimsContextKey` unexported** (unstaged) -- Renamed to `claimsContextKey`, matching the unexported context key pattern used in `httputil`.

- **`BeginTx` session leak fixed** (unstaged) -- `session.EndSession(ctx)` now called on the `StartTransaction` error path, preventing session leaks when transaction start fails.

### Current State

**20 remaining issues** (6 carried from 5/10; 14 newly identified; 19 resolved this cycle). 1 P0, 2 P1, 4 P2, 13 P3.

**Forms Service** remains the most mature with handler-level tests and a complete domain model. The `Form.Update` dirty-state mutation bug is the only remaining P2. In-memory version uniqueness enforcement now matches MongoDB. Service struct properly unexported. No service or domain layer tests exist. Production readiness improved to 8/10.

**Tenants Service** has made strong progress this cycle with 6 issues resolved. The background job processing pipeline and leader election are fully functional. Service structs properly unexported. Only two P3 items remain (pagination and `Lookup` validation). Production readiness at 7/10.

**Submissions Service** has a critical P0 bug that blocks all creation. The `FindByIdempotencyID` error handling assumes `nil, nil` for "not found" but repositories return `ErrNotFound`. Additionally, `Replay` has a tenant authorization bypass, the worker data flow is broken (submission not passed to `Process`), and both `Process` stubs remain. Unstaged improvements include `Statuses` filtering, unused BSON field cleanup, and service struct unexport. This service needs significant work before deployment.

**pkg/** provides well-designed abstractions (`Database`, `CacheManager`, `Elector`, `BackgroundWorker`) with functioning implementations. The `BeginTx` session leak is fixed (unstaged). Remaining gaps: cache entries never expire, cache miss is indistinguishable from zero-value, middleware errors produce 500s, and `CacheManager` conflates caching with distributed locking. Zero test coverage across the entire shared package.

**Hexagonal Architecture** -- All three services maintain correct dependency direction (adapters -> core, never core -> adapters). Port interfaces cleanly separate primary (driving) and secondary (driven) boundaries. The `pkg/` packages serve as infrastructure modules consumed by adapter layers. The unstaged service struct unexport reinforces this: consumers depend on port interfaces, not concrete types.

**DDD** -- Domain entities encapsulate state transitions (e.g., `Version.Publish`, `ScheduledDataSourceAttributes.RefreshData`). The remaining DDD gaps are structural (no domain events, no aggregate-level validation on `Lookup`).

**Idiomatic Go** -- Functional options are applied consistently. Small interfaces follow Go conventions. `errors.Is` now used correctly across all sentinel comparisons. Service structs are unexported with constructors returning interfaces -- the standard Go pattern for implementation hiding.

### Highest-Impact Improvements

1. **Fix `FindByIdempotencyID` error handling in submissions `Create`** (P0 -- all creation blocked)
2. **Add tenant authorization check to `Replay`** (P1 -- security)
3. **Fix submissions worker data flow** (P1 -- pass submission to `Process`)
4. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 -- validation no-op)
5. **Add `ErrMissingIdempotencyHeader`/`ErrMissingTenantID` mappings to `SendErrorResponse`** (P2 -- middleware errors produce 500s)
