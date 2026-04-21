# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/19 Review

1. ~~`ConditionalRule` is an empty stub~~ (Forms #12, P3) -- `conditional_rule.go` has been replaced by a full `Rule` domain object in `rule.go`. `Rule` has `ID`, `Type` (`visible`, `required`, `readonly`), and `Expression` fields. A `baseWithRules` mixin struct with a `Rules map[RuleType]*Rule` field and `SetRules()` method is embedded by `Page`, `Section`, and `Field`. New `dto/rule.go` provides `RuleRequest`, `RuleResponse`, and mappers. `SetRules()` prevents duplicate rule types via `ErrDuplicateRuleType`.

2. ~~`FieldResponse` DTO omits `Attributes`~~ (Forms #13, P3) -- `FieldResponse` now includes both `Attributes any` and `Rules []*RuleResponse` fields (`dto/field.go:20-21`). `FieldToResponse()` maps `field.Attributes` and `field.Rules` to the response.

3. ~~`DataSource` can be created without verifying its parent `Tenant` exists~~ (Tenants #32, P2) -- `DataSourcesService` now has a private `tenantExists()` helper (`data_sources_service.go:162-173`) that calls `s.tenantsRepository.Exists()`. The helper is invoked in `Create()` (line 60), `Update()` (line 96), and `Remove()` (line 128) before proceeding with the operation.

4. ~~No tenant middleware~~ (Submissions #16, P2) -- `routes.go` now applies `tenants.TenantMiddleware("X-Tenant-ID")` to the submissions router (line 15), consistent with the forms service.

8. ~~Tenant extraction performed in service layer via `baseService.getTenantFromContext()`~~ (All services) -- Tenant extraction has been moved from the service layer to the handler layer across all three services, following hexagonal architecture. Handlers now call `getTenantFromContext(r)` and pass the tenant ID into commands and queries. `base_service.go` has been deleted from forms, submissions, and tenants services. Specifically: all 10 forms handlers, the `getSubmissionByReferenceID` submissions handler, and all 6 tenants data-source handlers now extract the tenant and pass it via typed command/query structs. New types added: `RemoveDataSourceCommand`, `GetDataSourceLookupsCommand`, `FindDataSourceByIDQuery`. The `Lookup` method signature changed from `Lookup(context.Context, domain.DataSourceID)` to `Lookup(context.Context, *GetDataSourceLookupsCommand)`. `FindById` changed from `FindById(context.Context, domain.DataSourceID)` to `FindById(context.Context, *FindDataSourceByIDQuery)`. `Remove` changed from `Remove(context.Context, domain.DataSourceID)` to `Remove(context.Context, *RemoveDataSourceCommand)`. This resolves the prior issue of `getSubmissionByReferenceID` passing an empty TenantID (Submissions #16, P1). *(Unstaged.)*

9. ~~`FindByIdQuery` fields have no `validate` tags~~ (Submissions #19, P3) -- `FindByIdQuery` now carries a `TenantID string` field with `validate:"required"` alongside the existing `ID T` field with `validate:"required"`. The `NewFindByIdQuery` constructor accepts `tenantID` as its first parameter and assigns it to the struct. Validation is performed in the service layer via `validate.ValidateStruct(query)`. *(Unstaged.)*

10. ~~`SendErrorResponse` missing mappings for `ErrInvalidID` and `ErrUnauthorized`~~ (Shared #46, P1) -- `ErrInvalidID` has been added to the 400 Bad Request case alongside `ErrDecodeJSON` and validation errors. `ErrUnauthorized` has a new dedicated case returning 401 Unauthorized. Both the HTTP wire status and the response body `StatusCode` field are correct (`httputil/http.go:62-73`). *(Unstaged.)*

11. ~~`NewFindByIdQuery` creates `validator.New()` per call~~ (Submissions #17, P3) -- The constructor has been simplified to a plain struct builder that returns `*FindByIdQuery[T]` (no error). The `validator` import has been removed entirely. The `ID` field now has a `validate:"required"` tag, and validation is performed in the service layer via `validate.ValidateStruct(query)`, consistent with forms and tenants. *(Unstaged.)*

12. ~~`SendJsonResponse` accepts `headers` but never applies them~~ (Shared #43, P3) -- The function now iterates over the `headers` variadic parameter and sets each header key-value pair on the response (`httputil/http.go:45-51`). *(Unstaged.)*

13. ~~`w.Write` error ignored~~ (Shared #44, P3) -- The return value from `w.Write(out)` is now captured and returned from the function (`httputil/http.go:56`). Handlers that care about write failures can now handle the error. *(Unstaged.)*

14. ~~`IsValidationErr` uses type assertion instead of `errors.As`~~ (Shared #45, P1) -- Changed from `_, ok := err.(validator.ValidationErrors)` to the idiomatic `errors.As(err, &validator.ValidationErrors{})` pattern. This correctly handles wrapped validation errors (`httputil/http.go:10`). *(Unstaged.)*

15. ~~`SendErrorResponse` discards `SendJsonResponse` return errors~~ (Shared #46, P3) -- The function now returns `error` (signature changed from `void`). Every case branch uses `return SendJsonResponse(...)`, propagating write failures to callers. *(Unstaged.)*

16. ~~`ValidateStruct` redundant pattern~~ (Shared #47, P3) -- Simplified from `if err := v.Struct(s); err != nil { return err }; return nil` to `return v.Struct(s)` (`validate/validate.go:13`). *(Unstaged.)*

17. ~~`IsValidationErr` redundant boolean pattern~~ (Shared #48, P3) -- Simplified from multi-line `if/return` to single `return errors.As(...)` expression (`validate/validate.go:10`). *(Unstaged.)*

18. ~~`sendErrorResponse` method is dead code~~ (Submissions #20, P3) -- The `sendErrorResponse` method on the submissions handlers is now called in `getSubmissions` (line 44), `getSubmissionByReferenceID` (lines 60, 79). Previously it existed but was never invoked. *(Unstaged.)*

5. ~~Inconsistent multi-tenancy approach~~ (Tenants #35) -- The `/data-sources` route group now uses `tenants.TenantMiddleware("X-Tenant-ID")` (`routes.go:28`). The `/tenants` route group intentionally has no tenant-scoping middleware since tenant CRUD operations are administrative and not scoped to a single tenant (moved to Will Not Fix).

6. ~~`r.PathValue()` vs `chi.URLParam()` router mismatch~~ (All services) -- All three services now use `chi.URLParam()` instead of the stdlib `r.PathValue()` to extract path parameters. Forms: `getFormIdPathValue`, `getVersionIdPathValue`. Submissions: `getReferenceIdPathValue`. Tenants: `getDataSource`, `updateDataSource`, `removeDataSource`, `getDataSourceLookup`, `getTenantIDPathValue`. *(Unstaged.)*

7. ~~`form_service.go` naming inconsistency~~ (Forms) -- Renamed to `forms_service.go` to match the `FormsService` type name and the naming convention used by the other services (`tenants_service.go`, `data_sources_service.go`, `submissions_service.go`). Content is unchanged. *(Unstaged.)*

---

## Will Not Fix

These issues have been reviewed and accepted as intentional design decisions. They should not be flagged in future reviews.

1. **Goroutine/channel pattern in handlers** (Previously Forms #7, Submissions #27, Tenants #46) -- The `go func() -> chan -> select { case <-r.Context().Done(); case res := <-resultChan }` pattern in every handler is the recommended approach for respecting chi router's context-based request timeouts. Without this pattern, a handler performing a long-running service call would not be able to short-circuit when the request context is cancelled (e.g., client disconnect, server timeout). The `select` on `r.Context().Done()` enables cooperative cancellation at the handler level. The allocation overhead of a single goroutine and buffered channel per request is negligible relative to the I/O cost of a real database call.

2. **In-memory transactions are no-ops** (Previously Forms #17, Shared #57) -- `BeginTx`/`CommitTx`/`RollbackTx` in `inmemory_database.go` do nothing by design. The in-memory database is intended for local development and testing only; atomicity guarantees are not required in this context.

3. **Inconsistent response envelope** (Previously Tenants #51) -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. This is an intentional convention: GET list operations return the collection directly, while CUD operations return the response envelope.

4. **REST handlers hold a reference to the full `Application`** (Previously Forms #2, Submissions #16, Tenants #32) -- The `handlers` struct takes `*core.Application` rather than narrowed dependencies. `Application` acts as a dependency container assembled at the composition root (`main.go`) that groups the application's top-level dependencies. It exports only `Logger` and `Services` (the `repository` field is unexported and inaccessible to the adapter layer). Passing the container directly avoids cascading signature changes through `newHandlers` and `Routes` when new cross-cutting concerns (e.g., config, metrics) are added to `Application`. The surface area exposed to handlers is already minimal.

5. **Tenants route group has no tenant-scoping middleware** (Previously Tenants #35) -- The `/tenants` route group in the tenants service does not use `TenantMiddleware`. This is intentional: tenant CRUD operations are administrative endpoints that manage tenants themselves, so scoping them to a single tenant via `X-Tenant-ID` is not applicable. The `/data-sources` route group correctly uses tenant middleware because data sources belong to a specific tenant.

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

#### Architectural

16. **`Find()` has no tenant filtering** (`submissions_service.go:27-29`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18.)*

19. **Four handler stubs return 200 OK with empty body** (`handlers.go:84-91`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header, which is misleading to clients. These should either return 501 Not Implemented or not be registered. *(New.)*

#### Missing Functionality

20. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #30, 4/18 #29, 4/19 #21.)*

21. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:68-74`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22.)*

22. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23.)*

23. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24.)*

24. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25.)*

25. **No domain events** for cross-service communication. *(Unresolved from 4/17 #35, 4/18 #34, 4/19 #26.)*

26. **No real authentication.** *(Unresolved from 4/17 #36, 4/18 #35, 4/19 #27.)*

27. **`ReplaySubmissionCommand` is an empty struct** (`commands.go:3`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28.)*

#### Code Quality

28. **`Payload` typed as `any`** (`submission.go:18`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29.)*

29. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30.)*

30. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go:51-61`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31.)*

---

### Tenants Service

#### Architectural

31. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34.)*

32. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned in the data sources store. There is no cascade delete and no service-level cleanup. *(New.)*

#### Missing Functionality

33. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12, 4/17 #50, 4/18 #48, 4/19 #36.)*

34. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` has some validation (attribute type matching) but does not validate empty `TenantID`, empty `Type`, or field lengths. *(Unresolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37.)*

35. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14, 4/17 #52, 4/18 #50, 4/19 #38.)*

36. **Incomplete error-to-HTTP mapping** -- `ErrStrategyNotFound` and `ErrDataSourceAttrParse` fall through to 500. The service-level `sendErrorResponse` only maps `ErrInvalidSourceTypeAttributes`. (`ErrInvalidID` and `ErrUnauthorized` are now handled at the shared level -- see resolved issue #10.) *(Unresolved from 4/16 #15, 4/17 #53, 4/18 #51, 4/19 #39.)*

37. **No real authentication**. *(Unresolved from 4/16 #16, 4/17 #54, 4/18 #52, 4/19 #40.)*

38. **`Lookup` service method is a stub** (`data_sources_service.go:145-160`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41.)*

39. **`DataSourceAttributes` concrete types incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names. *(Unresolved from 4/16 #18, 4/17 #56, 4/18 #54, 4/19 #42.)*

40. **`DataSourceLookup` value object has no validation** (`data_source_lookup.go`) -- `NewDataSourceLookup(code, description)` constructor has been added, but it performs no validation and has no json tags on the struct fields. *(Partially resolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43.)*

41. **`DataSourceType` not validated in domain constructors** -- Command-level `oneof` validation exists, but `NewDataSource` still accepts any arbitrary string for `Type`. *(Unresolved from 4/16 #22, 4/17 #58, 4/18 #56, 4/19 #44.)*

#### Code Quality

42. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:69`). *(Unresolved from 4/16 #24, 4/17 #60, 4/18 #58, 4/19 #45.)*

43. **`ErrInvalidSourceType` defined but unused** (`data_source.go:19`) -- The error variable is declared but never referenced anywhere. Only `ErrInvalidSourceTypeAttributes` is used. *(New.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 2, 16 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 7, 34 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 23 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 13 | `ErrMissingTenantID` maps to 500 | Forms |
| **P2** | 32 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 19 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 5 | Redundant double-fetch in `Update()` | Forms |
| **P3** | 6, 20, 33 | Zero test files | All |
| **P3** | 8, 25, 35 | No domain events | All |
| **P3** | 9, 36 | Incomplete error-to-HTTP mapping | Forms, Tenants |
| **P3** | 28 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4, 42 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 15 | Map iteration non-deterministic in DTO mappers | Forms |
| **P3** | 14 | `publishVersion`/`retireVersion` discard returned version | Forms |


---

## Summary

### Progress Since 4/19

Nineteen issues from the prior review have been resolved:

- **`ConditionalRule` replaced with full `Rule` domain object** (Forms) -- The empty stub has been replaced with a proper `Rule` entity supporting three rule types (`visible`, `required`, `readonly`) and an `Expression` field. A `baseWithRules` mixin is embedded by `Page`, `Section`, and `Field`, with `SetRules()` enforcing duplicate-type prevention via `ErrDuplicateRuleType`. New DTOs (`RuleRequest`, `RuleResponse`) and mappers support rules throughout the request/response pipeline.
- **`FieldResponse` DTO now includes `Attributes` and `Rules`** (Forms) -- Field attribute data and rules are no longer silently dropped in API responses.
- **`DataSource` parent tenant existence check added** (Tenants) -- `DataSourcesService` now verifies the parent tenant exists before Create, Update, and Remove operations via a `tenantExists()` helper, preventing orphaned data sources.
- **Tenant middleware applied to submissions service** (Submissions) -- `routes.go` now uses `tenants.TenantMiddleware("X-Tenant-ID")`, consistent with the forms service.
- **Tenant middleware applied to tenants service data-sources routes** (Tenants) -- The `/data-sources` route group now uses the shared tenant middleware. The `/tenants` route group intentionally has no tenant middleware (moved to Will Not Fix).
- **`r.PathValue()` replaced with `chi.URLParam()` across all services** (All, unstaged) -- All path parameter extraction methods now use chi's `URLParam()` function, fixing a potential silent empty-string bug when using chi's router.
- **`form_service.go` renamed to `forms_service.go`** (Forms, unstaged) -- Naming now matches the `FormsService` type and the convention used by sibling services.
- **Hexagonal architecture: tenant extraction moved from service layer to handler layer** (All services, unstaged) -- `base_service.go` (containing `getTenantFromContext()`) has been deleted from all three services. Handlers now extract the tenant via `getTenantFromContext(r)` and pass it into commands/queries. All forms handlers (10), the submissions `getSubmissionByReferenceID` handler, and all tenants data-source handlers (6) have been updated. New types: `RemoveDataSourceCommand`, `GetDataSourceLookupsCommand`, `FindDataSourceByIDQuery`. The `DataSourcesService` interface methods `FindById`, `Remove`, and `Lookup` now accept typed command/query structs instead of bare `domain.DataSourceID`. This also resolves the prior `getSubmissionByReferenceID` empty TenantID bug (Submissions #16, P1).
- **`FindByIdQuery` now has validated `TenantID` field** (Submissions, unstaged) -- `FindByIdQuery` carries `TenantID string` with `validate:"required"`. The `NewFindByIdQuery` constructor accepts and assigns `tenantID`. Validation is performed in the service layer via `validate.ValidateStruct(query)`.
- **`SendErrorResponse` now maps `ErrInvalidID` and `ErrUnauthorized`** (Shared, unstaged) -- `ErrInvalidID` is mapped to 400 Bad Request alongside `ErrDecodeJSON` and validation errors. `ErrUnauthorized` has a new dedicated case returning 401 Unauthorized. This was the longest-standing P1 in the codebase (first identified 4/17).
- **`NewFindByIdQuery` no longer creates `validator.New()` per call** (Submissions, unstaged) -- The constructor is now a plain struct builder returning `*FindByIdQuery[T]` (no error). The `validator` import has been removed. Validation is performed in the service layer.
- **`SendJsonResponse` now applies custom headers** (Shared, unstaged) -- The function iterates over the `headers` variadic parameter and sets each header on the response.
- **`SendJsonResponse` returns write errors** (Shared, unstaged) -- The return value from `w.Write(out)` is now captured and returned.
- **`IsValidationErr` uses `errors.As`** (Shared, unstaged) -- Changed from direct type assertion to idiomatic `errors.As()`, correctly handling wrapped validation errors.
- **`SendErrorResponse` returns errors** (Shared, unstaged) -- The function signature changed from `void` to `error`, propagating write failures to callers.
- **`ValidateStruct` simplified** (Shared, unstaged) -- Redundant `if/return` pattern removed.
- **`IsValidationErr` simplified** (Shared, unstaged) -- Redundant boolean pattern removed.
- **`sendErrorResponse` no longer dead code** (Submissions, unstaged) -- The method is now called in `getSubmissions` and `getSubmissionByReferenceID` handlers instead of calling `httputil.SendErrorResponse` directly.
- **`DataSourceLookup` now has a constructor** (Tenants, unstaged) -- `NewDataSourceLookup(code, description)` has been added. Struct still lacks json tags and validation.

### New Issues Found

1. **Tenant removal doesn't cascade-delete DataSources** (Tenants #32, P2) -- Deleting a tenant leaves its data sources orphaned.
2. **Empty handler stubs return 200 OK** (Submissions #19, P2) -- Four unimplemented endpoints are reachable and return misleading 200 responses.
3. **Redundant double-fetch in forms `Update()`** (Forms #5, P2) -- `isValidAccess()` and the subsequent `FindById()` fetch the same form twice.
4. **Map iteration non-deterministic in DTO mappers** (Forms #15, P3) -- Pages, sections, fields, and rules appear in random order in responses.
5. **`publishVersion`/`retireVersion` discard returned version** (Forms #14, P3) -- Clients receive `nil` data despite a successful operation.
6. **`ErrInvalidSourceType` defined but unused** (Tenants #43, P3) -- Error variable declared but never referenced.

### Current State

**Forms Service** remains the most mature. The `Rule` domain object, `FieldResponse` DTO fixes, and the hexagonal architecture refactor (tenant extraction moved to handlers, passed via commands/queries) are significant improvements. The primary remaining gaps are: the hardcoded `"placeholder"` user IDs, `Find()` with no tenant filtering, the aggregate boundary ambiguity between `Form` and `Version`, the redundant double-fetch in `Update()`, the non-deterministic map iteration in response DTOs, and the continued absence of domain validation and test coverage.

**Tenants Service** has made meaningful progress with the parent tenant existence check, tenant middleware on data-source routes, and the hexagonal architecture refactor. All data-source operations now receive tenant via typed commands/queries, and `Lookup`, `FindById`, and `Remove` have been updated to accept structured types instead of bare IDs. The remaining gaps are: the cascade-delete problem when removing tenants, incomplete attribute types (missing json tags, empty `ScheduledDataSourceAttributes`), the stub `Lookup` method, unused `ErrInvalidSourceType`, and incomplete error-to-HTTP mapping.

**Submissions Service** has made notable progress. Tenant middleware is wired at the router level, and the hexagonal architecture refactor moves tenant extraction to the handler layer with tenant passed via `FindByIdQuery`. The `sendErrorResponse` method is now actively used. Four handler stubs still return misleading 200 OK responses. `Find()` still has no tenant filtering. The service still lacks write operations, domain constructors, request DTOs, and test coverage.

**Shared Package** (`pkg/common`) remains issue-free. All prior issues have been resolved.

### Highest-Impact Improvements

1. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
2. **Add tenant filtering to `Find()` methods** across forms and submissions (P2 -- data isolation)
3. **Add test coverage** starting with service and handler layers (P3 -- long-term reliability)
