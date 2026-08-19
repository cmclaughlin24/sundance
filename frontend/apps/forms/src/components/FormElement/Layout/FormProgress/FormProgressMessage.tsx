import type { IFormProgress } from "@/utils/progress";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion, type Variants } from "motion/react";
import { formProgressStyles } from "./FormProgress.style";
import { useTheme } from "@mui/material/styles";

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

export interface FormProgressMessageProps {
  progress: IFormProgress;

  variant?: "errorsOnly" | "remainWithErrors";
}

export const FormProgressMessage: React.FC<FormProgressMessageProps> =
  function ({ variant, progress }) {
    const theme = useTheme();

    let statusMessage =
      variant && variant === "remainWithErrors" ? (
        <Typography
          component={motion.span}
          key="remaining"
          variants={variants}
          initial="initial"
          animate="animate"
          exit="exit"
        >
          {progress!.total - progress!.filled} required field(s) reamining
        </Typography>
      ) : null;

    const handleErrorClk = () => {
      const element = document.querySelector('[data-error="true"]');

      if (!element) {
        return;
      }

      element.scrollIntoView({ behavior: "smooth", block: "center" });
      element.animate(
        [
          {
            background: `${theme.palette.error.light}15`,
            borderRadius: "16px",
            boxShadow: `0 0 0 5px ${theme.palette.error.main}80`,
          },
          {
            background: `${theme.palette.error.light}00`,
            borderRadius: "16px",
            boxShadow: `0 0 0 13px ${theme.palette.error.main}00`,
          },
        ],
        { duration: 1200, easing: "ease-out" },
      );
    };

    if (progress!.errors) {
      statusMessage = (
        <Typography
          component={motion.span}
          key="remaining"
          sx={formProgressStyles["error"]}
          variants={variants}
          initial="initial"
          animate="animate"
          onClick={handleErrorClk}
          role="button"
          exit="exit"
        >
          {progress!.errors} field(s) require attention
        </Typography>
      );
    }

    return <AnimatePresence mode="wait">{statusMessage}</AnimatePresence>;
  };
