import { Page } from "@/components/layout/Page/Page";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import { TabPanel } from "@/components/layout/Tab/TabPanel";
import { TabPanelGroup } from "@/components/layout/Tab/TabPanelGroup";
import { TENANT_ID } from "@/constants/tenant";
import { resolveHttpService } from "@/hooks/useHttpService";
import { FormsService } from "@/services/formsService";
import Box from "@mui/material/Box";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { formDesignerPageStyles } from "./-index.style";
import Button from "@mui/material/Button";
import type { DefaultRequestOptions } from "@/services/baseHttpService";
import { FormDesignerProvider } from "@/store/formDesigner";
import { FormBuilder } from "@/components/FormDesigner/FormBuilder";
import { FormVersionTag } from "@/components/FormVersionStatusTag";

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
    const options: DefaultRequestOptions = { tenantId: TENANT_ID, token };
    const service = resolveHttpService(FormsService);
    const [form, versions] = await Promise.all([
      service.getForm(context.params.formId, options),
      service.getFormVersions(context.params.formId, options),
    ]);

    return { form, versions };
  },
});

function RouteComponent() {
  const navigate = useNavigate();
  const { formId, tab = FormDesignerTab.Build } = Route.useParams();
  const { form, versions } = Route.useLoaderData();

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
      <Box sx={formDesignerPageStyles.header}>
        <Box sx={formDesignerPageStyles.headerTitle}>
          <PageTitle name={form.name} description={form.description} />
          <Box sx={formDesignerPageStyles.headerIcons}>
            <FormVersionTag status="draft" text="Draft v2" />
            <FormVersionTag status="active" text="Active v1" />
          </Box>
        </Box>
        <Box sx={formDesignerPageStyles.headerActions}>
          <Button variant="text">Save Draft</Button>
          <Button>Preview</Button>
          <Button>Publish</Button>
        </Box>
      </Box>
      <Box>
        <Tabs value={tab} onChange={handleTabChange}>
          <Tab label="Build" value={FormDesignerTab.Build} />
          <Tab label="Rules" value={FormDesignerTab.Rules} />
          <Tab label="Reference Data" value={FormDesignerTab.DataSources} />
          <Tab label="Version" value={FormDesignerTab.Versions} />
          <Tab label="Settings" value={FormDesignerTab.Settings} />
        </Tabs>
      </Box>
      <FormDesignerProvider form={form} version={versions[0]}>
        <TabPanelGroup active={tab} order={TAB_ORDER}>
          <TabPanel value={FormDesignerTab.Build}>
            <FormBuilder />
          </TabPanel>
          <TabPanel value={FormDesignerTab.Rules}>Rules Tab</TabPanel>
          <TabPanel value={FormDesignerTab.DataSources}>
            Reference Data Tab
          </TabPanel>
          <TabPanel value={FormDesignerTab.Versions}>Versions Tab</TabPanel>
          <TabPanel value={FormDesignerTab.Settings}>Settings Tab</TabPanel>
        </TabPanelGroup>
      </FormDesignerProvider>
    </Page>
  );
}
