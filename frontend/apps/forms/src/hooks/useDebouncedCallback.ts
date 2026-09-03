import { useMemo, useRef } from "react";

/**
 * A custom hook that returns a debounced version of the given callback function.
 * @param callback The callback function to debounce.
 * @param delay The debounce delay in milliseconds. Defaults to 500ms.
 * @returns An object containing the debounced function, a cancel function, and a flush function.
 */
export function useDebouncedCallback<A extends unknown[]>(
  callback: (...args: A) => void,
  delay = 500,
) {
  const timerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const callbackRef = useRef(callback);
  callbackRef.current = callback;

  const debounced = useMemo(
    () =>
      (...args: A) => {
        clearTimeout(timerRef.current);
        timerRef.current = setTimeout(
          () => callbackRef.current(...args),
          delay,
        );
      },
    [],
  );

  const cancel = () => clearTimeout(timerRef.current);

  const flush = (...args: A) => {
    if (!timerRef.current) {
      return;
    }

    cancel();
    callbackRef.current(...args);
  };

  return { debounced, cancel, flush };
}
