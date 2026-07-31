import type { SxProps, Theme } from "@mui/material/styles";

export const fieldElementLabelStyles: Record<string, SxProps<Theme>> = {
  label: {
    fontSize: "2.125rem",
    fontWeight: 300,
    letterSpacing: 0,
    lineHeight: 1.18,
    display: "inline-block",
    marginBottom: 1,
    "::after": {
      content: '"*"',
      color: "#971E28",
      marginLeft: 1,
    },
  },
  description: {
    fontSize: "1rem",
    fontWeight: 400,
    letterSpacing: 0,
    lineHeight: 1.25,
  },
};
