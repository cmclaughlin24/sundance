import type { ISection } from "@/types/section";
import { sortFormElements } from "@/utils/sort";
import { ElementRenderer } from "./ElementRenderer";

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const elements = sortFormElements(section.elements);

  return elements.map((element) => <ElementRenderer element={element} key={element.id} />);
};
