import { createContext, useRef } from "react";
import {
  createContextMenuStore,
  type ContextMenuApi,
} from "./contextMenuStore";

export const ContextMenuContext = createContext<ContextMenuApi | null>(null);

export const ContextMenuProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const storeRef = useRef<ContextMenuApi | null>(null);

    if (!storeRef.current) {
      storeRef.current = createContextMenuStore();
    }

    return (
      <ContextMenuContext value={storeRef.current}>
        {children}
      </ContextMenuContext>
    );
  };
