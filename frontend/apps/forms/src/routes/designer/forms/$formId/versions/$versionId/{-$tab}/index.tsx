import { Page } from "@/components/layout/Page/Page";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import { TabPanel } from "@/components/layout/Tab/TabPanel";
import { TabPanelGroup } from "@/components/layout/Tab/TabPanelGroup";
import { TENANT_ID } from "@/constants/tenant";
import { resolveHttpService, useFormsService } from "@/hooks/useHttpService";
import { FormsService } from "@/services/formsService";
import Box from "@mui/material/Box";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { formDesignerPageStyles } from "./-index.style";
import Button from "@mui/material/Button";
import type { DefaultRequestOptions } from "@/services/baseHttpService";
import {
  FormDesignerProvider,
  useFormDesignerHistory,
  useFormSnapshot,
} from "@/store/formDesigner";
import { FormBuilder } from "@/components/FormDesigner/FormBuilder";
import { FormVersionTag } from "@/components/FormVersionStatusTag";
import { FormRules } from "@/components/FormDesigner/FormRules";
import { useAsyncData } from "@/hooks/useAsyncData";
import type { IFormVersion } from "@/types/formVersion";
import {
  copyVersion,
  defaultFormVersion,
  getLatestFormVersionByStatus,
  isActiveVersion,
  isDraftVersion,
  versionToRequest,
} from "@/utils/form";
import { FormTags } from "@/components/FormDesigner/FormTags/FormTags";

const token = "placeholder";

export const Route = createFileRoute(
  "/designer/forms/$formId/versions/$versionId/{-$tab}/",
)({
  component: RouteComponent,
  loader: async (context) => {
    const options: DefaultRequestOptions = { tenantId: TENANT_ID, token };
    const service = resolveHttpService(FormsService);
    const [form, version] = await service.getFormAndVersion(
      context.params.formId,
      context.params.versionId,
      options,
    );

    return { form, version };
  },
});

enum FormDesignerTab {
  Build = "build",
  Rules = "rules",
  Settings = "settings",
  Tags = "tags",
  Versions = "versions",
}

const TAB_ORDER = [
  FormDesignerTab.Build,
  FormDesignerTab.Rules,
  FormDesignerTab.Tags,
  FormDesignerTab.Versions,
  FormDesignerTab.Settings,
];

function RouteComponent() {
  const { form, version } = Route.useLoaderData();
  const { formId, versionId, tab } = Route.useParams();
  const navigate = useNavigate();

  const handleTabChange = (tab: FormDesignerTab) => {
    navigate({
      to: "/designer/forms/$formId/versions/$versionId/{-$tab}",
      params: { formId, versionId, tab },
    });
  };

  return (
    <FormDesignerProvider form={form} version={version!} key={versionId}>
      <PageComponent
        tab={tab as FormDesignerTab}
        onTabChange={handleTabChange}
      />
    </FormDesignerProvider>
  );
}

const PageComponent: React.FC<{
  tab?: FormDesignerTab;
  onTabChange: (tab: FormDesignerTab) => void;
}> = function ({ tab = FormDesignerTab.Build, onTabChange }) {
  const navigate = useNavigate();
  const { commit } = useFormDesignerHistory();
  const formsService = useFormsService();
  const { form, version, rules } = useFormSnapshot();

  const { data: versions, refetch: refetchVersions } = useAsyncData<
    IFormVersion[]
  >(
    async (token) => {
      return await formsService.getFormVersions(form.id, {
        tenantId: TENANT_ID,
        token,
      });
    },
    [form.id!],
  );

  const handleTabChange = (
    _event: React.SyntheticEvent,
    tab: FormDesignerTab,
  ) => onTabChange(tab);

  const saveDraft = async () => {
    if (!isDraftVersion(version.status)) {
      throw new Error("cannot update a non-draft version");
    }

    const request = versionToRequest(version, rules);
    return await formsService.updateFormVersion(form.id, version.id, request, {
      tenantId: TENANT_ID,
      token: "placeholder",
    });
  };

  const handleSaveDraft = async () => {
    try {
      const version = await saveDraft();
      commit(version);
    } catch (error) {
      // TODO: Implement error handling if the request to create the draft fails.
    }
  };

  const handlePublish = async (
    version: IFormVersion,
    updateStore: boolean = true,
  ) => {
    if (!isDraftVersion(version.status)) {
      throw new Error("cannot publish a non-draft version");
    }

    try {
      const updated = await formsService.publishFormVersion(
        form.id,
        version.id,
        {
          tenantId: TENANT_ID,
          token: "placeholder",
        },
      );

      updateStore && commit(updated);
      refetchVersions();
    } catch (error) {
      // TODO: Implement error handling if the request to create the draft fails.
    }
  };

  const handleRetire = async (
    version: IFormVersion,
    updateStore: boolean = true,
  ) => {
    if (!isActiveVersion(version.status)) {
      throw new Error("cannot retire a non-active version");
    }

    try {
      const updated = await formsService.retireFormVersion(
        form.id,
        version.id,
        {
          tenantId: TENANT_ID,
          token: "placeholder",
        },
      );

      updateStore && commit(updated);
      refetchVersions();
    } catch (error) {
      // TODO: Implement error handling if the request to create the draft fails.
    }
  };

  const handleNewDraft = async (base?: IFormVersion) => {
    const request = base ? copyVersion(base) : defaultFormVersion();

    try {
      const resp = await formsService.createFormVersion(form.id, request, {
        tenantId: TENANT_ID,
        token: "placeholder",
      });

      navigate({
        to: "/designer/forms/$formId/versions/$versionId/{-$tab}",
        params: { formId: resp.formId, versionId: resp.id },
      });
    } catch (error) {
      // TODO: Implement error handling if the request to create the draft fails.
    }
  };

  const latestActive = getLatestFormVersionByStatus(versions, "active");

  return (
    <Page sx={formDesignerPageStyles.page}>
      <Box sx={formDesignerPageStyles.header}>
        <Box sx={formDesignerPageStyles.headerTitle}>
          <PageTitle name={form.name} description={form.description} />
          <Box sx={formDesignerPageStyles.headerIcons}>
            <FormVersionTag status={version.status} version={version.version} />
            {latestActive && latestActive.id !== version.id && (
              <FormVersionTag
                status={latestActive.status}
                version={latestActive.version}
              />
            )}
          </Box>
        </Box>
        <Box sx={formDesignerPageStyles.headerActions}>
          <Button
            variant="text"
            onClick={handleSaveDraft}
            disabled={!isDraftVersion(version.status)}
          >
            Save Draft
          </Button>
          <Button>Preview</Button>
          {isActiveVersion(version.status) && (
            <Button onClick={() => handleRetire(version)}>Retire</Button>
          )}
          {isDraftVersion(version.status) && (
            <Button onClick={() => handlePublish(version)}>Publish</Button>
          )}
        </Box>
      </Box>
      <Box>
        <Tabs value={tab} onChange={handleTabChange}>
          <Tab label="Build" value={FormDesignerTab.Build} />
          <Tab label="Rules" value={FormDesignerTab.Rules} />
          <Tab label="Tags" value={FormDesignerTab.Tags} />
          <Tab label="Version" value={FormDesignerTab.Versions} />
          <Tab label="Settings" value={FormDesignerTab.Settings} />
        </Tabs>
      </Box>
      <TabPanelGroup active={tab} order={TAB_ORDER} sx={{ flex: 1 }}>
        <TabPanel value={FormDesignerTab.Build}>
          <FormBuilder />
        </TabPanel>
        <TabPanel value={FormDesignerTab.Rules}>
          <FormRules />
        </TabPanel>
        <TabPanel value={FormDesignerTab.Tags}>
          <FormTags />
        </TabPanel>
        <TabPanel value={FormDesignerTab.Versions}>Versions Tab</TabPanel>
        <TabPanel value={FormDesignerTab.Settings}>Settings Tab</TabPanel>
      </TabPanelGroup>
    </Page>
  );
};
