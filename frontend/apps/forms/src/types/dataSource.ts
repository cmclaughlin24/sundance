import type {
  DataLakeDataSourceAttributes,
  DataSourceAttributes,
  ScheduledDataSourceAttributes,
  StaticDataSourceAttributes,
  WebhookDataSourceAttributes,
} from "./dataSourceAttributes";
import type { HasName } from "./formObject";

export type DataSourceType = "static" | "scheduled" | "webhook" | "data-lake";

export interface IDataSource<
  T extends DataSourceAttributes = DataSourceAttributes,
> extends HasName {
  id: string;
  tenantId: string;
  description: string;
  type: DataSourceType;
  attributes: T;
  createdAt: Date;
  updatedAt: Date;
}

export interface StaticDataSource extends IDataSource<StaticDataSourceAttributes> {
  type: "static";
}

export interface ScheduledDataSource extends IDataSource<ScheduledDataSourceAttributes> {
  type: "scheduled";
}

export interface WebhookDataSource extends IDataSource<WebhookDataSourceAttributes> {
  type: "webhook";
}

export interface DataLakeDataSource extends IDataSource<DataLakeDataSourceAttributes> {
  type: "data-lake";
}

export type DataSource =
  | StaticDataSource
  | ScheduledDataSource
  | WebhookDataSource
  | DataLakeDataSource;
