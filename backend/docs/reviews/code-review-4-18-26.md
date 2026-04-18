# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/17 Review

1. ~~`Tenant.DataSources` field declared but never populated~~ (Tenants #47) -- The `DataSources` field has been removed from the `Tenant` struct entirely (`tenant.go`). This also resolves the related issue of `getTenants` returning tenants with `DataSources` always `nil` (Tenants #62).

2. ~~`Remove` operations succeed silently for non-existent entities~~ (Tenants #42) -- While the repository layer still uses Go's `delete` (a no-op for missing keys), the service layer now guards both `Remove` methods with an `Exists` check before calling the repository, returning `common.ErrNotFound` for missing entities (`tenants_service.go:73-84`, `data_sources_service.go:91-102`).

3. ~~Stray semicolon in `data_sources_service.go`~~ (Tenants #63) -- Removed.

4. ~~`DataSourceResponse` DTO silently drops `Name` and `Description`~~ (New from this review) -- `DataSourceResponse` now includes `Name` and `Description` fields, and `DataSourceToResponse` maps them from the domain model (`dto/response.go`). *(Unstaged.)*

5. ~~`nil` route handlers will panic at runtime~~ (Submissions #20) -- All 6 routes now wire to actual handler methods on the `handlers` struct (`routes.go`).

6. ~~Services bootstrap does not wire `SubmissionsService`~~ (Submissions #22) -- `services.Bootstrap()` now returns `&ports.Services{Submissions: NewSubmissionsService(logger, repository)}` (`services/services.go`).

7. ~~Entry point is a stub~~ (Submissions #23) -- `main.go` now creates the `Application`, builds routes via `rest.NewRoutes(app)`, configures an `http.Server` with timeouts, and starts listening on port 80.

8. ~~REST uses `net/http.ServeMux` while other services use `go-chi/chi`~~ (Submissions #25) -- `routes.go` now uses `chi.NewRouter()`, consistent with forms and tenants.

9. ~~DTOs not implemented~~ (Submissions #32, partially) -- `dto/response.go` now defines `SubmissionResponse` and `SubmissionToResponse` mapper. However, `dto/request.go` remains empty (contains only the package declaration).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:328-330`, `363-365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1.)*

#### Architectural

2. **`core.go` imports from the adapters layer** (`core.go:7`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, and `NewApplication` calls `persistence.Bootstrap` directly. The core package should not depend on adapter types. *(Unresolved from 4/13 #7, 4/17 #2.)*

3. **Services receive the entire `Repository` struct** (`form_service.go:14-17`) -- The form service holds `*ports.Repository` which includes `Database` and the full repository aggregate rather than just the `FormsRepository` interface. Violates Interface Segregation. *(Unresolved from 4/13 #8, 4/17 #3.)*

4. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-19`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services` or the specific service interface. *(Unresolved from 4/13 #9, 4/17 #4.)*

5. **`Find()` has no tenant filtering** (`form_service.go:26-28`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. *(Unresolved from 4/13 #10, 4/17 #5.)*

6. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6.)*

7. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 10 handlers) -- `net/http` already runs each handler in its own goroutine. The extra goroutine + buffered channel + `select` pattern adds allocation overhead and complexity with no concurrency benefit. *(Unresolved from 4/13 #25, 4/17 #7.)*

8. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `form_service.go` in `PublishVersion` and `RetireVersion`) -- Couples persistence and business logic to wall-clock time. Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26, 4/17 #8.)*

9. **`internal/types/error.go` breaks domain purity** (`internal/types/error.go`) -- `ErrDuplicatePosition` is defined outside the domain package and imported by `domain/page.go`, `domain/section.go`, and `domain/version.go`. This creates an outward dependency from the domain layer. Domain errors should live in `core/domain/`. The tenants service correctly places domain errors inside the domain package. *(New.)*

#### Missing Functionality

10. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12, 4/17 #9.)*

11. **Domain validation unimplemented** (`form.go`, `field.go`, `page.go`, `section.go`, `version.go`) -- All entity constructors contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13, 4/17 #10.)*

12. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11.)*

13. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and domain errors (`ErrVersionLocked`, `ErrInvalidVersion`, etc.) all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` is a pass-through with only a `default` case. *(Unresolved from 4/13 #15, 4/17 #12.)*

14. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16, 4/17 #13.)*

15. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. Routes register only GET, POST, and PUT. *(Unresolved from 4/13 #17, 4/17 #14.)*

16. **`ConditionalRule` is an empty stub** (`conditional_rule.go`) -- Contains only an ID; no rule type, conditions, or actions. *(Unresolved from 4/13 #18, 4/17 #15.)*

#### Code Quality

17. **`FieldAttributes` is typed as `any`** (`field_attributes.go:3`) -- No compile-time type safety. Should use a sealed interface with a marker method (e.g., `isFieldAttributes()`). *(Unresolved from 4/13 #19, 4/17 #16.)*

18. **`FieldResponse` DTO omits `Attributes`** (`dto/field.go`) -- `FieldToResponse` maps all fields except `Attributes`, meaning field attribute data is silently lost in API responses. *(Unresolved from 4/17 #17.)*

19. **Inconsistent constructor signatures** -- Forms domain constructors return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. *(Unresolved from 4/13 #23, 4/17 #18.)*

20. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing; the `CreateVersion` flow has a race condition. *(Unresolved from 4/13 #29, 4/17 #19.)*

21. **Inconsistent attribute parsing pattern** -- Forms uses a plain `map[domain.FieldType]attributeParser` while tenants uses the generic `strategy.Strategies[K, U]` utility for the same problem. *(New.)*

---

### Submissions Service

#### Bugs

22. **`Database` field is nil in in-memory bootstrap** (`inmemory/inmemory.go`) -- `Bootstrap` returns `&ports.Repository{Submissions: ...}` without setting the `Database` field. `Application.Close()` calls `app.repository.Database.Close()`, which will panic with a nil pointer dereference. *(Unresolved from 4/17 #21.)*

#### Architectural

23. **`core.go` imports from the adapters layer** (`core.go:7`) -- Same hexagonal dependency violation as the other services. *(Unresolved from 4/17 #24.)*

24. **No tenant middleware** -- No `middleware.go` file exists. The `getSubmissionByReferenceID` handler passes an empty string `""` as the tenantID to `NewFindByIdQuery`. *(Unresolved from 4/17 #26.)*

25. **`Find()` has no tenant filtering** (`submissions_service.go:24-26`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27.)*

26. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-20`) -- Same pattern as other services. *(Unresolved from 4/17 #28.)*

27. **Unnecessary goroutine/channel pattern** in implemented handlers (`getSubmissions`, `getSubmissionByReferenceID`). *(Unresolved from 4/17 #29.)*

28. **`NewFindByIdQuery` creates `validator.New()` per call** (`ports/query.go:16`) -- Inconsistent with tenants and forms which use the shared `validate.ValidateStruct()` singleton from `pkg/common/validate`. *(New.)*

#### Missing Functionality

29. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #30.)*

30. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:56-62`) -- Return `nil, nil` and `nil` respectively. *(Unresolved from 4/17 #31.)*

31. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Partially unresolved from 4/17 #32.)*

32. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33.)*

33. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34.)*

34. **No domain events** for cross-service communication. *(Unresolved from 4/17 #35.)*

35. **No real authentication.** *(Unresolved from 4/17 #36.)*

36. **`ReplaySubmissionCommand` is an empty struct** (`commands.go:3`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/17 #37.)*

#### Code Quality

37. **`Payload` typed as `any`** (`submission.go:18`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38.)*

38. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39.)*

39. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go:51-61`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40.)*

---

### Tenants Service

#### Bugs

40. **`ErrDataSourceAttrParse` not wrapped on JSON unmarshal failure** (`dto/request.go:58-65`) -- `parseAttributes` returns raw `json.Unmarshal` errors without wrapping them in `ErrDataSourceAttrParse`. Only the unknown-type path (strategy lookup failure) wraps the sentinel. Malformed attribute JSON bypasses error mapping and produces 500. *(Unresolved from 4/17 #41.)*

#### Architectural

41. **`core.go` imports from the adapters layer** (`core.go:7`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, violating the hexagonal dependency rule. *(Unresolved from 4/16 #5, 4/17 #43.)*

42. **Services receive the entire `*ports.Repository` struct** (`tenants_service.go:15`; `data_sources_service.go:15`) -- Each service gets access to all repositories. Violates Interface Segregation. *(Unresolved from 4/16 #6, 4/17 #44.)*

43. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-19`) -- Same pattern as other services. *(Unresolved from 4/16 #7, 4/17 #45.)*

44. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:37-62`) -- `Create` validates the command and calls `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources. *(Unresolved from 4/16 #8, 4/17 #46.)*

45. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. *(Unresolved from 4/16 #10, 4/17 #48.)*

46. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 11 handlers). *(Unresolved from 4/16 #11, 4/17 #49.)*

47. **Tenants service has no tenant-scoping middleware** -- Unlike forms which extracts `X-Tenant-ID` via middleware into context, tenants relies on URL path params with no validation that the tenant exists. Inconsistent multi-tenancy approach across services. *(New.)*

#### Missing Functionality

48. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12, 4/17 #50.)*

49. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` has some validation (attribute type matching) but does not validate empty `TenantID`, empty `Type`, or field lengths. *(Unresolved from 4/16 #13, 4/17 #51.)*

50. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14, 4/17 #52.)*

51. **Incomplete error-to-HTTP mapping** -- `ErrInvalidID` and `ErrUnauthorized` fall through to 500. Service-specific errors like `ErrDataSourceAttrParse`, `ErrInvalidSourceTypeAttributes`, and `ErrStrategyNotFound` are also unmapped. The service-level `sendErrorResponse` is a pass-through with only a `default` case. *(Unresolved from 4/16 #15, 4/17 #53.)*

52. **No real authentication**. *(Unresolved from 4/16 #16, 4/17 #54.)*

53. **`Lookup` service method is a stub** (`data_sources_service.go:105-114`) -- Always returns `nil, nil` after verifying the data source exists. Contains `// TODO`. *(Unresolved from 4/16 #17, 4/17 #55.)*

54. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names. *(Unresolved from 4/16 #18, 4/17 #56.)*

55. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct with two fields, no json tags, no invariant enforcement. *(Unresolved from 4/16 #19, 4/17 #57.)*

56. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string for `Type`. The `isValidAttributeType` check implicitly catches unknown types but with a misleading error message. *(Unresolved from 4/16 #22, 4/17 #58.)*

#### Code Quality

57. **`DataSourceAttributes` typed as `any`** (`data_source_attributes.go:3`) -- Should be a sealed interface with a marker method. *(Unresolved from 4/16 #20, 4/17 #59.)*

58. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:69`). *(Unresolved from 4/16 #24, 4/17 #60.)*

59. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. *(Unresolved from 4/16 #29, 4/17 #61.)*

---

### Shared Package (`pkg/common`)

#### Bugs

60. **`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`** (`httputil/http.go:56-87`) -- These sentinel errors are defined in `error.go` but not handled in the switch. They fall through to the 500 default. `ErrInvalidID` should map to 400; `ErrUnauthorized` should map to 401 or 403. *(Unresolved from 4/17 #64.)*

61. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`httputil/http.go:41`) -- The `headers ...http.Header` variadic parameter is accepted but the body never iterates or sets them on the response. *(Unresolved from 4/17 #65.)*

62. **`w.Write` error ignored** (`httputil/http.go:50`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. *(Unresolved from 4/17 #66.)*

#### Code Quality

63. **`ValidateStruct` has a redundant pattern** (`validate/validate.go:19-25`) -- `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. *(Unresolved from 4/17 #67.)*

64. **In-memory transactions are no-ops** (`database/inmemory_database.go`) -- `BeginTx` stores an empty struct in context; `CommitTx`/`RollbackTx` return `nil` in all code paths. No atomicity guarantees. *(Unresolved from 4/17 #68.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 22 | `Database` nil -- `Close()` will panic | Submissions |
| **P1** | 60 | `ErrInvalidID`/`ErrUnauthorized` map to 500 | Shared |
| **P1** | 40 | `ErrDataSourceAttrParse` not wrapped on unmarshal | Tenants |
| **P2** | 2, 23, 41 | `core.go` imports adapters (hex violation) | All |
| **P2** | 3, 42 | Services receive full Repository (ISP) | Forms, Tenants |
| **P2** | 4, 26, 43 | Handlers receive full Application | All |
| **P2** | 5, 25 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 11, 49 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 32 | No domain constructors | Submissions |
| **P2** | 44 | DataSource created without parent Tenant check | Tenants |
| **P2** | 7, 27, 46 | Unnecessary goroutine/channel pattern | All |
| **P2** | 24 | No tenant middleware | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 9 | `internal/types/error.go` breaks domain purity | Forms |
| **P2** | 47 | Inconsistent multi-tenancy approach | Tenants |
| **P3** | 10, 29, 48 | Zero test files | All |
| **P3** | 12, 34, 50 | No domain events | All |
| **P3** | 13, 51 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 17, 37, 57 | `any`-typed attributes (no type safety) | All |
| **P3** | 8, 58 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 59 | Inconsistent response envelope | Tenants |
| **P3** | 20, 64 | In-memory transactions are no-ops | Forms, Shared |
| **P3** | 18 | `FieldResponse` DTO omits Attributes | Forms |
| **P3** | 63 | Redundant validation patterns | Shared |
| **P3** | 21 | Inconsistent attribute parsing pattern | Forms, Tenants |
| **P3** | 28 | `NewFindByIdQuery` creates validator per call | Submissions |

---

## Summary

### Progress Since 4/17

Nine issues from the prior review have been resolved (one partially):

- **Submissions service is now functional end-to-end** -- The entry point bootstraps an HTTP server, the service is wired, all 6 routes have real handlers, and the router was switched to `go-chi/chi` for consistency. The `Database` nil issue remains the sole P0 crash.
- **`Tenant.DataSources` field removed** -- Eliminates both the "never populated" issue and the confusing nil `DataSources` in API responses.
- **`Remove` operations now return `ErrNotFound`** -- Service-layer `Exists` guards prevent silent success on non-existent entities.
- **`DataSourceResponse` maps `Name` and `Description`** (unstaged) -- Eliminates silent data loss in data source API responses.
- **Stray semicolon removed** from `data_sources_service.go`.
- **Response DTOs partially implemented** for submissions -- `SubmissionResponse` and mapper exist; request DTOs still empty.

### New Issues Found

1. **`internal/types/error.go` breaks domain purity** (Forms #9) -- `ErrDuplicatePosition` is defined outside `core/domain/` and creates an outward dependency from the domain layer.
2. **`NewFindByIdQuery` creates `validator.New()` per call** (Submissions #28) -- Inconsistent with the shared singleton pattern used by other services.
3. **Inconsistent attribute parsing pattern** (Forms #21) -- Forms and tenants solve the same polymorphic attribute parsing problem with different approaches.
4. **Inconsistent multi-tenancy approach** (Tenants #47) -- Forms uses header middleware, tenants uses URL path params, submissions has neither.

### Current State

**Forms Service** remains the most mature. The primary gaps are: the hardcoded `"placeholder"` user IDs, the aggregate boundary ambiguity between `Form` and `Version`, the `FieldResponse` DTO silently dropping attributes, the domain error defined outside the domain package, and the continued absence of domain validation and test coverage. The hexagonal dependency violation in `core.go` and the ISP violation in the service layer remain the top architectural concerns.

**Tenants Service** has made incremental progress: the `DataSources` field was removed from `Tenant` (cleaning up the domain model), `Remove` operations are now guarded at the service layer, and the `DataSourceResponse` DTO fix is in progress (unstaged). The `DataSource` orphan problem (no tenant existence check), incomplete attribute types (missing json tags, empty `ScheduledDataSourceAttributes`), and the `parseAttributes` error wrapping gap remain.

**Submissions Service** has made significant progress -- the service is now functional with a real entry point, wired services, chi router, and actual handler implementations. The remaining P0 is the nil `Database` in the in-memory bootstrap. Beyond that, the service still lacks tenant middleware, write operations, domain constructors, request DTOs, and test coverage.

**Shared Package** (`pkg/common`) is unchanged. Missing `ErrInvalidID`/`ErrUnauthorized` HTTP mappings, the unused `headers` parameter in `SendJsonResponse`, the ignored `w.Write` error, redundant validation patterns, and the no-op in-memory transactions remain.

### Highest-Impact Improvements

1. **Set `Database` in submissions in-memory bootstrap** (P0 -- `Close()` will panic)
2. **Add `ErrInvalidID` and `ErrUnauthorized` mappings** to `SendErrorResponse` (P1 -- auth/validation errors produce 500s)
3. **Wrap `parseAttributes` errors with `ErrDataSourceAttrParse`** in `dto/request.go` (P1 -- malformed attribute JSON produces 500)
4. **Fix the hexagonal dependency violation** in `core.go` across all three services (P2 -- foundational architecture)
5. **Move `ErrDuplicatePosition` into `core/domain/`** in the forms service (P2 -- domain purity)
6. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
