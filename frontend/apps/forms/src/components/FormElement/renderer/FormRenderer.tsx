import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { PageTitle } from "../../layout/Page/PageTitle";
import { sortPositioned } from "@/utils/position";
import type { ISubmissionSection } from "@/types/submission";
import { PageRenderer } from "./PageRenderer";
import { useMemo, useRef, useState, type SubmitEvent } from "react";
import { useFormValues } from "@/store/submission/useSubmissionContext";
import { filterVisible } from "@/utils/filter";
import { EvalContextContext } from "@/store/submission/evalContext";
import { buildFieldEvalNamespace, type EvalContext } from "@/utils/evaluate";
import { rendererStyles } from "./renderer.style";
import { useForm, useFormVersion } from "@/store/formDefinition";
import { AnimatePresence } from "motion/react";
import { calculateProgress } from "@/utils/progress";
import { checkFormValidity } from "@/utils/form";
import { FormFooter } from "../layout/FormFooter/FormFooter";
import { FormFooterActions } from "../layout/FormFooter/FormFooterActions";
import type { FormRendererSlotProps } from "../FormElement.type";
import { getFormRendererSlotProp } from "./rendererDefaults";

export interface FormRendererProps {
  onSubmit: (values: ISubmissionSection[]) => void;
  onCancel: () => void;
  slotProps?: FormRendererSlotProps;
}

export const FormRenderer: React.FC<FormRendererProps> = function ({
  onSubmit,
  onCancel,
  slotProps,
}) {
  const formRef = useRef<HTMLFormElement>(null);
  const form = useForm();
  const version = useFormVersion();
  const formValues = useFormValues();
  const [pageIndex, setPageIndex] = useState(0);
  const evalCtx = useMemo<EvalContext>(() => {
    const pages = version?.pages ?? [];
    const field = buildFieldEvalNamespace(pages, formValues);
    // TODO: Add userClaims namespace.
    return { field };
  }, [version, formValues]);

  let pages = sortPositioned(version!.pages);
  pages = filterVisible(pages, evalCtx);

  // NOTE: While it is protected against during form design time, if the current page becomes invisible
  // the previous page will be show to the user.
  const page = pages[pageIndex] ?? pages[pageIndex - 1];
  const progress = calculateProgress(formValues, evalCtx, pages);
  const isLastPage = pageIndex === pages.length - 1;
  const isFormValid = checkFormValidity(formRef.current, progress);

  const handleSubmit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    const submission: ISubmissionSection[] = [];

    for (const sectionId in formValues) {
      submission.push({
        sectionId,
        instances: formValues[sectionId].map((i) => ({
          values: i.toSubmissionValues(),
        })),
      });
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
        ref={formRef}
      >
        <PageTitle name={form!.name} description={form!.description} />
        <AnimatePresence mode="wait">
          <PageRenderer page={page} key={page.id} />
        </AnimatePresence>
      </Box>
      <FormFooter progress={progress} variant="compact">
        <FormFooterActions>
          {getFormRendererSlotProp("buttons.cancel.show", slotProps) && (
            <Button variant="contained" color="secondary" onClick={onCancel}>
              {getFormRendererSlotProp("buttons.cancel.label", slotProps)}
            </Button>
          )}
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
          {isLastPage &&
            getFormRendererSlotProp("buttons.submit.show", slotProps) && (
              <Button
                variant="contained"
                type="submit"
                disabled={!isFormValid}
                form={form.id}
              >
                {getFormRendererSlotProp("buttons.submit.label", slotProps)}
              </Button>
            )}
        </FormFooterActions>
      </FormFooter>
    </EvalContextContext>
  );
};
