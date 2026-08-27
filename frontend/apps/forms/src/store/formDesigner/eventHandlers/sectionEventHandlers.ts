import type { ISection } from "@/types/section";
import type { AddSectionEvent, MoveSectionEvent, RemoveSectionEvent } from "../events";
import type { IFormAggregate } from "./eventHandler";
import type { IPage } from "@/types/page";
import { insertAtPosition, removeById } from "./utils";

export function onAddSection(
  aggregate: IFormAggregate,
  _event: AddSectionEvent,
) {
  return aggregate;
}

export function onMoveSection(
  aggregate: IFormAggregate,
  event: MoveSectionEvent,
): IFormAggregate {
  let section: ISection | undefined;
  let sourcePage: IPage | undefined;

  for (const page of aggregate.version.pages) {
    const found = page.sections.find((s) => s.id === event.sectionId);

    if (found) {
      section = found;
      sourcePage = page;
      break;
    }
  }

  if (!section || !sourcePage) {
    return aggregate;
  }

  const pages = aggregate.version.pages.map((page): IPage => {
    if (page.id === sourcePage!.id && page.id === event.targetPageId) {
      const without = removeById(page.sections, event.sectionId);

      return {
        ...page,
        sections: insertAtPosition(without, section!, event.position),
      };
    }

    if (page.id === sourcePage!.id) {
      return { ...page, sections: removeById(page.sections, event.sectionId) };
    }

    if (page.id === event.targetPageId) {
      return {
        ...page,
        sections: insertAtPosition(page.sections, section!, event.position),
      };
    }

    return page;
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onRemoveSection(
  aggregate: IFormAggregate,
  event: RemoveSectionEvent,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some((s) => s.id === event.sectionId);

    if (!hasSection) {
      return page;
    }

    return { ...page, sections: removeById(page.sections, event.sectionId) };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}
