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
import { DataSourceTypeTag } from "../../-components/DataSourceTypeTag";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { DataSourceAttributeForm } from "./-components/DataSourceAttributeForm/DataSourceAttributeForm";
import { TestBench } from "./-components/TestBench";
import { DataSourceCardList } from "../../-components/DataSourceCardList";
import { DataSourceSectionLayout } from "./-components/layout/DataSourceSectionLayout";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import type { DataSource } from "@/types/dataSource";
import { TabPanelGroup } from "@/components/layout/Tab/TabPanelGroup";
import { TabPanel } from "@/components/layout/Tab/TabPanel";

const token = "placeholder";

export const Route = createFileRoute("/data-sources/$dataSourceId/{-$tab}/")({
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

enum DataSourceTab {
  Configure = "configure",
  TestBench = "test-bench",
}

const TAB_ORDER = [DataSourceTab.Configure, DataSourceTab.TestBench];

function RouteComponent() {
  const { dataSource } = Route.useLoaderData();
  const { tab } = Route.useParams();
  const navigate = Route.useNavigate();

  const handleTabChange = (tab: DataSourceTab) => {
    navigate({
      to: "/data-sources/$dataSourceId/{-$tab}",
      params: { dataSourceId: dataSource.id, tab },
    });
  };

  return (
    <PageComponent
      dataSource={dataSource}
      tab={tab as DataSourceTab}
      onTabChange={handleTabChange}
    />
  );
}

const PageComponent: React.FC<{
  dataSource: DataSource;
  tab?: DataSourceTab;
  onTabChange: (tab: DataSourceTab) => void;
}> = function ({ dataSource, tab = DataSourceTab.Configure, onTabChange }) {
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

  const handleTabChange = (_event: React.SyntheticEvent, tab: DataSourceTab) =>
    onTabChange(tab);

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
      <Tabs value={tab} onChange={handleTabChange}>
        <Tab label="Configure" value={DataSourceTab.Configure} />
        <Tab label="Test Bench" value={DataSourceTab.TestBench} />
      </Tabs>
      <TabPanelGroup active={tab} order={TAB_ORDER} sx={{ flex: 1 }}>
        <TabPanel value={DataSourceTab.Configure}>
          <Box sx={styles.configureTab}>
            <DataSourceSectionLayout title="Where does the data come from?">
              <DataSourceCardList
                selection={dataSource!.type}
                onSelection={() => {}}
              />
            </DataSourceSectionLayout>
            <DataSourceAttributeForm dataSource={dataSource} />
          </Box>
        </TabPanel>
        <TabPanel value={DataSourceTab.TestBench}>
          <TestBench />
        </TabPanel>
      </TabPanelGroup>
    </Page>
  );
};
