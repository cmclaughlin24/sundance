import ReactDOM from "react-dom/client";
import { createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import App from "./App";

function getRouter() {
  const router = createRouter({
    routeTree: routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
  });

  return router;
}

// NOTE: To allow the host micro-frontend to navigate to remote routes, we need to disable the type safety for the router.
// This allows the host to call `router.navigate("/some-route")` without TypeScript errors, even if "/some-route" is
// not defined in the host's route tree.
// declare module "@tanstack/react-router" {
//   interface Register {
//     router: ReturnType<typeof getRouter>;
//   }
// }

export function bootstrap(rootElement: HTMLElement | null) {
  if (!rootElement || rootElement.innerHTML) {
    return;
  }

  const router = getRouter();
  const root = ReactDOM.createRoot(rootElement);

  root.render(<App router={router} />);
}
