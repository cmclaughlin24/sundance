import type { SxProps, Theme } from "@mui/material/styles";

export const formProgressStyles: Record<string, SxProps<Theme>> = {
  container: {
    marginBottom: 2.5,
  },
  required: {
    fontSize: "0.75rem",
    fontWeight: 600,
  },
  text: {
    fontWeight: 600,
  },
  error: {
    color: "#B42D19",
    ":hover": {
      cursor: "pointer",
      textDecoration: "underline",
    },
  },
};
