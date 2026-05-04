# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface~~ (Submissions, P2) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations are fully functional. The in-memory implementation stores by ID under a write lock. The MongoDB implementation uses `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a session. New `toSubmissionDocument` and `toSubmissionAttemptDocument` mapping functions handle domain-to-BSON conversion with `bson.Marshal` for the `Payload` and `ErrorDetails` fields. *(Resolves 4/28 #8.)*

3. ~~`CreateSubmissionCommand` is an empty struct~~ (Submissions, P3) -- `CreateSubmissionCommand` (`commands.go`) now has `TenantID`, `FormID`, `VersionID`, and `Payload` fields. A `Create` method was added to both the `SubmissionsService` interface (`ports/primary.go:17`) and implementation (`submissions_service.go:65-82`). The service validates the command, calls `domain.NewSubmission`, and persists via `repository.Upsert`. *(Resolves 4/28 #9 for `CreateSubmissionCommand`; `ReplaySubmissionCommand` remains empty.)*

4. ~~`UpdatePages`/`UpdateSections`/`UpdateFields` naming~~ (Forms, P3) -- Renamed to `ReplacePages` (`version.go:132`), `ReplaceSections` (`page.go:99`), and `ReplaceFields` (`section.go:99`). These methods clear and replace all children, so "Replace" is more semantically accurate than "Update" which implies a merge/patch. `FormsService.UpdateVersion` call site updated. *(Not previously tracked; naming improvement.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

11. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

12. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

13. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- The `go func() -> chan -> select { case <-r.Context().Done() }` pattern is the established approach across all three services for respecting Chi's context-based request timeouts. Chi's timeout middleware handles writing the timeout response; the handler just needs to stop work and return. *(Closed from 4/28 #14.)*

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:373`, `408`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/28 #1.)*

---

### Submissions Service

#### Bugs

2. **`NewSubmission` never sets `Status` or `Payload`** (`submission.go:35-44`) -- The constructor accepts `payload any` but never assigns it to the struct. `Status` is left as zero-value empty string with a TODO comment. Both fields persist as empty/nil despite being provided by the caller. The new `Create` service method calls `NewSubmission` and then `Upsert`, meaning submissions are now actively persisted with missing `Status` and `Payload`. *(Escalated from 4/28 #3 -- now exploitable via the write path.)*

3. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(New.)*

#### Architectural

4. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. The `getSubmissions` handler does not call `getTenantFromContext`. *(Unresolved from 4/28 #4.)*

5. **`createSubmission` handler partially implemented** (`handlers.go:87-93`) -- Extracts tenant from context and sends error response on failure, but does not deserialize a request body, construct a `CreateSubmissionCommand`, or call the service `Create` method. Returns 200 OK with empty body on success path. *(Improved from 4/28 #5 but still non-functional.)*

6. **Two handler stubs return 200 OK with empty body** -- `getSubmissionStatus` and `replaySubmission`. *(Reduced from 4/28 #5.)*

#### Missing Functionality

7. **`Replay` service method is a stub** (`submissions_service.go`) -- Returns `nil`. *(Unresolved from 4/28 #6.)*

8. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/28 #7.)*

9. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/28 #9.)*

10. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. *(Unresolved from 4/28 #10.)*

11. **`sendErrorResponse` has no domain error mapping** (`handlers.go:108-112`) -- Switch statement contains only a `default` case. All errors map to 500. *(New from current review cycle.)*

#### Code Quality

12. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 4/28 #11.)*

13. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/28 #12.)*

---

### Tenants Service

#### Architectural

14. **`Find()` has no pagination or filtering** (`tenants_service.go:25-27`). *(Unresolved from 4/28 #17.)*

#### Missing Functionality

15. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 4/28 #18.)*

---

### Shared Package

#### Bugs

16. **`ErrMissingTenantID` maps to 500** (`middleware.go:15`) -- `TenantMiddleware` calls `httputil.SendErrorResponse(w, ErrMissingTenantID)` when the `X-Tenant-ID` header is absent. `ErrMissingTenantID` doesn't match any case in `SendErrorResponse`, falling through to 500. Should be 400. *(Unresolved from 4/28 #19.)*

---

### Cross-Service

#### Architectural

17. **Zero test files** in all three services and shared packages. *(Unresolved from 4/28 #20.)*

18. **No domain events** for cross-service communication. *(Unresolved from 4/28 #21.)*

19. **No real authentication** -- `X-Tenant-ID` is blindly trusted. *(Unresolved from 4/28 #22.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 2 | `NewSubmission` never sets `Status` or `Payload` (now exploitable) | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 3 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 4 | `Find()` has no tenant filtering | Submissions |
| **P2** | 5 | `createSubmission` handler incomplete | Submissions |
| **P2** | 11 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 16 | `ErrMissingTenantID` maps to 500 | Shared |
| **P3** | 17 | Zero test files | All |
| **P3** | 18 | No domain events | All |
| **P3** | 10 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 12 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 13 | `SubmissionStatus` no constants | Submissions |
| **P3** | 15 | `Lookup` value object no validation | Tenants |
| **P3** | 19 | No real authentication | All |

---

## Summary

### Progress Since 4/28

One committed change (`62b64588 submissions | basic mongodb repository`) plus significant unstaged work across submissions and forms:

- **Submissions MongoDB repository fully implemented** -- `Find`, `FindByID`, and `FindByReferenceID` are real implementations backed by the generic `MongoDBRepository[submissionDocument]` base. `Upsert` is now functional using `FindOneAndUpdate` with `$set`, `SetUpsert(true)`, and `SetReturnDocument(After)` wrapped in a MongoDB session. Bidirectional BSON mapping (`toSubmissionDocument`/`fromSubmissionDocument`, `toSubmissionAttemptDocument`/`fromSubmissionAttemptDocument`) handles domain-to-document conversion. `Payload` and `ErrorDetails` are marshaled via `bson.Marshal`. This closes the P0 issue tracked since 4/25.

- **Submissions write path partially wired** -- `CreateSubmissionCommand` now carries `TenantID`, `FormID`, `VersionID`, and `Payload`. `SubmissionsService.Create` validates the command, constructs the domain aggregate via `NewSubmission`, and persists via `repository.Upsert`. The `createSubmission` handler extracts the tenant from context but does not yet deserialize a request body or invoke the service. The write path is functional at the service+repository layer but not yet exposed via HTTP.

- **In-memory `Upsert` implemented** -- Acquires a write lock and stores the submission by ID. No longer a stub.

- **Forms domain method rename** -- `UpdatePages`, `UpdateSections`, and `UpdateFields` renamed to `ReplacePages`, `ReplaceSections`, and `ReplaceFields`. The prior names implied a merge/patch semantic; the new names accurately reflect that these methods clear and replace all children. `FormsService.UpdateVersion` call site updated.

- **Service field naming cleanup** -- `SubmissionsService.submissionsRepository` renamed to `repository`, consistent with the single-repository pattern used by the forms service.

### Current State

**22 remaining issues (4/28) -> 19 remaining issues** (resolved 3: P0 MongoDB stubs, P2 Upsert stubs, P3 CreateSubmissionCommand empty; moved 3 to Will Not Fix; introduced 2 new issues: missing validation tags P2, `sendErrorResponse` no mapping P2).

**Forms Service** remains the strongest service. The `Replace*` rename improves semantic clarity. Rich domain model with version state machine, position-keyed sorted collections, complete error-to-HTTP mapping, and transactional version creation. Only remaining issue: hardcoded placeholder user ID (P2).

**Tenants Service** is stable with complete CRUD, cascade-delete, and lookup strategies. Minor formatting cleanup in `TenantsService.Create`. Remaining gaps: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has made meaningful progress: the persistence layer is now complete (both read and write), the service layer has a `Create` method, and the CQRS write-side command exists. However, the domain model remains fundamentally broken: `NewSubmission` doesn't assign `Status` or `Payload`, which means the new `Create` flow actively persists incomplete aggregates. The handler layer still needs to deserialize request bodies and invoke the service. Two handlers remain empty stubs.

**Hexagonal Architecture** -- Dependency direction remains correct. The `toSubmissionDocument`/`fromSubmissionDocument` mapping functions live in the persistence adapter, maintaining clean separation. The `Create` service method correctly depends on the repository port interface, not the concrete implementation. No cross-adapter imports.

**CQRS-Lite** -- `CreateSubmissionCommand` now follows the established command pattern with typed fields. The `Create` service method validates the command before delegating to domain construction. `ReplaySubmissionCommand` remains empty. The handler layer has not yet wired deserialization to command construction, creating a gap between the HTTP adapter and the application layer.

**DDD** -- The submissions write path now exists but the domain constructor is broken: `NewSubmission` produces an aggregate missing `Status` and `Payload`, violating the invariant that a newly-created aggregate should be in a valid state. This is now a live bug since `Create` -> `Upsert` will persist the incomplete object. The forms `Replace*` rename is a positive DDD signal -- method names now accurately describe aggregate behavior. Tenants `DataSourceAttributes` sealed interface and strategy pattern remain strong.

**Idiomatic Go** -- The `toSubmissionDocument` functions using `bson.Marshal` for untyped fields is a pragmatic approach but brittle -- if `Payload` contains a non-BSON-serializable type, it will fail at runtime with no compile-time safety. The service field rename from `submissionsRepository` to `repository` is cleaner when the service only depends on one repository. The `Replace*` naming follows Go's preference for precise verb choices.

### Highest-Impact Improvements

1. **Fix `NewSubmission` to assign `Status` and `Payload`** (P1 -- constructor produces incomplete domain objects, now actively persisted)
2. **Add validation tags to `CreateSubmissionCommand`** (P2 -- `validate.ValidateStruct` is a no-op without tags)
3. **Complete `createSubmission` handler** (P2 -- deserialize request body, construct command, call service)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
6. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
7. **Add test coverage** starting with domain constructors and the `Create` flow (P3 -- long-term reliability)
