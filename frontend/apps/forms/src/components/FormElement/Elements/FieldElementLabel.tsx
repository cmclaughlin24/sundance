import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

export interface FieldElementLabelProps {
  label: string;
  description: string | null;
  htmlFor: string;
}

export const FieldElementLabel: React.FC<{
  label: string;
  description: string | null;
  htmlFor: string;
}> = function ({ label, description, htmlFor }) {
  return (
    <Box>
      <Typography component="label" htmlFor={htmlFor}>
        {label}
      </Typography>
      {description && <Typography component="p">{description}</Typography>}
    </Box>
  );
};
