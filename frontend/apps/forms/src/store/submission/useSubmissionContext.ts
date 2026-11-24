import { useContext, useMemo } from "react";
import { useEvalContext } from "./evalContext";
import type { IElement } from "@/types/element";
import type { IRuleState } from "@/types/rule";
import { evaluateRules } from "@/utils/evaluate";
import { SubmissionStoreContext } from "./SubmissionProvider";
import { useStore } from "zustand";
import { useSectionInstance } from "./sectionInstanceContext";
import { findSectionInstance } from "./utils";

export function useSubmissionContext() {
  const store = useContext(SubmissionStoreContext);

  if (!store) {
    throw new Error("useFormStore be used within FormStoreContext");
  }

  return store;
}

export function useFormValues() {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.sections);
}

export function useSubmissionInstanceActions() {
  const store = useSubmissionContext();
  const { addInstance, removeInstance, ensureMinInstances } = store.getState();
  return { addInstance, removeInstance, ensureMinInstances };
}

export function useSubmissionSectionInstances(sectionId: string) {
  const store = useSubmissionContext();
  return useStore(store, (s) => s.sections[sectionId]);
}

export function useElementDispatch() {
  const store = useSubmissionContext();
  const { sectionId, instanceId } = useSectionInstance();

  return useMemo(() => {
    const { setValue, setError } = store.getState();

    return {
      setValue: setValue.bind(null, sectionId, instanceId),
      setError: setError.bind(null, sectionId, instanceId),
    };
  }, [store, sectionId, instanceId]);
}

export function useElementValue<T>(elementId: string, defaultValue?: T): T {
  const store = useSubmissionContext();
  const { sectionId, instanceId } = useSectionInstance();

  return (
    useStore(store, (s) => {
      const instance = findSectionInstance(s.sections, sectionId, instanceId);
      return instance?.getValue(elementId);
    }) ?? defaultValue
  );
}

export function useElementErrors(elementId: string): string[] | undefined {
  const store = useSubmissionContext();
  const { sectionId, instanceId } = useSectionInstance();

  return useStore(store, (s) => {
    const instance = findSectionInstance(s.sections, sectionId, instanceId);
    return instance?.getErrors(elementId);
  });
}

export function useElementRuleState(element: IElement): Readonly<IRuleState> {
  const evalCtx = useEvalContext();

  // TODO: Possible optimization to useMemo to prevent the rule state for each element from being computed each re-render.
  return evaluateRules(element.rules, evalCtx, {
    readonly: element.attributes.isReadOnly,
    required: element.attributes.isRequired,
  });
}
