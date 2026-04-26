# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/24 Review

1. ~~`Find()` has no tenant filtering~~ (Forms, P2) -- `getForms` handler (`handlers.go:32-38`) now calls `getTenantFromContext` and passes a `FindFormsQuery` to `Find()`. `FormsService.Find` (`forms_service.go:29-35`) now accepts `*ports.FindFormsQuery`, validates it, and passes `FormFilters{TenantID: query.TenantID}` to the repository. `FormsRepository.Find` (`ports/secondary.go:18`) now accepts `*FormFilters`. The MongoDB implementation (`mongodb/forms_repository.go:36-55`) applies the tenant filter to the BSON query. Tenant isolation is now enforced on the forms `Find` path. *(Resolves 4/24 #2.)*

2. ~~`publishVersion` and `retireVersion` discard the returned `*domain.Version`~~ (Forms, P3) -- Both handlers now pass `dto.VersionToResponse(res.data)` in the `Data` field (`handlers.go:354`, `389`) instead of `nil`. *(Resolves 4/24 #9.)*

3. ~~500 errors leak `err.Error()` to clients~~ (Shared, P2) -- The `default` case in `SendErrorResponse` (`httputil/http.go:97-100`) now returns `"An unexpected error occurred. Please contact support if the issue persists."` instead of `err.Error()`. *(Resolves 4/24 #34.)*

4. ~~`ReadJsonFile` and `FindById` naming conventions~~ (Shared, P3) -- `ReadJsonFile` renamed to `ReadJSONFile` (`pkg/common/utils.go:8`). All call sites updated. `FindById` no longer appears anywhere in the codebase; all ports and repositories use `FindByID`. Private handler helpers renamed: `getFormIdPathValue` → `getFormIDPathValue`, `getVersionIdPathValue` → `getVersionIDPathValue`, `getReferenceIdPathValue` → `getReferenceIDPathValue`. *(Resolves 4/24 #35.)*

5. ~~`mongodb.Bootstrap` returns empty `Repository{}`~~ (Submissions, P0) -- `Bootstrap` (`submissions/.../mongodb/mongodb.go:10-17`) now returns a fully wired `*ports.Repository` with `Database` via `database.NewMongoDBDatabase(client, db)` and `Submissions` via `newMongoDBSubmissionsRepository(db, logger)`. Service startup no longer panics. *(Resolves 4/24 #11. Note: repository methods are stubs -- see new issue #1.)*

6. ~~`time.Now()` called directly in the domain layer~~ (Forms, Tenants, P3) -- All three domain packages (`forms`, `tenants`, `submissions`) now declare a package-level `var Now = time.Now` in `domain.go`. All call sites in `form.go`, `version.go`, `tenant.go`, and `data_source.go` replaced `time.Now()` with `Now()`. This enables test injection of a mock clock without requiring a `Clock` interface. *(Resolves 4/24 #4.)*

7. ~~`toVersionDocument` and `rulesToDocuments` iterate over maps -- non-deterministic order~~ (Forms, P3) -- New domain methods `GetPagesSlice()` (`version.go`), `GetSectionsSlice()` (`page.go`), and `GetFieldsSlice()` (`section.go`) return children sorted by position key via `slices.Sorted(maps.Keys(...))`. `toVersionDocument`, `toPageDocument`, and `toSectionDocument` (`documents.go`) now iterate over these sorted slices instead of maps. DTO converters (`VersionToResponse`, `PageToResponse`, `SectionToResponse`) also use the new `Get*Slice()` methods, removing their own `slices.Sorted(maps.Keys(...))` logic. *(Resolves 4/25 new issue #2.)*

9. ~~Aggregate boundaries unclear~~ (Forms, P3) -- Form and Version are now explicitly modeled as separate aggregates. The `FormsRepository` and `VersionRepository` port interfaces are split (`ports/secondary.go`), with `Repository` struct wiring both. Both persistence adapters (in-memory and MongoDB) implement the split: `InMemoryFormsRepository` and `InMemoryVersionsRepository` are separate structs with independent state, and `mongoDBFormsRepository` and `mongoDBVersionRepository` are separate structs backed by distinct collections. `FormsService` holds separate `formsRepository` and `versionsRepository` fields. A unique compound index on `(form_id, version)` in the MongoDB versions collection prevents duplicate version numbers per form. Cross-aggregate reference is by ID (`Version.FormID`). *(Resolves 4/24 #3.)*

10. ~~No `Delete` operation for forms~~ (Forms, P2) -- `DELETE /api/v1/forms/{formId}` route added (`routes.go`). `deleteForm` handler (`handlers.go`) extracts tenant and form ID, builds a `RemoveFormCommand`, and calls `FormsService.Delete`. The service (`forms_service.go`) validates the command, checks tenant access, checks for active versions via `hasActiveVersion` (returns `ErrFormHasActiveVersion` if any version has `VersionStatusActive`), then delegates to `FormsRepository.Delete`. Both adapters implement `Delete`: in-memory removes from the map, MongoDB delegates to `MongoDBRepository.Remove` which calls `DeleteOne`. A new `ErrFormHasActiveVersion` sentinel error is defined in `domain/form.go`. *(Resolves 4/24 #6.)*

11. ~~Domain validation partially unimplemented~~ (Forms, P2) -- All domain constructors (`NewForm`, `NewVersion`, `NewPage`, `NewSection`, `NewField`, `NewRule`) now validate invariants. Exported fields use `validate:"required,notblank"` struct tags via `go-playground/validator` (called through `validate.ValidateStruct`). Enum-typed fields (`VersionStatus`, `FieldType`, `RuleType`) use manual checks with `validate.NewTypeValidator` before struct validation, returning dedicated sentinel errors (`ErrInvalidVersionStatus`, `ErrInvalidFieldType`, `ErrInvalidRuleType`). `Field` additionally validates attribute-type alignment via `isValidFieldAttributes`. Position is validated manually via `isValidPosition` (`position >= 0`) since the `withPosition.position` field is unexported and invisible to the validator. `Rule.Expression` has `validate:"required,notblank"`. `Form.Update` re-validates after mutation. The `// TODO: Implement domain specific validation` comments are removed. *(Resolves 4/24 #5. Also resolves 4/24 #8 -- constructors now return meaningful errors.)*

8. ~~Previously unstaged changes now committed~~ -- The following items listed as "unstaged" in the 4/24 review are now committed: Go naming conventions (`APIResponse`, `SendJSONResponse`, etc.), forms `isBadRequest` error-to-HTTP mapping, DTO sort order via `slices.Sorted(maps.Keys(...))`, `Position` encapsulation via `withPosition`, `HydratePage`/`HydrateSection` child map initialization, `HydrateVersion` `pages` map initialization, `GetDataSourceLookupsQuery` CQRS fix, `FindNextVersionNumber` session handling, dead timestamp parameters removed, redundant double-fetch eliminated, `FindVersions` adapter consistency, `newMongoDBFormsRepository` unexported, `cursor.Close` context fix, `Tenant.Update()` nil guard added with `ErrInvalidTenant` sentinel error.

---

## Will Not Fix

See [4/24 review](code-review-4-24-26.md) for the full Will Not Fix list (9 items).

10. **`Ping` uses `context.Background()` with no timeout** (`pkg/database/mongodb.go:43`) -- The `Ping` call occurs once at startup to verify connectivity. A hung ping will be caught by deployment health checks or process-level timeouts. Adding a context timeout here adds complexity for minimal benefit. *(4/24 #29.)*

---

## New Issues

1. **Submissions MongoDB repository methods are stubs with no document mapping** (`submissions_repository.go:29-37`) -- `Find`, `FindByID`, and `FindByReferenceID` all return `nil, nil`. The underlying `MongoDBRepository` uses `any` as its type parameter instead of a concrete document struct (unlike forms and tenants which use typed documents with BSON tags and domain mappers). `FindByID` and `FindByReferenceID` returning `nil, nil` will cause nil pointer dereference in the service layer when it accesses `submission.TenantID` (`submissions_service.go:44,60`). Needs: document structs with BSON tags, BSON-to-domain mappers, and actual query implementations. *(New, P0.)*

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:339`, `370`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/24 #1.)*

#### Code Quality

4. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go`, `handlers.go:250-254`) -- `createVersion` calls `ReadValidateJSONPayload(r, &body)` where `body` is `CreateVersionRequest struct{}`. *(Unresolved from 4/24 #10.)*

---

### Submissions Service

#### Bugs

5. **MongoDB repository methods are stubs returning `nil, nil`** (`submissions_repository.go:29-37`) -- See new issue #1. Any `FindByID` or `FindByReferenceID` call with the MongoDB driver will panic with nil pointer dereference in the service layer. *(New, P0.)*

#### Architectural

6. **`Find()` has no tenant filtering** (`submissions_service.go:27-29`) -- Returns all submissions across all tenants. *(Unresolved from 4/24 #12.)*

7. **Four handler stubs return 200 OK with empty body** (`handlers.go:95-101`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered but have empty function bodies. *(Unresolved from 4/24 #13.)*

#### Missing Functionality

8. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:69-75`) -- Return `nil, nil` and `nil`. *(Unresolved from 4/24 #14.)*

9. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/24 #15.)*

10. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, no business methods. *(Unresolved from 4/24 #16.)*

11. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindByID`, `FindByReferenceID`. *(Unresolved from 4/24 #17.)*

12. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/24 #18.)*

#### Code Quality

13. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/24 #19.)*

14. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/24 #20.)*

15. **`FindByReferenceID` does a linear scan** in the in-memory repository. *(Unresolved from 4/24 #21.)*

16. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/24 #22.)*

17. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`). *(Unresolved from 4/24 #23.)*

---

### Tenants Service

#### Architectural

18. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/24 #24.)*

19. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:73-84`). *(Unresolved from 4/24 #25.)*

#### Missing Functionality

20. **Domain validation unimplemented** (`tenant.go:17-22`) -- `NewTenant` never validates. *(Unresolved from 4/24 #26.)*

21. **`Lookup` service method is a stub** (`data_sources_service.go`). *(Unresolved from 4/24 #27.)*

22. **`Lookup` value object has no validation** (`lookup.go`). *(Unresolved from 4/24 #28.)*

---

### Cross-Service

#### Architectural

23. **Zero test files** in all three services and shared packages. *(Unresolved from 4/24 #30.)*

24. **No domain events** for cross-service communication. *(Unresolved from 4/24 #31.)*

25. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/24 #32.)*

26. **No graceful shutdown** (all services). *(Unresolved from 4/24 #33.)*

---

### Shared Package

27. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/24 #7.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 5 | MongoDB repository methods are stubs -- nil panic on FindByID/FindByReferenceID | Submissions |
| **P2** | 6 | `Find()` has no tenant filtering | Submissions |
| **P2** | 20 | Domain validation unimplemented | Tenants |
| **P2** | 10 | No domain constructors | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 27 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 19 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 7 | Empty handler stubs return 200 OK | Submissions |
| **P3** | 23 | Zero test files | All |
| **P3** | 24 | No domain events | All |
| **P3** | 13 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 4 | `CreateVersionDto` empty struct deserialized | Forms |
| **P3** | 26 | No graceful shutdown | All |

---

## Summary

### Progress Since 4/24

The primary focus since the last review was committing the large set of previously-unstaged fixes, resolving tenant isolation in the forms service, beginning submissions MongoDB persistence, introducing an idiomatic clock pattern, ensuring deterministic child ordering, and clarifying aggregate boundaries:

- **Forms `Find()` tenant filtering implemented** -- The `getForms` handler now extracts the tenant from context and passes a `FindFormsQuery` with `TenantID` to the service. The repository accepts `FormFilters` and applies the tenant ID to the MongoDB query. This closes the last tenant isolation gap in the forms service.
- **`publishVersion` and `retireVersion` now return version data** -- Both handlers populate the `Data` field with `dto.VersionToResponse(res.data)` instead of `nil`.
- **500 errors no longer leak internal error messages** -- `SendErrorResponse` default case returns a generic message. Internal error details are no longer exposed to clients.
- **Go naming conventions fully applied** -- `ReadJsonFile` → `ReadJSONFile`, `FindById` → `FindByID` across all ports, and private handler helpers (`getFormIdPathValue` → `getFormIDPathValue`, etc.) all corrected. No naming deviations remain in the codebase.
- **All previously-unstaged changes committed** -- 16 items from the 4/24 review that were listed as unstaged are now committed, including: Go naming conventions, forms error mapping, DTO sort ordering, domain encapsulation improvements (`withPosition`, `HydrateVersion` pages map init), CQRS naming fix, session handling, dead parameter removal, double-fetch elimination, adapter consistency, constructor visibility, cursor context fix, and `Tenant.Update()` nil guard with `ErrInvalidTenant`.
- **Submissions `mongodb.Bootstrap` P0 resolved** -- `Bootstrap` now wires both `Database` and `Submissions` fields. Service startup no longer panics. However, repository methods are stubs returning `nil, nil`, creating a new P0 (nil dereference in the service layer on `FindByID`/`FindByReferenceID`).
- **Idiomatic clock pattern adopted across all domains** -- All three domain packages (`forms`, `tenants`, `submissions`) now declare `var Now = time.Now` in `domain.go`. All `time.Now()` calls in `form.go`, `version.go`, `tenant.go`, and `data_source.go` replaced with `Now()`, enabling test injection of a mock clock without a `Clock` interface.
- **Deterministic child ordering via `Get*Slice()` methods** -- New domain methods `GetPagesSlice()`, `GetSectionsSlice()`, and `GetFieldsSlice()` return children sorted by position key. Document converters and DTO converters now use these sorted slices instead of iterating maps directly, eliminating non-deterministic BSON document order and removing duplicated `slices.Sorted(maps.Keys(...))` logic from the DTO layer.
- **Form delete operation implemented** -- `DELETE /api/v1/forms/{formId}` route, `deleteForm` handler, `FormsService.Delete` with active version guard (`ErrFormHasActiveVersion`), and `FormsRepository.Delete` in both adapters. A generic `MongoDBRepository.Remove` method was added to the shared database package.
- **Form and Version aggregate boundaries clarified** -- `FormsRepository` and `VersionRepository` are now separate port interfaces (`ports/secondary.go`). Both the in-memory and MongoDB adapters implement the split with independent structs. `FormsService` holds separate `formsRepository` and `versionsRepository` fields. A unique compound index on `(form_id, version)` in the MongoDB versions collection prevents duplicate version numbers per form. `Version` references `Form` by ID only, making the separate-aggregate design explicit.

### Current State

**35 remaining issues → 25 remaining issues** (resolved 12 from 4/24, added 1 new from 4/25, resolved 1 new from 4/25, moved 1 to Will Not Fix).

**Forms Service** is now the most complete service. Tenant filtering is enforced on all paths. MongoDB persistence is fully implemented with typed documents and BSON mappers. Error-to-HTTP mapping covers 8 domain errors plus shared errors. The version state machine is robust with proper aggregate encapsulation (`withPosition`, private `pages` map, duplicate guards). Remaining gaps: hardcoded placeholder user ID, `CreateVersionRequest` empty struct.

**Tenants Service** remains stable with no changes this cycle. MongoDB repositories are fully implemented. `Tenant.Update()` now has a nil guard with `ErrInvalidTenant`. Remaining gaps: no validation in `NewTenant`, cascade-delete on tenant removal, `Lookup` stub.

**Submissions Service** saw the most structural progress but remains the weakest. The P0 empty `Repository{}` is resolved -- `Bootstrap` now wires both `Database` and a `mongoDBSubmissionsRepository`. However, all three repository methods are stubs returning `nil, nil`, and the repository uses `MongoDBRepository[any]` instead of typed document structs. This creates a new P0: `FindByID` and `FindByReferenceID` return nil submissions, causing nil pointer dereference when the service checks `submission.TenantID`. The domain remains anemic, the repository interface is read-only, four handler stubs return 200 OK, and no request DTOs exist.

**Hexagonal Architecture** -- Dependency direction remains correct throughout. No cross-adapter imports. Domain layers are pure. All previously-unstaged architectural fixes are now committed.

**CQRS-Lite** -- Commands and queries are well-separated. The `GetDataSourceLookupsCommand` → `GetDataSourceLookupsQuery` fix is now committed. Forms queries use composition (`FindVersionsQuery` embeds `FindFormsByIDQuery`).

**DDD** -- Domain encapsulation improved: `withPosition` embedding, `HydrateVersion`/`HydratePage`/`HydrateSection` all initialize child maps, `Tenant.Update()` has nil guard, idiomatic `var Now = time.Now` clock pattern across all domains, and `Get*Slice()` methods encapsulate sorted iteration. Form and Version are explicitly modeled as separate aggregates with distinct repository interfaces and cross-aggregate reference by ID. No domain events exist. Submissions domain remains entirely anemic.

### Highest-Impact Improvements

1. **Implement submissions MongoDB repository methods** (P0 -- stub methods cause nil panic in service layer)
2. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
3. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
4. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
