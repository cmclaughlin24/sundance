import type { IPage } from "@/types/page";
import { sortPositioned } from "@/utils/position";
import { SectionRenderer } from "./SectionRenderer";
import { filterVisible } from "@/utils/filter";
import { useEvalContext } from "@/store/submission/evalContext";
import Box from "@mui/material/Box";
import { AnimatePresence, motion } from "motion/react";
import { pageVariants } from "./renderer.animations";
import { rendererStyles } from "./renderer.style";

export const PageRenderer: React.FC<{ page: IPage }> = function ({ page }) {
  const evalCtx = useEvalContext();
  let sections = sortPositioned(page.sections);
  sections = filterVisible(sections, evalCtx);

  return (
    <Box
      component={motion.div}
      variants={pageVariants}
      initial="initial"
      animate="animate"
      sx={rendererStyles["page"]}
    >
      <AnimatePresence initial={false}>
        {sections.map((section) => (
          <SectionRenderer section={section} key={section.id} />
        ))}
      </AnimatePresence>
    </Box>
  );
};
