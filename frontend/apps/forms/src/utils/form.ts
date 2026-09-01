import type { IElement } from "@/types/element";
import type { IFormVersion } from "@/types/formVersion";
import type { ISubmissionValue } from "@/types/submission";
import type { IFormProgress } from "./progress";
import type { IPage } from "@/types/page";
import type { SelectedItem } from "@/store/formDesigner";

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
 * Builds a complete set of submission values for the given form version, using prior submission
 * values where available and falling back to each element's default value for any gaps.
 *
 * @param raw - The prior submission values to rehydrate from, if any.
 * @param version - The form version whose elements define the expected shape.
 * @returns A complete array of `ISubmissionValue` entries, one per element in the version.
 */
export function hydrateSubmissionValues(
  raw: ISubmissionValue[] | undefined,
  version: IFormVersion,
): ISubmissionValue[] {
  const values: ISubmissionValue[] = [];
  const elements = getVersionElements(version);

  for (const element of elements) {
    const rawValue = raw?.find((val) => val.elementId === element.id);

    if (rawValue) {
      values.push(rawValue);
      continue;
    }

    const defaultValue = element.attributes.defaultValue;

    if (defaultValue != null) {
      values.push({ elementId: element.id, value: defaultValue });
    }
  }

  return values;
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
