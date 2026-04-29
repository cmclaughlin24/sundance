# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/26 Review

1. ~~No graceful shutdown~~ (All, P3) -- All three services (`forms/cmd/server/main.go`, `submissions/cmd/server/main.go`, `tenants/cmd/server/main.go`) now listen for `os.Interrupt`/`SIGTERM` and call `server.Shutdown(ctx)` with a 15-second timeout. *(Resolves 4/26 #23.)*

2. ~~`Lookup` service method is a stub~~ (Tenants, P3) -- `DataSourcesService.Lookup` (`data_sources_service.go:125-139`) now resolves a `LookupStrategy` from the strategy registry by data source type and delegates to `strategy.Lookup(ctx, ds)`. Three strategy implementations exist: `StaticLookupStrategy`, `ScheduledLookupStrategy`, and `WebhookLookupStrategy`. *(Resolves 4/26 #18.)*

3. ~~No domain constructors~~ (Submissions, P2) -- `NewSubmission` (`submission.go:33-50`) now generates IDs, sets `CreatedAt`, initializes `Attempts`, and validates via `validate.ValidateStruct`. `HydrateSubmission` (`submission.go:52-73`) provides a persistence reconstitution path. *(Partially resolves 4/26 #8 -- `SubmissionAttempt` still lacks a constructor; see #11 below.)*

4. ~~Bootstrapping inconsistencies~~ (All) -- `NewApplication` across all services now accepts pre-built `*ports.Repository` and `*ports.Services` rather than constructing them internally. Bootstrap functions in `services/` and `persistence/` packages are the sole composition points. *(Not previously tracked; cleanup resolved same day.)*

5. ~~`InMemoryCacheManager` is not thread-safe~~ (Shared, P2) -- `InMemoryCacheManager` (`inmemory_cache_manager.go`) now uses a `sync.RWMutex`. `Get` acquires `RLock`, `Set` and `Del` acquire `Lock`. *(Introduced and resolved same day.)*

6. ~~`InMemoryCacheManager.Set` swallows marshal errors~~ (Shared, P2) -- `Set` (`inmemory_cache_manager.go:38`) now returns the `json.Marshal` error instead of `nil`. *(Introduced and resolved same day.)*

7. ~~`DecodeJSONResponse` does not check HTTP status code~~ (Shared, P2) -- `DecodeJSONResponse` (`httputil/http.go`) now returns an error for any status code >= 300 before attempting to decode the body. *(Introduced and resolved same day.)*

8. ~~`boostrapInMemory` typo~~ (Shared, P3) -- Renamed to `bootstrapInMemory` (`cache.go:42`). *(Introduced and resolved same day.)*

9. ~~`CreateVersionRequest` is an empty struct deserialized from request body~~ (Forms, P3) -- `CreateVersionRequest` merged into `UpsertVersionRequest` (formerly `UpdateVersionRequest`). The `createVersion` handler now parses pages from the request body via `dto.RequestToPages`, and `CreateVersionCommand` carries `Pages []*domain.Page`. `FormsService.CreateVersion` calls `version.SetPages(command.Pages...)` before persisting. *(Resolves 4/26 #2.)*

10. ~~Tenant removal does not cascade-delete DataSources~~ (Tenants, P2) -- `TenantsService.Delete` (`tenants_service.go:81-107`) now wraps tenant deletion and data source deletion in a transaction via `database.BeginTx`/`CommitTx`/`RollbackTx`. A new `DeleteAll(ctx, tenantID)` method was added to the `DataSourcesRepository` port with both in-memory and MongoDB implementations. The MongoDB implementation uses `DeleteMany` within a session. *(Resolves 4/26 #17.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the full Will Not Fix list (10 items).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:367`, `402`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/26 #1.)*

---

### Submissions Service

#### Bugs

2. **MongoDB repository methods are stubs returning `nil, nil`** (`submissions_repository.go:26-37`) -- `Find`, `FindByID`, and `FindByReferenceID` all return `nil, nil`. `FindByID` and `FindByReferenceID` returning nil will cause nil pointer dereference in the service layer when it accesses `submission.TenantID` (`submissions_service.go:40,58`). *(Unresolved from 4/26 #3.)*

3. **`NewSubmission` never sets `Status` or `Payload`** (`submission.go:33-50`) -- The constructor accepts `payload any` but never assigns it to the struct. `Status` is left as zero-value empty string. Both fields will persist as empty/nil despite being provided by the caller. *(New.)*

#### Architectural

4. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/26 #4.)*

5. **Four handler stubs return 200 OK with empty body** (`handlers.go:95-115`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission`. *(Unresolved from 4/26 #5.)*

#### Missing Functionality

6. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:63-71`) -- Return `nil, nil` and `nil`. *(Unresolved from 4/26 #6.)*

7. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/26 #7.)*

8. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindByID`, `FindByReferenceID`. *(Unresolved from 4/26 #9.)*

9. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/26 #10.)*

10. **`SubmissionAttempt` has no constructor or factory function** -- Bare struct with no `NewSubmissionAttempt`, no validation, no business methods. *(Partially unresolved from 4/26 #8.)*

#### Code Quality

11. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/26 #11.)*

12. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/26 #12.)*

13. **`FindByReferenceID` does a linear scan** in the in-memory repository. *(Unresolved from 4/26 #13.)*

14. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/26 #14.)*

15. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`). *(Unresolved from 4/26 #15.)*

---

### Tenants Service

#### Architectural

17. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/26 #16.)*

#### Missing Functionality

18. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/26 #19.)*

---

### Shared Package

#### Bugs

19. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/26 #24.)*

---

### Cross-Service

#### Architectural

20. **Zero test files** in all three services and shared packages. *(Unresolved from 4/26 #20.)*

21. **No domain events** for cross-service communication. *(Unresolved from 4/26 #21.)*

22. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/26 #22.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 2 | MongoDB repository methods are stubs -- nil panic on FindByID/FindByReferenceID | Submissions |
| **P1** | 3 | `NewSubmission` never sets `Status` or `Payload` | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 4 | `Find()` has no tenant filtering | Submissions |
| **P2** | 5 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 19 | `ErrMissingTenantID` maps to 500 | Shared |
| **P3** | 20 | Zero test files | All |
| **P3** | 21 | No domain events | All |
| **P3** | 10 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 11 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 12 | `SubmissionStatus` no constants | Submissions |
| **P3** | 18 | `Lookup` value object no validation | Tenants |
| **P3** | 22 | No real authentication | All |

---

## Summary

### Progress Since 4/26

The focus since the last review was infrastructure maturity, implementing the lookup feature end-to-end, and closing gaps in the forms and tenants services:

- **Graceful shutdown implemented** -- All three services now handle `SIGTERM`/`os.Interrupt` with a 15-second shutdown timeout via `http.Server.Shutdown`. The pattern is consistent across services: goroutine for `ListenAndServe`, `select` on server error vs signal channel, then graceful shutdown with context deadline. This closes a long-standing P3 issue tracked since 4/25.
- **Lookup strategies fully implemented** -- Three concrete `LookupStrategy` implementations now exist: `StaticLookupStrategy` returns data directly from attributes, `ScheduledLookupStrategy` returns pre-loaded data (with a TODO for lazy loading), and `WebhookLookupStrategy` makes an outbound HTTP call and decodes the JSON response into `[]*Lookup`. The strategy registry (`stratreg.StrategyRegistry`) dispatches by `DataSourceType`. The `DataSourcesService.Lookup` method is no longer a stub.
- **Cache package introduced** -- `pkg/cache` adds a `CacheManager` interface with `Get`/`Set`/`Del` and an in-memory implementation. Currently unused by any service but establishes the caching infrastructure. Thread-safety (`sync.RWMutex`) and error handling (marshal error propagation) are correctly implemented.
- **Submissions domain model matured** -- `NewSubmission` factory function added with UUID generation, timestamp, and struct validation. `HydrateSubmission` provides the persistence reconstitution path. However, the constructor has a bug: it never assigns `Status` or `Payload` to the struct (see #3).
- **Bootstrapping refactored** -- `NewApplication` across all services simplified to accept pre-built dependencies. `strategies.Bootstrap` wires lookup strategies with an `HTTPClient` interface for testability. The `strategy` package was renamed to `stratreg` to avoid collision with the `strategies` domain package.
- **Shared package hardened** -- `DecodeJSONResponse` now validates HTTP status codes (>= 300 returns error) before decoding, protecting the `WebhookLookupStrategy` from silently consuming error responses. `InMemoryCacheManager` now guards map access with `sync.RWMutex` and correctly propagates `json.Marshal` errors. The `boostrapInMemory` typo was fixed.
- **`createVersion` handler now functional** -- `CreateVersionRequest` (empty struct) merged into `UpsertVersionRequest`. The `createVersion` handler now deserializes pages from the request body, converts them via `dto.RequestToPages`, passes them through `CreateVersionCommand.Pages`, and `FormsService.CreateVersion` calls `version.SetPages(command.Pages...)` before upserting. This closes a P3 issue tracked since 4/26.
- **Tenant cascade-delete implemented** -- `TenantsService.Delete` now wraps the operation in a transaction (`database.BeginTx`/`CommitTx`/`RollbackTx`), deleting the tenant and then calling `DataSourcesRepository.DeleteAll(ctx, tenantID)` to remove all associated data sources. The new `DeleteAll` method is implemented in both in-memory (prefix-scanned map delete) and MongoDB (`DeleteMany` within a session) repositories. This closes a P2 issue tracked since 4/26.
- **Repository naming consistency** -- `mongoDBVersionRepository` renamed to `mongoDBVersionsRepository` (file renamed `version_repository.go` -> `versions_repository.go`). Similarly `inMemoryTenantRepository` renamed to `inMemoryTenantsRepository` (file renamed `tenant_repository.go` -> `tenants_repository.go`). Aligns with the plural convention used by the port interfaces.

### Current State

**24 remaining issues (4/26) -> 22 remaining issues** (resolved 10 total: 6 from prior reviews, 4 introduced-and-resolved same day; introduced 1 new P1 bug).

**Forms Service** now has a functional `createVersion` endpoint that accepts pages in the request body, closing the empty-struct deserialization issue. The `UpdateVersionRequest` was consolidated into `UpsertVersionRequest` shared by both create and update handlers. Remaining gap: hardcoded placeholder user ID (P2).

**Tenants Service** now has a complete lookup feature and proper cascade-delete on tenant removal. The `Delete` method uses transactional semantics to ensure atomicity. The `DataSourcesRepository.DeleteAll` is cleanly implemented in both persistence backends. Remaining gaps: `Lookup` value object no validation (P3), no pagination on `Find()` (P3).

**Submissions Service** has improved with domain constructors but introduced a P1 bug: `NewSubmission` accepts `payload` and uses a `Status` field but assigns neither. The MongoDB repository methods remain stubs (P0). The domain is still largely anemic.

**Hexagonal Architecture** -- Dependency direction remains correct. The new `ports.HTTPClient` interface and `ports.LookupStrategy` interface maintain the port boundary. Strategy implementations live in `core/strategies/`, which is appropriate since they contain domain logic. The `stratreg` package in `pkg/common` is a clean generic utility with no domain coupling. The `TenantsService` now correctly depends on both `TenantsRepository` and `DataSourcesRepository` for cascade operations, wired through the `Repository` aggregate.

**CQRS-Lite** -- The new `GetDataSourceLookupsQuery` follows the established query pattern with constructor and validation. `CreateVersionCommand` now carries domain data (`Pages`), which is consistent with the command pattern. No violations.

**DDD** -- The `DataSourceAttributes` sealed interface hierarchy continues to be one of the strongest DDD patterns in the codebase. The strategy pattern for lookups aligns well with DDD's domain service concept. `Lookup` value object still lacks validation, which weakens the invariant guarantees. The transactional cascade-delete in `TenantsService` correctly enforces the aggregate boundary between tenants and data sources.

**Idiomatic Go** -- The `getDataSourceAttributes[T]` generic helper in `strategies.go` is clean and avoids repetitive type switches. The `StrategyRegistry` generic type is well-designed with a fluent builder API (`New().Set().Set()`). The `InMemoryCacheManager` now follows Go concurrency best practices with proper `sync.RWMutex` usage. The `UpsertVersionRequest` consolidation reduces DTO duplication.

### Highest-Impact Improvements

1. **Fix `NewSubmission` to assign `Status` and `Payload`** (P1 -- constructor produces incomplete domain objects)
2. **Implement submissions MongoDB repository methods** (P0 -- stub methods cause nil panic in service layer)
3. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
4. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
5. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
