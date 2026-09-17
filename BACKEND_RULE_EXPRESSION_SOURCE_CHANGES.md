# Backend Plan: RuleExprSource Domain Refactor

## Objective

Update the `RuleExpression` domain model in the Forms service to support multiple source operands (e.g. form fields, user profile claims) using an explicit `RuleExprSource` struct (`type`, `key`).
_Scope:_ **Backend only**. No backwards compatibility required.

---

## Architectural Strategy: Option B with Typed `EvaluationContext`

To achieve strong domain typing, 1:1 symmetry with `RuleExprSource`, and full extensibility for future sources without mutating `Submission`:

1. **Domain Model (`domain/submission.go`):**
   - Define strongly-typed `EvaluationContext`:
     ```go
     type EvaluationContext map[RuleExprSourceType]map[string]any
     ```
   - Symmetrically matches `RuleExprSource` access: `submission.EvaluationContext[source.Type][source.Key]`.

2. **At Save Time (`POST /submissions`):**
   - Ambient sources from the HTTP request (e.g. JWT claims) are captured into `domain.EvaluationContext`:
     ```go
     evalCtx := domain.EvaluationContext{
         domain.RuleExprSourceTypeUserClaim: extractClaimsMap(auth.GetClaimsFromContext(r.Context())),
     }
     ```
   - Saved on `domain.Submission` and persisted to MongoDB. Keeps HTTP `Create` ultra-fast without loading `FormVersion` on ingestion.

3. **When Job Processes (Worker `submissionResolver.resolve`):**
   - The async worker merges submitted form values under the `field` source type:
     ```go
     evalCtx[domain.RuleExprSourceTypeField] = fieldsMap
     ```

4. **At Evaluation Time (`ExprRuleEvaluator`):**
   - The evaluator adapter adapts domain enum keys (`RuleExprSourceType`) to string keys for `expr-lang/expr`'s reflection runtime.
   - Expressions evaluate via bracket notation: `field["annualIncome"] > 50000 && user_claim["department"] == "Underwriting"`.

---

## 1. Domain Layer (`backend/services/forms/internal/core/domain/`)

### File: `rule_expression.go`

- [x] Define `RuleExprSourceType` type and initial validation (`isValidExprSourceType`).
- [x] Add `RuleExprSourceTypeField` ("field") and `RuleExprSourceTypeUserClaim` ("user_claim") to `RuleExprSourceType` and `isValidExprSourceType`.
- [x] Define `RuleExprSource` struct (`Type RuleExprSourceType`, `Key string`).
- [x] Refactor `RuleExpression` struct: replace `FieldKey string` with `Source RuleExprSource`.
- [x] Update `NewRuleExpression` constructor to accept `RuleExprSource` and validate `isValidExprSourceType`.
- [ ] Add validation in `NewRuleExpression` ensuring `strings.TrimSpace(source.Key) != ""` (`ErrInvalidRuleExprSourceKey`).
- [x] Update `HydrateRuleExpression` constructor to accept `RuleExprSource`.

### File: `submission.go`

- [x] Define `EvaluationContext`:
  ```go
  type EvaluationContext map[RuleExprSourceType]map[string]any
  ```
- [x] Add `evalContext EvaluationContext` field to `Submission` aggregate.
- [x] Update `NewSubmission` and `HydrateSubmission` to accept and initialize `EvaluationContext`.
- [x] Add `GetEvalContext() EvaluationContext` accessor method.

---

## 2. Core Ports & Commands (`backend/services/forms/internal/core/ports/`)

### File: `commands/page_data.go`

- [x] Define `RuleExprSourceData`:
  ```go
  type RuleExprSourceData struct {
      Type string
      Key  string
  }
  ```
- [x] Update `RuleExpressionData`: replace `FieldKey string` with `Source RuleExprSourceData`.

### File: `commands/submission.go`

- [x] Add `EvalContext domain.EvaluationContext` to `CreateSubmissionCommand`.
- [x] Update `NewCreateSubmissionCommand` constructor.
- [x] Add `EvalContext domain.EvaluationContext` to `NormalizeSubmissionCommand` and update `NewNormalizeSubmissionCommand`.

### File: `secondary.go`

- [x] Update `RuleEvaluationContext` contract:
  ```go
  type RuleEvaluationContext = domain.EvaluationContext
  ```

---

## 3. Core Services & Processors (`backend/services/forms/internal/core/`)

### File: `services/form_definition_mapper.go`

- [x] In `createRule`: map `re.Source` command data to `domain.RuleExprSource`.
- [x] In `trackExpressionKeys`: only track keys when `expression.Source.Type == domain.RuleExprSourceTypeField`.
  - _Prevents referential integrity errors for user claims or external sources that do not correspond to form element IDs._

### File: `services/submissions_service.go`

- [x] In `Create()`: pass `cmd.EvalContext` to `domain.NewSubmission(...)`.
- [x] In `Normalize()`: pass `cmd.EvalContext` to `domain.NewSubmission(...)`.

### File: `processors/submission_resolver.go` (Option B Worker Processing)

- [x] In `resolve()`: initialize `evalContext` from `s.GetEvalContext()`.
- [x] Map submitted element values by human-readable key into `evalContext[domain.RuleExprSourceTypeField]`.
- [x] Pass composite `evalContext` into `shouldValidate()` and `isRequired()`.

---

## 4. Evaluators Adapter (`backend/services/forms/internal/adapters/evaluators/`)

### File: `expr_rule_evaluator.go`

- [ ] Update `newDefaultStatementFn` to format statement scope using `string(re.Source.Type)` with bracket notation:
  ```go
  return fmt.Sprintf(`%s[%q] %s %v`, string(re.Source.Type), re.Source.Key, operator, val)
  ```
  _(Bracket notation `%s[%q]` safely supports special characters in claim names, such as colons, dots, and hyphens)._
- [ ] In `Evaluate()`: adapt `domain.EvaluationContext` (keyed by `RuleExprSourceType`) to `map[string]map[string]any` before calling `expr.Run` to satisfy Go reflection string-key requirements:
  ```go
  env := make(map[string]map[string]any, len(evalCtx))
  for srcType, values := range evalCtx {
      env[string(srcType)] = values
  }
  ```

---

## 5. Persistence Adapter (`backend/services/forms/internal/adapters/persistence/mongodb/documents/`)

### File: `documents/rule.go`

- [x] Define `ruleExprSource`:
  ```go
  type ruleExprSource struct {
      Type string `bson:"type"`
      Key  string `bson:"key"`
  }
  ```
- [x] Update `ruleExpressionDocument`: replace `FieldKey string bson:"field_key"` with `Source ruleExprSource bson:"source"`.
- [x] Update `toRuleExpressionDocument` and `fromRuleExpressionDocument`.

### File: `documents/submission.go`

- [x] Add `EvaluationContext domain.EvaluationContext bson:"evaluation_context"` to `SubmissionDocument`.
- [x] In `ToSubmissionDocument`: map `s.GetEvalContext()` to `EvaluationContext`.
- [x] In `fromSubmissionDocument`: pass `doc.EvaluationContext` to `domain.HydrateSubmission(...)`.

---

## 6. REST Adapter (`backend/services/forms/internal/adapters/rest/`)

### File: `dto/rule_expression.go`

- [x] Define DTOs:
  ```go
  type RuleExprSourceRequest struct {
      Type string `json:"type" validate:"required"`
      Key  string `json:"key" validate:"required"`
  }

  type RuleExprSourceResponse struct {
      Type string `json:"type"`
      Key  string `json:"key"`
  }
  ```
- [x] Update `RuleExpressionRequest` and `RuleExpressionResponse` to use `Source`.
- [x] Update `requestsToRuleExpressionData` and `ruleExpressionsToResponse`.

### File: `handlers/submission_handlers.go` (Option B Save Time)

- [x] In `CreateSubmission`: assemble `domain.EvaluationContext` and pass to `commands.NewCreateSubmissionCommand`.
- [x] In `NormalizeSubmission`: assemble `domain.EvaluationContext` and pass to `commands.NewNormalizeSubmissionCommand`.
- [ ] Implement `extractClaimsMap(c auth.Claims) map[string]any` (replace current stub).
- [ ] Update test mocks in `handlers/mocks_test.go` if needed.

### File: Swagger / OpenAPI Docs

- [x] Regenerate swagger specs (`swagger.json`, `swagger.yaml`, `docs.go`).

---

## 7. Verification & Tests

- [ ] Add unit tests in `rule_expression_test.go` covering source type and non-empty key validation.
- [ ] Add unit tests in `expr_rule_evaluator_test.go` verifying compound evaluations across both `field` and `user_claim` sources.
- [ ] Add unit tests in `submission_resolver_test.go` verifying that ambient evaluation context from the submission is merged with field values.
- [ ] Run test suite: `go test -v ./...` in `backend/services/forms`.
