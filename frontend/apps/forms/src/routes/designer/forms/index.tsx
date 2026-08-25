import { Page } from "@/components/layout/Page/Page";
import { createFileRoute } from "@tanstack/react-router";
import { designerStyles } from "./-index.styles";
import { PageTitle } from "@/components/layout/Page/PageTitle";
import Box from "@mui/material/Box";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useFormsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";
import { FormList } from "./-components/FormList";
import type { IForm } from "@/types/form";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

export const Route = createFileRoute("/designer/forms/")({
  component: DesignerRouteComponent,
});

function DesignerRouteComponent() {
  const formsService = useFormsService();

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

  const handleFormClick = (form: IForm) => {
    console.log(form);
  };

  const handleNewClick = () => {
    console.log("new");
  };

  if (isLoading) {
    return <>Loading forms...</>;
  }

  if (error) {
    return <>Something went wrong...</>;
  }

  return (
    <Page sx={designerStyles["page"]}>
      <PageTitle
        name="Forms Hub"
        description="View and manage request forms for your assets"
      />
      <Box>FORM COUNTS</Box>
      <Box>
        <Box sx={designerStyles["toolbar"]}>
          <Box sx={designerStyles["toolbarLeft"]}>
            <Button onClick={handleFilterClick}>Filter</Button>
            <Button onClick={handleNewClick}>New Form +</Button>
          </Box>
          <Box>
            <TextField
              variant="outlined"
              placeholder="Search Forms"
              sx={designerStyles["searchInput"]}
            />
          </Box>
        </Box>
        <FormList forms={data} onClick={handleFormClick} />
      </Box>
    </Page>
  );
}
