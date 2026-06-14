# Full Codebase Review: Forms and Tenants Services

**Date:** June 14, 2026

---

## Issues Resolved Since 6/5 Review

1. ~~No real authentication -- all services use `PlaceholderAuthenticator`~~ (6/5 #11) -- `PingFedTokenValidator` added with JWK-based JWT validation: audience, issuer, expiry required, `WithIssuedAt()`, 5s leeway. `NewBearerAuthenticator[T]` is a generic wrapper that extracts the `Authorization: Bearer <token>` header and delegates to any `TokenValidator[T]`. Both `main.go` entrypoints now pass `settings.Auth` (from `settings.json`) to `rest.NewRoutes`. `PlaceholderTokenValidator` is retained for local/dev use. Unstaged: `autenticatorSettings` typo corrected to `authenticatorSettings` throughout `authenticators.go` and `pingfed_token_validator.go`.

2. ~~`auth.NewMiddleware` imports `reflect` for a nil/zero claims check~~ -- Unstaged: `reflect` import removed; `valOfClaims.IsValid() && !valOfClaims.IsZero()` replaced with `claims != nil`. Simpler, idiomatic, and unambiguous for interface values.

3. ~~Commands in monolithic `ports/commands.go` files~~ -- Forms `ports/commands.go` (269 lines) split into `commands/form.go`, `commands/form_version.go`, `commands/submission.go`, `commands/tag.go`, `commands/tag_version.go`, `commands/delete.go`. Tenants `ports/commands.go` split into `commands/tenant.go`, `commands/data_source.go`. All service methods, handlers, tests, and mocks updated to the new import paths. Clean build and all tests pass.

4. ~~`form_definition_mapper` did not preserve existing IDs on update and had no key uniqueness enforcement~~ -- `updateFormDefinition` now looks up existing pages/sections/fields by ID before creating new ones, preserving identity across update requests. A `formKeys` map tracks every page, section, and field key; `validate()` returns `ErrDuplicateFormKey` on collision. An `expressionKeys` map verifies that every rule expression `FieldKey` references an element that exists in the form definition, returning `ErrInvalidExpressionKey` otherwise. Both errors are surfaced as bad-request responses in `handlers.go`. Unstaged: error message inside `updatePage` corrected from `"page id=%s", *p.ID` to `"section id=%s", *s.ID`.

5. ~~`DataSourceRef.DataSourceId` non-idiomatic naming~~ -- Unstaged: field renamed to `DataSourceID`, consistent with all other ID fields in the codebase (`TagID`, `FormID`, `TenantID`, etc.).

6. ~~Tenants `main.go` logs graceful shutdown signal at `Error` level~~ -- Unstaged: `app.Logger.Error(fmt.Sprintf("application received shutdown signal: ..."))` changed to `app.Logger.Info`. Expected operational event, not an error condition.

7. ~~`DataSourceRef` / `BindingSource` domain types absent~~ -- New `data_source_ref.go` adds `BindingSource`, `BindingSourceType` constants (`field`, `static`), `NewBindingSourceType` constructor with `isValidBindingSourceType` predicate, and `DataSourceRef` for parameterizing data source lookups from field values or static values.

8. ~~`data-lake` data source type absent~~ -- `DataLakeDataSourceAttributes` added with `Query`, `RequiredKeys`, `OptionalKeys`, `Catalog`, `Schema`, `ValueField`, `LabelField`, and `TimeoutMs`. `DataLakeLookupStrategy` added, wired through the strategies bootstrap. `BigQueryDataLakeClient` (stub, see issue #10 below) added and registered. Type fully registered in `isValidSourceType` and `isValidAttributeType`.

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

`RuleExpression.FieldKey` has no referential integrity check against version fields -- Rules may be created when a field does not yet have an ID and cannot be associated at creation time. Expression field keys are resolved at evaluation time against the `RuleEvaluationContext`; invalid keys evaluate to nil/zero in the `expr` environment, which is acceptable behavior for conditional visibility rules. _(Carried from 5/20.)_

**`stratreg.ErrStrategyNotFound` produces `SubmissionStatusFailed` rather than `SubmissionStatusRejected`** -- `Failed` is the correct semantic: the submission itself is valid, the system was misconfigured. `Rejected` would misrepresent a valid submission as invalid. _(Closed as Will Not Fix in 5/29.)_

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **No `SubmissionAttempt` ever created -- `Reject`/`Fail` error params unused** (`submission_jobs_service.go:216-245`) -- `recordAttempt` transitions status and persists the submission but never calls `NewSubmissionAttempt` or appends to `submission.Attempts`. `Submission.Fail(err error)` and `Submission.Reject(err error)` accept an `err` parameter but discard it entirely (`submission.go:116-123`). The domain type, constructor, `Submission.Attempts` field, and mongo document mapper all exist as dead code. With the retry loop in place, operators have no audit trail of per-attempt outcomes or the errors that caused them. _(Carried from 6/5 #1, P2.)_

2. **`Tag` and `TagVersion` have no `validate` struct tags despite calling `validate.ValidateStruct`** (`tag.go:15-22`, `tag_version.go:29-39`) -- `Tag` fields (`TenantID`, `Key`, `DisplayName`) and `TagVersion` fields (`TagID`, `Version`, `Type`) carry no `validate:"required"` tags. Both `NewTag` and `NewTagVersion` call `validate.ValidateStruct` but the call is a no-op on untagged structs. An empty `tenantID`, empty `key`, or empty `tagType` passes construction without error. Compare to `Form`, `FormVersion`, `Submission`, and `DataSource`, which all tag their required fields. _(Carried from 6/5 #2, P2.)_

3. **`TagType` has no constants and no `isValidTagType` predicate** (`tag_version.go:18`) -- `type TagType string` is declared with no constants. `NewTagVersion` accepts any `TagType` string without validation. Every other constrained type in the forms service -- `FormVersionStatus`, `FieldType`, `SubmissionStatus`, `RuleType`, `ExprOperator` -- uses `validate.NewTypeValidator` with an explicit enum. `TagType` is the sole exception. _(Carried from 6/5 #3, P2.)_

4. **`FieldTagMapping` missing `validate` struct tags -- `NewFieldTagMapping` validates nothing** (`field_tag_mapping.go:23-38`) -- `NewFieldTagMapping` calls `validate.ValidateStruct(ftm)` but neither `FieldTagMapping` nor its embedded `FieldTagMappingConfig` carry any `validate` tags. `FieldID` and `TagVersionID` can be empty strings at construction time. Same pattern as issues #2 and #3 above; same fix: add `validate:"required"` to `FieldID` and `TagVersionID`. _(New, P2.)_

#### Missing Functionality

5. **Field validator strategies: select and checkbox remain stubs; date partial** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:37`) -- Both select and checkbox return `nil` without performing any validation. Date has `checkValueRequired` but no date-range validation (TODO comment present). Submissions with these field types pass validation unconditionally. _(Carried from 6/5 #4, P3.)_

6. **`tagsService.Delete` missing active-version guard** (`tags_service.go:138`) -- A `// FIXME` comment acknowledges the missing invariant. Deleting a tag with a historically active version should be prevented to preserve audit history, consistent with the guard on `formsService.Delete` via `hasActiveVersion`. Without it, a tag backing live `FieldTagMapping` records can be hard-deleted. _(Carried from 6/5 #5, P3.)_

---

### Tenants Service

#### Missing Functionality

7. **`Find()` has no pagination or filtering** (`tenants_service.go:30-40`) -- _(Carried from 6/5 #6, P3.)_

8. **`Lookup` value object has no validation** (`lookup.go`) -- _(Carried from 6/5 #7, P3.)_

---

### Cross-Service / pkg/

#### Bugs

9. **`defer close(pool)` race on shutdown** (`pkg/worker/background_worker.go:188-194`) -- `onLeader` creates the pool, spawns worker goroutines, and defers `close(pool)`. On context cancellation, `close(pool)` fires before the worker goroutines have observed the cancellation. A worker that reaches `w.WorkerPool <- w.JobChannel` after the channel is closed panics. The existing `sync.WaitGroup` in `Start` synchronizes `onLeader` itself but does not track the goroutines spawned inside it. Suggest dropping `defer close(pool)` (let GC reclaim) or adding a second `sync.WaitGroup` inside `onLeader` that waits on worker exits before closing. _(Carried from 6/5 #8, P2.)_

10. **`BigQueryDataLakeClient` is a live-registered stub** (`adapters/clients/big_query_data_lake_client.go:24`) -- `DataSourceTypeDataLake` is fully accepted by the domain, persistence, and REST API layers. Any `data-lake` data source created today will be persisted and processed by the worker, but `Query` unconditionally returns `ErrBigQueryDataLakeNotConfigured`. There is no indication at API ingestion time that the type is unimplemented. A `data-lake` data source will silently cycle through worker retries and exhaust the retry limit with no meaningful error message. Options: guard the type at the command/API layer until BigQuery is implemented, or surface `ErrBigQueryDataLakeNotConfigured` as a non-retryable error to avoid retry churn. _(New, P2.)_

#### Architectural

11. **Test coverage gaps (forms, `pkg/`)** -- No domain/service/strategy/repository tests in the forms service; zero test files in `pkg/`. `PingFedTokenValidator` handles security-critical JWT validation with no test coverage. Tenants service remains well-covered. _(Carried from 6/5 #9, P3.)_

12. **No domain events** for cross-service communication. _(Carried from 6/5 #10, P3.)_

---

## Priority Summary

| Priority | #   | Issue                                                                               | Service(s)  |
| -------- | --- | ----------------------------------------------------------------------------------- | ----------- |
| **P2**   | 1   | No `SubmissionAttempt` created; error params unused                                 | Forms       |
| **P2**   | 2   | `Tag`/`TagVersion` missing `validate` struct tags -- `ValidateStruct` is no-op      | Forms       |
| **P2**   | 3   | `TagType` has no constants and no `isValidTagType` predicate                        | Forms       |
| **P2**   | 4   | `FieldTagMapping` missing `validate` struct tags -- `ValidateStruct` is no-op       | Forms       |
| **P2**   | 9   | `defer close(pool)` shutdown race                                                   | pkg/worker  |
| **P2**   | 10  | `BigQueryDataLakeClient` stub accepts API requests, always fails at processing time | Tenants     |
| **P3**   | 5   | Select/checkbox validator stubs, date partial                                       | Forms       |
| **P3**   | 6   | `tagsService.Delete` missing active-version guard                                   | Forms       |
| **P3**   | 7   | Tenants `Find()` no pagination                                                      | Tenants     |
| **P3**   | 8   | `Lookup` no validation                                                              | Tenants     |
| **P3**   | 11  | Test coverage gaps (forms, pkg/)                                                    | Forms, pkg/ |
| **P3**   | 12  | No domain events                                                                    | All         |

---

## Production Readiness

| Service     | Rating                         | Assessment                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| ----------- | ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Forms**   | **8.5/10 -- Production-Ready** | Authentication is now real (PingFed JWT with JWK rotation). Form definition mapper reworked: existing IDs are preserved on update and key uniqueness/expression-reference violations are caught before persistence. Commands migrated to sub-packages. 4 active P2 gaps: `Tag`, `TagVersion`, and `FieldTagMapping` construction validation are all no-ops (missing `validate` struct tags); `TagType` accepts any string; `SubmissionAttempt` audit trail remains absent despite the retry loop making it more necessary than ever. |
| **Tenants** | **9/10 -- Production-Ready**   | Authentication now real. `DataLakeLookupStrategy` and `BigQueryDataLakeClient` added; however the client is a stub that always returns an error at processing time -- a visible production gap if any `data-lake` data sources are configured and processed by the worker. No other regressions.                                                                                                                                                                                                                                     |
| **pkg/**    | **8/10 -- Production-Ready**   | Auth package cleaned up: `reflect` removed, typo corrected, `claims != nil` check is now idiomatic. `defer close(pool)` race in the background worker is unchanged. Zero test coverage across `pkg/`.                                                                                                                                                                                                                                                                                                                                |

---

## Summary

### Progress Since 6/5

- **PingFed JWT authentication** (committed + unstaged) -- Both services now authenticate via `Authorization: Bearer` with JWK-based JWT validation (audience, issuer, expiry, `WithIssuedAt`, 5s leeway). `PlaceholderTokenValidator` remains available for local dev. `autenticatorSettings` typo corrected in unstaged. `reflect`-based claims check replaced with `claims != nil` in unstaged.

- **Commands migrated to sub-packages** (committed) -- `ports/commands.go` in both services replaced with per-entity files under `commands/`. All service methods, handlers, mocks, and tests updated. Reduces per-file responsibility and improves navigability at the cost of one additional import path per consumer.

- **`form_definition_mapper` reworked** (committed + unstaged) -- `updateFormDefinition` now looks up existing entities by ID rather than always creating new ones, preserving identity on update. `formKeys` and `expressionKeys` maps enforce key uniqueness and expression referential integrity across pages, sections, and fields before any persistence write. Unstaged corrects an error message inside `updatePage` that was referencing the wrong entity type and wrong ID field.

- **`DataSourceRef` / `BindingSource` / `DataLake`** (committed) -- Structured parameterized lookup binding added to the forms domain. `DataLakeDataSourceAttributes` and `DataLakeLookupStrategy` extend the tenants service; `BigQueryDataLakeClient` stubs the actual BigQuery integration pending credentials/configuration. `DataSourceID` naming corrected in unstaged.

- **Tenants `main.go` shutdown log level** (unstaged) -- Shutdown signal demoted from `Error` to `Info`.

### Current State

**12 remaining issues** (6 P2, 6 P3). 0 P0, 0 P1.

The authentication work is the most significant production-readiness improvement since 6/5, closing the last P3 placeholder gap. The mapper rework is a correctness improvement: form definitions with duplicate keys or invalid expression references are now rejected deterministically rather than silently persisted in a broken state.

The new `FieldTagMapping` P2 issue (#4) follows the exact same pattern as the outstanding `Tag`/`TagVersion` issues (#2, #3) -- a `validate.ValidateStruct` call with no struct tags to validate against. All three are one-line fixes per field. The `BigQueryDataLakeClient` stub (#10) is the only new issue that introduces visible production risk: the `data-lake` type is now live in the API and domain but always fails at the worker layer with no retry-bypass and no API-layer guard.

**Hexagonal Architecture** -- The `formDefinitionMapper` is correctly placed in the services layer as a domain mapping concern driven from the application service, not leaked into the adapter. The `DataLakeLookupStrategy` follows the established strategy pattern: a port (`DataLakeClient`) defined in the core, an adapter (`BigQueryDataLakeClient`) in the clients package, wired at the composition root. Commands are now cleanly organized in sub-packages that match the aggregate they operate on, making the port boundary more navigable.

**DDD** -- `BindingSource` and `DataSourceRef` are correctly modeled as value objects in the domain with a constrained type enum (`isValidBindingSourceType`). The mapper's `ErrDuplicateFormKey` and `ErrInvalidExpressionKey` sentinel errors correctly surface form definition invariants that were previously enforced nowhere. The ongoing DDD gap is the cluster of validation no-ops: `NewTag`, `NewTagVersion`, and `NewFieldTagMapping` all call `validate.ValidateStruct` without struct tags, meaning the invariants they are intended to enforce are entirely absent. `SubmissionAttempt` remains modeled but unused.

**Idiomatic Go** -- `NewBearerAuthenticator[T auth.Claims]` is a well-typed generic wrapper that avoids a type assertion at the call site. Replacing `reflect.ValueOf(claims).IsValid() && !valOfClaims.IsZero()` with `claims != nil` is the correct idiomatic form for an interface nil check. The `autenticatorSettings` typo and `DataSourceRef.DataSourceId` naming inconsistency are both corrected in unstaged changes.

### Highest-Impact Improvements

1. **Add `validate` struct tags to `Tag`, `TagVersion`, and `FieldTagMapping`** (P2) -- one-line fix per required field; immediately enforces construction-level invariants on the live tag and field-mapping API paths.
2. **Add `TagType` constants and `isValidTagType` predicate** (P2) -- define the enum, wire `validate.NewTypeValidator` into `NewTagVersion`; closes the unconstrained type string gap.
3. **Wire `SubmissionAttempt` into `recordAttempt()`** (P2) -- append `NewSubmissionAttempt(attempt, result, errorDetails)` before persisting; store the `err` that `Reject`/`Fail` currently discard.
4. **Fix the worker pool shutdown race** (P2) -- drop `defer close(pool)` or track worker goroutines in a second `sync.WaitGroup` that `onLeader` waits on before closing.
5. **Guard `DataSourceTypeDataLake` at the API or make `ErrBigQueryDataLakeNotConfigured` non-retryable** (P2) -- prevents `data-lake` data sources from silently consuming all worker retries before failing.
6. **Add active-version guard to `tagsService.Delete`** (P3) -- mirrors the existing `hasActiveVersion` check on `formsService.Delete`; FIXME comment already identifies the invariant.
7. **Implement select/checkbox/date field validators** (P3) -- stubs silently accept invalid field data.
8. **Backfill tests for forms domain/service/strategies and `pkg/`** (P3) -- `PingFedTokenValidator` is the most urgent gap given its security-critical role.
