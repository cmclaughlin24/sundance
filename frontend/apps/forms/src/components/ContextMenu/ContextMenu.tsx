import { createPortal } from "react-dom";
import { ContextMenuButton } from "./ContextMenuButton";
import { useContextMenuDispatch, useContextMenuTarget } from "./useContextMenu";
import { AnimatePresence, motion } from "motion/react";
import { contextMenuStyles } from "./ContextMenu.styles";
import Box from "@mui/material/Box";
import { useEffect, useRef } from "react";

export interface ContextMenuProps {
  container: Element | DocumentFragment;
  width?: string | number;
  children: (data: unknown) => React.ReactNode | undefined;
}

interface ContextMenuComponent extends React.FC<ContextMenuProps> {
  Button: typeof ContextMenuButton;
}

const MENU_ITEM_SELECTOR = '[role="menuitem"]:not([disabled])';

const ContextMenu: ContextMenuComponent = function ({
  container,
  children,
  width = "244px",
}) {
  const { close } = useContextMenuDispatch();
  const target = useContextMenuTarget();
  const menuRef = useRef<HTMLDivElement | null>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);

  const getMenuItems = (): HTMLElement[] => {
    const menu = menuRef.current;

    if (!menu) {
      return [];
    }

    return Array.from(menu.querySelectorAll<HTMLElement>(MENU_ITEM_SELECTOR));
  };

  useEffect(() => {
    if (target) {
      previousFocusRef.current = document.activeElement as HTMLElement | null;

      const frame = requestAnimationFrame(() => {
        getMenuItems()[0]?.focus();
      });

      return () => cancelAnimationFrame(frame);
    }

    if (previousFocusRef.current) {
      previousFocusRef.current.focus();
      previousFocusRef.current = null;
    }
  }, [target]);

  const handleBackdropClk = (event: React.MouseEvent) => {
    event.stopPropagation();
    close();
  };

  const handleMenuKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
    const items = getMenuItems();

    if (event.key === "Escape") {
      event.preventDefault();
      close();
      return;
    }

    if (items.length === 0) {
      return;
    }

    const currentIndex = items.indexOf(document.activeElement as HTMLElement);

    switch (event.key) {
      case "ArrowDown":
        event.preventDefault();
        items[(currentIndex + 1) % items.length].focus();
        break;
      case "ArrowUp":
        event.preventDefault();
        items[(currentIndex - 1 + items.length) % items.length].focus();
        break;
      case "Home":
        event.preventDefault();
        items[0].focus();
        break;
      case "End":
        event.preventDefault();
        items[items.length - 1].focus();
        break;
      case "Tab":
        event.preventDefault();
        const direction = event.shiftKey ? -1 : 1;
        items[(currentIndex + direction + items.length) % items.length].focus();
        break;
    }
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
            ref={menuRef}
            role="menu"
            aria-orientation="vertical"
            aria-label="Context menu"
            onKeyDown={handleMenuKeyDown}
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
