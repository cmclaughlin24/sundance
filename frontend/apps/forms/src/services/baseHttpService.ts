import axios, {
  AxiosHeaders,
  type AxiosInstance,
  type CreateAxiosDefaults,
} from "axios";

export interface DefaultRequestOptions {
  /**
   * The bearer token used for authorization in the request headers.
   */
  token: string;

  /**
   * The tenant ID used for multi-tenant applications, included in the request headers.
   */
  tenantId: string;
}

export interface ApiCUDResponse<T> {
  /**
   * The message of the response, providing additional information about the result of the request.
   */
  message: string;

  /**
   * The data returned from the API, which can be of any type `T`.
   */
  data: T;
}

/**
 * `BaseHttpService` is an abstract class that provides a base implementation for making HTTP requests. It
 * encapsulates an Axios instance and provides methods for sending `GET`, `POST`, `PUT`, and `DELETE` requests.
 */
export abstract class BaseHttpService {
  protected _client: AxiosInstance;

  constructor(baseURL: string, config?: CreateAxiosDefaults) {
    this._client = axios.create({ baseURL, ...config });
  }

  isBaseURL(url: string): boolean {
    return this._client.defaults.baseURL === url;
  }

  /**
   * Sends a `GET` request to the specified URL.
   * @param url The URL to send the `GET` request to.
   * @param options The default request options.
   * @returns The response ddata of type `R`.
   */
  protected async _get<R>(
    url: string,
    options: DefaultRequestOptions,
  ): Promise<R> {
    const headers = this._defaultRequestHeaders(options);
    const response = await this._client.get<R>(url, { headers });
    return response.data;
  }

  /**
   * Sends a `POST` request to the specified URL with the given payload.
   * @param url The URL to send the `POST` request to.
   * @param payload The paylaod to include in the `POST` request.
   * @param options The default request options.
   * @returns The response data of type `ApiCUDReponse<R>`.
   */
  protected async _post<P, R>(
    url: string,
    payload: P,
    options: DefaultRequestOptions,
  ): Promise<ApiCUDResponse<R>> {
    const headers = this._defaultRequestHeaders(options);
    const response = await this._client.post<ApiCUDResponse<R>>(url, payload, {
      headers,
    });
    return response.data;
  }

  /**
   * Sends a `PUT` request to the specified URL with the given payload.
   * @param url The URL to send the `PUT` request to.
   * @param payload The paylaod to include in the `PUT` request.
   * @param options The default request options.
   * @returns The response data of type `ApiCUDReponse<R>`.
   */
  protected async _put<P, R>(
    url: string,
    payload: P,
    options: DefaultRequestOptions,
  ): Promise<ApiCUDResponse<R>> {
    const headers = this._defaultRequestHeaders(options);
    const response = await this._client.put<ApiCUDResponse<R>>(url, payload, {
      headers,
    });
    return response.data;
  }

  /**
   * Sends a `DELETE` request to the specified URL with the given payload.
   * @param url The URL to send the `DELETE` request to.
   * @param options The default request options.
   * @returns
   */
  protected async _delete(
    url: string,
    options: DefaultRequestOptions,
  ): Promise<void> {
    const headers = this._defaultRequestHeaders(options);
    await this._client.delete<ApiCUDResponse<void>>(url, { headers });
  }

  /**
   * Generates the default request headers based on the provided options.
   * @param options The default request options.
   * @returns The Axios headers object containing the default headers.
   */
  protected _defaultRequestHeaders(
    options: DefaultRequestOptions,
  ): AxiosHeaders {
    let headers = new AxiosHeaders();
    headers = this._setTenantHeader(headers, options.tenantId);
    headers = this._setBearerHeader(headers, options.token);
    return headers;
  }

  /**
   * Set the tenant ID header
   * @param headers The Axios headers object.
   * @param tenantId The ID of the tenant.
   * @returns The updated Axios headers object.
   */
  protected _setTenantHeader(
    headers: AxiosHeaders,
    tenantId: string,
  ): AxiosHeaders {
    return headers.set("X-Tenant-ID", tenantId);
  }

  /**
   * Set the bearer token authorization header
   * @param headers The Axios headers object.
   * @param token The bearer token.
   * @returns The updated Axios headers object.
   */
  protected _setBearerHeader(
    headers: AxiosHeaders,
    token: string,
  ): AxiosHeaders {
    return headers.set("Authorization", `Bearer ${token}`);
  }
}
