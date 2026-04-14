import { MfeNavigationEvent } from "./mfe-navigation-event";

/**
 * Interface for callbacks that the micro-frontend (MFE) can use to communicate with
 * child applications from the host application.
 */
export interface MfeCallbacks {
  /**
   * A callback function that will be called when the micro-frontend application navigates to a new route. This
   * is required for the host application to be able to synchronize its own routing state with the micro-frontend application's
   * routing state.
   *
   * @param event - An object containing the action (e.g., "PUSH", "REPLACE") and the new path that the micro-frontend application has
   * navigated to.
   */
  onParentNavigate: (event: MfeNavigationEvent) => void;
}
