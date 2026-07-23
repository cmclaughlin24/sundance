import {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  type IRule,
  type IRuleExpression,
} from "@/types/rule";
import { sortPositioned } from "./sort";

type EvaluatorFn = (fieldValue: any, target: any) => boolean;

const evaluatorRegistry = new Map<RuleExpressionOp, EvaluatorFn>([
  [RuleExpressionOp.Equal, (a, b) => a === b],
  [RuleExpressionOp.NEqual, (a, b) => a !== b],
  [RuleExpressionOp.LessThan, (a, b) => a < b],
  [RuleExpressionOp.GreaterThan, (a, b) => a > b],
  [RuleExpressionOp.LessThanEqualTo, (a, b) => a <= b],
  [RuleExpressionOp.GreaterThanEqualTo, (a, b) => a >= b],
]);

export function evaluate(rule: IRule, values: Map<string, any>): boolean {
  const expressions = sortPositioned(rule.expressions);
  let result = false;

  for (let i = 0; i < expressions.length; i++) {
    const exp = expressions[i];
    const exprResult = evaluateExpression(exp, values);

    if (i === 0) {
      result = exprResult;
      continue;
    }

    result = applyJoinOp(result, exprResult, exp);
  }

  return result;
}

function evaluateExpression(
  exp: IRuleExpression,
  values: Map<string, any>,
): boolean {
  const evaluator = evaluatorRegistry.get(exp.operator);

  if (!evaluator) {
    throw new Error(`invalid expression operator: ${exp.operator}`);
  }

  const fieldValue = values.get(exp.fieldKey);

  return evaluator(fieldValue, exp.value);
}

function applyJoinOp(
  left: boolean,
  right: boolean,
  exp: IRuleExpression,
): boolean {
  switch (exp.joinWithPrevious) {
    case RuleExpressionJoinOp.And:
      return left && right;
    case RuleExpressionJoinOp.Or:
      return left || right;
    default:
      throw new Error(`invalid join operator: ${exp.joinWithPrevious}`);
  }
}
