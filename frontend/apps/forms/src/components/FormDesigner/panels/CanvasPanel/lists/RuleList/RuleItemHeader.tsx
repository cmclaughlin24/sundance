import ExpandMore from "@mui/icons-material/ExpandMore";
import ContentCopy from "@mui/icons-material/ContentCopy";
import Delete from "@mui/icons-material/Delete";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { motion } from "motion/react";
import { useState, type KeyboardEventHandler, type MouseEvent } from "react";
import { ruleItemHeaderStyles } from "./RuleItemHeader.style";
import { Tag } from "@/components/Tag";
import Tooltip from "@mui/material/Tooltip";
import IconButton from "@mui/material/IconButton";
import Snackbar from "@mui/material/Snackbar";

export interface RuleItemHeaderProps {
  id: string;
  isCollapsed: boolean;
  onCollapse: (isCollapsed: boolean) => void;
  onCopy: () => void;
  onDelete: () => void;
}

export const RuleItemHeader: React.FC<RuleItemHeaderProps> = function ({
  id,
  isCollapsed,
  onCollapse,
  onCopy,
  onDelete,
}) {
  const [isSnackbarOpen, setIsSnackbarOpen] = useState(false);

  const handle =
    (action: () => void) => (event: MouseEvent<HTMLButtonElement>) => {
      event.stopPropagation();
      action();
    };

  const handleToggle = () => onCollapse(!isCollapsed);

  const handleKeyDown: KeyboardEventHandler<HTMLDivElement> = (event) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      handleToggle();
    }
  };

  return (
    <>
      <Box
        sx={ruleItemHeaderStyles.ruleItemHeader}
        data-testid="rule-item-header"
      >
        <Box sx={{ display: "flex", gap: 1 }}>
          <Box
            sx={ruleItemHeaderStyles.toggle}
            onClick={handleToggle}
            onKeyDown={handleKeyDown}
            role="button"
            tabIndex={0}
            aria-expanded={!isCollapsed}
            aria-controls={id}
          >
            <ExpandMore
              fontSize="medium"
              component={motion.svg}
              animate={{ rotate: isCollapsed ? -90 : 0 }}
              transition={{ type: "spring", bounce: 0.6, duration: 0.4 }}
              aria-hidden="true"
            />
            <Typography>Rule</Typography>
          </Box>
          <Tag sx={{ alignSelf: "start" }}>Required</Tag>
        </Box>
        <Box>
          <Tooltip title="Copy">
            <IconButton
              size="small"
              aria-label="Copy"
              data-testid="item-toolbar-copy"
              sx={ruleItemHeaderStyles.button}
              onClick={handle(() => {
                onCopy();
                setIsSnackbarOpen(true);
              })}
            >
              <ContentCopy fontSize="inherit" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete">
            <IconButton
              size="small"
              aria-label="Delete"
              data-testid="item-toolbar-delete"
              sx={ruleItemHeaderStyles.deleteButton}
              onClick={handle(() => onDelete())}
            >
              <Delete fontSize="inherit" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      <Snackbar
        open={isSnackbarOpen}
        onClose={() => setIsSnackbarOpen(false)}
        message="Copied to Clipboard!"
        autoHideDuration={2500}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      />
    </>
  );
};
