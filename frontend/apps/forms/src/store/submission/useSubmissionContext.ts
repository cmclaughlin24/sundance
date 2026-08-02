import { useContext } from "react";
import { useEvalContext } from "./evalContext";
import type { IElement } from "@/types/element";
import type { IRuleState } from "@/types/rule";
import { evaluateRules } from "@/utils/evaluate";
import { SubmissionStoreContext } from "./SubmissionProvider";
import { useStore } from "zustand";
import { useShallow } from "zustand/react/shallow";

export function useSubmissionContext() {
  const store = useContext(SubmissionStoreContext);

  if (!store) {
    throw new Error("useFormStore be used within FormStoreContext");
  }

  return store;
}

export function useSubmissionDispatch() {
  const store = useSubmissionContext();
  return useStore(
    store,
    useShallow((s) => ({
      setValue: s.setValue,
      setError: s.setError,
    })),
  );
}

export function useFormValues() {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.values);
}

export function useFormErrors() {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.errors);
}

export function useElementValue<T>(elementId: string, defaultValue?: T): T {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.values[elementId]) ?? defaultValue;
}

export function useElementErrors(elementId: string): string[] {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.errors[elementId]);
}

export function useElementRuleState(element: IElement): Readonly<IRuleState> {
  const evalCtx = useEvalContext();

  // TODO: Possible optimization to useMemo to prevent the rule state for each element from being computed each re-render.
  return evaluateRules(element.rules, evalCtx, {
    readonly: element.attributes.isReadOnly,
    required: element.attributes.isRequired,
  });
}
