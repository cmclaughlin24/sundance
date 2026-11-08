import Box from "@mui/material/Box";
import { dataSourcesCatalogPanelStyles as styles } from "./DataSourcesCatalogPanel.styles";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useDataSourcesService } from "@/hooks/useHttpService";
import { useMemo, type ChangeEvent } from "react";
import TextField from "@mui/material/TextField";
import { useDebounce } from "@/hooks/useDebounce";
import * as ArrayUtils from "@/utils/array";
import { TENANT_ID } from "@/constants/tenant";
import { DataSourceCatalogItem } from "./DataSourceCatalogItem";
import type { DataSource } from "@/types/dataSource";
import { Panel } from "@/components/layout/Panel";

export const DataSourcesCatalogPanel: React.FC<{
  dataSource?: DataSource;
  onChange: (dataSource: DataSource) => void;
}> = function ({ dataSource: _, onChange }) {
  const {
    value: searchTerm,
    debounceValue: debounceSearchTerm,
    setValue: setSearchTerm,
  } = useDebounce("", []);

  const dataSourcesService = useDataSourcesService();

  const { data, isLoading, error } = useAsyncData(
    async (token: string) => {
      return await dataSourcesService.getDataSources({
        token,
        tenantId: TENANT_ID,
      });
    },
    [dataSourcesService],
  );

  const dataSources = useMemo(() => {
    if (!ArrayUtils.hasLengthGreaterThan(data, 0)) {
      return [];
    }

    return data!.filter(
      (ds) =>
        ds.name.includes(debounceSearchTerm) ||
        ds.description.includes(debounceSearchTerm),
    );
  }, [data, debounceSearchTerm]);

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    setSearchTerm(event.target.value);
  };

  if (isLoading) {
    return <>Loading Data Source Catalog.</>;
  }

  if (error) {
    return <>Something went wrong listing the data sources.</>;
  }

  return (
    <Panel sx={styles.panel}>
      <Panel.Title title="Data Sources Catalog" />
      <TextField
        placeholder="Search"
        value={searchTerm}
        onChange={handleChange}
      />
      <Box component="ul" sx={styles.catalog}>
        {ArrayUtils.hasExactLength(dataSources, 0) && (
          <>No Data Sources Configured</>
        )}
        {ArrayUtils.hasLengthGreaterThan(dataSources, 0) &&
          dataSources.map((ds) => (
            <Box component="li">
              <DataSourceCatalogItem dataSource={ds} onClick={onChange} />
            </Box>
          ))}
      </Box>
    </Panel>
  );
};
