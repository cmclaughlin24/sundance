import type { FormContentComponent } from "./FormFooter";
import { formFooterStyles } from "./FormFooter.style";
import Box from "@mui/material/Box";
import { FormProgressMessage } from "../FormProgress/FormProgressMessage";

export const CompactFormFooterContent: FormContentComponent = function ({
  progress,
  children,
}) {
  let statusMessage =
    !progress || progress.percentage === 100 ? null : (
      <Box component="span" sx={formFooterStyles["statusMessage"]}>
        <FormProgressMessage progress={progress} />
      </Box>
    );

  return (
    <Box sx={formFooterStyles["compactContent"]}>
      {statusMessage}
      {children}
    </Box>
  );
};
