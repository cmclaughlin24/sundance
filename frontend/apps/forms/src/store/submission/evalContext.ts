import type { EvalContext } from "@/utils/evaluate";
import { createContext, useContext } from "react";

export const EvalContextContext = createContext<EvalContext>({});

export function useEvalContext() {
  return useContext(EvalContextContext);
}
