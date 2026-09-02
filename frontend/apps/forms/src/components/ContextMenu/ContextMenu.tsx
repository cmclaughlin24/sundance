import { createPortal } from "react-dom";
import { ContextMenuButton } from "./ContextMenuButton";
import { useContextMenuDispatch, useContextMenuTarget } from "./useContextMenu";
import { motion } from "motion/react";
import { contextMenuStyles } from "./ContextMenu.styles";
import Box from "@mui/material/Box";

export interface ContextMenuProps {
  container: Element | DocumentFragment;
  width?: string | number;
  children: (data: unknown) => React.ReactNode | undefined;
}

interface ContextMenuComponent extends React.FC<ContextMenuProps> {
  Button: typeof ContextMenuButton;
}

const ContextMenu: ContextMenuComponent = function ({
  container,
  children,
  width = "244px",
}) {
  const { close } = useContextMenuDispatch();
  const target = useContextMenuTarget();

  const handleBackdropClk = (event: React.MouseEvent) => {
    event.stopPropagation();
    close();
  };

  if (!target) {
    return null;
  }

  // TODO: Handle offset for screen edges.
  const x = target.position.x;
  const y = target.position.y;

  return createPortal(
    <Box
      component={motion.div}
      sx={contextMenuStyles.backdrop}
      onClick={handleBackdropClk}
    >
      <Box sx={{ ...contextMenuStyles.contextMenu, width, top: y, left: x }}>
        {children(target.data)}
      </Box>
    </Box>,
    container,
  );
};

ContextMenu.Button = ContextMenuButton;

export default ContextMenu;
