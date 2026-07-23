import type { IPage } from "@/types/page";
import { sortFormElements } from "@/utils/sort";
import { SectionRenderer } from "./SectionRenderer";

export const PageRenderer: React.FC<{ page: IPage }> = function ({ page }) {
  const sections = sortFormElements(page.sections);

  return sections.map((section) => (
    <SectionRenderer section={section} key={section.id} />
  ));
};
