import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import { useTheme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion, type Variants } from "motion/react";
import { useMemo } from "react";

const styles: Styles = {
  dropZone: (theme) => ({
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    border: `1px dashed ${theme.palette.primary.main}`,
    borderRadius: 2.5,
    background: `${theme.palette.primary.main}20`,
    px: 3,
    height: "4.25rem",
  }),
  text: (theme) => ({
    fontWeight: 600,
    color: theme.palette.primary.main,
  }),
};

export const DropZoneIndicator: React.FC<{
  text: string;
  isVisible?: boolean;
  isDropTarget?: boolean;
}> = function ({ text, isVisible = false, isDropTarget = false }) {
  const theme = useTheme();
  const variants: Variants = useMemo(
    () => ({
      initial: {
        background: `${theme.palette.primary.main}20`,
      },
      animate: {
        background: `${theme.palette.primary.main}60`,
        transition: { duration: 0.25, ease: "easeInOut" },
      },
      exit: {
        background: `${theme.palette.primary.main}20`,
      },
    }),
    [theme],
  );

  return (
    <AnimatePresence>
      {isVisible && (
        <Box
          component={motion.div}
          sx={styles.dropZone}
          variants={variants}
          initial="initial"
          animate={isDropTarget ? "animate" : "initial"}
          exit="exit"
        >
          <Typography sx={styles.text}>{text}</Typography>
        </Box>
      )}
    </AnimatePresence>
  );
};
