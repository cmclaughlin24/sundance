import type { ClipboardEventType } from "@/types/clipboard";
import type {
  ElementType,
  IElement,
  IElementTagMapping,
} from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";
import type { IPage } from "@/types/page";
import type { IFlatRule, RuleType } from "@/types/rule";
import type { ISection } from "@/types/section";

export type FormDesignerEvent =
  | AddPageEvent
  | MovePageEvent
  | RemovePageEvent
  | PastePageEvent
  | AddSectionEvent
  | MoveSectionEvent
  | UpdateSectionEvent
  | RemoveSectionEvent
  | ReorderSectionEvent
  | PasteSectionEvent
  | CutSectionEvent
  | AddElementEvent
  | MoveElementEvent
  | UpdateElementEvent
  | RemoveElementEvent
  | ReorderElementEvent
  | PasteElementEvent
  | CutElementEvent
  | AddElementTagEvent
  | UpdateElementTagEvent
  | RemoveElementTagEvent
  | AddRuleEvent
  | UpdateRuleEvent
  | RemoveRuleEvent
  | PasteRuleEvent;

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
  id: string;
};

export type PastePageEvent = {
  type: "PastePage";
  page: IPage;
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

export type UpdateSectionEvent = {
  type: "UpdateSection";
  id: string;
  changes: Partial<Pick<ISection, "key" | "name">>;
};

export type RemoveSectionEvent = {
  type: "RemoveSection";
  id: string;
};

export type ReorderSectionEvent = {
  type: "ReorderSection";
  sectionId: string;
  inc?: -1 | 1;
  targetIndex?: number;
};

export type CutSectionEvent = {
  type: "CutSection";
  sectionId: string;
};

export type PasteSectionEvent = {
  type: "PasteSection";
  section: ISection;
  targetPageId: string;
  clipboardOp: ClipboardEventType.CutSection | ClipboardEventType.CopySection;
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

export type UpdateElementEvent = {
  type: "UpdateElement";
  id: string;
  changes: Partial<
    Pick<IElement, "key" | "name" | "description"> & {
      attributes: Partial<ElementAttributes>;
    }
  >;
};

export type RemoveElementEvent = {
  type: "RemoveElement";
  id: string;
};

export type ReorderElementEvent = {
  type: "ReorderElement";
  elementId: string;
  inc?: -1 | 1;
  targetIndex?: number;
};

export type CutElementEvent = {
  type: "CutElement";
  elementId: string;
};

export type PasteElementEvent = {
  type: "PasteElement";
  element: IElement;
  targetSectionId: string;
  clipboardOp: ClipboardEventType.CutElement | ClipboardEventType.CopyElement;
};

export type AddElementTagEvent = {
  type: "AddElementTag";
  id: string;
  mapping: IElementTagMapping;
};

export type UpdateElementTagEvent = {
  type: "UpdateElementTag";
  id: string;
  tagVersionId: string;
  mapping: Partial<IElementTagMapping>;
};

export type RemoveElementTagEvent = {
  type: "RemoveElementTag";
  id: string;
  tagVersionId: string;
};

export type AddRuleEvent = {
  type: "AddRule";
  id: string;
  ruleType: RuleType;
};

export type UpdateRuleEvent = {
  type: "UpdateRule";
  id: string;
  changes: Partial<Pick<IFlatRule, "expressions" | "parentId" | "parentType">>;
};

export type RemoveRuleEvent = {
  type: "RemoveRule";
  id: string;
};

export type PasteRuleEvent = {
  type: "PasteRule";
  rule: IFlatRule;
};
