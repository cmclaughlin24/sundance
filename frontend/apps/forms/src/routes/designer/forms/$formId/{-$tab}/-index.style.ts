import type { SxProps, Theme } from "@mui/material/styles";

export const formDesignerPageStyles: Readonly<Record<string, SxProps<Theme>>> =
  {
    page: {
      marginTop: 2.5,
      padding: 6.5,
      width: "100%",
      maxWidth: 1440,
      display: "flex",
      flexDirection: "column",
      alignItems: "start",
      gap: 2.5,
    },
  };
