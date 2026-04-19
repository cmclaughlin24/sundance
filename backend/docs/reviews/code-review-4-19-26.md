# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/18 Review

1. ~~`Database` field is nil in in-memory bootstrap~~ (Submissions #22, P0) -- `submissions/internal/adapters/persistence/inmemory/inmemory.go` now sets `Database: database.NewInMemoryDatabase()`. `Application.Close()` will no longer panic.

2. ~~`core.go` imports from the adapters layer~~ (Forms #2, Submissions #23, Tenants #41, P2) -- All three `core.go` files no longer import from `adapters/persistence`. `NewApplication` now accepts `*log.Logger` and `*ports.Repository` as parameters. Persistence bootstrapping has been moved to `cmd/server/main.go` in each service.

3. ~~Services receive the entire `*ports.Repository` struct~~ (Forms #3, Tenants #42, P2) -- Each service now stores only the specific repository interface(s) it needs:
   - `FormsService` -> `database.Database` + `ports.FormsRepository`
   - `TenantsService` -> `ports.TenantsRepository`
   - `DataSourcesService` -> `ports.DataSourcesRepository`
   - `SubmissionsService` -> `ports.SubmissionsRepository`

4. ~~`DataSourceAttributes` typed as `any`~~ (Tenants #57, P3) -- `data_source_attributes.go` now defines a sealed interface with a `isDataSourceAttributes()` marker method. Each concrete type (`StaticDataSourceAttributes`, `ScheduledDataSourceAttributes`, `QueryDataSourceAttributes`) implements the marker.

---

## Will Not Fix

These issues have been reviewed and accepted as intentional design decisions. They should not be flagged in future reviews.

1. **Goroutine/channel pattern in handlers** (Previously Forms #7, Submissions #27, Tenants #46) -- The `go func() -> chan -> select { case <-r.Context().Done(); case res := <-resultChan }` pattern in every handler is the recommended approach for respecting chi router's context-based request timeouts. Without this pattern, a handler performing a long-running service call would not be able to short-circuit when the request context is cancelled (e.g., client disconnect, server timeout). The `select` on `r.Context().Done()` enables cooperative cancellation at the handler level. The allocation overhead of a single goroutine and buffered channel per request is negligible relative to the I/O cost of a real database call.

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:328-330`, `363-365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1.)*

#### Architectural

2. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-19`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services` or the specific service interface. *(Unresolved from 4/13 #9, 4/17 #4, 4/18 #4.)*

3. **`Find()` has no tenant filtering** (`form_service.go:26-28`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5.)*

4. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6.)*

5. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `form_service.go` in `PublishVersion` and `RetireVersion`) -- Couples persistence and business logic to wall-clock time. Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8.)*

6. **`internal/types/error.go` breaks domain purity** (`internal/types/error.go`) -- `ErrDuplicatePosition` is defined outside the domain package and imported by `domain/page.go`, `domain/section.go`, and `domain/version.go`. This creates an outward dependency from the domain layer. Domain errors should live in `core/domain/`. *(Unresolved from 4/18 #9.)*

#### Missing Functionality

7. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10.)*

8. **Domain validation unimplemented** (`form.go`, `field.go`, `page.go`, `section.go`, `version.go`) -- All entity constructors contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13, 4/17 #10, 4/18 #11.)*

9. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12.)*

10. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and domain errors (`ErrVersionLocked`, `ErrInvalidVersion`, etc.) all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` is a pass-through with only a `default` case. *(Unresolved from 4/13 #15, 4/17 #12, 4/18 #13.)*

11. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14.)*

12. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15.)*

13. **`ConditionalRule` is an empty stub** (`conditional_rule.go`) -- Contains only an ID; no rule type, conditions, or actions. *(Unresolved from 4/13 #18, 4/17 #15, 4/18 #16.)*

#### Code Quality

14. **`FieldAttributes` is typed as `any`** (`field_attributes.go:3`) -- No compile-time type safety. Should use a sealed interface with a marker method (e.g., `isFieldAttributes()`), consistent with the tenants service which now uses this pattern for `DataSourceAttributes`. *(Unresolved from 4/13 #19, 4/17 #16, 4/18 #17.)*

15. **`FieldResponse` DTO omits `Attributes`** (`dto/field.go`) -- `FieldToResponse` maps all fields except `Attributes`, meaning field attribute data is silently lost in API responses. *(Unresolved from 4/17 #17, 4/18 #18.)*

16. **Inconsistent constructor signatures** -- Forms domain constructors return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. *(Unresolved from 4/13 #23, 4/17 #18, 4/18 #19.)*

17. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing; the `CreateVersion` flow has a race condition. *(Unresolved from 4/13 #29, 4/17 #19, 4/18 #20.)*

18. **Inconsistent attribute parsing pattern** -- Forms uses a plain `map[domain.FieldType]attributeParser` while tenants uses the generic `strategy.Strategies[K, U]` utility for the same problem. *(Unresolved from 4/18 #21.)*

19. **`ErrMissingTenantID` maps to 500** (`middleware.go:24`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` is a local error (`errors.New("X-Tenant-ID header is required")`), it doesn't match any case in `SendErrorResponse` and falls through to the 500 default. Should be 400. *(New.)*

---

### Submissions Service

#### Architectural

20. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-20`) -- Same pattern as other services. *(Unresolved from 4/17 #28, 4/18 #26.)*

21. **No tenant middleware** -- No `middleware.go` implementation exists (file is empty). The `getSubmissionByReferenceID` handler passes an empty string `""` as the tenantID to `NewFindByIdQuery`. *(Unresolved from 4/17 #26, 4/18 #24.)*

22. **`Find()` has no tenant filtering** (`submissions_service.go:24-26`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25.)*

23. **`NewFindByIdQuery` creates `validator.New()` per call** (`ports/query.go:16`) -- Inconsistent with tenants and forms which use the shared `validate.ValidateStruct()` singleton from `pkg/common/validate`. *(Unresolved from 4/18 #28.)*

24. **`sendErrorResponse` method is dead code** (`handlers.go:97-101`) -- The `sendErrorResponse` method exists but is never called. All handlers call `httputil.SendErrorResponse` directly (lines 42, 61, 76). *(New.)*

#### Missing Functionality

25. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #30, 4/18 #29.)*

26. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:56-62`) -- Return `nil, nil` and `nil` respectively. *(Unresolved from 4/17 #31, 4/18 #30.)*

27. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31.)*

28. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32.)*

29. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33.)*

30. **No domain events** for cross-service communication. *(Unresolved from 4/17 #35, 4/18 #34.)*

31. **No real authentication.** *(Unresolved from 4/17 #36, 4/18 #35.)*

32. **`ReplaySubmissionCommand` is an empty struct** (`commands.go:3`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36.)*

#### Code Quality

33. **`Payload` typed as `any`** (`submission.go:18`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37.)*

34. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38.)*

35. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go:51-61`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39.)*

---

### Tenants Service

#### Bugs

36. **`ErrDataSourceAttrParse` not wrapped on JSON unmarshal failure** (`dto/request.go:61-68`) -- `parseAttributes` returns raw `json.Unmarshal` errors without wrapping them in `ErrDataSourceAttrParse`. Only the unknown-type path (strategy lookup failure) wraps the sentinel. Malformed attribute JSON bypasses error mapping and produces 500. *(Unresolved from 4/17 #41, 4/18 #40.)*

#### Architectural

37. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-19`) -- Same pattern as other services. *(Unresolved from 4/16 #7, 4/17 #45, 4/18 #43.)*

38. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:37-62`) -- `Create` validates the command and calls `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources. *(Unresolved from 4/16 #8, 4/17 #46, 4/18 #44.)*

39. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45.)*

40. **Tenants service has no tenant-scoping middleware** -- Unlike forms which extracts `X-Tenant-ID` via middleware into context, tenants relies on URL path params with no validation that the tenant exists. Inconsistent multi-tenancy approach across services. *(Unresolved from 4/18 #47.)*

#### Missing Functionality

41. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12, 4/17 #50, 4/18 #48.)*

42. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` has some validation (attribute type matching) but does not validate empty `TenantID`, empty `Type`, or field lengths. *(Unresolved from 4/16 #13, 4/17 #51, 4/18 #49.)*

43. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14, 4/17 #52, 4/18 #50.)*

44. **Incomplete error-to-HTTP mapping** -- `ErrInvalidID` and `ErrUnauthorized` fall through to 500. Service-specific errors like `ErrDataSourceAttrParse` and `ErrStrategyNotFound` are also unmapped. The service-level `sendErrorResponse` only maps `ErrInvalidSourceTypeAttributes`. *(Unresolved from 4/16 #15, 4/17 #53, 4/18 #51.)*

45. **No real authentication**. *(Unresolved from 4/16 #16, 4/17 #54, 4/18 #52.)*

46. **`Lookup` service method is a stub** (`data_sources_service.go:105-114`) -- Always returns `nil, nil` after verifying the data source exists. Contains `// TODO`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53.)*

47. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names. *(Unresolved from 4/16 #18, 4/17 #56, 4/18 #54.)*

48. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct with two fields, no json tags, no invariant enforcement. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55.)*

49. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string for `Type`. *(Unresolved from 4/16 #22, 4/17 #58, 4/18 #56.)*

#### Code Quality

50. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:69`). *(Unresolved from 4/16 #24, 4/17 #60, 4/18 #58.)*

51. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. *(Unresolved from 4/16 #29, 4/17 #61, 4/18 #59.)*

52. **Error wrapping format inconsistency** (`dto/request.go:49`) -- Uses `fmt.Errorf("%w, %w", ErrDataSourceAttrParse, err)` (comma separator) while the shared `ReadJsonPayload` uses `fmt.Errorf("%w: %w", ErrDecodeJSON, err)` (colon separator). Minor but inconsistent. *(New.)*

---

### Shared Package (`pkg/common`)

#### Bugs

53. **`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`** (`httputil/http.go:56-87`) -- These sentinel errors are defined in `error.go` but not handled in the switch. They fall through to the 500 default. `ErrInvalidID` should map to 400; `ErrUnauthorized` should map to 401 or 403. *(Unresolved from 4/17 #64, 4/18 #60.)*

54. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`httputil/http.go:41`) -- The `headers ...http.Header` variadic parameter is accepted but the body never iterates or sets them on the response. *(Unresolved from 4/17 #65, 4/18 #61.)*

55. **`w.Write` error ignored** (`httputil/http.go:50`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. *(Unresolved from 4/17 #66, 4/18 #62.)*

#### Code Quality

56. **`ValidateStruct` has a redundant pattern** (`validate/validate.go:19-25`) -- `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. *(Unresolved from 4/17 #67, 4/18 #63.)*

57. **In-memory transactions are no-ops** (`database/inmemory_database.go`) -- `BeginTx` stores an empty struct in context; `CommitTx`/`RollbackTx` return `nil` in all code paths. No atomicity guarantees. *(Unresolved from 4/17 #68, 4/18 #64.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 53 | `ErrInvalidID`/`ErrUnauthorized` map to 500 | Shared |
| **P1** | 36 | `ErrDataSourceAttrParse` not wrapped on unmarshal | Tenants |
| **P2** | 2, 20, 37 | Handlers receive full Application | All |
| **P2** | 3, 22 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 8, 42 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 28 | No domain constructors | Submissions |
| **P2** | 38 | DataSource created without parent Tenant check | Tenants |
| **P2** | 21 | No tenant middleware | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 6 | `internal/types/error.go` breaks domain purity | Forms |
| **P2** | 40 | Inconsistent multi-tenancy approach | Tenants |
| **P2** | 19 | `ErrMissingTenantID` maps to 500 | Forms |
| **P3** | 7, 25, 41 | Zero test files | All |
| **P3** | 9, 30, 43 | No domain events | All |
| **P3** | 10, 44 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 14, 33 | `any`-typed attributes (no type safety) | Forms, Submissions |
| **P3** | 5, 50 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 51 | Inconsistent response envelope | Tenants |
| **P3** | 17, 57 | In-memory transactions are no-ops | Forms, Shared |
| **P3** | 15 | `FieldResponse` DTO omits Attributes | Forms |
| **P3** | 56 | Redundant validation patterns | Shared |
| **P3** | 18 | Inconsistent attribute parsing pattern | Forms, Tenants |
| **P3** | 23 | `NewFindByIdQuery` creates validator per call | Submissions |
| **P3** | 24 | `sendErrorResponse` dead code | Submissions |
| **P3** | 52 | Error wrapping format inconsistency | Tenants |

---

## Summary

### Progress Since 4/18

Four issues from the prior review have been resolved:

- **Submissions P0 crash resolved** -- `Database` is now set in the in-memory bootstrap, so `Application.Close()` no longer panics with a nil pointer dereference.
- **Hexagonal dependency violation fixed across all three services** -- `core.go` no longer imports from the adapters layer. `ApplicationSettings` with its `persistence.PersistenceSettings` reference has been removed. `NewApplication` now accepts a `*log.Logger` and `*ports.Repository`, with persistence bootstrapping moved to `cmd/server/main.go`. This is a significant structural improvement.
- **Interface Segregation violation fixed** -- All service structs now store only the specific repository interfaces they need (`FormsRepository`, `TenantsRepository`, `DataSourcesRepository`, `SubmissionsRepository`) rather than the full `*ports.Repository` struct.
- **`DataSourceAttributes` sealed interface** (Tenants) -- No longer typed as `any`. Now uses a marker method pattern (`isDataSourceAttributes()`) for compile-time type safety.

### New Issues Found

1. **`ErrMissingTenantID` maps to 500** (Forms #19) -- The forms tenant middleware sends this error through `SendErrorResponse`, where it hits the 500 default. A missing header should produce 400.
2. **`sendErrorResponse` is dead code** (Submissions #24) -- The method exists but handlers call `httputil.SendErrorResponse` directly.
3. **Error wrapping format inconsistency** (Tenants #52) -- `dto/request.go` uses comma separator while the shared package uses colon separator in `fmt.Errorf` wrapping.

### Issues Moved to Will Not Fix

The goroutine/channel pattern in handlers (previously Forms #7, Submissions #27, Tenants #46) has been accepted as a deliberate design decision. It is the recommended pattern for cooperative cancellation with chi router's context-based timeouts.

### Current State

**Forms Service** remains the most mature. The P2 architectural violations (core importing adapters, ISP) are now resolved. The primary remaining gaps are: the hardcoded `"placeholder"` user IDs, the aggregate boundary ambiguity between `Form` and `Version`, the `FieldResponse` DTO silently dropping attributes, the domain error defined outside the domain package (`ErrDuplicatePosition`), the `ErrMissingTenantID`-to-500 mapping bug, and the continued absence of domain validation and test coverage. The handlers still receive the full `Application` struct.

**Tenants Service** has made meaningful progress: the hexagonal violation is fixed, services now follow ISP, and `DataSourceAttributes` uses a proper sealed interface. The `DataSource` orphan problem (no tenant existence check), incomplete attribute types (missing json tags, empty `ScheduledDataSourceAttributes`), and the `parseAttributes` error wrapping gap remain. The `sendErrorResponse` now correctly maps `ErrInvalidSourceTypeAttributes` to 400.

**Submissions Service** has resolved its sole P0 (nil `Database`). Beyond that, the service still lacks tenant middleware, write operations, domain constructors, request DTOs, and test coverage. The `sendErrorResponse` method is defined but never used.

**Shared Package** (`pkg/common`) is unchanged. Missing `ErrInvalidID`/`ErrUnauthorized` HTTP mappings, the unused `headers` parameter in `SendJsonResponse`, the ignored `w.Write` error, redundant validation patterns, and the no-op in-memory transactions remain.

### Highest-Impact Improvements

1. **Add `ErrInvalidID` and `ErrUnauthorized` mappings** to `SendErrorResponse` (P1 -- auth/validation errors produce 500s)
2. **Wrap `parseAttributes` errors with `ErrDataSourceAttrParse`** in `dto/request.go` (P1 -- malformed attribute JSON produces 500)
3. **Fix `ErrMissingTenantID` mapping** in forms middleware or `SendErrorResponse` (P2 -- missing tenant header produces 500)
4. **Narrow handler dependencies** to `*ports.Services` instead of `*core.Application` across all three services (P2 -- principle of least privilege)
5. **Move `ErrDuplicatePosition` into `core/domain/`** in the forms service (P2 -- domain purity)
6. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
