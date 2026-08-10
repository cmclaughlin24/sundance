import type { Theme } from "@emotion/react";
import type { SxProps } from "@mui/material/styles";

export const formFooterActionsStyles: Record<string, SxProps<Theme>> = {
  container: {
    display: "flex",
    gap: 2.5,
    alignSelf: "flex-end",
  },
};
