# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/6 Review

1. ~~`MongoDBDatabase` lacks structured logging~~ (Shared, P3) -- `BeginTx`, `CommitTx`, and `RollbackTx` now log `DebugContext` at entry and `ErrorContext` on failure. `NewMongoDBDatabase` accepts a `*slog.Logger` parameter. All three service bootstrap call sites updated.

2. ~~`NewFindDataSourceByID` naming inconsistent with convention~~ (Tenants, P3) -- Renamed to `NewFindDataSourceByIDQuery` to match the naming pattern of all other query constructors (e.g., `NewFindSubmissionsQuery`, `NewFindSubmissionByIDQuery`).

3. ~~Test coverage limited to tenants REST adapter~~ (All, P3) -- Partially resolved further. Forms service now has comprehensive handler unit tests (834 lines, `handlers_test.go`) with mock service port using the function-field pattern (`mocks_test.go`). Tenants service now has service-layer unit tests: `tenants_service_test.go` (325 lines) and `data_sources_service_test.go` (431 lines) with repository mocks (`mocks_test.go`). Coverage now spans: route walking (all 3), handler tests (forms + tenants), service tests (tenants). Submissions still has only route walking tests. *(Partially resolves 5/6 #9.)*

4. ~~Submissions MongoDB repository lacks unique index for idempotency~~ (Submissions, P2) -- `indexes` var added with a unique index on `idempotency_id`. The `Upsert` method now checks `mongo.IsDuplicateKeyError(err)` and returns `domain.ErrDuplicateSubmission`, providing database-level deduplication enforcement. *(New issue found and fixed same cycle.)*

5. ~~Tenants handler test used `len` as struct field name~~ (Tenants, P3) -- Renamed to `count` to avoid shadowing the builtin. Test helper `newTestHandlers` now provides a logger to `core.Application` (previously nil, could panic on handler log calls). *(Not previously tracked.)*

6. ~~`CreateSubmissionCommand` has no validation tags~~ (Submissions, P2) -- All fields (`TenantID`, `FormID`, `VersionID`, `IdempotencyID`, `Payload`) now have `validate:"required"` struct tags. `validate.ValidateStruct(command)` in the service will now reject empty commands. *(Resolves 5/6 #1.)*

7. ~~`SubmissionRequest` DTO has no validation tags~~ (Submissions, P2) -- `FormID` and `VersionID` now have `validate:"required,uuidv7"` and `Payload` has `validate:"required"`. `ReadValidateJSONPayload` will reject invalid submissions at the adapter boundary. *(Resolves 5/6 #2.)*

8. ~~`Submission` domain struct has no validation tags~~ (Submissions, P2) -- `TenantID`, `FormID`, `VersionID` now have `validate:"required"`, `IdempotencyID` has `validate:"required"`, and `Payload` has `validate:"required"`. `NewSubmission` now enforces domain invariants at construction time. *(Resolves new issue from this review cycle.)*

---

## Will Not Fix

See [5/4 review](code-review-5-4-26.md) for the full Will Not Fix list (items #25-27 covering `FindByReferenceID` linear scan, in-memory map key type, and context cancel response pattern).

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`createSubmission` handler hardcodes `IdempotencyID("")`** (`handlers.go`) -- With the addition of `validate:"required"` on `IdempotencyID` at command and domain levels, the empty string now **fails validation**, meaning all submission creation requests are rejected. The handler must extract an idempotency key from a request header (e.g., `Idempotency-Key`) or DTO field. This is now a **blocking bug** -- no submissions can be created. *(Elevated from 5/6 #4; was latent, now active.)*

#### Architectural

2. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. `common.ErrNotFound` maps to 500 instead of 404; `common.ErrUnauthorized` maps to 500 instead of 403; validation errors map to 500 instead of 400. Unlike forms and tenants handlers, there's no `isBadRequest` helper. *(Unresolved from 5/6 #3.)*

3. **`Replay` service method is a stub** (`submissions_service.go`) -- Validates the command and calls `repository.FindByID` to verify the submission exists, then returns nil without performing any replay logic. The handler returns 201 ("Successfully replayed") despite no replay occurring. *(Unresolved from 5/6 #5.)*

#### Code Quality

4. **`Payload` typed as `any`** (`submission.go`) -- No type safety across DTO, command, domain, and persistence layers. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 5/6 #6.)*

---

### Tenants Service

#### Architectural

5. **`Find()` has no pagination or filtering** (`tenants_service.go`). The `TODO` comment in `ListDataSourceQuery` acknowledges this. *(Unresolved from 5/6 #7.)*

#### Missing Functionality

6. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 5/6 #8.)*

---

### Cross-Service

#### Architectural

7. **Test coverage gaps remain** -- Submissions has no handler or service tests. Forms has no service tests. No domain-layer or repository-layer tests exist anywhere. *(Partially resolved from 5/6 #9.)*

8. **No domain events** for cross-service communication. *(Unresolved from 5/6 #10.)*

9. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placeholder"`). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Unresolved from 5/6 #11.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 1 | `createSubmission` hardcodes empty `IdempotencyID` (blocks all creation) | Submissions |
| **P2** | 2 | `sendErrorResponse` no domain error mapping | Submissions |
| **P3** | 3 | `Replay` is a stub | Submissions |
| **P3** | 4 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 5 | `Find()` no pagination | Tenants |
| **P3** | 6 | `Lookup` value object no validation | Tenants |
| **P3** | 7 | Test coverage gaps | All |
| **P3** | 8 | No domain events | All |
| **P3** | 9 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 5/6

Six commits since the last review, focused on test coverage expansion, database-layer observability, idempotency enforcement, and validation completion:

- **Forms handler unit tests introduced** -- 834-line `handlers_test.go` covering all form and version CRUD handlers (create, update, delete, find, publish, retire) with table-driven tests, `httptest.NewRecorder`, and function-field mocks implementing `FormsService`. Follows the same pattern established in tenants handler tests. `mocks_test.go` provides the mock struct.

- **Tenants service-layer unit tests introduced** -- `tenants_service_test.go` (325 lines) and `data_sources_service_test.go` (431 lines) test the business logic layer directly, mocking repository ports. This validates domain authorization checks, error propagation, and service orchestration independently of the HTTP adapter. `mocks_test.go` (98 lines) provides repository mock structs.

- **`MongoDBDatabase` transaction logging** -- `BeginTx`, `CommitTx`, and `RollbackTx` now log `DebugContext` at entry and `ErrorContext` on failure. The `NewMongoDBDatabase` constructor accepts a `*slog.Logger`, completing the structured logging story for the shared database infrastructure. Direct-return patterns replaced with explicit error handling for full observability.

- **Submissions idempotency index** -- A unique MongoDB index on `idempotency_id` added to the submissions collection. The `Upsert` method now detects `mongo.IsDuplicateKeyError` and returns `domain.ErrDuplicateSubmission`. This provides database-level enforcement as a safety net beyond the application-level `FindByIdempotencyID` check.

- **Query constructor naming consistency** -- `NewFindDataSourceByID` renamed to `NewFindDataSourceByIDQuery` aligning with the established convention across all services where query constructors are suffixed with `Query`.

- **Test infrastructure hardening** -- Tenants handler tests fixed: `len` field renamed to `count` (avoids builtin shadowing), test helper `newTestHandlers` now injects a logger into `core.Application` preventing nil pointer panics when handlers log.

- **Submissions validation finalized** -- `CreateSubmissionCommand` fields (`TenantID`, `FormID`, `VersionID`, `IdempotencyID`, `Payload`) now have `validate:"required"` tags. `SubmissionRequest` DTO has `validate:"required,uuidv7"` on `FormID`/`VersionID` and `validate:"required"` on `Payload`. `Submission` domain struct has `validate:"required"` on `TenantID`, `FormID`, `VersionID`, `IdempotencyID`, and `Payload`. Validation is now enforced at all three layers (adapter, application, domain). However, this introduces a regression: the handler hardcodes `IdempotencyID("")` which now fails the `required` validation, blocking all submission creation.

### Current State

**11 remaining issues (5/6) -> 9 remaining issues** (resolved 8 from prior review; introduced 1 new P1 regression; carried forward 8 unchanged).

**Forms Service** remains fully mature. No remaining issues. Handler tests now provide adapter-boundary coverage validating HTTP status codes, response structure, and error propagation through the service port.

**Tenants Service** is now tested at two layers: handler tests (adapter boundary) and service tests (domain orchestration). The service tests validate tenant authorization, error propagation from repositories, and correct delegation to repository methods. Remaining gaps are feature-level: pagination (P3) and `Lookup` validation (P3).

**Submissions Service** has made significant progress this cycle. Validation is now complete at all layers (DTO, command, domain), and the idempotency infrastructure has database-level enforcement via a unique index. However, there is a **P1 regression**: the handler hardcodes an empty `IdempotencyID("")` which now fails `validate:"required"` at the command layer, meaning **no submissions can be created**. The fix is straightforward -- extract the idempotency key from a request header (e.g., `Idempotency-Key`) or add it to the DTO. Beyond this blocker: `sendErrorResponse` still lacks domain error mapping (P2), `Replay` is a stub (P3), and no handler or service tests exist.

**Hexagonal Architecture** -- The test structure correctly validates port boundaries at two levels: handler tests mock service ports (testing the REST adapter in isolation), and service tests mock repository ports (testing domain orchestration in isolation). This confirms the hexagonal architecture enables independent testability at each layer. The `MongoDBDatabase` logging addition is appropriate since `Database` is an infrastructure port, not a domain concern.

**DDD** -- All three services now enforce domain invariants at construction time via `validate` struct tags. The submissions domain model gap identified in the prior review cycle is closed: `NewSubmission` validates `TenantID`, `FormID`, `VersionID`, `IdempotencyID`, and `Payload` are present. The `validate:"required,uuidv7"` tags on the DTO layer provide format validation (UUIDv7) at the adapter boundary, while the domain layer enforces presence. This is the correct layering: format validation at the edge, invariant enforcement in the domain.

**Idiomatic Go** -- The test patterns continue to follow Go conventions: table-driven subtests, `httptest.NewRecorder`, function-field mocks avoiding external dependencies, and `t.Run` for subtest naming. The `count` rename from `len` correctly avoids builtin shadowing. The explicit error handling in `MongoDBDatabase` (replacing direct-return patterns) follows the Go principle of handling errors at the point of occurrence for observability. The `validate:"required,uuidv7"` usage on DTO fields is idiomatic `go-playground/validator` -- format constraints belong at the adapter boundary, not the domain.

### Highest-Impact Improvements

1. **Wire `IdempotencyID` from request header to handler** (P1 -- blocks all submission creation; the empty string fails `validate:"required"`)
2. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- `ErrNotFound`, `ErrUnauthorized`, and validation errors all map to 500)
3. **Add handler and service unit tests for submissions** (P3 -- pattern established in forms and tenants; replicate)
4. **Implement actual `Replay` logic** (P3 -- handler returns success but nothing happens)
