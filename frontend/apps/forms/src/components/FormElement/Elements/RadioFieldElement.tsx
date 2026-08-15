import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import type { RadioElementAttributes } from "@/types/elementAttributes";
import {
  useElementErrors,
  useElementValue,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
import type { LookupValue } from "@/types/data";
import z from "zod";
import { FieldElementContainer } from "../Layout/FieldElementContainer";
import FormControl from "@mui/material/FormControl";
import MuiRadioGroup from "@mui/material/RadioGroup";
import MuiRadio from "@mui/material/Radio";
import FormControlLabel from "@mui/material/FormControlLabel";

export const RadioFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "radio");

  const attr = element.attributes as RadioElementAttributes;
  const { setError } = useSubmissionDispatch();
  const value = useElementValue<LookupValue>(
    element.id,
    attr.defaultValue ?? "",
  );
  const errors = useElementErrors(element.id);
  const validationSchema = buildRadioValidationSchema({
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
        <MuiRadioGroup
          id={element.id}
          value={value}
          row={attr.orientation === "horizontal"}
          onChange={handleChange}
          onBlur={handleBlur}
          data-testid={`radio-field-${element.id}`}
        >
          {attr.data.map((lookup) => (
            <FormControlLabel
              key={`${lookup.value}=${lookup.value}`}
              value={lookup.value}
              label={lookup.label}
              disabled={ruleState.readonly}
              control={
                <MuiRadio data-testid={`radio-option-${lookup.value}`} />
              }
            />
          ))}
        </MuiRadioGroup>
      </FormControl>
    </FieldElementContainer>
  );
};

function buildRadioValidationSchema(options: {
  required: boolean;
}): z.ZodType<any> {
  const item = z.union([z.string(), z.number()]);

  if (!options.required) {
    return item.optional();
  }

  return item.refine((v) => v !== "" && v != null, "This field is required");
}
