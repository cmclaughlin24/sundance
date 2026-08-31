import { useEffect, type DependencyList } from "react";
import { useState } from "react";

/**
 * Custom hook that debounces a value.
 * @param initial The initial value.
 * @param delay The debounce value in milliseconds.
 * @returns An object containing the current value, the debounce value, and a setter for the value.
 */
export function useDebounce<T>(
  initial: T,
  deps: DependencyList,
  delay: number = 500,
) {
  const [value, setValue] = useState(initial);
  const [debounceValue, setDebounceValue] = useState(initial);

  useEffect(() => {
    const timerRef = setTimeout(() => setDebounceValue(value), delay);
    return () => clearTimeout(timerRef);
  }, [value]);

  useEffect(() => {
    setValue(initial);
    setDebounceValue(initial);
  }, deps);

  return { value, debounceValue, setValue };
}
