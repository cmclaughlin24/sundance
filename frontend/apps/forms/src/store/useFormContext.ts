import { useContext } from "react";
import { FormDispatchContext, FormStateContext } from "./FormContext";

export function useFormState() {
  return useContext(FormStateContext);
}

export function useFormDispatch() {
  return useContext(FormDispatchContext);
}
