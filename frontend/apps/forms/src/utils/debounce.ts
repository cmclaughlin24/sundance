export interface DebouncedFunction<A extends unknown[]> {
  (...args: A): void;
  cancel: () => void;
  flush: () => void;
}

/**
 * Creates a debounced version of a function that delays execution until after `delay` milliseconds
 * have elapsed since the last time it was invoked. Includes `cancel` and `flush` methods.
 *
 * @param fn The callback function to debounce.
 * @param delay The delay in milliseconds. Defaults to 300ms.
 * @returns A debounced function with `cancel()` and `flush()` controls.
 */
export function debounce<A extends unknown[]>(
  fn: (...args: A) => void | Promise<void>,
  delay = 300,
): DebouncedFunction<A> {
  let timer: ReturnType<typeof setTimeout> | undefined;
  let lastArgs: A | undefined;

  const debounced = (...args: A) => {
    lastArgs = args;
    timer && clearTimeout(timer);

    timer = setTimeout(() => {
      timer = undefined;
      const currentArgs = lastArgs;
      lastArgs = undefined;
      currentArgs && fn(...currentArgs);
    }, delay);
  };

  debounced.cancel = () => {
    if (timer) {
      clearTimeout(timer);
      timer = undefined;
    }
    lastArgs = undefined;
  };

  debounced.flush = () => {
    if (!timer || !lastArgs) {
      return;
    }

    clearTimeout(timer);
    timer = undefined;
    const currentArgs = lastArgs;
    lastArgs = undefined;
    fn(...currentArgs);
  };

  return debounced;
}
