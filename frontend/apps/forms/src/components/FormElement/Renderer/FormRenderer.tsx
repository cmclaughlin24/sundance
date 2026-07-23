import type { IForm } from "@/types/form";
import { FormTitle } from "../FormTitle";
import type { IFormVersion } from "@/types/formVersion";
import { sortFormElements } from "@/utils/sort";
import type { ISubmissionValue } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import type { SubmitEvent } from "react";

export interface FormRendererProps {
  form: IForm;
  version: IFormVersion;
  onSubmit: (values: ISubmissionValue[]) => void;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  form,
  version,
  onSubmit,
}) {
  const pages = sortFormElements(version.pages);

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    onSubmit([]);
  };

  return (
    <>
      <FormTitle name={form.name} description={form.description} />
      <form onSubmit={handleSubmit}>
        {pages.map((page) => (
          <PageRenderer page={page} key={page.id} />
        ))}
      </form>
    </>
  );
};
