import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import { createContext, useEffect, useRef } from "react";
import {
  createFormDesignerStore,
  type FormDesignerStoreApi,
} from "./formDesignerStore";
import { formsDraftStorage } from "@/services/storage/formDesignerStorage";
import type { FormDesignerDraft } from "./formDesigner.type";
import * as ArrayUtils from "@/utils/array";

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
    storeRef.current = createFormDesignerStore(
      form,
      version,
      formsDraftStorage,
    );
  }

  useEffect(() => {
    let isMounted = true;

    formsDraftStorage.findByKey(version.id).then((draft) => {
      if (!isMounted || !draft) {
        return;
      }

      if (!isValidDraft(draft, version)) {
        return formsDraftStorage.delete(version.id);
      }

      storeRef.current!.getState().hydrate(draft);
    });

    return () => {
      isMounted = false;
    };
  }, [version.id]);

  return (
    <FormDesignerContext value={storeRef.current}>
      {children}
    </FormDesignerContext>
  );
};

function isValidDraft(
  draft: FormDesignerDraft,
  version: IFormVersion,
): boolean {
  if (!ArrayUtils.hasLengthGreaterThan(draft.events, 0)) {
    return false;
  }

  return draft.updatedAt > new Date(version.updatedAt).getTime();
}
