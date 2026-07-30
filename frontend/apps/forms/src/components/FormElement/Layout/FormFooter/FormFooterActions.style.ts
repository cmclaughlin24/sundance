import type { Theme } from "@emotion/react";
import type { SxProps } from "@mui/material/styles";

export const formFooterActionsStyles: Record<string, SxProps<Theme>> = {
  container: {
    "*:not(:last-child)": {
      marginRight: 2.5,
    },
  },
};
