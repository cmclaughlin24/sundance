import type { TextElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./FieldElementContainer";
import MuiTextField from "@mui/material/TextField";
import { type ChangeEvent } from "react";
import { checkElementType } from "@/utils/error";
import {
  useElementErrors,
  useElementValue,
  useFormDispatch,
} from "@/store/useFormStoreContext";
import { z, ZodString } from "zod";

export const TextFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "text");

  const { setError } = useFormDispatch();
  const value = useElementValue<string>(element.id, "");
  const errors = useElementErrors(element.id);
  const attr = element.attributes as TextElementAttributes;
  const validationSchema = buildTextValidationSchema({
    required: ruleState.required,
    minLength: attr.minLength,
    maxLength: attr.maxLength,
    pattern: attr.pattern,
  });

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    onChange(event.target.value);
  };

  const handleBlur = () => {
    const result = validationSchema.safeParse(value);
    let errors: string[] = [];

    if (!result.success) {
      errors = result.error.issues.map((e) => e.message);
    }

    setError(element.id, errors);
  };

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        value={value}
        required={ruleState.required}
        disabled={ruleState.readonly}
        placeholder={attr.placeholder}
        error={errors && errors.length > 0}
        helperText={errors?.[0]}
        onBlur={handleBlur}
        onChange={handleChange}
      />
    </FieldElementContainer>
  );
};

function buildTextValidationSchema(options: {
  required: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: string;
}): ZodString {
  let schema = z.string();

  if (options.required) {
    schema = schema.min(1, "This field is required");
  }

  if (options.minLength) {
    schema = schema.min(
      options.minLength,
      `Minimum ${options.minLength} characters`,
    );
  }

  if (options.maxLength) {
    schema = schema.max(
      options.maxLength,
      `Maximum ${options.minLength} characters`,
    );
  }

  if (options.pattern != null) {
    schema = schema.regex(
      new RegExp(options.pattern),
      `Value does not match required format`,
    );
  }

  return schema;
}
