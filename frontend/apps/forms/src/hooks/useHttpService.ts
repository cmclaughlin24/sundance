import type { BaseHttpService } from "@/services/baseHttpService";
import { DataSourcesService } from "@/services/dataSourcesService";
import { FormsService } from "@/services/formsService";
import { SubmissionsService } from "@/services/submissionService";
import type { CreateAxiosDefaults } from "axios";

/**
 */
interface HttpServiceClass<T extends BaseHttpService = BaseHttpService> {
  readonly serviceKey: string;
  new (baseURL: string, config?: CreateAxiosDefaults): T;
}

/**
 * `instanceCache` is a memozation cache that stores instances of HTTP service classes. It maps a service's unique `serviceKey`
 * to an instance of the service. This allows for reuse of service instances when the same service is requested.
 */
const instanceCache = new Map<string, BaseHttpService>();

/**
 * `resolveHttpService` provides an instance of the specified HTTP service class. It checks the `instanceCache`
 * to see if an instance already exists for the given `baseURL` and returns it if available. Otherwise, it creates a new instance, stores
 * it in the cache, and returns it.
 * @param Ctor The HTTP service class to instantiate.
 * @param baseURL The base URL for the HTTP service.
 * @returns An instance of the specified HTTP service.
 */
function resolveHttpService<T extends BaseHttpService>(
  Ctor: HttpServiceClass<T>,
  baseURL: string,
) {
  const cached = instanceCache.get(Ctor.serviceKey);

  if (cached && cached.isBaseURL(baseURL)) {
    return cached as T;
  }

  const service = new Ctor(baseURL);
  instanceCache.set(Ctor.serviceKey, service);
  return service;
}

/**
 * `useDataSourcesService` is a custom hook that provides an instance of the `DataSourcesService` class.
 * @returns An instance of the `DataSourcesService` class.
 */
export function useDataSourcesService() {
  return resolveHttpService(
    DataSourcesService,
    import.meta.env.VITE_TENANTS_API_URL,
  );
}

/**
 * `useFormsService` is a custom hook that provides an instance of the `FormsService` class.
 * @returns An instance of the `FormsService` class.
 */
export function useFormsService() {
  return resolveHttpService(FormsService, "/forms-api");
}

/**
 * `useSubmissionsService` is a custom hook that provides an instance of the `SubmissionsService` class.
 * @returns An instance of the `SubmissionsService` class.
 */
export function useSubmissionsService() {
  return resolveHttpService(SubmissionsService, "/tenants-api");
}
