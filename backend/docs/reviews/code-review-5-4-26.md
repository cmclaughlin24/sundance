# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface / `Upsert` stubs~~ (Submissions, P2) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations are fully functional. The in-memory implementation stores by ID under a write lock. The MongoDB implementation uses `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a session. New `toSubmissionDocument` and `toSubmissionAttemptDocument` mapping functions handle domain-to-BSON conversion with `bson.Marshal` for the `Payload` and `ErrorDetails` fields. *(Resolves 4/28 #8.)*

3. ~~`CreateSubmissionCommand` is an empty struct~~ (Submissions, P3) -- `CreateSubmissionCommand` (`commands.go`) now has `TenantID`, `FormID`, `VersionID`, and `Payload` fields with a `NewCreateSubmissionCommand` constructor. A `Create` method was added to both the `SubmissionsService` interface (`ports/primary.go:17`) and implementation (`submissions_service.go:67-82`). The service validates the command (pointer receiver), constructs the domain aggregate via `NewSubmission`, and persists via `repository.Upsert`. *(Resolves 4/28 #9 for `CreateSubmissionCommand`.)*

4. ~~`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID~~ (Forms, P2) -- Both handlers now use `auth.GetClaimsFromContext(r.Context()).GetSubject()` to obtain the user identity. The `// FIXME` comments are removed. A `PlaceholderAuthenticator` is used temporarily (with a TODO to replace), but the handler-level hardcoding is eliminated. The underlying "no real authentication" problem is tracked separately in #12. *(Resolves 4/28 #1.)*

5. ~~`Find()` has no tenant filtering~~ (Submissions, P2) -- `getSubmissions` handler now extracts tenant from context, constructs a `FindSubmissionsQuery` (validated with `required` tag on `TenantID`), and the service passes a `FindSubmissionsFilter{TenantID}` to the repository. The MongoDB implementation filters by `bson.M{"tenant_id": filter.TenantID}`. The in-memory implementation skips non-matching submissions. *(Resolves 4/28 #4.)*

6. ~~`UpdatePages`/`UpdateSections`/`UpdateFields` naming~~ (Forms, P3) -- Renamed to `ReplacePages` (`version.go:132`), `ReplaceSections` (`page.go:99`), and `ReplaceFields` (`section.go:99`). These methods clear and replace all children, so "Replace" is more semantically accurate. `FormsService.UpdateVersion` call site updated. *(Not previously tracked; naming improvement.)*

7. ~~`FindByIDQuery` generic naming~~ (Submissions, P3) -- Renamed to `FindSubmissionByIDQuery[T]`. A new `FindSubmissionsQuery` struct with `TenantID` validation was added for list operations, along with `FindSubmissionsFilter` for the repository layer. Follows proper CQRS query separation. *(Not previously tracked; architectural improvement.)*

8. ~~Service refactoring~~ (Forms, Tenants, DataSources) -- Redundant `err` check + re-assignment patterns replaced with direct `return repository.Upsert(...)` in `Update`, `PublishVersion`, `RetireVersion`, `UpdateVersion` (forms), `Create`, `Update` (tenants), and `Create`, `Update` (data sources). Reduces boilerplate without changing behavior. *(Not previously tracked; code quality improvement.)*

9. ~~`TenantMiddleware`/`WithTenant` naming~~ (Shared, P3) -- `TenantMiddleware` renamed to `tenants.NewMiddleware`, `WithTenant` renamed to `SetTenantContext`. Consistent with `auth.NewMiddleware` and `auth.SetClaimsContext` naming conventions. More idiomatic Go. *(Not previously tracked.)*

10. ~~MongoDB connection string hardcodes credentials~~ (Shared, P3) -- `ConnectMongoDB` now uses a `createMongoURI` helper that conditionally includes `username:password@` only when both are non-empty. Settings files no longer contain hardcoded credentials. *(Not previously tracked; security improvement.)*

11. ~~`createSubmission` handler not implemented~~ (Submissions, P2) -- The handler now deserializes a `SubmissionRequest` DTO via `httputil.ReadValidateJSONPayload`, constructs a `CreateSubmissionCommand` via `NewCreateSubmissionCommand`, calls `services.Submissions.Create`, and returns 201 with a `SubmissionResponse`. Follows the established goroutine/channel/select pattern. *(Resolves 4/28 #5 for `createSubmission`.)*

12. ~~Request DTOs not implemented~~ (Submissions, P3) -- `dto/request.go` now defines `SubmissionRequest` with `FormID`, `VersionID`, and `Payload` fields (JSON-tagged). *(Resolves 4/28 #7.)*

13. ~~Handler `resultChan` ordering~~ (Forms, Tenants) -- Across `createForm`, `updateForm`, `createVersion`, `updateVersion`, `createTenant`, `updateTenant`, `createDataSource`, `updateDataSource`, the `resultChan` declaration is moved after request body deserialization/validation. Avoids allocating the channel if validation fails early. *(Not previously tracked; performance improvement.)*

14. ~~Tenants handlers pointer in `APIResponse` type param~~ (Tenants) -- `APIResponse[dto.TenantResponse]` changed to `APIResponse[*dto.TenantResponse]`, removing unnecessary struct copy via dereference. Same for `DataSourceResponse`. *(Not previously tracked; code quality improvement.)*

15. ~~`NewSubmission` sets `Status` to empty string~~ (Submissions, P2) -- `Status` is now set to `SubmissionStatusPending`. Constants defined: `SubmissionStatusPending`, `SubmissionStatusAccepted`, `SubmissionStatusRejected`. The aggregate is now created in a valid initial state. *(Resolves 4/28 #3.)*

16. ~~`SubmissionStatus` has no defined constants~~ (Submissions, P3) -- `const` block with `SubmissionStatusPending`, `SubmissionStatusAccepted`, `SubmissionStatusRejected` added to `submission.go`. *(Resolves 4/28 #12.)*

17. ~~Handler stubs `getSubmissionStatus` and `replaySubmission`~~ (Submissions, P2) -- Both handlers are now fully implemented. `getSubmissionStatus` extracts tenant, looks up submission by reference ID, returns `{status: "..."}`. `replaySubmission` extracts tenant, constructs `ReplaySubmissionCommand` with `TenantID` and `ID`, calls service `Replay`, returns 201 on success. *(Resolves 4/28 #5 remainder.)*

18. ~~`ReplaySubmissionCommand` is an empty struct~~ (Submissions, P3) -- Now has `TenantID string` and `ID domain.SubmissionID` fields with a `NewReplaySubmissionCommand` constructor. *(Resolves 4/28 #9 for `ReplaySubmissionCommand`.)*

19. ~~UUIDv4 -> UUIDv7 migration~~ (All) -- All three services now use `NewID()` (which calls `uuid.NewV7()`) instead of `uuid.NewString()` (v4). A `NewID()` helper is added to each service's `domain/domain.go`. A new `uuidv7` custom validator is registered in `pkg/common/validate/validate.go` for validating UUIDv7 strings. IDs are now time-ordered, improving database index performance. *(Not previously tracked; infrastructure improvement.)*

20. ~~Submissions routes restructured~~ -- Routes changed from `/{referenceId}` to `/{submissionId}` for replay, and `getSubmissionByReferenceID`/`getSubmissionStatus` moved under `/by-reference/{referenceId}`. Better REST resource design separating identity-based operations from reference-based lookups. *(Not previously tracked; API design improvement.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

21. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

22. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

23. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- The `go func() -> chan -> select { case <-r.Context().Done() }` pattern is the established approach across all three services for respecting Chi's context-based request timeouts. Chi's timeout middleware handles writing the timeout response; the handler just needs to stop work and return. *(Closed from 4/28 #14.)*

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(New from current review cycle.)*

#### Architectural

2. **`Replay` service method is a stub** (`submissions_service.go`) -- Returns `nil`. The handler and command are now fully wired, but the service does nothing. The handler returns 201 ("Successfully replayed") despite no replay occurring. *(Unresolved from 4/28 #6 -- handler wired but service logic missing.)*

#### Missing Functionality

3. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. *(Unresolved from 4/28 #10.)*

4. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. All errors map to 500. *(New from current review cycle.)*

#### Code Quality

5. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 4/28 #11.)*

---

### Tenants Service

#### Architectural

6. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/28 #17.)*

#### Missing Functionality

7. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/28 #18.)*

---

### Shared Package

#### Bugs

8. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `tenants.NewMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/28 #19.)*

9. **`GetClaimsFromContext` panics on missing claims** (`claims.go:18`) -- Uses an unchecked type assertion `ctx.Value(ClaimsKey).(Claims)`. If auth middleware is not wired for a service (submissions and tenants currently don't use it) and this function is called, it will panic with a nil interface assertion. Should use comma-ok pattern. *(New from current review cycle.)*

---

### Cross-Service

#### Architectural

10. **Zero test files** in all three services and shared packages. *(Unresolved from 4/28 #20.)*

11. **No domain events** for cross-service communication. *(Unresolved from 4/28 #21.)*

12. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placholder"` -- note typo). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Updated from 4/28 #22 -- auth infrastructure now exists but is placeholder-only.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 2 | `Replay` service method is a stub (handler wired) | Submissions |
| **P2** | 4 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 8 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 9 | `GetClaimsFromContext` unchecked type assertion | Shared |
| **P3** | 10 | Zero test files | All |
| **P3** | 11 | No domain events | All |
| **P3** | 3 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 5 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 7 | `Lookup` value object no validation | Tenants |
| **P3** | 12 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 4/28

Sixteen commits since the last review. The net changes represent a major maturation of the submissions service and infrastructure improvements across all services:

- **Submissions persistence layer fully implemented** -- All repository methods (`Find`, `FindByID`, `FindByReferenceID`, `Upsert`) are functional in both in-memory and MongoDB implementations. Bidirectional BSON mapping (`to*`/`from*` functions) handles domain-to-document conversion. This closes the P0 issue tracked since 4/25.

- **Submissions write path fully wired end-to-end** -- `createSubmission` handler deserializes a `SubmissionRequest` DTO via `ReadValidateJSONPayload`, constructs a `CreateSubmissionCommand` via `NewCreateSubmissionCommand`, calls `SubmissionsService.Create`, and returns 201 with a `SubmissionResponse`. The service validates the command, constructs the aggregate via `NewSubmission`, and persists via `repository.Upsert`. The full HTTP -> handler -> command -> service -> domain -> repository -> response flow is operational.

- **All handler stubs eliminated** -- `getSubmissionStatus` looks up by reference ID and returns the status string. `replaySubmission` constructs a `ReplaySubmissionCommand` and calls the service (though the service itself is still a stub). No handlers return empty 200 OK responses.

- **Submission domain model completed** -- `NewSubmission` now correctly assigns all fields including `Status: SubmissionStatusPending` and `Payload`. Three status constants defined (`pending`, `accepted`, `rejected`). The aggregate is created in a valid initial state. `ReplaySubmissionCommand` has proper fields.

- **UUIDv7 migration** -- All domain constructors across all three services now use `NewID()` (UUIDv7) instead of `uuid.NewString()` (UUIDv4). A `NewID()` helper in each service's `domain/domain.go` generates time-ordered UUIDs, improving database index locality. A custom `uuidv7` validator is registered for input validation.

- **Tenant filtering implemented** -- `getSubmissions` handler extracts tenant, constructs a validated `FindSubmissionsQuery`, and the service passes a `FindSubmissionsFilter` to the repository. Both implementations filter correctly. This closes a P2 data isolation issue tracked since 4/25.

- **Authentication infrastructure introduced** -- New `pkg/auth` package provides a `Claims` interface, `Authenticator` function type, and `NewMiddleware` that chains authenticators and returns 401 if none succeed. Forms service wires the middleware with a `PlaceholderAuthenticator` and uses `claims.GetSubject()` for publish/retire operations, removing the hardcoded `"placeholder"` user ID.

- **Routes restructured** -- Submissions routes now separate identity-based operations (`/{submissionId}/replay`) from reference-based lookups (`/by-reference/{referenceId}/` and `/by-reference/{referenceId}/status`). Better REST resource design.

- **Forms service refactored** -- `ReplacePages`/`ReplaceSections`/`ReplaceFields` naming. Direct-return patterns. Handler `resultChan` ordering after validation.

- **Tenants/DataSources services refactored** -- Same direct-return pattern. Handler ordering. Pointer-based `APIResponse` type params.

- **Shared package improvements** -- `tenants.NewMiddleware`/`SetTenantContext` naming. MongoDB `createMongoURI` helper. `uuidv7` custom validator.

### Current State

**22 remaining issues (4/28) -> 12 remaining issues** (resolved 13 from prior reviews including P0 MongoDB stubs, P2 Upsert stubs, P2 placeholder user ID, P2 Find no tenant filtering, P2 createSubmission handler, P2 handler stubs, P2 Status empty, P3 Status no constants, P3 CreateSubmissionCommand empty, P3 ReplaySubmissionCommand empty, P3 Request DTOs; moved 3 to Will Not Fix; introduced 3 new issues: missing validation tags P2, `GetClaimsFromContext` panic P2, `sendErrorResponse` no mapping P2).

**Forms Service** is the most complete service. No remaining issues specific to forms. Rich domain model with version state machine, UUIDv7 IDs, position-keyed sorted collections, complete error-to-HTTP mapping, transactional version creation, authentication integration, and clean code patterns throughout.

**Tenants Service** is stable with complete CRUD, cascade-delete, lookup strategies, and UUIDv7 migration. Remaining gaps are minor: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has undergone the most dramatic improvement this review cycle. From 14 issues down to 5. All handlers are implemented. The create flow is end-to-end functional. The domain model has proper status constants and a valid initial state. The persistence layer is complete. Primary remaining gaps: `CreateSubmissionCommand` missing validation tags (P2), `Replay` service still a stub (P2), `sendErrorResponse` no domain error mapping (P2), `SubmissionAttempt` no constructor (P3), and `Payload` as `any` (P3).

**Hexagonal Architecture** -- Dependency direction remains correct throughout. The `SubmissionRequest` DTO in the REST adapter correctly converts to a domain-agnostic command before crossing the port boundary. The route restructure cleanly separates resource identity (`/{submissionId}`) from alternative lookups (`/by-reference/{referenceId}`). The `NewID()` helper in the domain layer keeps UUID generation as a domain concern.

**CQRS-Lite** -- All commands now have fields and constructors: `CreateSubmissionCommand`, `ReplaySubmissionCommand`. The handler -> command -> service -> repository flow is fully operational for create. `Replay` has the command and handler wired but the service is a stub. Queries (`FindSubmissionsQuery`, `FindSubmissionByIDQuery`) are validated and separated from repository filters.

**DDD** -- The submissions domain model has matured significantly. `NewSubmission` now produces a valid aggregate with `SubmissionStatusPending`, correct `Payload` assignment, UUIDv7 IDs, and struct validation. The three status constants define a clear lifecycle. The route restructure correctly models submissions as resources with their own identity, separate from the reference-based lookup path. `SubmissionAttempt` remains the only entity without a proper constructor. The forms domain (`Version` state machine, `Replace*` methods, `withPosition` mixin) and tenants domain (`DataSourceAttributes` sealed interface, strategy pattern) continue to exemplify strong DDD patterns.

**Idiomatic Go** -- `NewID()` panicking on UUID generation failure follows the accepted Go pattern of panicking on impossible errors (random source failure). The `uuidv7` custom validator is cleanly integrated into the existing validation infrastructure. The `NewReplaySubmissionCommand` constructor follows the established pattern. All handler patterns are now consistent across services.

### Highest-Impact Improvements

1. **Add validation tags to `CreateSubmissionCommand`** (P2 -- `validate.ValidateStruct` is a no-op without tags)
2. **Implement `Replay` service logic** (P2 -- handler returns 201 "success" but nothing happens)
3. **Fix `GetClaimsFromContext` type assertion** (P2 -- use comma-ok pattern to avoid panic)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
6. **Add test coverage** starting with domain constructors and the `Create` flow (P3 -- long-term reliability)
