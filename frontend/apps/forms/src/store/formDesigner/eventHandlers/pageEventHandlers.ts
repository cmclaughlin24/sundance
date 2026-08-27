import type { AddPageEvent, MovePageEvent, RemovePageEvent } from "../events";
import type { IFormAggregate } from "./eventHandler";
import { insertAtPosition, removeById } from "./utils";

export function onAddPage(aggregate: IFormAggregate, _event: AddPageEvent) {
  return aggregate;
}

export function onMovePage(
  aggregate: IFormAggregate,
  event: MovePageEvent,
): IFormAggregate {
  const page = aggregate.version.pages.find((p) => p.id === event.pageId);

  if (!page) {
    return aggregate;
  }

  const without = removeById(aggregate.version.pages, event.pageId);
  const pages = insertAtPosition(without, page, event.position);

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onRemovePage(
  aggregate: IFormAggregate,
  event: RemovePageEvent,
): IFormAggregate {
  const pages = removeById(aggregate.version.pages, event.pageId);

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}
