import { useReducer } from "react";
import {
  formReducer,
  initialFormState,
  initializeForm,
  type FormState,
} from "./formReducer";
import { FormDispatchContext, FormStateContext } from "./FormContext";
import type { ISubmissionValue } from "@/types/submission";

export type FormProviderProps = React.PropsWithChildren<{
  rawSubmission: ISubmissionValue[] | undefined;
}>;

export const FormProvider: React.FC<FormProviderProps> = ({
  children,
  rawSubmission,
}) => {
  const [state, dispatch] = useReducer(
    formReducer,
    initialFormState,
    (state: FormState): FormState => {
      return rawSubmission ? initializeForm(state, rawSubmission) : state;
    },
  );

  return (
    <FormStateContext value={state}>
      <FormDispatchContext value={dispatch}>{children}</FormDispatchContext>
    </FormStateContext>
  );
};
