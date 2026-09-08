import type { IPage } from "@/types/page";
import type { IFlatRule } from "@/types/rule";
import * as ArrayUtils from "./array";
import type { HasID } from "@/types/formObject";

/**
 * Walks the full page/section/element hierarchy of a form version and returns a flat array
 * of all rules, each annotated with the ID and type of its owner.
 * @param pages - The pages to extract rules from.
 * @returns A flat array of `IFlatRule` instances.
 */
export function extractFlatRules(pages: IPage[]): IFlatRule[] {
  const rules: IFlatRule[] = [];

  for (const page of pages) {
    for (const rule of page.rules) {
      rules.push({ ...rule, parentId: page.id, parentType: "page" });
    }

    for (const section of page.sections) {
      for (const rule of section.rules) {
        rules.push({ ...rule, parentId: section.id, parentType: "section" });
      }

      for (const element of section.elements) {
        for (const rule of element.rules) {
          rules.push({ ...rule, parentId: element.id, parentType: "element" });
        }
      }
    }
  }

  return rules;
}

export function removeFlatRulesByIDs(
  rules: IFlatRule[],
  ...ids: (string | HasID)[]
): IFlatRule[] {
  if (!ArrayUtils.hasLengthGreaterThan(rules, 0)) {
    return [];
  }

  const targetIds = new Set(
    ids.map((id) => (typeof id === "string" ? id : id.id)),
  );

  return rules.filter((rule) => rule.parentId && !targetIds.has(rule.parentId));
}
