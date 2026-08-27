import type { ElementType } from "@/types/element";

export interface SelectedItem {
  type: "page" | "section" | ElementType;
  id: string;
}
