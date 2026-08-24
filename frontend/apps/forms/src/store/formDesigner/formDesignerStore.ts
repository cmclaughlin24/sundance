import { createStore } from "zustand";
import type { FormDesignerEvent } from "./events";
import type { IFormVersion } from "@/types/formVersion";
import type { IForm } from "@/types/form";
import { apply, type IFormAggregate } from "./eventHandlers";

export interface IFormDesignerStore {
  snapshot: IFormAggregate;
  events: FormDesignerEvent[];
  dispatch: (event: FormDesignerEvent) => void;
}

export type FormDesignerStoreApi = ReturnType<typeof createFormDesignerStore>;

export function createFormDesignerStore(form: IForm, version: IFormVersion) {
  return createStore<IFormDesignerStore>((set) => ({
    snapshot: { form, version },
    events: [],
    dispatch: (event) => {
      return set((s) => {
        const events = [...s.events, event];
        const snapshot = apply(s.snapshot!, event);

        return { ...s, events, snapshot };
      });
    },
  }));
}
