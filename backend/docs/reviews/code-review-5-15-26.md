# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/14 Review

1. ~~`Payload` typed as `any`~~ (Submissions/Forms, P3) -- `Payload` field in `Submission` domain entity and `CreateSubmissionCommand` is now `map[string]any`, properly constraining it to a JSON object structure. DTOs updated to match. *(Committed. Resolves 5/10 #9.)*

2. ~~`Form.Update` mutates fields before validation~~ (Forms, P2) -- All domain entity `Update` methods across both services now use the dereference copy pattern: `cpy := *entity`, mutate the copy, validate the copy, then `*entity = cpy` only on success. Applied to `Form.Update`, `Tenant.Update`, and `DataSource.Update`. *(Committed.)*

3. ~~Swagger host hardcoded to `localhost`~~ (All, P3) -- All services now derive the Swagger host from configuration, stripping the protocol scheme via regex. Swagger URL uses the configured host dynamically. *(Committed.)*

4. ~~`sendErrorResponse` wrapper adds no value~~ (Submissions/Forms, P3) -- The merged forms service `sendErrorResponse` now maps domain-specific errors (e.g., `ErrVersionLocked`, `ErrInvalidPosition`, `ErrDuplicateVersion`) to 400 Bad Request via a local `isBadRequest` helper, before falling through to `httputil.SendErrorResponse` for other cases. Adds clear value. *(Committed. Resolves 5/10 #7.)*

5. ~~Submissions service as separate bounded context~~ (Architecture) -- The submissions service has been intentionally folded into the forms service. Submissions require form definitions to validate payload structure against the version schema, making them part of the same "Form Builder" bounded context. This reduces operational complexity (one fewer service to deploy) while maintaining hexagonal architecture internally. *(Committed.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`submissionJob.Process` does not pass the held `*domain.Submission` to the service** (`workers/submissions_worker.go:27`) -- The `submissionJob` holds a `*domain.Submission` in field `s`, but `Process` calls `j.service.Process(ctx)` without passing it. The `SubmissionJobsService.Process` method has no way to know which submission to process. Worker data flow is fundamentally broken. *(P1 -- worker non-functional.)*

2. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go`) -- `TenantID` and `ID` fields have no `validate:"required"` tags, so `validate.ValidateStruct(command)` always passes even with empty values. *(P2.)*

#### Missing Functionality

3. **`SubmissionJobsService.Process` is a stub** (`submission_jobs_service.go`) -- Returns `nil` without doing anything. The worker pipeline is incomplete without an actual processing implementation. *(P3.)*

---

### Tenants Service

#### Missing Functionality

4. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

5. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Cross-Service

#### Architectural

6. **Test coverage gaps** -- No domain-layer, service-layer, or repository-layer tests exist across services. Zero test files in entire `pkg/` directory. Forms service has handler tests only. *(Carried from 5/10 #10, P3.)*

7. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

8. **No real authentication** -- Placeholder only. All services have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 1 | Worker data flow broken -- submission not passed to `Process` | Forms |
| **P2** | 2 | `ReplaySubmissionCommand` no validation tags | Forms |
| **P3** | 3 | `SubmissionJobsService.Process` is a stub | Forms |
| **P3** | 4 | `Find()` no pagination | Tenants |
| **P3** | 5 | `Lookup` no validation | Tenants |
| **P3** | 6 | Test coverage gaps | All |
| **P3** | 7 | No domain events | All |
| **P3** | 8 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **7/10 -- Beta** | Unified "Form Builder" bounded context encompassing forms, versions, and submissions. Complete CRUD with versioning lifecycle. Idempotent submission creation, tenant-authorized replay, reference ID lookup. OpenAPI documentation with Swagger UI. All domain `Update` methods use validate-before-mutate pattern. Handler-level `sendErrorResponse` maps domain errors to proper HTTP status codes. Worker data flow still broken (P1). `SubmissionJobsService.Process` is a stub. Handler tests exist but no service/domain/repository tests. |
| **Tenants** | **8/10 -- Production-Ready** | Fully functional including background job processing pipeline, leader election, and data source strategies. OpenAPI documentation with Swagger UI. All domain `Update` methods use validate-before-mutate pattern. Service structs properly unexported. Only P3 gaps remain (pagination and `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | All previously identified bugs and architectural issues resolved. `CacheManager` cleanly separated from `CacheLocker`. `CacheCloser` enables clean shutdown. TTL parameter on `Set`. `ErrCacheMiss` for cache miss detection. Session leak fixed. Middleware errors correctly map to 400. Worker failover with configurable failure limits. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/14

- **Submissions service merged into Forms** (committed) -- The entire `services/submissions/` directory has been deleted and its domain, ports, services, repositories, workers, and handlers merged into the forms service. This is architecturally correct: submissions require form definitions to validate payload structure against the version schema, making them part of the same bounded context. Reduces operational complexity from 3 services to 2.

- **`Payload` typed as `map[string]any`** (committed) -- Previously typed as `any`, now properly constrained to a JSON object structure matching form field definitions.

- **Validate-before-mutate applied across all domain entities** (committed) -- `Form.Update`, `Tenant.Update`, and `DataSource.Update` all use the dereference copy pattern (`cpy := *entity`, mutate copy, validate copy, `*entity = cpy` on success). No domain entity can be left in a dirty state on validation failure.

- **Hardcoded Swagger host removed** (committed) -- All services derive the host from configuration.

- **`sendErrorResponse` now maps domain errors** (committed) -- The forms service `sendErrorResponse` maps 14 domain-specific errors to 400 Bad Request via `isBadRequest`, then delegates to `httputil.SendErrorResponse` for other cases. No longer a no-op wrapper.

### Current State

**8 remaining issues** (3 carried from 5/10; 5 newly identified in prior cycles). 0 P0, 1 P1, 1 P2, 6 P3.

**Forms Service** (now includes submissions) at 7/10. The merge into a unified "Form Builder" bounded context is architecturally sound and reduces deployment complexity. Core functionality is complete: form CRUD, versioning lifecycle, idempotent submission creation, tenant-authorized replay, and reference ID lookup. The remaining P1 is the worker data flow — `submissionJob.Process` does not pass the held submission to the service. `SubmissionJobsService.Process` remains a stub. OpenAPI documentation, validate-before-mutate pattern, and proper error mapping all in place.

**Tenants Service** at 8/10. All domain `Update` methods now use the dereference copy pattern. Background job pipeline, leader election, and data source strategies all functional. Only P3 gaps remain.

**pkg/** at 8/10 — all bugs and architectural issues resolved. Only test coverage remains.

**Hexagonal Architecture** -- Both services maintain correct dependency direction (adapters -> core, never core -> adapters). Port interfaces cleanly separate primary and secondary boundaries. The forms-submissions merge maintains internal hexagonal structure: submission ports, services, and repositories are distinct from form equivalents within the same service boundary.

**DDD** -- The forms service now represents a single "Form Builder" bounded context. Submissions are a downstream lifecycle concern of forms — they reference form definitions and version schemas for payload validation. Domain entities encapsulate state transitions (`Version.Publish`, `Submission.Reset`, `ScheduledDataSourceAttributes.RefreshData`). All `Update` methods protect invariants via validate-before-mutate.

**Idiomatic Go** -- Functional options applied consistently. Small interfaces follow Go conventions. `errors.Is` used correctly. Service structs unexported with constructors returning interfaces. Dereference copy pattern for safe mutation.

### Highest-Impact Improvements

1. **Fix submissions worker data flow** (P1 -- pass submission to `Process`)
2. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 -- validation no-op)
3. **Implement `SubmissionJobsService.Process`** (P3 -- worker pipeline incomplete)
4. **Add test coverage** (P3 -- zero tests in `pkg/`, no service/domain tests)
5. **Replace placeholder authentication** (P3 -- all services wired but using `PlaceholderAuthenticator`)
