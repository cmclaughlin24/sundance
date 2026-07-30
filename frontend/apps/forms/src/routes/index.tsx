import { FormElement } from "@/components/FormElement/FormElement";
import type { ISyncSubmitEvent } from "@/components/FormElement/FormElement.type";
import { MainContainer } from "@/components/MainContainer/MainContainer";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const handleSubmit = (event: ISyncSubmitEvent) => console.log(event);

  const handleCancel = () => console.log("cancel");

  return (
    <div>
      <MainContainer>
        <FormElement
          tenantId="019f9f3e-ff1d-7363-9ade-9604833c0d9e"
          formId="019f9f40-22de-7c12-b867-d22c77470217"
          versionId="019f9f40-4df0-7a80-bd84-4c947d09aac2"
          onSubmit={handleSubmit}
          onCancel={handleCancel}
        />
      </MainContainer>
    </div>
  );
}
