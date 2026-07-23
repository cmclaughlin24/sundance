import type { IPage } from "@/types/page";
import { sortPositioned } from "@/utils/sort";
import { SectionRenderer } from "./SectionRenderer";
import { filterVisible } from "@/utils/filter";
import { useEvalContext } from "@/store/evalContext";

export const PageRenderer: React.FC<{ page: IPage }> = function ({ page }) {
  const evalCtx = useEvalContext();
  let sections = sortPositioned(page.sections);
  sections = filterVisible(sections, evalCtx);

  return sections.map((section) => (
    <SectionRenderer section={section} key={section.id} />
  ));
};
