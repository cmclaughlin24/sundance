import {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  type IRule,
  type IRuleExpression,
  type IRuleStates,
} from "@/types/rule";
import { sortPositioned } from "./sort";
import type { FormState } from "@/store/formReducer";

type EvaluatorFn = (fieldValue: any, target: any) => boolean;

export type EvalContext = Map<string, any>;

const evaluatorRegistry = new Map<RuleExpressionOp, EvaluatorFn>([
  [RuleExpressionOp.Equal, (a, b) => a === b],
  [RuleExpressionOp.NEqual, (a, b) => a !== b],
  [RuleExpressionOp.LessThan, (a, b) => a < b],
  [RuleExpressionOp.GreaterThan, (a, b) => a > b],
  [RuleExpressionOp.LessThanEqualTo, (a, b) => a <= b],
  [RuleExpressionOp.GreaterThanEqualTo, (a, b) => a >= b],
]);

export function buildEvalContext(state: FormState): EvalContext {
  const evalCtx = new Map<string, any>();

  if (!state.version) {
    return evalCtx;
  }

  for (const page of state.version.pages) {
    for (const section of page.sections) {
      for (const element of section.elements) {
        evalCtx.set(element.key, state.values.get(element.id));
      }
    }
  }

  return evalCtx;
}

export function evaluateRules(
  rules: IRule[],
  evalCtx: EvalContext,
  defaultState?: IRuleStates,
): Readonly<IRuleStates> {
  let state: IRuleStates = {
    readonly: false,
    required: false,
    visible: true,
  };

  if (defaultState) {
    state = { ...state, ...defaultState };
  }

  for (const rule of rules) {
    const result = evaluateRule(rule, evalCtx);

    switch (rule.type) {
      case "visible":
        state.visible = result;
        break;
      case "required":
        state.required = result;
        break;
      case "readonly":
        state.readonly = result;
        break;
    }
  }

  return state;
}

export function evaluateRule(rule: IRule, evalCtx: EvalContext): boolean {
  const expressions = sortPositioned(rule.expressions);
  let result = false;

  for (let i = 0; i < expressions.length; i++) {
    const exp = expressions[i];
    const exprResult = evaluateExpression(exp, evalCtx);

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
