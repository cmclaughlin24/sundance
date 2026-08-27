import { createStore } from "zustand";
import type { FormDesignerEvent } from "./events";
import type { IFormVersion } from "@/types/formVersion";
import type { IForm } from "@/types/form";
import type { SelectedItem } from "./formDesigner.type";
import { apply, reduce, type IFormAggregate } from "./eventHandlers/eventHandler";

export interface IFormDesignerStore {
  snapshot: IFormAggregate;
  events: FormDesignerEvent[];
  cursor: number;
  selected: SelectedItem | null;
  dispatch: (event: FormDesignerEvent) => void;
  undo: () => void;
  select: (item: SelectedItem | null) => void;
}

export type FormDesignerStoreApi = ReturnType<typeof createFormDesignerStore>;

export function createFormDesignerStore(form: IForm, version: IFormVersion) {
  const initial: Readonly<IFormAggregate> = { form, version };

  return createStore<IFormDesignerStore>((set) => ({
    snapshot: { form, version },
    events: [],
    cursor: -1,
    selected: null,
    dispatch: (event) =>
      set((s) => {
        const cursor = s.cursor + 1;
        const events = [...s.events.slice(0, cursor), event];
        const snapshot = apply(s.snapshot!, event);

        return { ...s, events, cursor, snapshot };
      }),
    undo: () =>
      set((s) => {
        const cursor = s.cursor >= 0 ? s.cursor - 1 : s.cursor;
        const events = cursor !== -1 ? s.events.slice(0, s.cursor) : [];
        const snapshot = reduce(initial, events);

        return { ...s, cursor, snapshot };
      }),
    select: (item) => set((s) => ({ ...s, selected: item })),
  }));
}
