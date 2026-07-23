import { FormTitle } from "../FormTitle";
import { sortPositioned } from "@/utils/sort";
import type { ISubmissionValue } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import { useMemo, type SubmitEvent } from "react";
import { useFormState } from "@/store/useFormContext";
import { filterVisible } from "@/utils/filter";
import { EvalContextContext } from "@/store/evalContext";
import { buildEvalContext, type EvalContext } from "@/utils/evaluate";

export interface FormRendererProps {
  onSubmit: (values: ISubmissionValue[]) => void;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  onSubmit,
}) {
  const state = useFormState();
  const evalCtx = useMemo<EvalContext>(() => {
    const pages = state.version?.pages ?? [];
    return buildEvalContext(pages, state.values);
  }, [state.version, state.values]);

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    const values: ISubmissionValue[] = [];

    for (const [elementId, value] of Object.entries(state.values)) {
      values.push({ elementId, value });
    }

    onSubmit(values);
  };

  const { form, version } = state;

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
