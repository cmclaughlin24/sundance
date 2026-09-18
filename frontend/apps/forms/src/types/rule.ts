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

export type RuleExprSourceType = "field" | "userClaim";

export interface IRuleExprSource {
  type: RuleExprSourceType;
  key: string;
}

export interface IRuleExpression extends HasPosition {
  source: IRuleExprSource;
  operator: RuleExpressionOp;
  value: any;
  joinWithPrevious: RuleExpressionJoinOp;
}

export interface HasRules {
  rules: IRule[];
}

export interface IRuleState {
  required: boolean;
  readonly: boolean;
  visible: boolean;
}

export type RuleParentType = "element" | "section" | "page";

export interface IFlatRule extends IRule {
  parentId?: string;
  parentType?: RuleParentType;
}
