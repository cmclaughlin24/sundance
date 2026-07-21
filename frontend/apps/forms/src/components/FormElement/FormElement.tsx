import { useFormsService, useSubmissionsService } from "@/hooks/useHttpService";
import type { FormElementProps } from "./FormElement.type";
import { useAsyncData } from "@/hooks/useAsyncData";

export const FormElement: React.FC<FormElementProps> = function ({
  tenantId,
  formId,
  versionId,
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

  console.log(data);

  return <form onSubmit={handleSubmit}></form>;
};
