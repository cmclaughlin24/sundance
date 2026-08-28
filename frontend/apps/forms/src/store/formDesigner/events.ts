import type { ElementType } from "@/types/element";

export type FormDesignerEvent =
  | AddPageEvent
  | MovePageEvent
  | RemovePageEvent
  | AddSectionEvent
  | MoveSectionEvent
  | RemoveSectionEvent
  | ReorderSectionEvent
  | AddElementEvent
  | MoveElementEvent
  | RemoveElementEvent
  | ReorderElementEvent;

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
