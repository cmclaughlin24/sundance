# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/23 Review

1. ~~Forms `mongodb.Bootstrap` returns empty `Repository{}`~~ (Forms, P0) -- `Bootstrap` (`services/forms/.../mongodb/mongodb.go:10-17`) now returns a fully wired `*ports.Repository` with `Database` via `database.NewMongoDBDatabase(client, db)` and `Forms` via `newMongoDBFormsRepository(db, logger)`. MongoDB repos implement all 7 `FormsRepository` interface methods: `Find`, `FindById`, `Upsert`, `FindVersions`, `FindVersion`, `FindNextVersionNumber`, `UpsertVersion`. Document mapping covers forms, versions, pages, sections, fields (with polymorphic attribute strategy), and rules. `CreateVersion` correctly uses transactions via `BeginTx`/`CommitTx`/`RollbackTx`. *(Partially resolves 4/23 #15 -- forms half only. Submissions still returns empty `Repository{}`.)*

2. ~~`UpdatePages` calls `SetPages()` with no arguments -- data loss~~ (Forms, P0) -- `v.SetPages()` at `version.go:120` now correctly forwards the variadic argument as `v.SetPages(pages...)`. *(Introduced and resolved same day.)*

3. ~~Docker Compose configured for MongoDB replica set~~ (All services) -- `docker-compose` updated to support `replSet` for MongoDB transactions.

4. ~~`cursor.Close(ctx)` uses outer context instead of session context~~ (Shared, P3) -- `cursor.Close(ctx)` in `MongoDBRepository.Find` (`mongodb_repository.go:45`) now correctly uses `sctx` (the session-aware context) instead of the outer `ctx`. *(Unstaged.)*

5. ~~`FindVersions` behavioral inconsistency between adapters~~ (Forms, P3) -- In-memory `FindVersions` (`inmemory/forms_repository.go:70`) now returns `make([]*domain.Version, 0), nil` instead of `common.ErrNotFound` when a form has no versions, matching the MongoDB adapter behavior. *(Unstaged.)*

6. ~~`NewMongoDBFormsRepository` is exported~~ (Forms, P3) -- Renamed to `newMongoDBFormsRepository` (`forms_repository.go:20`), consistent with the tenants service's unexported constructor pattern. Call site in `mongodb.go` updated. *(Unstaged.)*

7. ~~`FindNextVersionNumber` bypasses `mongo.WithSession`~~ (Forms, P2) -- Now wrapped in `mongo.WithSession` with correct session context used throughout (`forms_repository.go:113-141`). Cursor operations and `cursor.All` both use `sctx`. *(Unstaged.)*

8. ~~`Publish` and `Retire` ignore their timestamp parameters~~ (Forms, P2) -- Dead `publishedAt` and `retiredAt` parameters removed from method signatures (`version.go:130`, `149`). Callers in `forms_service.go:207,235` updated to pass only `command.UserID`. *(Unstaged.)*

9. ~~Redundant double-fetch in `Update()`~~ (Forms, P2) -- `isValidAccess()` call removed from `Update()` (`forms_service.go:71-86`). The form is now fetched once via `FindById`, and tenant ownership is checked inline (`form.TenantID != command.TenantID`). *(Unstaged. Resolved from 4/20 #5, 4/22 #5, 4/23 #5.)*

10. ~~Incomplete error-to-HTTP mapping~~ (Forms, P3) -- Forms `sendErrorResponse` (`handlers.go:411-432`) now has an `isBadRequest` helper that maps `ErrVersionLocked`, `ErrDuplicatePosition`, `ErrInvalidRuleType`, `ErrDuplicateRuleType`, `ErrPublishedByRequired`, `ErrRetiredByRequired`, `ErrInvalidFieldType`, and `ErrInvalidFieldAttributes` to 400 Bad Request. *(Unstaged. Partially resolves 4/13 #15, 4/17 #12, 4/18 #13, 4/19 #9, 4/20 #9, 4/22 #9, 4/23 #7. Note: `ErrMissingTenantID` in shared middleware is not addressed -- see #7.)*

11. ~~Map iteration order non-deterministic in DTO response mappers~~ (Forms, P3) -- `VersionToResponse` (`dto/version.go:37-40`), `PageToResponse` (`dto/page.go:83-87`), and `SectionToResponse` (`dto/section.go:67-71`) now all use `slices.Sorted(maps.Keys(...))` to iterate by position key, producing deterministic JSON array ordering. *(Unstaged. Resolved from 4/20 #15, 4/22 #15, 4/23 #12.)*

12. ~~`GetDataSourceLookupsCommand` is a read operation misclassified as a command~~ (Tenants, CQRS) -- Moved from `commands.go` to `query.go` as `GetDataSourceLookupsQuery`. `DataSourcesService.Lookup` signature (`ports/primary.go:28`) and handler call site (`handlers.go:378`) updated. *(Unstaged.)*

13. ~~Go naming convention deviations in shared HTTP utilities~~ (Shared, P3) -- `ApiResponse` → `APIResponse`, `ApiErrorResponse` → `APIErrorResponse`, `ReadJsonPayload` → `ReadJSONPayload`, `ReadValidateJsonPayload` → `ReadValidateJSONPayload`, `SendJsonResponse` → `SendJSONResponse` (`httputil/http.go`). All call sites across all three services updated. *(Unstaged.)*

14. ~~`Position` field publicly accessible on `Page`, `Section`, and `Field`~~ (Forms, DDD) -- `Position` is now a private field encapsulated via a `withPosition` embedded struct with a `GetPosition()` accessor. `SetPages`, `SetSections`, and `SetFields` now use `GetPosition()` instead of direct field access. MongoDB document mappers and DTO response mappers updated to use `GetPosition()`. `baseWithRules` renamed to `withRules` for consistency. *(Unstaged.)*

15. ~~`HydratePage` and `HydrateSection` don't initialize child maps~~ (Forms, P3) -- Both now include `sections: make(map[int]*Section)` and `fields: make(map[int]*Field)` respectively in their initializers, consistent with `NewPage`/`NewSection`. *(Unstaged.)*

---

## Will Not Fix

1. **`mongo.WithSession` called with potentially nil session from `SessionFromContext`** (Previously Shared #43) -- `mongo.SessionFromContext(ctx)` returns `nil` when no session exists in the context (non-transactional calls). The MongoDB Go driver v2 handles a nil session by continuing without a session, executing the callback with the original context. This is safe and correct behavior for non-transactional read paths.

See [4/23 review](code-review-4-23-26.md) for the prior Will Not Fix list (8 items).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:330`, `365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1, 4/19 #1, 4/20 #1, 4/22 #1, 4/23 #1.)*

#### Architectural

2. **`Find()` has no tenant filtering** (`forms_service.go:29-31`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. The `getForms` handler (`handlers.go:30-55`) is also the only handler that never calls `getTenantFromContext`, meaning tenant context is never even read on this path. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5, 4/19 #2, 4/20 #2, 4/22 #2, 4/23 #2.)*

3. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6, 4/19 #4, 4/20 #3, 4/22 #3, 4/23 #3.)*

4. **`time.Now()` called directly in the domain layer** (`form.go:29,55`; `version.go:48,125,143,165`) -- Should be injected via a `Clock` interface or function for testability. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8, 4/19 #5, 4/20 #4, 4/22 #4, 4/23 #4. Note: `forms_service.go` no longer calls `time.Now()` directly after #8 resolution.)*

#### Missing Functionality

5. **Domain validation partially unimplemented** (`form.go:29,42`, `version.go:48`, `page.go:31`, `section.go:31`) -- `NewField` and `NewRule` now validate their types, but `NewForm`, `NewVersion`, `NewPage`, and `NewSection` constructors still contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13, 4/17 #10, 4/18 #11, 4/19 #7, 4/20 #7, 4/22 #7, 4/23 #6.)*

6. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15, 4/19 #11, 4/20 #11, 4/22 #11, 4/23 #8.)*

#### Code Quality

7. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, it falls through to the 500 default. Should be 400. *(Unresolved from 4/19 #15, 4/20 #13, 4/22 #13, 4/23 #10.)*

8. **Inconsistent constructor signatures** -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. `NewField` and `NewRule` now correctly use the error return. *(Unresolved from 4/13 #23, 4/17 #18, 4/18 #19, 4/19 #14, 4/20 #12, 4/22 #12, 4/23 #9.)*

9. **`publishVersion` and `retireVersion` discard the returned `*domain.Version`** (`handlers.go:347-350`, `382-385`) -- Both handlers call the service, receive an updated `*domain.Version`, but respond with `Data: nil`. The mutated version state is never returned to the client. *(Unresolved from 4/20 #14, 4/22 #14, 4/23 #11.)*

10. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:242-246`) -- `createVersion` calls `ReadValidateJSONPayload(r, &body)` where `body` is `CreateVersionRequest struct{}`. The request body is read and decoded into a type with no fields, meaning any payload is silently discarded. Either the DTO should carry fields or the `ReadValidateJSONPayload` call should be removed. *(Unresolved from 4/22 #16, 4/23 #13.)*

---

### Submissions Service

#### Bugs

11. **`mongodb.Bootstrap` returns empty `&ports.Repository{}`** (`submissions/.../mongodb/mongodb.go:10-12`) -- `Bootstrap` returns `&ports.Repository{}` with nil interface fields. Any repository call will panic with a nil pointer dereference when the `mongodb` driver is selected. *(Unresolved from 4/23 #15, submissions half.)*

#### Architectural

12. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18, 4/20 #16, 4/22 #17, 4/23 #16.)*

13. **Four handler stubs return 200 OK with empty body** (`handlers.go:87-93`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header. These should either return 501 Not Implemented or not be registered. *(Unresolved from 4/20 #17, 4/22 #18, 4/23 #17.)*

#### Missing Functionality

14. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:65-71`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22, 4/20 #19, 4/22 #20, 4/23 #18.)*

15. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23, 4/20 #20, 4/22 #21, 4/23 #19.)*

16. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24, 4/20 #21, 4/22 #22, 4/23 #20.)*

17. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25, 4/20 #22, 4/22 #23, 4/23 #21.)*

18. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28, 4/20 #25, 4/22 #26, 4/23 #22.)*

#### Code Quality

19. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29, 4/20 #26, 4/22 #27, 4/23 #23.)*

20. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30, 4/20 #27, 4/22 #28, 4/23 #24.)*

21. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31, 4/20 #28, 4/22 #29, 4/23 #25.)*

22. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When the request context is cancelled, the `select` on `r.Context().Done()` returns without writing any HTTP response. The client receives a connection drop with no status code. *(Unresolved from 4/23 #26.)*

23. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`) -- The map key is `string`, but the domain uses `domain.SubmissionID` (a named `string` type). The map should be `map[domain.SubmissionID]*domain.Submission` for type safety. *(Unresolved from 4/23 #27.)*

---

### Tenants Service

#### Architectural

24. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34, 4/20 #29, 4/22 #30, 4/23 #28.)*

25. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned. *(Unresolved from 4/20 #30, 4/22 #31, 4/23 #29.)*

#### Missing Functionality

26. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` now validates type and attribute-type agreement but does not validate empty `TenantID` or field lengths. *(Unresolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37, 4/20 #32, 4/22 #33, 4/23 #30.)*

27. **`Lookup` service method is a stub** (`data_sources_service.go:131-145`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41, 4/20 #36, 4/22 #37, 4/23 #31.)*

28. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup(value, label)` constructor performs no validation and has no json tags on the struct fields. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43, 4/20 #38, 4/22 #39, 4/23 #32.)*

#### Code Quality

29. **`Ping` uses `context.Background()` with no timeout** (`pkg/database/mongodb.go:43`) -- If MongoDB is unreachable, the service hangs forever at startup. Should use `context.WithTimeout`. *(Unresolved from 4/23 #33.)*

---

### Cross-Service

#### Architectural

30. **Zero test files** in all three services and shared packages. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10, 4/19 #6, 4/20 #6, 4/22 #6, 4/23 #34.)*

31. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12, 4/19 #8, 4/20 #8, 4/22 #8, 4/23 #35.)*

32. **No real authentication** -- `X-Tenant-ID` is blindly trusted across all services. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14, 4/19 #10, 4/20 #10, 4/22 #10, 4/23 #36.)*

33. **No graceful shutdown** (all services, e.g. submissions `cmd/server/main.go:58-60`) -- `server.ListenAndServe()` blocks until error, and `log.Fatal` calls `os.Exit`, preventing `defer app.Close()` from executing. There is no signal handling (`os.Signal`) or `server.Shutdown(ctx)` call. *(Unresolved from 4/23 #37.)*

---

### Shared Package

34. **500 errors leak `err.Error()` to clients** (`httputil/http.go:103-107`) -- The `default` case in `SendErrorResponse` returns the raw error string in the JSON response body. In production, this could expose internal details. Should return a generic message and log the real error server-side. *(Unresolved from 4/23 #38.)*

35. **`ReadJsonFile` naming convention** (`pkg/common/utils.go:8`) -- Should be `ReadJSONFile` per Go initialism conventions. `FindById` method names throughout ports should be `FindByID`. *(Remaining from naming review.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 11 | `bootstrapMongoDB` returns empty `Repository{}` -- crash on any repo call | Submissions |
| **P2** | 2, 12 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 5, 26 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 16 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 7 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 25 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 13 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 34 | 500 errors leak `err.Error()` to clients | Shared |
| **P3** | 30 | Zero test files | All |
| **P3** | 31 | No domain events | All |
| **P3** | 19 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4 | `time.Now()` in domain layer | Forms |
| **P3** | 9 | `publishVersion`/`retireVersion` discard returned version | Forms |
| **P3** | 10 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 33 | No graceful shutdown | All |
| **P3** | 29 | `Ping` uses `context.Background()` with no timeout | All |
| **P3** | 35 | `ReadJsonFile` and `FindById` naming conventions | Shared |

---

## Summary

### Progress Since 4/23

The primary focus since the last review was implementing MongoDB persistence for the forms service, followed by two rounds of fixes addressing architectural, code quality, DDD, and Go idiom issues:

- **Forms MongoDB repositories fully implemented** -- All 7 `FormsRepository` interface methods now have MongoDB backing: `Find`, `FindById`, `Upsert`, `FindVersions`, `FindVersion`, `FindNextVersionNumber`, `UpsertVersion`. BSON document mapping covers the full entity graph (forms, versions, pages, sections, fields with polymorphic attribute strategies, and rules). `CreateVersion` uses `BeginTx`/`CommitTx`/`RollbackTx` for transactional version number generation and persistence.
- **Docker Compose updated** for MongoDB replica set support, enabling transaction support.
- **`UpdatePages` data loss bug introduced and fixed** -- `SetPages()` was called without forwarding the variadic `pages` argument, silently wiping all pages on every `UpdateVersion` call. Corrected to `SetPages(pages...)`.
- **`FindNextVersionNumber` session handling fixed** (unstaged) -- Now wrapped in `mongo.WithSession` with correct session context throughout, ensuring transactional consistency.
- **Dead timestamp parameters removed** (unstaged) -- `Publish` and `Retire` no longer accept unused `time.Time` parameters.
- **Redundant double-fetch eliminated** (unstaged) -- `Update()` no longer calls `isValidAccess()` separately; the form is fetched once and tenant ownership checked inline.
- **`FindVersions` adapter inconsistency fixed** (unstaged) -- In-memory implementation now returns an empty slice instead of `ErrNotFound` for forms with no versions.
- **MongoDB repository constructor unexported** (unstaged) -- `NewMongoDBFormsRepository` renamed to `newMongoDBFormsRepository`, consistent with tenants service pattern.
- **`cursor.Close` context fixed** (unstaged) -- Shared `MongoDBRepository.Find` now uses the session context for cursor close.
- **Forms error-to-HTTP mapping implemented** (unstaged) -- `sendErrorResponse` now maps 8 domain errors to 400 Bad Request via `isBadRequest` helper, matching the tenants service pattern.
- **DTO response mapper sort order fixed** (unstaged) -- `VersionToResponse`, `PageToResponse`, and `SectionToResponse` now use `slices.Sorted(maps.Keys(...))` for deterministic JSON array ordering by position.
- **CQRS naming violation corrected** (unstaged) -- `GetDataSourceLookupsCommand` moved to `query.go` as `GetDataSourceLookupsQuery` with all references updated.
- **Go naming conventions applied** (unstaged) -- `ApiResponse` → `APIResponse`, `ApiErrorResponse` → `APIErrorResponse`, `ReadJsonPayload` → `ReadJSONPayload`, `ReadValidateJsonPayload` → `ReadValidateJSONPayload`, `SendJsonResponse` → `SendJSONResponse`. All call sites across all three services updated.
- **Domain encapsulation improved** (unstaged) -- `Position` field on `Page`, `Section`, and `Field` is now private via a `withPosition` embedded struct with `GetPosition()` accessor. `baseWithRules` renamed to `withRules`. All external access updated to use getters.
- **`HydratePage` and `HydrateSection` now initialize child maps** (unstaged) -- Both include `sections: make(map[int]*Section)` and `fields: make(map[int]*Field)` respectively.

### Issues Moved to Will Not Fix

1. **`mongo.WithSession` with nil session** -- The MongoDB Go driver v2 handles nil sessions by continuing without a session. Safe for non-transactional paths.

### Current State

**Forms Service** now has full MongoDB support with CRUD for forms and versions, BSON document mappers, and transactional `CreateVersion`. Unstaged changes address three P2 issues (session handling, dead parameters, double-fetch), resolve the incomplete error-to-HTTP mapping, fix DTO sort ordering, and improve domain encapsulation. The forms domain layer is the strongest in the codebase: `Version` is a well-designed aggregate root with a proper state machine (`draft -> active -> retired`), private child collections (`pages`, `sections`, `fields`), private position fields via `withPosition` embedding, duplicate position guards, and field type/attribute validation. Several `New*` constructors have unfulfilled validation TODOs (`form.go:32`, `version.go:51`, `page.go:33`, `section.go:33`). Aggregate getters (`GetPages()`, `GetSections()`, `GetFields()`) still return direct map references, allowing callers to mutate internal state without going through the aggregate root. CQRS-lite is clean: commands and queries are separate types in separate files.

**Tenants Service** saw one structural improvement: `GetDataSourceLookupsCommand` corrected to `GetDataSourceLookupsQuery` and moved to `query.go`. Domain entities have controlled mutation methods with validation on `DataSource` type/attribute alignment. MongoDB repositories are fully implemented. The REST adapter has partial error translation (`isBadRequest` helper covers three domain errors), making it the most complete error boundary in the codebase. `NewTenant` has no validation and `Tenant.Update()` lacks a nil guard.

**Submissions Service** is the weakest across all evaluation criteria. The domain is entirely anemic: `Submission` and `SubmissionAttempt` are bare structs with no factory methods, no `Hydrate*` functions, no validation, no business methods, and no state machine. `SubmissionStatus` is declared as a type but has no defined constants. The repository interface is read-only (no `Create`/`Update`/`Delete`). Four handler stubs return 200 OK with empty bodies. The MongoDB `Bootstrap` still returns an empty `Repository{}` (P0 crash). The REST adapter's `sendErrorResponse` is a no-op.

**Shared Package** -- `pkg/database` provides clean infrastructure abstractions (`Database` interface, `MongoDBDatabase`, `InMemoryDatabase`, `MongoDBRepository[T]`). `pkg/common/httputil` now uses idiomatic Go naming (`APIResponse`, `SendJSONResponse`, etc.). `pkg/common` remains a grab-bag package (errors, HTTP utils, validation, tenant middleware) that would benefit from splitting into focused packages. Remaining naming deviations: `ReadJsonFile` should be `ReadJSONFile`, and `FindById` method names throughout ports should be `FindByID`.

**Hexagonal Architecture** -- Dependency direction is correct throughout: adapters import core, never reversed. No cross-adapter imports exist. Domain layers are pure (stdlib + uuid only, with one acceptable `validate` utility). DTO and document mapping is explicit and lives in the adapter layer. The `Database` interface in `ports/secondary.go` represents a small leak of transaction-management (infrastructure concern) into the core, but is a pragmatic trade-off.

**CQRS-Lite** -- Commands and queries are well-separated into distinct types and files across all services. The service interfaces combine read and write methods (standard CQRS-lite). The `GetDataSourceLookupsCommand` naming violation has been corrected.

**DDD** -- No domain events exist anywhere. Aggregate boundaries are clear in Forms (`Version` as root over pages/sections/fields) with improved encapsulation via private `withPosition` fields, but getters still leak internal map references. Submissions has no aggregate pattern, no invariant enforcement, and no factory methods.

### Highest-Impact Improvements

1. **Implement MongoDB repositories for submissions** (P0 -- `Bootstrap` still returns empty `Repository{}`)
2. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
3. **Add tenant filtering to `Find()` methods** across forms and submissions (P2 -- data isolation)
4. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
