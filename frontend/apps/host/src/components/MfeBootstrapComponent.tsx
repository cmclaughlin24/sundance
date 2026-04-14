import type { MfeBootstrapFn } from "@sundance/mfe";
import { useRouter } from "@tanstack/react-router";
import { useCallback, type FC } from "react";

export function MfeBootstrapComponent(
  basePath: string,
  bootstrap: MfeBootstrapFn,
): FC {
  return function () {
    const router = useRouter();

    const appRef = useCallback((node: HTMLDivElement) => {
      if (!node) {
        return;
      }

      const { onParentNavigate } = bootstrap(node, {
        basePath: basePath,
        initialPath: router.history.location.pathname,
        onNavigate: (event) => {
          if (router.history.location.pathname === event.pathname) {
            return;
          }

          router.history.push(event.pathname);
        },
      });

      const unsubscribe = router.history.subscribe((arg: any) =>
        onParentNavigate({
          action: arg.action.type,
          pathname: arg.location.pathname,
        }),
      );

      return () => unsubscribe();
    }, []);

    return <div ref={appRef}></div>;
  };
}
