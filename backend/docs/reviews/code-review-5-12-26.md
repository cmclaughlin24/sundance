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

5. **`InMemorySubmissionsRepository.Find` ignores `Statuses` filter** (`inmemory/submissions_repository.go:26-39`) -- The in-memory implementation only filters by `TenantID` but completely ignores `filter.Statuses`. The MongoDB implementation correctly handles both. The background worker's `Find` call for `pending` submissions returns all submissions in in-memory mode. *(P2 -- dev/test parity.)*

#### Missing Functionality

6. **`sendErrorResponse` wrapper adds no value** (`handlers.go:194-199`) -- *(Carried from 5/10 #7, P3.)*

7. **`Replay` service method is a stub** (`submissions_service.go:130-145`) -- *(Carried from 5/10 #8, P3.)*

8. **`SubmissionJobsService.Process` is a stub** (`submission_jobs_service.go:35-37`) -- Returns `nil` without doing anything. Separate from the `Replay` stub. *(P3.)*

9. **`submissionAttemptDocument` has unused `IdempotencyID` field** (`mongodb/documents.go:80`) -- Field exists in BSON document struct but is never populated in `toSubmissionAttemptDocument` or read in `fromSubmissionAttemptDocument`. *(P3 -- dead code.)*

#### Code Quality

10. **`Payload` typed as `any`** (`submission.go:36`) -- *(Carried from 5/10 #9, P3.)*

---

### Tenants Service

#### Missing Functionality

11. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

12. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Forms Service

#### Bugs

13. **`Form.Update` mutates fields before validation** (`domain/form.go:53-66`) -- The `Update` method sets fields, then calls `validate.ValidateStruct(f)`. If validation fails, the entity is left in a dirty state with the invalid data already applied. Should validate first or operate on a copy. *(P2.)*

#### Missing Functionality

14. **`InMemoryVersionsRepository` lacks duplicate `(form_id, version)` enforcement** (`inmemory/versions_repository.go:84-97`) -- MongoDB enforces uniqueness via index; in-memory does not. *(P3 -- dev/test parity.)*

#### Code Quality

15. **`FormsService` struct is exported but returned as interface** (`services/forms_service.go:13`) -- `NewFormsService` returns `ports.FormsService`. The struct should be unexported (`formsService`). *(P3.)*

16. **Typo: `"placholder"` in authenticator init** (`routes.go:18`) -- *(P3.)*

---

### pkg/

#### Bugs

17. **`MongoDBDatabase.BeginTx` leaks session on `StartTransaction` failure** (`database/mongodb_database.go:26-41`) -- If `StartSession` succeeds but `StartTransaction` fails, the session is never ended. No `defer session.EndSession(ctx)` on the error path. *(P2.)*

18. **`CacheManager.Get` cannot distinguish cache miss from cached zero-value** (`redis_cache_manager.go:45-69`, `inmemory_cache_manager.go:28-43`) -- Returns `nil` error for both cache miss and cached zero-value. Should return a `bool` hit/miss indicator or a sentinel `ErrCacheMiss`. *(P2.)*

19. **`ErrMissingIdempotencyHeader` and `ErrMissingTenantID` map to 500** (`httputil/http.go:104-109`) -- Neither error matches any `errors.Is` check in `SendErrorResponse`, falling through to the default 500 case. Should be 400 Bad Request. *(P2.)*

#### Architectural

20. **`BackgroundWorker.Start` releases leadership on shutdown but not on `workFn` errors** (`background_worker.go:98`) -- *(Carried from 5/10 #5, P3.)*

21. **Worker pool dispatch blocks indefinitely if all workers are busy** (`background_worker.go:131-135`) -- *(Carried from 5/10 #6, P3.)*

22. **`CacheManager` interface conflates caching with distributed locking** (`cache/cache.go:29-36`) -- `Get/Set/Del` (caching) bundled with `AcquireLock/RenewLock/ReleaseLock` (distributed locking). These are separate concerns; violates ISP. *(P2.)*

#### Missing Functionality

23. **No `Close()` on `CacheManager`** (`cache/cache.go`) -- Redis connection never cleanly shut down. Compare with `Database` which has `Close`. *(P3.)*

24. **`RedisCacheManager.Set` hardcodes TTL=0** (`redis_cache_manager.go:79`) -- All cache entries stored with no expiry. Cache grows unboundedly. *(P3.)*

#### Code Quality

25. **`ClaimsContextKey` is exported** (`auth/claims.go:9-11`) -- Context key types should be unexported to prevent external manipulation. Other context key types in `httputil` correctly use unexported types. *(P3.)*

26. **Zero test files in entire `pkg/` directory.** *(P3.)*

---

### Cross-Service

#### Architectural

27. **Test coverage gaps** -- Submissions has no handler or service tests. No domain-layer or repository-layer tests exist. *(Carried from 5/10 #10, P3.)*

28. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

29. **No real authentication** -- Placeholder only. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 1 | `Create` blocked by `FindByIdempotencyID` error handling | Submissions |
| **P1** | 2 | `Replay` tenant authorization bypass | Submissions |
| **P1** | 3 | Worker data flow broken -- submission not passed to `Process` | Submissions |
| **P2** | 4 | `ReplaySubmissionCommand` no validation tags | Submissions |
| **P2** | 5 | In-memory `Find` ignores `Statuses` filter | Submissions |
| **P2** | 13 | `Form.Update` mutates before validation | Forms |
| **P2** | 17 | Session leak on `StartTransaction` failure | pkg/database |
| **P2** | 18 | Cache miss indistinguishable from zero-value | pkg/cache |
| **P2** | 19 | Middleware errors map to 500 instead of 400 | pkg/common |
| **P2** | 22 | `CacheManager` conflates caching with locking (ISP) | pkg/cache |
| **P3** | 6 | `sendErrorResponse` wrapper is no-op | Submissions |
| **P3** | 7 | `Replay` is a stub | Submissions |
| **P3** | 8 | `SubmissionJobsService.Process` is a stub | Submissions |
| **P3** | 9 | Unused `IdempotencyID` in attempt document | Submissions |
| **P3** | 10 | `Payload` typed as `any` | Submissions |
| **P3** | 11 | `Find()` no pagination | Tenants |
| **P3** | 12 | `Lookup` no validation | Tenants |
| **P3** | 14 | In-memory lacks unique version enforcement | Forms |
| **P3** | 15 | `FormsService` exported unnecessarily | Forms |
| **P3** | 16 | Typo: `"placholder"` | Forms |
| **P3** | 20 | Leadership held despite failures | pkg/worker |
| **P3** | 21 | Pool dispatch blocks indefinitely | pkg/worker |
| **P3** | 23 | No `Close()` on `CacheManager` | pkg/cache |
| **P3** | 24 | Redis cache has no TTL | pkg/cache |
| **P3** | 25 | Exported `ClaimsContextKey` | pkg/auth |
| **P3** | 26 | Zero tests in `pkg/` | pkg/ |
| **P3** | 27 | Test coverage gaps | All |
| **P3** | 28 | No domain events | All |
| **P3** | 29 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **7/10 -- Beta** | Most mature service. Complete CRUD with versioning lifecycle (create, publish, retire). Handler tests provide good HTTP layer coverage. Blocked from production by: no service/domain layer tests, `Form.Update` dirty-state bug, no real authentication. |
| **Tenants** | **7/10 -- Beta** | Fully functional including background job processing pipeline and leader election infrastructure. Clean hexagonal structure with strategies pattern. Multiple issues resolved this cycle. Remaining gaps are pagination and `Lookup` validation -- both P3. |
| **Submissions** | **2/10 -- Alpha** | **Critical P0 bug blocks all submission creation.** Authorization bypass in `Replay`. Worker data flow fundamentally broken. `Process` is a stub. Nearly zero test coverage (only route registration tested). Not deployable in any environment. |
| **pkg/** | **5/10 -- Beta** | `BeginTx` leaks sessions on transaction start failure. Cache has no TTL, no miss detection, no shutdown. Middleware errors incorrectly map to 500. Zero test coverage. The abstractions are well-designed but the implementations have gaps. |

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

### Current State

**29 remaining issues** (6 carried from 5/10; 23 newly identified; 10 resolved this cycle). 1 P0, 2 P1, 6 P2, 20 P3.

**Forms Service** remains the most mature with handler-level tests and a complete domain model. The `Form.Update` dirty-state mutation bug is the only remaining P2. No service or domain layer tests exist.

**Tenants Service** has made strong progress this cycle with 6 issues resolved. The background job processing pipeline and leader election are fully functional. Only two P3 items remain (pagination and `Lookup` validation). Production readiness improved to 7/10.

**Submissions Service** has a critical P0 bug that blocks all creation. The `FindByIdempotencyID` error handling assumes `nil, nil` for "not found" but repositories return `ErrNotFound`. Additionally, `Replay` has a tenant authorization bypass, the worker data flow is broken (submission not passed to `Process`), and both `Process` stubs remain. This service needs significant work before deployment.

**pkg/** provides well-designed abstractions (`Database`, `CacheManager`, `Elector`, `BackgroundWorker`) with functioning implementations. `BeginTx` leaks sessions on transaction start failure, cache entries never expire, and middleware errors produce 500s. Zero test coverage across the entire shared package.

**Hexagonal Architecture** -- All three services maintain correct dependency direction (adapters -> core, never core -> adapters). Port interfaces cleanly separate primary (driving) and secondary (driven) boundaries. The `pkg/` packages serve as infrastructure modules consumed by adapter layers.

**DDD** -- Domain entities encapsulate state transitions (e.g., `Version.Publish`, `ScheduledDataSourceAttributes.RefreshData`). The remaining DDD gaps are structural (no domain events, no aggregate-level validation on `Lookup`).

**Idiomatic Go** -- Functional options are applied consistently. Small interfaces follow Go conventions. `errors.Is` now used correctly across all sentinel comparisons.

### Highest-Impact Improvements

1. **Fix `FindByIdempotencyID` error handling in submissions `Create`** (P0 -- all creation blocked)
2. **Add tenant authorization check to `Replay`** (P1 -- security)
3. **Fix submissions worker data flow** (P1 -- pass submission to `Process`)
4. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 -- validation no-op)
5. **Fix `BeginTx` session leak** (P2 -- resource leak on transaction failure)
