import ExpandMore from "@mui/icons-material/ExpandMore";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion } from "motion/react";
import { useId, useState, type KeyboardEventHandler } from "react";
import { collapsibleStyles } from "./Collapsible.style";

export type CollapisbleProps = React.PropsWithChildren<{
  summary: string;
  defaultCollapsed?: boolean;
}>;

export const Collapisble: React.FC<CollapisbleProps> = function ({
  summary,
  children,
  defaultCollapsed,
}) {
  const [isCollapsed, setIsCollapsed] = useState(defaultCollapsed ?? false);
  const contentId = useId();

  const handleToggle = () => {
    setIsCollapsed((current) => !current);
  };

  const handleKeyDown: KeyboardEventHandler<HTMLDivElement> = (event) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      handleToggle();
    }
  };

  return (
    <Box sx={collapsibleStyles.collapisble}>
      <Box
        sx={collapsibleStyles.summary}
        onClick={handleToggle}
        onKeyDown={handleKeyDown}
        role="button"
        tabIndex={0}
        aria-expanded={!isCollapsed}
        aria-controls={contentId}
      >
        <ExpandMore
          fontSize="medium"
          component={motion.svg}
          animate={{ rotate: isCollapsed ? -90 : 0 }}
          transition={{ type: "spring", bounce: 0.6, duration: 0.4 }}
          aria-hidden="true"
        />
        <Typography>{summary}</Typography>
      </Box>
      <AnimatePresence initial={false}>
        {!isCollapsed && (
          <Box
            component={motion.div}
            id={contentId}
            key="content"
            initial={{ height: 0, opacity: 0 }}
            exit={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            transition={{ type: "spring", bounce: 0, duration: 0.4 }}
            sx={collapsibleStyles.content}
          >
            {children}
          </Box>
        )}
      </AnimatePresence>
    </Box>
  );
};
