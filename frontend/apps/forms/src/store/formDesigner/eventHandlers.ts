import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import type {
  AddElementEvent,
  AddPageEvent,
  AddSectionEvent,
  FormDesignerEvent,
  MoveElementEvent,
  MovePageEvent,
  MoveSectionEvent,
  RemoveElementEvent,
  RemovePageEvent,
  RemoveSectionEvent,
} from "./events";

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
  AddPage: onAddPage,
  MovePage: onMovePage,
  RemovePage: onRemovePage,
  AddSection: onAddSection,
  MoveSection: onMoveSection,
  RemoveSection: onRemoveSection,
  AddElement: onAddElement,
  MoveElement: onMoveElement,
  RemoveElement: onRemoveElement,
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

function onAddPage(aggregate: IFormAggregate, _event: AddPageEvent) {
  return aggregate;
}

function onMovePage(aggregate: IFormAggregate, _event: MovePageEvent) {
  return aggregate;
}

function onRemovePage(aggregate: IFormAggregate, _event: RemovePageEvent) {
  return aggregate;
}

function onAddSection(aggregate: IFormAggregate, _event: AddSectionEvent) {
  return aggregate;
}

function onMoveSection(aggregate: IFormAggregate, _event: MoveSectionEvent) {
  return aggregate;
}

function onRemoveSection(
  aggregate: IFormAggregate,
  _event: RemoveSectionEvent,
) {
  return aggregate;
}

function onAddElement(aggregate: IFormAggregate, _event: AddElementEvent) {
  return aggregate;
}

function onMoveElement(aggregate: IFormAggregate, _event: MoveElementEvent) {
  return aggregate;
}

function onRemoveElement(
  aggregate: IFormAggregate,
  _event: RemoveElementEvent,
) {
  return aggregate;
}
