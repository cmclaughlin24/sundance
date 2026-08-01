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
  },
  page: { ...baseStyles },
  section: { ...baseStyles },
};
