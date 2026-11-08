import type { ILookup } from "./data";

export interface DataSourceHttpRequest {
  url: string;
  method: string;
  headers?: Record<string, string>;
  valueField: string;
  labelField: string;
}

export interface StaticDataSourceAttributes {
  data: ILookup[];
}

export interface ScheduledDataSourceAttributes extends DataSourceHttpRequest {
  data: ILookup[];
  intervalHours: number;
  expirationDate: Date;
  attempts?: number;
}

export interface WebhookDataSourceAttributes extends DataSourceHttpRequest {
  requiredKeys: string[];
}

export interface DataLakeDataSourceAttributes {
  query: string;
  requiredKeys: string[];
  optionalKeys?: string[];
  catalog: string;
  schema: string;
  valueField: string;
  labelField: string;
  timeoutMs: number;
}

export type DataSourceAttributes =
  | StaticDataSourceAttributes
  | ScheduledDataSourceAttributes
  | WebhookDataSourceAttributes
  | DataLakeDataSourceAttributes;
