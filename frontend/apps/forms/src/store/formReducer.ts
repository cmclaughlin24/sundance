import type { ISubmissionValue } from "@/types/submission";

export interface FormState {
  values: Map<string, any>;
}

export const initialFormState: FormState = {
  values: new Map<string, any>(),
};

export type FormAction =
  | { type: "INITIALIZE"; values: ISubmissionValue[] }
  | { type: "SET_VALUE"; elementId: string; value: any }
  | { type: "SET_ERROR"; elementId: string; errors: string[] };

export function formReducer(state: FormState, action: FormAction) {
  switch (action.type) {
    case "INITIALIZE":
      return initializeForm(state, action.values);
    case "SET_ERROR":
      return setError(state, action.elementId, action.errors);
    case "SET_VALUE":
      return setValue(state, action.elementId, action.value);
  }
}

export function initializeForm(
  state: FormState,
  raw: ISubmissionValue[],
): FormState {
  const values = new Map<string, any>();

  for (const { elementId, value } of raw) {
    values.set(elementId, value);
  }

  return { ...state, values };
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
