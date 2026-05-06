# Full Codebase Review: Forms, Submissions, and Tenants Services

## Issues Resolved Since 5/4 Review

1. ~~Zero test files~~ (All, P3) -- Partially resolved. Route walking tests (`routes_test.go`) added for all three services verifying all expected routes and HTTP methods exist via `chi.Walk`. Comprehensive handler unit tests for tenants service (`handlers_test.go`) covering all 9 handler methods (getTenants, getTenant, createTenant, updateTenant, deleteTenant, getDataSources, getDataSource, createDataSource, updateDataSource, deleteDataSource, getLookups) plus `isBadRequest` helper. Clean mock pattern in `mocks_test.go` using function fields on mock structs implementing the port interfaces. *(Partially resolves 5/4 #8 -- forms and submissions services still lack handler/service tests.)*

2. ~~`os.Exit(1)` in `GetClaimsFromContext` and `TenantFromContext` is untestable~~ (Shared, P2) -- Both functions now use `panic(err)` instead of `os.Exit(1)`. This preserves fail-fast behavior for misconfiguration while making the functions testable (panics are recoverable via `recover()`) and not bypassing deferred cleanup. Both now create a proper `errors.New(...)` error value, log it via `slog.ErrorContext`, then panic. *(Resolves concern raised in 5/4 review cycle regarding testability of fail-fast infrastructure.)*

3. ~~Forms service has no query validation~~ (Forms, P2) -- `query.go` now defines `FindFormsQuery`, `FindFormsByIDQuery`, `FindVersionsQuery`, and `FindVersionByIDQuery` with `validate:"required"` tags and `Validate()` methods. `FormsService` methods (`Find`, `FindByID`, `FindVersions`, `FindVersion`, `CreateVersion`, `UpdateVersion`, `PublishVersion`, `RetireVersion`) call `query.Validate()` or `command.Validate()` at entry. Commands already had `Validate()` methods but `CreateVersionCommand` and `UpdateVersionCommand` now use named field initialization (`baseVersionCommand: baseVersionCommand{...}`) for correctness. *(Not previously tracked.)*

4. ~~Tenants service commands/queries missing validation~~ (Tenants, P2) -- `commands.go` now has proper `validate` struct tags on all command fields (`required` on `Name`, `TenantID`, `Type`, `Attributes`, `ID`). All commands have `Validate()` methods. `query.go` adds `Validate()` methods to `ListDataSourceQuery`, `FindDataSourceByIDQuery`, and `GetDataSourceLookupsQuery`. Services call `Validate()` at entry consistently. *(Not previously tracked.)*

5. ~~Tenants handlers lack structured logging~~ (Tenants, P3) -- Handlers now log `WarnContext` on context cancellation and invalid request bodies. `core.Application` exposes `Logger *slog.Logger` for handler access. *(Not previously tracked.)*

6. ~~Structured JSON logging not configured~~ (All, P3) -- All three `main.go` files now create `slog.New(slog.NewJSONHandler(os.Stdout, nil))` with `RequestContextHandler` wrapping to inject `request_id` from Chi's middleware into every log record. *(Enhancement to 5/4 #23's text handler.)*

7. ~~Forms DTOs lack validation struct tags~~ (Forms, P3) -- `PageRequest`, `SectionRequest`, `FieldRequest`, `RuleRequest`, `UpsertFormRequest` all have `validate` tags (`required`, `max`, `gte`, `lte`). Nested slice fields (`Pages`, `Sections`, `Fields`, `Rules`) now use `validate:"dive"` to recursively validate child elements. `UpsertVersionRequest.Pages` has `validate:"dive"`. `FieldRequest.Attributes` has `validate:"required"`. `ReadValidateJSONPayload` now rejects malformed request bodies at the adapter boundary. *(Not previously tracked.)*

8. ~~Tenants DTOs lack validation struct tags~~ (Tenants, P3) -- `TenantRequest` and `DataSourceRequest` have `validate` tags. *(Not previously tracked.)*

---

## Will Not Fix

See [5/4 review](code-review-5-4-26.md) for the full Will Not Fix list (items #25-27 covering `FindByReferenceID` linear scan, in-memory map key type, and context cancel response pattern).

---

## Remaining Issues

### Submissions Service

#### Bugs

1. **`CreateSubmissionCommand` has no validation tags** (`commands.go`) -- `SubmissionsService.Create` calls `validate.ValidateStruct(command)` but none of the fields (`TenantID`, `FormID`, `VersionID`, `Payload`) have `validate` struct tags. Validation is a no-op; empty commands will pass through to `NewSubmission`. *(Unresolved from 5/4 #1.)*

2. **`SubmissionRequest` DTO has no validation tags** (`dto/request.go`) -- `FormID`, `VersionID`, and `Payload` have only `json` tags. `ReadValidateJSONPayload` will not reject empty submissions at the adapter boundary. *(New.)*

#### Architectural

3. **`Replay` service method only validates and checks existence** (`submissions_service.go:79-88`) -- Validates the command and calls `repository.FindByID` to verify the submission exists, then returns nil without performing any replay logic. The handler returns 201 ("Successfully replayed") despite no replay occurring. *(Unresolved from 5/4 #2.)*

4. **`sendErrorResponse` has no domain error mapping** (`handlers.go`) -- Switch statement contains only a `default` case. All errors (including `common.ErrNotFound`, validation errors) fall through to `httputil.SendErrorResponse` which may not map them correctly. Unlike forms and tenants handlers, there's no `isBadRequest` helper. *(Unresolved from 5/4 #4.)*

#### Missing Functionality

5. **`SubmissionAttempt` has no constructor or factory function** -- Only `HydrateSubmissionAttempt` exists for reconstitution; no `NewSubmissionAttempt` for creation. Needed for replay implementation. *(Unresolved from 5/4 #3.)*

#### Code Quality

6. **`Payload` typed as `any`** (`submission.go`) -- No type safety. `ErrorDetails` on `SubmissionAttempt` also `any`. The `toSubmissionDocument` mapping uses `bson.Marshal` on `Payload` which will fail at runtime if `Payload` is not BSON-serializable. *(Unresolved from 5/4 #5.)*

---

### Tenants Service

#### Architectural

7. **`Find()` has no pagination or filtering** (`tenants_service.go:30-39`). The `TODO` comment in `ListDataSourceQuery` acknowledges this. *(Unresolved from 5/4 #6.)*

#### Missing Functionality

8. **`Lookup` value object has no validation** (`lookup.go`) -- `NewLookup` accepts any strings without checking for blank `Value` or `Label`. *(Unresolved from 5/4 #7.)*

---

### Shared Package

#### Code Quality

9. **`RequestContextHandler` does not implement `Enabled` method** (`logger/request_context_handler.go`) -- The custom `slog.Handler` wraps the inner handler but does not override `Enabled(context.Context, slog.Level) bool`. While `slog.Handler` interface only requires `Handle`, `WithAttrs`, `WithGroup`, and `Enabled`, the embedded `slog.Handler` will delegate `Enabled` correctly. However, the `Handle` method will be called even when the level is disabled because the type assertion in `slog.Logger` dispatches to the wrapper's `Handle` directly. Minor performance concern at high log volumes. *(New.)*

---

### Cross-Service

#### Architectural

10. **Test coverage limited to tenants REST adapter** -- Forms and submissions services have only route walking tests. No service layer, domain layer, or repository tests exist anywhere. *(Partially resolved from 5/4 #8.)*

11. **No domain events** for cross-service communication. *(Unresolved from 5/4 #9.)*

12. **No real authentication** -- `PlaceholderAuthenticator` always returns a fixed subject (`"placeholder"`). Only the forms service wires auth middleware. Submissions and tenants services have no authentication. *(Unresolved from 5/4 #10.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `CreateSubmissionCommand` no validation tags | Submissions |
| **P2** | 2 | `SubmissionRequest` DTO no validation tags | Submissions |
| **P2** | 4 | `sendErrorResponse` no domain error mapping | Submissions |
| **P3** | 3 | `Replay` validates but doesn't replay | Submissions |
| **P3** | 5 | `SubmissionAttempt` has no constructor | Submissions |
| **P3** | 6 | `any`-typed attributes (no type safety) | Submissions |
| **P3** | 7 | `Find()` no pagination | Tenants |
| **P3** | 8 | `Lookup` value object no validation | Tenants |
| **P3** | 9 | `RequestContextHandler` Enabled method | Shared |
| **P3** | 10 | Test coverage limited to tenants REST | All |
| **P3** | 11 | No domain events | All |
| **P3** | 12 | No real authentication (placeholder only) | All |

---

## Summary

### Progress Since 5/4

Eight commits since the last review plus unstaged changes, focused on validation hardening, structured logging, and test coverage:

- **Validation finalized across forms and tenants** -- All commands, queries, and DTOs now have proper `validate` struct tags. Service methods consistently call `Validate()` at entry. The forms service added four query types (`FindFormsQuery`, `FindFormsByIDQuery`, `FindVersionsQuery`, `FindVersionByIDQuery`) with embedded composition for shared fields. Forms DTOs use `validate:"dive"` on nested slice fields (`Pages`, `Sections`, `Fields`, `Rules`) to recursively validate children. `FieldRequest.Attributes` is `validate:"required"`. The tenants service commands now validate all required fields. This closes a class of issues where `validate.ValidateStruct` was being called but had nothing to validate.

- **Test coverage introduced** -- Route walking tests verify all three services register expected routes. The tenants service has comprehensive handler unit tests (9 test functions, table-driven, covering success/error/not-found paths). Mock pattern uses function fields on struct types implementing port interfaces -- simple, explicit, no external mock library needed. Tests verify HTTP status codes and response body structure.

- **Structured JSON logging with request context** -- All services now use `slog.NewJSONHandler` wrapped with a custom `RequestContextHandler` that injects Chi's `request_id` into every log record. Tenants handlers log warnings on context cancellation and invalid request bodies. The `core.Application` struct exposes `Logger` for handler-level logging.

- **Fail-fast changed from `os.Exit(1)` to `panic`** -- `GetClaimsFromContext` and `TenantFromContext` now panic instead of calling `os.Exit(1)` on misconfiguration. This preserves fail-fast semantics while enabling testability (panics are recoverable) and not bypassing deferred cleanup (e.g., database connection close, graceful shutdown).

### Current State

**10 remaining issues (5/4) -> 12 remaining issues** (resolved 4 from prior review + 4 newly tracked improvements; introduced 2 new issues; carried forward 6 unchanged).

**Forms Service** is fully mature. Validation complete at all layers: DTO tags with `validate:"dive"` on nested slices (`Pages`, `Sections`, `Fields`, `Rules`), `validate:"required"` on `FieldRequest.Attributes`, command `Validate()` methods, query `Validate()` methods, and domain constructor validation. No remaining issues.

**Tenants Service** is the best-tested service with comprehensive handler unit tests. Validation finalized. Structured logging in handlers. Remaining gaps are feature-level: pagination (P3) and `Lookup` validation (P3).

**Submissions Service** remains the weakest. It has functional persistence and a working create flow, but lacks validation at both the DTO and command layers (P2). The `sendErrorResponse` has no domain error mapping (P2). `Replay` is effectively a no-op (P3). No handler or service tests beyond route walking.

**Hexagonal Architecture** -- The test structure correctly tests at the adapter boundary: handler tests mock the service port interface, not the repository. This validates that the port boundary is well-defined and testable. The `core.Application` exposing `Logger` to handlers is acceptable since logging is infrastructure, not domain.

**CQRS-Lite** -- Forms now has complete query/command separation with dedicated query types, constructors, and validation. Tenants follows the same pattern. Submissions queries use a generic `FindSubmissionByIDQuery[T]` which is clever but slightly unusual -- the type parameter enables reuse for both `SubmissionID` and `ReferenceID` lookups without duplicating the struct.

**DDD** -- Domain validation is now enforced at construction time across all aggregates (forms, tenants, data sources). The `validate` tags on domain struct fields (`form.go`, `tenant.go`, `data_source.go`, `version.go`, `page.go`, `section.go`, `field.go`) ensure invariants are checked via `validate.ValidateStruct` in constructors. Submissions domain is the exception -- `NewSubmission` validates the struct but `CreateSubmissionCommand` does not validate its inputs before reaching the domain.

**Idiomatic Go** -- The test pattern (table-driven, `httptest.NewRecorder`, function-field mocks) is idiomatic and avoids external dependencies. The `panic` for impossible-state failures aligns with Go conventions (cf. `regexp.MustCompile`). `slog` usage with custom handlers follows Go 1.21+ best practices. The `validate` tag approach is the standard pattern with `go-playground/validator`.

### Highest-Impact Improvements

1. **Add validation tags to `CreateSubmissionCommand` and `SubmissionRequest`** (P2 -- entire submissions write path accepts empty input)
2. **Add `sendErrorResponse` domain error cases** in submissions handlers (P2 -- `ErrNotFound` and validation errors map to 500)
3. **Add handler unit tests for forms and submissions** (P3 -- tenants pattern is established; replicate)
4. **Implement actual `Replay` logic** (P3 -- handler returns success but nothing happens)
5. **Add `NewSubmissionAttempt` constructor** (P3 -- needed for replay implementation)
