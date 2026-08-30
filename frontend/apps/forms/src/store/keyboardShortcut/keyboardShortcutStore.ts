import { createStore } from "zustand";
import type { KeyboardShortcut } from "./keyboardShortcut";

export interface IKeyboardShortcutStore {
  shortcuts: KeyboardShortcut[];
  register: (shortcut: KeyboardShortcut) => void;
  unregister: (name: string) => void;
}

export type KeyboardShortcutApi = ReturnType<
  typeof createKeyboardShortcutStore
>;

export function createKeyboardShortcutStore() {
  return createStore<IKeyboardShortcutStore>((set) => ({
    shortcuts: [],
    register: (shortcut: KeyboardShortcut) =>
      set((s) => {
        const filtered = s.shortcuts.filter((sc) => sc.name !== shortcut.name);
        return { shortcuts: [...filtered, shortcut] };
      }),
    unregister: (name: string) =>
      set((s) => {
        const filtered = s.shortcuts.filter((sc) => sc.name !== name);
        return { shortcuts: filtered };
      }),
  }));
}
