import { Page } from "@/components/Layout/Page/Page";
import { createFileRoute } from "@tanstack/react-router";
import { designerStyles } from "./-index.styles";
import { PageTitle } from "@/components/Layout/Page/PageTitle";

export const Route = createFileRoute("/designer/")({
  component: DesignerRouteComponent,
});

function DesignerRouteComponent() {
  return (
    <Page sx={designerStyles["page"]}>
      <PageTitle
        name="Forms Hub"
        description="View and manage request forms for your assets"
      />
    </Page>
  );
}
