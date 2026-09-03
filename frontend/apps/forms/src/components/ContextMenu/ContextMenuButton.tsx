import { Background } from "@/constants/colors";
import type { ButtonProps } from "@mui/material/Button";
import Button from "@mui/material/Button";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";

const styles: SxProps<Theme> = (theme) => ({
  boxShadow: "none",
  border: "none",
  borderRadius: 1.5,
  textTransform: "none",
  justifyContent: "flex-start",
  width: "100%",
  px: 1.5,
  py: 1,
  minWidth: 0,
  fontSize: "0.875rem",
  fontWeight: 400,
  lineHeight: 1.4,
  color: "font.primary",
  background: Background.Primary,
  ":hover": {
    boxShadow: "none",
    background: `${theme.palette.primary.main}70`,
  },
});

export const ContextMenuButton: React.FC<ButtonProps> = function ({
  children,
  sx,
  ...props
}) {
  return (
    <Button
      role="menuitem"
      tabIndex={-1}
      disableElevation
      disableRipple
      variant="text"
      sx={mergeSx(styles, sx)}
      {...props}
    >
      {children}
    </Button>
  );
};
