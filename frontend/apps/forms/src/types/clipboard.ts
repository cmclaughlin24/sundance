import type { IElement } from "./element";
import type { ISection } from "./section";

export enum ClipboardEventType {
  CopyElement = "copy-element",
  CopySection = "copy-section",
}

export interface CopyElementClipboardData {
  type: ClipboardEventType.CopyElement;
  element: IElement;
}

export interface CopySectionClipboardData {
  type: ClipboardEventType.CopySection;
  section: ISection;
  pageId: string;
}

export type ClipboardData = CopyElementClipboardData | CopySectionClipboardData;
