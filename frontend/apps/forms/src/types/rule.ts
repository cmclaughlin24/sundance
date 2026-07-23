import type { HasPosition } from "./hasPosition";

export type RuleType = "visible" | "required" | "readonly";

export interface IRule {
  id: string;
  type: RuleType;
  expressions: IRuleExpression[];
}

export enum RuleExpressionOp {
  Equal = "equal",
  NEqual = "nequal",
  LessThan = "lt",
  GreaterThan = "gt",
  LessThanEqualTo = "lte",
  GreaterThanEqualTo = "gte",
}

export enum RuleExpressionJoinOp {
  And = "and",
  Or = "or",
}

export interface IRuleExpression extends HasPosition {
  fieldKey: string;
  operator: RuleExpressionOp;
  value: any;
  joinWithPrevious: RuleExpressionJoinOp;
}

export interface HasRules {
  rules: IRule[];
}

export interface IRuleStates {
  required: boolean;
  readonly: boolean;
  visible: boolean;
}
