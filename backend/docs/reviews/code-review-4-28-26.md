# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/26 Review

1. ~~No graceful shutdown~~ (All, P3) -- All three services (`forms/cmd/server/main.go`, `submissions/cmd/server/main.go`, `tenants/cmd/server/main.go`) now listen for `os.Interrupt`/`SIGTERM` and call `server.Shutdown(ctx)` with a 15-second timeout. *(Resolves 4/26 #23.)*

2. ~~`Lookup` service method is a stub~~ (Tenants, P3) -- `DataSourcesService.Lookup` (`data_sources_service.go:125-139`) now resolves a `LookupStrategy` from the strategy registry by data source type and delegates to `strategy.Lookup(ctx, ds)`. Three strategy implementations exist: `StaticLookupStrategy`, `ScheduledLookupStrategy`, and `WebhookLookupStrategy`. *(Resolves 4/26 #18.)*

3. ~~No domain constructors~~ (Submissions, P2) -- `NewSubmission` (`submission.go:33-50`) now generates IDs, sets `CreatedAt`, initializes `Attempts`, and validates via `validate.ValidateStruct`. `HydrateSubmission` (`submission.go:52-73`) provides a persistence reconstitution path. *(Partially resolves 4/26 #8 -- `SubmissionAttempt` still lacks a constructor; see #10 below.)*

4. ~~Bootstrapping inconsistencies~~ (All) -- `NewApplication` across all services now accepts pre-built `*ports.Repository` and `*ports.Services` rather than constructing them internally. Bootstrap functions in `services/` and `persistence/` packages are the sole composition points. *(Not previously tracked; cleanup resolved same day.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the full Will Not Fix list (10 items).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:367`, `402`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/26 #1.)*

#### Code Quality

2. **`CreateVersionRequest` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:280-283`) -- `createVersion` calls `ReadValidateJSONPayload(r, &body)` where `body` is `CreateVersionRequest struct{}`. *(Unresolved from 4/26 #2.)*

---

### Submissions Service

#### Bugs

3. **MongoDB repository methods are stubs returning `nil, nil`** (`submissions_repository.go:26-37`) -- `Find`, `FindByID`, and `FindByReferenceID` all return `nil, nil`. `FindByID` and `FindByReferenceID` returning nil will cause nil pointer dereference in the service layer when it accesses `submission.TenantID` (`submissions_service.go:40,58`). *(Unresolved from 4/26 #3.)*

4. **`NewSubmission` never sets `Status` or `Payload`** (`submission.go:33-50`) -- The constructor accepts `payload any` but never assigns it to the struct. `Status` is left as zero-value empty string. Both fields will persist as empty/nil despite being provided by the caller. *(New.)*

#### Architectural

5. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/26 #4.)*

6. **Four handler stubs return 200 OK with empty body** (`handlers.go:95-115`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission`. *(Unresolved from 4/26 #5.)*

#### Missing Functionality

7. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:63-71`) -- Return `nil, nil` and `nil`. *(Unresolved from 4/26 #6.)*

8. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/26 #7.)*

9. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindByID`, `FindByReferenceID`. *(Unresolved from 4/26 #9.)*

10. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/26 #10.)*

11. **`SubmissionAttempt` has no constructor or factory function** -- Bare struct with no `NewSubmissionAttempt`, no validation, no business methods. *(Partially unresolved from 4/26 #8.)*

#### Code Quality

12. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/26 #11.)*

13. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/26 #12.)*

14. **`FindByReferenceID` does a linear scan** in the in-memory repository. *(Unresolved from 4/26 #13.)*

15. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/26 #14.)*

16. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`). *(Unresolved from 4/26 #15.)*

---

### Tenants Service

#### Architectural

17. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/26 #16.)*

18. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:75-87`). *(Unresolved from 4/26 #17.)*

#### Missing Functionality

19. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/26 #19.)*

---

### Shared Package

#### Bugs

20. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/26 #24.)*

21. **`InMemoryCacheManager` is not thread-safe** (`inmemory_cache_manager.go`) -- The `cache` map is accessed without a `sync.RWMutex`. Concurrent `Get`/`Set`/`Del` calls will cause a data race. *(New.)*

22. **`InMemoryCacheManager.Set` swallows marshal errors** (`inmemory_cache_manager.go:34`) -- When `json.Marshal(data)` fails, the method returns `nil` instead of the error. Callers will believe the write succeeded. *(New.)*

23. **`DecodeJSONResponse` does not check HTTP status code** (`httputil/http.go`) -- `DecodeJSONResponse` decodes the body regardless of status. A 4xx/5xx response will either decode into an unexpected shape or silently produce a zero-value result. The `WebhookLookupStrategy` calls this on the external HTTP response without status validation. *(New.)*

24. **`boostrapInMemory` typo** (`cache.go:39`) -- Function name should be `bootstrapInMemory`. *(New.)*

---

### Cross-Service

#### Architectural

25. **Zero test files** in all three services and shared packages. *(Unresolved from 4/26 #20.)*

26. **No domain events** for cross-service communication. *(Unresolved from 4/26 #21.)*

27. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/26 #22.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 3 | MongoDB repository methods are stubs -- nil panic on FindByID/FindByReferenceID | Submissions |
| **P1** | 4 | `NewSubmission` never sets `Status` or `Payload` | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 5 | `Find()` has no tenant filtering | Submissions |
| **P2** | 6 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 18 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 20 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 21 | `InMemoryCacheManager` not thread-safe | Shared |
| **P2** | 22 | `InMemoryCacheManager.Set` swallows marshal errors | Shared |
| **P2** | 23 | `DecodeJSONResponse` no status check | Shared |
| **P3** | 25 | Zero test files | All |
| **P3** | 26 | No domain events | All |
| **P3** | 11 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 12 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 13 | `SubmissionStatus` no constants | Submissions |
| **P3** | 19 | `Lookup` value object no validation | Tenants |
| **P3** | 2 | `CreateVersionRequest` empty struct deserialized | Forms |
| **P3** | 24 | `boostrapInMemory` typo | Shared |
| **P3** | 27 | No real authentication | All |

---

## Summary

### Progress Since 4/26

The focus since the last review was infrastructure maturity and implementing the lookup feature end-to-end:

- **Graceful shutdown implemented** -- All three services now handle `SIGTERM`/`os.Interrupt` with a 15-second shutdown timeout via `http.Server.Shutdown`. The pattern is consistent across services: goroutine for `ListenAndServe`, `select` on server error vs signal channel, then graceful shutdown with context deadline. This closes a long-standing P3 issue tracked since 4/25.
- **Lookup strategies fully implemented** -- Three concrete `LookupStrategy` implementations now exist: `StaticLookupStrategy` returns data directly from attributes, `ScheduledLookupStrategy` returns pre-loaded data (with a TODO for lazy loading), and `WebhookLookupStrategy` makes an outbound HTTP call and decodes the JSON response into `[]*Lookup`. The strategy registry (`stratreg.StrategyRegistry`) dispatches by `DataSourceType`. The `DataSourcesService.Lookup` method is no longer a stub.
- **Cache package introduced** -- `pkg/cache` adds a `CacheManager` interface with `Get`/`Set`/`Del` and an in-memory implementation. Currently unused by any service but establishes the caching infrastructure. Has thread-safety and error-handling issues (see #21, #22).
- **Submissions domain model matured** -- `NewSubmission` factory function added with UUID generation, timestamp, and struct validation. `HydrateSubmission` provides the persistence reconstitution path. However, the constructor has a bug: it never assigns `Status` or `Payload` to the struct (see #4).
- **Bootstrapping refactored** -- `NewApplication` across all services simplified to accept pre-built dependencies. `strategies.Bootstrap` wires lookup strategies with an `HTTPClient` interface for testability. The `strategy` package was renamed to `stratreg` to avoid collision with the `strategies` domain package.

### Current State

**24 remaining issues (4/26) -> 27 remaining issues** (resolved 4 from prior reviews; introduced 5 new issues including 1 P1 bug).

**Forms Service** remains stable. No changes since 4/26. Remaining gaps: hardcoded placeholder user ID (P2), `CreateVersionRequest` empty struct (P3).

**Tenants Service** now has a complete lookup feature. The strategy pattern is well-implemented: `LookupStrategy` interface defined as a port, concrete implementations in the `strategies` package, and the registry wired via `Bootstrap`. The `WebhookLookupStrategy` correctly uses `ports.HTTPClient` (an interface wrapping `*http.Client`) for testability. The sealed `DataSourceAttributes` interface with the `getDataSourceAttributes[T]` generic helper in the strategies package is a clean approach to type-safe attribute dispatch. Remaining gaps: no cascade-delete on tenant removal (P2), `Lookup` value object no validation (P3), no pagination on `Find()` (P3).

**Submissions Service** has improved with domain constructors but introduced a P1 bug: `NewSubmission` accepts `payload` and uses a `Status` field but assigns neither. The MongoDB repository methods remain stubs (P0). The domain is still largely anemic.

**Hexagonal Architecture** -- Dependency direction remains correct. The new `ports.HTTPClient` interface and `ports.LookupStrategy` interface maintain the port boundary. Strategy implementations live in `core/strategies/`, which is appropriate since they contain domain logic. The `stratreg` package in `pkg/common` is a clean generic utility with no domain coupling.

**CQRS-Lite** -- The new `GetDataSourceLookupsQuery` follows the established query pattern with constructor and validation. No violations.

**DDD** -- The `DataSourceAttributes` sealed interface hierarchy continues to be one of the strongest DDD patterns in the codebase. The strategy pattern for lookups aligns well with DDD's domain service concept. `Lookup` value object still lacks validation, which weakens the invariant guarantees.

**Idiomatic Go** -- The `getDataSourceAttributes[T]` generic helper in `strategies.go` is clean and avoids repetitive type switches. The `StrategyRegistry` generic type is well-designed with a fluent builder API (`New().Set().Set()`). The `boostrapInMemory` typo in `cache.go` should be fixed. The `InMemoryCacheManager` violates Go concurrency best practices by accessing a map without synchronization.

### Highest-Impact Improvements

1. **Fix `NewSubmission` to assign `Status` and `Payload`** (P1 -- constructor produces incomplete domain objects)
2. **Implement submissions MongoDB repository methods** (P0 -- stub methods cause nil panic in service layer)
3. **Add `sync.RWMutex` to `InMemoryCacheManager`** (P2 -- data race on concurrent access)
4. **Fix `InMemoryCacheManager.Set` to return marshal errors** (P2 -- silent data loss)
5. **Add HTTP status check in `DecodeJSONResponse`** (P2 -- webhook strategy will silently produce empty results on error responses)
