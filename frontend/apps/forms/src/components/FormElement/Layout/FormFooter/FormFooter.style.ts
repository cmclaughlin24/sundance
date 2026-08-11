import type { SxProps, Theme } from "@mui/material/styles";

const baseContentStyles: SxProps<Theme> = {
  paddingX: 8,
  paddingY: 2.5,
  display: "flex",
  alignItems: "center",
  borderTop: "1px solid #2b2b2b",
};

export const formFooterStyles: Record<string, SxProps<Theme>> = {
  footer: {
    marginTop: 2.25,
    position: "sticky",
    bottom: 0,
  },
  fullContent: {
    ...baseContentStyles,
    justifyContent: "space-between",
  },
  compactContent: {
    ...baseContentStyles,
    justifyContent: "flex-end",
    gap: 2.5
  },
  name: {
    fontWeight: 600,
    letterSpacing: 0,
  },
};
