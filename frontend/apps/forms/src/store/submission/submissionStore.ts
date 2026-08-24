import { createStore } from "zustand";
import type { ISubmissionValue } from "@/types/submission";

export type FormValues = Record<string, any>;

export interface ISubmissionStore {
  values: FormValues;
  errors: Record<string, string[]>;
  setValue: (elementId: string, value: any) => void;
  setError: (elementId: string, errors: string[]) => void;
}

export type SubmissionStoreApi = ReturnType<typeof createSubmissionStore>;

export function createSubmissionStore(raw: ISubmissionValue[] = []) {
  return createStore<ISubmissionStore>((set) => ({
    values: createFormValues(raw),
    errors: {},
    setValue(elementId, value) {
      return set((s) => ({ values: { ...s.values, [elementId]: value } }));
    },
    setError(elementId, errors) {
      return set((s) => ({ errors: { ...s.errors, [elementId]: errors } }));
    },
  }));
}

export function createFormValues(
  raw: ISubmissionValue[] | undefined,
): FormValues {
  if (!raw) {
    return {};
  }

  return Object.fromEntries(
    raw.map(({ elementId, value }) => [elementId, value]),
  );
}
