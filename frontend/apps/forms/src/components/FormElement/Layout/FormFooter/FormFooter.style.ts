import type { SxProps, Theme } from "@mui/material/styles";

export const formFooterStyles: Record<string, SxProps<Theme>> = {
  footer: {
    marginTop: 2.25,
    position: "sticky",
    bottom: 0,
  },
  content: {
    paddingX: 8,
    paddingY: 2.5,
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    borderTop: "1px solid #2b2b2b",
  },
  name: {
    fontWeight: 600,
    letterSpacing: 0,
  },
};
