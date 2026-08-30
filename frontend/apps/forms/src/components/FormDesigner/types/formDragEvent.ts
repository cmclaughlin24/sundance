import type { ElementType, IElement } from "@/types/element";
import type { ISection } from "@/types/section";

export enum FormDragEventSource {
  Palette = "palette",
  Canvas = "canvas",
}

export enum PaletteItemDragType {
  Element = "palette-element",
  Section = "palette-section",
}

export interface PaletteDragEventData {
  source: FormDragEventSource.Palette;
  objectType: ElementType | "section";
}

export enum CanvasDragType {
  Element = "canvas-element",
  Section = "canvas-section",
}

export interface CanvasSectionDragEventData {
  source: FormDragEventSource.Canvas;
  type: "section";
  section: ISection;
  fromPageId: string;
}

export interface CanvasElementDragEventData {
  source: FormDragEventSource.Canvas;
  type: "element";
  element: IElement;
  fromSectionId: string;
}

export type FormDragEventData =
  | PaletteDragEventData
  | CanvasElementDragEventData
  | CanvasSectionDragEventData;
