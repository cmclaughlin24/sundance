import { Page } from "@/components/layout/Page/Page";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { designerStyles } from "./-index.styles";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import Box from "@mui/material/Box";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useFormsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";
import { FormList } from "./-components/FormList";
import type { IForm } from "@/types/form";
import { FormListToolbar } from "./-components/FormListToolbar";
import { useRef, useState } from "react";
import {
  CreateFormDrawer,
  type CreateFormEvent,
} from "./-components/CreateFormDrawer";
import type { DrawerHandle } from "@/components/Drawer";
import type { DefaultRequestOptions } from "@/services/baseHttpService";
import { defaultFormVersion } from "@/utils/form";

export const Route = createFileRoute("/designer/forms/")({
  component: DesignerRouteComponent,
});

function DesignerRouteComponent() {
  const formsService = useFormsService();
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const drawerRef = useRef<DrawerHandle>(null);

  const { data, isLoading, error } = useAsyncData(
    async (accessToken) => {
      if (!accessToken) {
        return null;
      }

      return await formsService.getForms({
        tenantId: TENANT_ID,
        token: accessToken,
      });
    },
    [formsService],
  );

  const handleFilterClick = () => {
    console.log("filters");
  };

  const handleFormClk = (form: IForm) => {
    navigate({
      to: "/designer/forms/$formId",
      params: { formId: form.id },
    });
  };

  const handleNewClk = () => {
    drawerRef.current?.open();
  };

  const handleCreate = async ({ name, description }: CreateFormEvent) => {
    try {
      const opts: DefaultRequestOptions = {
        token: "placeholder",
        tenantId: TENANT_ID,
      };

      const form = await formsService.createForm({ name, description }, opts);
      const version = await formsService.createFormVersion(
        form.id,
        defaultFormVersion(),
        opts,
      );

      navigate({
        to: "/designer/forms/$formId/versions/$versionId/{-$tab}",
        params: { formId: form.id, versionId: version.id },
      });
    } catch (error) {
      // TOOD: Implement error handling.
    }
  };

  if (isLoading) {
    return <>Loading forms...</>;
  }

  if (error) {
    return <>Something went wrong...</>;
  }

  return (
    <>
      <Page sx={designerStyles["page"]}>
        <PageTitle
          name="Forms Hub"
          description="View and manage request forms for your assets"
        />
        <Box>FORM COUNTS</Box>
        <Box>
          <FormListToolbar
            onNew={handleNewClk}
            search={search}
            onSearch={(value) => setSearch(value)}
          />
          <FormList forms={data} onClick={handleFormClk} />
        </Box>
      </Page>
      <CreateFormDrawer ref={drawerRef} onCreate={handleCreate} />
    </>
  );
}
