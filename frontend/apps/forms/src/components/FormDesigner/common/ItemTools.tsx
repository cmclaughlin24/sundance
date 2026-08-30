import { useState, type MouseEvent } from "react";
import { itemToolsStyles } from "./ItemTools.style";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import ArrowUpward from "@mui/icons-material/ArrowUpward";
import ArrowDownward from "@mui/icons-material/ArrowDownward";
import ContentCopy from "@mui/icons-material/ContentCopy";
import Delete from "@mui/icons-material/Delete";
import Tooltip from "@mui/material/Tooltip";
import Snackbar from "@mui/material/Snackbar";

export interface ItemToolsConfig {
  tooltip: {
    moveUp?: string;
    moveDown?: string;
    copy?: string;
    delete?: string;
  };
  copy: {
    enable?: boolean;
    message?: string;
  };
}

const defaultConfig: ItemToolsConfig = {
  tooltip: {
    moveUp: "Move Up",
    moveDown: "Move Down",
    copy: "Copy",
    delete: "Delete",
  },
  copy: {
    enable: true,
    message: "Copied to Clipboard!",
  },
};

export interface ItemToolsProps {
  onReorder: (inc: -1 | 1) => void;
  onCopy: () => void;
  onDelete: () => void;
  config?: ItemToolsConfig;
}

export const ItemTools: React.FC<ItemToolsProps> = function ({
  onReorder,
  onCopy,
  onDelete,
  config: propsConfig = defaultConfig,
}) {
  const [isSnackbarOpen, setIsSnackbarOpen] = useState(false);

  const config: ItemToolsConfig = {
    tooltip: { ...defaultConfig.tooltip, ...propsConfig.tooltip },
    copy: { ...defaultConfig.copy, ...propsConfig.copy },
  };

  const handle =
    (action: () => void) => (event: MouseEvent<HTMLButtonElement>) => {
      event.stopPropagation();
      action();
    };

  return (
    <>
      <Box sx={itemToolsStyles.toolbar} data-testid="item-toolbar">
        <Tooltip title={config.tooltip.moveUp}>
          <IconButton
            size="small"
            aria-label="Move up"
            data-testid="item-toolbar-move-up"
            sx={itemToolsStyles.button}
            onClick={handle(() => onReorder(-1))}
          >
            <ArrowUpward fontSize="inherit" />
          </IconButton>
        </Tooltip>
        <Tooltip title={config.tooltip.moveDown}>
          <IconButton
            size="small"
            aria-label="Move down"
            data-testid="item-toolbar-move-down"
            sx={itemToolsStyles.button}
            onClick={handle(() => onReorder(1))}
          >
            <ArrowDownward fontSize="inherit" />
          </IconButton>
        </Tooltip>
        <Tooltip title={config.tooltip.copy}>
          <IconButton
            size="small"
            aria-label="Copy"
            data-testid="item-toolbar-copy"
            sx={itemToolsStyles.button}
            onClick={handle(() => {
              onCopy();
              setIsSnackbarOpen(true);
            })}
          >
            <ContentCopy fontSize="inherit" />
          </IconButton>
        </Tooltip>
        <Tooltip title={config.tooltip.delete}>
          <IconButton
            size="small"
            aria-label="Delete"
            data-testid="item-toolbar-delete"
            sx={itemToolsStyles.deleteButton}
            onClick={handle(() => onDelete())}
          >
            <Delete fontSize="inherit" />
          </IconButton>
        </Tooltip>
      </Box>
      {config.copy.enable && (
        <Snackbar
          open={isSnackbarOpen}
          onClose={() => setIsSnackbarOpen(false)}
          message={config.copy.message}
          autoHideDuration={2500}
          anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
        />
      )}
    </>
  );
};
