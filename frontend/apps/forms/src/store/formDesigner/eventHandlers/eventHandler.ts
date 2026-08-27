import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import type { FormDesignerEvent } from "../events";
import * as pageHandlers from "./pageEventHandlers";
import * as sectionHandlers from "./sectionEventHandlers";
import * as elementHandlers from "./elementEventHandlers";

export interface IFormAggregate {
  form: IForm;
  version: IFormVersion;
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
  AddSection: sectionHandlers.onAddSection,
  MoveSection: sectionHandlers.onMoveSection,
  RemoveSection: sectionHandlers.onRemoveSection,
  AddElement: elementHandlers.onAddElement,
  MoveElement: elementHandlers.onMoveElement,
  RemoveElement: elementHandlers.onRemoveElement,
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
