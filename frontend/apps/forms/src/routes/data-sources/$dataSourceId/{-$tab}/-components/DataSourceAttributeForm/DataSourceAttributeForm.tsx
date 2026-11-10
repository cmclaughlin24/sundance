import type { DataSourceType, DataSource } from "@/types/dataSource";
import { DataLakeAttributeForm } from "./DataLakeAttributeForm";
import { ScheduledAttributeForm } from "./ScheduledAttributeForm";
import { StaticAttributeForm } from "./StaticAttributeForm";
import { WebhookAttributeForm } from "./WebhookAttributeForm";

export type DataSourceAttributeFormComponent = React.FC<{
  dataSource: DataSource;
}>;

const registry = new Map<DataSourceType, DataSourceAttributeFormComponent>([
  ["data-lake", DataLakeAttributeForm],
  ["scheduled", ScheduledAttributeForm],
  ["static", StaticAttributeForm],
  ["webhook", WebhookAttributeForm],
]);

export const DataSourceAttributeForm: React.FC<{ dataSource: DataSource }> =
  function ({ dataSource }) {
    const FormComponent = registry.get(dataSource.type);

    if (!FormComponent) {
      return <>data source attribute form not defined!</>;
    }

    return <FormComponent dataSource={dataSource} />;
  };
