import type { ISection } from "@/types/section";
import type {
  AddSectionEvent,
  MoveSectionEvent,
  PasteSectionEvent,
  ReorderSectionEvent,
  RemoveSectionEvent,
  UpdateSectionEvent,
  CutSectionEvent,
} from "../events";
import type { IFormAggregate } from "./eventHandler";
import type { IPage } from "@/types/page";
import { insertAtPosition, removeById } from "./utils";
import { createEmptySection } from "@/factories/sectionFactory";
import { swapPositions, getNextPosition } from "@/utils/position";
import { generatedID } from "@/utils/id";
import { copyKey, copyName } from "@/utils/copy";
import { ClipboardEventType } from "@/types/clipboard";

export function onAddSection(
  aggregate: IFormAggregate,
  event: AddSectionEvent,
): IFormAggregate {
  const section = createEmptySection(event.id);

  const pages = aggregate.version.pages.map((page): IPage => {
    if (page.id !== event.pageId) {
      return page;
    }

    return {
      ...page,
      sections: insertAtPosition(page.sections, section, event.position),
    };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onUpdateSection(
  aggregate: IFormAggregate,
  event: UpdateSectionEvent,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some((s) => s.id === event.sectionId);

    if (!hasSection) {
      return page;
    }

    const sections = page.sections.map((section) => {
      if (section.id !== event.sectionId) {
        return section;
      }

      return { ...section, ...event.changes };
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onReorderSection(
  aggregate: IFormAggregate,
  event: ReorderSectionEvent,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some((s) => s.id === event.sectionId);

    if (!hasSection) {
      return page;
    }

    return {
      ...page,
      sections: swapPositions(page.sections, event.sectionId, event.inc),
    };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
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
  return removeSectionById(aggregate, event.sectionId);
}

export function onCutSection(
  aggregate: IFormAggregate,
  event: CutSectionEvent,
): IFormAggregate {
  return removeSectionById(aggregate, event.sectionId);
}

function removeSectionById(
  aggregate: IFormAggregate,
  sectionId: string,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some((s) => s.id === sectionId);

    if (!hasSection) {
      return page;
    }

    return { ...page, sections: removeById(page.sections, sectionId) };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onPasteSection(
  aggregate: IFormAggregate,
  event: PasteSectionEvent,
): IFormAggregate {
  const perserve = event.clipboardOp === ClipboardEventType.CutSection;

  const section: ISection = {
    ...event.section,
    id: perserve ? event.section.id : generatedID(),
    key: perserve ? event.section.key : copyKey(event.section.key),
    name: perserve ? event.section.name : copyName(event.section.name),
    elements: event.section.elements.map((element) => ({
      ...element,
      id: perserve ? element.id : generatedID(),
      key: perserve ? element.key : copyKey(element.key),
      name: perserve ? element.name : copyName(element.name),
    })),
  };

  const pages = aggregate.version.pages.map((page): IPage => {
    if (page.id !== event.targetPageId) {
      return page;
    }

    return {
      ...page,
      sections: insertAtPosition(
        page.sections,
        section,
        getNextPosition(page.sections),
      ),
    };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}
