import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import {
  createMemoryHistory,
  createRouter,
  type RouterHistory,
} from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen.ts";
import "./index.css";
import App from "./App.tsx";
import {
  parsePath,
  type MfeBootstrapFn,
  type MfeBootstrapOptions,
} from "@sundance/mfe";

export function getRouter(history: RouterHistory, basePath?: string) {
  const router = createRouter({
    routeTree: routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
    history,
    basepath: basePath,
  });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}

export const bootstrap: MfeBootstrapFn = (
  rootElement: HTMLElement | null,
  options: MfeBootstrapOptions,
) => {
  if (!rootElement) {
    throw new Error("rootElement is required to bootstrap the application.");
  }

  const history =
    options.defaultHistory ||
    createMemoryHistory({
      initialEntries: [parsePath(options.initialPath!, options.basePath!)],
    });
  const router = getRouter(history, options.basePath);
  const root = ReactDOM.createRoot(rootElement);

  root.render(
    <StrictMode>
      <App router={router} options={options} />
    </StrictMode>,
  );

  return {
    onParentNavigate(arg) {
      if (history.location.pathname === arg.pathname) {
        return;
      }

      history.push(arg.pathname);
    },
  };
};
