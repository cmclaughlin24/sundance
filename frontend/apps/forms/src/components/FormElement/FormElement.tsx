import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type { FormElementProps } from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";
import { FormProvider } from "@/store/FormProvider";
import Box from "@mui/material/Box";
import { FormRenderer } from "./Renderer/FormRenderer";
import type { ISubmissionValue } from "@/types/submission";

export const FormElement: React.FC<FormElementProps> = function ({
  tenantId,
  formId,
  versionId,
  rawSubmission,
  onSubmit,
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

  const handleSubmit = async (values: ISubmissionValue[]) => {
    try {
      const result = await submissionService.normalize(
        formId,
        versionId,
        values,
        {
          tenantId,
          token: accessToken,
        },
      );

      onSubmit({ raw: values, normalized: result });
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
    <FormProvider rawSubmission={rawSubmission}>
      <Box>
        <FormRenderer form={form} version={version} onSubmit={handleSubmit} />
      </Box>
    </FormProvider>
  );
};
