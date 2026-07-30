import Box from "@mui/material/Box";
import { formFooterActionsStyles } from "./FormFooterActions.style";

export const FormFooterActions: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    return <Box sx={formFooterActionsStyles["container"]}>{children}</Box>;
  };
