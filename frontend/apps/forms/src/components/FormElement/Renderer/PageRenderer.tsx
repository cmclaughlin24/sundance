import type { IPage } from "@/types/page";
import { sortPositioned } from "@/utils/sort";
import { SectionRenderer } from "./SectionRenderer";

export const PageRenderer: React.FC<{ page: IPage }> = function ({ page }) {
  const sections = sortPositioned(page.sections);

  return sections.map((section) => (
    <SectionRenderer section={section} key={section.id} />
  ));
};
