import { createContext, useRef } from "react";
import type { ISubmissionValue } from "@/types/submission";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import { createFormStore, type FormStoreApi } from "./formStore";

export type FormProviderProps = React.PropsWithChildren<{
  form: IForm;
  version: IFormVersion;
  rawSubmission: ISubmissionValue[] | undefined;
}>;

export const FormStoreContext = createContext<FormStoreApi | null>(null);

export const FormStoreProvider: React.FC<FormProviderProps> = ({
  children,
  form,
  version,
  rawSubmission,
}) => {
  const storeRef = useRef<FormStoreApi | null>(null);

  if (!storeRef.current) {
    storeRef.current = createFormStore(form, version, rawSubmission);
  }

  return (
    <FormStoreContext value={storeRef.current}>{children}</FormStoreContext>
  );
};
