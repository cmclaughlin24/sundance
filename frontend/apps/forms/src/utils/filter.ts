import type { HasRules } from "@/types/rule";
import { evaluateRules, type EvalContext } from "./evaluate";

export function filterVisible<T extends HasRules>(
  hasRule: T[],
  evalCtx: EvalContext,
): T[] {
  return hasRule.filter((hr) => {
    const { visible } = evaluateRules(hr.rules, evalCtx);
    return visible;
  });
}
