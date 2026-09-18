import {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  type IRule,
  type IRuleExpression,
  type IRuleState,
  type RuleExprSourceType,
} from "@/types/rule";
import { sortPositioned } from "./position";
import type { IPage } from "@/types/page";
import type { FormValues } from "@/store/submission/submissionStore";

type EvaluatorFn = (fieldValue: any, target: any) => boolean;

const evaluatorRegistry = new Map<RuleExpressionOp, EvaluatorFn>([
  [RuleExpressionOp.Equal, (a, b) => a === b],
  [RuleExpressionOp.NEqual, (a, b) => a !== b],
  [RuleExpressionOp.LessThan, (a, b) => a < b],
  [RuleExpressionOp.GreaterThan, (a, b) => a > b],
  [RuleExpressionOp.LessThanEqualTo, (a, b) => a <= b],
  [RuleExpressionOp.GreaterThanEqualTo, (a, b) => a >= b],
]);

export type EvalNamespace = Record<string, any>;

export type EvalContext = Partial<Record<RuleExprSourceType, EvalNamespace>>;

export function buildFieldEvalNamespace(
  pages: IPage[] | null,
  values: FormValues,
): EvalNamespace {
  const namespace: EvalNamespace = {};

  if (!pages || pages.length === 0) {
    return namespace;
  }

  pageLoop: for (const page of pages) {
    if (!page.sections) {
      continue pageLoop;
    }

    sectionLoop: for (const section of page.sections) {
      if (!section.elements) {
        continue sectionLoop;
      }

      for (const element of section.elements) {
        namespace[element.key] = values[element.id];
      }
    }
  }

  return namespace;
}

export function evaluateRules(
  rules: IRule[],
  evalCtx: EvalContext,
  defaultState?: Partial<IRuleState>,
): Readonly<IRuleState> {
  let state: IRuleState = {
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
  let ctx = evalCtx ?? {};

  for (let i = 0; i < expressions.length; i++) {
    const exp = expressions[i];
    const exprResult = evaluateExpression(exp, ctx);

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
  evalCtx: EvalContext,
): boolean {
  const evaluator = evaluatorRegistry.get(exp.operator);

  if (!evaluator) {
    throw new Error(`invalid expression operator: ${exp.operator}`);
  }

  const namespace: EvalNamespace = evalCtx[exp.source.type] ?? {};
  const value = namespace[exp.source.key];

  return evaluator(value, exp.value);
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
