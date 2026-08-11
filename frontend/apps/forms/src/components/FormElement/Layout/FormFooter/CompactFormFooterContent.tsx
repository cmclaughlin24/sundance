import type { FormContentComponent } from "./FormFooter";
import { formFooterStyles } from "./FormFooter.style";
import Box from "@mui/material/Box";
import { FormProgressMessage } from "../FormProgress/FormProgressMessage";

export const CompactFormFooterContent: FormContentComponent = function ({
  progress,
  children,
}) {
  const handleClick = () => {
    if (!progress?.errors) {
      return;
    }

    // TODO: Implement scroll into view on the first element with errors.
    console.log("scrolling invalid field into view");
  };

  let statusMessage =
    !progress || progress.percentage === 100 ? null : (
      <Box component="span" onClick={handleClick}>
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
