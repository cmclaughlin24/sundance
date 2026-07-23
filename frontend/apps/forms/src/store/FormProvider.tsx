import { useReducer } from "react";
import {
  formReducer,
  initialFormState,
  initializeForm,
  type FormState,
} from "./formReducer";
import { FormDispatchContext, FormStateContext } from "./formContext";
import type { ISubmissionValue } from "@/types/submission";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";

export type FormProviderProps = React.PropsWithChildren<{
  form: IForm;
  version: IFormVersion;
  rawSubmission: ISubmissionValue[] | undefined;
}>;

export const FormProvider: React.FC<FormProviderProps> = ({
  children,
  form,
  version,
  rawSubmission,
}) => {
  const [state, dispatch] = useReducer(
    formReducer,
    initialFormState,
    (state: FormState): FormState => {
      return initializeForm(state, form, version, rawSubmission ?? []);
    },
  );

  return (
    <FormStateContext value={state}>
      <FormDispatchContext value={dispatch}>{children}</FormDispatchContext>
    </FormStateContext>
  );
};
