import type { IRule } from "@/types/rule";

export const RuleItem: React.FC<{ rule: IRule }> = function ({ rule }) {
  return rule.id;
};
