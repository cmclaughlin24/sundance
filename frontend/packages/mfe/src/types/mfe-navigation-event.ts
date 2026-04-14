/**
 * Type of routing action that occurred. This can be either a `PUSH` action, which indicates
 * that a new entry was added to the history stack, or a `REPLACE` action, which indicates that
 * the current entry in the history stack was replaced with a new one.
 */
export type NavigationAction = "PUSH" | "REPLACE";

/**
 * Standardized event object for micro-frontend navigation events. This interface is used to synchronize
 * routing state between the host application and remote application(s).
 */
export interface MfeNavigationEvent {
  /**
   * The type of navigation action that occurred.
   */
  action: NavigationAction;

  /**
   * The new path that the application's router has navigated to.
   */
  pathname: string;
}
