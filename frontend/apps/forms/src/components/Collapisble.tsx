import ExpandMore from "@mui/icons-material/ExpandMore";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion } from "motion/react";
import { useState } from "react";
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

  const handleToggle = () => {
    setIsCollapsed((current) => !current);
  };

  return (
    <Box sx={collapsibleStyles.collapisble}>
      <Box sx={collapsibleStyles.summary} onClick={handleToggle} role="button">
        <ExpandMore
          fontSize="medium"
          component={motion.svg}
          animate={{ rotate: isCollapsed ? -90 : 0 }}
          transition={{ type: "spring", bounce: 0.6, duration: 0.4 }}
        />
        <Typography>{summary}</Typography>
      </Box>
      <AnimatePresence initial={false}>
        {!isCollapsed && (
          <Box
            component={motion.div}
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
