import { FormTitle } from "../FormTitle";
import { sortPositioned } from "@/utils/sort";
import type { ISubmissionValue } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import { useMemo, type SubmitEvent } from "react";
import {
  useForm,
  useFormValues,
  useFormVersion,
} from "@/store/useFormStoreContext";
import { filterVisible } from "@/utils/filter";
import { EvalContextContext } from "@/store/evalContext";
import { buildEvalContext, type EvalContext } from "@/utils/evaluate";

export interface FormRendererProps {
  onSubmit: (values: ISubmissionValue[]) => void;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  onSubmit,
}) {
  const form = useForm();
  const version = useFormVersion();
  const values = useFormValues();
  const evalCtx = useMemo<EvalContext>(() => {
    const pages = version?.pages ?? [];
    return buildEvalContext(pages, values);
  }, [version, values]);

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    const submission: ISubmissionValue[] = [];

    for (const [elementId, value] of Object.entries(submission)) {
      submission.push({ elementId, value });
    }

    onSubmit(submission);
  };

  if (!form || !version) {
    return <>Missing form and version</>;
  }

  let pages = sortPositioned(version!.pages);
  pages = filterVisible(pages, evalCtx);

  return (
    <EvalContextContext value={evalCtx}>
      <FormTitle name={form!.name} description={form!.description} />
      <form onSubmit={handleSubmit}>
        {pages.map((page) => (
          <PageRenderer page={page} key={page.id} />
        ))}
        <button type="submit">submit</button>
      </form>
    </EvalContextContext>
  );
};
