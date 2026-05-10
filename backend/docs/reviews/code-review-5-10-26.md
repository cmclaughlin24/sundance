# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/8 Review

1. ~~`sendErrorResponse` has no domain error mapping~~ (Submissions, P2) -- The submissions handler's `sendErrorResponse` now delegates directly to `httputil.SendErrorResponse`, which correctly maps `ErrNotFound` → 404, `ErrUnauthorized` → 401, validation errors → 400, `ErrExists` → 409. The empty switch with only a default case is effectively a passthrough to the shared utility that handles all cases. *(Resolves 5/8 #4.)*

2. ~~`DataSourceJobsService.Find` has no filtering~~ (Tenants, P3) -- `DataSourcesJobService.Find` now calls `repository.FindJobs(ctx, &ports.FindDataSourceJobsFilter{})`. The MongoDB implementation filters by `type: "scheduled"` and expired/missing `expirationDate`, ensuring only data sources due for refresh are returned. *(Resolves 5/8 #2.)*

3. ~~`DataSourceJobsService.Process` is a stub~~ (Tenants, P3) -- `Process` now fetches lookups via `LookupClient.FetchLookups`, updates the data source's `Data` and `ExpirationDate`, and persists via `Upsert`. Background job processing is functional end-to-end. *(Resolves 5/8 highest-impact #2 and #3.)*

4. ~~`ExpirationDate` drift risk~~ (Tenants, P2) -- `data_source_jobs_service.go:60` now uses `time.Now().Add(...)` instead of `attr.ExpirationDate.Add(...)`, preventing backward drift when the worker runs late. *(Found and fixed same cycle.)*

5. ~~`DataSourcesJobService.Process` missing error context~~ (Tenants, P3) -- Error log at `data_source_jobs_service.go:49` now includes `data_source_id` and `ds.Type` for observability when `GetDataSourceAttributes` fails. *(Found and fixed same cycle.)*

6. ~~`NewDataSourcesBackgroundWorker` panics on bootstrap error~~ (Tenants, P3) -- `NewDataSourcesBackgroundWorker` now returns `(*BackgroundWorker, error)`. The adapter no longer panics; `main.go` handles the error at the composition root, consistent with the persistence bootstrap pattern. *(Found and fixed same cycle.)*

---

## Will Not Fix

See [5/8 review](code-review-5-8-26.md) for the full Will Not Fix list.

`LookupClient` used by both webhook strategy and job service -- The webhook strategy performs a live on-demand lookup while the scheduled job service performs a fetch-and-persist-with-expiration cycle. These are fundamentally different operations sharing only the HTTP client port, not a DRY concern.

---

## Remaining Issues

### Tenants Service

#### Bugs

1. **`DataSourcesJobService.Process` mutates domain entity in-place without domain method** (`data_source_jobs_service.go:53-56`) -- The service directly modifies `attr.Data`, `attr.ExpirationDate`, and `ds.Attributes` rather than using a domain method (e.g., `ds.RefreshLookups(lookups, nextExpiration)`). This bypasses any future domain invariants and violates DDD's encapsulation principle. The `FIXME` comment acknowledges this. *(P2 -- domain model integrity.)*

2. **`InMemoryDataSourceRepository.FindJobs` does not filter** (`inmemory/data_sources_repository.go:64-72`) -- Returns all data sources regardless of type or expiration, diverging from the MongoDB implementation's behavior. Tests against in-memory will not catch filtering bugs. *(P2 -- dev/test parity.)*

#### Missing Functionality

3. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/8 #1, P3.)*

4. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/8 #3, P3.)*

### pkg/worker

#### Architectural

5. **`Elector` interface has no distributed implementation** (`elector.go`) -- `InMemoryElector` always returns `true`, meaning multiple replicas will all process jobs concurrently. The interface is well-designed for future extension (MongoDB/Redis-based election), but there is no protection against duplicate job processing in a multi-replica deployment today. *(P3 -- noted as in-progress/planned.)*

6. **`BackgroundWorker.Start` releases leadership on shutdown but not on `workFn` errors** (`background_worker.go:98`) -- If `workFn` consistently fails, the worker holds leadership indefinitely without doing useful work. Other replicas cannot take over. Consider releasing leadership after N consecutive failures. *(P3 -- resilience.)*

7. **Worker pool dispatch blocks indefinitely if all workers are busy** (`background_worker.go:131-135`) -- The `w := <-pool` receive will block until a worker becomes available. If all workers are processing long-running jobs, the ticker tick is effectively "lost" but the goroutine is stuck waiting. This is acceptable for the current use case (small pool, fast jobs) but could benefit from a `select` with a timeout or skip for observability. *(P3 -- future concern.)*

### Submissions Service

#### Architectural

8. **`sendErrorResponse` wrapper adds no value** (`handlers.go:194-198`) -- The switch statement has only a `default` case delegating to `httputil.SendErrorResponse`. This indirection is unnecessary. Either add service-specific error cases (justifying the wrapper) or call `httputil.SendErrorResponse` directly. *(P3 -- dead code.)*

#### Missing Functionality

9. **`Replay` service method is a stub** (`submissions_service.go`) -- *(Carried from 5/8 #5, P3.)*

#### Code Quality

10. **`Payload` typed as `any`** (`submission.go`) -- *(Carried from 5/8 #6, P3.)*

### Cross-Service

#### Architectural

11. **Test coverage gaps** -- Submissions has no handler or service tests. No domain-layer or repository-layer tests exist. *(Carried from 5/8 #7, P3.)*

12. **No domain events** for cross-service communication. *(Carried from 5/8 #8, P3.)*

13. **No real authentication** -- Placeholder only. *(Carried from 5/8 #9, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | Domain mutation without domain method | Tenants |
| **P2** | 2 | In-memory `FindJobs` doesn't filter | Tenants |
| **P3** | 3 | `Find()` no pagination | Tenants |
| **P3** | 4 | `Lookup` no validation | Tenants |
| **P3** | 5 | No distributed elector implementation | pkg/worker |
| **P3** | 6 | Leadership held despite failures | pkg/worker |
| **P3** | 7 | Pool dispatch blocks indefinitely | pkg/worker |
| **P3** | 8 | `sendErrorResponse` wrapper is no-op | Submissions |
| **P3** | 9 | `Replay` is a stub | Submissions |
| **P3** | 10 | `Payload` typed as `any` | Submissions |
| **P3** | 11 | Test coverage gaps | All |
| **P3** | 12 | No domain events | All |
| **P3** | 13 | No real authentication | All |

---

## Summary

### Progress Since 5/8

- **Leader election introduced** (`pkg/worker/elector.go`) -- New `Elector` interface with `TryAcquire`, `Renew`, and `Release` methods. `InMemoryElector` provides a single-instance default. `BackgroundWorker` integrates election: acquires leadership before dispatching jobs, renews each tick, releases on shutdown. This correctly separates the election concern from job processing.

- **Clients adapter extracted** (Tenants) -- New `adapters/clients/` package with `LookupClient` implementing the `ports.LookupClient` interface. Uses `http.Client` with configurable timeout. Decodes JSON responses into domain `Lookup` entities. Bootstrap function uses functional options.

- **Data source job processing fully implemented** (Tenants) -- `DataSourcesJobService.Process` now performs the full cycle: validate command → extract scheduled attributes → fetch lookups via HTTP client → update attributes with new data and next expiration → persist via repository. The background worker is now a complete, functional system for refreshing scheduled data sources.

- **`FindJobs` repository query implemented** (MongoDB) -- Filters to `type: "scheduled"` and `expirationDate` that is null, missing, or past. Supports optional limit via `FindDataSourceJobsFilter`.

- **Webhook lookup strategy refactored** -- Extracted HTTP fetch logic into `LookupClient` port, making `WebhookLookupStrategy` a thin adapter that delegates to the client. Reduces duplication and improves testability.

- **`sendErrorResponse` effectively resolved** (Submissions) -- While the wrapper function still exists as a no-op passthrough, `httputil.SendErrorResponse` now handles all domain error mapping centrally. The P2 issue (incorrect status codes) is functionally resolved.

- **Bootstrap error handling improved** (Tenants) -- `NewDataSourcesBackgroundWorker` now returns `(*BackgroundWorker, error)` instead of panicking, consistent with the persistence adapter pattern.

### Current State

**13 remaining issues** (9 from 5/8 → resolved 3, carried 6 unchanged, introduced 10 new issues; 3 fixed same-cycle; 1 moved to Will Not Fix).

**Forms Service** remains fully mature. No remaining issues.

**Tenants Service** has taken a major step forward with a fully functional background job processing pipeline. The `DataSourcesJobService` fetches scheduled data sources due for refresh, delegates to `LookupClient` to retrieve fresh lookup data via HTTP, and persists the updated attributes. Leader election ensures only one replica processes jobs (once a distributed elector is implemented). The `Clients` adapter cleanly separates HTTP concerns from business logic. The primary remaining concern is domain encapsulation (direct attribute mutation).

**Submissions Service** error mapping is now correct via the shared `httputil.SendErrorResponse`. The remaining wrapper function is dead code but harmless. Core issues remain the `Replay` stub and untyped `Payload`.

**pkg/worker** now includes leader election infrastructure. The `Elector` interface is cleanly designed for pluggable implementations (in-memory for dev, MongoDB/Redis for production). The `BackgroundWorker` correctly integrates election into its tick loop with acquire/renew/release lifecycle. No distributed implementation exists yet, but the architecture is ready for one.

**Hexagonal Architecture** -- The new `adapters/clients/` directory correctly implements a driven adapter. It depends inward on `core/ports/LookupClient` interface. The `strategies` and `services` packages consume it through the port boundary, never directly importing the adapter. The `Elector` interface in `pkg/worker` follows the same port pattern — infrastructure-agnostic with a pluggable implementation.

**DDD** -- The main DDD concern is the direct mutation of `ScheduledDataSourceAttributes` in `DataSourcesJobService.Process`. The `FIXME` comment indicates awareness. A proper approach would be a domain method like `DataSource.RefreshLookups(data []*Lookup, nextExpiration time.Time)` that encapsulates the state transition and could enforce invariants (e.g., non-empty data, future expiration).

**Idiomatic Go** -- The `Elector` interface follows Go's small-interface preference (3 methods, single responsibility). Functional options continue to be applied consistently. The `LookupClient` adapter correctly uses `http.NewRequestWithContext` for cancellation propagation. `context.Background()` is appropriately used for the leadership release on shutdown (parent context is already cancelled).

### Highest-Impact Improvements

1. **Encapsulate data source refresh in a domain method** (P2 -- fixes #1; use `domain.Now()` for next expiration)
2. **Implement in-memory `FindJobs` filtering** (P2 -- dev/test parity with MongoDB)
3. **Add a distributed `Elector` implementation** (P3 -- required before multi-replica deployment)
