import Box from "@mui/material/Box";
import { mainContainerStyles } from "./MainContainer.styles";

export const MainContainer: React.FC<React.PropsWithChildren<{}>> = function ({
  children,
}) {
  return (
    <Box component="main" sx={mainContainerStyles["container"]}>
      {children}
    </Box>
  );
};
