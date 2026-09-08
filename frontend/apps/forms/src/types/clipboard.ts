import type { IElement } from "./element";
import type { IFlatRule } from "./rule";
import type { IPage } from "./page";
import type { ISection } from "./section";

export enum ClipboardEventType {
  CopyElement = "copy-element",
  CopySection = "copy-section",
  CopyPage = "copy-page",
  CopyRule = "copy-rule",
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

export interface RuleClipboardData {
  type: ClipboardEventType.CopyRule;
  rule: IFlatRule;
}

export type ClipboardData =
  | ElementClipboardData
  | SectionClipboardData
  | PagesClipboardData
  | RuleClipboardData;
