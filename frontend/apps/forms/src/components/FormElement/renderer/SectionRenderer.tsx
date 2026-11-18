import type { ISection } from "@/types/section";
import { sortPositioned } from "@/utils/position";
import { ElementRenderer } from "./ElementRenderer";
import { filterVisible } from "@/utils/filter";
import { useEvalContext } from "@/store/submission/evalContext";
import { AnimatePresence, motion } from "motion/react";
import Box from "@mui/material/Box";
import { rendererStyles } from "./renderer.style";
import { sectionVariants } from "./renderer.animations";
import {
  useSubmissionInstanceActions,
  useSubmissionSectionInstances,
} from "@/store/submission/useSubmissionContext";
import { useEffect } from "react";
import * as ArrayUtils from "@/utils/array";
import { SectionInstanceProvider } from "@/store/submission/sectionInstanceContext";

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const evalCtx = useEvalContext();
  const instances = useSubmissionSectionInstances(section.id);
  const { ensureInstance } = useSubmissionInstanceActions();

  useEffect(() => {
    if (!ArrayUtils.hasLengthGreaterThan(instances, 0)) {
      ensureInstance(section.id);
    }
  }, [ensureInstance, section.id, instances.length]);

  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

  if (!ArrayUtils.hasLengthGreaterThan(instances, 0)) {
    return null;
  }

  return (
    <>
      {instances.map((i) => (
        <SectionInstanceProvider
          sectionId={section.id}
          instanceId={i.id}
          key={i.id}
        >
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
        </SectionInstanceProvider>
      ))}
    </>
  );
};
