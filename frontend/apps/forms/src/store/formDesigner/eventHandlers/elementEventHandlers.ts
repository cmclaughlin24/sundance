import type { IElement } from "@/types/element";
import type {
  AddElementEvent,
  MoveElementEvent,
  PasteElementEvent,
  ReorderElementEvent,
  RemoveElementEvent,
  UpdateElementEvent,
  CutElementEvent,
} from "../events";
import type { IFormAggregate } from "./eventHandler";
import type { ISection } from "@/types/section";
import type { IPage } from "@/types/page";
import { insertAtPosition, removeById } from "./utils";
import { createElementFromType } from "@/factories/elementFactory";
import { swapPositions, getNextPosition } from "@/utils/position";
import { generatedID } from "@/utils/id";
import { copyKey, copyName } from "@/utils/copy";
import { ClipboardEventType } from "@/types/clipboard";

export function onAddElement(
  aggregate: IFormAggregate,
  event: AddElementEvent,
): IFormAggregate {
  const element = createElementFromType(event.elementType, event.id);

  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some((s) => s.id === event.sectionId);

    if (!hasSection) {
      return page;
    }

    const sections = page.sections.map((section): ISection => {
      if (section.id !== event.sectionId) {
        return section;
      }

      return {
        ...section,
        elements: insertAtPosition(section.elements, element, event.position),
      };
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onReorderElement(
  aggregate: IFormAggregate,
  event: ReorderElementEvent,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const sections = page.sections.map((section): ISection => {
      const hasElement = section.elements.some((e) => e.id === event.elementId);

      if (!hasElement) {
        return section;
      }

      return {
        ...section,
        elements: swapPositions(section.elements, event.elementId, event.inc),
      };
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onMoveElement(
  aggregate: IFormAggregate,
  event: MoveElementEvent,
): IFormAggregate {
  let element: IElement | undefined;
  let sourceSection: ISection | undefined;

  outer: for (const page of aggregate.version.pages) {
    for (const section of page.sections) {
      const found = section.elements.find((e) => e.id === event.elementId);

      if (found) {
        element = found;
        sourceSection = section;
        break outer;
      }
    }
  }

  if (!element || !sourceSection) {
    return aggregate;
  }

  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSourceSection = page.sections.some(
      (s) => s.id === sourceSection!.id,
    );
    const hasTargetSection = page.sections.some(
      (s) => s.id === event.targetSectionId,
    );

    if (!hasSourceSection && !hasTargetSection) {
      return page;
    }

    const sections = page.sections.map((section): ISection => {
      if (
        section.id === sourceSection!.id &&
        section.id === event.targetSectionId
      ) {
        const without = removeById(section.elements, event.elementId);
        return {
          ...section,
          elements: insertAtPosition(without, element!, event.position),
        };
      }

      if (section.id === sourceSection!.id) {
        return {
          ...section,
          elements: removeById(section.elements, event.elementId),
        };
      }

      if (section.id === event.targetSectionId) {
        return {
          ...section,
          elements: insertAtPosition(
            section.elements,
            element!,
            event.position,
          ),
        };
      }

      return section;
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onUpdateElement(
  aggregate: IFormAggregate,
  event: UpdateElementEvent,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const sections = page.sections.map((section): ISection => {
      const hasElement = section.elements.some((e) => e.id === event.id);

      if (!hasElement) {
        return section;
      }

      const elements = section.elements.map((element) => {
        if (element.id !== event.id) {
          return element;
        }

        return {
          ...element,
          ...event.changes,
          attributes: event.changes.attributes
            ? { ...element.attributes, ...event.changes.attributes }
            : element.attributes,
        };
      });

      return { ...section, elements };
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onRemoveElement(
  aggregate: IFormAggregate,
  event: RemoveElementEvent,
): IFormAggregate {
  return removeElementById(aggregate, event.id);
}

export function onCutElement(
  aggregate: IFormAggregate,
  event: CutElementEvent,
): IFormAggregate {
  return removeElementById(aggregate, event.elementId);
}

function removeElementById(
  aggregate: IFormAggregate,
  elementId: string,
): IFormAggregate {
  const pages = aggregate.version.pages.map((page): IPage => {
    const sections = page.sections.map((section): ISection => {
      const hasElement = section.elements.some((e) => e.id === elementId);

      if (!hasElement) {
        return section;
      }

      return {
        ...section,
        elements: removeById(section.elements, elementId),
      };
    });

    const changed = sections.some((s, i) => s !== page.sections[i]);

    if (!changed) {
      return page;
    }

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}

export function onPasteElement(
  aggregate: IFormAggregate,
  event: PasteElementEvent,
): IFormAggregate {
  const perserve = event.clipboardOp === ClipboardEventType.CutElement;

  const element: IElement = {
    ...event.element,
    id: perserve ? event.element.id : generatedID(),
    key: perserve ? event.element.key : copyKey(event.element.key),
    name: perserve ? event.element.name : copyName(event.element.name),
  };

  const pages = aggregate.version.pages.map((page): IPage => {
    const hasSection = page.sections.some(
      (s) => s.id === event.targetSectionId,
    );

    if (!hasSection) {
      return page;
    }

    const sections = page.sections.map((section): ISection => {
      if (section.id !== event.targetSectionId) {
        return section;
      }

      return {
        ...section,
        elements: insertAtPosition(
          section.elements,
          element,
          getNextPosition(section.elements),
        ),
      };
    });

    return { ...page, sections };
  });

  return {
    ...aggregate,
    version: { ...aggregate.version, pages },
  };
}
