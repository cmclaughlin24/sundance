import { createStore } from "zustand";

export interface ContextMenuTarget {
  position: { x: number; y: number };
  data: unknown;
}

export interface ContextMenuStore {
  target: ContextMenuTarget | null;
  open: (target: ContextMenuTarget) => void;
  close: () => void;
}

export type ContextMenuApi = ReturnType<typeof createContextMenuStore>;

export function createContextMenuStore() {
  return createStore<ContextMenuStore>((set) => ({
    target: null,
    open: (target) => {
      set((_s) => ({ target }));
    },
    close: () => {
      set((_s) => ({ target: null }));
    },
  }));
}
