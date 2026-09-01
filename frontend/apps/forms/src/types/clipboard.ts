import type { IElement } from "./element";
import type { ISection } from "./section";

export enum ClipboardEventType {
  CopyElement = "copy-element",
  CopySection = "copy-section",
  CutElement = "cut-element",
  CutSection = "cut-section",
}

export interface ElementClipboardData {
  type: ClipboardEventType.CopyElement | ClipboardEventType.CutElement;
  element: IElement;
}

export interface SectionClipboardData {
  type: ClipboardEventType.CopySection | ClipboardEventType.CutSection;
  section: ISection;
}

export type ClipboardData = ElementClipboardData | SectionClipboardData;
