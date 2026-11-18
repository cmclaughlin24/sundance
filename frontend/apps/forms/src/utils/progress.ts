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

      if (!ArrayUtils.hasLengthGreaterThan(instances, 0)) {
        continue;
      }

      const visibleElements = filterVisible(section.elements, evalCtx);

      for (const element of visibleElements) {
        const { required } = evaluateRules(element.rules, evalCtx, {
          required: element.attributes.isRequired,
        });

        if (!required) {
          continue;
        }

        total++;

        const hasErrors = instances.some((instance) =>
          ArrayUtils.hasLengthGreaterThan(instance.getErrors(element.id)!, 0),
        );

        const hasValue = instances.some((instance) => {
          const value = instance.getValue(element.id);
          return value != undefined;
        });

        if (!hasValue || hasErrors) {
          hasErrors && errorCount++;
          continue;
        }

        filled++;
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
