import type { FormValues } from "@/store/submission/submissionStore";
import type { EvalContext } from "./evaluate";
import { evaluateRules } from "./evaluate";
import { filterVisible } from "./filter";
import type { IPage } from "@/types/page";

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
  values: FormValues,
  errors: Record<string, string[]>,
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
      const visibleElements = filterVisible(section.elements, evalCtx);

      for (const element of visibleElements) {
        const { required } = evaluateRules(element.rules, evalCtx, {
          required: element.attributes.isRequired,
        });

        if (!required) {
          continue;
        }

        total++;
        const hasErrors = errors[element.id] && errors[element.id].length > 0;

        if (!(element.id in values) || hasErrors) {
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
