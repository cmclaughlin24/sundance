import type { ISection } from "@/types/section";

export function createEmptySection(id: string): ISection {
  return {
    id,
    key: "",
    name: "Section",
    position: 0,
    rules: [],
    elements: [],
  };
}
