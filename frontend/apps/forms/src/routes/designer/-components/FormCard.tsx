import type { IForm } from "@/types/form";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import { formCardStyles } from "./FormCard.style";

export const FormCard: React.FC<{ form: IForm }> = function ({ form }) {
  return (
    <Card variant="outlined" sx={formCardStyles["card"]}>
      <Typography sx={formCardStyles["name"]}>{form.name}</Typography>
      <Typography sx={formCardStyles["description"]}>
        {form.description}
      </Typography>
      <Box></Box>
    </Card>
  );
};
