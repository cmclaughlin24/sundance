import { Page } from "@/components/layout/Page/Page";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import Box from "@mui/material/Box";
import { createFileRoute } from "@tanstack/react-router";
import { dataSourcesStyles as styles } from "./-index.styles";
import { useDataSourcesService } from "@/hooks/useHttpService";
import { useAsyncData } from "@/hooks/useAsyncData";
import { TENANT_ID } from "@/constants/tenant";
import { CardList } from "@/components/CardList";
import { DataSourceCard } from "./-components/DataSourceCard";
import type { DataSource } from "@/types/dataSource";
import {
  CreateDataSourceDrawer,
  type CreateDataSourceEvent,
} from "./-components/CreateDataSourceDrawer";
import { useRef, useState } from "react";
import type { DrawerHandle } from "@/components/Drawer";
import {
  DataSourcesToolbar,
  type DataSourcesToolbarFilters,
} from "./-components/DataSourcesToolbar";
import type { DefaultRequestOptions } from "@/services/baseHttpService";

export const Route = createFileRoute("/data-sources/")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = Route.useNavigate();
  const dataSourcesService = useDataSourcesService();
  const drawerRef = useRef<DrawerHandle>(null);
  const [filters, setFilters] = useState<DataSourcesToolbarFilters>({
    search: "",
  });

  const { data, isLoading, error } = useAsyncData(
    async (token: string) => {
      return await dataSourcesService.getDataSources({
        token,
        tenantId: TENANT_ID,
      });
    },
    [dataSourcesService],
  );

  const handleCardClk = (dataSource: DataSource) => {
    navigate({
      to: "/data-sources/$dataSourceId",
      params: { dataSourceId: dataSource.id },
    });
  };

  const handleCreate = async ({
    name,
    description,
    type,
  }: CreateDataSourceEvent) => {
    try {
      const opts: DefaultRequestOptions = {
        token: "placeholder",
        tenantId: TENANT_ID,
      };

      const dataSource = await dataSourcesService.create(
        { name, description, type },
        opts,
      );

      navigate({
        to: "/data-sources/$dataSourceId",
        params: { dataSourceId: dataSource.id },
      });
    } catch (error) {
      // FIXME: Correctly handle the error when something goes wrong.
    }
  };

  const handleNewClk = () => drawerRef.current?.open();

  if (isLoading) {
    return <>Loading Data Source Catalog.</>;
  }

  if (error) {
    return <>Something went wrong listing the data sources.</>;
  }

  return (
    <>
      <Page sx={styles.page}>
        <PageTitle
          name="Data Sources"
          description="Manage and monitor global lookup endpoints, background sync workers, and SQL connectors."
        />
        <Box>
          <DataSourcesToolbar
            filters={filters}
            onFilter={(event) => setFilters(event)}
            onNew={handleNewClk}
          />
          <CardList data={data} onClick={handleCardClk} keyFn={(ds) => ds.id}>
            {(dataSource) => <DataSourceCard dataSource={dataSource} />}
          </CardList>
        </Box>
      </Page>
      <CreateDataSourceDrawer ref={drawerRef} onCreate={handleCreate} />
    </>
  );
}
