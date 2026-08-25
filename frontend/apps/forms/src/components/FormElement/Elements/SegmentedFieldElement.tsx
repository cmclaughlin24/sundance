import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../renderer/ElementRenderer";
import type { SegmentedElementAttributes } from "@/types/elementAttributes";
import {
  useElementErrors,
  useElementValue,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
import type { LookupValue } from "@/types/data";
import z from "zod";
import MuiToggleButton from "@mui/material/ToggleButton";
import { FieldElementContainer } from "../layout/FieldElementContainer";
import FormControl from "@mui/material/FormControl";
import ToggleButtonGroup from "@mui/material/ToggleButtonGroup";
import FormHelperText from "@mui/material/FormHelperText";

export const SegmentedFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "segmented");

  const attr = element.attributes as SegmentedElementAttributes;
  const { setError } = useSubmissionDispatch();
  const value = useElementValue<LookupValue>(
    element.id,
    attr.defaultValue ?? "",
  );
  const errors = useElementErrors(element.id);
  const validationSchema = buildSegmentedValidationSchema({
    required: ruleState.required,
  });

  const handleChange = (
    _event: React.MouseEvent<HTMLElement>,
    newValue: LookupValue | null,
  ) => {
    if (newValue === null) {
      return;
    }

    onChange(newValue);
    const result = validationSchema.safeParse(newValue);
    setError(
      element.id,
      result.success ? [] : result.error.issues.map((e) => e.message),
    );
  };

  const handleBlur = () => {
    const result = validationSchema.safeParse(value);
    setError(
      element.id,
      result.success ? [] : result.error.issues.map((e) => e.message),
    );
  };

  const content = attr.data.map((lookup) => (
    <MuiToggleButton
      value={lookup.value}
      key={`${lookup.value}=${lookup.label}`}
      data-testid={`segmented-option-${lookup.value}`}
    >
      {lookup.label}
    </MuiToggleButton>
  ));

  return (
    <FieldElementContainer element={element}>
      <FormControl error={errors?.length > 0}>
        <ToggleButtonGroup
          value={value}
          exclusive
          onChange={handleChange}
          onBlur={handleBlur}
          disabled={ruleState.readonly}
          data-testid={`segmented-field-${element.id}`}
        >
          {content}
        </ToggleButtonGroup>
        {errors && errors[0] && (
          <FormHelperText data-testid="segmented-field-error">
            {errors[0]}
          </FormHelperText>
        )}
      </FormControl>
    </FieldElementContainer>
  );
};

function buildSegmentedValidationSchema(options: {
  required: boolean;
}): z.ZodType<any> {
  const item = z.union([z.string(), z.number()]);

  if (!options.required) {
    return item.optional();
  }

  return item.refine((v) => v !== "" && v != null, "This field is required");
}
