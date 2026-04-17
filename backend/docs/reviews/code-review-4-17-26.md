# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/16 Review

1. ~~Handlers swallow JSON parse errors~~ (Tenants #1) -- All four mutating tenant handlers (`createTenant`, `updateTenant`, `createDataSource`, `updateDataSource`) now call `h.sendErrorResponse(w, err)` when `ReadJsonPayload` fails (`handlers.go:84-87`, `122-125`, `228-231`, `269-272`).

2. ~~`updateTenant` success message typo `"Success updated!"`~~ (Tenants #4) -- Now correctly says `"Successfully updated!"` (`handlers.go:141`).

3. ~~Request DTO `Attributes` field uses `any` with no mapping logic~~ (Tenants #21) -- A strategy pattern in `dto/request.go` now maps inbound JSON to typed `DataSourceAttributes` implementations based on the `Type` field via `RequestToDataSourceAttributes`.

4. ~~`CreateDataSourceCommand` / `UpdateDataSourceCommand` have no `validate` struct tags~~ (Tenants #23) -- `baseDataSourceCommand` now has `validate:"required"` on `TenantID`, `validate:"oneof=static scheduled query"` on `Type`, and `validate:"required"` on `Attributes` (`commands.go`).

5. ~~`"Persistance"` typo throughout~~ (Tenants #25) -- All names now use the correct spelling: `PersistenceDriver`, `PersistenceSettings`, `PersistenceOptions` (`persistence.go`).

6. ~~`updateDataSource` sends `nil` to error handler~~ -- The `select` case now correctly passes `res.err` (`handlers.go:292`).

7. ~~`createTenant` / `updateTenant` write two HTTP responses for validation errors~~ -- The redundant inner `if validate.IsValidationErr(err)` blocks have been removed. Both handlers now delegate directly to `sendErrorResponse`.

8. ~~Inconsistent response mapper naming~~ -- `DataSourceToResponseDto` renamed to `DataSourceToResponse` for consistency with `TenantToResponse` (`dto/response.go:36`).

9. ~~`getForms` sends `res.data` instead of `dtos`~~ (Forms #1) -- `handlers.go:51` now correctly sends the mapped `dtos` slice.

10. ~~Receiver name mismatch in `section.go`~~ (Forms #6) -- `UpdateFields` now uses receiver `s`.

11. ~~`DateFieldAttributes` missing `BaseFieldAttributes` embedding~~ (Forms #20) -- Now properly embeds `BaseFieldAttributes` (`field_attributes.go:37-41`).

12. ~~Services silently discard `domain.New*` constructor errors~~ (Tenants #2) -- Both `tenants_service.go` and `data_sources_service.go` now properly check and return errors from `domain.NewTenant` and `domain.NewDataSource`.

---

## Remaining Issues

### Forms Service

#### Bugs

1. **Handlers swallow JSON parse errors** (`handlers.go:96-98`, `135-137`, `240-242`) -- `createForm`, `updateForm`, and `createVersion` `return` without writing an HTTP response when `ReadJsonPayload` fails. The client receives an empty response with a 200 status code. `updateVersion` (`handlers.go:280-283`) correctly calls `SendErrorResponse`. *(Partially unresolved from 4/13 #3.)*

2. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:329`, `363`) -- The publish/retire state transitions record a fake user. Not a real auth integration.

#### Architectural

3. **`core.go` imports from the adapters layer** (`core.go:7`, `13`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, violating the hexagonal dependency rule. The core package should not depend on adapter types. Configuration should be defined in core or injected as plain values. *(Unresolved from 4/13 #7.)*

4. **Services receive the entire `Repository` struct** (`form_service.go`) -- The form service gets access to `Database` and the full repository aggregate rather than just the `FormsRepository` interface. Violates Interface Segregation. *(Unresolved from 4/13 #8.)*

5. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-20`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services` or the specific service interface, giving the HTTP layer access to the logger, repository, and any future internal state. *(Unresolved from 4/13 #9.)*

6. **`Find()` has no tenant filtering** (`form_service.go`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. *(Unresolved from 4/13 #10.)*

7. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. From a DDD perspective, if `Version` is within the `Form` aggregate, mutations should go through the aggregate root. If it's its own aggregate, the relationship should be modeled as a reference only. *(Unresolved from 4/13 #11.)*

8. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 10 handlers) -- `net/http` already runs each handler in its own goroutine. The extra goroutine + buffered channel + `select` pattern adds allocation overhead and complexity with no concurrency benefit. *(Unresolved from 4/13 #25.)*

9. **`time.Now()` called directly in the repository layer** (`forms_repository.go`) -- Couples persistence to wall-clock time, making deterministic tests impossible. Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26.)*

#### Missing Functionality

10. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12.)*

11. **Domain validation unimplemented** (`form.go`, `field.go`, `page.go`, `section.go`) -- All entity constructors contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13.)*

12. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14.)*

13. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and domain errors (`ErrVersionLocked`, `ErrInvalidVersion`, etc.) all fall through to the default 500 case in `common.SendErrorResponse`. *(Unresolved from 4/13 #15.)*

14. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16.)*

15. **No `Delete` operation for forms.** *(Unresolved from 4/13 #17.)*

16. **`ConditionalRule` is an empty stub** (`conditional_rule.go`) -- Contains only an ID; no rule type, conditions, or actions. *(Unresolved from 4/13 #18.)*

#### Code Quality

17. **`FieldAttributes` is typed as `any`** (`field_attributes.go:3`) -- No type safety; any value can be assigned. Should use a sealed interface with a marker method (e.g., `isFieldAttributes()`) to restrict implementations to known types. *(Unresolved from 4/13 #19.)*

18. **Forms commands have no `validate` struct tags** (`commands.go`) -- `CreateFormCommand`, `UpdateFormCommand`, `CreateVersionCommand`, `PublishVersionCommand`, and `RetireVersionCommand` have no tags, so the `ValidateStruct` call in the service layer is a no-op. *(Unresolved from 4/13 #22.)*

19. **`validator.New()` created per `NewUpdateVersionCommand` call** (`commands.go:71`) -- Expensive; should use the package-level singleton from `pkg/common/validate`. *(Unresolved from 4/13 #21.)*

20. **Inconsistent constructor signatures** -- Forms domain constructors return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. *(Unresolved from 4/13 #23.)*

21. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing; the `CreateVersion` flow has a race condition since read and write locks are not held atomically. *(Unresolved from 4/13 #29.)*

---

### Submissions Service

#### Bugs

22. **`nil` route handlers will panic at runtime** (`routes.go`) -- 4 of 6 route registrations pass `nil` as the handler function: `POST /submissions`, `GET .../attempts`, `GET .../status`, `POST .../replay`. Any request to these endpoints will cause a nil pointer dereference panic.

23. **`Database` field is nil in in-memory bootstrap** (`inmemory/inmemory.go`) -- `Repository.Database` is not set, so `Application.Close()` will panic with a nil pointer dereference since `Close()` calls `a.repository.Database.Close()`.

24. **Services bootstrap does not wire `SubmissionsService`** (`services/services.go`) -- `Bootstrap()` returns `ports.Services{Submissions: nil}`. The implementation exists in `submissions_service.go` but is never instantiated. Any handler calling `app.Services.Submissions.*` will nil-pointer panic.

#### Architectural

25. **Entry point is a stub** (`cmd/server/main.go`) -- The `main()` function only prints `"Submissions Service"` and exits. No HTTP server is started, no application is bootstrapped.

26. **`core.go` imports from the adapters layer** (`core.go:7`, `13`) -- Same hexagonal dependency violation as the other services: `ApplicationSettings` references `persistence.PersistenceSettings`. *(Same pattern as Forms #3, Tenants #46.)*

27. **REST uses `net/http.ServeMux` while other services use `go-chi/chi`** (`routes.go`) -- Inconsistent router choice across services. Forms and tenants use chi; submissions uses the stdlib mux. This also means submissions lacks chi's middleware and path parameter features.

28. **No tenant middleware** -- `middleware.go` is empty. Unlike forms and tenants, there is no `X-Tenant-ID` extraction. Yet the service layer's `FindById` and `FindByReferenceId` enforce tenant authorization via the query object's `TenantID` field. Without middleware, tenant ID is never available.

29. **`Find()` has no tenant filtering** (`submissions_service.go`) -- Returns all submissions across all tenants. *(Same pattern as Forms #6.)*

30. **REST handlers hold a reference to the full `Application`** (`handlers.go`) -- Same pattern as other services.

31. **Unnecessary goroutine/channel pattern** in implemented handlers -- Same pattern as other services.

#### Missing Functionality

32. **Zero test files** in the entire submissions service.

33. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go`) -- Return `nil, nil` and `nil` respectively.

34. **DTOs not implemented** -- `dto.go` is empty. Handlers have `// TODO` comments noting domain-to-DTO conversion is missing. Handlers currently send raw domain objects.

35. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods.

36. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`.

37. **No domain events** for cross-service communication.

38. **No real authentication.**

#### Code Quality

39. **`Payload` typed as `any`** (`submission.go`) -- Same `any` issue as `FieldAttributes` and `DataSourceAttributes` in other services. No type safety.

40. **`SubmissionStatus` has no defined constants** -- Unlike `VersionStatus` in forms or `DataSourceType` in tenants, there are no `const` values for valid submission statuses.

41. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- No secondary index on `ReferenceID`. Acceptable for in-memory dev but will not scale.

---

### Tenants Service

#### Bugs

42. **Route bug: `removeDataSource` mapped to `PUT` instead of `DELETE`** (`routes.go:31`) -- `dataSourceRoutes.Put("/", h.removeDataSource)` should be `dataSourceRoutes.Delete("/", h.removeDataSource)`. The delete endpoint is unreachable; `PUT` requests to this path will hit `removeDataSource` instead of `updateDataSource` since chi will use the last registered handler for the same method+path.

43. **JSON parse errors map to 500 instead of 400** (`handlers.go:84-87`, `122-125`, `228-231`, `269-272`) -- While handlers now call `sendErrorResponse`, the `json.Decoder` error from `ReadJsonPayload` is not a `validator.ValidationErrors` and is not `ErrDataSourceAttrParse`, so it falls through to `common.SendErrorResponse` which defaults to 500. Malformed JSON from clients produces a 500 Internal Server Error instead of a 400 Bad Request. *(Unresolved from 4/17 #1.)*

44. **`ErrDataSourceAttrParse` is defined but never returned** (`dto/request.go:9`) -- The sentinel error is defined and checked in `sendErrorResponse` (`handlers.go:359`), but `RequestToDataSourceAttributes` and `parseAttributes` never return or wrap it. The `errors.Is(err, dto.ErrDataSourceAttrParse)` check is dead code. *(Unresolved from 4/17 #4.)*

45. **`Remove` operations succeed silently for non-existent entities** (`tenant_repository.go:78-84`; `data_sources_repository.go:84-91`) -- Go's `delete` on a map is a no-op for missing keys, so deleting a non-existent tenant or data source returns `nil` (HTTP 204) instead of `ErrNotFound` (HTTP 404). *(Unresolved from 4/16 #3.)*

#### Architectural

46. **`core.go` imports from the adapters layer** (`core.go:7`, `13`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, violating the hexagonal dependency rule. *(Unresolved from 4/16 #5.)*

47. **Services receive the entire `*ports.Repository` struct** (`tenants_service.go:14`; `data_sources_service.go:14`) -- Each service gets access to `Database`, `Tenants`, AND `DataSources` repositories. Violates Interface Segregation. *(Unresolved from 4/16 #6.)*

48. **REST handlers hold a reference to the full `Application`** (`handlers.go:20-22`) -- Same pattern as other services. *(Unresolved from 4/16 #7.)*

49. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:36-54`) -- `Create` and `Update` call `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources. *(Unresolved from 4/16 #8.)*

50. **`Tenant.DataSources` field is declared but never populated** (`tenant.go`) -- The `Tenant` struct previously had a `DataSources` field that is never set in any code path, giving a false impression that `Tenant` is an aggregate root containing its data sources. Either populate it or remove it. *(Unresolved from 4/16 #9.)*

51. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:24-26`) -- Returns every tenant in a single unbounded response. *(Unresolved from 4/16 #10.)*

52. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 11 handlers) -- Same pattern as other services. *(Unresolved from 4/16 #11.)*

#### Missing Functionality

53. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12.)*

54. **Domain validation unimplemented** (`tenant.go:17-23`; `data_source.go`) -- Constructors return `(*Entity, error)` but never validate. *(Unresolved from 4/16 #13.)*

55. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14.)*

56. **Incomplete error-to-HTTP mapping** (`http.go:48-72`) -- `ErrInvalidID` and `ErrUnauthorized` are defined in `error.go` but not mapped in `SendErrorResponse`. They fall through to the default 500 case. *(Unresolved from 4/16 #15.)*

57. **No real authentication**. *(Unresolved from 4/16 #16.)*

58. **`Lookup` service method is a stub** (`data_sources_service.go:80-89`) -- Always returns `nil, nil` after verifying the data source exists. *(Unresolved from 4/16 #17.)*

59. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names (e.g., `"Data"` instead of `"data"`). *(Unresolved from 4/16 #18.)*

60. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct, no invariant enforcement. *(Unresolved from 4/16 #19.)*

61. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string. *(Partially resolved from 4/16 #22.)*

#### Code Quality

62. **`DataSourceAttributes` typed as `any`** (`data_source_attributes.go:3`) -- Equivalent to `any`. Should be a sealed interface with a marker method. *(Unresolved from 4/16 #20.)*

63. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:57`; `data_sources_repository.go:61`) -- Couples persistence to wall-clock time. *(Unresolved from 4/16 #24.)*

64. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. Clients must handle two different response shapes. *(Unresolved from 4/16 #29.)*

65. **`getTenants` returns tenants with `DataSources` always `nil`** (`handlers.go:48-51`) -- The domain model declares the field but no code path populates it. *(Unresolved from 4/16 #30.)*

66. **Stray semicolon** (`data_sources_service.go:28`) -- Bare `;` on its own line.

---

### Shared Package (`pkg/common`)

#### Bugs

67. **`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`** (`http.go:48-72`) -- These sentinel errors are defined in `error.go` but not handled in the switch. They fall through to the 500 default case. `ErrInvalidID` should map to 400; `ErrUnauthorized` should map to 401 or 403.

68. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`http.go:33`) -- The `headers ...http.Header` variadic parameter is accepted but the body never iterates or sets them on the response. Callers expecting custom headers will be silently ignored.

69. **`w.Write` error ignored** (`http.go:42`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. Should at minimum log it.

#### Code Quality

70. **`ValidateStruct` has a redundant pattern** (`validate.go:19-25`) -- `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. Idiomatic Go prefers the simpler form.

71. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing. No atomicity guarantees for multi-step operations like `CreateVersion` in the forms service.

72. **No `json.Decoder` error mapping** -- `json.SyntaxError`, `json.UnmarshalTypeError`, and related decoding errors are not mapped in `SendErrorResponse`. Any service delegating JSON parse errors to the shared handler will produce a 500 for malformed client input. Should map these to HTTP 400.

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 22 | `nil` route handlers will panic at runtime | Submissions |
| **P0** | 23 | `Database` nil -- `Close()` will panic | Submissions |
| **P0** | 24 | `SubmissionsService` never wired -- nil pointer | Submissions |
| **P0** | 42 | `removeDataSource` mapped to `PUT` not `DELETE` | Tenants |
| **P0** | 25 | Entry point is a stub (no HTTP server) | Submissions |
| **P1** | 1 | Handlers swallow JSON parse errors (empty 200) | Forms |
| **P1** | 43 | JSON parse errors map to 500 | Tenants |
| **P1** | 72 | No `json.Decoder` error mapping in shared pkg | Shared |
| **P1** | 67 | `ErrInvalidID`/`ErrUnauthorized` map to 500 | Shared |
| **P1** | 44 | `ErrDataSourceAttrParse` dead code | Tenants |
| **P1** | 45 | `Remove` returns 204 for non-existent entities | Tenants |
| **P2** | 3, 26, 46 | `core.go` imports adapters (hex violation) | All |
| **P2** | 4, 47 | Services receive full Repository (ISP) | Forms, Tenants |
| **P2** | 5, 30, 48 | Handlers receive full Application | All |
| **P2** | 6, 29 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 11, 54 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 35 | No domain constructors | Submissions |
| **P2** | 49 | DataSource created without parent Tenant check | Tenants |
| **P2** | 18 | Forms commands have no `validate` struct tags | Forms |
| **P2** | 8, 31, 52 | Unnecessary goroutine/channel pattern | All |
| **P2** | 27 | Inconsistent router (stdlib vs chi) | Submissions |
| **P2** | 28 | No tenant middleware | Submissions |
| **P3** | 10, 32, 53 | Zero test files | All |
| **P3** | 12, 37, 55 | No domain events | All |
| **P3** | 13, 56 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 17, 39, 62 | `any`-typed attributes (no type safety) | All |
| **P3** | 9, 63 | `time.Now()` in repository layer | Forms, Tenants |
| **P3** | 64 | Inconsistent response envelope | Tenants |
| **P3** | 19 | `validator.New()` per command constructor | Forms |
| **P3** | 21, 71 | In-memory transactions are no-ops | Forms, Shared |

---

## Summary

### Progress Since 4/16

Twelve issues from the prior review cycle have been resolved. The tenants service continues its accelerated improvement pace, fixing JSON parse error handling in handlers, adding validation struct tags to data source commands, implementing the strategy pattern for attribute deserialization, correcting the persistence typo, and properly propagating domain constructor errors. The forms service also resolved three issues: the `getForms` DTO bug, the `section.go` receiver mismatch, and the missing `DateFieldAttributes` base embedding.

### Current State

**Forms Service** is the most mature of the three. Its domain model is well-designed with a clear version lifecycle state machine, proper page/section/field hierarchy with position enforcement, and bidirectional DTO mapping. The primary gaps are three handlers that still swallow JSON parse errors, no validation struct tags on form commands, and the aggregate boundary between `Form` and `Version` being unclear from a DDD perspective.

**Tenants Service** has a clean hexagonal structure and has resolved most of the DTO/validation issues from prior reviews. A new route registration bug was found (`removeDataSource` on `PUT` instead of `DELETE`), and the JSON-parse-to-500 pipeline remains the most impactful open bug. The `DataSource` orphan problem (no tenant existence check) is the top architectural concern.

**Submissions Service** is in early development and has the most critical issues: the entry point is a stub, the service implementation is never wired, `nil` handlers will panic, and the `Database` nil reference will panic on shutdown. The service needs fundamental wiring before it can function.

**Shared Package** (`pkg/common`) has three issues that affect all services: missing `ErrInvalidID`/`ErrUnauthorized` error-to-HTTP mappings, no `json.Decoder` error handling (making JSON parse errors 500s globally), and the unused `headers` parameter in `SendJsonResponse`.

### Highest-Impact Improvements

1. **Wire the submissions service** -- fix the entry point, bootstrap wiring, nil handlers, and nil database (P0 -- service is non-functional)
2. **Fix the tenants `removeDataSource` route** from `PUT` to `DELETE` (P0 -- wrong HTTP method)
3. **Map `json.Decoder` errors to HTTP 400** in `pkg/common/http.go` and fix the three forms handlers that swallow parse errors (P1 -- client errors produce 500s or empty responses)
4. **Add `ErrInvalidID` and `ErrUnauthorized` mappings** to `SendErrorResponse` (P1 -- auth/validation errors produce 500s)
5. **Fix the hexagonal dependency violation** in `core.go` across all three services (P2 -- foundational architecture)
6. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
