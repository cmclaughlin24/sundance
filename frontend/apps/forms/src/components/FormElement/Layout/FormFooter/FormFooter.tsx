import Box from "@mui/material/Box";
import { formFooterStyles } from "./FormFooter.style";
import { useForm } from "@/store/useFormStoreContext";
import Typography from "@mui/material/Typography";
import type { FormFooterActions } from "./FormFooterActions";

export interface FormFooterProps {
  actions: React.ReactElement<typeof FormFooterActions>;
}

export const FormFooter: React.FC<FormFooterProps> = function ({ actions }) {
  const form = useForm();

  return (
    <Box sx={formFooterStyles["footer"]}>
      <Typography variant="body2" sx={formFooterStyles["name"]}>
        {form?.name}
      </Typography>
      {actions}
    </Box>
  );
};
