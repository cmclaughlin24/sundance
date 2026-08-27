import type { MouseEvent } from "react";
import { itemToolsStyles } from "./ItemTools.style";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import ArrowUpward from "@mui/icons-material/ArrowUpward";
import ArrowDownward from "@mui/icons-material/ArrowDownward";
import ContentCopy from "@mui/icons-material/ContentCopy";
import Delete from "@mui/icons-material/Delete";

export interface ItemToolsProps {
  onReorder: (inc: -1 | 1) => void;
  onCopy: () => void;
  onDelete: () => void;
}

export const ItemTools: React.FC<ItemToolsProps> = function ({
  onReorder,
  onCopy,
  onDelete,
}) {
  const handle =
    (action: () => void) => (event: MouseEvent<HTMLButtonElement>) => {
      event.stopPropagation();
      action();
    };

  return (
    <Box sx={itemToolsStyles.toolbar} data-testid="item-toolbar">
      <IconButton
        size="small"
        aria-label="Move up"
        data-testid="item-toolbar-move-up"
        sx={itemToolsStyles.button}
        onClick={handle(() => onReorder(-1))}
      >
        <ArrowUpward fontSize="inherit" />
      </IconButton>
      <IconButton
        size="small"
        aria-label="Move down"
        data-testid="item-toolbar-move-down"
        sx={itemToolsStyles.button}
        onClick={handle(() => onReorder(1))}
      >
        <ArrowDownward fontSize="inherit" />
      </IconButton>
      <IconButton
        size="small"
        aria-label="Copy"
        data-testid="item-toolbar-copy"
        sx={itemToolsStyles.button}
        onClick={handle(() => onCopy())}
      >
        <ContentCopy fontSize="inherit" />
      </IconButton>
      <IconButton
        size="small"
        aria-label="Delete"
        data-testid="item-toolbar-delete"
        sx={itemToolsStyles.deleteButton}
        onClick={handle(() => onDelete())}
      >
        <Delete fontSize="inherit" />
      </IconButton>
    </Box>
  );
};
