import type { ILookup } from "@/types/data";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";

export class DataSourcesService extends BaseHttpService {
  static readonly serviceKey = "DataSourcesService";

  constructor(baseURL: string) {
    super(baseURL);
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
      `/data-sources/${dataSourceId}/look-ups`,
      { headers, params },
    );

    return resp.data;
  }
}
