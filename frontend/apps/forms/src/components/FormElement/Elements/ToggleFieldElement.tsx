import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../renderer/ElementRenderer";
import {
  useElementErrors,
  useElementValue,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
import z from "zod";
import { FieldElementContainer } from "../layout/FieldElementContainer";
import FormControl from "@mui/material/FormControl";
import Switch from "@mui/material/Switch";
import type { ToggleElementAttributes } from "@/types/elementAttributes";

export const ToggleFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "toggle");

  const attr = element.attributes as ToggleElementAttributes;
  const { setError } = useSubmissionDispatch();
  const value = useElementValue<boolean>(element.id, attr.defaultValue);
  const errors = useElementErrors(element.id);
  const validationSchema = buildToggleValidationSchema({
    required: ruleState.required,
  });

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange(event.target.value);
  };

  const handleBlur = () => {
    const result = validationSchema.safeParse(value);
    setError(
      element.id,
      result.success ? [] : result.error.issues.map((e) => e.message),
    );
  };

  return (
    <FieldElementContainer element={element}>
      <FormControl error={errors?.length > 0}>
        <Switch value={value} onChange={handleChange} onBlur={handleBlur} />
      </FormControl>
    </FieldElementContainer>
  );
};

function buildToggleValidationSchema(options: {
  required: boolean;
}): z.ZodType<any> {
  const schema = z.boolean();

  if (options.required) {
    return schema.refine((val) => val === true, {
      error: "This field is required",
    });
  }

  return schema;
}
