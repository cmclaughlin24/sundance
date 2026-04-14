import type { RouterHistory } from "@tanstack/react-router";
import { MfeNavigationEvent } from "./mfe-navigation-event";

/**
 * Options for bootstrapping a micro-frontend application.
 */
export interface MfeBootstrapOptions {
  /**
   * An optional `RouterHistory` instance to be used by the micro-frontend application. The remote
   * application will create its own history instance if this is not provide.
   *
   * @note This commonly provide by th remote application when it's bootstrapped as a standalone application.
   */
  defaultHistory?: RouterHistory;

  /**
   * An optional base path for the micro-frontend application. This is used to ensure that `@tanstack/react-router` can
   * properly handle routing with the micro-frontend application.
   */
  basePath?: string;

  /**
   * An optional initial path to navigate to when the micro-frontend application is bootstrapped.
   */
  initialPath?: string;

  /**
   * An optional callback function that will be called when the micro-frontend application navigates to a new route. This
   * is required for the host application to be able to synchronize its own routing state with the micro-frontend application's
   * routing state.
   */
  onNavigate?: (event: MfeNavigationEvent) => void;
}
