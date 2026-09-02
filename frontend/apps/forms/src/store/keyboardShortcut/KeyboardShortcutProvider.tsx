import { createContext, useEffect, useRef } from "react";
import {
  createKeyboardShortcutStore,
  type KeyboardShortcutApi,
} from "./keyboardShortcutStore";
import { isMatch, isValidTarget } from "./keyboardShortcut";

export const KeyboardShortcutStoreContext =
  createContext<KeyboardShortcutApi | null>(null);

export const KeyboardShortcutProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const storeRef = useRef<KeyboardShortcutApi | null>(null);

    if (!storeRef.current) {
      storeRef.current = createKeyboardShortcutStore();
    }

    useEffect(() => {
      const handler = (event: KeyboardEvent) => {
        const { shortcuts } = storeRef.current!.getState();

        for (const shortcut of shortcuts) {
          if (
            isMatch(event, shortcut) &&
            isValidTarget(shortcut, event.target)
          ) {
            shortcut.action();
            break;
          }
        }
      };

      document.addEventListener("keydown", handler);
      return () => document.removeEventListener("keydown", handler);
    }, []);

    return (
      <KeyboardShortcutStoreContext value={storeRef.current}>
        {children}
      </KeyboardShortcutStoreContext>
    );
  };
