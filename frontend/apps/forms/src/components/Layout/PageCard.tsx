import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";

export type PageCardProps = React.PropsWithChildren<{
  sx?: SxProps<Theme>;
}>;

const styles: SxProps<Theme> = {
  padding: 6.5,
  paddingBottom: 15,
  boxShadow: "0 1px 3px 0 rgba(0, 0, 0, .1)",
  border: "1px solid #2b2b2b",
  position: "relative",
};

export const PageCard: React.FC<PageCardProps> = function ({ sx, children }) {
  return <Box sx={sx ? { ...styles, ...sx } : styles}>{children}</Box>;
};
