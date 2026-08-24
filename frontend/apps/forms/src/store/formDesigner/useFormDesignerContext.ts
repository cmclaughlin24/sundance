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
