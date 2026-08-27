import type { ElementType } from "@/types/element";

export interface PaletteDragEventData {
  source: "palette";
  type: ElementType | "section";
}

export type FormDragEventData = PaletteDragEventData;
