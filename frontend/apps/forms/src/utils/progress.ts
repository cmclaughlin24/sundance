import type { FormValues } from "@/store/submission/submissionStore";
import type { EvalContext } from "./evaluate";
import { evaluateRules } from "./evaluate";
import { filterVisible } from "./filter";
import type { IPage } from "@/types/page";

export function calculateProgress(
  values: FormValues,
  errors: Record<string, string[]>,
  evalCtx: EvalContext,
  pages: IPage[],
): number {
  if (!pages || pages.length === 0) {
    return 100;
  }

  let total = 0;
  let filled = 0;
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

        if (
          !(element.id in values) ||
          (errors[element.id] && errors[element.id].length > 0)
        ) {
          continue;
        }

        filled++;
      }
    }
  }

  return total === 0 ? 100 : (filled / total) * 100;
}
