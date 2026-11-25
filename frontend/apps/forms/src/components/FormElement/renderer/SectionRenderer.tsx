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
import * as ArrayUtils from "@/utils/array";

function requiredInstances(repeat: ISectionRepeat | undefined): number {
  return repeat ? (repeat.minInstances ?? 0) : 1;
}

export const SectionRenderer: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const { addInstance, ensureMinInstances, removeInstance } =
    useSubmissionInstanceActions();
  const evalCtx = useEvalContext();
  const instances = useSubmissionSectionInstances(section.id);
  const repeat = section.repeat;
  const min = requiredInstances(repeat);
  const needsInstances = instances.length < min;

  useEffect(() => {
    if (!needsInstances) {
      ensureMinInstances(section.id, repeat?.minInstances ?? min);
    }
  }, [ensureMinInstances, section.id, needsInstances, repeat]);

  let elements = sortPositioned(section.elements);
  elements = filterVisible(elements, evalCtx);

  const hasInstances = ArrayUtils.hasLengthGreaterThan(instances, 0);
  const instanceLabel = repeat?.instanceLabel ?? "Item";

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
      {repeat && !hasInstances && (
        <Box sx={rendererStyles["repeatableSectionNoInstances"]}>
          <Typography sx={{ textAlign: "center", fontSize: "1rem" }}>
            No {instanceLabel} added.
          </Typography>
        </Box>
      )}
      <AnimatePresence initial={false}>
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
                    {instanceLabel} {idx + 1}
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
      </AnimatePresence>
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
