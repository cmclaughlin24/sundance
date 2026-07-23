import type { IPage } from "@/types/page";
import { sortPositioned } from "@/utils/sort";
import { SectionRenderer } from "./SectionRenderer";
import { useFormState } from "@/store/useFormContext";
import { filterVisible } from "@/utils/filter";

export const PageRenderer: React.FC<{ page: IPage }> = function ({ page }) {
  const state = useFormState();
  let sections = sortPositioned(page.sections);
  sections = filterVisible(sections, state);

  return sections.map((section) => (
    <SectionRenderer section={section} key={section.id} />
  ));
};
