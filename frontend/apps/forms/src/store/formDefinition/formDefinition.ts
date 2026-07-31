import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import { createStore } from "zustand";

export interface IFormDefinitionStore {
  form: IForm | null;
  version: IFormVersion | null;
}

export type FormDefinitionStoreApi = ReturnType<
  typeof createFormDefinitionStore
>;

export function createFormDefinitionStore(form: IForm, version: IFormVersion) {
  return createStore<IFormDefinitionStore>(() => ({
    form,
    version,
  }));
}
