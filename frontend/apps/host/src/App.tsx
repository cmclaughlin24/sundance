import { RouterProvider } from "@tanstack/react-router";
import "./App.css";
import { SundanceThemeProvider } from "@sundance/common";

export interface AppProps {
  router: any;
}

function App({ router }: AppProps) {
  return (
    <SundanceThemeProvider>
      <RouterProvider router={router} />
    </SundanceThemeProvider>
  );
}

export default App;
