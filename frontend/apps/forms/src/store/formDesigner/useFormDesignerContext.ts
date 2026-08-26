import { useContext } from "react";
import { FormDesignerContext } from "./FormDesignerProvider";
import { useStore } from "zustand";
import { useShallow } from "zustand/react/shallow";

export function useFormDesignerContext() {
  const store = useContext(FormDesignerContext);

  if (!store) {
    throw new Error(
      "useFormDesignerContext be used within FormDesignerContext",
    );
  }

  return store;
}

export function useFormDesignerDispatch() {
  const store = useFormDesignerContext();
  return useStore(
    store,
    useShallow((s) => ({
      dispatch: s.dispatch,
      undo: s.undo,
      select: s.select,
    })),
  );
}

export function useFormSnapshot() {
  const store = useFormDesignerContext();
  return useStore(
    store,
    useShallow((s) => s.snapshot),
  );
}

export function useFormPagesSnapshot() {
  const store = useFormDesignerContext();
  return useStore(
    store,
    useShallow((s) => s.snapshot.version.pages),
  );
}

export function useSelectedItem() {
  const store = useFormDesignerContext();
  return useStore(store, (s) => s.selected);
}
