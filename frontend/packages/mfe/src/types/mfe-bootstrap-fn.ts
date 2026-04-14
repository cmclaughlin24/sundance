import { MfeBootstrapOptions } from "./mfe-bootstrap-options";
import { MfeCallbacks } from "./mfe-callbacks";

/**
 * A function type that defines the signature for bootstrapping a micro-frontend application.
 *
 * @param rootElement - The HTML element where the micro-frontend will be mounted.
 * @param options - Configuration options for bootstrapping the micro-frontend.
 * @returns An object containing callbacks for lifecycle events of the micro-frontend.
 */
export type MfeBootstrapFn = (
  rootElement: HTMLElement | null,
  options: MfeBootstrapOptions,
) => MfeCallbacks;
