import { createStore } from "zustand";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import type { ISubmissionValue } from "@/types/submission";

export type FormValues = Record<string, any>;

export interface IFormStore {
  form: Readonly<IForm> | null;
  version: Readonly<IFormVersion> | null;
  values: FormValues;
  setValue: (elementId: string, value: any) => void;
  setError: (elementId: string, errors: string[]) => void;
}

export type FormStoreApi = ReturnType<typeof createFormStore>;

export function createFormStore(
  form: IForm,
  version: IFormVersion,
  raw: ISubmissionValue[] = [],
) {
  return createStore<IFormStore>((set) => ({
    form,
    version,
    values: createFormValues(raw),
    setValue(elementId, value) {
      return set((s) => ({ values: { ...s.values, [elementId]: value } }));
    },
    setError(_elementId, _errors) {},
  }));
}

function createFormValues(raw: ISubmissionValue[]): FormValues {
  if (!raw) {
    return {};
  }

  return Object.fromEntries(
    raw.map(({ elementId, value }) => [elementId, value]),
  );
}
