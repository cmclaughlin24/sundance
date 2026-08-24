interface CacheEntry<T> {
  value: T;
  expires: number | null;
}

export interface CachedOptions<Args> {
  /**
   * The time-to-livve (TTL) for the cached value in milliseconds. After this duration, the cached value
   * will be considered expired and will be removed from the cache.
   */
  ttl?: number;

  /**
   * A function to generate a unqiue cache key based on teh function arguments.
   * @param args - The arguments passed to the cached function.
   * @returns A string representing the cache key.
   */
  key?: (args: Args) => string;
}

/**
 * A decorator to cache the result of an asynchronous function based on its arguments.
 * @param options - The caching options
 * @returns A decorator function that can be applied to class methods.
 */
export function Cached<This, Args extends unknown[], Return>(
  options: CachedOptions<Args> = { ttl: 60_000 },
) {
  const store = new Map<string, CacheEntry<Awaited<Return>>>();

  return function (
    target: (this: This, ...args: Args) => Return,
    context: ClassMethodDecoratorContext<
      This,
      (this: This, ...args: Args) => Return
    >,
  ) {
    return async function (
      this: This,
      ...args: Args
    ): Promise<Awaited<Return>> {
      const key = options.key
        ? options.key(args)
        : `${String(context.name)}:${JSON.stringify(args)}`;

      const hit = store.get(key);

      if (hasHit<Awaited<Return>>(hit)) {
        return hit!.value;
      }

      const value = (await target.apply(this, args)) as Awaited<Return>;
      store.set(key, {
        value,
        expires: options.ttl ? Date.now() + options.ttl : null,
      });
      return value;
    };
  };
}

/**
 * Legacy (experimentalDecorators) variant of {@link Cached}.
 * Uses the TypeScript stage-1 method decorator signature.
 * @param options - The caching options
 * @returns A method decorator that can be applied to class methods.
 */
export function LegacyCache<Args extends unknown[], Return>(
  options: CachedOptions<Args> = { ttl: 60_000 },
) {
  const store = new Map<string, CacheEntry<Awaited<Return>>>();

  return function (
    _target: object,
    propertyKey: string | symbol,
    descriptor: PropertyDescriptor,
  ): PropertyDescriptor {
    const original = descriptor.value as (
      this: unknown,
      ...args: Args
    ) => Return;

    descriptor.value = async function (
      this: unknown,
      ...args: Args
    ): Promise<Awaited<Return>> {
      const key = options.key
        ? options.key(args)
        : `${String(propertyKey)}:${JSON.stringify(args)}`;

      const hit = store.get(key);

      if (hasHit<Awaited<Return>>(hit)) {
        return hit!.value;
      }

      const value = (await original.apply(this, args)) as Awaited<Return>;
      store.set(key, {
        value,
        expires: options.ttl ? Date.now() + options.ttl : null,
      });
      return value;
    };

    return descriptor;
  };
}

function hasHit<T>(hit: CacheEntry<T> | undefined): boolean {
  if (!hit) {
    return false;
  }

  return hit.expires === null || hit.expires > Date.now();
}
