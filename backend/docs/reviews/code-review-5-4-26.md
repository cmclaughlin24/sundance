# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 4/28 Review

1. ~~MongoDB repository methods are stubs returning `nil, nil`~~ (Submissions, P0) -- `Find`, `FindByID`, and `FindByReferenceID` (`submissions_repository.go:27-60`) are now real implementations using the `MongoDBRepository[submissionDocument]` base. They query with `bson.M{}`, `bson.M{"_id": id}`, and `bson.M{"reference_id": id}` respectively. A new `documents.go` file defines `submissionDocument` and `submissionAttemptDocument` BSON structs with `fromSubmissionDocument` and `fromSubmissionAttemptDocument` mapping functions. Attempts are embedded as a subdocument array. The nil-pointer dereference path is eliminated because `MongoDBRepository.FindOne` maps `ErrNoDocuments` to `common.ErrNotFound`, which is checked before tenant access. *(Resolves 4/28 #2.)*

2. ~~No write operations in the repository interface~~ (Submissions, P3) -- `SubmissionsRepository` (`ports/secondary.go:19`) now defines `Upsert`. Both in-memory and MongoDB implementations exist, though both are stubs returning `nil, nil`. *(Partially resolves 4/28 #8 -- interface defined but implementations are stubs; see #7 below.)*

---

## Will Not Fix

See [4/25 review](code-review-4-25-26.md) and [4/24 review](code-review-4-24-26.md) for the prior Will Not Fix list (10 items).

11. **`FindByReferenceID` does a linear scan** in the in-memory repository -- acceptable for a repository not intended for production use. *(Closed from 4/28 #13.)*

12. **In-memory submissions repository map keyed by `string` instead of `SubmissionID`** -- consistent with how the forms and tenants in-memory repositories are implemented. *(Closed from 4/28 #15.)*

---

## Remaining Issues

### Forms Service

#### Bugs

1. **`publishVersion` and `retireVersion` use hardcoded `"placeholder"` user ID** (`handlers.go:373`, `408`) -- The publish/retire state transitions record a fake user. Both have `// FIXME` comments. *(Unresolved from 4/28 #1.)*

---

### Submissions Service

#### Bugs

2. **`NewSubmission` never sets `Status` or `Payload`** (`submission.go:35-44`) -- The constructor accepts `payload any` but never assigns it to the struct. `Status` is left as zero-value empty string with a TODO comment. Both fields persist as empty/nil despite being provided by the caller. *(Unresolved from 4/28 #3.)*

#### Architectural

3. **`Find()` has no tenant filtering** (`submissions_service.go:25-27`) -- Returns all submissions across all tenants. The `getSubmissions` handler does not call `getTenantFromContext`. *(Unresolved from 4/28 #4.)*

4. **Three handler stubs return 200 OK with empty body** (`handlers.go:87-91`) -- `createSubmission`, `getSubmissionStatus`, and `replaySubmission`. *(Reduced from 4/28 #5 -- `getSubmissionAttempts` now implemented.)*

#### Missing Functionality

5. **`Replay` service method is a stub** (`submissions_service.go`) -- Returns `nil`. *(Unresolved from 4/28 #6.)*

6. **Request DTOs not implemented** -- `dto/request.go` contains only the package declaration. *(Unresolved from 4/28 #7.)*

7. **`Upsert` is a stub in both repositories** -- In-memory (`return nil, nil`) and MongoDB (`return nil, nil`). Blocks the entire write path. *(New -- interface added per 4/28 #8 but implementations remain stubs.)*

8. **`ReplaySubmissionCommand` is an empty struct** (`commands.go`) -- Has no fields. *(Unresolved from 4/28 #9.)*

9. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. *(Unresolved from 4/28 #10.)*

10. **`sendErrorResponse` has no domain error mapping** (`handlers.go:108-112`) -- Switch statement contains only a `default` case. All errors map to 500. *(New.)*

#### Code Quality

11. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. *(Unresolved from 4/28 #11.)*

12. **`SubmissionStatus` has no defined constants** -- `type SubmissionStatus string` declared but no `const` block. *(Unresolved from 4/28 #12.)*

13. **Context cancel drops response silently** (`handlers.go:39-41`, `74-76`) -- When context is cancelled, `select` on `r.Context().Done()` returns without writing any HTTP response. *(Unresolved from 4/28 #14.)*

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
| **P1** | 2 | `NewSubmission` never sets `Status` or `Payload` | Submissions |
| **P2** | 1 | Hardcoded `"placeholder"` user ID | Forms |
| **P2** | 3 | `Find()` has no tenant filtering | Submissions |
| **P2** | 4 | Handler stubs return 200 OK | Submissions |
| **P2** | 7 | `Upsert` stubs block write path | Submissions |
| **P2** | 10 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 16 | `ErrMissingTenantID` maps to 500 | Shared |
| **P3** | 17 | Zero test files | All |
| **P3** | 18 | No domain events | All |
| **P3** | 9 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 11 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 12 | `SubmissionStatus` no constants | Submissions |
| **P3** | 15 | `Lookup` value object no validation | Tenants |
| **P3** | 19 | No real authentication | All |

---

## Summary

### Progress Since 4/28

The sole functional commit since the 4/28 review (`62b64588 submissions | basic mongodb repository`) implemented the MongoDB persistence layer for submissions:

- **Submissions MongoDB repository implemented** -- `Find`, `FindByID`, and `FindByReferenceID` are now real implementations backed by the generic `MongoDBRepository[submissionDocument]` base. A `documents.go` file provides the BSON document structures (`submissionDocument`, `submissionAttemptDocument`) with proper `bson` tags and bidirectional mapping functions (`fromSubmissionDocument`, `fromSubmissionAttemptDocument`). Attempts are modeled as an embedded array subdocument within the submission document. This closes the P0 issue that has been tracked since 4/25.

- **`Upsert` method added to repository interface** -- `SubmissionsRepository` now includes an `Upsert(ctx, *Submission) (*Submission, error)` method. However, both the in-memory and MongoDB implementations are stubs returning `nil, nil`. The write path remains non-functional end-to-end.

### Current State

**22 remaining issues (4/28) -> 19 remaining issues** (resolved 1 P0 fully, partially resolved 1 P3; moved 2 to Will Not Fix; introduced 2 new issues: `Upsert` stubs P2, `sendErrorResponse` no mapping P2).

**Forms Service** remains the strongest service. Rich domain model with version state machine, position-keyed sorted collections, complete error-to-HTTP mapping, and transactional version creation. Only remaining issue: hardcoded placeholder user ID (P2).

**Tenants Service** is stable with complete CRUD, cascade-delete, and lookup strategies. Remaining gaps are minor: `Lookup` value object validation (P3), `Find()` pagination (P3).

**Submissions Service** has improved at the persistence layer but remains the weakest service by a significant margin. The P0 nil-panic issue is resolved, but the domain model is still largely anemic: `NewSubmission` doesn't assign its own fields, `SubmissionStatus` has no constants, `Payload` is untyped, the write path is entirely stubbed, and 3 of 6 handlers are empty bodies. The service needs focused attention to reach parity with forms and tenants.

**Hexagonal Architecture** -- Dependency direction remains correct. The new `submissionDocument` and mapping functions live in the persistence adapter, keeping domain types clean of infrastructure concerns. The `MongoDBRepository[T]` generic base continues to reduce adapter boilerplate effectively. No cross-adapter imports.

**CQRS-Lite** -- The generic `FindByIDQuery[T]` with constructor validation follows the established pattern cleanly. Commands and queries remain well-separated. However, the lack of a `CreateSubmissionCommand` (to complement the `createSubmission` handler) is a gap -- the write side of CQRS is essentially unimplemented for submissions.

**DDD** -- The submissions domain remains the weakest DDD implementation. Key deficiencies: (1) `NewSubmission` produces an incomplete aggregate -- missing `Status` and `Payload` violates the invariant that a newly-created aggregate should be in a valid state; (2) No `SubmissionStatus` constants means the state lifecycle is undefined; (3) `Payload` as `any` provides no bounded context contract; (4) No domain events means aggregate state transitions are invisible to other services. By contrast, forms (`Version` state machine, `withPosition` mixin, sorted iteration) and tenants (`DataSourceAttributes` sealed interface, strategy pattern) continue to exemplify strong DDD patterns.

**Idiomatic Go** -- The BSON document mapping approach (separate document structs with `from*` converters) is idiomatic and keeps domain types free of persistence tags. The `MongoDBRepository[T]` generic continues to work well. The `sendErrorResponse` function in submissions handlers should follow the same switch-on-error pattern established in the forms service and shared `httputil.SendErrorResponse`.

### Highest-Impact Improvements

1. **Fix `NewSubmission` to assign `Status` and `Payload`** (P1 -- constructor produces incomplete domain objects)
2. **Implement `Upsert` in both repositories** (P2 -- blocks entire write path)
3. **Implement `createSubmission` handler end-to-end** (P2 -- requires command, service method, and repository write)
4. **Fix `ErrMissingTenantID` mapping** in `SendErrorResponse` (P2 -- missing tenant header produces 500)
5. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- all errors currently map to 500)
6. **Add tenant filtering to submissions `Find()`** (P2 -- data isolation)
7. **Add test coverage** starting with domain constructors and repository implementations (P3 -- long-term reliability)
