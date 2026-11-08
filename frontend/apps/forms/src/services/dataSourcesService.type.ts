import type { DataSourceType } from "@/types/dataSource";
import type { DataSourceAttributes } from "@/types/dataSourceAttributes";

export interface DataSourceRequest {
  name: string;
  description: string;
  type: DataSourceType;
  attributes?: DataSourceAttributes;
}
