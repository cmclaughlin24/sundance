import { createStore } from "zustand";
import type { FormDesignerEvent } from "./events";
import type { IFormVersion } from "@/types/formVersion";
import type { IForm } from "@/types/form";
import type { SelectedItem } from "./formDesigner.type";
import {
  apply,
  reduce,
  type IFormAggregate,
} from "./eventHandlers/eventHandler";
import { findSelectedById } from "@/utils/form";
import { extractFlatRules } from "@/utils/rule";

export interface IFormDesignerStore {
  baseline: IFormAggregate;
  snapshot: IFormAggregate;
  events: FormDesignerEvent[];
  cursor: number;
  selected: SelectedItem | null;
  dispatch: (event: FormDesignerEvent) => void;
  undo: () => void;
  redo: () => void;
  commit: (version: IFormVersion) => void;
  select: (item: SelectedItem | null) => void;
}

export type FormDesignerStoreApi = ReturnType<typeof createFormDesignerStore>;

export function createFormDesignerStore(form: IForm, version: IFormVersion) {
  return createStore<IFormDesignerStore>((set) => {
    const rules = extractFlatRules(version.pages);

    return {
      baseline: { form, version, rules },
      snapshot: { form, version, rules },
      events: [],
      cursor: -1,
      selected: null,
      dispatch: (event) =>
        set((s) => {
          const cursor = s.cursor + 1;
          const events = [...s.events.slice(0, cursor), event];
          const snapshot = apply(s.snapshot!, event);
          let selected = s.selected;

          if (shouldUpdateSelection(event)) {
            selected = findSelectedById(
              snapshot.version.pages,
              (event as { id: string }).id,
            );
          } else if (isRemoveEvent(event)) {
            selected = null;
          }

          return { ...s, events, cursor, snapshot, selected };
        }),
      commit: (version) =>
        set((s) => {
          const selected = s.selected
            ? findSelectedById(version.pages, s.selected.item.id)
            : null;
          const rules = extractFlatRules(version.pages);

          return {
            baseline: { ...s.snapshot, version, rules },
            snapshot: { ...s.snapshot, version, rules },
            events: [],
            cursor: -1,
            selected,
          };
        }),
      undo: () =>
        set((s) => {
          const cursor = s.cursor >= 0 ? s.cursor - 1 : s.cursor;
          const events = cursor !== -1 ? s.events.slice(0, s.cursor) : [];
          const snapshot = reduce(s.baseline, events);

          return { ...s, cursor, snapshot };
        }),
      redo: () =>
        set((s) => {
          const cursor =
            s.cursor + 1 < s.events.length ? s.cursor + 1 : s.cursor;
          const events = s.events.slice(0, cursor + 1);
          const snapshot = reduce(s.baseline, events);

          return { ...s, cursor, snapshot };
        }),
      select: (item) => set((s) => ({ ...s, selected: item })),
    };
  });
}

function shouldUpdateSelection(event: FormDesignerEvent): boolean {
  return !isRuleEvent(event) && (isAddEvent(event) || isUpdateEvent(event));
}

function isAddEvent(event: FormDesignerEvent): boolean {
  return (
    event.type === "AddPage" ||
    event.type === "AddSection" ||
    event.type === "AddElement"
  );
}

function isRemoveEvent(event: FormDesignerEvent): boolean {
  return (
    event.type === "RemoveSection" ||
    event.type === "CutSection" ||
    event.type === "RemoveElement" ||
    event.type === "CutElement"
  );
}

function isUpdateEvent(event: FormDesignerEvent): boolean {
  return event.type === "UpdateElement" || event.type === "UpdateSection";
}

function isRuleEvent(event: FormDesignerEvent): boolean {
  return (
    event.type === "AddRule" ||
    event.type === "UpdateRule" ||
    event.type === "RemoveRule" ||
    event.type === "PasteRule"
  );
}
