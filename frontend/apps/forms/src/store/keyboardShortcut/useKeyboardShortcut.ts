import { useContext, useEffect } from "react";
import { KeyboardShortcutStoreContext } from "./KeyboardShortcutProvider";
import type { KeyboardShortcut } from "./keyboardShortcut";

export function useKeyboardShortcutContext() {
  const store = useContext(KeyboardShortcutStoreContext);

  if (!store) {
    throw new Error(
      "useFormDesignerContext be used within FormDesignerContext",
    );
  }

  return store;
}

export function useKeyboardShortcut(
  shortcut: KeyboardShortcut,
  deps: React.DependencyList,
) {
  const store = useKeyboardShortcutContext();

  useEffect(() => {
    store.getState().register(shortcut);
    return () => store.getState().unregister(shortcut.name);
  }, deps);
}
