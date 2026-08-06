import type { IElement } from "@/types/element";
import type { IFormVersion } from "@/types/formVersion";
import type { ISubmissionValue } from "@/types/submission";

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

    values.push({
      elementId: element.id,
      value: element.attributes.defaultValue,
    });
  }

  return values;
}
