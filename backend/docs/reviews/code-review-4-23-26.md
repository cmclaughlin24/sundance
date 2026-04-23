# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/22 Review

1. ~~`DataSourceType` not validated in domain constructors~~ (Tenants, P3) -- `NewDataSource` now calls `isValidSourceType()` and returns `ErrInvalidSourceType` when invalid (`data_source.go:44-46`). Was incorrectly flagged as unresolved in the 4/22 review; already fixed at that time.

2. ~~`ErrInvalidSourceType` defined but unused~~ (Tenants, P3) -- Now used in `NewDataSource` (`data_source.go:45`) and checked in `isBadRequest()` (`handlers.go:434`). Was incorrectly flagged as unresolved in the 4/22 review; already fixed at that time.

3. ~~`ErrDataSourceAttrParse` falls through to 500~~ (Tenants, P3) -- New `isBadRequest()` helper (`handlers.go:432-436`) checks `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes`. All three now correctly map to 400 Bad Request.

4. ~~Inconsistent error response delegation in forms handlers~~ (Forms, P3) -- All `getTenantFromContext` error paths (`handlers.go:60,89,129,169,234,275,322,357`) now route through `h.sendErrorResponse` instead of calling `httputil.SendErrorResponse` directly.

5. ~~`NewTypeValidator` captures mutable slice~~ (Shared, P3) -- Now uses `slices.Clone(types)` before closing over the slice (`validate/validate.go:20-23`).

6. ~~Typo "tyep" in doc comment~~ (Shared, P3) -- Fixed to "type" (`httputil/http.go:39`).

7. **Tenants DTO file reorganization** (Tenants) -- `response.go` and `request.go` split into per-entity files: `tenant.go`, `data_source.go`, `data_source_attributes.go`, `data_source_lookup.go`. Purely structural, no functional changes.

8. ~~Incomplete error-to-HTTP mapping in tenants~~ (Tenants, P3) -- `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes` now map to 400 via `isBadRequest()`. `ErrStrategyNotFound` moved to Will Not Fix (see below). *(Resolved from 4/16 #15, 4/17 #53, 4/18 #51, 4/19 #39, 4/20 #34, 4/22 #35.)*

---

## Will Not Fix

These issues have been reviewed and accepted as intentional design decisions. They should not be flagged in future reviews.

1. **Goroutine/channel pattern in handlers** (Previously Forms #7, Submissions #27, Tenants #46) -- The `go func() -> chan -> select { case <-r.Context().Done(); case res := <-resultChan }` pattern in every handler is the recommended approach for respecting chi router's context-based request timeouts. Without this pattern, a handler performing a long-running service call would not be able to short-circuit when the request context is cancelled (e.g., client disconnect, server timeout). The `select` on `r.Context().Done()` enables cooperative cancellation at the handler level. The allocation overhead of a single goroutine and buffered channel per request is negligible relative to the I/O cost of a real database call.

2. **In-memory transactions are no-ops** (Previously Forms #17, Shared #57) -- `BeginTx`/`CommitTx`/`RollbackTx` in `inmemory_database.go` do nothing by design. The in-memory database is intended for local development and testing only; atomicity guarantees are not required in this context.

3. **Inconsistent response envelope** (Previously Tenants #51) -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. This is an intentional convention: GET list operations return the collection directly, while CUD operations return the response envelope.

4. **REST handlers hold a reference to the full `Application`** (Previously Forms #2, Submissions #16, Tenants #32) -- The `handlers` struct takes `*core.Application` rather than narrowed dependencies. `Application` acts as a dependency container assembled at the composition root (`main.go`) that groups the application's top-level dependencies. It exports only `Logger` and `Services` (the `repository` field is unexported and inaccessible to the adapter layer). Passing the container directly avoids cascading signature changes through `newHandlers` and `Routes` when new cross-cutting concerns (e.g., config, metrics) are added to `Application`. The surface area exposed to handlers is already minimal.

5. **Tenants route group has no tenant-scoping middleware** (Previously Tenants #35) -- The `/tenants` route group in the tenants service does not use `TenantMiddleware`. This is intentional: tenant CRUD operations are administrative endpoints that manage tenants themselves, so scoping them to a single tenant via `X-Tenant-ID` is not applicable. The `/data-sources` route group correctly uses tenant middleware because data sources belong to a specific tenant.

6. **`getTenantFromContext` return types differ across services** -- Forms and submissions return `(string, error)`, while tenants returns `(domain.TenantID, error)`. This is intentional: the tenants service owns the `TenantID` domain type and should use it; the other services do not own that type and correctly use `string`.

7. **`ErrStrategyNotFound` falls through to 500** (Previously Tenants #35) -- `ErrStrategyNotFound` indicates a developer failed to register a strategy for a data source type, which is an internal programming error, not a client error. A 500 is the correct response. This should not be flagged in future reviews.

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:330`, `365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1, 4/19 #1, 4/20 #1, 4/22 #1.)*

#### Architectural

2. **`Find()` has no tenant filtering** (`forms_service.go:29-31`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. The `getForms` handler (`handlers.go:30-55`) is also the only handler that never calls `getTenantFromContext`, meaning tenant context is never even read on this path. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5, 4/19 #2, 4/20 #2, 4/22 #2.)*

3. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6, 4/19 #4, 4/20 #3, 4/22 #3.)*

4. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `forms_service.go:210`, `238`) -- Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8, 4/19 #5, 4/20 #4, 4/22 #4.)*

5. **Redundant double-fetch in `Update()`** (`forms_service.go:78-86`) -- `isValidAccess()` fetches the form from the repository to check tenant ownership, then `Update()` immediately fetches the same form again via `FindById()`. The form should be fetched once and reused. *(Unresolved from 4/20 #5, 4/22 #5.)*

#### Missing Functionality

6. **Zero test files** in the entire forms service. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10, 4/19 #6, 4/20 #6, 4/22 #6.)*

7. **Domain validation partially unimplemented** (`form.go:29,42`, `version.go:48`, `page.go:31`, `section.go:31`) -- `NewField` and `NewRule` now validate their types, but `NewForm`, `NewVersion`, `NewPage`, and `NewSection` constructors still contain `// TODO: Implement domain specific validation`. *(Partially resolved from 4/13 #13, 4/17 #10, 4/18 #11, 4/19 #7, 4/20 #7, 4/22 #7.)*

8. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12, 4/19 #8, 4/20 #8, 4/22 #8.)*

9. **Incomplete error-to-HTTP mapping** -- `ErrMissingTenantID`, `ErrVersionLocked`, `ErrDuplicatePosition`, `ErrDuplicateRuleType`, `ErrPublishedByRequired`, `ErrRetiredByRequired`, `ErrInvalidFieldType`, `ErrInvalidFieldAttributes`, and `ErrInvalidRuleType` all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` (`handlers.go:409-413`) is an empty switch that delegates directly to `httputil.SendErrorResponse`. *(Unresolved from 4/13 #15, 4/17 #12, 4/18 #13, 4/19 #9, 4/20 #9, 4/22 #9.)*

10. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14, 4/19 #10, 4/20 #10, 4/22 #10.)*

11. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15, 4/19 #11, 4/20 #11, 4/22 #11.)*

#### Code Quality

12. **Inconsistent constructor signatures** -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. `NewField` and `NewRule` now correctly use the error return. *(Partially resolved from 4/13 #23, 4/17 #18, 4/18 #19, 4/19 #14, 4/20 #12, 4/22 #12.)*

13. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, it falls through to the 500 default. Should be 400. *(Unresolved from 4/19 #15, 4/20 #13, 4/22 #13.)*

14. **`publishVersion` and `retireVersion` discard the returned `*domain.Version`** (`handlers.go:347-350`, `382-385`) -- Both handlers call the service, receive an updated `*domain.Version`, but respond with `Data: nil`. The mutated version state is never returned to the client. *(Unresolved from 4/20 #14, 4/22 #14.)*

15. **Map iteration order non-deterministic in DTO response mappers** (`dto/version.go:34-36`, `dto/page.go:77-79`, `dto/section.go`) -- `VersionToResponse()`, `PageToResponse()`, and `SectionToResponse()` iterate over Go maps. The order of items in JSON array responses will be non-deterministic across requests. Pages, sections, and fields should be sorted by position key. *(Unresolved from 4/20 #15, 4/22 #15.)*

16. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:242-246`) -- `createVersion` calls `ReadValidateJsonPayload(r, &body)` where `body` is `CreateVersionDto struct{}`. The request body is read and decoded into a type with no fields, meaning any payload is silently discarded. Either the DTO should carry fields or the `ReadValidateJsonPayload` call should be removed. *(Unresolved from 4/22 #16.)*

41. **`FindVersions` returns `ErrNotFound` when form has no versions** (`forms_repository.go:93-96`) -- If a form exists but has no entry in the versions map (which happens when a form is created before any version is added), `FindVersions` returns `common.ErrNotFound` rather than an empty slice. This will produce a 404 when listing versions of a new form. *(New.)*

---

### Submissions Service

#### Architectural

17. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18, 4/20 #16, 4/22 #17.)*

18. **Four handler stubs return 200 OK with empty body** (`handlers.go:87-93`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header. These should either return 501 Not Implemented or not be registered. *(Unresolved from 4/20 #17, 4/22 #18.)*

#### Missing Functionality

19. **Zero test files** in the entire submissions service. *(Unresolved from 4/17 #30, 4/18 #29, 4/19 #21, 4/20 #18, 4/22 #19.)*

20. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:65-71`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22, 4/20 #19, 4/22 #20.)*

21. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23, 4/20 #20, 4/22 #21.)*

22. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24, 4/20 #21, 4/22 #22.)*

23. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25, 4/20 #22, 4/22 #23.)*

24. **No domain events** for cross-service communication. *(Unresolved from 4/17 #35, 4/18 #34, 4/19 #26, 4/20 #23, 4/22 #24.)*

25. **No real authentication.** *(Unresolved from 4/17 #36, 4/18 #35, 4/19 #27, 4/20 #24, 4/22 #25.)*

26. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28, 4/20 #25, 4/22 #26.)*

#### Code Quality

27. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29, 4/20 #26, 4/22 #27.)*

28. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30, 4/20 #27, 4/22 #28.)*

29. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31, 4/20 #28, 4/22 #29.)*

42. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When the request context is cancelled, the `select` on `r.Context().Done()` returns without writing any HTTP response. The client receives a connection drop with no status code. *(New.)*

45. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`) -- The map key is `string`, but the domain uses `domain.SubmissionID` (a named `string` type). `FindById` does `r.submissions[string(id)]`. The map should be `map[domain.SubmissionID]*domain.Submission` for type safety. *(New.)*

---

### Tenants Service

#### Architectural

30. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34, 4/20 #29, 4/22 #30.)*

31. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned in the data sources store. There is no cascade delete and no service-level cleanup. *(Unresolved from 4/20 #30, 4/22 #31.)*

#### Missing Functionality

32. **Zero test files** in the entire tenants service. *(Unresolved from 4/16 #12, 4/17 #50, 4/18 #48, 4/19 #36, 4/20 #31, 4/22 #32.)*

33. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` now validates type and attribute-type agreement but does not validate empty `TenantID` or field lengths. *(Partially resolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37, 4/20 #32, 4/22 #33.)*

34. **No domain events** for cross-service communication. *(Unresolved from 4/16 #14, 4/17 #52, 4/18 #50, 4/19 #38, 4/20 #33, 4/22 #34.)*

36. **No real authentication**. *(Unresolved from 4/16 #16, 4/17 #54, 4/18 #52, 4/19 #40, 4/20 #35, 4/22 #36.)*

37. **`Lookup` service method is a stub** (`data_sources_service.go:131-145`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41, 4/20 #36, 4/22 #37.)*

38. **`DataSourceAttributes` concrete types lack json tags** (`data_source_attributes.go`) -- `StaticDataSourceAttributes`, `ScheduledDataSourceAttributes`, and `QueryDataSourceAttributes` have no `json` struct tags, so JSON marshaling uses Go's default capitalized field names (`Data`, `Type`, `Endpoint`). *(Unresolved from 4/16 #18, 4/17 #56, 4/18 #54, 4/19 #42, 4/20 #37, 4/22 #38.)*

39. **`DataSourceLookup` value object has no validation** (`data_source_lookup.go`) -- `NewDataSourceLookup(code, description)` constructor performs no validation and has no json tags on the struct fields. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43, 4/20 #38, 4/22 #39.)*

#### Code Quality

40. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:64`; `data_sources_repository.go:69`). *(Unresolved from 4/16 #24, 4/17 #60, 4/18 #58, 4/19 #45, 4/20 #40, 4/22 #41.)*

---

### Cross-Service

43. **No graceful shutdown** (all services, e.g. submissions `cmd/server/main.go:58-60`) -- `server.ListenAndServe()` blocks until error, and `log.Fatal` calls `os.Exit`, preventing `defer app.Close()` from executing. There is no signal handling (`os.Signal`) or `server.Shutdown(ctx)` call. *(New.)*

---

### Shared Package

44. **500 errors leak `err.Error()` to clients** (`httputil/http.go:103-107`) -- The `default` case in `SendErrorResponse` returns the raw error string in the JSON response body. In production, this could expose internal details (database errors, file paths, stack traces). Should return a generic message and log the real error server-side. *(New.)*

46. **Zero test files in `pkg/common/`** -- The shared package imported by all three services has no test coverage. `SendErrorResponse`, `ReadValidateJsonPayload`, `TenantMiddleware`, and `NewTypeValidator` are all untested. *(New.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 2, 17 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 7, 33 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 22 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 13 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 31 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 18 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 5 | Redundant double-fetch in `Update()` | Forms |
| **P2** | 44 | 500 errors leak `err.Error()` to clients | Shared |
| **P3** | 6, 19, 32, 46 | Zero test files | All |
| **P3** | 8, 24, 34 | No domain events | All |
| **P3** | 9 | Incomplete error-to-HTTP mapping | Forms |
| **P3** | 27 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4, 40 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 15 | Map iteration non-deterministic in DTO mappers | Forms |
| **P3** | 14 | `publishVersion`/`retireVersion` discard returned version | Forms |
| **P3** | 16 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 38 | `DataSourceAttributes` types lack json tags | Tenants |
| **P3** | 41 | `FindVersions` returns 404 for empty version list | Forms |
| **P3** | 43 | No graceful shutdown | All |

---

## Summary

### Progress Since 4/22

One commit landed (`d25c645c tenants | refactor dto package`) reorganizing the tenants DTO package by splitting `response.go` and `request.go` into per-entity files (`tenant.go`, `data_source.go`, `data_source_attributes.go`, `data_source_lookup.go`). This is a structural improvement with no functional changes.

Additionally, unstaged changes address several issues:

- **`ErrDataSourceAttrParse` now maps to 400** (Tenants) -- New `isBadRequest()` helper in `handlers.go` checks `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes`.
- **Consistent error response delegation in forms handlers** (Forms) -- All `getTenantFromContext` error paths now route through `h.sendErrorResponse` instead of bypassing it.
- **`NewTypeValidator` defensively clones the slice** (Shared) -- `slices.Clone(types)` prevents mutation of the captured slice.
- **Typo fix** (Shared) -- "tyep" corrected to "type" in `ReadValidateJsonPayload` doc comment.
- **Issues #40 and #42 from the 4/22 review were already resolved** -- `DataSourceType` validation in `NewDataSource` and usage of `ErrInvalidSourceType` were both present at the time of the 4/22 review but incorrectly flagged as unresolved.

### New Issues Found

1. **`FindVersions` returns `ErrNotFound` when form has no versions** (Forms #41, P3) -- A newly created form with no versions produces a 404 instead of an empty list.
2. **Context cancel drops response silently** (Submissions #42, P3) -- The `r.Context().Done()` branch returns without writing any HTTP response.
3. **No graceful shutdown** (All #43, P3) -- No signal handling; `defer app.Close()` never executes due to `log.Fatal` calling `os.Exit`.
4. **500 errors leak `err.Error()` to clients** (Shared #44, P2) -- Internal error strings returned verbatim in 500 JSON responses.
5. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (Submissions #45, P3) -- Breaks type safety of the named ID type.
6. **Zero test files in `pkg/common/`** (Shared #46, P3) -- Shared package used by all three services has no test coverage.

### Current State

**Forms Service** remains the most mature. Error response delegation is now consistent across all handlers (unstaged). The primary remaining gaps are: the hardcoded `"placeholder"` user IDs in publish/retire, `Find()` with no tenant filtering (and `getForms` not even reading tenant context), the redundant double-fetch in `Update()`, the `ErrMissingTenantID` -> 500 mapping bug, non-deterministic map iteration in DTO mappers, publish/retire discarding the returned version, `FindVersions` returning 404 for new forms, and the continued absence of test coverage.

**Tenants Service** has improved with the DTO file reorganization and error mapping for `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes` via the new `isBadRequest()` helper. The remaining gaps are: cascade-delete on tenant removal, missing json tags on attribute structs, the stub `Lookup` method, and no test coverage.

**Submissions Service** is unchanged since 4/20. It has tenant middleware wired and tenant isolation on read paths, but four handler stubs return misleading 200 OK responses, `Find()` has no tenant filtering, and the domain model is entirely anemic (no constructors, no validation, no write operations, no status constants). Request DTOs are still empty.

**Shared Package** (`pkg/common`) gained the defensive `slices.Clone` in `NewTypeValidator` and the typo fix. `ErrMissingTenantID` still maps to 500. The `SendErrorResponse` default case leaks internal error details to clients.

### Highest-Impact Improvements

1. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500 instead of 400)
2. **Stop leaking `err.Error()` in 500 responses** (P2 -- security risk in production)
3. **Add tenant filtering to `Find()` methods** across forms and submissions (P2 -- data isolation)
4. **Eliminate redundant double-fetch** in forms `Update()` (P2 -- unnecessary repository call)
5. **Add test coverage** starting with `pkg/common` and service layers (P3 -- long-term reliability)
