import { useContext } from "react";
import { FormDispatchContext, FormStateContext } from "./formContext";
import { useEvalContext } from "./evalContext";
import type { IElement } from "@/types/element";
import type { IRuleState } from "@/types/rule";
import { evaluateRules } from "@/utils/evaluate";

export function useFormState() {
  return useContext(FormStateContext);
}

export function useFormDispatch() {
  return useContext(FormDispatchContext);
}

export function useElementRuleState(element: IElement): Readonly<IRuleState> {
  const evalCtx = useEvalContext();

  return evaluateRules(element.rules, evalCtx, {
    readonly: element.attributes.isReadOnly,
    required: element.attributes.isRequired,
  });
}
