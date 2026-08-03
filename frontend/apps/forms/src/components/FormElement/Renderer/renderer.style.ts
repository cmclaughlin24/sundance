import type { SxProps, Theme } from "@mui/material/styles";

const baseStyles: SxProps<Theme> = {
  display: "flex",
  flexDirection: "column",
  alignItems: "start",
  gap: 5,
};

export const rendererStyles: Record<string, SxProps<Theme>> = {
  form: {
    ...baseStyles,
    maxWidth: "901px",
    position: "relative",
    flex: 1,
    padding: 6.5,
    paddingBottom: 0,
    minHeight: 0,
  },
  page: { ...baseStyles },
  section: { ...baseStyles },
};
