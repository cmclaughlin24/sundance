import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import { createContext, useRef } from "react";
import {
  createFormDesignerStore,
  type FormDesignerStoreApi,
} from "./formDesignerStore";

export type FormDesignerProps = React.PropsWithChildren<{
  form: IForm;
  version: IFormVersion;
}>;

export const FormDesignerContext = createContext<FormDesignerStoreApi | null>(
  null,
);

export const FormDesignerProvider: React.FC<FormDesignerProps> = ({
  children,
  form,
  version,
}) => {
  const storeRef = useRef<FormDesignerStoreApi | null>(null);

  if (!storeRef.current) {
    storeRef.current = createFormDesignerStore(form, version);
  }

  return (
    <FormDesignerContext value={storeRef.current}>
      {children}
    </FormDesignerContext>
  );
};
