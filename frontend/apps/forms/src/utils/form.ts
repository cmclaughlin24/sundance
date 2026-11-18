import type { IElement } from "@/types/element";
import type { FormVersionStatus, IFormVersion } from "@/types/formVersion";
import type {
  ISubmissionSection,
  ISubmissionSectionInstance,
  ISubmissionValue,
} from "@/types/submission";
import type { IFormProgress } from "./progress";
import type { IPage } from "@/types/page";
import type { SelectedItem } from "@/store/formDesigner";
import type {
  FormVersionRequest,
  RuleRequest,
} from "@/services/formService.type";
import type { IFlatRule, RuleParentType } from "@/types/rule";
import { stripID, stripTemporaryID } from "./id";
import type { ISection } from "@/types/section";
import * as ArrayUtils from "./array";
import { sortVersions } from "./version";

/**
 * Returns a flat array of all elements across every page and section of the given form version.
 *
 * @param version - The form version whose elements should be collected.
 * @returns A flat array of all `IElement` instances in the version.
 */
export function getVersionElements(version: IFormVersion): IElement[] {
  return version.pages.flatMap((page) =>
    page.sections.flatMap((section) => section.elements),
  );
}

/**
 * Returns all sections across every page of the given form version.
 *
 * @param version - The form version whose sections should be collected.
 * @returns A flat array of all `ISection` instances in the version.
 */
export function getVersionSections(version: IFormVersion): ISection[] {
  return version.pages.flatMap((page) => page.sections);
}

/**
 * Builds the starting submission for the given form version, grouped by section. Saved instances
 * are kept in order, and any element missing from an instance gets its default value. A section
 * with no saved data gets one instance holding only default values.
 *
 * Saved sections that don't exist in the version are dropped.
 *
 * @param raw - The prior submission to rehydrate from, if any.
 * @param version - The form version whose sections and elements define the expected shape.
 * @returns One `ISubmissionSection` per section in the version.
 */
export function hydrateSubmissionValues(
  raw: ISubmissionSection[] | undefined,
  version: IFormVersion,
): ISubmissionSection[] {
  const saved = new Map(raw?.map((s) => [s.sectionId, s.instances]) ?? []);

  return getVersionSections(version).map((section) => {
    const instances = saved.get(section.id);

    if (!instances || instances.length === 0) {
      return {
        sectionId: section.id,
        instances: [hydrateInstance(undefined, section.elements)],
      };
    }

    return {
      sectionId: section.id,
      instances: instances.map((inst) =>
        hydrateInstance(inst, section.elements),
      ),
    };
  });
}

/**
 * Keeps the saved values for an instance and adds default values for any element that has none.
 */
function hydrateInstance(
  instance: ISubmissionSectionInstance | undefined,
  elements: IElement[],
): ISubmissionSectionInstance {
  const saved = instance?.values ?? [];
  const answered = new Set(saved.map((v) => v.elementId));
  const values: ISubmissionValue[] = [...saved];

  for (const element of elements) {
    if (answered.has(element.id)) {
      continue;
    }

    const defaultValue = element.attributes.defaultValue;

    if (defaultValue != null) {
      values.push(...toSubmissionValues(element.id, defaultValue));
    }
  }

  return { values };
}

/**
 * Converts one element's value into API entries: an array becomes one entry per item, each
 * with its `itemIndex`; any other value becomes a single entry.
 */
function toSubmissionValues(elementId: string, value: any): ISubmissionValue[] {
  if (Array.isArray(value)) {
    return value.map((v, itemIndex) => ({ elementId, value: v, itemIndex }));
  }

  return [{ elementId, value }];
}

/**
 * Checks the validity of the given form element, returning true if it is valid or false if it is invalid.
 * @param formEl - The HTML form element to check.
 * @param progress - The current form element to check.
 * @returns `true` if the form is valid, `false` otherwise.
 */
export function checkFormValidity(
  formEl: HTMLFormElement | null,
  progress: IFormProgress | undefined,
): boolean {
  if (!formEl || !progress) {
    return false;
  }

  return (
    formEl.checkValidity() &&
    progress.errors === 0 &&
    progress.percentage === 100
  );
}

/**
 * Finds a selected item by its ID within the given pages.
 * @param pages - The array of pages to search through.
 * @param id - The ID of the item to find.
 * @returns The selected item if found, or `null` if not found.
 */
export function findSelectedById(
  pages: IPage[],
  id: string,
): SelectedItem | null {
  for (const page of pages) {
    if (page.id === id) {
      return { type: "page", item: page };
    }

    for (const section of page.sections) {
      if (section.id === id) {
        return { type: "section", item: section };
      }

      for (const element of section.elements) {
        if (element.id === id) {
          return { type: "element", item: element };
        }
      }
    }
  }

  return null;
}

export function versionToRequest(
  version: IFormVersion,
  rules: IFlatRule[],
): FormVersionRequest {
  const pages = version.pages.map((page) => ({
    ...stripTemporaryID(page),
    sections: page.sections.map((section) => ({
      ...stripTemporaryID(section),
      elements: section.elements.map((element) => ({
        ...stripTemporaryID(element),
        rules: setRules("element", element.id, rules),
      })),
      rules: setRules("section", section.id, rules),
    })),
    rules: setRules("page", page.id, rules),
  }));

  return { metadata: {}, pages };
}

export function copyVersion(version: IFormVersion): FormVersionRequest {
  const pages = version.pages.map((page) => ({
    ...stripID(page),
    sections: page.sections.map((section) => ({
      ...stripID(section),
      elements: section.elements.map((element) => stripID(element)),
    })),
  }));

  return { metadata: {}, pages };
}

function setRules(
  parentType: RuleParentType,
  parentId: string,
  rules: IFlatRule[],
): RuleRequest[] {
  const toRule = ({
    parentId: _parentId,
    parentType: _parentType,
    ...rule
  }: IFlatRule) => stripTemporaryID(rule);

  return rules
    .filter((r) => r.parentType === parentType && r.parentId === parentId)
    .map(toRule);
}

export function isDraftVersion(status: FormVersionStatus): boolean {
  return status === "draft";
}

export function isActiveVersion(status: FormVersionStatus): boolean {
  return status === "active";
}

export function groupFormObjects(pages: IPage[]) {
  const objects: {
    pages: IPage[];
    sections: ISection[];
    elements: IElement[];
  } = {
    pages: [],
    sections: [],
    elements: [],
  };

  for (const page of pages) {
    for (const section of page.sections) {
      for (const element of section.elements) {
        objects.elements.push(element);
      }
      objects.sections.push(section);
    }
    objects.pages.push(page);
  }

  return objects;
}

/**
 * Returns a flat array of all elements across every page and section of the given pages.
 *
 * @param pages - The pages whose elements should be collected.
 * @returns A flat array of all `IElement` instances.
 */
export function getFlattenedElements(pages: IPage[]): IElement[] {
  return pages.flatMap((page) =>
    (page.sections ?? []).flatMap((section) => section.elements ?? []),
  );
}

/**
 * Returns a default form version object with empty metadata and pages.
 * @returns A default `FormVersionRequest` objet with empty metadata and pages.
 */
export function defaultFormVersion(): FormVersionRequest {
  return {
    metadata: {},
    pages: [],
  };
}

export function getLatestFormVersionByStatus(
  versions: IFormVersion[] | null,
  status: FormVersionStatus,
): IFormVersion | undefined {
  if (!ArrayUtils.hasLengthGreaterThan(versions, 0)) {
    return undefined;
  }

  return sortVersions(versions!).find((v) => v.status === status);
}
