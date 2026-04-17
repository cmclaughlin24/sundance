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

## Issues Resolved Since 4/17 Review

13. ~~Handlers swallow JSON parse errors~~ (Forms #1) -- `createForm` (`handlers.go:97-100`), `updateForm` (`handlers.go:137-140`), and `createVersion` (`handlers.go:243-246`) now call `h.sendErrorResponse(w, err)` when `ReadJsonPayload` fails, instead of returning without writing an HTTP response.

14. ~~Forms commands have no `validate` struct tags~~ (Forms #18) -- All command structs now have validation tags: `baseFormCommand` has `validate:"required"` on `TenantID`, `validate:"required,max=75"` on `Name`, `validate:"required,max=250"` on `Description`; `baseVersionCommand` has `validate:"required"` on `FormID` and `TenantID`; `PublishVersionCommand` and `RetireVersionCommand` have `validate:"required"` on `VersionID` and `UserID`.

15. ~~`validator.New()` created per `NewUpdateVersionCommand` call~~ (Forms #19) -- Command constructors no longer perform validation themselves. Validation is handled in the service layer via `validate.ValidateStruct(command)`, which uses the package-level singleton validator in `pkg/common/validate/validate.go`.

16. ~~Route bug: `removeDataSource` mapped to `PUT` instead of `DELETE`~~ (Tenants #42) -- `routes.go:31` now correctly uses `dataSourceRoutes.Delete("/", h.removeDataSource)`.

17. ~~`ErrDataSourceAttrParse` is defined but never returned~~ (Tenants #44, partially) -- `RequestToDataSourceAttributes` now wraps errors from `strategies.Get()` (unknown type) with `ErrDataSourceAttrParse` (`request.go:53`). However, errors from `parseAttributes` (JSON unmarshal failures) are still returned unwrapped. See remaining issue #41 below.

18. ~~JSON parse errors map to 500 instead of 400~~ (Tenants #43) -- `ReadJsonPayload` in `pkg/common/http.go` now wraps `json.Decoder` errors with `ErrDecodeJSON` (`http.go:30`), and `SendErrorResponse` maps `ErrDecodeJSON` to HTTP 400 (`http.go:59-64`). The tenants `sendErrorResponse` delegates unmatched errors to `common.SendErrorResponse`, so JSON parse errors from `ReadJsonPayload` now correctly produce 400 responses.

19. ~~No `json.Decoder` error mapping~~ (Shared #72) -- `ReadJsonPayload` wraps all decode errors with `ErrDecodeJSON` via `fmt.Errorf("%w: %w", ErrDecodeJSON, err)` (`http.go:30`), and `SendErrorResponse` checks for `ErrDecodeJSON` as its first case (`http.go:59`). This resolves the global JSON parse-to-500 issue for any service that uses the shared `ReadJsonPayload`/`SendErrorResponse` pipeline.

20. ~~`Exists` key mismatch in `data_sources_repository.go`~~ (New from this review) -- `Exists` was looking up by `string(id)` instead of the composite `getDataSourceKey(tenantID, id)` key used by all other methods. Fixed in the working tree to use the composite key (`data_sources_repository.go:60`).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:329-330`, `364-365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME: Remove temporary placeholder for user ID.` comments but remain unresolved. Not a real auth integration. *(Unresolved from 4/13 #3, 4/17 #2.)*

#### Architectural

2. **`core.go` imports from the adapters layer** (`core.go:7`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, and `NewApplication` calls `persistence.Bootstrap` directly. The core package should not depend on adapter types. Configuration should be defined in core or injected as plain values. *(Unresolved from 4/13 #7.)*

3. **Services receive the entire `Repository` struct** (`form_service.go:14-17`) -- The form service holds `*ports.Repository` which includes `Database` and the full repository aggregate rather than just the `FormsRepository` interface. Violates Interface Segregation. *(Unresolved from 4/13 #8.)*

4. **REST handlers hold a reference to the full `Application`** (`handlers.go:19-21`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services` or the specific service interface, giving the HTTP layer access to the logger, repository, and any future internal state. *(Unresolved from 4/13 #9.)*

5. **`Find()` has no tenant filtering** (`form_service.go:26-28`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. *(Unresolved from 4/13 #10.)*

6. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. From a DDD perspective, if `Version` is within the `Form` aggregate, mutations should go through the aggregate root. If it's its own aggregate, the relationship should be modeled as a reference only. *(Unresolved from 4/13 #11.)*

7. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 10 handlers) -- `net/http` already runs each handler in its own goroutine. The extra goroutine + buffered channel + `select` pattern adds allocation overhead and complexity with no concurrency benefit. *(Unresolved from 4/13 #25.)*

8. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `form_service.go` in `PublishVersion` and `RetireVersion`) -- Couples persistence and business logic to wall-clock time, making deterministic tests impossible. Should be injected via a `Clock` interface or function. *(Expanded from 4/13 #26.)*

#### Missing Functionality

9. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12.)*

10. **Domain validation unimplemented** (`form.go`, `field.go`, `page.go`, `section.go`, `version.go`) -- All entity constructors contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13.)*

11. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14.)*

12. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and domain errors (`ErrVersionLocked`, `ErrInvalidVersion`, etc.) all fall through to the default 500 case in `common.SendErrorResponse`. *(Unresolved from 4/13 #15.)*

13. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16.)*

14. **No `Delete` operation for forms.** *(Unresolved from 4/13 #17.)*

15. **`ConditionalRule` is an empty stub** (`conditional_rule.go`) -- Contains only an ID; no rule type, conditions, or actions. *(Unresolved from 4/13 #18.)*

#### Code Quality

16. **`FieldAttributes` is typed as `any`** (`field_attributes.go:3`) -- No type safety; any value can be assigned. Concrete types (`TextFieldAttributes`, `NumberFieldAttributes`, etc.) are defined but the type alias itself is `any`, so `Field.Attributes` has no compile-time safety. Should use a sealed interface with a marker method (e.g., `isFieldAttributes()`). *(Unresolved from 4/13 #19.)*

17. **`FieldResponse` DTO omits `Attributes`** (`dto/field.go`) -- `FieldToResponse` maps all fields except `Attributes`, meaning field attribute data is silently lost in API responses. *(New.)*

18. **Inconsistent constructor signatures** -- Forms domain constructors return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. *(Unresolved from 4/13 #23.)*

19. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing; the `CreateVersion` flow has a race condition since read and write locks are not held atomically. *(Unresolved from 4/13 #29.)*

---

### Submissions Service

#### Bugs

20. **`nil` route handlers will panic at runtime** (`routes.go`) -- 4 of 6 route registrations pass `nil` as the handler function: `POST /submissions`, `GET .../attempts`, `GET .../status`, `POST .../replay`. Any request to these endpoints will cause a nil pointer dereference panic. *(Unresolved from 4/17 #22.)*

21. **`Database` field is nil in in-memory bootstrap** (`inmemory/inmemory.go`) -- `Repository.Database` is not set, so `Application.Close()` will panic with a nil pointer dereference since `Close()` calls `a.repository.Database.Close()`. *(Unresolved from 4/17 #23.)*

22. **Services bootstrap does not wire `SubmissionsService`** (`services/services.go`) -- `Bootstrap()` returns `ports.Services{Submissions: nil}`. The implementation exists in `submissions_service.go` but is never instantiated. Any handler calling `app.Services.Submissions.*` will nil-pointer panic. *(Unresolved from 4/17 #24.)*

#### Architectural

23. **Entry point is a stub** (`cmd/server/main.go`) -- The `main()` function only prints `"Submissions Service"` and exits. No HTTP server is started, no application is bootstrapped. *(Unresolved from 4/17 #25.)*

24. **`core.go` imports from the adapters layer** (`core.go`) -- Same hexagonal dependency violation as the other services: `ApplicationSettings` references `persistence.PersistenceSettings`. *(Same pattern as Forms #2, Tenants #43. Unresolved from 4/17 #26.)*

25. **REST uses `net/http.ServeMux` while other services use `go-chi/chi`** (`routes.go`) -- Inconsistent router choice across services. Forms and tenants use chi; submissions uses the stdlib mux. This also means submissions lacks chi's middleware and path parameter features. *(Unresolved from 4/17 #27.)*

26. **No tenant middleware** -- `middleware.go` is empty (0 lines, no package declaration). Unlike forms and tenants, there is no `X-Tenant-ID` extraction. Yet the service layer's `FindById` and `FindByReferenceId` enforce tenant authorization via the query object's `TenantID` field. Without middleware, tenant ID is never available. *(Unresolved from 4/17 #28.)*

27. **`Find()` has no tenant filtering** (`submissions_service.go`) -- Returns all submissions across all tenants. *(Same pattern as Forms #5. Unresolved from 4/17 #29.)*

28. **REST handlers hold a reference to the full `Application`** (`handlers.go`) -- Same pattern as other services. *(Unresolved from 4/17 #30.)*

29. **Unnecessary goroutine/channel pattern** in implemented handlers -- Same pattern as other services. *(Unresolved from 4/17 #31.)*

#### Missing Functionality

30. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #32.)*

31. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go`) -- Return `nil, nil` and `nil` respectively. *(Unresolved from 4/17 #33.)*

32. **DTOs not implemented** -- `dto.go` is empty (0 lines). Handlers have `// TODO` comments noting domain-to-DTO conversion is missing. Handlers currently send raw domain objects. *(Unresolved from 4/17 #34.)*

33. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #35.)*

34. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #36.)*

35. **No domain events** for cross-service communication. *(Unresolved from 4/17 #37.)*

36. **No real authentication.** *(Unresolved from 4/17 #38.)*

37. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields, making it impossible to specify what to replay. *(New.)*

#### Code Quality

38. **`Payload` typed as `any`** (`submission.go`) -- Same `any` issue as `FieldAttributes` and `DataSourceAttributes` in other services. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. No type safety. *(Unresolved from 4/17 #39.)*

39. **`SubmissionStatus` has no defined constants** -- Unlike `VersionStatus` in forms or `DataSourceType` in tenants, there are no `const` values for valid submission statuses. *(Unresolved from 4/17 #40.)*

40. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- No secondary index on `ReferenceID`. Acceptable for in-memory dev but will not scale. *(Unresolved from 4/17 #41.)*

---

### Tenants Service

#### Bugs

41. **`ErrDataSourceAttrParse` not wrapped on JSON unmarshal failure** (`dto/request.go:59-66`) -- `parseAttributes` returns raw `json.Unmarshal` errors without wrapping them in `ErrDataSourceAttrParse`. Only the `strategies.Get()` error path (unknown type) wraps the sentinel. When attribute JSON is malformed, the error bypasses the `errors.Is(err, dto.ErrDataSourceAttrParse)` check in `sendErrorResponse` and falls through to `common.SendErrorResponse`. Since `ReadJsonPayload` is not involved here (the body was already decoded), `ErrDecodeJSON` is also not present, so the error hits the 500 default. *(Partially resolved from 4/17 #44 -- unknown type path fixed, unmarshal path still broken.)*

42. **`Remove` operations succeed silently for non-existent entities at the repository layer** (`tenant_repository.go:85-91`; `data_sources_repository.go:91-98`) -- Go's `delete` on a map is a no-op for missing keys, so the repository returns `nil` regardless. The service layer mitigates this by calling `Exists` first, but the repository contract is still incorrect -- a `Remove` for a non-existent key should return `ErrNotFound`. *(Unresolved from 4/16 #3. Service-layer mitigation noted.)*

#### Architectural

43. **`core.go` imports from the adapters layer** (`core.go:7`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, violating the hexagonal dependency rule. *(Unresolved from 4/16 #5.)*

44. **Services receive the entire `*ports.Repository` struct** (`tenants_service.go:14`; `data_sources_service.go:14`) -- Each service gets access to `Database`, `Tenants`, AND `DataSources` repositories. Violates Interface Segregation. *(Unresolved from 4/16 #6.)*

45. **REST handlers hold a reference to the full `Application`** (`handlers.go:20-22`) -- Same pattern as other services. *(Unresolved from 4/16 #7.)*

46. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:37-55`) -- `Create` and `Update` call `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources. *(Unresolved from 4/16 #8.)*

47. **`Tenant.DataSources` field is declared but never populated** (`tenant.go`) -- The `Tenant` struct has a `DataSources` field that is never set in any code path, giving a false impression that `Tenant` is an aggregate root containing its data sources. Either populate it or remove it. *(Unresolved from 4/16 #9.)*

48. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. *(Unresolved from 4/16 #10.)*

49. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 11 handlers) -- Same pattern as other services. *(Unresolved from 4/16 #11.)*

#### Missing Functionality

50. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12.)*

51. **Domain validation unimplemented** (`tenant.go:17-23`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` has *some* validation (attribute type matching) but does not validate empty `TenantID`, empty `Type`, or field lengths. *(Partially unresolved from 4/16 #13.)*

52. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14.)*

53. **Incomplete error-to-HTTP mapping** (`http.go:58-83`) -- `ErrInvalidID` and `ErrUnauthorized` are defined in `error.go` but not mapped in `SendErrorResponse`. They fall through to the default 500 case. *(Unresolved from 4/16 #15.)*

54. **No real authentication**. *(Unresolved from 4/16 #16.)*

55. **`Lookup` service method is a stub** (`data_sources_service.go:91-101`) -- Always returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern based on the type of data source.` *(Unresolved from 4/16 #17.)*

56. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names (e.g., `"Data"` instead of `"data"`). *(Unresolved from 4/16 #18.)*

57. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct with two fields (`Code`, `Description`), no json tags, no invariant enforcement. *(Unresolved from 4/16 #19.)*

58. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string for `Type`. *(Partially resolved from 4/16 #22.)*

#### Code Quality

59. **`DataSourceAttributes` typed as `any`** (`data_source_attributes.go:3`) -- Equivalent to `any`. Should be a sealed interface with a marker method. *(Unresolved from 4/16 #20.)*

60. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:68`) -- Couples persistence to wall-clock time. *(Unresolved from 4/16 #24.)*

61. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. Clients must handle two different response shapes. *(Unresolved from 4/16 #29.)*

62. **`getTenants` returns tenants with `DataSources` always `nil`** (`handlers.go:48-51`) -- The domain model declares the field but no code path populates it. *(Unresolved from 4/16 #30.)*

63. **Stray semicolon** (`data_sources_service.go:29`) -- Bare `;` on its own line in the `Find` method. *(Unresolved from 4/17 #66.)*

---

### Shared Package (`pkg/common`)

#### Bugs

64. **`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`** (`http.go:58-83`) -- These sentinel errors are defined in `error.go` but not handled in the switch. They fall through to the 500 default case. `ErrInvalidID` should map to 400; `ErrUnauthorized` should map to 401 or 403. *(Unresolved from 4/17 #67.)*

65. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`http.go:38`) -- The `headers ...http.Header` variadic parameter is accepted but the body never iterates or sets them on the response. Callers expecting custom headers will be silently ignored. *(Unresolved from 4/17 #68.)*

66. **`w.Write` error ignored** (`http.go:47`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. Should at minimum log it. *(Unresolved from 4/17 #69.)*

#### Code Quality

67. **`ValidateStruct` has a redundant pattern** (`validate/validate.go:19-23`) -- `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. Similarly, `IsValidationErr` uses `if !ok { return false }; return true` instead of just `return ok`. Idiomatic Go prefers the simpler forms. *(Unresolved from 4/17 #70.)*

68. **In-memory transactions are no-ops** (`database/inmemory_database.go`) -- `BeginTx` stores an empty struct in context; `CommitTx`/`RollbackTx` return `nil` in all code paths. No atomicity guarantees for multi-step operations like `CreateVersion` in the forms service. *(Unresolved from 4/17 #71.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 20 | `nil` route handlers will panic at runtime | Submissions |
| **P0** | 21 | `Database` nil -- `Close()` will panic | Submissions |
| **P0** | 22 | `SubmissionsService` never wired -- nil pointer | Submissions |
| **P0** | 23 | Entry point is a stub (no HTTP server) | Submissions |
| **P1** | 64 | `ErrInvalidID`/`ErrUnauthorized` map to 500 | Shared |
| **P1** | 41 | `ErrDataSourceAttrParse` not wrapped on unmarshal | Tenants |
| **P1** | 42 | `Remove` returns nil for non-existent entities (repo layer) | Tenants |
| **P2** | 2, 24, 43 | `core.go` imports adapters (hex violation) | All |
| **P2** | 3, 44 | Services receive full Repository (ISP) | Forms, Tenants |
| **P2** | 4, 28, 45 | Handlers receive full Application | All |
| **P2** | 5, 27 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 10, 51 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 33 | No domain constructors | Submissions |
| **P2** | 46 | DataSource created without parent Tenant check | Tenants |
| **P2** | 7, 29, 49 | Unnecessary goroutine/channel pattern | All |
| **P2** | 25 | Inconsistent router (stdlib vs chi) | Submissions |
| **P2** | 26 | No tenant middleware | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P3** | 9, 30, 50 | Zero test files | All |
| **P3** | 11, 35, 52 | No domain events | All |
| **P3** | 12, 53 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 16, 38, 59 | `any`-typed attributes (no type safety) | All |
| **P3** | 8, 60 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 61 | Inconsistent response envelope | Tenants |
| **P3** | 19, 68 | In-memory transactions are no-ops | Forms, Shared |
| **P3** | 17 | `FieldResponse` DTO omits Attributes | Forms |
| **P3** | 63 | Stray semicolon | Tenants |
| **P3** | 67 | Redundant validation patterns | Shared |

---

## Summary

### Progress Since 4/17

Eight issues from the prior review cycle have been resolved, with one additional partial resolution. The most impactful fixes were:

- **Forms handlers no longer swallow JSON parse errors** -- all three problematic handlers (`createForm`, `updateForm`, `createVersion`) now call `sendErrorResponse`, eliminating the empty-200 bug that had persisted since the 4/13 review.
- **Forms commands now have validation struct tags** -- `baseFormCommand`, `baseVersionCommand`, `PublishVersionCommand`, and `RetireVersionCommand` all enforce required fields and length constraints via `validate` tags.
- **`validator.New()` per-call eliminated** -- validation is now handled by the service layer using the package-level singleton in `pkg/common/validate`.
- **`removeDataSource` route fixed** -- correctly mapped to `DELETE` instead of `PUT`.
- **`ReadJsonPayload` now wraps decode errors with `ErrDecodeJSON`** -- combined with the new `ErrDecodeJSON` case in `SendErrorResponse`, this resolves the global JSON-parse-to-500 issue for all services using the shared pipeline.
- **`ErrDataSourceAttrParse` partially fixed** -- the unknown-type error path now wraps the sentinel, but the JSON unmarshal path still does not.
- **`Exists` key mismatch fixed** (unstaged) -- `data_sources_repository.go:Exists` now uses the composite key `getDataSourceKey(tenantID, id)` instead of `string(id)`.

### Current State

**Forms Service** is the most mature of the three. With JSON parse error handling and command validation now resolved, the primary gaps are: the hardcoded `"placeholder"` user IDs in publish/retire, the aggregate boundary ambiguity between `Form` and `Version`, the `FieldResponse` DTO silently dropping attributes, and the continued absence of domain validation and test coverage. The hexagonal dependency violation in `core.go` and the ISP violation in the service layer remain the top architectural concerns.

**Tenants Service** has resolved the critical `removeDataSource` route bug, the JSON-parse-to-500 pipeline, and the `Exists` key mismatch. The `DataSource` orphan problem (no tenant existence check), incomplete attribute types (missing json tags, empty `ScheduledDataSourceAttributes`), the `parseAttributes` error wrapping gap, and the stray semicolon remain.

**Submissions Service** has not changed since the prior review. It remains non-functional end-to-end: the entry point is a stub, the service is never wired, 4 routes have `nil` handlers, and the `Database` nil reference will panic on shutdown. No write operations, no DTOs, no middleware, and no tests exist.

**Shared Package** (`pkg/common`) resolved the `json.Decoder` error mapping issue by wrapping decode errors in `ReadJsonPayload` with `ErrDecodeJSON`. The remaining issues are: missing `ErrInvalidID`/`ErrUnauthorized` HTTP mappings, the unused `headers` parameter in `SendJsonResponse`, the ignored `w.Write` error, and the no-op in-memory transactions.

### New Issues Found

1. **`FieldResponse` DTO omits `Attributes`** (Forms #17) -- Field attribute data is silently dropped in API responses.
2. **`ReplaySubmissionCommand` is an empty struct** (Submissions #37) -- Cannot specify replay parameters.
3. **`time.Now()` also in forms service layer** (Forms #8, expanded) -- `PublishVersion` and `RetireVersion` call `time.Now()` directly in business logic, not just the repository layer.

### Highest-Impact Improvements

1. **Wire the submissions service** -- fix the entry point, bootstrap wiring, nil handlers, and nil database (P0 -- service is non-functional)
2. **Add `ErrInvalidID` and `ErrUnauthorized` mappings** to `SendErrorResponse` (P1 -- auth/validation errors produce 500s)
3. **Wrap `parseAttributes` errors with `ErrDataSourceAttrParse`** in `dto/request.go` (P1 -- malformed attribute JSON produces 500)
4. **Fix the hexagonal dependency violation** in `core.go` across all three services (P2 -- foundational architecture)
5. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
