import Typography from "@mui/material/Typography";
import { fieldElementLabelStyles } from "./FieldElementLabel.style";
import type { SxProps, Theme } from "@mui/material/styles";

export interface FieldElementLabelProps {
  label: string;
  description: string | null;
  htmlFor: string;
  required: boolean;
}

export const FieldElementLabel: React.FC<FieldElementLabelProps> = function ({
  label,
  description,
  htmlFor,
  required,
}) {
  const labelStyles: SxProps<Theme> = { ...fieldElementLabelStyles["label"] };

  if (!required) {
    delete labelStyles["::after" as keyof SxProps<Theme>];
  }

  return (
    <>
      <Typography component="label" htmlFor={htmlFor} sx={labelStyles}>
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
