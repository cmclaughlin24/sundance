<div align="center">

# #007 `expr-lang/expr` for Rule Evaluation

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

Form fields, sections, and pages carry configurable visibility, required, and read-only rules that are evaluated at runtime against submitted field values. These rules are user-defined and must be evaluated dynamically without requiring application code changes. Evaluating arbitrary expressions safely is a known risk — naive approaches using `eval`-style execution or string interpolation expose the system to injection attacks and unpredictable behaviour.

## Decision

`expr-lang/expr` is used to evaluate all runtime rule expressions. Rule expressions are compiled into type-safe boolean programs at evaluation time and executed against a `RuleEvaluationContext` — a `map[string]any` of field key to submitted value. The `ExprRuleEvaluator` adapter in `adapters/evaluators/` implements the `RuleEvaluator` port interface, keeping the expression engine isolated from the domain and application layers.

## Consequences

- Rule expressions are evaluated in a sandboxed environment; arbitrary code execution and side effects are not possible.
- Expressions are compiled before evaluation, surfacing syntax errors early rather than at the point of field traversal.
- The `RuleEvaluator` port interface means the expression engine can be swapped without changes to domain or application logic.
- Rule syntax is constrained to the `expr` language; form designers and integrators must author rules within its supported operators and functions.
