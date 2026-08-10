import Box from "@mui/material/Box";
import { formFooterStyles } from "./FormFooter.style";
import Typography from "@mui/material/Typography";
import type { FormFooterActions } from "./FormFooterActions";
import { useForm } from "@/store/formDefinition";
import type { IFormProgress } from "@/utils/progress";
import LinearProgress from "@mui/material/LinearProgress";
import { FormProgress } from "../FormProgress";

export interface FormFooterProps {
  children: React.ReactElement<typeof FormFooterActions>;
  progress?: IFormProgress | undefined;
}

export const FormFooter: React.FC<FormFooterProps> = function ({
  children,
  progress,
}) {
  const form = useForm();

  let progressBar: React.ReactElement<typeof LinearProgress> | undefined =
    undefined;
  let text: React.ReactNode;

  if (progress) {
    progressBar = (
      <LinearProgress
        variant="determinate"
        sx={formFooterStyles["progressBar"]}
        value={progress.percentage}
      />
    );

    text = <FormProgress progress={progress} />;
  }

  return (
    <>
      <Box sx={formFooterStyles["footer"]}>
        {progressBar}
        <Box
          sx={{
            ...formFooterStyles["content"],
            paddingTop: progress ? 1 : 2.5,
          }}
        >
          <Box>
            {text}
            <Typography variant="body2" sx={formFooterStyles["name"]}>
              {form?.name}
            </Typography>
          </Box>
          {children}
        </Box>
      </Box>
    </>
  );
};
