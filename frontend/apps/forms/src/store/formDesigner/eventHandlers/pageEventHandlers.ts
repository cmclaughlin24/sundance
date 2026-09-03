import type {
  AddPageEvent,
  MovePageEvent,
  PastePageEvent,
  RemovePageEvent,
} from "../events";
import type { IFormAggregate } from "./eventHandler";
import { insertAtPosition, removeById } from "./utils";
import { generatedID } from "@/utils/id";
import { copyKey, copyName } from "@/utils/copy";
import { getNextPosition } from "@/utils/position";

export function onAddPage(aggregate: IFormAggregate, _event: AddPageEvent) {
  return aggregate;
}

export function onPastePage(
  aggregate: IFormAggregate,
  event: PastePageEvent,
): IFormAggregate {
  const page = {
    ...event.page,
    id: generatedID(),
    key: copyKey(event.page.key),
    name: copyName(event.page.name),
    sections: event.page.sections.map((section) => ({
      ...section,
      id: generatedID(),
      key: copyKey(section.key),
      name: copyName(section.name),
      elements: section.elements.map((element) => ({
        ...element,
        id: generatedID(),
        key: copyKey(element.key),
        name: copyName(element.name),
      })),
    })),
  };

  const pages = insertAtPosition(
    aggregate.version.pages,
    page,
    getNextPosition(aggregate.version.pages),
  );

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
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
  const pages = removeById(aggregate.version.pages, event.id);

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}
