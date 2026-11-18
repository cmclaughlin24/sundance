import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type {
  FormElementProps,
  IAsyncSubmitEvent,
  ISyncSubmitEvent,
} from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";
import { SubmissionProvider } from "@/store/submission/SubmissionProvider";
import { FormRenderer } from "./renderer/FormRenderer";
import type { ISubmissionSection, ISubmissionValue } from "@/types/submission";
import { formElementStyles } from "./FormElement.style";
import { FormDefinitionProvider } from "@/store/formDefinition/FormDefinitionProvider";
import { Page } from "../layout/Page/Page";
import { hydrateSubmissionValues } from "@/utils/form";
import { useMemo } from "react";

export const FormElement: React.FC<FormElementProps> = function ({
  tenantId,
  formId,
  versionId,
  submitType,
  rawSubmission,
  onSubmit,
  onCancel,
}) {
  const formsService = useFormsService();
  const submissionService = useSubmissionsService();
  const token = "placeholder";

  const { data, isLoading, error } = useAsyncData(
    async (accessToken) => {
      if (!accessToken) {
        return null;
      }

      return await formsService.getFormAndVersion(formId, versionId, {
        tenantId,
        token: accessToken,
      });
    },
    [formsService, tenantId, formId, versionId],
  );

  const asyncSubmit = async (
    sections: ISubmissionSection[],
  ): Promise<IAsyncSubmitEvent> => {
    const result = await submissionService.submit(formId, versionId, sections, {
      tenantId,
      token,
    });

    return { referenceId: result.referenceId };
  };

  const syncSubmit = async (
    sections: ISubmissionSection[],
  ): Promise<ISyncSubmitEvent> => {
    const result = await submissionService.normalize(
      formId,
      versionId,
      sections,
      {
        tenantId,
        token,
      },
    );

    return { raw: values, normalized: result };
  };

  const handleSubmit = async (_values: ISubmissionValue[]) => {
    // FIXME: Correctly pass values
    try {
      if (submitType === "async") {
        const result = await asyncSubmit([]);
        onSubmit(result);
      } else if (!submitType || submitType === "sync") {
        const result = await syncSubmit([]);
        onSubmit(result);
      }
    } catch (error) {}
  };

  const values = useMemo(
    () =>
      data && data.length == 2
        ? hydrateSubmissionValues(rawSubmission, data[1])
        : [],
    [rawSubmission, data],
  );

  if (isLoading) {
    return <>Loading the form...</>;
  }

  if (error) {
    return <>Something went wrong...</>;
  }

  if (!data) {
    return <>Not found...</>;
  }

  const [form, version] = data;

  return (
    <FormDefinitionProvider form={form} version={version}>
      <SubmissionProvider rawSubmission={values}>
        <Page sx={formElementStyles["page"]}>
          <FormRenderer onSubmit={handleSubmit} onCancel={onCancel} />
        </Page>
      </SubmissionProvider>
    </FormDefinitionProvider>
  );
};
