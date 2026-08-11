import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { PageTitle } from "../../Layout/Page/PageTitle";
import { sortPositioned } from "@/utils/sort";
import type { ISubmissionValue } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import { useMemo, useState, type SubmitEvent } from "react";
import {
  useFormErrors,
  useFormValues,
} from "@/store/submission/useSubmissionContext";
import { filterVisible } from "@/utils/filter";
import { EvalContextContext } from "@/store/submission/evalContext";
import { buildEvalContext, type EvalContext } from "@/utils/evaluate";
import { rendererStyles } from "./renderer.style";
import { FormFooter } from "../Layout/FormFooter/FormFooter";
import { FormFooterActions } from "../Layout/FormFooter/FormFooterActions";
import { useForm, useFormVersion } from "@/store/formDefinition";
import { AnimatePresence } from "motion/react";
import { calculateProgress } from "@/utils/progress";

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
  const errors = useFormErrors();
  const [pageIndex, setPageIndex] = useState(0);
  const evalCtx = useMemo<EvalContext>(() => {
    const pages = version?.pages ?? [];
    return buildEvalContext(pages, values);
  }, [version, values]);

  let pages = sortPositioned(version!.pages);
  pages = filterVisible(pages, evalCtx);

  // NOTE: While it is protected against during form design time, if the current page becomes invisible
  // the previous page will be show to the user.
  const page = pages[pageIndex] ?? pages[pageIndex - 1];
  const progress = calculateProgress(values, errors, evalCtx, pages);
  const isLastPage = pageIndex === pages.length - 1;

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    const submission: ISubmissionValue[] = [];

    for (const [elementId, value] of Object.entries(values)) {
      submission.push({ elementId, value });
    }

    onSubmit(submission);
  };

  const handlePageChange = (increment: number) => {
    if (increment < 0) {
      setPageIndex((idx) => Math.max(idx + increment, 0));
    } else if (increment > 0) {
      setPageIndex((idx) => Math.min(idx + increment, pages.length - 1));
    }
  };

  if (!form || !version) {
    return <>Missing form and version</>;
  }

  return (
    <EvalContextContext value={evalCtx}>
      <Box
        component="form"
        sx={rendererStyles["form"]}
        onSubmit={handleSubmit}
        id={form.id}
      >
        <PageTitle name={form!.name} description={form!.description} />
        <AnimatePresence mode="wait">
          <PageRenderer page={page} key={page.id} />
        </AnimatePresence>
      </Box>
      <FormFooter progress={progress} variant="compact">
        <FormFooterActions>
          <Button variant="contained" color="secondary" onClick={onCancel}>
            Cancel
          </Button>
          {pages.length > 1 && pageIndex > 0 && (
            <Button variant="outlined" onClick={() => handlePageChange(-1)}>
              Previous
            </Button>
          )}
          {!isLastPage && (
            <Button variant="outlined" onClick={() => handlePageChange(1)}>
              Next
            </Button>
          )}
          {isLastPage && (
            <Button variant="contained" type="submit" form={form.id}>
              Submit
            </Button>
          )}
        </FormFooterActions>
      </FormFooter>
    </EvalContextContext>
  );
};
