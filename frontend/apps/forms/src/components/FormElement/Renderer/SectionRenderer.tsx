import type { ISection } from "@/types/section";
import { sortPositioned } from "@/utils/position";
import { ElementRenderer } from "./ElementRenderer";
import { filterVisible } from "@/utils/filter";
import { useEvalContext } from "@/store/submission/evalContext";
import { AnimatePresence, motion } from "motion/react";
import Box from "@mui/material/Box";
import { rendererStyles } from "./renderer.style";
import { sectionVariants } from "./renderer.animations";

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const evalCtx = useEvalContext();
  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

  return (
    <Box
      component={motion.section}
      variants={sectionVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      sx={rendererStyles["section"]}
    >
      <AnimatePresence initial={false}>
        {elements.map((element) => (
          <ElementRenderer element={element} key={element.id} />
        ))}
      </AnimatePresence>
    </Box>
  );
};
