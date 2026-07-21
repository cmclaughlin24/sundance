import { BaseHttpService } from "./baseHttpService";

export class DataSourcesService extends BaseHttpService {
  static readonly serviceKey = "DataSourcesService";

  constructor(baseURL: string) {
    super(baseURL);
  }
}
