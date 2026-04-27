# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/25 Review

1. ~~Domain validation unimplemented~~ (Tenants, P2) -- `NewTenant` (`tenant.go:23-36`) now calls `validate.ValidateStruct(t)` with `Name` tagged `validate:"required,notblank"`. The constructor returns an error on invalid input. *(Resolves 4/25 #20, previously tracked from 4/24 #26.)*

2. ~~`NewFormsService` does not assign `versionsRepository`~~ (Forms, P0) -- The constructor (`forms_service.go:21-27`) now includes `versionsRepository: repository.Versions`. All version-related operations are properly wired. *(Introduced and resolved same day.)*

3. ~~`isBadRequest` missing domain error mappings~~ (Forms, P2) -- `isBadRequest` (`handlers.go:460-478`) now includes `ErrInvalidVersion`, `ErrInvalidPosition`, `ErrInvalidForm`, `ErrFormHasActiveVersion`, `ErrInvalidPage`, and `ErrInvalidSection`. These domain validation errors now correctly map to 400 Bad Request instead of falling through to 500. *(Introduced and resolved same day.)*

4. ~~`VersionResponseDto` naming~~ (Forms, P3) -- Renamed to `VersionResponse` (`dto/version.go:15`). All call sites in `handlers.go` updated. *(Introduced and resolved same day.)*

5. ~~`Tenant.Update()` does not re-validate after mutation~~ (Tenants, P2) -- `Update` (`tenant.go:56-61`) now calls `validate.ValidateStruct(t)` after setting fields. A blank name is now rejected on the update path, matching the validation enforced in `NewTenant`. *(Introduced and resolved same day.)*

6. ~~`DataSource.Update()` does not re-validate after mutation~~ (Tenants, P2) -- `Update` (`data_source.go:106-108`) now calls `validate.ValidateStruct(ds)` after setting fields. *(Not previously tracked; resolved same day as tenant equivalent.)*

7. ~~`NewCreateVersionCommand` parameter `formId` naming~~ (Forms, P3) -- Renamed to `formID` (`ports/commands.go:64`). *(Introduced and resolved same day.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the full Will Not Fix list (10 items).

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:367`, `402`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/25 #1.)*

#### Code Quality

2. **`CreateVersionDto` is an empty struct deserialized from request body** (`dto/version.go:9`, `handlers.go:280-283`) -- `createVersion` calls `ReadValidateJSONPayload(r, &body)` where `body` is `CreateVersionRequest struct{}`. *(Unresolved from 4/25 #4.)*

---

### Submissions Service

#### Bugs

3. **MongoDB repository methods are stubs returning `nil, nil`** (`submissions_repository.go:26-37`) -- `Find`, `FindByID`, and `FindByReferenceID` all return `nil, nil`. `FindByID` and `FindByReferenceID` returning nil will cause nil pointer dereference in the service layer when it accesses `submission.TenantID` (`submissions_service.go:40,58`). Needs: document structs with BSON tags, BSON-to-domain mappers, and actual query implementations. *(Unresolved from 4/25 #5.)*

#### Architectural

4. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. The `getSubmissions` handler does not extract tenant ID from context. *(Unresolved from 4/25 #6.)*

5. **Four handler stubs return 200 OK with empty body** (`handlers.go:95-115`) -- `createSubmission`, `getSubmissionAttempts`, `getSubmissionStatus`, and `replaySubmission` are registered but have empty function bodies. *(Unresolved from 4/25 #7.)*

#### Missing Functionality

6. **`FindAttempts` and `Replay` service methods are stubs** (`submissions_service.go:63-71`) -- Return `nil, nil` and `nil`. *(Unresolved from 4/25 #8.)*

7. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/25 #9.)*

8. **No domain constructors** -- `Submission` and `SubmissionAttempt` are bare structs with no factory functions, no validation, no business methods. *(Unresolved from 4/25 #10.)*

9. **No write operations in the repository interface** -- `SubmissionsRepository` only defines `Find`, `FindByID`, `FindByReferenceID`. *(Unresolved from 4/25 #11.)*

10. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/25 #12.)*

#### Code Quality

11. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/25 #13.)*

12. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/25 #14.)*

13. **`FindByReferenceID` does a linear scan** in the in-memory repository. *(Unresolved from 4/25 #15.)*

14. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/25 #16.)*

15. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** (`submissions_repository.go:14`). *(Unresolved from 4/25 #17.)*

---

### Tenants Service

#### Architectural

16. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/25 #18.)*

17. **Tenant removal does not cascade-delete DataSources** (`tenants_service.go:75-87`). *(Unresolved from 4/25 #19.)*

#### Missing Functionality

18. **`Lookup` service method is a stub** (`data_sources_service.go:125-139`). *(Unresolved from 4/25 #21.)*

19. **`Lookup` value object has no validation** (`lookup.go`). *(Unresolved from 4/25 #22.)*

---

### Cross-Service

#### Architectural

20. **Zero test files** in all three services and shared packages. *(Unresolved from 4/25 #23.)*

21. **No domain events** for cross-service communication. *(Unresolved from 4/25 #24.)*

22. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/25 #25.)*

23. **No graceful shutdown** (all services). *(Unresolved from 4/25 #26.)*

---

### Shared Package

24. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/25 #27.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P0** | 3 | MongoDB repository methods are stubs -- nil panic on FindByID/FindByReferenceID | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 4 | `Find()` has no tenant filtering | Submissions |
| **P2** | 5 | Empty handler stubs return 200 OK | Submissions |
| **P2** | 8 | No domain constructors | Submissions |
| **P2** | 17 | Tenant removal doesn't cascade-delete DataSources | Tenants |
| **P2** | 24 | `ErrMissingTenantID` maps to 500 | Shared |
| **P3** | 20 | Zero test files | All |
| **P3** | 21 | No domain events | All |
| **P3** | 11 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 2 | `CreateVersionRequest` empty struct deserialized | Forms |
| **P3** | 23 | No graceful shutdown | All |

---

## Summary

### Progress Since 4/25

The focus since the last review was closing validation gaps and fixing a blocking regression from the aggregate boundary refactor:

- **Forms P0 regression fixed** -- `NewFormsService` (`forms_service.go:21-27`) now wires `versionsRepository: repository.Versions`. The aggregate boundary refactor that split `FormsRepository` and `VersionRepository` had left this field nil, meaning every version operation would panic. Resolved before any deployment impact.
- **Forms error-to-HTTP mapping completed** -- `isBadRequest` (`handlers.go:460-478`) now covers all domain validation errors: `ErrInvalidVersion`, `ErrInvalidPosition`, `ErrInvalidForm`, `ErrFormHasActiveVersion`, `ErrInvalidPage`, and `ErrInvalidSection`. No domain error should fall through to 500 anymore.
- **`VersionResponseDto` renamed to `VersionResponse`** -- Dropped the `Dto` suffix entirely rather than capitalizing to `DTO`, which is cleaner and more idiomatic Go. All call sites updated.
- **Tenants `NewTenant` validation implemented** -- `NewTenant` now validates via `validate.ValidateStruct()` with `Name` tagged `required,notblank`. This closes the last domain constructor without validation in the tenants service.
- **`Tenant.Update()` and `DataSource.Update()` now re-validate** -- Both `Update` methods call `validate.ValidateStruct()` after mutation, closing a bypass where invalid data could be persisted through the update path without hitting constructor validation.
- **Go naming convention fix** -- `NewCreateVersionCommand` parameter renamed from `formId` to `formID`.

### Current State

**25 remaining issues → 24 remaining issues** (resolved 7 including 1 from 4/25 and 6 introduced-and-resolved same day; no new unresolved issues).

**Forms Service** is now fully operational again after the P0 regression fix. The domain layer is the strongest in the codebase: all constructors validate, the version state machine is robust, position fields are encapsulated via `withPosition`, sorted iteration methods (`Get*Slice()`) are used throughout, and all domain errors map to appropriate HTTP status codes. Form and Version are cleanly modeled as separate aggregates with distinct repository interfaces. Remaining gaps: hardcoded placeholder user ID (P2), `CreateVersionRequest` empty struct (P3).

**Tenants Service** is now the most validation-complete service. All domain constructors validate (`NewTenant`, `NewDataSource`), and both `Update` methods re-validate after mutation. MongoDB repositories are fully implemented. Remaining gaps: no cascade-delete on tenant removal (P2), `Lookup` stub (P3), `Lookup` value object no validation (P3), no pagination on `Find()` (P3).

**Submissions Service** remains unchanged and is still the weakest service. The P0 MongoDB stubs persist: all three repository methods return `nil, nil`, causing nil dereference in the service layer on `FindByID`/`FindByReferenceID`. The domain is entirely anemic (no constructors, no validation, no status constants, no business methods). The repository interface is read-only. Four handler stubs return 200 OK with empty bodies. No request DTOs exist.

**Hexagonal Architecture** -- Dependency direction remains correct throughout. No cross-adapter imports. Domain layers are pure.

**CQRS-Lite** -- Commands and queries are well-separated. No violations.

**DDD** -- Domain encapsulation continues to strengthen in forms and tenants. Both services now have complete validation chains (constructor + update). Submissions domain remains entirely anemic. No domain events exist anywhere.

### Highest-Impact Improvements

1. **Implement submissions MongoDB repository methods** (P0 -- stub methods cause nil panic in service layer)
2. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
3. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
4. **Add test coverage** starting with `pkg/database` and domain layers (P3 -- long-term reliability)
