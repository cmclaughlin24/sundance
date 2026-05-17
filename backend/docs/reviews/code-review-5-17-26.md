# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/15 Review

1. ~~`Process` returns empty error on draft version check~~ (Forms, P3) -- `Process` now returns `ErrVersionStatus` sentinel error instead of `fmt.Errorf("")`. Descriptive error with proper sentinel. *(Committed. Resolves 5/15 #2.)*

2. ~~`SubmissionJobsService.Process` payload validation incomplete~~ (Forms, P3) -- `Process` now iterates version pages → sections → fields, resolves each field's validator strategy via `FieldValidatorRegistry`, retrieves the submitted value via `submission.GetFieldValue`, and calls `fieldValidator.Validate()`. Required-field check added for missing values. Strategy pattern fully wired end-to-end. Rule evaluation (page/section/field conditionals) remains TODO. *(Committed. Partially resolves 5/15 #3.)*

3. ~~`validateField` logs `submission.ID` as `version_id`~~ (Forms, P3) -- Fixed. *(Resolves new issue found this cycle.)*

4. ~~`Process` logs stale `err` variable~~ (Forms, P2) -- Fixed. *(Resolves new issue found this cycle.)*

5. ~~`routes_test.go` won't compile — missing `host` argument~~ (Forms, P2) -- Fixed. *(Resolves new issue found this cycle.)*

6. ~~`SubmissionToResponse` iterates over empty slice — values always empty~~ (Forms, P1) -- `for _, value := range values` was iterating the just-created empty `values` slice instead of `s.Values`. Fixed to iterate `s.Values`. *(Resolves new issue found this cycle.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`TextFieldValidatorStrategy.Validate` never uses submitted value** (`strategies/text_field_validator.go:23-51`) -- `MinLength` check (line 32) returns error unconditionally if the attribute is non-nil, without comparing against the actual value length. `MaxLength` check (line 36) same. `Pattern` match (line 45) tests against empty string `""` instead of the submitted `value`. All error messages are empty strings (`fmt.Errorf("")`). The `value` parameter is received but never referenced. *(P1 — text field validation is non-functional.)*

2. **`Process` never updates submission status after validation** (`submission_jobs_service.go:93-98`) -- Validation errors are collected in a `[]error` slice but: (a) never joined into a single error, (b) never persisted as a `SubmissionAttempt`, (c) submission status is never changed to `accepted` or `rejected`, (d) function returns `nil` regardless of validation outcome. Submissions remain in `pending` status forever after processing. *(P2 — submission lifecycle incomplete.)*

#### Missing Functionality

3. **Field validator strategies are stubs** (`number_field_validator.go:28`, `select_field_validator.go:28`, `checkbox_field_validator.go:28`, `date_field_validator.go:28`) -- All four non-text field validators return `nil` without performing any validation. Each contains a `// TODO: Implement validation.` comment. *(P3.)*

4. **Rule evaluation not implemented** (`submission_jobs_service.go:78,81,84`) -- Page, section, and field rule checks remain TODO. Conditional visibility/required rules are not evaluated during submission processing. Fields behind inactive rules will still be validated. *(P3.)*

5. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go`) -- *(Carried from 5/15 #1, P2.)*

---

### Tenants Service

#### Missing Functionality

6. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

7. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Cross-Service

#### Architectural

8. **Test coverage gaps** -- No domain-layer, service-layer, or repository-layer tests exist in the forms service. Zero test files in entire `pkg/` directory. Tenants service has service-layer tests but delete and lookup test cases are empty (TODO). *(Carried from 5/10 #10, P3.)*

9. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

10. **No real authentication** -- Placeholder only. All services have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P1** | 1 | `TextFieldValidatorStrategy.Validate` broken | Forms |
| **P2** | 2 | `Process` never updates submission status | Forms |
| **P2** | 5 | `ReplaySubmissionCommand` no validation tags | Forms |
| **P3** | 3 | Field validator strategies are stubs (4 of 5 types) | Forms |
| **P3** | 4 | Rule evaluation not implemented | Forms |
| **P3** | 6 | `Find()` no pagination | Tenants |
| **P3** | 7 | `Lookup` no validation | Tenants |
| **P3** | 8 | Test coverage gaps | All |
| **P3** | 9 | No domain events | All |
| **P3** | 10 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **7/10 -- Beta** | Submission processing pipeline is now wired end-to-end: `Process` iterates version structure, resolves field validators via strategy registry, checks required fields. However, a P1 bug undermines it: the text field validator never uses the submitted value (MinLength/MaxLength/Pattern checks are broken). Additionally, `Process` collects validation errors but never updates the submission's status or persists attempt records — submissions stay `pending` forever. Core form/version CRUD with lifecycle management remains solid. Handler tests exist but no service/domain/repository tests. |
| **Tenants** | **8/10 -- Production-Ready** | No changes since 5/15. Fully functional including background job processing pipeline, leader election, and data source strategies. Service-layer tests exist (delete and lookup cases are TODO). Only P3 gaps remain (pagination and `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | No changes since 5/15. All previously identified bugs and architectural issues resolved. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/15

- **Submission field validation wired end-to-end** (committed) -- `SubmissionJobsService.Process` now iterates the version's page → section → field hierarchy, resolves each field's `FieldValidatorStrategy` from the strategy registry, retrieves the submitted value via `submission.GetFieldValue(field.ID)`, and invokes `fieldValidator.Validate()`. Required-field enforcement added: if a field's `Attributes.GetIsRequired()` returns true and no submitted value exists, an error is returned. The `validateField` helper cleanly separates per-field validation from the orchestration loop.

- **`ErrVersionStatus` sentinel error introduced** (committed) -- Replaces the previous `fmt.Errorf("")` on draft version check. `Process` now returns a named error that callers can inspect with `errors.Is`.

- **Strategy pattern for field validation** (committed) -- `FieldValidatorRegistry` (a `stratreg.StrategyRegistry[domain.FieldType, FieldValidatorStrategy]`) wired in `strategies.go`. Five strategies registered: text, number, select, checkbox, date. `TextFieldValidatorStrategy` has partial implementation (attribute extraction works, validation logic is broken). Other four are stubs returning `nil`.

- **`SubmissionToResponse` DTO mapper fixed** -- Iteration bug corrected: `for _, value := range s.Values` instead of `for _, value := range values`. Submission API responses now include field values.

### Current State

**10 remaining issues** (5 carried from 5/15; 5 newly identified, 6 resolved this cycle). 0 P0, 1 P1, 2 P2, 7 P3.

**Forms Service** drops to 7/10. The submission processing pipeline has made significant structural progress — the strategy-based field validation architecture is in place and the orchestration loop correctly traverses the version hierarchy. However, the P1 text field validator bug prevents validation from working: `TextFieldValidatorStrategy` never references the submitted value (MinLength/MaxLength always error when set, Pattern matches against empty string). The `Process` method also does not update submission status or create attempt records after validation, leaving submissions permanently pending.

**Tenants Service** at 8/10. No changes since 5/15.

**pkg/** at 8/10. No changes since 5/15.

**Hexagonal Architecture** -- Both services maintain correct dependency direction. The new `FieldValidatorStrategy` port interface and strategy registry follow the hexagonal pattern: strategies are registered at the composition root and injected into the service via the `Strategies` port struct. The service never imports adapter packages.

**DDD** -- The field validation strategies correctly model domain-specific validation rules per field type. The `GetFieldValue` method on `Submission` provides aggregate-level access to field data. The `SubmissionFieldValue` value object and `FieldAttributes` polymorphic interface maintain type boundaries. The gap is that `Process` doesn't complete the domain lifecycle — status transitions and attempt recording are missing.

**Idiomatic Go** -- Strategy pattern uses the generic `stratreg.StrategyRegistry` from `pkg/common`. `GetFieldAttributes[T]` generic helper provides type-safe attribute extraction. Error sentinel `ErrVersionStatus` follows Go naming conventions. The empty `fmt.Errorf("")` calls in the text validator are non-idiomatic and need descriptive messages.

### Highest-Impact Improvements

1. **Fix `TextFieldValidatorStrategy` to use submitted value** (P1 — compare against `value` parameter, add descriptive error messages)
2. **Complete `Process` submission lifecycle** (P2 — update status to accepted/rejected, create `SubmissionAttempt`, persist via repository)
3. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 — validation no-op)
4. **Implement remaining field validator strategies** (P3 — number, select, checkbox, date)
5. **Add test coverage** (P3 — zero tests in `pkg/`, no service/domain tests in forms)
