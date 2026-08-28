import type { ISection } from "@/types/section";
import { generatedID } from "@/utils/id";

export function createEmptySection(): ISection {
  return {
    id: generatedID(),
    key: "",
    name: "Section",
    position: 0,
    rules: [],
    elements: [],
  };
}
