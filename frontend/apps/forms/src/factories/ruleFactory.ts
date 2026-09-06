import type { IRule, RuleType } from "@/types/rule";

export function createRule(id: string, type: RuleType): IRule {
  return {
    id,
    type,
    expressions: [],
  };
}
