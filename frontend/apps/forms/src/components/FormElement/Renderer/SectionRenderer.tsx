import type { ISection } from "@/types/section";
import { sortPositioned } from "@/utils/sort";
import { ElementRenderer } from "./ElementRenderer";
import { filterVisible } from "@/utils/filter";
import { useEvalContext } from "@/store/evalContext";

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const evalCtx = useEvalContext();
  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

  return elements.map((element) => (
    <ElementRenderer element={element} key={element.id} />
  ));
};
