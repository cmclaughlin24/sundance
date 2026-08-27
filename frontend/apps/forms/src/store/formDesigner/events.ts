export type FormDesignerEvent =
  | AddPageEvent
  | MovePageEvent
  | RemovePageEvent
  | AddSectionEvent
  | MoveSectionEvent
  | RemoveSectionEvent
  | AddElementEvent
  | MoveElementEvent
  | RemoveElementEvent;

export type AddPageEvent = {
  type: "AddPage";
  pageId: string;
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
  sectionId: string;
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

export type AddElementEvent = {
  type: "AddElement";
  elementId: string;
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
