import { useContext } from "react";
import { ContextMenuContext } from "./ContextMenuProvider";
import { useStore } from "zustand";
import { useShallow } from "zustand/react/shallow";

export function useContextMenuContext() {
  const store = useContext(ContextMenuContext);

  if (!store) {
    throw new Error(
      "useContextMenuContext must be used within ContextMenuContext",
    );
  }

  return store;
}

export function useContextMenuDispatch() {
  const store = useContextMenuContext();
  return useStore(
    store,
    useShallow((s) => ({
      open: s.open,
      close: s.close,
    })),
  );
}

export function useContextMenuTarget() {
  const store = useContextMenuContext();
  return useStore(store, (s) => s.target);
}
