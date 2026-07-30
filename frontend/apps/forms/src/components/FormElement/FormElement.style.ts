import type { SxProps, Theme } from "@mui/material/styles";

export const formElementStyles: Record<string, SxProps<Theme>> = {
  container: {
    marginX: "auto",
    padding: 6.5,
    paddingBottom: 15,
    maxWidth: 1440,
    boxShadow: "0 1px 3px 0 rgba(0, 0, 0, .1)",
    border: "1px solid #2b2b2b",
    position: "relative",
  },
};
