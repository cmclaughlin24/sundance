import { DataSourceSectionLayout } from "../layout/DataSourceSectionLayout";
import type { DataSourceAttributeFormComponent } from "./DataSourceAttributeForm";

export const StaticAttributeForm: DataSourceAttributeFormComponent =
  function ({ dataSource }) {
    return (
      <DataSourceSectionLayout title="Add your list of options.">
        static form not implemented
      </DataSourceSectionLayout>
    );
  };
