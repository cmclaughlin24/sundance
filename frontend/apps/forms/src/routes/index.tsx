import { FormElement } from "@/components/FormElement/FormElement";
import type { IFormElementSubmitEvent } from "@/components/FormElement/FormElement.type";
import { MainContainer } from "@/components/MainContainer/MainContainer";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const handleSubmit = (event: IFormElementSubmitEvent) => console.log(event);

  return (
    <div>
      <MainContainer>
        <FormElement
          tenantId="019f8b42-c81a-7c4f-afaa-beb2b04f9ef6"
          formId="019f8b43-712c-7788-ab1e-0728bae405f4"
          versionId="019f8b43-f86c-7e4b-9046-1937037154af"
          onSubmit={handleSubmit}
        />
      </MainContainer>
    </div>
  );
}
