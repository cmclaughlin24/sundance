import type { SxProps, Theme } from "@mui/material/styles";

export const formRendererStyles: Record<string, SxProps<Theme>> = {
  form: {
    display: "flex",
    flexDirection: "column",
    alignItems: "start",
    gap: 5,
    maxWidth: "901px",
    position: "relative",
  },
};
