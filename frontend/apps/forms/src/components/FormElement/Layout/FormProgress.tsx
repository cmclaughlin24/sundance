import type { IFormProgress } from "@/utils/progress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion, type Variants } from "motion/react";
import { formProgressStyles } from "./FormProgress.style";

const variants: Variants = {
  initial: {
    opacity: 0,
  },
  animate: {
    opacity: 1,
    transition: { duration: 0.2 },
  },
  exit: {
    opacity: 0,
  },
};

export const FormProgress: React.FC<{ progress: IFormProgress }> = function ({
  progress,
}) {
  let statusMessage = (
    <Typography
      component={motion.span}
      key="remaining"
      variants={variants}
      initial="initial"
      animate="animate"
      exit="exit"
    >
      {progress.total - progress.filled} required field(s) reamining
    </Typography>
  );

  if (progress.errors) {
    statusMessage = (
      <Typography
        component={motion.span}
        key="remaining"
        sx={formProgressStyles["error"]}
        variants={variants}
        initial="initial"
        animate="animate"
        exit="exit"
      >
        {progress.errors} field(s) require attention
      </Typography>
    );
  }

  return (
    <Box sx={formProgressStyles["container"]}>
      <Typography variant="body2" sx={formProgressStyles["required"]}>
        Required Information
      </Typography>
      <Typography variant="body1" sx={formProgressStyles['text']}>
        {progress.filled}/{progress.total} completed ·{" "}
        <AnimatePresence mode="wait">{statusMessage}</AnimatePresence>
      </Typography>
    </Box>
  );
};
