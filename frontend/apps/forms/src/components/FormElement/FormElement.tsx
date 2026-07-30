import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type {
  FormElementProps,
  IAsyncSubmitEvent,
  ISyncSubmitEvent,
} from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";
import { FormStoreProvider } from "@/store/FormStoreProvider";
import Box from "@mui/material/Box";
import { FormRenderer } from "./Renderer/FormRenderer";
import type { ISubmissionValue } from "@/types/submission";
import { formElementStyles } from "./FormElement.style";

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
  const accessToken = "placeholder";

  const { data, isLoading, error } = useAsyncData(async () => {
    if (!accessToken) {
      return null;
    }

    return await formsService.getForm(formId, versionId, {
      tenantId,
      token: accessToken,
    });
  }, [formsService, tenantId, formId, versionId, accessToken]);

  const asyncSubmit = async (
    values: ISubmissionValue[],
  ): Promise<IAsyncSubmitEvent> => {
    const result = await submissionService.submit(formId, versionId, values, {
      tenantId,
      token: accessToken,
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
      {
        tenantId,
        token: accessToken,
      },
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

  return (
    <FormStoreProvider
      form={form}
      version={version}
      rawSubmission={rawSubmission}
    >
      <Box sx={formElementStyles["container"]}>
        <FormRenderer onSubmit={handleSubmit} onCancel={onCancel} />
      </Box>
    </FormStoreProvider>
  );
};
