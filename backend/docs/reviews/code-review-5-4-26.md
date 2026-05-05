# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface / `Upsert` stubs~~ (Submissions, P2) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations are fully functional. The in-memory implementation stores by ID under a write lock. The MongoDB implementation uses `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a session. New `toSubmissionDocument` and `toSubmissionAttemptDocument` mapping functions handle domain-to-BSON conversion with `bson.Marshal` for the `Payload` and `ErrorDetails` fields. *(Resolves 4/28 #8.)*

3. ~~`CreateSubmissionCommand` is an empty struct~~ (Submissions, P3) -- `CreateSubmissionCommand` (`commands.go`) now has `TenantID`, `FormID`, `VersionID`, and `Payload` fields with a `NewCreateSubmissionCommand` constructor. A `Create` method was added to both the `SubmissionsService` interface (`ports/primary.go:17`) and implementation (`submissions_service.go:67-82`). The service validates the command (pointer receiver), constructs the domain aggregate via `NewSubmission`, and persists via `repository.Upsert`. *(Resolves 4/28 #9 for `CreateSubmissionCommand`.)*

4. ~~`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID~~ (Forms, P2) -- Both handlers now use `auth.GetClaimsFromContext(r.Context()).GetSubject()` to obtain the user identity. The `// FIXME` comments are removed. A `PlaceholderAuthenticator` is used temporarily (with a TODO to replace), but the handler-level hardcoding is eliminated. The underlying "no real authentication" problem is tracked separately in #10. *(Resolves 4/28 #1.)*

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

21. ~~`GetClaimsFromContext` panics on missing claims~~ (Shared, P2) -- Now uses comma-ok pattern. On failure, logs error via `slog.ErrorContext` and calls `os.Exit(1)` instead of panicking with a nil interface assertion. Deliberate fail-fast on misconfiguration -- if auth middleware is correctly wired, claims will always be in context. *(Resolves issue from current review cycle.)*

22. ~~`ErrMissingTenantID` maps to 500~~ (Shared, P2) -- The issue is eliminated by design. `TenantFromContext` no longer returns an error -- it logs and calls `os.Exit(1)` on failure (misconfiguration). The tenant middleware guarantees the header is present before handlers run; the middleware itself returns the error response to the client. The `getTenantFromContext` handler helper is removed from all services; handlers call `tenants.TenantFromContext(r.Context())` directly. *(Resolves 4/28 #19.)*

23. ~~`log.Logger` -> `slog.Logger` migration~~ (All) -- All three services now use Go's structured logging (`log/slog`). `main.go` creates `slog.New(slog.NewTextHandler(os.Stdout, nil))`. Service constructors, repositories, persistence bootstrap functions, and strategy implementations all accept `*slog.Logger`. `logger.Fatal` replaced with `logger.Error` + `os.Exit(1)`. *(Not previously tracked; infrastructure improvement.)*

24. ~~`getTenantFromContext` handler helper removed~~ (All) -- Handlers now call `tenants.TenantFromContext(r.Context())` directly. Since the function exits on failure (misconfiguration), error handling boilerplate is eliminated (~5 lines per handler). Consistent with the fail-fast philosophy for infrastructure invariants guaranteed by middleware. *(Not previously tracked; code quality improvement.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

25. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

26. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

27. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- The `go func() -> chan -> select { case <-r.Context().Done() }` pattern is the established approach across all three services for respecting Chi's context-based request timeouts. Chi's timeout middleware handles writing the timeout response; the handler just needs to stop work and return. *(Closed from 4/28 #14.)*

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(New from current review cycle.)*

#### Architectural

2. **`Replay` service method only validates and checks existence** (`submissions_service.go`) -- Validates the command and calls `repository.FindByID` to verify the submission exists, then returns nil without performing any replay logic. The handler returns 201 ("Successfully replayed") despite no replay occurring. *(Partially resolved from 4/28 #6 -- no longer a pure stub but the actual replay logic is unimplemented.)*

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

### Cross-Service

#### Architectural

8. **Zero test files** in all three services and shared packages. *(Unresolved from 4/28 #20.)*

9. **No domain events** for cross-service communication. *(Unresolved from 4/28 #21.)*

10. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placholder"` -- note typo). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Updated from 4/28 #22 -- auth infrastructure now exists but is placeholder-only.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 4 | `sendErrorResponse` no domain error mapping | Submissions |
| **P3** | 2 | `Replay` validates but doesn't replay | Submissions |
| **P3** | 3 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 5 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 7 | `Lookup` value object no validation | Tenants |
| **P3** | 8 | Zero test files | All |
| **P3** | 9 | No domain events | All |
| **P3** | 10 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 4/28

Sixteen commits since the last review representing a major maturation across all services:

- **Submissions persistence layer fully implemented** -- All repository methods (`Find`, `FindByID`, `FindByReferenceID`, `Upsert`) are functional in both in-memory and MongoDB implementations. Bidirectional BSON mapping handles domain-to-document conversion. This closes the P0 issue tracked since 4/25.

- **Submissions write path fully wired end-to-end** -- `createSubmission` handler deserializes a `SubmissionRequest` DTO via `ReadValidateJSONPayload`, constructs a `CreateSubmissionCommand`, calls `SubmissionsService.Create`, and returns 201 with a `SubmissionResponse`. The full HTTP -> handler -> command -> service -> domain -> repository -> response flow is operational.

- **All handler stubs eliminated** -- Every route has a real implementation. `getSubmissionStatus` returns status by reference ID. `replaySubmission` validates and checks existence (though actual replay logic is unimplemented).

- **Submission domain model completed** -- `NewSubmission` correctly assigns all fields including `Status: SubmissionStatusPending` and `Payload`. Three status constants defined. `ReplaySubmissionCommand` has proper fields. The aggregate is created in a valid initial state.

- **UUIDv7 migration** -- All domain constructors across all three services now use `NewID()` (UUIDv7) for time-ordered IDs. Custom `uuidv7` validator registered for input validation.

- **Structured logging migration** -- All three services converted from `log.Logger` to `slog.Logger`. Consistent structured logging throughout services, repositories, and strategies.

- **Fail-fast on misconfiguration** -- `TenantFromContext` and `GetClaimsFromContext` now use the comma-ok pattern and call `os.Exit(1)` with an error log on failure, rather than returning errors or panicking. This eliminates handler-level error handling boilerplate for invariants guaranteed by middleware. The `getTenantFromContext` helper is removed from all handlers.

- **Authentication infrastructure introduced** -- `pkg/auth` package with `Claims` interface, `Authenticator` type, `NewMiddleware`. Forms service uses `claims.GetSubject()` for publish/retire, removing the hardcoded user ID.

- **Routes restructured** -- Submissions routes separate identity-based operations (`/{submissionId}/replay`) from reference-based lookups (`/by-reference/{referenceId}`).

- **Forms/Tenants services refactored** -- Direct-return patterns, `Replace*` naming, handler `resultChan` ordering, pointer `APIResponse` type params.

### Current State

**22 remaining issues (4/28) -> 10 remaining issues** (resolved 15 from prior reviews; moved 3 to Will Not Fix; introduced 3 new issues during the review cycle of which 2 are now resolved). Only 2 P2 issues remain; the rest are P3.

**Forms Service** is the most complete service. No remaining issues. Rich domain model with version state machine, UUIDv7 IDs, structured logging, authentication integration, complete error-to-HTTP mapping, and clean handler patterns with fail-fast tenant/claims extraction.

**Tenants Service** is stable with complete CRUD, cascade-delete, lookup strategies, structured logging, and UUIDv7 migration. Remaining gaps are minor: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has undergone the most dramatic improvement this review cycle -- from 14 issues down to 5. All handlers implemented. Create flow end-to-end functional. Domain model has proper status constants and valid initial state. Persistence complete. Primary remaining gaps: `CreateSubmissionCommand` missing validation tags (P2), `sendErrorResponse` no domain error mapping (P2), `Replay` logic unimplemented (P3), `SubmissionAttempt` no constructor (P3), `Payload` as `any` (P3).

**Hexagonal Architecture** -- Dependency direction correct throughout. The fail-fast approach for `TenantFromContext`/`GetClaimsFromContext` enforces that middleware correctly establishes invariants before handlers run -- this is an infrastructure boundary concern, not a domain concern. The `SubmissionRequest` DTO converts to a command before crossing the port boundary. Routes cleanly separate resource identity from alternative lookups.

**CQRS-Lite** -- All commands have fields and constructors. The handler -> command -> service -> repository flow is fully operational for create and partially for replay. Queries are validated and separated from repository filters. No violations.

**DDD** -- The submissions domain model is now functional: valid initial state, status constants, proper constructors. `SubmissionAttempt` remains the only entity without a constructor. Forms domain (state machine, `Replace*`, `withPosition`) and tenants domain (`DataSourceAttributes` sealed interface, strategy pattern) continue to exemplify strong DDD patterns. No domain events exist for cross-aggregate communication.

**Idiomatic Go** -- The `slog` migration follows Go 1.21+ best practices. The fail-fast pattern (`os.Exit(1)` on misconfiguration) is appropriate for invariants that should never be violated at runtime. `NewID()` panicking on UUID failure follows the accepted pattern for impossible errors. All handler patterns are consistent across services. The `uuidv7` custom validator integrates cleanly.

### Highest-Impact Improvements

1. **Add validation tags to `CreateSubmissionCommand`** (P2 -- `validate.ValidateStruct` is a no-op without tags)
2. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
3. **Implement actual `Replay` logic** (P3 -- handler returns success but nothing happens)
4. **Add test coverage** starting with domain constructors and the `Create` flow (P3 -- long-term reliability)
5. **Add `NewSubmissionAttempt` constructor** (P3 -- needed for replay implementation)
