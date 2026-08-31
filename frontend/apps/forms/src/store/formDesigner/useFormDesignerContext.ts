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
    })),
  );
}

export function useFormDesignerUndo() {
  const store = useFormDesignerContext();
  return useStore(
    store,
    useShallow((s) => ({
      undo: s.undo,
      redo: s.redo,
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

export function useFormDesignerSelect(id?: string) {
  const store = useFormDesignerContext();
  return useStore(
    store,
    useShallow((s) => ({
      selected: s.selected,
      isSelected: !!s.selected && s.selected.item.id === id,
      select: s.select,
    })),
  );
}
