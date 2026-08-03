import type { SxProps } from "@mui/material/styles";
import type { Theme } from "node_modules/@emotion/react/dist/declarations/src";

export const designerStyles: Record<string, SxProps<Theme>> = {
  page: {
    maxWidth: "1440px",
    display: "flex",
    flexDirection: "column",
    gap: 5,
  },
  toolbar: {
    marginBottom: 3,
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
};
