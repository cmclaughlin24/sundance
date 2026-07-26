import { useContext } from "react";
import { useEvalContext } from "./evalContext";
import type { IElement } from "@/types/element";
import type { IRuleState } from "@/types/rule";
import { evaluateRules } from "@/utils/evaluate";
import { FormStoreContext } from "./FormStoreProvider";
import { useStore } from "zustand";
import { useShallow } from "zustand/react/shallow";

export function useFormStoreContext() {
  const store = useContext(FormStoreContext);

  if (!store) {
    throw new Error("useFormStore be used within FormStoreContext");
  }

  return store;
}

export function useFormDispatch() {
  const store = useFormStoreContext();
  return useStore(
    store,
    useShallow((s) => ({
      setValue: s.setValue,
      setError: s.setError,
    })),
  );
}

export function useForm() {
  const store = useFormStoreContext();
  return useStore(store, (s) => s.form);
}

export function useFormVersion() {
  const store = useFormStoreContext();
  return useStore(store, (s) => s.version);
}

export function useFormValues() {
  const store = useFormStoreContext();
  return useStore(store, (s) => s.values);
}

export function useElementValue<T>(elementId: string, defaultValue?: T): T {
  const store = useFormStoreContext();
  return useStore(store, (s) => s.values[elementId]) ?? defaultValue;
}

export function useElementErrors(elementId: string): string[] {
  const store = useFormStoreContext();
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

export function useTenantId() {
  const store = useFormStoreContext();
  return useStore(store, (s) => s.form?.tenantId);
}
