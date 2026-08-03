import type { SxProps, Theme } from "@mui/material/styles";

export const formElementStyles: Record<string, SxProps<Theme>> = {
  page: {
    marginTop: 2.5,
    padding: 0,
    width: "100%",
    maxWidth: 1440,
    display: "flex",
    flexDirection: "column",
  },
};
