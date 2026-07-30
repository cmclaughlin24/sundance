import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { formTitleStyles } from "./FormTitle.style";

export const FormTitle: React.FC<{ name: string; description: string }> =
  function ({ name, description }) {
    return (
      <Box component="section" sx={formTitleStyles["container"]}>
        <Typography component="h1" sx={formTitleStyles["name"]}>
          {name}
        </Typography>
        <Typography component="p" sx={formTitleStyles["description"]}>
          {description}
        </Typography>
      </Box>
    );
  };
