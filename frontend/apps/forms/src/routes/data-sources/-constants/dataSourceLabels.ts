import type { DataSourceType } from "@/types/dataSource";

export const dataSourceTypeLabels: Readonly<Record<DataSourceType, string>> = {
  "data-lake": "Data Lake",
  scheduled: "Scheduled",
  webhook: "Webhook",
  static: "Static",
};
