import { createPortal } from "react-dom";
import { ContextMenuButton } from "./ContextMenuButton";
import { useContextMenuDispatch, useContextMenuTarget } from "./useContextMenu";
import { AnimatePresence, motion } from "motion/react";
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

  // TODO: Handle offset for screen edges.
  const x = target?.position.x;
  const y = target?.position.y;

  return createPortal(
    <AnimatePresence>
      {target && (
        <Box
          component={motion.div}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          sx={contextMenuStyles.backdrop}
          onClick={handleBackdropClk}
        >
          <Box
            sx={{ ...contextMenuStyles.contextMenu, width, top: y, left: x }}
          >
            {children(target.data)}
          </Box>
        </Box>
      )}
    </AnimatePresence>,
    container,
  );
};

ContextMenu.Button = ContextMenuButton;

export default ContextMenu;
