import type { EvalContext } from "./evaluate";
import { evaluateRules } from "./evaluate";
import { filterVisible } from "./filter";
import type { IPage } from "@/types/page";
import type { SectionInstanceState } from "@/store/submission/SectionInstanceState";
import * as ArrayUtils from "@/utils/array";

export interface IFormProgress {
  /**
   * The percentage of the form that has been filled out, calculated as (filled / total ) * 100.
   */
  percentage: number;

  /**
   * The number of required fields that have been filled out.
   */
  filled: number;

  /**
   * The total number oof required fields in the form.
   */
  total: number;

  /**
   * The number of fields that have errors, which can be used to provide feedback to the user
   * about the form's completion.
   */
  errors: number;
}

export function calculateProgress(
  formValues: Record<string, SectionInstanceState[]>,
  evalCtx: EvalContext,
  pages: IPage[],
): IFormProgress | undefined {
  if (!pages || pages.length === 0) {
    return undefined;
  }

  let total = 0;
  let filled = 0;
  let errorCount = 0;
  const visiblePages = filterVisible(pages, evalCtx);

  for (const page of visiblePages) {
    const visibleSections = filterVisible(page.sections, evalCtx);

    for (const section of visibleSections) {
      const instances = formValues[section.id];

      const requiredElements = filterVisible(section.elements, evalCtx).filter(
        (element) =>
          evaluateRules(element.rules, evalCtx, {
            required: element.attributes.isRequired,
          }).required,
      );

      if (!ArrayUtils.hasLengthGreaterThan(instances, 0)) {
        total += requiredElements.length;
        continue;
      }

      // NOTE: Each instance must fill in its own required fields, so they are counted once per instance.
      for (const instance of instances) {
        for (const element of requiredElements) {
          total++;

          if (
            ArrayUtils.hasLengthGreaterThan(instance.getErrors(element.id)!, 0)
          ) {
            errorCount++;
            continue;
          }

          if (isEmpty(instance.getValue(element.id))) {
            continue;
          }

          filled++;
        }
      }
    }
  }

  return {
    percentage: total === 0 ? 100 : (filled / total) * 100,
    filled: filled,
    total: total,
    errors: errorCount,
  };
}

function isEmpty(value: any): boolean {
  return (
    value == null ||
    value === "" ||
    (Array.isArray(value) && value.length === 0)
  );
}
