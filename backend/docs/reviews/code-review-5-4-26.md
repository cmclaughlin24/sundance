# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface / `Upsert` stubs~~ (Submissions, P2) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations are fully functional. The in-memory implementation stores by ID under a write lock. The MongoDB implementation uses `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a session. New `toSubmissionDocument` and `toSubmissionAttemptDocument` mapping functions handle domain-to-BSON conversion with `bson.Marshal` for the `Payload` and `ErrorDetails` fields. *(Resolves 4/28 #8.)*

3. ~~`CreateSubmissionCommand` is an empty struct~~ (Submissions, P3) -- `CreateSubmissionCommand` (`commands.go`) now has `TenantID`, `FormID`, `VersionID`, and `Payload` fields. A `Create` method was added to both the `SubmissionsService` interface (`ports/primary.go:17`) and implementation (`submissions_service.go:67-82`). The service validates the command, constructs the domain aggregate via `NewSubmission`, and persists via `repository.Upsert`. *(Resolves 4/28 #9 for `CreateSubmissionCommand`; `ReplaySubmissionCommand` remains empty.)*

4. ~~`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID~~ (Forms, P2) -- Both handlers now use `auth.GetClaimsFromContext(r.Context()).GetSubject()` to obtain the user identity. The `// FIXME` comments are removed. A `PlaceholderAuthenticator` is used temporarily (with a TODO to replace), but the handler-level hardcoding is eliminated. The underlying "no real authentication" problem is tracked separately in #18. *(Resolves 4/28 #1.)*

5. ~~`Find()` has no tenant filtering~~ (Submissions, P2) -- `getSubmissions` handler now extracts tenant from context, constructs a `FindSubmissionsQuery` (validated with `required` tag on `TenantID`), and the service passes a `FindSubmissionsFilter{TenantID}` to the repository. The MongoDB implementation filters by `bson.M{"tenant_id": filter.TenantID}`. The in-memory implementation skips non-matching submissions. *(Resolves 4/28 #4.)*

6. ~~`UpdatePages`/`UpdateSections`/`UpdateFields` naming~~ (Forms, P3) -- Renamed to `ReplacePages` (`version.go:132`), `ReplaceSections` (`page.go:99`), and `ReplaceFields` (`section.go:99`). These methods clear and replace all children, so "Replace" is more semantically accurate. `FormsService.UpdateVersion` call site updated. *(Not previously tracked; naming improvement.)*

7. ~~`FindByIDQuery` generic naming~~ (Submissions, P3) -- Renamed to `FindSubmissionByIDQuery[T]`. A new `FindSubmissionsQuery` struct with `TenantID` validation was added for list operations, along with `FindSubmissionsFilter` for the repository layer. Follows proper CQRS query separation. *(Not previously tracked; architectural improvement.)*

8. ~~Service refactoring~~ (Forms, Tenants, DataSources) -- Redundant `err` check + re-assignment patterns replaced with direct `return repository.Upsert(...)` in `Update`, `PublishVersion`, `RetireVersion`, `UpdateVersion` (forms), `Create`, `Update` (tenants), and `Create`, `Update` (data sources). Reduces boilerplate without changing behavior. *(Not previously tracked; code quality improvement.)*

9. ~~`TenantMiddleware`/`WithTenant` naming~~ (Shared, P3) -- `TenantMiddleware` renamed to `tenants.NewMiddleware`, `WithTenant` renamed to `SetTenantContext`. Consistent with `auth.NewMiddleware` and `auth.SetClaimsContext` naming conventions. More idiomatic Go. *(Not previously tracked.)*

10. ~~MongoDB connection string hardcodes credentials~~ (Shared, P3) -- `ConnectMongoDB` now uses a `createMongoURI` helper that conditionally includes `username:password@` only when both are non-empty. Settings files no longer contain hardcoded credentials. *(Not previously tracked; security improvement.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

11. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

12. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

13. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- The `go func() -> chan -> select { case <-r.Context().Done() }` pattern is the established approach across all three services for respecting Chi's context-based request timeouts. Chi's timeout middleware handles writing the timeout response; the handler just needs to stop work and return. *(Closed from 4/28 #14.)*

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`NewSubmission` never sets `Status` or `Payload`** (`submission.go:35-44`) -- The constructor accepts `payload any` but never assigns it to the struct. `Status` is left as zero-value empty string with a TODO comment. Both fields persist as empty/nil despite being provided by the caller. The `Create` service method calls `NewSubmission` and then `Upsert`, meaning submissions are now actively persisted with missing `Status` and `Payload`. *(Escalated from 4/28 #3 -- now exploitable via the write path.)*

2. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(New.)*

#### Architectural

3. **`createSubmission` handler partially implemented** (`handlers.go:87-93`) -- Extracts tenant from context and sends error response on failure, but does not deserialize a request body, construct a `CreateSubmissionCommand`, or call the service `Create` method. Returns 200 OK with empty body on success path. *(Improved from 4/28 #5 but still non-functional.)*

4. **Two handler stubs return 200 OK with empty body** -- `getSubmissionStatus` and `replaySubmission`. *(Reduced from 4/28 #5.)*

#### Missing Functionality

5. **`Replay` service method is a stub** (`submissions_service.go`) -- Returns `nil`. *(Unresolved from 4/28 #6.)*

6. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/28 #7.)*

7. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/28 #9.)*

8. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. *(Unresolved from 4/28 #10.)*

9. **`sendErrorResponse` has no domain error mapping** (`handlers.go:108-112`) -- Switch statement contains only a `default` case. All errors map to 500. *(New from current review cycle.)*

#### Code Quality

10. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 4/28 #11.)*

11. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/28 #12.)*

---

### Tenants Service

#### Architectural

12. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/28 #17.)*

#### Missing Functionality

13. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/28 #18.)*

---

### Shared Package

#### Bugs

14. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `tenants.NewMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/28 #19.)*

15. **`GetClaimsFromContext` panics on missing claims** (`claims.go:18`) -- Uses an unchecked type assertion `ctx.Value(ClaimsKey).(Claims)`. If auth middleware is not wired for a service (submissions and tenants currently don't use it) and this function is called, it will panic with a nil interface assertion. Should use comma-ok pattern. *(New.)*

---

### Cross-Service

#### Architectural

16. **Zero test files** in all three services and shared packages. *(Unresolved from 4/28 #20.)*

17. **No domain events** for cross-service communication. *(Unresolved from 4/28 #21.)*

18. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placholder"` -- note typo). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Updated from 4/28 #22 -- auth infrastructure now exists but is placeholder-only.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 1 | `NewSubmission` never sets `Status` or `Payload` (now exploitable) | Submissions |
| **P2** | 2 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 3 | `createSubmission` handler incomplete | Submissions |
| **P2** | 9 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 14 | `ErrMissingTenantID` maps to 500 | Shared |
| **P2** | 15 | `GetClaimsFromContext` unchecked type assertion | Shared |
| **P3** | 16 | Zero test files | All |
| **P3** | 17 | No domain events | All |
| **P3** | 8 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 10 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 11 | `SubmissionStatus` no constants | Submissions |
| **P3** | 13 | `Lookup` value object no validation | Tenants |
| **P3** | 18 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 4/28

Twelve commits since the last review (including reverts and reapplies for an authentication refactor). The net changes are substantial:

- **Submissions persistence layer fully implemented** -- All repository methods (`Find`, `FindByID`, `FindByReferenceID`, `Upsert`) are functional in both in-memory and MongoDB implementations. Bidirectional BSON mapping (`to*`/`from*` functions) handles domain-to-document conversion. This closes the P0 issue tracked since 4/25.

- **Submissions write path partially wired** -- `CreateSubmissionCommand` carries `TenantID`, `FormID`, `VersionID`, and `Payload`. `SubmissionsService.Create` validates the command, constructs the domain aggregate, and persists via `repository.Upsert`. The `createSubmission` handler extracts tenant but does not yet deserialize a request body or invoke the service.

- **Tenant filtering implemented** -- `getSubmissions` handler extracts tenant, constructs a validated `FindSubmissionsQuery`, and the service passes a `FindSubmissionsFilter` to the repository. Both implementations filter correctly. This closes a P2 data isolation issue tracked since 4/25.

- **Authentication infrastructure introduced** -- New `pkg/auth` package provides a `Claims` interface, `Authenticator` function type, and `NewMiddleware` that chains authenticators and returns 401 if none succeed. `SetClaimsContext`/`GetClaimsFromContext` handle context propagation. Forms service wires the middleware with a `PlaceholderAuthenticator` and uses `claims.GetSubject()` for publish/retire operations, removing the hardcoded `"placeholder"` user ID.

- **Forms service refactored** -- `Update`, `PublishVersion`, `RetireVersion`, and `UpdateVersion` now return `repository.Upsert(...)` directly, eliminating redundant error check + re-assignment patterns. `ReplacePages`/`ReplaceSections`/`ReplaceFields` naming replaces the misleading `Update*` names.

- **Tenants/DataSources services refactored** -- Same direct-return pattern applied to `Create` and `Update` in both services. Blank line cleanup throughout.

- **Shared package improvements** -- `TenantMiddleware` renamed to `tenants.NewMiddleware`, `WithTenant` renamed to `SetTenantContext`. MongoDB `createMongoURI` helper conditionally includes credentials. Settings files no longer contain hardcoded usernames/passwords.

- **Submissions settings switched to MongoDB** -- `settings.json` now uses `"driver": "mongodb"` instead of `"in-memory"`.

### Current State

**22 remaining issues (4/28) -> 18 remaining issues** (resolved 5 from prior reviews: P0 MongoDB stubs, P2 Upsert stubs, P2 placeholder user ID, P2 Find no tenant filtering, P3 CreateSubmissionCommand empty; moved 3 to Will Not Fix; introduced 3 new issues: missing validation tags P2, `GetClaimsFromContext` panic P2, `sendErrorResponse` no mapping P2).

**Forms Service** is now the most complete service. The authentication integration (even if placeholder) demonstrates the full request lifecycle: middleware extracts claims, handlers use claims for domain operations. Rich domain model with version state machine, position-keyed sorted collections, complete error-to-HTTP mapping, transactional version creation, and direct-return patterns that reduce boilerplate. No remaining issues specific to forms.

**Tenants Service** is stable with complete CRUD, cascade-delete, and lookup strategies. Clean refactoring to direct-return patterns. Remaining gaps are minor: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has made significant progress: persistence is complete, tenant filtering works, and the write path exists at the service layer. However, the domain model remains fundamentally broken: `NewSubmission` doesn't assign `Status` or `Payload`, which means the `Create` flow actively persists incomplete aggregates. The handler layer still needs to deserialize request bodies. Two handlers remain empty stubs.

**Hexagonal Architecture** -- Dependency direction remains correct throughout. The new `pkg/auth` package acts as a shared infrastructure port: `Claims` is the interface, `Authenticator` is a function type (port), and `PlaceholderAuthenticator` is an adapter. Services depend on the interface, not the implementation. The `FindSubmissionsFilter` as a repository-layer DTO (separate from the application-layer `FindSubmissionsQuery`) is a clean port boundary.

**CQRS-Lite** -- `FindSubmissionsQuery` (application layer, validated) vs `FindSubmissionsFilter` (repository layer, plain struct) demonstrates proper CQRS query/filter separation. `CreateSubmissionCommand` follows the command pattern. The handler -> query/command -> service -> repository flow is well-established. `ReplaySubmissionCommand` remains empty.

**DDD** -- The submissions domain remains the weakest DDD implementation. `NewSubmission` produces an aggregate missing `Status` and `Payload`, violating the invariant that a newly-created aggregate should be in a valid state. This is now a live bug. The forms `Replace*` rename is a positive DDD signal -- method names accurately describe aggregate behavior. The auth `Claims` interface as a cross-cutting concern is appropriately modeled outside any specific bounded context.

**Idiomatic Go** -- The direct-return refactoring (`return s.repository.Upsert(ctx, entity)`) is idiomatic and reduces unnecessary variable declarations. `tenants.NewMiddleware` and `auth.NewMiddleware` follow the standard `New*` constructor convention. The `Authenticator` type alias (`type Authenticator = func(*http.Request) (Claims, error)`) is clean. However, `GetClaimsFromContext` should use the comma-ok pattern for the type assertion to avoid panics -- this is a deviation from idiomatic Go error handling.

### Highest-Impact Improvements

1. **Fix `NewSubmission` to assign `Status` and `Payload`** (P1 -- constructor produces incomplete domain objects, now actively persisted)
2. **Add validation tags to `CreateSubmissionCommand`** (P2 -- `validate.ValidateStruct` is a no-op without tags)
3. **Complete `createSubmission` handler** (P2 -- deserialize request body, construct command, call service)
4. **Fix `GetClaimsFromContext` type assertion** (P2 -- use comma-ok pattern to avoid panic)
5. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
6. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
7. **Add test coverage** starting with domain constructors and the `Create` flow (P3 -- long-term reliability)
