import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { FormTitle } from "../Layout/FormTitle";
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
import { formRendererStyles } from "./FormRenderer.style";
import { FormFooter } from "../Layout/FormFooter/FormFooter";
import { FormFooterActions } from "../Layout/FormFooter/FormFooterActions";

export interface FormRendererProps {
  onSubmit: (values: ISubmissionValue[]) => void;
  onCancel: () => void;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  onSubmit,
  onCancel,
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

    for (const [elementId, value] of Object.entries(values)) {
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
      <Box
        component="form"
        sx={formRendererStyles["form"]}
        onSubmit={handleSubmit}
        id={form.id}
      >
        <FormTitle name={form!.name} description={form!.description} />
        {pages.map((page) => (
          <PageRenderer page={page} key={page.id} />
        ))}
      </Box>
      <FormFooter
        actions={
          <FormFooterActions>
            <Button variant="contained" color="secondary" onClick={onCancel}>
              Cancel
            </Button>
            <Button variant="contained" type="submit" form={form.id}>
              Submit
            </Button>
          </FormFooterActions>
        }
      />
    </EvalContextContext>
  );
};
