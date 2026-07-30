import type { SxProps, Theme } from "@mui/material/styles";

export const fieldElementContainerStyles: Record<string, SxProps<Theme>> = {
  container: {
    display: "flex",
    gap: 1,
    paddingBottom: 1.25,
  },
  label: {
    flex: 1,
    marginRight: 4,
  },
  input: {
    flex: 1,
    width: "100%",
    "> *": {
      width: "100%",
    },
  },
};
