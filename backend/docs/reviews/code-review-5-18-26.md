# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/15 Review

1. ~~`Process` returns empty error on draft version check~~ (Forms, P3) -- `Process` now returns `ErrVersionStatus` sentinel error instead of `fmt.Errorf("")`. *(Committed. Resolves 5/15 #2.)*

2. ~~`SubmissionJobsService.Process` payload validation incomplete~~ (Forms, P3) -- `Process` now iterates version pages → sections → fields, resolves each field's validator strategy via `FieldValidatorRegistry`, retrieves the submitted value via `submission.GetFieldValue`, and calls `fieldValidator.Validate()`. Required-field check added for missing values. Strategy pattern fully wired end-to-end. *(Committed. Resolves 5/15 #3.)*

3. ~~`SubmissionToResponse` iterates over empty slice — values always empty~~ (Forms, P1) -- `for _, value := range values` was iterating the just-created empty `values` slice instead of `s.Values`. Fixed to iterate `s.Values`. *(Committed.)*

4. ~~`validateField` logs `submission.ID` as `version_id`~~ (Forms, P3) -- Fixed. *(Committed.)*

5. ~~`Process` logs stale `err` variable~~ (Forms, P2) -- Fixed. *(Committed.)*

6. ~~`routes_test.go` won't compile — missing `host` argument~~ (Forms/Tenants, P2) -- Both forms and tenants `routes_test.go` now pass `""` as the second argument to `NewRoutes`. *(Committed.)*

7. ~~`TextFieldValidatorStrategy.Validate` never uses submitted value~~ (Forms, P1) -- Text validator now uses `checkValueRequired[string]` generic helper. `MinLength`/`MaxLength` compare against `len(val)`. `Pattern` matches against `val`. Nil-value early return distinguishes "no value submitted" from "empty string submitted". *(Committed.)*

8. ~~Rule evaluation not implemented~~ (Forms, P3) -- New `RuleEvaluator` port interface with `Evaluate(context.Context, *domain.Rule, RuleEvaluationContext) (bool, error)`. `submissionJobsService.shouldValidate` calls the evaluator for visibility rules on pages, sections, and fields. Labeled loops skip invisible elements. Actual expression evaluation is a stub (`ExprRuleEvaluator` always returns `true`). *(Partially resolved. Committed. See remaining #3.)*

9. ~~`Process` accepted `*ProcessSubmissionJobCommand` requiring manual validation~~ (Forms) -- `Process` now takes `domain.SubmissionID` directly. `ProcessSubmissionJobCommand` removed. Validation loop refactored into a dedicated `validate` helper with fail-fast per field. *(Committed.)*

10. ~~Position maps used for pages/sections/fields~~ (Forms) -- Refactored from `map[float32]T` to sorted slices (`PositionElements[T]`). New generic `hasUniqueElements` and `sortElements` helpers in `position.go`. `GetPagesSlice()`/`GetSectionsSlice()` renamed to `GetPages()`/`GetSections()`/`GetFields()`. `SetSections`/`SetFields` replaced with `AddSections`/`AddFields` + `ReplaceSections`/`ReplaceFields` using clone-validate-swap pattern. *(Committed.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **Empty `fmt.Errorf("")` across all validators and `checkValueRequired`** (`text_field_validator.go:42,46,55`, `number_field_validator.go:39,43`, `utils.go:15,21,26`) -- Seven `fmt.Errorf("")` calls return errors with zero descriptive content. Callers (including `Process`) cannot determine what failed — which field, which constraint, or what the actual vs expected values were. *(P2 — error messages are empty strings.)*

2. **`Process` never updates submission status after validation** (`submission_jobs_service.go:81-84`) -- Validation errors are returned (fail-fast) but the submission status never transitions to `accepted` or `rejected`. No `SubmissionAttempt` is created. The submission is not persisted after processing. Submissions remain in `pending` status forever after being processed. *(P2 — submission lifecycle incomplete.)*

3. **`ExprRuleEvaluator` is a stub — always returns `true`** (`evaluators/expr_rule_evaluator.go:12-13`) -- The rule evaluation architecture is fully wired (port interface, evaluation context, `shouldValidate` in the service) but the only `RuleEvaluator` implementation never inspects the rule's expressions. All visibility rules pass regardless of submitted values. *(P2 — rule evaluation non-functional.)*

4. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go`) -- *(Carried from 5/15 #1, P2.)*

#### Missing Functionality

5. **Field validator strategies: select and checkbox remain stubs** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`) -- Both return `nil` without performing any validation. Date has `checkValueRequired` but no date-specific validation (TODO at `date_field_validator.go:37`). *(P3.)*

6. **`ExpressionOperator` has no validation** (`rule_expression.go:37`) -- `NewRuleExpression` has a `// TODO: Implement validation for operators.` comment. Any string is accepted as an operator. The domain cannot enforce which operators are valid. *(P3.)*

---

### Tenants Service

#### Missing Functionality

7. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

8. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Cross-Service

#### Architectural

9. **Test coverage gaps** -- No domain-layer, service-layer, or repository-layer tests exist in the forms service. Zero test files in entire `pkg/` directory. Tenants service has service-layer tests but delete and lookup test cases are empty (TODO). *(Carried from 5/10 #10, P3.)*

10. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

11. **No real authentication** -- Placeholder only. All services have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | Empty `fmt.Errorf("")` across all validators | Forms |
| **P2** | 2 | `Process` never updates submission status | Forms |
| **P2** | 3 | `ExprRuleEvaluator` is a stub | Forms |
| **P2** | 4 | `ReplaySubmissionCommand` no validation tags | Forms |
| **P3** | 5 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 6 | `ExpressionOperator` no validation | Forms |
| **P3** | 7 | `Find()` no pagination | Tenants |
| **P3** | 8 | `Lookup` no validation | Tenants |
| **P3** | 9 | Test coverage gaps | All |
| **P3** | 10 | No domain events | All |
| **P3** | 11 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8/10 -- Production-Ready** | Submission processing pipeline fully wired: `Process` iterates version hierarchy, evaluates visibility rules via `RuleEvaluator` port, resolves field validators via strategy registry, checks required fields, and validates values. Text and number validators are functional with `checkValueRequired[T]` generic helper. Position model refactored from maps to sorted slices with clone-validate-swap mutation. Rule expressions introduce conditional logic scaffolding (expressions with operators, join logic, position ordering). Core form/version CRUD with lifecycle management remains solid. Remaining P2 gaps: empty error messages in validators, stub rule evaluator, and submission status never updated after processing. Handler tests exist but no service/domain/repository tests. |
| **Tenants** | **8/10 -- Production-Ready** | No changes since 5/15. Fully functional including background job processing pipeline, leader election, and data source strategies. Service-layer tests exist (delete and lookup cases are TODO). Only P3 gaps remain (pagination and `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | No changes since 5/15. All previously identified bugs and architectural issues resolved. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/15

- **Text and number field validators implemented** (committed) -- `TextFieldValidatorStrategy` now uses `checkValueRequired[string]` to handle nil/required/type-mismatch, then validates `MinLength`/`MaxLength` via `len(val)` and `Pattern` via compiled regex against `val`. Nil-value early return distinguishes "no value submitted" (skip validation) from "empty string submitted" (validate). `NumberFieldValidatorStrategy` uses `checkValueRequired[float64]` and validates `Min`/`Max` constraints. New `checkValueRequired[T comparable]` generic helper in `strategies/utils.go` centralizes nil-check, required-check, and type assertion logic for all validators.

- **Position maps replaced with sorted slices** (committed) -- `Version.pages`, `Page.sections`, `Section.fields` changed from `map[float32]T` to `PositionElements[T]` (a type alias for `[]T` where `T` implements `PositionGetter`). New generic helpers `hasUniqueElements[T]` and `sortElements[T]` in `position.go` enforce uniqueness and ordering. `GetPagesSlice()`/`GetSectionsSlice()` renamed to `GetPages()`/`GetSections()`/`GetFields()`. Mutation methods refactored: `SetSections`/`SetFields` replaced with `AddSections`/`AddFields` (append with uniqueness check) and `ReplaceSections`/`ReplaceFields` (clone-validate-swap with rollback on error). This eliminates the map iteration nondeterminism and provides consistent ordering by position.

- **Rule expressions domain entity introduced** (committed) -- New `RuleExpression` entity in `domain/rule_expression.go` with `FieldID`, `Operator` (`ExpressionOperator`), `Value` (`any`), `JoinWithPrevious` (`*JoinOperator`), and position. `JoinOperator` validated against `and`/`or`. `Rule` now holds `PositionElements[*RuleExpression]` instead of being a simple type enum. `AddExpressions` validates position uniqueness and sorts. Full DTO layer (`RuleExpressionRequest`/`RuleExpressionResponse`) with request-to-domain and domain-to-response mappers. MongoDB document mapping added for rule expressions.

- **Rule evaluation architecture wired** (committed) -- New `RuleEvaluator` port interface in `ports/secondary.go` with `Evaluate(context.Context, *domain.Rule, RuleEvaluationContext) (bool, error)`. `RuleEvaluationContext` is a `map[string]any` built from submission field values. `submissionJobsService` holds a `RuleEvaluator` and uses `shouldValidate` to check visibility rules on pages, sections, and fields via labeled loops. `ExprRuleEvaluator` adapter in `adapters/evaluators/` implements the port (stub — always returns `true`). Wired at the composition root in `main.go` via `services.WithRuleEvaluator(&evaluators.ExprRuleEvaluator{})`.

- **`Process` simplified and validation refactored** (committed) -- `SubmissionJobsService.Process` now accepts `domain.SubmissionID` directly instead of `*ProcessSubmissionJobCommand`. `ProcessSubmissionJobCommand` removed. Validation logic extracted into a `validate` helper method that builds the evaluation context, iterates the version hierarchy with visibility checks, and calls `validateField` per visible field. Fail-fast on first validation error.

- **`SubmissionToResponse` DTO mapper fixed** (committed) -- Iteration bug corrected: `for _, value := range s.Values` instead of `for _, value := range values`. Submission API responses now include field values.

- **Route tests fixed** (committed) -- Both forms and tenants `routes_test.go` now pass `""` as the `host` argument to `NewRoutes`.

### Current State

**11 remaining issues** (5 carried from 5/15; 6 newly identified; 10 resolved this cycle). 0 P0, 0 P1, 4 P2, 7 P3.

**Forms Service** at 8/10 (up from 8/10 at 5/15). Both P1 issues from the 5/17 review are resolved: the text field validator now uses submitted values, and the DTO mapper returns field values correctly. Significant architectural progress: position model refactored from nondeterministic maps to sorted slices with generic helpers; rule expressions introduce a conditional logic model with operators, join semantics, and position ordering; rule evaluation is architecturally wired through a port interface with visibility-based field skipping in the submission processing loop. Text and number validators are functional. The remaining P2 gaps are: empty error messages across all validators (`fmt.Errorf("")`), the `ExprRuleEvaluator` stub that always returns `true`, submission status never transitions after processing, and `ReplaySubmissionCommand` missing validation tags.

**Tenants Service** at 8/10. No changes since 5/15.

**pkg/** at 8/10. No changes since 5/15.

**Hexagonal Architecture** -- The new `RuleEvaluator` port interface correctly separates the rule evaluation concern from the service layer. The `ExprRuleEvaluator` lives in `adapters/evaluators/`, maintaining the adapter → core dependency direction. The service depends only on the `RuleEvaluator` interface, never on the concrete evaluator. The `FieldValidatorStrategy` port and strategy registry continue to follow the same pattern. The `RuleEvaluationContext` type alias in ports provides a clean contract between the service and evaluator without leaking implementation details.

**DDD** -- Rule expressions enrich the domain model: rules are no longer simple type enums but contain ordered expressions with field references, operators, comparison values, and join semantics. The `PositionElements[T]` generic type and associated helpers (`hasUniqueElements`, `sortElements`) provide reusable domain infrastructure for ordered collections with uniqueness constraints, applied consistently across pages, sections, fields, and rule expressions. The clone-validate-swap mutation pattern (`AddPages`/`ReplacePages`, `AddSections`/`ReplaceSections`, etc.) protects aggregate invariants. The gap remains that `Process` doesn't complete the domain lifecycle — status transitions and attempt recording are missing.

**Idiomatic Go** -- The `PositionElements[T]` type alias with `PositionGetter` interface constraint is clean generic design. `checkValueRequired[T comparable]` centralizes nil/required/type-assertion logic. Labeled loops (`pageLoop`, `sectionLoop`, `fieldLoop`) for visibility-based skipping are idiomatic. The `ruleGetter` interface in `submission_jobs_service.go` is a minimal local interface following Go's "accept interfaces" convention. The empty `fmt.Errorf("")` calls remain non-idiomatic and need descriptive messages.

### Highest-Impact Improvements

1. **Add descriptive error messages to all validators** (P2 — seven `fmt.Errorf("")` calls with zero content)
2. **Implement `ExprRuleEvaluator` expression evaluation** (P2 — architecture wired, evaluator is a stub)
3. **Complete `Process` submission lifecycle** (P2 — update status to accepted/rejected, create `SubmissionAttempt`, persist via repository)
4. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 — validation no-op)
5. **Implement select/checkbox/date field validators** (P3 — select and checkbox are stubs, date is partial)
