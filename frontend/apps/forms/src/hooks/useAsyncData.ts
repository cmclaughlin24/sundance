import { useEffect, useState, type DependencyList } from "react";

/**
 * `resolveHttpService` returns an instance of the specified HTTP service class. It checks the `instanceCache`
 * to see if a valid instance already exists for the given service and `baseURL`. If a cached instance is found
 * with a matching `baseURL`, it is returned directly. Otherwise, a new instance is created, stored in the cache,
 * and returned.
 * @param Ctor The HTTP service class to resolve.
 * @param baseURL The base URL for the HTTP service. If the cached instance has a different base URL, it will be replaced.
 * @returns An instance of the specified HTTP service.
 */
export function useAsyncData<T>(
  operation: () => Promise<T>,
  deps?: DependencyList,
) {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<unknown>(null);

  useEffect(() => {
    let isCancelled = false;

    const run = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const result = await operation();

        if (isCancelled) {
          return;
        }

        setData(result);
      } catch (error) {
        if (isCancelled) {
          return;
        }
        setError(error);
      } finally {
        if (isCancelled) {
          return;
        }
        setIsLoading(false);
      }
    };

    run();

    return () => {
      isCancelled = true;
    };
  }, deps);

  return { data, isLoading, error };
}
