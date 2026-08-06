import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type {
  FormElementProps,
  IAsyncSubmitEvent,
  ISyncSubmitEvent,
} from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";
import { SubmissionProvider } from "@/store/submission/SubmissionProvider";
import { FormRenderer } from "./Renderer/FormRenderer";
import type { ISubmissionValue } from "@/types/submission";
import { formElementStyles } from "./FormElement.style";
import { FormDefinitionProvider } from "@/store/formDefinition/FormDefinitionProvider";
import { Page } from "../Layout/Page/Page";
import { hydrateSubmissionValues } from "@/utils/form";

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

      return await formsService.getForm(formId, versionId, {
        tenantId,
        token: accessToken,
      });
    },
    [formsService, tenantId, formId, versionId],
  );

  const asyncSubmit = async (
    values: ISubmissionValue[],
  ): Promise<IAsyncSubmitEvent> => {
    const result = await submissionService.submit(formId, versionId, values, {
      tenantId,
      token,
    });

    return { referenceId: result.referenceId };
  };

  const syncSubmit = async (
    values: ISubmissionValue[],
  ): Promise<ISyncSubmitEvent> => {
    const result = await submissionService.normalize(
      formId,
      versionId,
      values,
      { tenantId, token },
    );

    return { raw: values, normalized: result };
  };

  const handleSubmit = async (values: ISubmissionValue[]) => {
    try {
      if (submitType === "async") {
        const result = await asyncSubmit(values);
        onSubmit(result);
      } else if (!submitType || submitType === "sync") {
        const result = await syncSubmit(values);
        onSubmit(result);
      }
    } catch (error) {}
  };

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
  const values = hydrateSubmissionValues(rawSubmission, version);

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
