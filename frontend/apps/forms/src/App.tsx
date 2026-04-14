import { useEffect } from "react";
import "./App.css";
import { RouterProvider } from "@tanstack/react-router";
import { SundanceThemeProvider } from "@sundance/common";
import type { MfeBootstrapOptions } from "@sundance/mfe";

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
      <RouterProvider router={router} />
    </SundanceThemeProvider>
  );
}

export default App;
