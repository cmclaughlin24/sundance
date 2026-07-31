import { useContext } from "react";
import { FormDefinitionContext } from "./FormDefinitionProvider";
import { useStore } from "zustand";

export function useFormDefinitionContext() {
  const store = useContext(FormDefinitionContext);

  if (!store) {
    throw new Error(
      "useFormDefinitionStore be used within FormDefinitionContext",
    );
  }

  return store;
}

export function useForm() {
  const store = useFormDefinitionContext();
  return useStore(store, (s) => s.form);
}

export function useFormVersion() {
  const store = useFormDefinitionContext();
  return useStore(store, (s) => s.version);
}

export function useTenantId() {
  const store = useFormDefinitionContext();
  return useStore(store, (s) => s.form?.tenantId);
}
