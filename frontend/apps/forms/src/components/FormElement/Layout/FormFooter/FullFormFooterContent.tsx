import Typography from "@mui/material/Typography";
import { FormProgress } from "../FormProgress/FormProgress";
import type { FormContentComponent } from "./FormFooter";
import { formFooterStyles } from "./FormFooter.style";
import Box from "@mui/material/Box";

export const FullFormFooterContent: FormContentComponent = function ({
  form,
  progress,
  children,
}) {
  const text: React.ReactNode = !progress ? null : (
    <FormProgress progress={progress} />
  );

  return (
    <Box
      sx={{
        ...formFooterStyles["fullContent"],
        paddingTop: progress ? 1 : 2.5,
      }}
    >
      <Box sx={formFooterStyles["statusMessage"]}>
        {text}
        <Typography variant="body2" sx={formFooterStyles["name"]}>
          {form?.name}
        </Typography>
      </Box>
      {children}
    </Box>
  );
};
