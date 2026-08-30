import type { ElementType, IElement } from "@/types/element";
import type { ISection } from "@/types/section";

export type FormDesignerEvent =
  | AddPageEvent
  | MovePageEvent
  | RemovePageEvent
  | AddSectionEvent
  | MoveSectionEvent
  | RemoveSectionEvent
  | ReorderSectionEvent
  | PasteSectionEvent
  | AddElementEvent
  | MoveElementEvent
  | RemoveElementEvent
  | ReorderElementEvent
  | PasteElementEvent;

export type AddPageEvent = {
  type: "AddPage";
  position: number;
};

export type MovePageEvent = {
  type: "MovePage";
  pageId: string;
  position: number;
};

export type RemovePageEvent = {
  type: "RemovePage";
  pageId: string;
};

export type AddSectionEvent = {
  type: "AddSection";
  id: string;
  pageId: string;
  position: number;
};

export type MoveSectionEvent = {
  type: "MoveSection";
  sectionId: string;
  targetPageId: string;
  position: number;
};

export type RemoveSectionEvent = {
  type: "RemoveSection";
  sectionId: string;
};

export type ReorderSectionEvent = {
  type: "ReorderSection";
  sectionId: string;
  inc: -1 | 1;
};

export type PasteSectionEvent = {
  type: "PasteSection";
  section: ISection;
  targetPageId: string;
};

export type AddElementEvent = {
  type: "AddElement";
  id: string;
  sectionId: string;
  elementType: ElementType;
  position: number;
};

export type MoveElementEvent = {
  type: "MoveElement";
  elementId: string;
  targetSectionId: string;
  position: number;
};

export type RemoveElementEvent = {
  type: "RemoveElement";
  elementId: string;
};

export type ReorderElementEvent = {
  type: "ReorderElement";
  elementId: string;
  inc: -1 | 1;
};

export type PasteElementEvent = {
  type: "PasteElement";
  element: IElement;
  targetSectionId: string;
};
