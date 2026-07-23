import type { FormState } from "@/store/formReducer";
import type { HasRules } from "@/types/rule";
import { buildEvalContext, evaluateRules } from "./evaluate";

export function filterVisible<T extends HasRules>(
  hasRule: T[],
  state: FormState,
): T[] {
  return hasRule.filter((hr) => {
    const evalCtx = buildEvalContext(state);
    const { visible } = evaluateRules(hr.rules, evalCtx);
    return visible;
  });
}
