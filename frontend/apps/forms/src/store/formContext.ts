import { createContext } from "react";
import {
  initialFormState,
  type FormAction,
  type FormState,
} from "./formReducer";

export const FormDispatchContext = createContext<React.Dispatch<FormAction>>(
  () => {},
);

export const FormStateContext = createContext<FormState>(initialFormState);
