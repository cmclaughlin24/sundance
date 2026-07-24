import { useEffect } from "react";
import "./App.css";
import { RouterProvider } from "@tanstack/react-router";
import { SundanceThemeProvider } from "@sundance/common";
import type { MfeBootstrapOptions } from "@sundance/mfe";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";

export interface AppProps {
  router: any;
  options: MfeBootstrapOptions;
}

function App({ router, options }: AppProps) {
  useEffect(() => {
    let unsubscribe: () => void;

    if (options.onNavigate) {
      unsubscribe = router.history.subscribe((arg: any) =>
        options.onNavigate!({
          action: arg.action.type,
          pathname: arg.location.pathname,
        }),
      );
    }

    return () => unsubscribe && unsubscribe();
  }, [router, options]);

  return (
    <SundanceThemeProvider>
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <RouterProvider router={router} />
      </LocalizationProvider>
    </SundanceThemeProvider>
  );
}

export default App;
