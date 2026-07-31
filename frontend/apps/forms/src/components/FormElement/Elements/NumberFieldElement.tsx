import type { NumberElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "../Layout/FieldElementContainer";
import MuiTextField from "@mui/material/TextField";
import type { ChangeEvent } from "react";
import { checkElementType } from "@/utils/error";
import {
  useElementErrors,
  useElementValue,
  useFormDispatch,
} from "@/store/useFormStoreContext";
import { z } from "zod";

export const NumberFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "number");

  const { setError } = useFormDispatch();
  const value = useElementValue<number | "">(element.id, "");
  const errors = useElementErrors(element.id);
  const attr = element.attributes as NumberElementAttributes;
  const validationSchema = buildNumberValidationSchema({
    required: ruleState.required,
    min: attr.min,
    max: attr.max,
  });

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const raw = event.target.value;
    onChange(raw === "" ? "" : Number(raw));
  };

  const handleBlur = () => {
    const result = validationSchema.safeParse(value === "" ? undefined : value);
    setError(
      element.id,
      result.success ? [] : result.error.issues.map((e) => e.message),
    );
  };

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        type="number"
        value={value}
        required={ruleState.required}
        disabled={ruleState.readonly}
        error={errors.length > 0}
        helperText={errors[0]}
        onChange={handleChange}
        onBlur={handleBlur}
        slotProps={{
          htmlInput: {
            min: attr.min,
            max: attr.max,
            step: attr.step,
          },
        }}
      />
    </FieldElementContainer>
  );
};

function buildNumberValidationSchema(options: {
  required: boolean;
  min?: number;
  max?: number;
}): z.ZodOptional<z.ZodNumber> | z.ZodNumber {
  let schema = z.number({ error: "Must be a valid number" });

  if (options.min != null) {
    schema = schema.min(options.min, `Must be at least ${options.min}`);
  }

  if (options.max != null) {
    schema = schema.max(options.max, `Must be no more than ${options.max}`);
  }

  if (!options.required) {
    return schema.optional();
  }

  return schema;
}
