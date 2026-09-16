import { useImperativeHandle, useState } from "react";
import CloseIcon from "@mui/icons-material/Close";
import MuiDrawer from "@mui/material/Drawer";
import { drawerStyles } from "./Drawer.style";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";

export interface DrawerHandle {
  open: () => void;
  close: () => void;
}

export interface DrawerProps extends React.PropsWithChildren {
  ref: React.Ref<DrawerHandle>;
  title: string;
}

export const Drawer: React.FC<DrawerProps> = function ({
  children,
  ref,
  title,
}) {
  const [isOpen, setIsOpen] = useState(false);

  useImperativeHandle(ref, () => ({
    open: () => setIsOpen(true),
    close: () => setIsOpen(false),
  }));

  return (
    <MuiDrawer
      anchor="right"
      open={isOpen}
      sx={{ zIndex: 1202 }}
      slotProps={{ paper: { sx: { width: 727, maxWidth: "100%" } } }}
      onClose={() => setIsOpen(false)}
      data-testid="drawer"
    >
      <Box sx={drawerStyles.header}>
        <IconButton
          size="small"
          sx={drawerStyles.closeIcon}
          onClick={() => setIsOpen(false)}
          data-testid="close-drawer-btn"
        >
          <CloseIcon fontSize="inherit" />
        </IconButton>
        <Typography component="h3" sx={drawerStyles.title}>
          {title}
        </Typography>
        <Box sx={drawerStyles.closeIcon} />
      </Box>
      {children}
    </MuiDrawer>
  );
};
