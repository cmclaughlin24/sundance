import Box from "@mui/material/Box";
import { formFooterStyles } from "./FormFooter.style";
import Typography from "@mui/material/Typography";
import type { FormFooterActions } from "./FormFooterActions";
import type LinearProgress from "@mui/material/LinearProgress";
import { useForm } from "@/store/formDefinition";

export interface FormFooterProps {
  actions: React.ReactElement<typeof FormFooterActions>;
  progress?: React.ReactElement<typeof LinearProgress>;
}

export const FormFooter: React.FC<FormFooterProps> = function ({
  actions,
  progress,
}) {
  const form = useForm();

  return (
    <>
      <Box sx={formFooterStyles["footer"]}>
        {progress}
        <Box sx={formFooterStyles["content"]}>
          <Typography variant="body2" sx={formFooterStyles["name"]}>
            {form?.name}
          </Typography>
          {actions}
        </Box>
      </Box>
    </>
  );
};
