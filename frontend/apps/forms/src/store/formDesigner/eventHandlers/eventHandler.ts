import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import type { IFlatRule } from "@/types/rule";
import type { FormDesignerEvent } from "../events";
import * as pageHandlers from "./pageEventHandlers";
import * as sectionHandlers from "./sectionEventHandlers";
import * as elementHandlers from "./elementEventHandlers";
import * as ruleHandlers from "./ruleEventHandlers";

export interface IFormAggregate {
  form: IForm;
  version: IFormVersion;
  rules: IFlatRule[];
}

type Handlers = {
  [E in FormDesignerEvent as E["type"]]: (
    state: IFormAggregate,
    event: E,
  ) => IFormAggregate;
};

const handlers: Readonly<Handlers> = {
  AddPage: pageHandlers.onAddPage,
  MovePage: pageHandlers.onMovePage,
  RemovePage: pageHandlers.onRemovePage,
  PastePage: pageHandlers.onPastePage,
  AddSection: sectionHandlers.onAddSection,
  MoveSection: sectionHandlers.onMoveSection,
  UpdateSection: sectionHandlers.onUpdateSection,
  RemoveSection: sectionHandlers.onRemoveSection,
  ReorderSection: sectionHandlers.onReorderSection,
  PasteSection: sectionHandlers.onPasteSection,
  CutSection: sectionHandlers.onCutSection,
  AddElement: elementHandlers.onAddElement,
  MoveElement: elementHandlers.onMoveElement,
  UpdateElement: elementHandlers.onUpdateElement,
  RemoveElement: elementHandlers.onRemoveElement,
  ReorderElement: elementHandlers.onReorderElement,
  PasteElement: elementHandlers.onPasteElement,
  CutElement: elementHandlers.onCutElement,
  AddRule: ruleHandlers.onAddRule,
  UpdateRule: ruleHandlers.onUpdateRule,
  RemoveRule: ruleHandlers.onRemoveRule,
  PasteRule: ruleHandlers.onPasteRule,
};

export function reduce(
  aggregate: IFormAggregate,
  events: FormDesignerEvent[],
): IFormAggregate {
  if (!events) {
    return aggregate;
  }

  let cpy = { ...aggregate };

  for (const event of events) {
    cpy = apply(cpy, event);
  }

  return cpy;
}

export function apply(
  aggregate: IFormAggregate,
  event: FormDesignerEvent,
): IFormAggregate {
  const handler = handlers[event.type] as (
    aggregate: IFormAggregate,
    event: FormDesignerEvent,
  ) => IFormAggregate;

  if (!handler) {
    return aggregate;
  }

  return handler(aggregate, event);
}
