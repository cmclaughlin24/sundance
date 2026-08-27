import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";

export type TagProps = React.PropsWithChildren<{ sx?: SxProps<Theme> }>;

const styles: SxProps<Theme> = (theme) => ({
  display: "inline-block",
  px: 1,
  fontSize: "0.75rem",
  background: theme.palette.secondary.main,
  borderRadius: "9999px",
});

export const Tag: React.FC<TagProps> = function ({ children, sx }) {
  return (
    <Box sx={mergeSx(styles, sx)}>
      {children}
    </Box>
  );
};
