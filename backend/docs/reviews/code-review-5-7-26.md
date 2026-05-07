# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/6 Review

1. ~~`MongoDBDatabase` lacks structured logging~~ (Shared, P3) -- `BeginTx`, `CommitTx`, and `RollbackTx` now log `DebugContext` at entry and `ErrorContext` on failure. `NewMongoDBDatabase` accepts a `*slog.Logger` parameter. All three service bootstrap call sites updated.

2. ~~`NewFindDataSourceByID` naming inconsistent with convention~~ (Tenants, P3) -- Renamed to `NewFindDataSourceByIDQuery` to match the naming pattern of all other query constructors (e.g., `NewFindSubmissionsQuery`, `NewFindSubmissionByIDQuery`).

3. ~~Test coverage limited to tenants REST adapter~~ (All, P3) -- Partially resolved further. Forms service now has comprehensive handler unit tests (834 lines, `handlers_test.go`) with mock service port using the function-field pattern (`mocks_test.go`). Tenants service now has service-layer unit tests: `tenants_service_test.go` (325 lines) and `data_sources_service_test.go` (431 lines) with repository mocks (`mocks_test.go`). Coverage now spans: route walking (all 3), handler tests (forms + tenants), service tests (tenants). Submissions still has only route walking tests. *(Partially resolves 5/6 #9.)*

4. ~~Submissions MongoDB repository lacks unique index for idempotency~~ (Submissions, P2) -- `indexes` var added with a unique index on `idempotency_id`. The `Upsert` method now checks `mongo.IsDuplicateKeyError(err)` and returns `domain.ErrDuplicateSubmission`, providing database-level deduplication enforcement. *(New issue found and fixed same cycle.)*

5. ~~Tenants handler test used `len` as struct field name~~ (Tenants, P3) -- Renamed to `count` to avoid shadowing the builtin. Test helper `newTestHandlers` now provides a logger to `core.Application` (previously nil, could panic on handler log calls). *(Not previously tracked.)*

---

## Will Not Fix

See [5/4 review](code-review-5-4-26.md) for the full Will Not Fix list (items #25-27 covering `FindByReferenceID` linear scan, in-memory map key type, and context cancel response pattern).

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `TenantID`, `FormID`, `VersionID`, `IdempotencyID`, `Payload` have no `validate` struct tags. `validate.ValidateStruct(command)` is a no-op; empty commands will pass through to `NewSubmission`. *(Unresolved from 5/6 #1.)*

2. **`SubmissionRequest` DTO has no validation tags** (`dto/request.go`) -- `FormID`, `VersionID`, and `Payload` have only `json` tags. `ReadValidateJSONPayload` will not reject empty submissions at the adapter boundary. *(Unresolved from 5/6 #2.)*

3. **`Submission` domain struct has no validation tags** (`submission.go`) -- `validate.ValidateStruct(s)` in `NewSubmission` is also a no-op. Empty `TenantID`, `FormID`, and `VersionID` will pass domain construction without error. *(New -- previously obscured by the command-level issue.)*

#### Architectural

4. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. `common.ErrNotFound` maps to 500 instead of 404; `common.ErrUnauthorized` maps to 500 instead of 403; validation errors map to 500 instead of 400. Unlike forms and tenants handlers, there's no `isBadRequest` helper. *(Unresolved from 5/6 #3.)*

5. **`createSubmission` handler hardcodes `IdempotencyID("")`** (`handlers.go`) -- Despite the unique index and `FindByIdempotencyID` infrastructure, the handler never extracts an idempotency key from request headers or body. The empty string means all submissions share the same key -- the unique index would reject the second submission entirely, or `FindByIdempotencyID` would always find the first-ever submission. The idempotency key should come from a request header (e.g., `Idempotency-Key`) or a DTO field. *(Unresolved from 5/6 #4.)*

6. **`Replay` service method is a stub** (`submissions_service.go`) -- Validates the command and calls `repository.FindByID` to verify the submission exists, then returns nil without performing any replay logic. The handler returns 201 ("Successfully replayed") despite no replay occurring. *(Unresolved from 5/6 #5.)*

#### Code Quality

7. **`Payload` typed as `any`** (`submission.go`) -- No type safety across DTO, command, domain, and persistence layers. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 5/6 #6.)*

---

### Tenants Service

#### Architectural

8. **`Find()` has no pagination or filtering** (`tenants_service.go`). The `TODO` comment in `ListDataSourceQuery` acknowledges this. *(Unresolved from 5/6 #7.)*

#### Missing Functionality

9. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 5/6 #8.)*

---

### Cross-Service

#### Architectural

10. **Test coverage gaps remain** -- Submissions has no handler or service tests. Forms has no service tests. No domain-layer or repository-layer tests exist anywhere. *(Partially resolved from 5/6 #9.)*

11. **No domain events** for cross-service communication. *(Unresolved from 5/6 #10.)*

12. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placeholder"`). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Unresolved from 5/6 #11.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 2 | `SubmissionRequest` DTO no validation tags | Submissions |
| **P2** | 3 | `Submission` domain struct no validation tags | Submissions |
| **P2** | 4 | `sendErrorResponse` no domain error mapping | Submissions |
| **P2** | 5 | `createSubmission` hardcodes empty `IdempotencyID` | Submissions |
| **P3** | 6 | `Replay` is a stub | Submissions |
| **P3** | 7 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 8 | `Find()` no pagination | Tenants |
| **P3** | 9 | `Lookup` value object no validation | Tenants |
| **P3** | 10 | Test coverage gaps | All |
| **P3** | 11 | No domain events | All |
| **P3** | 12 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 5/6

Five commits since the last review, focused on test coverage expansion, database-layer observability, and idempotency enforcement:

- **Forms handler unit tests introduced** -- 834-line `handlers_test.go` covering all form and version CRUD handlers (create, update, delete, find, publish, retire) with table-driven tests, `httptest.NewRecorder`, and function-field mocks implementing `FormsService`. Follows the same pattern established in tenants handler tests. `mocks_test.go` provides the mock struct.

- **Tenants service-layer unit tests introduced** -- `tenants_service_test.go` (325 lines) and `data_sources_service_test.go` (431 lines) test the business logic layer directly, mocking repository ports. This validates domain authorization checks, error propagation, and service orchestration independently of the HTTP adapter. `mocks_test.go` (98 lines) provides repository mock structs.

- **`MongoDBDatabase` transaction logging** -- `BeginTx`, `CommitTx`, and `RollbackTx` now log `DebugContext` at entry and `ErrorContext` on failure. The `NewMongoDBDatabase` constructor accepts a `*slog.Logger`, completing the structured logging story for the shared database infrastructure. Direct-return patterns replaced with explicit error handling for full observability.

- **Submissions idempotency index** -- A unique MongoDB index on `idempotency_id` added to the submissions collection. The `Upsert` method now detects `mongo.IsDuplicateKeyError` and returns `domain.ErrDuplicateSubmission`. This provides database-level enforcement as a safety net beyond the application-level `FindByIdempotencyID` check.

- **Query constructor naming consistency** -- `NewFindDataSourceByID` renamed to `NewFindDataSourceByIDQuery` aligning with the established convention across all services where query constructors are suffixed with `Query`.

- **Test infrastructure hardening** -- Tenants handler tests fixed: `len` field renamed to `count` (avoids builtin shadowing), test helper `newTestHandlers` now injects a logger into `core.Application` preventing nil pointer panics when handlers log.

### Current State

**11 remaining issues (5/6) -> 12 remaining issues** (resolved 5 from prior review; introduced 1 new issue; carried forward 11 unchanged).

**Forms Service** remains fully mature. No remaining issues. Handler tests now provide adapter-boundary coverage validating HTTP status codes, response structure, and error propagation through the service port.

**Tenants Service** is now tested at two layers: handler tests (adapter boundary) and service tests (domain orchestration). The service tests validate tenant authorization, error propagation from repositories, and correct delegation to repository methods. Remaining gaps are feature-level: pagination (P3) and `Lookup` validation (P3).

**Submissions Service** remains the primary area needing attention. 7 of 12 remaining issues are concentrated here. The validation gap (issues #1-3) is the highest-impact fix: it spans DTO, command, and domain layers, meaning the entire write path accepts empty/invalid input without rejection. The idempotency infrastructure is now enforced at the database level (unique index + duplicate key detection) but remains non-functional from the API layer because the handler hardcodes an empty key. No handler or service tests exist.

**Hexagonal Architecture** -- The test structure correctly validates port boundaries at two levels: handler tests mock service ports (testing the REST adapter in isolation), and service tests mock repository ports (testing domain orchestration in isolation). This confirms the hexagonal architecture enables independent testability at each layer. The `MongoDBDatabase` logging addition is appropriate since `Database` is an infrastructure port, not a domain concern.

**DDD** -- The submissions domain model has a structural gap: `NewSubmission` calls `validate.ValidateStruct` but the `Submission` struct has no `validate` tags, making domain construction accept any input. This violates the DDD principle that aggregates should enforce their own invariants at construction time. Forms and tenants services correctly enforce invariants via `validate` tags on their domain structs.

**Idiomatic Go** -- The test patterns continue to follow Go conventions: table-driven subtests, `httptest.NewRecorder`, function-field mocks avoiding external dependencies, and `t.Run` for subtest naming. The `count` rename from `len` correctly avoids builtin shadowing. The explicit error handling in `MongoDBDatabase` (replacing direct-return patterns) follows the Go principle of handling errors at the point of occurrence for observability.

### Highest-Impact Improvements

1. **Add validation tags to `CreateSubmissionCommand`, `SubmissionRequest`, and `Submission`** (P2 -- entire submissions write path accepts empty input at every layer)
2. **Wire `IdempotencyID` from request header to handler** (P2 -- infrastructure exists at service, repository, and database levels but is never activated from the API)
3. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- `ErrNotFound`, `ErrUnauthorized`, and validation errors all map to 500)
4. **Add handler and service unit tests for submissions** (P3 -- pattern established in forms and tenants; replicate)
5. **Implement actual `Replay` logic** (P3 -- handler returns success but nothing happens)
