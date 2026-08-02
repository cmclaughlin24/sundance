import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";

export type PageProps = React.PropsWithChildren<{
  sx?: SxProps<Theme>;
}>;

const styles: SxProps<Theme> = {
  marginX: "auto",
  padding: 6.5,
  boxShadow: "0 1px 3px 0 rgba(0, 0, 0, .1)",
  border: "1px solid #2b2b2b",
  position: "relative",
};

export const Page: React.FC<PageProps> = function ({ sx, children }) {
  return <Box sx={sx ? { ...styles, ...sx } : styles}>{children}</Box>;
};
