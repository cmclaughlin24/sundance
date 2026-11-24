import type { ISection, ISectionRepeat } from "@/types/section";
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

function hasMinInstances(
  count: number,
  repeat: ISectionRepeat | undefined,
): boolean {
  if (!repeat || !repeat.minInstances) {
    return count > 0;
  }

  return count >= repeat.minInstances;
}

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const { addInstance, ensureMinInstances, removeInstance } =
    useSubmissionInstanceActions();
  const evalCtx = useEvalContext();
  const instances = useSubmissionSectionInstances(section.id);
  const repeat = section.repeat;
  const hasInstances = instances && hasMinInstances(instances.length, repeat);

  useEffect(() => {
    if (!hasInstances) {
      ensureMinInstances(section.id, repeat?.minInstances ?? 0);
    }
  }, [ensureMinInstances, section.id, hasInstances, repeat]);

  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

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
