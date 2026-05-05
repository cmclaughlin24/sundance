# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface / `Upsert` stubs~~ (Submissions, P2) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations are fully functional. The in-memory implementation stores by ID under a write lock. The MongoDB implementation uses `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a session. New `toSubmissionDocument` and `toSubmissionAttemptDocument` mapping functions handle domain-to-BSON conversion with `bson.Marshal` for the `Payload` and `ErrorDetails` fields. *(Resolves 4/28 #8.)*

3. ~~`CreateSubmissionCommand` is an empty struct~~ (Submissions, P3) -- `CreateSubmissionCommand` (`commands.go`) now has `TenantID`, `FormID`, `VersionID`, and `Payload` fields with a `NewCreateSubmissionCommand` constructor. A `Create` method was added to both the `SubmissionsService` interface (`ports/primary.go:17`) and implementation (`submissions_service.go:67-82`). The service validates the command (pointer receiver), constructs the domain aggregate via `NewSubmission`, and persists via `repository.Upsert`. *(Resolves 4/28 #9 for `CreateSubmissionCommand`; `ReplaySubmissionCommand` remains empty.)*

4. ~~`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID~~ (Forms, P2) -- Both handlers now use `auth.GetClaimsFromContext(r.Context()).GetSubject()` to obtain the user identity. The `// FIXME` comments are removed. A `PlaceholderAuthenticator` is used temporarily (with a TODO to replace), but the handler-level hardcoding is eliminated. The underlying "no real authentication" problem is tracked separately in #16. *(Resolves 4/28 #1.)*

5. ~~`Find()` has no tenant filtering~~ (Submissions, P2) -- `getSubmissions` handler now extracts tenant from context, constructs a `FindSubmissionsQuery` (validated with `required` tag on `TenantID`), and the service passes a `FindSubmissionsFilter{TenantID}` to the repository. The MongoDB implementation filters by `bson.M{"tenant_id": filter.TenantID}`. The in-memory implementation skips non-matching submissions. *(Resolves 4/28 #4.)*

6. ~~`UpdatePages`/`UpdateSections`/`UpdateFields` naming~~ (Forms, P3) -- Renamed to `ReplacePages` (`version.go:132`), `ReplaceSections` (`page.go:99`), and `ReplaceFields` (`section.go:99`). These methods clear and replace all children, so "Replace" is more semantically accurate. `FormsService.UpdateVersion` call site updated. *(Not previously tracked; naming improvement.)*

7. ~~`FindByIDQuery` generic naming~~ (Submissions, P3) -- Renamed to `FindSubmissionByIDQuery[T]`. A new `FindSubmissionsQuery` struct with `TenantID` validation was added for list operations, along with `FindSubmissionsFilter` for the repository layer. Follows proper CQRS query separation. *(Not previously tracked; architectural improvement.)*

8. ~~Service refactoring~~ (Forms, Tenants, DataSources) -- Redundant `err` check + re-assignment patterns replaced with direct `return repository.Upsert(...)` in `Update`, `PublishVersion`, `RetireVersion`, `UpdateVersion` (forms), `Create`, `Update` (tenants), and `Create`, `Update` (data sources). Reduces boilerplate without changing behavior. *(Not previously tracked; code quality improvement.)*

9. ~~`TenantMiddleware`/`WithTenant` naming~~ (Shared, P3) -- `TenantMiddleware` renamed to `tenants.NewMiddleware`, `WithTenant` renamed to `SetTenantContext`. Consistent with `auth.NewMiddleware` and `auth.SetClaimsContext` naming conventions. More idiomatic Go. *(Not previously tracked.)*

10. ~~MongoDB connection string hardcodes credentials~~ (Shared, P3) -- `ConnectMongoDB` now uses a `createMongoURI` helper that conditionally includes `username:password@` only when both are non-empty. Settings files no longer contain hardcoded credentials. *(Not previously tracked; security improvement.)*

11. ~~`createSubmission` handler not implemented~~ (Submissions, P2) -- The handler now deserializes a `SubmissionRequest` DTO via `httputil.ReadJSONPayload`, constructs a `CreateSubmissionCommand` via `NewCreateSubmissionCommand`, calls `services.Submissions.Create`, and returns 201 with a `SubmissionResponse`. Follows the established goroutine/channel/select pattern. *(Resolves 4/28 #5 for `createSubmission`.)*

12. ~~Request DTOs not implemented~~ (Submissions, P3) -- `dto/request.go` now defines `SubmissionRequest` with `FormID`, `VersionID`, and `Payload` fields (JSON-tagged). *(Resolves 4/28 #7.)*

13. ~~Handler `resultChan` ordering~~ (Forms, Tenants) -- Across `createForm`, `updateForm`, `createVersion`, `updateVersion`, `createTenant`, `updateTenant`, `createDataSource`, `updateDataSource`, the `resultChan` declaration is moved after request body deserialization/validation. Avoids allocating the channel if validation fails early. *(Not previously tracked; performance improvement.)*

14. ~~Tenants handlers pointer in `APIResponse` type param~~ (Tenants) -- `APIResponse[dto.TenantResponse]` changed to `APIResponse[*dto.TenantResponse]`, removing unnecessary struct copy via dereference. Same for `DataSourceResponse`. *(Not previously tracked; code quality improvement.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

15. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

16. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

17. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- The `go func() -> chan -> select { case <-r.Context().Done() }` pattern is the established approach across all three services for respecting Chi's context-based request timeouts. Chi's timeout middleware handles writing the timeout response; the handler just needs to stop work and return. *(Closed from 4/28 #14.)*

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`NewSubmission` sets `Status` to empty string** (`submission.go:41`) -- The constructor hardcodes `Status: ""` with a TODO comment about implementing a state machine. `Payload` is correctly assigned. Submissions are persisted with an undefined status, which means consumers cannot determine submission state. Tied to the absence of `SubmissionStatus` constants (#9). *(Downgraded from 4/28 #3 -- `Payload` is now correctly assigned; only `Status` remains a gap.)*

2. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(New.)*

#### Architectural

3. **Two handler stubs return 200 OK with empty body** -- `getSubmissionStatus` and `replaySubmission`. *(Reduced from 4/28 #5.)*

#### Missing Functionality

4. **`Replay` service method is a stub** (`submissions_service.go`) -- Returns `nil`. *(Unresolved from 4/28 #6.)*

5. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/28 #9.)*

6. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. *(Unresolved from 4/28 #10.)*

7. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. All errors map to 500. *(New from current review cycle.)*

#### Code Quality

8. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 4/28 #11.)*

9. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. Tied to #1. *(Unresolved from 4/28 #12.)*

---

### Tenants Service

#### Architectural

10. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/28 #17.)*

#### Missing Functionality

11. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/28 #18.)*

---

### Shared Package

#### Bugs

12. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `tenants.NewMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/28 #19.)*

13. **`GetClaimsFromContext` panics on missing claims** (`claims.go:18`) -- Uses an unchecked type assertion `ctx.Value(ClaimsKey).(Claims)`. If auth middleware is not wired for a service (submissions and tenants currently don't use it) and this function is called, it will panic with a nil interface assertion. Should use comma-ok pattern. *(New.)*

---

### Cross-Service

#### Architectural

14. **Zero test files** in all three services and shared packages. *(Unresolved from 4/28 #20.)*

15. **No domain events** for cross-service communication. *(Unresolved from 4/28 #21.)*

16. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placholder"` -- note typo). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Updated from 4/28 #22 -- auth infrastructure now exists but is placeholder-only.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `NewSubmission` sets `Status` to empty string | Submissions |
| **P2** | 2 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 7 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 12 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 13 | `GetClaimsFromContext` unchecked type assertion | Shared |
| **P3** | 14 | Zero test files | All |
| **P3** | 15 | No domain events | All |
| **P3** | 6 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 8 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 9 | `SubmissionStatus` no constants | Submissions |
| **P3** | 11 | `Lookup` value object no validation | Tenants |
| **P3** | 16 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 4/28

Twelve commits since the last review (including reverts and reapplies for an authentication refactor) plus unstaged work. The net changes are substantial:

- **Submissions persistence layer fully implemented** -- All repository methods (`Find`, `FindByID`, `FindByReferenceID`, `Upsert`) are functional in both in-memory and MongoDB implementations. Bidirectional BSON mapping (`to*`/`from*` functions) handles domain-to-document conversion. This closes the P0 issue tracked since 4/25.

- **Submissions write path fully wired end-to-end** -- `createSubmission` handler deserializes a `SubmissionRequest` DTO, constructs a `CreateSubmissionCommand` via `NewCreateSubmissionCommand`, calls `SubmissionsService.Create`, and returns 201 with a `SubmissionResponse`. The service validates the command, constructs the aggregate via `NewSubmission`, and persists via `repository.Upsert`. The full HTTP -> handler -> command -> service -> domain -> repository -> response flow is operational.

- **Tenant filtering implemented** -- `getSubmissions` handler extracts tenant, constructs a validated `FindSubmissionsQuery`, and the service passes a `FindSubmissionsFilter` to the repository. Both implementations filter correctly. This closes a P2 data isolation issue tracked since 4/25.

- **Authentication infrastructure introduced** -- New `pkg/auth` package provides a `Claims` interface, `Authenticator` function type, and `NewMiddleware` that chains authenticators and returns 401 if none succeed. `SetClaimsContext`/`GetClaimsFromContext` handle context propagation. Forms service wires the middleware with a `PlaceholderAuthenticator` and uses `claims.GetSubject()` for publish/retire operations, removing the hardcoded `"placeholder"` user ID.

- **Forms service refactored** -- `Update`, `PublishVersion`, `RetireVersion`, and `UpdateVersion` now return `repository.Upsert(...)` directly, eliminating redundant error check + re-assignment patterns. `ReplacePages`/`ReplaceSections`/`ReplaceFields` naming replaces the misleading `Update*` names. Handler `resultChan` declarations moved after validation to avoid unnecessary allocation on early-return paths.

- **Tenants/DataSources services refactored** -- Same direct-return pattern applied to `Create` and `Update` in both services. Handler `resultChan` ordering fixed. `APIResponse` type params changed to pointer types, removing unnecessary struct copies.

- **Shared package improvements** -- `TenantMiddleware` renamed to `tenants.NewMiddleware`, `WithTenant` renamed to `SetTenantContext`. MongoDB `createMongoURI` helper conditionally includes credentials. Settings files no longer contain hardcoded usernames/passwords.

- **Submissions settings switched to MongoDB** -- `settings.json` now uses `"driver": "mongodb"` instead of `"in-memory"`.

### Current State

**22 remaining issues (4/28) -> 16 remaining issues** (resolved 7 from prior reviews: P0 MongoDB stubs, P2 Upsert stubs, P2 placeholder user ID, P2 Find no tenant filtering, P2 createSubmission handler, P3 CreateSubmissionCommand empty, P3 Request DTOs; moved 3 to Will Not Fix; introduced 3 new issues: missing validation tags P2, `GetClaimsFromContext` panic P2, `sendErrorResponse` no mapping P2; downgraded 1: `NewSubmission` Status from P1 to P2).

**Forms Service** is now the most complete service. The authentication integration (even if placeholder) demonstrates the full request lifecycle: middleware extracts claims, handlers use claims for domain operations. Rich domain model with version state machine, position-keyed sorted collections, complete error-to-HTTP mapping, transactional version creation, and direct-return patterns that reduce boilerplate. No remaining issues specific to forms.

**Tenants Service** is stable with complete CRUD, cascade-delete, and lookup strategies. Clean refactoring to direct-return patterns. Pointer-based `APIResponse` type params avoid struct copies. Remaining gaps are minor: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has made the most progress this review cycle. The write path is now functional end-to-end: `createSubmission` handler -> `SubmissionRequest` DTO -> `CreateSubmissionCommand` -> `SubmissionsService.Create` -> `NewSubmission` -> `repository.Upsert` -> 201 response. Persistence is complete with bidirectional BSON mapping. Tenant filtering works. The primary gaps are now: `Status` hardcoded to empty string (P2), missing validation tags (P2), no domain error mapping in `sendErrorResponse` (P2), and two remaining handler stubs for replay/status.

**Hexagonal Architecture** -- Dependency direction remains correct throughout. The new `pkg/auth` package acts as a shared infrastructure port: `Claims` is the interface, `Authenticator` is a function type (port), and `PlaceholderAuthenticator` is an adapter. Services depend on the interface, not the implementation. The `FindSubmissionsFilter` as a repository-layer DTO (separate from the application-layer `FindSubmissionsQuery`) is a clean port boundary. The `SubmissionRequest` DTO in the REST adapter correctly converts to a domain-agnostic command before crossing the port boundary.

**CQRS-Lite** -- `FindSubmissionsQuery` (application layer, validated) vs `FindSubmissionsFilter` (repository layer, plain struct) demonstrates proper CQRS query/filter separation. `CreateSubmissionCommand` follows the command pattern with a constructor. The handler -> query/command -> service -> repository flow is now fully operational for the create path. `ReplaySubmissionCommand` remains empty.

**DDD** -- The submissions domain model is still the weakest, but improving. `NewSubmission` correctly assigns all fields except `Status` (empty string with TODO). The absence of `SubmissionStatus` constants means the aggregate lifecycle is undefined -- there's no state machine to enforce valid transitions. `SubmissionAttempt` remains a bare struct with no constructor. The forms domain (`Version` state machine, `Replace*` methods) and tenants domain (`DataSourceAttributes` sealed interface, strategy pattern) continue to exemplify strong DDD patterns.

**Idiomatic Go** -- The direct-return refactoring (`return s.repository.Upsert(ctx, entity)`) is idiomatic. `NewCreateSubmissionCommand` constructor follows the standard pattern. Moving `resultChan` after validation is a sensible optimization. Pointer-based generic type params (`APIResponse[*dto.TenantResponse]`) avoid unnecessary copies. `GetClaimsFromContext` should use the comma-ok pattern for type assertions.

### Highest-Impact Improvements

1. **Add `SubmissionStatus` constants and set initial status in `NewSubmission`** (P2 -- aggregate has no defined lifecycle)
2. **Add validation tags to `CreateSubmissionCommand`** (P2 -- `validate.ValidateStruct` is a no-op without tags)
3. **Fix `GetClaimsFromContext` type assertion** (P2 -- use comma-ok pattern to avoid panic)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
6. **Add test coverage** starting with domain constructors and the `Create` flow (P3 -- long-term reliability)
