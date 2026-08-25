import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import { createContext, useRef } from "react";
import {
  createFormDefinitionStore,
  type FormDefinitionStoreApi,
} from "./formDefinitionStore";

export type FormDefinitionProviderProps = React.PropsWithChildren<{
  form: IForm;
  version: IFormVersion;
}>;

export const FormDefinitionContext =
  createContext<FormDefinitionStoreApi | null>(null);

export const FormDefinitionProvider: React.FC<FormDefinitionProviderProps> = ({
  children,
  form,
  version,
}) => {
  const storeRef = useRef<FormDefinitionStoreApi | null>(null);

  if (!storeRef.current) {
    storeRef.current = createFormDefinitionStore(form, version);
  }

  return (
    <FormDefinitionContext value={storeRef.current}>
      {children}
    </FormDefinitionContext>
  );
};
