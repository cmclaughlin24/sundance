# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/23 Review

1. ~~Forms `mongodb.Bootstrap` returns empty `Repository{}`~~ (Forms, P0) -- `Bootstrap` (`services/forms/.../mongodb/mongodb.go:10-17`) now returns a fully wired `*ports.Repository` with `Database` via `database.NewMongoDBDatabase(client, db)` and `Forms` via `NewMongoDBFormsRepository(db, logger)`. MongoDB repos implement all 7 `FormsRepository` interface methods: `Find`, `FindById`, `Upsert`, `FindVersions`, `FindVersion`, `FindNextVersionNumber`, `UpsertVersion`. Document mapping covers forms, versions, pages, sections, fields (with polymorphic attribute strategy), and rules. `CreateVersion` correctly uses transactions via `BeginTx`/`CommitTx`/`RollbackTx`. *(Partially resolves 4/23 #15 -- forms half only. Submissions still returns empty `Repository{}`.)*

2. ~~`UpdatePages` calls `SetPages()` with no arguments -- data loss~~ (Forms, P0) -- `v.SetPages()` at `version.go:120` now correctly forwards the variadic argument as `v.SetPages(pages...)`. *(Introduced and resolved same day.)*

3. ~~Docker Compose configured for MongoDB replica set~~ (All services) -- `docker-compose` updated to support `replSet` for MongoDB transactions.

---

## Will Not Fix

No new items moved to Will Not Fix in this review. See [4/23 review](code-review-4-23-26.md) for the full list (8 items).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:330`, `365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1, 4/19 #1, 4/20 #1, 4/22 #1, 4/23 #1.)*

#### Architectural

2. **`Find()` has no tenant filtering** (`forms_service.go:29-31`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. The `getForms` handler (`handlers.go:30-55`) is also the only handler that never calls `getTenantFromContext`, meaning tenant context is never even read on this path. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5, 4/19 #2, 4/20 #2, 4/22 #2, 4/23 #2.)*

3. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6, 4/19 #4, 4/20 #3, 4/22 #3, 4/23 #3.)*

4. **`time.Now()` called directly in the domain and service layers** (`form.go:29,55`; `version.go:48,125,143,165`; `forms_service.go:210,238`) -- Should be injected via a `Clock` interface or function for testability. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8, 4/19 #5, 4/20 #4, 4/22 #4, 4/23 #4.)*

5. **Redundant double-fetch in `Update()`** (`forms_service.go:78-86`) -- `isValidAccess()` fetches the form from the repository to check tenant ownership, then `Update()` immediately fetches the same form again via `FindById()`. The form should be fetched once and reused. *(Unresolved from 4/20 #5, 4/22 #5, 4/23 #5.)*

6. **(New) `FindNextVersionNumber` bypasses `mongo.WithSession`** (`forms_repository.go:114`) -- Calls `r.versions.Collection().Find(ctx, ...)` directly instead of going through the base `r.versions.Find()` or using `mongo.WithSession`. Since this is called within a transaction in `CreateVersion` (`forms_service.go:142`), the query may not see uncommitted writes within the same transaction, potentially causing duplicate version numbers under concurrent access.

7. **(New) `Publish` and `Retire` ignore their timestamp parameters** (`version.go:130-150`, `152-172`) -- Both methods accept a `time.Time` parameter (`publishedAt`/`retiredAt`) but call `time.Now()` internally and use that instead. The caller in `forms_service.go:210,238` passes `time.Now()` anyway so the practical impact is nil currently, but the API contract is misleading and the parameters are dead code. Either use the parameters or remove them from the signatures.

#### Missing Functionality

8. **Domain validation partially unimplemented** (`form.go:29,42`, `version.go:48`, `page.go:31`, `section.go:31`) -- `NewField` and `NewRule` now validate their types, but `NewForm`, `NewVersion`, `NewPage`, and `NewSection` constructors still contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/13 #13, 4/17 #10, 4/18 #11, 4/19 #7, 4/20 #7, 4/22 #7, 4/23 #6.)*

9. **Incomplete error-to-HTTP mapping** -- `ErrMissingTenantID`, `ErrVersionLocked`, `ErrDuplicatePosition`, `ErrDuplicateRuleType`, `ErrPublishedByRequired`, `ErrRetiredByRequired`, `ErrInvalidFieldType`, `ErrInvalidFieldAttributes`, and `ErrInvalidRuleType` all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` (`handlers.go:409-413`) is an empty switch that delegates directly to `httputil.SendErrorResponse`. *(Unresolved from 4/13 #15, 4/17 #12, 4/18 #13, 4/19 #9, 4/20 #9, 4/22 #9, 4/23 #7.)*

10. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15, 4/19 #11, 4/20 #11, 4/22 #11, 4/23 #8.)*

#### Code Quality

11. **Inconsistent constructor signatures** -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. `NewField` and `NewRule` now correctly use the error return. *(Unresolved from 4/13 #23, 4/17 #18, 4/18 #19, 4/19 #14, 4/20 #12, 4/22 #12, 4/23 #9.)*

12. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, it falls through to the 500 default. Should be 400. *(Unresolved from 4/19 #15, 4/20 #13, 4/22 #13, 4/23 #10.)*

13. **`publishVersion` and `retireVersion` discard the returned `*domain.Version`** (`handlers.go:347-350`, `382-385`) -- Both handlers call the service, receive an updated `*domain.Version`, but respond with `Data: nil`. The mutated version state is never returned to the client. *(Unresolved from 4/20 #14, 4/22 #14, 4/23 #11.)*

14. **Map iteration order non-deterministic in DTO response mappers** (`dto/version.go:34-36`, `dto/page.go:77-79`, `dto/section.go`) -- `VersionToResponse()`, `PageToResponse()`, and `SectionToResponse()` iterate over Go maps. The order of items in JSON array responses will be non-deterministic across requests. Pages, sections, and fields should be sorted by position key. *(Unresolved from 4/20 #15, 4/22 #15, 4/23 #12.)*

15. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:242-246`) -- `createVersion` calls `ReadValidateJsonPayload(r, &body)` where `body` is `CreateVersionDto struct{}`. The request body is read and decoded into a type with no fields, meaning any payload is silently discarded. Either the DTO should carry fields or the `ReadValidateJsonPayload` call should be removed. *(Unresolved from 4/22 #16, 4/23 #13.)*

16. **`FindVersions` returns `ErrNotFound` when form has no versions (in-memory only)** (`inmemory/forms_repository.go:67-71`) -- The in-memory implementation returns `common.ErrNotFound` rather than an empty slice for a form with no versions. The MongoDB implementation correctly returns `([], nil)`. Behavioral inconsistency between adapters. *(Unresolved from 4/23 #14, updated with cross-adapter comparison.)*

17. **(New) `NewMongoDBFormsRepository` is exported** (`forms_repository.go:20`) -- Tenants service unexports repository constructors (`newMongoDBTenantsRepository`, `newMongoDBDataSourcesRepository`). Forms exports `NewMongoDBFormsRepository` unnecessarily since it's only called from `Bootstrap` within the same package.

---

### Submissions Service

#### Bugs

18. **`mongodb.Bootstrap` returns empty `&ports.Repository{}`** (`submissions/.../mongodb/mongodb.go:10-12`) -- `Bootstrap` returns `&ports.Repository{}` with nil interface fields. Any repository call will panic with a nil pointer dereference when the `mongodb` driver is selected. *(Unresolved from 4/23 #15, submissions half.)*

#### Architectural

19. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18, 4/20 #16, 4/22 #17, 4/23 #16.)*

20. **Four handler stubs return 200 OK with empty body** (`handlers.go:87-93`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header. These should either return 501 Not Implemented or not be registered. *(Unresolved from 4/20 #17, 4/22 #18, 4/23 #17.)*

#### Missing Functionality

21. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:65-71`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22, 4/20 #19, 4/22 #20, 4/23 #18.)*

22. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23, 4/20 #20, 4/22 #21, 4/23 #19.)*

23. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24, 4/20 #21, 4/22 #22, 4/23 #20.)*

24. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25, 4/20 #22, 4/22 #23, 4/23 #21.)*

25. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28, 4/20 #25, 4/22 #26, 4/23 #22.)*

#### Code Quality

26. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29, 4/20 #26, 4/22 #27, 4/23 #23.)*

27. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30, 4/20 #27, 4/22 #28, 4/23 #24.)*

28. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31, 4/20 #28, 4/22 #29, 4/23 #25.)*

29. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When the request context is cancelled, the `select` on `r.Context().Done()` returns without writing any HTTP response. The client receives a connection drop with no status code. *(Unresolved from 4/23 #26.)*

30. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`) -- The map key is `string`, but the domain uses `domain.SubmissionID` (a named `string` type). The map should be `map[domain.SubmissionID]*domain.Submission` for type safety. *(Unresolved from 4/23 #27.)*

---

### Tenants Service

#### Architectural

31. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34, 4/20 #29, 4/22 #30, 4/23 #28.)*

32. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned. *(Unresolved from 4/20 #30, 4/22 #31, 4/23 #29.)*

#### Missing Functionality

33. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` now validates type and attribute-type agreement but does not validate empty `TenantID` or field lengths. *(Unresolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37, 4/20 #32, 4/22 #33, 4/23 #30.)*

34. **`Lookup` service method is a stub** (`data_sources_service.go:131-145`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41, 4/20 #36, 4/22 #37, 4/23 #31.)*

35. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup(value, label)` constructor performs no validation and has no json tags on the struct fields. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43, 4/20 #38, 4/22 #39, 4/23 #32.)*

#### Code Quality

36. **`Ping` uses `context.Background()` with no timeout** (`pkg/database/mongodb.go:43`) -- If MongoDB is unreachable, the service hangs forever at startup. Should use `context.WithTimeout`. *(Unresolved from 4/23 #33.)*

---

### Cross-Service

#### Architectural

37. **Zero test files** in all three services and shared packages. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10, 4/19 #6, 4/20 #6, 4/22 #6, 4/23 #34.)*

38. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12, 4/19 #8, 4/20 #8, 4/22 #8, 4/23 #35.)*

39. **No real authentication** -- `X-Tenant-ID` is blindly trusted across all services. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14, 4/19 #10, 4/20 #10, 4/22 #10, 4/23 #36.)*

40. **No graceful shutdown** (all services, e.g. submissions `cmd/server/main.go:58-60`) -- `server.ListenAndServe()` blocks until error, and `log.Fatal` calls `os.Exit`, preventing `defer app.Close()` from executing. There is no signal handling (`os.Signal`) or `server.Shutdown(ctx)` call. *(Unresolved from 4/23 #37.)*

---

### Shared Package

41. **500 errors leak `err.Error()` to clients** (`httputil/http.go:103-107`) -- The `default` case in `SendErrorResponse` returns the raw error string in the JSON response body. In production, this could expose internal details. Should return a generic message and log the real error server-side. *(Unresolved from 4/23 #38.)*

42. **(New) `cursor.Close(ctx)` uses outer context instead of session context** (`mongodb_repository.go:45`) -- Inside the `WithSession` callback, the session-aware context is `sctx`, but `cursor.Close(ctx)` uses the outer `ctx`. Should be `cursor.Close(sctx)` to ensure the close operation participates in the session.

43. **(New) `mongo.WithSession` called with potentially nil session** (`mongodb_repository.go:38,64,83`) -- `mongo.SessionFromContext(ctx)` returns `nil` when no session exists in the context (non-transactional calls). Behavior of `mongo.WithSession` with a nil session depends on driver implementation -- may panic or silently proceed. Should guard with a nil check or fall back to a non-session call path.

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 18 | `bootstrapMongoDB` returns empty `Repository{}` -- crash on any repo call | Submissions |
| **P2** | 2, 19 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 6 | `FindNextVersionNumber` bypasses `mongo.WithSession` | Forms |
| **P2** | 7 | `Publish`/`Retire` ignore timestamp parameters | Forms |
| **P2** | 8, 33 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 23 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 12 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 32 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 20 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 5 | Redundant double-fetch in `Update()` | Forms |
| **P2** | 41 | 500 errors leak `err.Error()` to clients | Shared |
| **P3** | 37 | Zero test files | All |
| **P3** | 38 | No domain events | All |
| **P3** | 9 | Incomplete error-to-HTTP mapping | Forms |
| **P3** | 26 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4 | `time.Now()` in domain/service layers | Forms |
| **P3** | 14 | Map iteration non-deterministic in DTO mappers | Forms |
| **P3** | 13 | `publishVersion`/`retireVersion` discard returned version | Forms |
| **P3** | 15 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 16 | `FindVersions` behavioral inconsistency between adapters | Forms |
| **P3** | 40 | No graceful shutdown | All |
| **P3** | 36 | `Ping` uses `context.Background()` with no timeout | All |
| **P3** | 42 | `cursor.Close` uses wrong context | Shared |
| **P3** | 43 | `WithSession` with potentially nil session | Shared |
| **P3** | 17 | Exported repository constructor inconsistency | Forms |

---

## Summary

### Progress Since 4/23

The primary focus since the last review was implementing MongoDB persistence for the forms service:

- **Forms MongoDB repositories fully implemented** -- All 7 `FormsRepository` interface methods now have MongoDB backing: `Find`, `FindById`, `Upsert`, `FindVersions`, `FindVersion`, `FindNextVersionNumber`, `UpsertVersion`. BSON document mapping covers the full entity graph (forms, versions, pages, sections, fields with polymorphic attribute strategies, and rules). `CreateVersion` uses `BeginTx`/`CommitTx`/`RollbackTx` for transactional version number generation and persistence.
- **Docker Compose updated** for MongoDB replica set support, enabling transaction support.
- **`UpdatePages` data loss bug introduced and fixed** -- `SetPages()` was called without forwarding the variadic `pages` argument, silently wiping all pages on every `UpdateVersion` call. Corrected to `SetPages(pages...)`.

### New Issues Found

1. **`FindNextVersionNumber` bypasses `mongo.WithSession`** (#6, P2) -- May cause duplicate version numbers under concurrent transactional access.
2. **`Publish`/`Retire` ignore timestamp parameters** (#7, P2) -- Dead parameters create a misleading API contract.
3. **`cursor.Close` uses wrong context** (#42, P3) -- Minor correctness issue in shared `MongoDBRepository.Find`.
4. **`WithSession` with potentially nil session** (#43, P3) -- Risk depends on driver behavior for non-transactional calls.
5. **Exported constructor inconsistency** (#17, P3) -- `NewMongoDBFormsRepository` is exported; tenants pattern unexports constructors.

### Current State

**Forms Service** now has full MongoDB support, resolving the P0 bootstrap crash from the 4/23 review. The service is functionally complete for CRUD operations on forms and versions with persistent storage. Two new P2 architectural issues were found in the MongoDB implementation (`FindNextVersionNumber` session handling, dead timestamp parameters). All prior code quality and missing functionality issues remain open.

**Tenants Service** is unchanged since 4/23. Remains the most stable service with the most complete MongoDB implementation.

**Submissions Service** is unchanged. The P0 MongoDB bootstrap crash remains -- `Bootstrap` still returns an empty `Repository{}` with nil fields.

**Shared Package** (`pkg/database`) -- Two new minor issues found in `MongoDBRepository`: wrong context for cursor close and nil session risk with `WithSession`.

### Highest-Impact Improvements

1. **Implement MongoDB repositories for submissions** (P0 -- `Bootstrap` still returns empty `Repository{}`)
2. **Fix `FindNextVersionNumber` to use `mongo.WithSession`** (P2 -- transaction correctness)
3. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
4. **Add tenant filtering to `Find()` methods** across forms and submissions (P2 -- data isolation)
5. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
