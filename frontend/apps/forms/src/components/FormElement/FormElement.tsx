import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type { FormElementProps } from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";
import { FormProvider } from "@/store/FormProvider";

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

  const { isLoading, error } = useAsyncData(async () => {
    if (!accessToken) {
      return null;
    }

    return await formsService.getForm(formId, versionId, {
      tenantId,
      token: accessToken,
    });
  }, [formsService, tenantId, formId, versionId, accessToken]);

  const handleSubmit = async (event: React.SubmitEvent<any>) => {
    event.preventDefault();

    try {
      const result = await submissionService.normalize(formId, versionId, [], {
        tenantId,
        token: accessToken,
      });

      onSubmit({ raw: [], normalized: result });
    } catch (error) {}
  };

  if (isLoading) {
    return <>Loading the form...</>;
  }

  if (error) {
    return <>Something went wrong...</>;
  }

  return (
    <FormProvider rawSubmission={rawSubmission}>
      <form onSubmit={handleSubmit}></form>
    </FormProvider>
  );
};
