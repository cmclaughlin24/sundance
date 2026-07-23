import { FormTitle } from "../FormTitle";
import { sortPositioned } from "@/utils/sort";
import type { ISubmissionValue } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import type { SubmitEvent } from "react";
import { useFormState } from "@/store/useFormContext";
import { filterVisible } from "@/utils/filter";

export interface FormRendererProps {
  onSubmit: (values: ISubmissionValue[]) => void;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  onSubmit,
}) {
  const state = useFormState();
  const { form, version } = state;

  if (!form || !version) {
    return <></>;
  }

  let pages = sortPositioned(version!.pages);
  pages = filterVisible(pages, state);

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    const values: ISubmissionValue[] = [];

    for (const [elementId, value] of state.values) {
      values.push({ elementId, value });
    }

    onSubmit(values);
  };

  return (
    <>
      <FormTitle name={form!.name} description={form!.description} />
      <form onSubmit={handleSubmit}>
        {pages.map((page) => (
          <PageRenderer page={page} key={page.id} />
        ))}
        <button type="submit">submit</button>
      </form>
    </>
  );
};
