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
import { SectionInstanceProvider } from "@/store/submission/sectionInstanceContext";
import Button from "@mui/material/Button";
import { mergeSx } from "merge-sx";
import Typography from "@mui/material/Typography";

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const evalCtx = useEvalContext();
  const instances = useSubmissionSectionInstances(section.id);
  const hasInstances = instances && instances.length > 0;
  const { addInstance, ensureInstance, removeInstance } =
    useSubmissionInstanceActions();
  const repeat = section.repeat;

  useEffect(() => {
    if (!hasInstances) {
      ensureInstance(section.id);
    }
  }, [ensureInstance, section.id, hasInstances]);

  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

  if (!hasInstances) {
    return null;
  }

  return (
    <>
      {repeat && (
        <Box sx={rendererStyles["repeatableSectionActionTitle"]}>
          <Typography component="h2" sx={{ fontSize: "2.5rem" }}>
            {section.name}
          </Typography>
          {repeat.maxInstances && (
            <Typography sx={{ fontWeight: 600 }}>
              {instances.length} of {repeat.maxInstances}
            </Typography>
          )}
        </Box>
      )}
      {instances.map((i, idx) => (
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
            sx={mergeSx(
              rendererStyles["section"],
              repeat && rendererStyles["repeatableSection"],
            )}
          >
            {repeat && (
              <Box sx={rendererStyles["repeatableSectionActionTitle"]}>
                <Typography
                  component="h3"
                  sx={{ fontSize: "1.25rem", fontWeight: 600 }}
                >
                  {repeat.instanceLabel ?? "Item"} {idx + 1}
                </Typography>
                <Button
                  variant="contained"
                  color="secondary"
                  size="small"
                  onClick={() => removeInstance(section.id, i.id)}
                  disabled={
                    repeat.minInstances != null && idx < repeat.minInstances
                  }
                >
                  Remove
                </Button>
              </Box>
            )}
            <AnimatePresence initial={false}>
              {elements.map((element) => (
                <ElementRenderer element={element} key={element.id} />
              ))}
            </AnimatePresence>
          </Box>
        </SectionInstanceProvider>
      ))}
      {repeat && (
        <Button
          variant="outlined"
          onClick={() => addInstance(section.id)}
          disabled={
            repeat.maxInstances != null &&
            instances.length === repeat.maxInstances
          }
        >
          {repeat.addLabel ?? "Add Instance"}
        </Button>
      )}
    </>
  );
};
