import type { IElement } from "./element";
import type { IPage } from "./page";
import type { ISection } from "./section";

export enum ClipboardEventType {
  CopyElement = "copy-element",
  CopySection = "copy-section",
  CopyPage = "copy-page",
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

export interface PagesClipboardData {
  type: ClipboardEventType.CopyPage;
  page: IPage;
}

export type ClipboardData = ElementClipboardData | SectionClipboardData | PagesClipboardData;
