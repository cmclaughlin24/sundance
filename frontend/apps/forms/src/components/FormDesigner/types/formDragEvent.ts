import type { ElementType, IElement } from "@/types/element";
import type { ISection } from "@/types/section";

export enum FormDragEventSource {
  Palette = "palette",
  Canvas = "canvas",
}

export enum BuilderItemDragType {
  Element = "palette-element",
  Section = "palette-section",
}

export enum RuleItemDragType {
  Rule = "palette-rule",
}

export interface PaletteDragEventData<T> {
  source: FormDragEventSource.Palette;
  itemType: T;
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

export type BuilderDragEventData =
  | PaletteDragEventData<ElementType | "section">
  | CanvasElementDragEventData
  | CanvasSectionDragEventData;
