# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/7 Review

1. ~~`createSubmission` handler hardcodes `IdempotencyID("")`~~ (Submissions, P1) -- The handler now extracts the idempotency key from context via `httputil.IdempotencyFromContext(r.Context())`. A new `IdempotencyMiddleware` is applied specifically to the `POST /submissions` route, extracting the `Idempotency-Key` header and rejecting requests that omit it. This resolves the P1 blocking bug from 5/7. *(Resolves 5/7 #1.)*

2. ~~`pkg/common/tenants` package had split concerns~~ (Shared, P3) -- The tenant middleware and context functions consolidated into `pkg/common/httputil/tenant_middleware.go`. The separate `pkg/common/tenants` package eliminated. `NewMiddleware` renamed to `NewTenantMiddleware` for clarity. All service import paths updated. *(Not previously tracked.)*

3. ~~`IdempotencyMiddlware` typo~~ (Shared, P3) -- Renamed to `IdempotencyMiddleware`. Call site in submissions routes updated. *(Found and fixed same cycle.)*

4. ~~`Worker.Start` passes original `ctx` to `job.Process` instead of `wctx`~~ (pkg/worker, P3) -- Now correctly passes `wctx` (with worker ID) to `job.Process`. Jobs have access to the worker ID in their context. *(Found and fixed same cycle.)*

5. ~~`pkg/worker/go.mod` missing `google/uuid` dependency~~ (pkg/worker, P2) -- `go.mod` already had `require github.com/google/uuid v1.6.0`. Mistakenly flagged. *(False positive.)*

6. ~~`context.go` uses wrong key name~~ (pkg/worker, P3) -- `workerIDKey` string value changed from `"tenantID"` to `"workerID"`. *(Found and fixed same cycle.)*

7. ~~`BackgroundWorkerBuilder.Build()` has no validation~~ (pkg/worker, P2) -- `Build()` now returns `(*BackgroundWorker[J], error)`. Validates: nil `logger` returns `ErrLoggerIsRequired`, nil `workFn` returns `ErrWorkFnIsRequired`. Zero `interval` defaults to 1 minute, zero `size` defaults to 5. Eliminates panic/deadlock risk from misconfiguration. Builder consolidated from separate file into `background_worker.go`. *(Resolves 5/8 #1.)*

8. ~~Worker lacks error handling, timeout, and panic resilience~~ (pkg/worker, P2) -- `Job.Process` now returns `error`; workers log job failures via `ErrorContext`. Per-job timeout support added via `SetTimeout` on the builder, applied with `context.WithTimeout`. `Worker` refactored to functional options pattern (`WorkerWithPool`, `WorkerWithLogger`, `WorkerWithTimeout`). Panic recovery moved inside the `for` loop so a panicking job does not kill the worker goroutine — the worker logs the panic and continues processing. *(Found and fixed same cycle.)*

9. ~~`BackgroundWorker.Start` has no error logging on `workFn` failure~~ (pkg/worker, P2) -- `workFn` errors now logged via `wp.logger.WarnContext(ctx, "failed to fetch jobs", "error", err)` before continuing. Job dispatch count logged at debug level. Worker start logged at info level with pool size and interval. *(Resolves 5/8 #1.)*

10. ~~`dataSourceJob.Process` is a no-op~~ (Tenants, P3) -- The job now delegates to `DataSourceJobsService.Process` via the port interface. `newDataSourceJob` factory injects the service dependency. The `WorkFn` calls `app.Services.DataSourceJobs.Find(ctx)` to produce jobs. A new `DataSourcesJobService` implementation provides structured logging and command validation. *(Resolves 5/8 #4.)*

11. ~~Bootstrap functions use positional parameters~~ (All, P3) -- All three services (`forms`, `submissions`, `tenants`) and `strategies` refactored to use functional options for `NewApplication`, `services.Bootstrap`, and `strategies.Bootstrap`. More idiomatic Go and enables optional/extensible configuration. *(Not previously tracked.)*

---

## Will Not Fix

See [5/4 review](code-review-5-4-26.md) for the full Will Not Fix list (items #25-27 covering `FindByReferenceID` linear scan, in-memory map key type, and context cancel response pattern).

`IdempotencyFromContext` panics on missing context value -- by design. Matches the `TenantFromContext` fail-fast pattern. The middleware guarantees the value is present via `Idempotency-Key` header enforcement; a panic indicates programmer error (handler called without middleware), not a runtime concern.

---

## Remaining Issues

### Tenants Service

#### Architectural

1. **`Find()` has no pagination or filtering** (`tenants_service.go`). *(Carried from 5/7, P3.)*

2. **`DataSourceJobsService.Find` has no filtering** (`data_source_jobs_service.go`) -- Returns all data sources regardless of type or scheduling state. The worker will re-process all data sources every tick rather than only those due for refresh. *(P3 -- design consideration, acceptable if in-progress.)*

#### Missing Functionality

3. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Carried from 5/7, P3.)*

---

### Submissions Service

#### Architectural

4. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. `common.ErrNotFound` maps to 500 instead of 404; `common.ErrUnauthorized` maps to 500 instead of 403; validation errors map to 500 instead of 400. *(Carried from 5/7, P2.)*

5. **`Replay` service method is a stub** (`submissions_service.go`) -- Validates and fetches but performs no replay. Handler returns 201 "Successfully replayed" misleadingly. *(Carried from 5/7, P3.)*

#### Code Quality

6. **`Payload` typed as `any`** (`submission.go`) -- No type safety across DTO, command, domain, and persistence layers. *(Carried from 5/7, P3.)*

---

### Cross-Service

#### Architectural

7. **Test coverage gaps** -- Submissions has no handler or service tests. Forms has no service tests. No domain-layer or repository-layer tests exist anywhere. *(Carried from 5/7, P3.)*

8. **No domain events** for cross-service communication. *(Carried from 5/7, P3.)*

9. **No real authentication** -- Placeholder only. *(Carried from 5/7, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 4 | `sendErrorResponse` no domain error mapping | Submissions |
| **P3** | 1 | `Find()` no pagination | Tenants |
| **P3** | 2 | `DataSourceJobsService.Find` no filtering | Tenants |
| **P3** | 3 | `Lookup` no validation | Tenants |
| **P3** | 5 | `Replay` is a stub | Submissions |
| **P3** | 6 | `Payload` typed as `any` | Submissions |
| **P3** | 7 | Test coverage gaps | All |
| **P3** | 8 | No domain events | All |
| **P3** | 9 | No real authentication | All |

---

## Summary

### Progress Since 5/7

Multiple commits and working changes since the last review, focused on idempotency support, background worker infrastructure, functional options refactoring, and data source job processing:

- **Idempotency middleware introduced** -- New `IdempotencyMiddleware` in `pkg/common/httputil` extracts the `Idempotency-Key` header, rejects requests missing it, and stores the value in context. Applied per-route to `POST /submissions` via Chi's `With()`. The handler retrieves the key via `IdempotencyFromContext` and passes it through the command to the domain. This completes the idempotency story end-to-end: request header → middleware → handler → command → domain → repository (with unique index enforcement). The P1 blocking bug from 5/7 (empty idempotency key failing validation) is fully resolved.

- **Tenant middleware consolidated** -- The `pkg/common/tenants` package (which contained `context.go` and `middleware.go`) has been eliminated. Both the context helpers (`SetTenantContext`, `TenantFromContext`) and the middleware (`NewTenantMiddleware`, renamed from `NewMiddleware`) now live in `pkg/common/httputil`. This reduces the number of shared packages and co-locates all HTTP-layer context utilities (tenant, idempotency) in one place. All three services updated to import from `httputil` instead of `tenants`.

- **Background worker package introduced** (`pkg/worker`) -- A new shared module providing a generic worker pool with ticker-based job dispatch:
  - `Worker[J Job]` -- individual goroutine worker that pulls jobs from a shared pool channel
  - `BackgroundWorker[J Job]` -- orchestrator that ticks on an interval, calls a `WorkFn` to produce jobs, and dispatches them to available workers
  - `BackgroundWorkerBuilder[J Job]` -- fluent builder for constructing `BackgroundWorker` instances
  - `WorkerContextHandler` -- `slog.Handler` wrapper that injects `worker_id` into log records
  - `Job` interface with a single `Process(context.Context)` method

- **Data sources background worker fully wired** (Tenants) -- `workers/data_sources_worker.go` defines a `dataSourceJob` implementing `Job` that delegates to `DataSourceJobsService.Process` through the port boundary. The `WorkFn` calls `DataSourceJobsService.Find` to produce jobs. A new `DataSourceJobsService` interface added to `ports/primary.go` with `Find` and `Process` methods. `DataSourcesJobService` implementation in the services layer provides structured logging and command validation. `ProcessDataSourceJobCommand` added to `ports/commands.go` with `validate:"required"` on the `DataSource` field.

- **Functional options refactoring** -- All three services refactored: `core.NewApplication`, `services.Bootstrap`, and `strategies.Bootstrap` now use functional options (`WithLogger`, `WithRepository`, `WithServices`, `WithStrategies`, `WithHTTPClient`). Replaces positional parameter constructors with the idiomatic Go options pattern, enabling extensible configuration without breaking changes.

### Current State

**9 remaining issues (5/7) -> 9 remaining issues** (resolved 11; introduced 10 new issues primarily in `pkg/worker`; 9 fixed same-cycle; carried forward 7 unchanged; 1 new design note added).

**Forms Service** remains fully mature. No remaining issues.

**Tenants Service** has a new adapter layer (`adapters/workers/`) introducing background processing capability. The worker infrastructure is fully wired: graceful shutdown via `context.WithCancel`, job dispatch through a `DataSourceJobsService` port, and proper command validation. A new `DataSourceJobsService` interface and implementation provide the service-layer boundary for background job processing. The `Find` and `Process` methods are stubs but correctly structured for incremental implementation.

**Submissions Service** idempotency is now fully functional end-to-end. The `Idempotency-Key` header is required for submission creation, extracted by middleware, passed through the service layer, and enforced at the database level via a unique index. Remaining issues are `sendErrorResponse` mapping (P2), `Replay` stub (P3), and `Payload` typing (P3).

**pkg/worker** has been substantially hardened this cycle and has **no remaining issues**. The builder validates required fields and provides sensible defaults. Workers handle job errors, support per-job timeouts, use functional options, and recover from panics without dying. `BackgroundWorker.Start` logs `workFn` errors at warn level and job dispatch counts at debug level.

**Hexagonal Architecture** -- The new `adapters/workers/` directory in tenants follows the hexagonal pattern correctly: it's an adapter (driving adapter, triggered by time rather than HTTP) that depends inward on `core`. The `dataSourceJob` struct holds a `*domain.DataSource`, keeping the domain model at the center. The `pkg/worker` package is pure infrastructure with no domain knowledge, analogous to `pkg/database`.

**DDD** -- The `Job` interface (`Process(context.Context) error`) is minimal and domain-agnostic. The tenant-specific `dataSourceJob` adapter wraps the domain entity and delegates to `DataSourceJobsService.Process` through the port boundary. The new `ProcessDataSourceJobCommand` carries the domain entity with validation. This correctly separates the scheduling concern (infrastructure) from the business logic (domain/service). The `DataSourceJobsService` port is a clean driving port for the worker adapter, mirroring how `TenantsService` and `DataSourcesService` serve the REST adapter.

**Idiomatic Go** -- The worker pool pattern using channels is a well-established Go concurrency pattern. The `Job` interface with a single method follows Go's preference for small interfaces. The bootstrap functions and `NewApplication` now use functional options (`WithLogger`, `WithRepository`, `WithServices`, `WithStrategies`) — the most idiomatic Go pattern for optional/extensible configuration. The `contextKey` type-safety pattern (unexported type prevents key collisions) is correctly applied in both `httputil` and `worker` packages.

### Highest-Impact Improvements

1. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- error semantics incorrect)
2. **Implement `DataSourceJobsService.Find` with scheduling-aware filtering** (P3 -- currently returns nil)
3. **Implement `DataSourceJobsService.Process` with actual data source refresh logic** (P3 -- currently a stub)
