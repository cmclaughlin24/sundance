import type { Theme } from "@emotion/react";
import type { SxProps } from "@mui/material/styles";

export const formCardStyles: Record<string, SxProps<Theme>> = {
  card: {
    padding: 2,
  },
  name: {
    fontSize: "1rem",
    fontWeight: 600,
    letterSpacing: 0,
    lineHeight: 1.5,
    marginBottom: 1,
  },
  description: {
    fontSize: "1rem",
    fontWeight: 300,
    letterSpacing: 0,
    lineHeight: 1.25,
  },
};
