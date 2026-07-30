import Typography from "@mui/material/Typography";
import { fieldElementLabelStyles } from "./FieldElementLabel.style";

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
    <>
      <Typography
        component="label"
        htmlFor={htmlFor}
        sx={fieldElementLabelStyles["label"]}
      >
        {label}
      </Typography>
      {description && (
        <Typography component="p" sx={fieldElementLabelStyles["description"]}>
          {description}
        </Typography>
      )}
    </>
  );
};
