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
          tenantId="019fbd44-3fd5-727a-9ff0-40eade1e297a"
          formId="019fbd45-6995-7457-a072-370251975daa"
          versionId="019fbdad-fb8a-7386-a070-1b9b5077c527"
          onSubmit={handleSubmit}
          onCancel={handleCancel}
        />
      </MainContainer>
    </div>
  );
}
