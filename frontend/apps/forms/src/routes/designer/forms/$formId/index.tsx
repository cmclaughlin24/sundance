import { TENANT_ID } from "@/constants/tenant";
import { resolveHttpService } from "@/hooks/useHttpService";
import { FormsService } from "@/services/formsService";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/designer/forms/$formId/")({
  component: RouteComponent,
  loader: async (context) => {
    const service = resolveHttpService(FormsService);
    const form = await service.getFormAndVersion(context.params.formId, "", {
      tenantId: TENANT_ID,
      token: "",
    });

    return { form };
  },
});

function RouteComponent() {
  return <div>Hello "/designer/forms/$formId/"!</div>;
}
