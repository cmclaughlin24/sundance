# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/24 Review

1. ~~`Find()` has no tenant filtering~~ (Forms, P2) -- `getForms` handler (`handlers.go:32-38`) now calls `getTenantFromContext` and passes a `FindFormsQuery` to `Find()`. `FormsService.Find` (`forms_service.go:29-35`) now accepts `*ports.FindFormsQuery`, validates it, and passes `FormFilters{TenantID: query.TenantID}` to the repository. `FormsRepository.Find` (`ports/secondary.go:18`) now accepts `*FormFilters`. The MongoDB implementation (`mongodb/forms_repository.go:36-55`) applies the tenant filter to the BSON query. Tenant isolation is now enforced on the forms `Find` path. *(Resolves 4/24 #2.)*

2. ~~`publishVersion` and `retireVersion` discard the returned `*domain.Version`~~ (Forms, P3) -- Both handlers now pass `dto.VersionToResponse(res.data)` in the `Data` field (`handlers.go:354`, `389`) instead of `nil`. *(Resolves 4/24 #9.)*

3. ~~500 errors leak `err.Error()` to clients~~ (Shared, P2) -- The `default` case in `SendErrorResponse` (`httputil/http.go:97-100`) now returns `"An unexpected error occurred. Please contact support if the issue persists."` instead of `err.Error()`. *(Resolves 4/24 #34.)*

4. ~~`ReadJsonFile` and `FindById` naming conventions~~ (Shared, P3) -- `ReadJsonFile` renamed to `ReadJSONFile` (`pkg/common/utils.go:8`). All call sites updated. `FindById` no longer appears anywhere in the codebase; all ports and repositories use `FindByID`. Private handler helpers renamed: `getFormIdPathValue` → `getFormIDPathValue`, `getVersionIdPathValue` → `getVersionIDPathValue`, `getReferenceIdPathValue` → `getReferenceIDPathValue`. *(Resolves 4/24 #35.)*

5. ~~`mongodb.Bootstrap` returns empty `Repository{}`~~ (Submissions, P0) -- `Bootstrap` (`submissions/.../mongodb/mongodb.go:10-17`) now returns a fully wired `*ports.Repository` with `Database` via `database.NewMongoDBDatabase(client, db)` and `Submissions` via `newMongoDBSubmissionsRepository(db, logger)`. Service startup no longer panics. *(Resolves 4/24 #11. Note: repository methods are stubs -- see new issue #1.)*

6. ~~Previously unstaged changes now committed~~ -- The following items listed as "unstaged" in the 4/24 review are now committed: Go naming conventions (`APIResponse`, `SendJSONResponse`, etc.), forms `isBadRequest` error-to-HTTP mapping, DTO sort order via `slices.Sorted(maps.Keys(...))`, `Position` encapsulation via `withPosition`, `HydratePage`/`HydrateSection` child map initialization, `HydrateVersion` `pages` map initialization, `GetDataSourceLookupsQuery` CQRS fix, `FindNextVersionNumber` session handling, dead timestamp parameters removed, redundant double-fetch eliminated, `FindVersions` adapter consistency, `newMongoDBFormsRepository` unexported, `cursor.Close` context fix, `Tenant.Update()` nil guard added with `ErrInvalidTenant` sentinel error.

---

## Will Not Fix

See [4/24 review](code-review-4-24-26.md) for the full Will Not Fix list (9 items).

---

## New Issues

1. **Submissions MongoDB repository methods are stubs with no document mapping** (`submissions_repository.go:29-37`) -- `Find`, `FindByID`, and `FindByReferenceID` all return `nil, nil`. The underlying `MongoDBRepository` uses `any` as its type parameter instead of a concrete document struct (unlike forms and tenants which use typed documents with BSON tags and domain mappers). `FindByID` and `FindByReferenceID` returning `nil, nil` will cause nil pointer dereference in the service layer when it accesses `submission.TenantID` (`submissions_service.go:44,60`). Needs: document structs with BSON tags, BSON-to-domain mappers, and actual query implementations. *(New, P0.)*

2. **`toVersionDocument` and `rulesToDocuments` iterate over maps -- non-deterministic order** (`documents.go:59`, `documents.go:309`) -- `toVersionDocument` iterates over `v.GetPages()` (a `map[int]*Page`) to build the pages slice. Map iteration in Go is non-deterministic, so the page order in the BSON document varies across writes. Similarly, `rulesToDocuments` iterates over a `map[RuleType]*Rule`. While MongoDB doesn't guarantee array order for queries, inconsistent document order can cause unnecessary writes when using `$set`. *(New, P3.)*

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:339`, `370`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/24 #1.)*

#### Architectural

2. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. *(Unresolved from 4/24 #3.)*

3. **`time.Now()` called directly in the domain layer** (`form.go:29,55`; `version.go:48,125,143,165`) -- Should be injected via a `Clock` interface or function for testability. *(Unresolved from 4/24 #4.)*

#### Missing Functionality

4. **Domain validation partially unimplemented** (`form.go:29,42`, `version.go:48`, `page.go`, `section.go`) -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` constructors still contain `// TODO: Implement domain specific validation`. *(Unresolved from 4/24 #5.)*

5. **No `Delete` operation for forms.** *(Unresolved from 4/24 #6.)*

#### Code Quality

6. **Inconsistent constructor signatures** -- `NewForm`, `NewVersion`, `NewPage`, and `NewSection` return `(*Entity, error)` but never return errors. *(Unresolved from 4/24 #8.)*

7. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go`, `handlers.go:250-254`) -- `createVersion` calls `ReadValidateJSONPayload(r, &body)` where `body` is `CreateVersionRequest struct{}`. *(Unresolved from 4/24 #10.)*

---

### Submissions Service

#### Bugs

8. **MongoDB repository methods are stubs returning `nil, nil`** (`submissions_repository.go:29-37`) -- See new issue #1. Any `FindByID` or `FindByReferenceID` call with the MongoDB driver will panic with nil pointer dereference in the service layer. *(New, P0.)*

#### Architectural

9. **`Find()` has no tenant filtering** (`submissions_service.go:27-29`) -- Returns all submissions across all tenants. *(Unresolved from 4/24 #12.)*

10. **Four handler stubs return 200 OK with empty body** (`handlers.go:95-101`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered but have empty function bodies. *(Unresolved from 4/24 #13.)*

#### Missing Functionality

11. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:69-75`) -- Return `nil, nil` and `nil`. *(Unresolved from 4/24 #14.)*

12. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/24 #15.)*

13. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, no business methods. *(Unresolved from 4/24 #16.)*

14. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindByID`, `FindByReferenceID`. *(Unresolved from 4/24 #17.)*

15. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/24 #18.)*

#### Code Quality

16. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/24 #19.)*

17. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/24 #20.)*

18. **`FindByReferenceID` does a linear scan** in the in-memory repository. *(Unresolved from 4/24 #21.)*

19. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/24 #22.)*

20. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`). *(Unresolved from 4/24 #23.)*

---

### Tenants Service

#### Architectural

21. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/24 #24.)*

22. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`). *(Unresolved from 4/24 #25.)*

#### Missing Functionality

23. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` never validates. *(Unresolved from 4/24 #26.)*

24. **`Lookup` service method is a stub** (`data_sources_service.go`). *(Unresolved from 4/24 #27.)*

25. **`Lookup` value object has no validation** (`lookup.go`). *(Unresolved from 4/24 #28.)*

#### Code Quality

26. **`Ping` uses `context.Background()` with no timeout** (`pkg/database/mongodb.go:43`). *(Unresolved from 4/24 #29.)*

---

### Cross-Service

#### Architectural

27. **Zero test files** in all three services and shared packages. *(Unresolved from 4/24 #30.)*

28. **No domain events** for cross-service communication. *(Unresolved from 4/24 #31.)*

29. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/24 #32.)*

30. **No graceful shutdown** (all services). *(Unresolved from 4/24 #33.)*

---

### Shared Package

31. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/24 #7.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 8 | MongoDB repository methods are stubs -- nil panic on FindByID/FindByReferenceID | Submissions |
| **P2** | 9 | `Find()` has no tenant filtering | Submissions |
| **P2** | 4, 23 | Domain validation unimplemented | Forms, Tenants |
| **P2** | 13 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 31 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 22 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 10 | Empty handler stubs return 200 OK | Submissions |
| **P3** | 27 | Zero test files | All |
| **P3** | 28 | No domain events | All |
| **P3** | 16 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 3 | `time.Now()` in domain layer | Forms |
| **P3** | 7 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 30 | No graceful shutdown | All |
| **P3** | 26 | `Ping` uses `context.Background()` with no timeout | All |

---

## Summary

### Progress Since 4/24

The primary focus since the last review was committing the large set of previously-unstaged fixes, resolving tenant isolation in the forms service, and beginning submissions MongoDB persistence:

- **Forms `Find()` tenant filtering implemented** -- The `getForms` handler now extracts the tenant from context and passes a `FindFormsQuery` with `TenantID` to the service. The repository accepts `FormFilters` and applies the tenant ID to the MongoDB query. This closes the last tenant isolation gap in the forms service.
- **`publishVersion` and `retireVersion` now return version data** -- Both handlers populate the `Data` field with `dto.VersionToResponse(res.data)` instead of `nil`.
- **500 errors no longer leak internal error messages** -- `SendErrorResponse` default case returns a generic message. Internal error details are no longer exposed to clients.
- **Go naming conventions fully applied** -- `ReadJsonFile` → `ReadJSONFile`, `FindById` → `FindByID` across all ports, and private handler helpers (`getFormIdPathValue` → `getFormIDPathValue`, etc.) all corrected. No naming deviations remain in the codebase.
- **All previously-unstaged changes committed** -- 16 items from the 4/24 review that were listed as unstaged are now committed, including: Go naming conventions, forms error mapping, DTO sort ordering, domain encapsulation improvements (`withPosition`, `HydrateVersion` pages map init), CQRS naming fix, session handling, dead parameter removal, double-fetch elimination, adapter consistency, constructor visibility, cursor context fix, and `Tenant.Update()` nil guard with `ErrInvalidTenant`.
- **Submissions `mongodb.Bootstrap` P0 resolved** -- `Bootstrap` now wires both `Database` and `Submissions` fields. Service startup no longer panics. However, repository methods are stubs returning `nil, nil`, creating a new P0 (nil dereference in the service layer on `FindByID`/`FindByReferenceID`).

### Current State

**35 remaining issues → 31 remaining issues** (resolved 6 from 4/24, added 2 new).

**Forms Service** is now the most complete service. Tenant filtering is enforced on all paths. MongoDB persistence is fully implemented with typed documents and BSON mappers. Error-to-HTTP mapping covers 8 domain errors plus shared errors. The version state machine is robust with proper aggregate encapsulation (`withPosition`, private `pages` map, duplicate guards). Remaining gaps: hardcoded placeholder user ID, unimplemented domain validation TODOs, no delete operation, `time.Now()` coupling, and `CreateVersionRequest` empty struct.

**Tenants Service** remains stable with no changes this cycle. MongoDB repositories are fully implemented. `Tenant.Update()` now has a nil guard with `ErrInvalidTenant`. Remaining gaps: no validation in `NewTenant`, cascade-delete on tenant removal, `Lookup` stub.

**Submissions Service** saw the most structural progress but remains the weakest. The P0 empty `Repository{}` is resolved -- `Bootstrap` now wires both `Database` and a `mongoDBSubmissionsRepository`. However, all three repository methods are stubs returning `nil, nil`, and the repository uses `MongoDBRepository[any]` instead of typed document structs. This creates a new P0: `FindByID` and `FindByReferenceID` return nil submissions, causing nil pointer dereference when the service checks `submission.TenantID`. The domain remains anemic, the repository interface is read-only, four handler stubs return 200 OK, and no request DTOs exist.

**Hexagonal Architecture** -- Dependency direction remains correct throughout. No cross-adapter imports. Domain layers are pure. All previously-unstaged architectural fixes are now committed.

**CQRS-Lite** -- Commands and queries are well-separated. The `GetDataSourceLookupsCommand` → `GetDataSourceLookupsQuery` fix is now committed. Forms queries use composition (`FindVersionsQuery` embeds `FindFormsByIDQuery`).

**DDD** -- Domain encapsulation improved: `withPosition` embedding, `HydrateVersion`/`HydratePage`/`HydrateSection` all initialize child maps, `Tenant.Update()` has nil guard. No domain events exist. Submissions domain remains entirely anemic.

### Highest-Impact Improvements

1. **Implement submissions MongoDB repository methods** (P0 -- stub methods cause nil panic in service layer)
2. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
3. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
4. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
