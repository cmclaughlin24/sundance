# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/19 Review

1. ~~`ConditionalRule` is an empty stub~~ (Forms #12, P3) -- `conditional_rule.go` has been replaced by a full `Rule` domain object in `rule.go`. `Rule` has `ID`, `Type` (`visible`, `required`, `readonly`), and `Expression` fields. A `baseWithRules` mixin struct with a `Rules map[RuleType]*Rule` field and `SetRules()` method is embedded by `Page`, `Section`, and `Field`. New `dto/rule.go` provides `RuleRequest`, `RuleResponse`, and mappers. `SetRules()` prevents duplicate rule types via `ErrDuplicateRuleType`.

2. ~~`FieldResponse` DTO omits `Attributes`~~ (Forms #13, P3) -- `FieldResponse` now includes both `Attributes any` and `Rules []*RuleResponse` fields (`dto/field.go:20-21`). `FieldToResponse()` maps `field.Attributes` and `field.Rules` to the response.

3. ~~`DataSource` can be created without verifying its parent `Tenant` exists~~ (Tenants #32, P2) -- `DataSourcesService` now has a private `tenantExists()` helper (`data_sources_service.go:162-173`) that calls `s.tenantsRepository.Exists()`. The helper is invoked in `Create()` (line 60), `Update()` (line 96), and `Remove()` (line 128) before proceeding with the operation.

4. ~~No tenant middleware~~ (Submissions #16, P2) -- `routes.go` now applies `tenants.TenantMiddleware("X-Tenant-ID")` to the submissions router (line 15), consistent with the forms service. *(Note: the `getSubmissionByReferenceID` handler still passes an empty string for TenantID instead of extracting it from context -- see new issue #17.)*

5. ~~Inconsistent multi-tenancy approach~~ (Tenants #35, partially) -- The `/data-sources` route group now uses `tenants.TenantMiddleware("X-Tenant-ID")` (`routes.go:28`). The `/tenants` route group still has no tenant-scoping middleware. This is a partial fix -- see remaining issue #32 for the residual gap.

6. ~~`r.PathValue()` vs `chi.URLParam()` router mismatch~~ (All services) -- All three services now use `chi.URLParam()` instead of the stdlib `r.PathValue()` to extract path parameters. Forms: `getFormIdPathValue`, `getVersionIdPathValue`. Submissions: `getReferenceIdPathValue`. Tenants: `getDataSource`, `updateDataSource`, `removeDataSource`, `getDataSourceLookup`, `getTenantIDPathValue`. *(Unstaged.)*

7. ~~`form_service.go` naming inconsistency~~ (Forms) -- Renamed to `forms_service.go` to match the `FormsService` type name and the naming convention used by the other services (`tenants_service.go`, `data_sources_service.go`, `submissions_service.go`). Content is unchanged. *(Unstaged.)*

---

## Will Not Fix

These issues have been reviewed and accepted as intentional design decisions. They should not be flagged in future reviews.

1. **Goroutine/channel pattern in handlers** (Previously Forms #7, Submissions #27, Tenants #46) -- The `go func() -> chan -> select { case <-r.Context().Done(); case res := <-resultChan }` pattern in every handler is the recommended approach for respecting chi router's context-based request timeouts. Without this pattern, a handler performing a long-running service call would not be able to short-circuit when the request context is cancelled (e.g., client disconnect, server timeout). The `select` on `r.Context().Done()` enables cooperative cancellation at the handler level. The allocation overhead of a single goroutine and buffered channel per request is negligible relative to the I/O cost of a real database call.

2. **In-memory transactions are no-ops** (Previously Forms #17, Shared #57) -- `BeginTx`/`CommitTx`/`RollbackTx` in `inmemory_database.go` do nothing by design. The in-memory database is intended for local development and testing only; atomicity guarantees are not required in this context.

3. **Inconsistent response envelope** (Previously Tenants #51) -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. This is an intentional convention: GET list operations return the collection directly, while CUD operations return the response envelope.

4. **REST handlers hold a reference to the full `Application`** (Previously Forms #2, Submissions #16, Tenants #32) -- The `handlers` struct takes `*core.Application` rather than narrowed dependencies. `Application` acts as a dependency container assembled at the composition root (`main.go`) that groups the application's top-level dependencies. It exports only `Logger` and `Services` (the `repository` field is unexported and inaccessible to the adapter layer). Passing the container directly avoids cascading signature changes through `newHandlers` and `Routes` when new cross-cutting concerns (e.g., config, metrics) are added to `Application`. The surface area exposed to handlers is already minimal.

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:281`, `310`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1, 4/19 #1.)*

#### Architectural

2. **`Find()` has no tenant filtering** (`forms_service.go:30-31`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5, 4/19 #2.)*

3. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6, 4/19 #4.)*

4. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `forms_service.go:249`, `282`) -- Repository calls `time.Now()` for `CreatedAt`/`UpdatedAt` timestamps. Service calls `time.Now()` when passing time to `Publish()` and `Retire()`. Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8, 4/19 #5.)*

5. **Redundant double-fetch in `Update()`** (`forms_service.go:92-96`) -- `isValidAccess()` fetches the form from the repository to check tenant ownership, then `Update()` immediately fetches the same form again via `FindById()`. The form should be fetched once and reused. *(New.)*

#### Missing Functionality

6. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10, 4/19 #6.)*

7. **Domain validation unimplemented** (`form.go:29,42`, `version.go:48`, `page.go:31`, `section.go:31`, `field.go:35`, `rule.go:32`) -- All entity constructors contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13, 4/17 #10, 4/18 #11, 4/19 #7.)*

8. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12, 4/19 #8.)*

9. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and domain errors (`ErrVersionLocked`, `ErrDuplicatePosition`, `ErrDuplicateRuleType`, `ErrPublishedByRequired`, `ErrRetiredByRequired`) all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` (`handlers.go:344-349`) is an empty switch that delegates directly to `httputil.SendErrorResponse`. *(Unresolved from 4/13 #15, 4/17 #12, 4/18 #13, 4/19 #9.)*

10. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14, 4/19 #10.)*

11. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15, 4/19 #11.)*

#### Code Quality

12. **Inconsistent constructor signatures** -- Forms domain constructors return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. *(Unresolved from 4/13 #23, 4/17 #18, 4/18 #19, 4/19 #14.)*

13. **`ErrMissingTenantID` maps to 500** (`middleware.go:24`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, it falls through to the 500 default. Should be 400. *(Unresolved from 4/19 #15.)*

14. **`publishVersion` and `retireVersion` discard the returned `*domain.Version`** (`handlers.go:287`, `320`) -- Both handlers call the service, receive an updated `*domain.Version`, but respond with `Data: nil`. The mutated version state is never returned to the client. *(New.)*

15. **Map iteration order non-deterministic in DTO response mappers** (`dto/version.go`, `dto/page.go`, `dto/section.go`, `dto/rule.go`) -- `VersionToResponse()`, `PageToResponse()`, `SectionToResponse()`, and `RuleToResponse()` iterate over Go maps (`map[int]*Page`, `map[int]*Section`, `map[int]*Field`, `map[RuleType]*Rule`). The order of items in JSON array responses will be non-deterministic across requests. Pages, sections, and fields should be sorted by position key. *(New.)*

---

### Submissions Service

#### Bugs

16. **`getSubmissionByReferenceID` passes empty string for TenantID** (`handlers.go:59`) -- The handler calls `ports.NewFindByIdQuery(referenceID, "")` with a hardcoded empty string. The middleware injects the tenant into context, but the handler never calls `tenants.TenantFromContext(r.Context())` to retrieve it. The service layer check `submission.TenantID != query.TenantID` will always fail for submissions with a real tenant, returning `ErrUnauthorized`. *(New.)*

#### Architectural

17. **`Find()` has no tenant filtering** (`submissions_service.go:24-26`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18.)*

18. **`NewFindByIdQuery` creates `validator.New()` per call** (`ports/query.go:16`) -- Inconsistent with tenants and forms which use the shared `validate.ValidateStruct()` singleton from `pkg/common/validate`. *(Unresolved from 4/18 #28, 4/19 #19.)*

19. **`FindByIdQuery` fields have no `validate` tags** (`ports/query.go:6-8`) -- Neither `ID` nor `TenantID` has `validate:"required"` tags. The validation call in `NewFindByIdQuery` will always pass regardless of input, making it pure overhead with no benefit. *(New.)*

20. **`sendErrorResponse` method is dead code** (`handlers.go:97-101`) -- The method exists but is never called. All handlers call `httputil.SendErrorResponse` directly. *(Unresolved from 4/19 #20.)*

21. **Four handler stubs return 200 OK with empty body** (`handlers.go:84-91`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header, which is misleading to clients. These should either return 501 Not Implemented or not be registered. *(New.)*

#### Missing Functionality

22. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #30, 4/18 #29, 4/19 #21.)*

23. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:56-62`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22.)*

24. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23.)*

25. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24.)*

26. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25.)*

27. **No domain events** for cross-service communication. *(Unresolved from 4/17 #35, 4/18 #34, 4/19 #26.)*

28. **No real authentication.** *(Unresolved from 4/17 #36, 4/18 #35, 4/19 #27.)*

29. **`ReplaySubmissionCommand` is an empty struct** (`commands.go:3`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28.)*

#### Code Quality

30. **`Payload` typed as `any`** (`submission.go:18`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29.)*

31. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30.)*

32. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go:51-61`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31.)*

---

### Tenants Service

#### Architectural

33. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34.)*

34. **Tenants route group has no tenant-scoping middleware** -- The `/data-sources` route group now uses `tenants.TenantMiddleware`, but the `/tenants` route group still has no middleware. This creates inconsistency: data source operations are tenant-scoped, but tenant CRUD operations are open. *(Partially unresolved from 4/18 #47, 4/19 #35.)*

35. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned in the data sources store. There is no cascade delete and no service-level cleanup. *(New.)*

#### Missing Functionality

36. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12, 4/17 #50, 4/18 #48, 4/19 #36.)*

37. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` has some validation (attribute type matching) but does not validate empty `TenantID`, empty `Type`, or field lengths. *(Unresolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37.)*

38. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14, 4/17 #52, 4/18 #50, 4/19 #38.)*

39. **Incomplete error-to-HTTP mapping** -- `ErrInvalidID` and `ErrUnauthorized` fall through to 500. `ErrStrategyNotFound` and `ErrDataSourceAttrParse` also fall through to 500. The service-level `sendErrorResponse` only maps `ErrInvalidSourceTypeAttributes`. *(Unresolved from 4/16 #15, 4/17 #53, 4/18 #51, 4/19 #39.)*

40. **No real authentication**. *(Unresolved from 4/16 #16, 4/17 #54, 4/18 #52, 4/19 #40.)*

41. **`Lookup` service method is a stub** (`data_sources_service.go:145-160`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41.)*

42. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names. *(Unresolved from 4/16 #18, 4/17 #56, 4/18 #54, 4/19 #42.)*

43. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct with two fields, no json tags, no invariant enforcement. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43.)*

44. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string for `Type`. *(Unresolved from 4/16 #22, 4/17 #58, 4/18 #56, 4/19 #44.)*

#### Code Quality

45. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:69`). *(Unresolved from 4/16 #24, 4/17 #60, 4/18 #58, 4/19 #45.)*

---

### Shared Package (`pkg/common`)

#### Bugs

46. **`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`** (`httputil/http.go:56-87`) -- These sentinel errors are defined in `error.go` but not handled in the switch. They fall through to the 500 default. `ErrInvalidID` should map to 400; `ErrUnauthorized` should map to 401 or 403. *(Unresolved from 4/17 #64, 4/18 #60, 4/19 #46.)*

47. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`httputil/http.go:41`) -- The `headers ...http.Header` variadic parameter is accepted but the body never iterates or sets them on the response. *(Unresolved from 4/17 #65, 4/18 #61, 4/19 #47.)*

48. **`w.Write` error ignored** (`httputil/http.go:50`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. *(Unresolved from 4/17 #66, 4/18 #62, 4/19 #48.)*

49. **`IsValidationErr` uses type assertion instead of `errors.As`** (`validate/validate.go:10`) -- `_, ok := err.(validator.ValidationErrors)` is a direct type assertion that will fail if the `ValidationErrors` is wrapped inside another error (e.g., via `fmt.Errorf("%w", err)`). The idiomatic Go approach is `errors.As(err, &ve)`. This matters because `IsValidationErr` is called by `SendErrorResponse` to decide the 400 mapping -- if any layer ever wraps a validation error, it would fall through to 500. *(New.)*

50. **`SendErrorResponse` discards `SendJsonResponse` return errors** (`httputil/http.go:63,69,75,81`) -- Every call to `SendJsonResponse` inside `SendErrorResponse` ignores the returned `error`. If JSON marshaling or writing fails, the error is silently swallowed. Since `SendErrorResponse` returns nothing, callers have no way to know the error response failed to send. *(New.)*

#### Code Quality

51. **`ValidateStruct` has a redundant pattern** (`validate/validate.go:19-25`) -- `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. *(Unresolved from 4/17 #67, 4/18 #63, 4/19 #49.)*

52. **`IsValidationErr` has the same redundant boolean pattern** (`validate/validate.go:9-17`) -- `if !ok { return false }; return true` is equivalent to `return ok`. *(New.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 46 | `ErrInvalidID`/`ErrUnauthorized` map to 500 | Shared |
| **P1** | 16 | `getSubmissionByReferenceID` passes empty TenantID | Submissions |
| **P1** | 49 | `IsValidationErr` uses type assertion instead of `errors.As` | Shared |
| **P2** | 2, 17 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 7, 37 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 25 | No domain constructors | Submissions |
| **P2** | 34 | Tenants route group has no middleware | Tenants |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 13 | `ErrMissingTenantID` maps to 500 | Forms |
| **P2** | 35 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 21 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 5 | Redundant double-fetch in `Update()` | Forms |
| **P3** | 6, 22, 36 | Zero test files | All |
| **P3** | 8, 27, 38 | No domain events | All |
| **P3** | 9, 39 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 30 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4, 45 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 15 | Map iteration non-deterministic in DTO mappers | Forms |
| **P3** | 14 | `publishVersion`/`retireVersion` discard returned version | Forms |
| **P3** | 50 | `SendErrorResponse` discards `SendJsonResponse` errors | Shared |
| **P3** | 51, 52 | Redundant validation patterns | Shared |
| **P3** | 18, 19 | `NewFindByIdQuery` validator issues | Submissions |
| **P3** | 20 | `sendErrorResponse` dead code | Submissions |

---

## Summary

### Progress Since 4/19

Seven issues from the prior review have been resolved:

- **`ConditionalRule` replaced with full `Rule` domain object** (Forms) -- The empty stub has been replaced with a proper `Rule` entity supporting three rule types (`visible`, `required`, `readonly`) and an `Expression` field. A `baseWithRules` mixin is embedded by `Page`, `Section`, and `Field`, with `SetRules()` enforcing duplicate-type prevention via `ErrDuplicateRuleType`. New DTOs (`RuleRequest`, `RuleResponse`) and mappers support rules throughout the request/response pipeline.
- **`FieldResponse` DTO now includes `Attributes` and `Rules`** (Forms) -- Field attribute data and rules are no longer silently dropped in API responses.
- **`DataSource` parent tenant existence check added** (Tenants) -- `DataSourcesService` now verifies the parent tenant exists before Create, Update, and Remove operations via a `tenantExists()` helper, preventing orphaned data sources.
- **Tenant middleware applied to submissions service** (Submissions) -- `routes.go` now uses `tenants.TenantMiddleware("X-Tenant-ID")`, consistent with the forms service. However, the handler for `getSubmissionByReferenceID` does not yet extract the tenant from context (new issue #16).
- **Tenant middleware partially applied to tenants service** (Tenants) -- The `/data-sources` route group now uses the shared tenant middleware. The `/tenants` route group remains without middleware.
- **`r.PathValue()` replaced with `chi.URLParam()` across all services** (All, unstaged) -- All path parameter extraction methods now use chi's `URLParam()` function, fixing a potential silent empty-string bug when using chi's router.
- **`form_service.go` renamed to `forms_service.go`** (Forms, unstaged) -- Naming now matches the `FormsService` type and the convention used by sibling services.

### New Issues Found

1. **`getSubmissionByReferenceID` passes empty TenantID** (Submissions #16, P1) -- The handler hardcodes `""` as the tenant ID instead of extracting it from context. The tenant middleware is wired but its value is never consumed, causing the tenant ownership check to always fail.
2. **`IsValidationErr` uses type assertion instead of `errors.As`** (Shared #49, P1) -- A wrapped validation error would bypass the 400 mapping and produce a 500.
3. **Tenant removal doesn't cascade-delete DataSources** (Tenants #35, P2) -- Deleting a tenant leaves its data sources orphaned.
4. **Empty handler stubs return 200 OK** (Submissions #21, P2) -- Four unimplemented endpoints are reachable and return misleading 200 responses.
5. **Redundant double-fetch in forms `Update()`** (Forms #5, P2) -- `isValidAccess()` and the subsequent `FindById()` fetch the same form twice.
6. **`FindByIdQuery` fields have no validate tags** (Submissions #19, P3) -- Validation call is a no-op.
7. **Map iteration non-deterministic in DTO mappers** (Forms #15, P3) -- Pages, sections, fields, and rules appear in random order in responses.
8. **`publishVersion`/`retireVersion` discard returned version** (Forms #14, P3) -- Clients receive `nil` data despite a successful operation.
9. **`SendErrorResponse` discards `SendJsonResponse` errors** (Shared #50, P3) -- Failed error responses are silently swallowed.
10. **`IsValidationErr` redundant boolean pattern** (Shared #52, P3) -- Same style issue as `ValidateStruct`.

### Current State

**Forms Service** remains the most mature. The `Rule` domain object and `FieldResponse` DTO fixes are significant additions that bring the forms domain model closer to feature-complete. The chi `URLParam()` fix and file rename (both unstaged) are minor but welcome cleanups. The primary remaining gaps are: the hardcoded `"placeholder"` user IDs, the aggregate boundary ambiguity between `Form` and `Version`, the redundant double-fetch in `Update()`, the non-deterministic map iteration in response DTOs, and the continued absence of domain validation and test coverage.

**Tenants Service** has made meaningful progress with the parent tenant existence check on DataSource operations and tenant middleware on the data-sources route group. The remaining gaps are: the unprotected tenants route group, the cascade-delete problem when removing tenants, incomplete attribute types (missing json tags, empty `ScheduledDataSourceAttributes`), the stub `Lookup` method, and incomplete error-to-HTTP mapping.

**Submissions Service** has the new tenant middleware wired at the router level, but the handler-level integration is broken (`getSubmissionByReferenceID` passes empty tenant ID). Four handler stubs return misleading 200 OK responses. The service still lacks write operations, domain constructors, request DTOs, and test coverage. The `FindByIdQuery` validation is a no-op due to missing struct tags.

**Shared Package** (`pkg/common`) is largely unchanged. The `IsValidationErr` type assertion issue is newly identified as P1 because it affects the correctness of error-to-HTTP mapping for wrapped validation errors. Missing `ErrInvalidID`/`ErrUnauthorized` HTTP mappings remain P1. The unused `headers` parameter in `SendJsonResponse`, the ignored `w.Write` error, and the redundant validation patterns remain.

### Highest-Impact Improvements

1. **Fix `getSubmissionByReferenceID` to extract TenantID from context** (P1 -- endpoint is functionally broken)
2. **Change `IsValidationErr` to use `errors.As`** (P1 -- wrapped validation errors silently produce 500s)
3. **Add `ErrInvalidID` and `ErrUnauthorized` mappings** to `SendErrorResponse` (P1 -- auth/validation errors produce 500s)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
