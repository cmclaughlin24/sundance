# Full Codebase Review: Forms and Tenants Services

## Issues Resolved Since 5/19 Review

1. ~~`Process` does not update submission status to accepted/rejected~~ (Forms, P2) -- `Process` now calls `recordAttempt()` which implements tripartite error handling: `submission.Accept()` on success, `submission.Reject(err)` on validation/version-status errors, `submission.Fail(err)` on infrastructure errors. Status is transitioned and the submission is persisted in a transaction. *(Committed. Resolves 5/19 #1.)*

2. ~~`ReplaySubmissionCommand` has no `validate` tags~~ (Forms, P2) -- Both `TenantID` and `ID` now have `validate:"required"` tags. `Validate()` method added for consistency with other command structs. *(Committed. Resolves 5/19 #2.)*

3. ~~`joinOperator` default case returns `fmt.Errorf("")`~~ (Forms, P3) -- Default case now returns `domain.ErrInvalidJoinOperator` sentinel error. *(Committed. Resolves 5/19 #3.)*

4. ~~`Fail()` sets `SubmissionStatusRejected` instead of `SubmissionStatusFailed`~~ (Forms, P2) -- `Fail()` now correctly sets `SubmissionStatusFailed`, enabling retry semantics for infrastructure/transient errors. *(Committed. New issue found and fixed this cycle.)*

5. ~~`ReplaySubmissionCommand` missing `Validate()` method~~ (Forms, P3) -- `Validate()` method added for consistency with other command structs. *(Committed. New issue found and fixed this cycle.)*

---

## Will Not Fix

See [5/10 review](code-review-5-10-26.md) for the full Will Not Fix list.

`RuleExpression.FieldKey` has no referential integrity check against version fields -- Rules may be created when a field does not yet have an ID and therefore cannot be associated at creation time. Expression field keys are resolved at evaluation time against the `RuleEvaluationContext`; invalid keys evaluate to nil/zero in the `expr` environment, which is acceptable behavior for conditional visibility rules. *(Carried from 5/19.)*

---

## Remaining Issues

### Forms Service (includes Submissions)

#### Bugs

1. **No `SubmissionAttempt` ever created** (`submission_jobs_service.go:210-239`) -- Despite the method being named `recordAttempt`, no `SubmissionAttempt` is appended to `submission.Attempts`. The domain model (`submission_attempt.go`), the `Submission.Attempts` field, and the MongoDB document mapper all exist but are dead code. Additionally, `Reject(err error)` and `Fail(err error)` accept an `err` parameter but never store it — once attempts are created, the error details should be recorded in `SubmissionAttempt.ErrorDetails`. *(P2 — audit trail gap, unused error parameters.)*

#### Missing Functionality

2. **Field validator strategies: select and checkbox remain stubs** (`select_field_validator.go:28`, `checkbox_field_validator.go:28`) -- Both return `nil` without performing any validation. Date has `checkValueRequired` but no date-specific validation (TODO at `date_field_validator.go:37`). *(Carried from 5/19 #4, P3.)*

3. **`FindJobs` Limit not honored** (`inmemory/submissions_repository.go:87-106`, `mongodb/submissions_repository.go:71-86`) -- `FindSubmissionsFilter.Limit` field exists but neither the in-memory nor MongoDB `FindJobs` implementations apply it. All matching submissions are returned regardless of limit. Not critical for a background worker but could cause unbounded result sets. *(P3.)*

---

### Tenants Service

#### Missing Functionality

4. **`Find()` has no pagination or filtering** (`tenants_service.go`) -- *(Carried from 5/10 #2, P3.)*

5. **`Lookup` value object has no validation** (`lookup.go`) -- *(Carried from 5/10 #3, P3.)*

---

### Cross-Service

#### Architectural

6. **Test coverage gaps** -- No domain-layer, service-layer, or repository-layer tests exist in the forms service. Zero test files in entire `pkg/` directory. Tenants service has service-layer tests but delete and lookup test cases are empty (TODO). *(Carried from 5/10 #10, P3.)*

7. **No domain events** for cross-service communication. *(Carried from 5/10 #11, P3.)*

8. **No real authentication** -- Placeholder only. All services have auth middleware wired, but it uses `PlaceholderAuthenticator`. *(Carried from 5/10 #12, P3.)*

---

## Priority Summary

| Priority | # | Issue | Service(s) |
|----------|---|-------|------------|
| **P2** | 1 | No `SubmissionAttempt` created; error params unused | Forms |
| **P3** | 2 | Select/checkbox validator stubs, date partial | Forms |
| **P3** | 3 | `FindJobs` Limit not honored | Forms |
| **P3** | 4 | `Find()` no pagination | Tenants |
| **P3** | 5 | `Lookup` no validation | Tenants |
| **P3** | 6 | Test coverage gaps | All |
| **P3** | 7 | No domain events | All |
| **P3** | 8 | No real authentication | All |

---

## Production Readiness

| Service | Rating | Assessment |
|---------|--------|------------|
| **Forms** | **9/10 -- Production-Ready** | Submission processing pipeline is now complete end-to-end: `Process` fetches the submission, guards against non-pending status, evaluates visibility rules via `ExprRuleEvaluator`, dynamically evaluates required rules via `isRequired()`, validates visible fields via strategy-pattern validators, and records the outcome via `recordAttempt()` with tripartite error handling (Accept/Reject/Fail) and transactional persistence. `Fail()` correctly sets `SubmissionStatusFailed` for retry semantics. `validate()` extracted as a clean helper. Rule evaluation now influences both visibility and required-field semantics. Core form/version CRUD with lifecycle management remains solid. Remaining P2 gap: `SubmissionAttempt` records are never created despite the domain model existing. Handler tests exist but no service/domain/repository tests. |
| **Tenants** | **8/10 -- Production-Ready** | No changes since 5/15 (package rename only). Fully functional including background job processing pipeline, leader election, and data source strategies. Service-layer tests exist (delete and lookup cases are TODO). Only P3 gaps remain (pagination and `Lookup` validation). |
| **pkg/** | **8/10 -- Production-Ready** | No changes since 5/15 (package rename only). All previously identified bugs and architectural issues resolved. Only remaining gap: zero test coverage. |

---

## Summary

### Progress Since 5/19

- **Submission lifecycle completed** (committed) -- `Process` now calls `recordAttempt()` after validation. `recordAttempt()` implements tripartite error handling: `submission.Accept()` on success (sets `SubmissionStatusAccepted`), `submission.Reject(err)` on `ErrVersionStatus` or `ErrFieldValidation` (sets `SubmissionStatusRejected`), `submission.Fail(err)` on infrastructure errors (sets `SubmissionStatusFailed`). Opens a transaction, upserts the submission, and commits. Submissions no longer remain `pending` after processing.

- **`Fail()` correctly uses `SubmissionStatusFailed`** (committed) -- `Fail()` now sets the correct status constant, distinguishing infrastructure errors (retryable) from validation failures (terminal). The `SubmissionStatusFailed` constant was already defined but was not being used.

- **Dynamic required-field evaluation** (committed) -- New `isRequired()` method on `submissionJobsService` evaluates `RuleTypeRequired` rules against the submission context and dynamically sets `field.Attributes.SetIsRequired()`. The `FieldAttributes` interface now includes `SetIsRequired(bool)`. This means required-field status can be conditional based on other field values.

- **`validate()` extracted as helper** (committed) -- Validation logic is now a clean, separate method on `submissionJobsService`, improving readability and testability of the `Process` flow.

- **`ReplaySubmissionCommand` validation complete** (committed) -- `validate:"required"` tags added to `TenantID` and `ID` fields. `Validate()` method added for consistency with other command structs.

- **`joinOperator` uses sentinel error** (committed) -- Default case in the `joinOperator` switch now returns `domain.ErrInvalidJoinOperator` instead of `fmt.Errorf("")`.

### Current State

**8 remaining issues** (5 resolved from 5/19; 2 newly introduced; 1 new carry). 0 P0, 0 P1, 1 P2, 7 P3.

**Forms Service** at 9/10 (up from 8/10). The submission processing pipeline is now feature-complete: rule evaluation drives both visibility and required-field semantics, field validation is strategy-based with descriptive errors, and processing outcomes are persisted transactionally with correct status transitions. `Fail()` correctly distinguishes infrastructure errors from validation rejections. The sole remaining P2 is that `SubmissionAttempt` records are never created — the domain model, struct field, and persistence mapper all exist but are unused. `Reject(err)` and `Fail(err)` accept error parameters that are never stored.

**Tenants Service** at 8/10. No changes.

**pkg/** at 8/10. No changes.

**Hexagonal Architecture** -- The `recordAttempt()` method correctly lives in the service layer, orchestrating domain method calls (`Accept`/`Reject`/`Fail`) and infrastructure concerns (transaction management) without leaking adapter details. The `SetIsRequired` addition to the `FieldAttributes` interface maintains the port boundary — adapters implement the interface, the service consumes it. The `isRequired()` method reuses the existing `RuleEvaluator` port for required-rule evaluation, avoiding duplication.

**DDD** -- The submission lifecycle is now modeled through explicit domain methods: `Accept()`, `Reject(err)`, `Fail(err)`, and `Reset()`. Status transitions are owned by the aggregate. The `isRequired()` evaluation correctly treats required rules as a domain concern resolved at processing time, not at definition time. The gap is that `SubmissionAttempt` — a value object designed to record processing outcomes — is defined but never populated, making it dead code in the domain model.

**Idiomatic Go** -- The tripartite error handling in `recordAttempt()` uses `errors.Is` for type-safe error classification. `SetIsRequired` on the interface follows Go's convention of small, focused interface methods. The `validate()` extraction follows Go's preference for small, named functions over large method bodies.

### Highest-Impact Improvements

1. **Create `SubmissionAttempt` records in `recordAttempt()`** (P2 — store attempt number, result, and error details; make `Reject(err)`/`Fail(err)` error parameters meaningful)
2. **Implement select/checkbox/date field validators** (P3 — stubs silently accept invalid data)
3. **Add test coverage** (P3 — zero tests in `pkg/`, no service/domain tests in forms)
4. **Apply `FindJobs` limit in both repository implementations** (P3 — unbounded result sets)
5. **Add pagination to tenants `Find()`** (P3 — unbounded result sets)
