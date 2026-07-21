import type { HasPosition } from "./hasPosition";

export type RuleType = "visible" | "required" | "readonly";

export interface Rule {
  id: string;
  type: RuleType;
  expressions: RuleExpression[];
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

export interface RuleExpression extends HasPosition {
  fieldKey: string;
  operator: RuleExpressionOp;
  value: any;
  joinWithPrevious: RuleExpressionJoinOp;
}

export interface HasRules {
  rules: Rule[];
}
