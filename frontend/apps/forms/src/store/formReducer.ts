import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import type { ISubmissionValue } from "@/types/submission";

export interface FormState {
  form: Readonly<IForm> | null;
  version: Readonly<IFormVersion> | null;
  values: Map<string, any>;
}

export const initialFormState: FormState = {
  form: null,
  version: null,
  values: new Map<string, any>(),
};

export type FormAction =
  | {
      type: "INITIALIZE";
      form: IForm;
      version: IFormVersion;
      values: ISubmissionValue[];
    }
  | { type: "SET_VALUE"; elementId: string; value: any }
  | { type: "SET_ERROR"; elementId: string; errors: string[] };

export function formReducer(state: FormState, action: FormAction) {
  switch (action.type) {
    case "INITIALIZE":
      return initializeForm(state, action.form, action.version, action.values);
    case "SET_ERROR":
      return setError(state, action.elementId, action.errors);
    case "SET_VALUE":
      return setValue(state, action.elementId, action.value);
  }
}

export function initializeForm(
  state: FormState,
  form: IForm,
  version: IFormVersion,
  raw: ISubmissionValue[],
): FormState {
  const values = new Map<string, any>();

  for (const { elementId, value } of raw) {
    values.set(elementId, value);
  }

  return { ...state, form, version, values };
}

export function setError(
  state: FormState,
  elementId: string,
  errors: string[],
): FormState {
  console.log(elementId, errors);
  return state;
}

export function setValue(
  state: FormState,
  elementId: string,
  value: any,
): FormState {
  const values = new Map(state.values);
  values.set(elementId, value);
  return { ...state, values };
}
