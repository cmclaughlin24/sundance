# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/22 Review

1. ~~`DataSourceType` not validated in domain constructors~~ (Tenants, P3) -- `NewDataSource` now calls `isValidSourceType()` and returns `ErrInvalidSourceType` when invalid (`data_source.go:44-46`). Was incorrectly flagged as unresolved in the 4/22 review; already fixed at that time.

2. ~~`ErrInvalidSourceType` defined but unused~~ (Tenants, P3) -- Now used in `NewDataSource` (`data_source.go:45`) and checked in `isBadRequest()` (`handlers.go:434`). Was incorrectly flagged as unresolved in the 4/22 review; already fixed at that time.

3. ~~`ErrDataSourceAttrParse` falls through to 500~~ (Tenants, P3) -- New `isBadRequest()` helper (`handlers.go:432-436`) checks `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes`. All three now correctly map to 400 Bad Request.

4. ~~Inconsistent error response delegation in forms handlers~~ (Forms, P3) -- All `getTenantFromContext` error paths (`handlers.go:60,89,129,169,234,275,322,357`) now route through `h.sendErrorResponse` instead of calling `httputil.SendErrorResponse` directly.

5. ~~`NewTypeValidator` captures mutable slice~~ (Shared, P3) -- Now uses `slices.Clone(types)` before closing over the slice (`validate/validate.go:20-23`).

6. ~~Typo "tyep" in doc comment~~ (Shared, P3) -- Fixed to "type" (`httputil/http.go:39`).

7. **Tenants DTO file reorganization** (Tenants) -- `response.go` and `request.go` split into per-entity files: `tenant.go`, `data_source.go`, `data_source_attributes.go`, `lookup.go`. Purely structural, no functional changes.

8. ~~Incomplete error-to-HTTP mapping in tenants~~ (Tenants, P3) -- `ErrDataSourceAttrParse`, `ErrInvalidSourceType`, and `ErrInvalidSourceTypeAttributes` now map to 400 via `isBadRequest()`. `ErrStrategyNotFound` moved to Will Not Fix (see below). *(Resolved from 4/16 #15, 4/17 #53, 4/18 #51, 4/19 #39, 4/20 #34, 4/22 #35.)*

9. **MongoDB connection infrastructure added** (All services) -- `Connect()` function with functional options pattern (`WithHost`, `WithPort`, `WithUsername`, `WithPassword`) added to all three services. Tenants service has full MongoDB repository implementations.

10. **DI/persistence bootstrap pattern implemented** (All services) -- `persistence.go` in each service now supports switching between `in-memory` and `mongodb` drivers via `settings.json` configuration.

11. **Forms field attribute response DTOs added** (Forms) -- Explicit response structs with json tags for all five field attribute types (`TextFieldAttributeResponse`, `NumberFieldAttributeResponse`, etc.) and a `fieldAttributesToResponse` mapper using type-switch (`dto/field_attribute.go`).

12. **Tenants data source attribute response DTOs added** (Tenants) -- Response structs with json tags for all three attribute types and a `dataSourceAttributesToResponse` mapper (`dto/data_source_attributes.go`).

13. **`MongoDBDatabase` implements `database.Database` interface** (Tenants) -- Full implementation with `BeginTx`/`CommitTx`/`RollbackTx` backed by MongoDB sessions, and `Close` method (`mongodb/mongodb_database.go`).

14. **MongoDB document mappers implemented** (Tenants) -- `tenantDocument`, `dataSourceDocument` BSON structs with bidirectional mapping functions and attribute strategy deserialization (`mongodb/documents.go`).

15. ~~`mongo.ErrNoDocuments` not translated to `common.ErrNotFound`~~ (Tenants, P1) -- Generic `findById` method in `base_repository.go` now checks `errors.Is(err, mongo.ErrNoDocuments)` and returns `common.ErrNotFound`. Applied to both tenants and data sources repos. *(Resolved from 4/23 #29.)*

16. ~~`MongoDBDatabase.Close()` is a no-op~~ (Tenants, P1) -- `Close()` now calls `db.client.Disconnect(context.Background())` (`mongodb_database.go`). *(Resolved from 4/23 #38.)*

17. ~~Error message says "field type" instead of "data source type"~~ (Tenants, P3) -- `dto/data_source_attributes.go` now reads `"data source type is required"`. *(Resolved from 4/23 #40.)*

18. ~~Exported struct types with constructors returning interfaces~~ (Tenants, P3) -- `MongoDBTenantsRepository`, `MongoDBDataSourcesRepository`, `InmemoryTenantRepository`, and `InmemoryDataSourceRepository` are now unexported with unexported constructors. Only the `Bootstrap` factory functions are exported. *(Resolved from 4/23 #41.)*

19. ~~`Inmemory` naming inconsistency~~ (Tenants, P3) -- `InmemoryTenantRepository` renamed to `inMemoryTenantRepository`, `InmemoryDataSourceRepository` renamed to `inMemoryDataSourceRepository`. Now consistent with `InMemory` casing in forms, submissions, and shared package. *(Resolved from 4/22, 4/23 #42.)*

20. ~~`base_repository.go` is a 4-line file with no methods~~ (Tenants, P3) -- Now a generic `mongodbBaseRepository[T any]` with shared `findById` and `exists` methods. Both MongoDB repos embed the generic type and use these methods, eliminating duplicated boilerplate. *(Resolved from 4/23 #43.)*

21. ~~`omitempty` on `_id` BSON tag is a latent risk~~ (Tenants, P3) -- Removed from both `tenantDocument` and `dataSourceDocument`. Tags now read `bson:"_id"`. *(Resolved from 4/23 #44.)*

22. ~~`MongoDBDatabase` stores unused `db *mongo.Database` field~~ (Tenants, P3) -- The `db` field has been removed. `BeginTx` now uses `db.client.StartSession()` directly. Constructor accepts `_ *mongo.Database` (ignored). *(Resolved from 4/23 #45.)*

23. ~~`DataSourceAttributes` concrete types lack json tags~~ (Tenants, P3) -- The DTO adapter layer now defines explicit response structs with json tags for all three attribute types (`dto/data_source_attributes.go`), with a `dataSourceAttributesToResponse` mapper. Domain types correctly do not carry serialization tags -- the adapter layer owns that concern, which is the proper hexagonal architecture approach. *(Resolved from 4/16 #18, 4/17 #56, 4/18 #54, 4/19 #42, 4/20 #37, 4/22 #38, 4/23 #31.)*

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

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:330`, `365`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments but remain unresolved. *(Unresolved from 4/13 #3, 4/17 #1, 4/18 #1, 4/19 #1, 4/20 #1, 4/22 #1, 4/23 #1.)*

#### Architectural

2. **`Find()` has no tenant filtering** (`forms_service.go:29-31`) -- Returns all forms across all tenants. Every other query enforces tenant isolation. The `getForms` handler (`handlers.go:30-55`) is also the only handler that never calls `getTenantFromContext`, meaning tenant context is never even read on this path. *(Unresolved from 4/13 #10, 4/17 #5, 4/18 #5, 4/19 #2, 4/20 #2, 4/22 #2, 4/23 #2.)*

3. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/13 #11, 4/17 #6, 4/18 #6, 4/19 #4, 4/20 #3, 4/22 #3, 4/23 #3.)*

4. **`time.Now()` called directly in the repository and service layers** (`forms_repository.go:59`, `83`; `forms_service.go:210`, `238`) -- Should be injected via a `Clock` interface or function. *(Unresolved from 4/13 #26, 4/17 #8, 4/18 #8, 4/19 #5, 4/20 #4, 4/22 #4, 4/23 #4.)*

5. **Redundant double-fetch in `Update()`** (`forms_service.go:78-86`) -- `isValidAccess()` fetches the form from the repository to check tenant ownership, then `Update()` immediately fetches the same form again via `FindById()`. The form should be fetched once and reused. *(Unresolved from 4/20 #5, 4/22 #5, 4/23 #5.)*

#### Missing Functionality

6. **Domain validation partially unimplemented** (`form.go:29,42`, `version.go:48`, `page.go:31`, `section.go:31`) -- `NewField` and `NewRule` now validate their types, but `NewForm`, `NewVersion`, `NewPage`, and `NewSection` constructors still contain `// TODO: Implement domain specific validation`. *(Partially resolved from 4/13 #13, 4/17 #10, 4/18 #11, 4/19 #7, 4/20 #7, 4/22 #7, 4/23 #6.)*

7. **Incomplete error-to-HTTP mapping** -- `ErrMissingTenantID`, `ErrVersionLocked`, `ErrDuplicatePosition`, `ErrDuplicateRuleType`, `ErrPublishedByRequired`, `ErrRetiredByRequired`, `ErrInvalidFieldType`, `ErrInvalidFieldAttributes`, and `ErrInvalidRuleType` all fall through to the default 500 case in `common.SendErrorResponse`. The service-level `sendErrorResponse` (`handlers.go:409-413`) is an empty switch that delegates directly to `httputil.SendErrorResponse`. *(Unresolved from 4/13 #15, 4/17 #12, 4/18 #13, 4/19 #9, 4/20 #9, 4/22 #9, 4/23 #7.)*

8. **No `Delete` operation for forms.** No delete handler, service method, or repository method exists. *(Unresolved from 4/13 #17, 4/17 #14, 4/18 #15, 4/19 #11, 4/20 #11, 4/22 #11, 4/23 #8.)*

#### Code Quality

9. **Inconsistent constructor signatures** -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` return `(*Entity, error)` but never return errors (validation is TODO). Either implement validation or simplify the signature. `NewField` and `NewRule` now correctly use the error return. *(Partially resolved from 4/13 #23, 4/17 #18, 4/18 #19, 4/19 #14, 4/20 #12, 4/22 #12, 4/23 #9.)*

10. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. Since `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, it falls through to the 500 default. Should be 400. *(Unresolved from 4/19 #15, 4/20 #13, 4/22 #13, 4/23 #10.)*

11. **`publishVersion` and `retireVersion` discard the returned `*domain.Version`** (`handlers.go:347-350`, `382-385`) -- Both handlers call the service, receive an updated `*domain.Version`, but respond with `Data: nil`. The mutated version state is never returned to the client. *(Unresolved from 4/20 #14, 4/22 #14, 4/23 #11.)*

12. **Map iteration order non-deterministic in DTO response mappers** (`dto/version.go:34-36`, `dto/page.go:77-79`, `dto/section.go`) -- `VersionToResponse()`, `PageToResponse()`, and `SectionToResponse()` iterate over Go maps. The order of items in JSON array responses will be non-deterministic across requests. Pages, sections, and fields should be sorted by position key. *(Unresolved from 4/20 #15, 4/22 #15, 4/23 #12.)*

13. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:242-246`) -- `createVersion` calls `ReadValidateJsonPayload(r, &body)` where `body` is `CreateVersionDto struct{}`. The request body is read and decoded into a type with no fields, meaning any payload is silently discarded. Either the DTO should carry fields or the `ReadValidateJsonPayload` call should be removed. *(Unresolved from 4/22 #16, 4/23 #13.)*

14. **`FindVersions` returns `ErrNotFound` when form has no versions** (`forms_repository.go:93-96`) -- If a form exists but has no entry in the versions map (which happens when a form is created before any version is added), `FindVersions` returns `common.ErrNotFound` rather than an empty slice. This will produce a 404 when listing versions of a new form. *(Unresolved from 4/23 #14.)*

---

### Submissions Service

#### Bugs

15. **Forms and Submissions `bootstrapMongoDB` connects but returns empty `&ports.Repository{}`** (`persistence.go:57-68`) -- `bootstrapMongoDB` calls `mongodb.Connect(...)` but assigns the returned client to `_` and returns `&ports.Repository{}` with nil interface fields. Any repository call will panic with a nil pointer dereference when the `mongodb` driver is selected. The tenants service correctly passes the client to `mongodb.Bootstrap(client, logger)`. *(New.)*

#### Architectural

16. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. *(Unresolved from 4/17 #27, 4/18 #25, 4/19 #18, 4/20 #16, 4/22 #17, 4/23 #15.)*

17. **Four handler stubs return 200 OK with empty body** (`handlers.go:87-93`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered in the router but have empty function bodies. They return HTTP 200 with zero-length body and no `Content-Type` header. These should either return 501 Not Implemented or not be registered. *(Unresolved from 4/20 #17, 4/22 #18, 4/23 #16.)*

#### Missing Functionality

18. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:65-71`) -- Return `nil, nil` and `nil` respectively. There is also no `SubmissionAttemptsRepository` in the secondary ports to back these methods. *(Unresolved from 4/17 #31, 4/18 #30, 4/19 #22, 4/20 #19, 4/22 #20, 4/23 #17.)*

19. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. No request DTOs exist for create/replay operations. *(Unresolved from 4/17 #32, 4/18 #31, 4/19 #23, 4/20 #20, 4/22 #21, 4/23 #18.)*

20. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, and no business methods. *(Unresolved from 4/17 #33, 4/18 #32, 4/19 #24, 4/20 #21, 4/22 #22, 4/23 #19.)*

21. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindById`, `FindByReferenceId`. No `Create`, `Update`, or `Delete`. *(Unresolved from 4/17 #34, 4/18 #33, 4/19 #25, 4/20 #22, 4/22 #23, 4/23 #20.)*

22. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields, making it impossible to specify what to replay. *(Unresolved from 4/18 #36, 4/19 #28, 4/20 #25, 4/22 #26, 4/23 #21.)*

#### Code Quality

23. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` is also typed as `any`. *(Unresolved from 4/17 #38, 4/18 #37, 4/19 #29, 4/20 #26, 4/22 #27, 4/23 #22.)*

24. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` is declared but no `const` block with valid status values exists. *(Unresolved from 4/17 #39, 4/18 #38, 4/19 #30, 4/20 #27, 4/22 #28, 4/23 #23.)*

25. **`SubmissionsRepository.FindByReferenceId` does a linear scan** (`submissions_repository.go`) -- Iterates over all entries comparing `ReferenceID`. No secondary index. *(Unresolved from 4/17 #40, 4/18 #39, 4/19 #31, 4/20 #28, 4/22 #29, 4/23 #24.)*

26. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When the request context is cancelled, the `select` on `r.Context().Done()` returns without writing any HTTP response. The client receives a connection drop with no status code. *(Unresolved from 4/23 #25.)*

27. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`) -- The map key is `string`, but the domain uses `domain.SubmissionID` (a named `string` type). `FindById` does `r.submissions[string(id)]`. The map should be `map[domain.SubmissionID]*domain.Submission` for type safety. *(Unresolved from 4/23 #26.)*

---

### Tenants Service

#### Bugs

28. **MongoDB `Upsert` overwrites `CreatedAt` with zero time on update path** (`tenants_repository.go:94-102`; `data_sources_repository.go:100-108`) -- On the update branch (`t.ID != ""`), `CreatedAt` is not set. The `$set` operation via `toTenantDocument(t)` / `toDataSourceDocument(ds)` writes the zero value of `time.Time` to `created_at`. The in-memory implementation preserves `CreatedAt` from the existing record, but the MongoDB implementation does not. The `// TODO: Move to the domain layer` comments acknowledge the underlying issue but the current code is actively destructive on updates. *(New.)*

#### Architectural

29. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:25-27`) -- Returns every tenant in a single unbounded response. `ListDataSourceQuery` in `query.go` is an empty struct with a `// TODO: Add pagination support` comment. *(Unresolved from 4/16 #10, 4/17 #48, 4/18 #45, 4/19 #34, 4/20 #29, 4/22 #30, 4/23 #27.)*

30. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`) -- When a tenant is removed, only the tenant record is deleted. Any `DataSource` records associated with that tenant remain orphaned in the data sources store. There is no cascade delete and no service-level cleanup. *(Unresolved from 4/20 #30, 4/22 #31, 4/23 #28.)*

31. **UUID generation and timestamp management duplicated across persistence layers with inconsistent behavior** (`tenants_repository.go:98-101`; `data_sources_repository.go:104-107`; `inmemory/tenant_repository.go:67-69`; `inmemory/data_sources_repository.go:73-76`) -- The same domain concern (identity assignment, lifecycle timestamps) is now implemented in both in-memory and MongoDB repositories with different behavior: in-memory preserves `CreatedAt` on update, MongoDB overwrites it with zero time. ID generation is a domain concern that should live in domain factories, not persistence adapters. The `// TODO: Move to the domain layer` comments acknowledge this. *(Unresolved from 4/23 #28, expanded with MongoDB duplication.)*

#### Missing Functionality

32. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` returns `(*Tenant, error)` but never validates or returns an error. `NewDataSource` now validates type and attribute-type agreement but does not validate empty `TenantID` or field lengths. *(Partially resolved from 4/16 #13, 4/17 #51, 4/18 #49, 4/19 #37, 4/20 #32, 4/22 #33, 4/23 #29.)*

33. **`Lookup` service method is a stub** (`data_sources_service.go:131-145`) -- Returns `nil, nil` after verifying the data source exists. Contains `// TODO: Implement data source lookup strategy pattern`. *(Unresolved from 4/16 #17, 4/17 #55, 4/18 #53, 4/19 #41, 4/20 #36, 4/22 #37, 4/23 #30.)*

34. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup(value, label)` constructor performs no validation and has no json tags on the struct fields. *(Unresolved from 4/16 #19, 4/17 #57, 4/18 #55, 4/19 #43, 4/20 #38, 4/22 #39, 4/23 #32.)*

#### Code Quality

35. **`time.Now()` called directly in the repository layer** (`inmemory/tenant_repository.go:64`; `inmemory/data_sources_repository.go:69`; `tenants_repository.go:101`; `data_sources_repository.go:107`). Now present in both in-memory and MongoDB repositories. *(Unresolved from 4/16 #24, 4/17 #60, 4/18 #58, 4/19 #45, 4/20 #40, 4/22 #41, 4/23 #33.)*

36. **`Ping` uses `context.Background()` with no timeout** (`mongodb/mongodb.go:38`) -- If MongoDB is unreachable, the service hangs forever at startup. Should use `context.WithTimeout`. *(New.)*

---

### Cross-Service

#### Architectural

37. **Triplicated MongoDB connection code** (`forms/.../mongodb/mongodb.go`; `submissions/.../mongodb/mongodb.go`; `tenants/.../mongodb/mongodb.go`) -- `MongoDBOpts` struct, `Connect()` function, and all `With*` option functions are copy-pasted identically across all three services (~71 lines each). Should be extracted to `pkg/common/database/` alongside the `Database` interface and `InMemoryDatabase`. *(New.)*

38. **Triplicated `persistence.go` bootstrap pattern** (`forms/.../persistence/persistence.go`; `submissions/.../persistence/persistence.go`; `tenants/.../persistence/persistence.go`) -- `PersistenceDriver` type, constants, `PersistenceSettings`, `PersistenceOptions`, `Bootstrap()`, and `parseOptions[T]()` are nearly identical across all three services. The only difference is which `inmemory.Bootstrap` and `mongodb.Bootstrap` they call. Should be generalized. *(New.)*

39. **`MongoDBDatabase` only implemented for tenants** -- Forms and submissions have no `MongoDBDatabase` implementation. Their `bootstrapMongoDB` returns a `Repository` with a nil `Database` field, which will panic when `Application.Close()` calls `repository.Database.Close()`. *(New.)*

40. **Zero test files** in all three services and `pkg/common/`. *(Unresolved from 4/13 #12, 4/17 #9, 4/18 #10, 4/19 #6, 4/20 #6, 4/22 #6 (Forms); 4/17 #30, 4/18 #29, 4/19 #21, 4/20 #18, 4/22 #19 (Submissions); 4/16 #12, 4/17 #50, 4/18 #48, 4/19 #36, 4/20 #31, 4/22 #32 (Tenants); 4/23 #34.)*

41. **No domain events** for cross-service communication. *(Unresolved from 4/13 #14, 4/17 #11, 4/18 #12, 4/19 #8, 4/20 #8, 4/22 #8 (Forms); 4/17 #35, 4/18 #34, 4/19 #26, 4/20 #23, 4/22 #24 (Submissions); 4/16 #14, 4/17 #52, 4/18 #50, 4/19 #38, 4/20 #33, 4/22 #34 (Tenants); 4/23 #35.)*

42. **No real authentication** -- `X-Tenant-ID` is blindly trusted across all services. *(Unresolved from 4/13 #16, 4/17 #13, 4/18 #14, 4/19 #10, 4/20 #10, 4/22 #10 (Forms); 4/17 #36, 4/18 #35, 4/19 #27, 4/20 #24, 4/22 #25 (Submissions); 4/16 #16, 4/17 #54, 4/18 #52, 4/19 #40, 4/20 #35, 4/22 #36 (Tenants); 4/23 #36.)*

43. **No graceful shutdown** (all services, e.g. submissions `cmd/server/main.go:58-60`) -- `server.ListenAndServe()` blocks until error, and `log.Fatal` calls `os.Exit`, preventing `defer app.Close()` from executing. There is no signal handling (`os.Signal`) or `server.Shutdown(ctx)` call. *(Unresolved from 4/23 #37.)*

---

### Shared Package

44. **500 errors leak `err.Error()` to clients** (`httputil/http.go:103-107`) -- The `default` case in `SendErrorResponse` returns the raw error string in the JSON response body. In production, this could expose internal details (database errors, file paths, stack traces). Should return a generic message and log the real error server-side. *(Unresolved from 4/23 #38.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 15 | `bootstrapMongoDB` returns empty `Repository{}` -- crash on any repo call | Forms, Submissions |
| **P1** | 28 | MongoDB `Upsert` overwrites `CreatedAt` with zero time on update | Tenants |
| **P2** | 2, 16 | `Find()` has no tenant filtering | Forms, Submissions |
| **P2** | 6, 32 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 20 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 10 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 30 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 17 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 5 | Redundant double-fetch in `Update()` | Forms |
| **P2** | 44 | 500 errors leak `err.Error()` to clients | Shared |
| **P2** | 31 | UUID/timestamp management duplicated with inconsistent behavior | Tenants |
| **P2** | 37 | Triplicated MongoDB connection code | All |
| **P2** | 38 | Triplicated `persistence.go` bootstrap pattern | All |
| **P2** | 39 | `MongoDBDatabase` only implemented for tenants | Forms, Submissions |
| **P3** | 40 | Zero test files | All |
| **P3** | 41 | No domain events | All |
| **P3** | 7 | Incomplete error-to-HTTP mapping | Forms |
| **P3** | 23 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4, 35 | `time.Now()` in repository/service layers | Forms, Tenants |
| **P3** | 12 | Map iteration non-deterministic in DTO mappers | Forms |
| **P3** | 11 | `publishVersion`/`retireVersion` discard returned version | Forms |
| **P3** | 13 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 14 | `FindVersions` returns 404 for empty version list | Forms |
| **P3** | 43 | No graceful shutdown | All |
| **P3** | 36 | `Ping` uses `context.Background()` with no timeout | All |

---

## Summary

### Progress Since 4/22

Eight commits landed since the initial 4/23 review was written, representing the most significant infrastructure push of the project to date:

- **MongoDB connection infrastructure** established across all three services with a functional options pattern (`Connect`, `WithHost`, `WithPort`, `WithUsername`, `WithPassword`).
- **Tenants service MongoDB repositories** fully implemented with CRUD operations for both tenants and data sources, BSON document mapping, attribute strategy deserialization, and a `MongoDBDatabase` struct implementing the `database.Database` interface with real MongoDB session-based transactions.
- **DI/persistence bootstrap** updated in all three services to support switching between `in-memory` and `mongodb` drivers via `settings.json` configuration.
- **Forms field attribute response DTOs** added with explicit json-tagged response structs for all five attribute types and a type-switch mapper.
- **Tenants data source attribute response DTOs** added with the same pattern.
- **Forms naming convention fixes** applied.
- **Tenants DTO reorganization** continued with `data_source_lookup.go` renamed to `lookup.go`.

Additionally, unstaged changes address several issues:

- **`mongo.ErrNoDocuments` now maps to `common.ErrNotFound`** (Tenants) -- Generic `findById` in `base_repository.go` translates the driver error, applied to both repos.
- **`MongoDBDatabase.Close()` now disconnects the client** (Tenants) -- Calls `db.client.Disconnect(context.Background())`.
- **Copy-paste "field type" error message fixed** (Tenants) -- Now reads "data source type is required".
- **MongoDB and in-memory repo types unexported** (Tenants) -- Struct types and constructors are now unexported; only `Bootstrap` factory functions are exported. Idiomatic Go convention.
- **`Inmemory` naming inconsistency fixed** (Tenants) -- Renamed to `inMemory`, consistent with forms, submissions, and shared package.
- **`base_repository.go` now generic with shared methods** (Tenants) -- `mongodbBaseRepository[T any]` with `findById` and `exists`, eliminating duplicated boilerplate in both MongoDB repos.
- **`omitempty` removed from `_id` BSON tags** (Tenants) -- Both document structs now use `bson:"_id"`.
- **Unused `db *mongo.Database` field removed from `MongoDBDatabase`** (Tenants) -- `BeginTx` now uses `db.client.StartSession()` directly.
- **`DataSourceAttributes` json serialization handled by DTO layer** (Tenants) -- Explicit response structs with json tags added in `dto/data_source_attributes.go`. Domain types correctly omit serialization tags per hexagonal architecture.

### New Issues Found

1. **Forms/Submissions `bootstrapMongoDB` returns empty `Repository{}`** (#15, P0) -- Client is discarded; any MongoDB-backed call will panic.
2. **MongoDB `Upsert` overwrites `CreatedAt` with zero time** (#28, P1) -- Destructive on every update operation.
3. **Triplicated MongoDB connection code** (#37, P2) -- ~71 identical lines in each service.
4. **Triplicated `persistence.go` bootstrap** (#38, P2) -- Nearly identical factory pattern in each service.
5. **`MongoDBDatabase` only implemented for tenants** (#39, P2) -- Forms/submissions have no MongoDB Database implementation.
6. **UUID/timestamp duplicated with inconsistent behavior** (#31, P2) -- In-memory preserves `CreatedAt`, MongoDB does not.
7. **`Ping` with no timeout** (#36, P3) -- Services hang forever if MongoDB unreachable.

### Current State

**Forms Service** gained explicit field attribute response DTOs and naming convention fixes. The service remains the most functionally complete, but its `bootstrapMongoDB` is non-functional (discards the client). All prior issues remain open.

**Tenants Service** saw the largest improvement with full MongoDB repository implementations, document mappers, and a `MongoDBDatabase` with real transaction support. Unstaged changes further improve code quality: `ErrNoDocuments` now maps to `ErrNotFound`, `Close()` properly disconnects, repo types are unexported, naming is consistent, `base_repository.go` is a useful generic, and `omitempty` is removed from BSON `_id` tags. The remaining MongoDB-layer bug is `CreatedAt` destruction on update. Code duplication between in-memory and MongoDB repositories (UUID generation, timestamp management) with inconsistent behavior is a growing concern.

**Submissions Service** is unchanged functionally. Its `bootstrapMongoDB` has the same empty-Repository problem as forms. The domain model remains entirely anemic.

**Shared Package** (`pkg/common`) is unchanged. The triplicated MongoDB connection code and bootstrap pattern across all three services should be extracted here.

### Highest-Impact Improvements

1. **Fix `bootstrapMongoDB` in forms and submissions** (P0 -- both crash if mongodb driver is selected)
2. **Preserve `CreatedAt` on MongoDB update path** (P1 -- data corruption on every update)
3. **Extract shared MongoDB connection code to `pkg/common/database/`** (P2 -- eliminate ~140 lines of duplication)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add tenant filtering to `Find()` methods** across forms and submissions (P2 -- data isolation)
