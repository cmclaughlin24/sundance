import { Page } from "@/components/layout/Page/Page";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import { TENANT_ID } from "@/constants/tenant";
import {
  resolveHttpService,
  useDataSourcesService,
} from "@/hooks/useHttpService";
import type { DefaultRequestOptions } from "@/services/baseHttpService";
import { DataSourcesService } from "@/services/dataSourcesService";
import { createFileRoute } from "@tanstack/react-router";
import { dataSourcePageStyles as styles } from "./-index.style";
import { DataSourceTypeTag } from "../-components/DataSourceTypeTag";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { DataSourceAttributeForm } from "./-components/DataSourceAttributeForm/DataSourceAttributeForm";
import { TestBench } from "./-components/TestBench";
import { CardList } from "@/components/CardList";
import { DataSourceCardList } from "../-components/DataSourceCardList";
import { DataSourceSectionLayout } from "./-components/layout/DataSourceSectionLayout";

const token = "placeholder";

export const Route = createFileRoute("/data-sources/$dataSourceId/")({
  component: RouteComponent,
  loader: async (context) => {
    const options: DefaultRequestOptions = { tenantId: TENANT_ID, token };
    const service = resolveHttpService(DataSourcesService);
    const dataSource = await service.getDataSource(
      context.params.dataSourceId,
      options,
    );

    return { dataSource };
  },
});

function RouteComponent() {
  const { dataSource } = Route.useLoaderData();
  const navigate = Route.useNavigate();
  const dataSourcesService = useDataSourcesService();

  const handleDeleteClk = async () => {
    try {
      await dataSourcesService.delete(dataSource!.id, {
        token,
        tenantId: TENANT_ID,
      });

      navigate({ to: "/data-sources" });
    } catch (error) {
      // FIXME: Properly handle error
    }
  };

  const handleSaveClk = async () => {
    try {
      await dataSourcesService.update(dataSource!.id, dataSource, {
        token,
        tenantId: TENANT_ID,
      });

      navigate({ to: "/data-sources" });
    } catch (error) {
      // FIXME: Properly handle error
    }
  };

  return (
    <Page sx={styles.page}>
      <Box sx={styles.header}>
        <Box sx={styles.headerTitle}>
          <PageTitle
            name={dataSource.name}
            description={dataSource.description}
          />
          <Box sx={styles.headerIcons}>
            <DataSourceTypeTag type={dataSource.type} />
          </Box>
        </Box>
        <Box sx={styles.headerActions}>
          <Button
            variant="text"
            onClick={handleDeleteClk}
            data-testid="delete-data-source-btn"
          >
            Delete
          </Button>
          <Button
            variant="text"
            onClick={handleSaveClk}
            data-testid="save-data-source-btn"
          >
            Save
          </Button>
        </Box>
      </Box>
      <DataSourceSectionLayout title="Where does the data come from?">
        <DataSourceCardList
          selection={dataSource!.type}
          onSelection={() => {}}
        />
      </DataSourceSectionLayout>
      <DataSourceAttributeForm dataSource={dataSource} />
      <TestBench />
    </Page>
  );
}
