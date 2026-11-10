import type { ILookup } from "@/types/data";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";
import { CONFIG } from "@/constants/config";
import type { DataSource } from "@/types/dataSource";
import type { DataSourceRequest } from "./dataSourcesService.type";
import type { DataSourceAttributes } from "@/types/dataSourceAttributes";

export class DataSourcesService extends BaseHttpService {
  static readonly serviceKey = "DataSourcesService";

  constructor() {
    super(CONFIG.tenantsUrl);
  }

  async getDataSources(options: DefaultRequestOptions): Promise<DataSource[]> {
    return this._get<DataSource[]>("/api/v1/data-sources", options);
  }

  async getDataSource(
    id: string,
    options: DefaultRequestOptions,
  ): Promise<DataSource> {
    return this._get<DataSource>(`/api/v1/data-sources/${id}`, options);
  }

  async getLookups(
    dataSourceId: string,
    filters: Record<string, any> | null,
    options: DefaultRequestOptions,
  ): Promise<ILookup[]> {
    let params: { params?: string } | undefined = undefined;
    if (filters) {
      params = { params: JSON.stringify(filters) };
    }

    const headers = this._defaultRequestHeaders(options);
    const resp = await this._client.get<ILookup[]>(
      `/api/v1/data-sources/${dataSourceId}/look-ups`,
      { headers, params },
    );

    return resp.data;
  }

  /**
   * Creates a new data source.
   * @param dataSource The data source to create.
   * @param options The default request options.
   * @returns A promise that resolves to the created data source.
   */
  async create(
    dataSource: DataSourceRequest,
    options: DefaultRequestOptions,
  ): Promise<DataSource> {
    const resp = await this._post<DataSourceRequest, DataSource>(
      `/api/v1/data-sources`,
      { attributes: {} as DataSourceAttributes, ...dataSource },
      options,
    );

    return resp.data;
  }

  async update(
    id: string,
    dataSource: DataSourceRequest,
    options: DefaultRequestOptions,
  ): Promise<DataSource> {
    const resp = await this._put<DataSourceRequest, DataSource>(
      `/api/v1/data-sources/${id}`,
      dataSource,
      options,
    );

    return resp.data;
  }

  async delete(id: string, options: DefaultRequestOptions): Promise<void> {
    return this._delete(`/api/v1/data-sources/${id}`, options);
  }
}
