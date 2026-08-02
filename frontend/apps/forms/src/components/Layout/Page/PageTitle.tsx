import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { pageTitleStyles } from "./PageTitle.style";

export const PageTitle: React.FC<{ name: string; description: string }> =
  function ({ name, description }) {
    return (
      <Box component="section" sx={pageTitleStyles["container"]}>
        <Typography component="h1" sx={pageTitleStyles["name"]}>
          {name}
        </Typography>
        <Typography component="p" sx={pageTitleStyles["description"]}>
          {description}
        </Typography>
      </Box>
    );
  };
