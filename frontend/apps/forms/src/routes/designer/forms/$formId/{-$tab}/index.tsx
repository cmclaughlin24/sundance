import { FormDesigner } from "@/components/FormDesigner/FormDesigner";
import { Page } from "@/components/layout/Page/Page";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import { TabPanel } from "@/components/layout/Tab/TabPanel";
import { TabPanelGroup } from "@/components/layout/Tab/TabPanelGroup";
import { TENANT_ID } from "@/constants/tenant";
import { resolveHttpService } from "@/hooks/useHttpService";
import { FormsService } from "@/services/formsService";
import type { IFormVersion } from "@/types/formVersion";
import Box from "@mui/material/Box";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { formDesignerPageStyles } from "./-index.style";

const token = "placeholder";

enum FormDesignerTab {
  Build = "build",
  Rules = "rules",
  DataSources = "dataSources",
  Versions = "versions",
  Settings = "settings",
}

const TAB_ORDER = [
  FormDesignerTab.Build,
  FormDesignerTab.Rules,
  FormDesignerTab.DataSources,
  FormDesignerTab.Versions,
  FormDesignerTab.Settings,
];

export const Route = createFileRoute("/designer/forms/$formId/{-$tab}/")({
  component: RouteComponent,
  loader: async (context) => {
    const service = resolveHttpService(FormsService);
    const form = await service.getForm(context.params.formId, {
      tenantId: TENANT_ID,
      token,
    });

    return { form };
  },
});

function RouteComponent() {
  const navigate = useNavigate();
  const { formId, tab = FormDesignerTab.Build } = Route.useParams();
  const { form } = Route.useLoaderData();

  const handleTabChange = (
    _event: React.SyntheticEvent,
    tab: FormDesignerTab,
  ) => {
    navigate({
      to: "/designer/forms/$formId/{-$tab}",
      params: { formId, tab },
    });
  };

  return (
    <Page sx={formDesignerPageStyles.page}>
      <PageTitle name={form.name} description={form.description} />
      <Box>
        <Tabs value={tab} onChange={handleTabChange}>
          <Tab label="Build" value={FormDesignerTab.Build} />
          <Tab label="Rules" value={FormDesignerTab.Rules} />
          <Tab label="Reference Data" value={FormDesignerTab.DataSources} />
          <Tab label="Version" value={FormDesignerTab.Versions} />
          <Tab label="Settings" value={FormDesignerTab.Settings} />
        </Tabs>
      </Box>
      <TabPanelGroup active={tab} order={TAB_ORDER}>
        <TabPanel value={FormDesignerTab.Build}>
          <FormDesigner form={form} version={{} as IFormVersion} />
        </TabPanel>
        <TabPanel value={FormDesignerTab.Rules}>Rules Tab</TabPanel>
        <TabPanel value={FormDesignerTab.DataSources}>
          Reference Data Tab
        </TabPanel>
        <TabPanel value={FormDesignerTab.Versions}>Versions Tab</TabPanel>
        <TabPanel value={FormDesignerTab.Settings}>Settings Tab</TabPanel>
      </TabPanelGroup>
    </Page>
  );
}
