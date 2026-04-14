import { ThemeProvider } from "@emotion/react";
import { createTheme } from "@mui/material";
import { PropsWithChildren } from "react";

// Gold: #FFC425

export const SundanceTheme = createTheme({
  palette: {
    primary: { main: "#0C1F40" },
    secondary: { main: "#E51837" },
  },
});

type SundanceThemeType = typeof SundanceTheme;

declare module "@emotion/react" {
  export interface Theme extends SundanceThemeType {}
}

const SundanceThemeProvider: React.FC<PropsWithChildren> = function ({
  children,
}) {
  return <ThemeProvider theme={SundanceTheme}>{children}</ThemeProvider>;
};

export default SundanceThemeProvider;
