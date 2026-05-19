# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/18 Review

1. ~~Empty `fmt.Errorf("")` across all validators and `checkValueRequired`~~ (Forms, P2) -- `checkValueRequired` now returns sentinel errors (`ErrFieldRequired`, `ErrFieldTypeValue`). All validators use `newValidationErr(field.Key, err)` which wraps errors as `ErrFieldValidation: '<key>' failed on <cause>`. Text validator: "min length", "max length", "does not match pattern %s". Number validator: "min value", "max value". `ErrFieldRequired` and `ErrFieldTypeValue` defined in `strategies/utils.go`. *(Committed. Resolves 5/18 #1.)*

2. ~~`ExprRuleEvaluator` is a stub — always returns `true`~~ (Forms, P2) -- The evaluator now builds expression statements from rule expressions using the `expr-lang/expr` library. Compiles and runs boolean expressions against the `RuleEvaluationContext`. Operator registry maps domain operators (`equal`, `nequal`, `lt`, `gt`, `lte`, `gte`) to expression syntax (`==`, `!=`, `<`, `>`, `<=`, `>=`). Join operators (`and`/`or`) mapped to `&&`/`||`. Statement construction iterates expressions in position order. *(Committed. Resolves 5/18 #3.)*

3. ~~`ExpressionOperator` has no validation~~ (Forms, P3) -- Renamed to `ExprOperator`. `NewRuleExpression` now validates the operator against `isValidExprOperator` with six defined operators: `equal`, `nequal`, `lt`, `gt`, `lte`, `gte`. `ErrInvalidExprOperator` sentinel error added. *(Committed. Resolves 5/18 #6.)*

4. ~~`Process` never updates submission status after validation~~ (Forms, P2) -- `Process` now begins a transaction after validation, upserts the submission, and commits. Transactional persistence of submission state is in place. However, submission status is still not explicitly transitioned to `accepted`/`rejected` and no `SubmissionAttempt` is created. *(Partially resolved. Committed. See remaining #1.)*

5. ~~Packages used non-standard import paths~~ (All) -- All import paths renamed from `github.com/cmclaughlin24/sundance/backend/...` to `sundance/backend/...`. *(Committed.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

`RuleExpression.FieldKey` has no referential integrity check against version fields -- Rules may be created when a field does not yet have an ID and therefore cannot be associated at creation time. Expression field keys are resolved at evaluation time against the `RuleEvaluationContext`; invalid keys evaluate to nil/zero in the `expr` environment, which is acceptable behavior for conditional visibility rules.

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **`Process` does not update submission status to accepted/rejected** (`submission_jobs_service.go:85-97`) -- Validation runs, a transaction is opened, the submission is upserted, and the transaction commits. However, `submission.Status` is never changed from `pending` — there is no call to set status to `accepted` on success or `rejected` on failure. No `SubmissionAttempt` is created to record the processing outcome. The FIXME at line 79 acknowledges that draft-version submissions should be rejected. Submissions will be re-processed on every worker tick since they remain `pending`. *(P2 — submission lifecycle incomplete.)*

2. **`ReplaySubmissionCommand` has no `validate` tags** (`ports/commands.go:179-182`) -- `TenantID` and `ID` fields have no `validate:"required"` tags. `validate.ValidateStruct(command)` always passes even with empty values. *(Carried from 5/15 #1, P2.)*

3. **`joinOperator` default case returns `fmt.Errorf("")`** (`evaluators/expr_rule_evaluator.go:79`) -- The default case in the switch returns an empty error instead of using `domain.ErrInvalidJoinOperator`. This path is unreachable since `NewRuleExpression` already validates join operators, but the empty error is non-idiomatic. *(P3.)*

#### Missing Functionality

4. **Field validator strategies: select and checkbox remain stubs** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`) -- Both return `nil` without performing any validation. Date has `checkValueRequired` but no date-specific validation (TODO at `date_field_validator.go:37`). *(P3.)*

---

### Tenants Service

#### Missing Functionality

5. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

6. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Cross-Service

#### Architectural

7. **Test coverage gaps** -- No domain-layer, service-layer, or repository-layer tests exist in the forms service. Zero test files in entire `pkg/` directory. Tenants service has service-layer tests but delete and lookup test cases are empty (TODO). *(Carried from 5/10 #10, P3.)*

8. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

9. **No real authentication** -- Placeholder only. All services have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | `Process` does not update submission status | Forms |
| **P2** | 2 | `ReplaySubmissionCommand` no validation tags | Forms |
| **P3** | 3 | `joinOperator` default returns empty error | Forms |
| **P3** | 4 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 5 | `Find()` no pagination | Tenants |
| **P3** | 6 | `Lookup` no validation | Tenants |
| **P3** | 7 | Test coverage gaps | All |
| **P3** | 8 | No domain events | All |
| **P3** | 9 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **8/10 -- Production-Ready** | Submission processing pipeline is now functional end-to-end: `Process` fetches the submission and version, evaluates visibility rules via `ExprRuleEvaluator` (using `expr-lang/expr` to compile and run boolean expressions), validates visible fields via strategy-pattern validators, and persists the submission in a transaction. Text and number validators produce descriptive errors via `newValidationErr`. Rule expressions support six comparison operators with and/or join semantics. `Version.FlatFields()` provides a convenience accessor for building evaluation context. Remaining P2 gaps: submission status never transitions to accepted/rejected, no `SubmissionAttempt` created, `ReplaySubmissionCommand` missing validation tags. Core form/version CRUD with lifecycle management remains solid. Handler tests exist but no service/domain/repository tests. |
| **Tenants** | **8/10 -- Production-Ready** | No changes since 5/15 (package rename only). Fully functional including background job processing pipeline, leader election, and data source strategies. Service-layer tests exist (delete and lookup cases are TODO). Only P3 gaps remain (pagination and `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | No changes since 5/15 (package rename only). All previously identified bugs and architectural issues resolved. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/18

- **Rule evaluator fully implemented** (committed) -- `ExprRuleEvaluator` uses the `expr-lang/expr` library to compile and execute boolean expressions built from rule expressions. The `statement` method iterates rule expressions in position order, resolves each operator via an `exprRegistry` (a `stratreg.StrategyRegistry` mapping `ExprOperator` to statement builder functions), applies join operators (`&&`/`||`), and concatenates into a single expression string. `expr.Compile` with `expr.AsBool()` ensures type safety. `expr.Run` evaluates against the `RuleEvaluationContext` (submission field values keyed by field key). Sentinel errors `ErrInvalidExpression` and `ErrInvalidExpressionOutput` added.

- **Expression operators defined and validated** (committed) -- `ExpressionOperator` renamed to `ExprOperator`. Six operators defined: `equal`, `nequal`, `lt`, `gt`, `lte`, `gte`. `NewRuleExpression` validates against `isValidExprOperator`. `ErrInvalidExprOperator` sentinel error added. The operator registry in the evaluator maps each to its expression syntax (`==`, `!=`, `<`, `>`, `<=`, `>=`).

- **Validation error messages made descriptive** (committed) -- All `fmt.Errorf("")` calls replaced with meaningful errors. New sentinel errors in `strategies/utils.go`: `ErrFieldRequired`, `ErrFieldValidation`, `ErrFieldTypeValue`. `newValidationErr(field, err)` produces wrapped errors like `field validation: 'fieldKey' failed on min length`. `checkValueRequired[T]` returns `ErrFieldRequired` for nil/empty required values and `ErrFieldTypeValue` for type mismatches.

- **`Process` persists submission in a transaction** (committed) -- After validation, `Process` begins a transaction via `s.database.BeginTx`, upserts the submission via `s.submissionRepository.Upsert`, and commits via `s.database.CommitTx`. Deferred `RollbackTx` ensures cleanup on error. The `submissionJobsService` now holds a `database.Database` reference. Status transition and attempt creation remain TODO.

- **`Version.FlatFields()` convenience method** (committed) -- Returns all fields across all pages and sections as a flat slice. Used by `validate` to build the `RuleEvaluationContext` by iterating all version fields and looking up submitted values by field ID, keying by field key.

- **`RuleExpression.FieldID` renamed to `FieldKey`** (committed) -- Rule expressions now reference fields by their human-readable key (e.g., `"firstName"`) rather than by UUID. This aligns with the `expr` evaluation context, which is keyed by field key. DTOs updated to use `fieldKey` in JSON.

- **Package import paths renamed** (committed) -- All imports changed from `github.com/cmclaughlin24/sundance/backend/...` to `sundance/backend/...` across both services and `pkg/`.

### Current State

**9 remaining issues** (5 carried from 5/18; 4 resolved, 0 newly introduced that weren't already tracked). 0 P0, 0 P1, 2 P2, 7 P3.

**Forms Service** at 8/10. The submission processing pipeline is now functional: rule evaluation uses `expr-lang/expr` to compile and execute boolean visibility expressions, field validators produce descriptive wrapped errors, and processed submissions are persisted transactionally. The remaining P2 gap is that `Process` persists the submission but never transitions its status — submissions remain `pending` after processing and will be re-fetched by the worker on the next tick. No `SubmissionAttempt` is created to record the outcome. `ReplaySubmissionCommand` still lacks validation tags.

**Tenants Service** at 8/10. Package rename only. No functional changes.

**pkg/** at 8/10. Package rename only. No functional changes.

**Hexagonal Architecture** -- The `ExprRuleEvaluator` correctly lives in `adapters/evaluators/`, implementing the `ports.RuleEvaluator` interface. It imports only domain types and port types from core, plus the `expr-lang/expr` third-party library and `stratreg` from `pkg/common`. The dependency direction (adapter → core) is maintained. The operator registry (`exprRegistry`) is an adapter-layer concern — it maps domain operators to the `expr` library's syntax, keeping the domain free of expression-engine specifics.

**DDD** -- Rule expressions now form a complete conditional logic model: expressions reference fields by key, use validated comparison operators, and join with and/or semantics. The `ExprOperator` type with validation enforces which operators exist in the domain. The `FieldKey` rename correctly reflects that expressions operate on field keys (semantic names) rather than field IDs (opaque UUIDs). The `FlatFields()` method on `Version` provides aggregate-level traversal without exposing internal structure. The gap remains that `Process` doesn't complete the submission lifecycle — no status transition, no attempt recording.

**Idiomatic Go** -- `newValidationErr` uses `%w` verb twice for multi-error wrapping (`fmt.Errorf("%w: '%s' failed on %w", ...)`), enabling `errors.Is` checks against both `ErrFieldValidation` and the underlying cause. Sentinel errors follow Go naming conventions. The `exprRegistry` uses method chaining (`New().Set().Set()...`) for clean initialization. The `statementFn` type alias keeps the registry declaration readable.

### Highest-Impact Improvements

1. **Complete `Process` submission lifecycle** (P2 — transition status to accepted/rejected, create `SubmissionAttempt`, handle validation failure path)
2. **Add `validate` tags to `ReplaySubmissionCommand`** (P2 — validation no-op)
3. **Implement select/checkbox/date field validators** (P3 — select and checkbox are stubs, date is partial)
4. **Add test coverage** (P3 — zero tests in `pkg/`, no service/domain tests in forms)
5. **Replace `fmt.Errorf("")` in `joinOperator` default case** (P3 — use `ErrInvalidJoinOperator`)
