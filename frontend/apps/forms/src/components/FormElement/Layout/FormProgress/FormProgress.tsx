import type { IFormProgress } from "@/utils/progress";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { formProgressStyles } from "./FormProgress.style";
import { FormProgressMessage } from "./FormProgressMessage";

export const FormProgress: React.FC<{ progress: IFormProgress }> = function ({
  progress,
}) {
  let statusMessage = !progress ? null : (
    <FormProgressMessage progress={progress} />
  );

  return (
    <Box sx={formProgressStyles["container"]}>
      <Typography variant="body2" sx={formProgressStyles["required"]}>
        Required Information
      </Typography>
      <Typography variant="body1" sx={formProgressStyles["text"]}>
        {progress.filled}/{progress.total} completed ·{" "}
        {statusMessage}
      </Typography>
    </Box>
  );
};
