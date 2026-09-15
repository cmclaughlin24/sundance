import { TENANT_ID } from "@/constants/tenant";
import { resolveHttpService } from "@/hooks/useHttpService";
import type { DefaultRequestOptions } from "@/services/baseHttpService";
import { FormsService } from "@/services/formsService";
import type { IFormVersion } from "@/types/formVersion";
import { sortFormVersions, getLatestVersionByStatus } from "@/utils/form";
import { createFileRoute, notFound, redirect } from "@tanstack/react-router";

function pickDefaultVersion(versions: IFormVersion[]): IFormVersion {
  const sorted = sortFormVersions(versions);

  const highestDraft = getLatestVersionByStatus(sorted, "draft");
  if (highestDraft) {
    return highestDraft;
  }

  const highestActive = getLatestVersionByStatus(sorted, "active");
  if (highestActive) {
    return highestActive;
  }

  return sorted[0];
}

export const Route = createFileRoute("/designer/forms/$formId/")({
  beforeLoad: async ({ params }) => {
    const { formId } = params;
    const options: DefaultRequestOptions = {
      tenantId: TENANT_ID,
      token: "placeholder",
    };
    const service = resolveHttpService(FormsService);

    const versions = await service.getFormVersions(formId, options);
    const activeVersion = pickDefaultVersion(versions);

    if (!activeVersion) {
      throw notFound();
    }

    throw redirect({
      to: "/designer/forms/$formId/versions/$versionId/{-$tab}",
      params: {
        formId,
        versionId: activeVersion.id,
      },
      replace: true,
    });
  },
});
