import Box from "@mui/material/Box";
import { formFooterStyles } from "./FormFooter.style";
import type { FormFooterActions } from "./FormFooterActions";
import { useForm } from "@/store/formDefinition";
import type { IFormProgress } from "@/utils/progress";
import LinearProgress from "@mui/material/LinearProgress";
import type { IForm } from "@/types/form";
import { FullFormFooterContent } from "./FullFormFooterContent";
import { CompactFormFooterContent } from "./CompactFormFooterContent";

export interface FormFooterProps {
  children: React.ReactElement<typeof FormFooterActions>;
  variant?: "full" | "compact";
  progress?: IFormProgress | undefined;
}

export type FormContentComponent = React.FC<
  { form: IForm | null } & Exclude<FormFooterProps, "variant">
>;

export const FormFooter: React.FC<FormFooterProps> = function ({
  children,
  variant,
  progress,
}) {
  const form = useForm();

  let progressBar: React.ReactElement<typeof LinearProgress> | undefined =
    undefined;

  if (progress) {
    progressBar = (
      <LinearProgress
        variant="determinate"
        sx={formFooterStyles["progressBar"]}
        value={progress.percentage}
      />
    );
  }

  let Content: FormContentComponent = FullFormFooterContent;

  if (variant === "compact") {
    Content = CompactFormFooterContent;
  }

  return (
    <>
      <Box sx={formFooterStyles["footer"]}>
        {progressBar}
        <Content form={form} progress={progress}>
          {children}
        </Content>
      </Box>
    </>
  );
};
