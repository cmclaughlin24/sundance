# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/12 Review

1. ~~`replaySubmission` returns 201 instead of 202~~ (Submissions, P3) -- `replaySubmission` handler now returns `202 Accepted` instead of `201 Created`, which is semantically correct for an asynchronous operation that triggers reprocessing. *(Committed.)*

2. ~~Submissions service missing auth middleware~~ (Submissions, P3) -- `PlaceholderAuthenticator` and `auth.NewMiddleware` now applied to submissions routes, matching the pattern in forms and tenants. *(Committed.)*

3. ~~Route middleware applied globally blocks non-API paths~~ (Forms/Submissions, P3) -- `TenantMiddleware` and `AuthMiddleware` moved from global mux-level to scoped within `/api/v1` routes. Swagger endpoints at `/swagger/*` no longer require auth/tenant headers. *(Committed.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`submissionJob.Process` does not pass the held `*domain.Submission` to the service** (`workers/submissions_worker.go:26-28`) -- The `submissionJob` holds a `*domain.Submission` in field `s`, but `Process` calls `j.service.Process(ctx)` without passing it. The `SubmissionJobsService.Process` method has no way to know which submission to process. Worker data flow is fundamentally broken. *(P1 -- worker non-functional.)*

2. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go:29-32`) -- `TenantID` and `ID` fields have no `validate:"required"` tags, so `validate.ValidateStruct(command)` always passes even with empty values. *(P2.)*

#### Missing Functionality

3. **`sendErrorResponse` wrapper adds no value** (`handlers.go:194-199`) -- *(Carried from 5/10 #7, P3.)*

4. **`SubmissionJobsService.Process` is a stub** (`submission_jobs_service.go:35-37`) -- Returns `nil` without doing anything. *(P3.)*

#### Code Quality

5. **`Payload` typed as `any`** (`submission.go:36`) -- *(Carried from 5/10 #9, P3.)*

---

### Tenants Service

#### Missing Functionality

6. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

7. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Forms Service

#### Bugs

8. **`Form.Update` mutates fields before validation** (`domain/form.go:53-66`) -- The `Update` method sets fields, then calls `validate.ValidateStruct(f)`. If validation fails, the entity is left in a dirty state with the invalid data already applied. Should validate first or operate on a copy. *(P2.)*

---

### Cross-Service

#### Architectural

9. **Test coverage gaps** -- Submissions has no handler or service tests. No domain-layer or repository-layer tests exist across services. Zero test files in entire `pkg/` directory. *(Carried from 5/10 #10, P3.)*

10. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

11. **No real authentication** -- Placeholder only. All three services now have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 1 | Worker data flow broken -- submission not passed to `Process` | Submissions |
| **P2** | 2 | `ReplaySubmissionCommand` no validation tags | Submissions |
| **P2** | 8 | `Form.Update` mutates before validation | Forms |
| **P3** | 3 | `sendErrorResponse` wrapper is no-op | Submissions |
| **P3** | 4 | `SubmissionJobsService.Process` is a stub | Submissions |
| **P3** | 5 | `Payload` typed as `any` | Submissions |
| **P3** | 6 | `Find()` no pagination | Tenants |
| **P3** | 7 | `Lookup` no validation | Tenants |
| **P3** | 9 | Test coverage gaps | All |
| **P3** | 10 | No domain events | All |
| **P3** | 11 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8/10 -- Beta** | Most mature service. Complete CRUD with versioning lifecycle (create, publish, retire). Handler tests provide good HTTP layer coverage. OpenAPI documentation with Swagger UI. In-memory adapter enforces version uniqueness. Service structs properly unexported. Blocked from production by: no service/domain layer tests, `Form.Update` dirty-state bug, no real authentication. |
| **Tenants** | **7/10 -- Beta** | Fully functional including background job processing pipeline, leader election, and data source strategies. OpenAPI documentation with Swagger UI. Clean hexagonal structure. Service structs properly unexported. Remaining gaps are pagination and `Lookup` validation -- both P3. |
| **Submissions** | **5/10 -- Beta** | Core submission lifecycle functional: create (with idempotency), find, replay (with tenant auth). OpenAPI documentation with Swagger UI. Auth middleware now wired. Worker data flow still broken (P1). `SubmissionJobsService.Process` is a stub. Nearly zero test coverage. Not deployable due to P1 worker issue. |
| **pkg/** | **8/10 -- Production-Ready** | All previously identified bugs and architectural issues resolved. `CacheManager` cleanly separated from `CacheLocker` (ISP). `CacheCloser` enables clean shutdown. TTL parameter on `Set`. `ErrCacheMiss` for cache miss detection. Session leak fixed. Middleware errors correctly map to 400. Worker failover with configurable failure limits. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/12

- **OpenAPI documentation added to all three services** (committed) -- Swagger annotations on all handler methods using `swag` comment syntax. Generated `docs.go`, `swagger.json`, and `swagger.yaml` for each service. Swagger UI served at `/swagger/*` with proper service titles, descriptions, and base paths. This significantly improves API discoverability and developer onboarding.

- **Route middleware scoping** (committed) -- `TenantMiddleware` and `AuthMiddleware` moved from global to `/api/v1` scope in forms and submissions. Swagger endpoints are accessible without auth/tenant headers. This is the correct pattern for serving documentation alongside protected API routes.

- **Submissions auth middleware** (committed) -- `PlaceholderAuthenticator` and `auth.NewMiddleware` added to submissions routes. All three services now consistently apply authentication middleware (still placeholder implementation).

- **`replaySubmission` returns 202** (committed) -- Corrected from `201 Created` to `202 Accepted`, semantically appropriate for triggering async reprocessing.

### Current State

**11 remaining issues** (4 carried from 5/10; 7 newly identified in 5/12 cycle). 0 P0, 1 P1, 2 P2, 8 P3.

**Forms Service** remains the most mature at 8/10. OpenAPI documentation adds API discoverability. The `Form.Update` dirty-state mutation bug is the only remaining P2. No service or domain layer tests exist.

**Tenants Service** at 7/10 with full background job pipeline, leader election, and now OpenAPI documentation. Only two P3 items remain (pagination and `Lookup` validation).

**Submissions Service** at 5/10. Core submission lifecycle is functional with idempotent creation, tenant-authorized replay, and reference ID lookup. Auth middleware now wired. OpenAPI documentation added. The remaining P1 is the worker data flow -- `submissionJob.Process` does not pass the held submission to the service. `SubmissionJobsService.Process` remains a stub.

**pkg/** at 8/10 -- all bugs and architectural issues from prior reviews resolved. Clean interface separation, proper shutdown, TTL support, cache miss detection, worker failover. Only test coverage remains.

**Hexagonal Architecture** -- All three services maintain correct dependency direction (adapters -> core, never core -> adapters). Port interfaces cleanly separate primary (driving) and secondary (driven) boundaries. The `pkg/` packages serve as infrastructure modules consumed by adapter layers. OpenAPI annotations live correctly in the REST adapter layer.

**DDD** -- Domain entities encapsulate state transitions (e.g., `Version.Publish`, `Submission.Reset`, `ScheduledDataSourceAttributes.RefreshData`). The remaining DDD gaps are structural (no domain events, no aggregate-level validation on `Lookup`).

**Idiomatic Go** -- Functional options applied consistently. Small interfaces follow Go conventions. `errors.Is` used correctly. Service structs unexported with constructors returning interfaces. OpenAPI annotations follow `swag` conventions with proper tag grouping.

### Highest-Impact Improvements

1. **Fix submissions worker data flow** (P1 -- pass submission to `Process`)
2. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 -- validation no-op)
3. **Fix `Form.Update` dirty-state mutation** (P2 -- validate before mutating)
4. **Add test coverage** (P3 -- zero tests in `pkg/`, no service/domain tests across services)
5. **Implement `SubmissionJobsService.Process`** (P3 -- worker pipeline incomplete without it)
