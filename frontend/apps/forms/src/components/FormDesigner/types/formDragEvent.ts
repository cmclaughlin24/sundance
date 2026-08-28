import type { ElementType } from "@/types/element";

export enum FormDragEventSource {
  Palette = "palette",
}

export interface PaletteDragEventData {
  source: FormDragEventSource.Palette;
  type: ElementType | "section";
}

export type FormDragEventData = PaletteDragEventData;
