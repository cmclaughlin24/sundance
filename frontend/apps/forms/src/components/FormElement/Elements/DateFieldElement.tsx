import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "../Layout/FieldElementContainer";
import type { DateElementAttributes } from "@/types/elementAttributes";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import {
  useElementErrors,
  useElementValue,
  useFormDispatch,
} from "@/store/useFormStoreContext";
import { z } from "zod";

export const DateFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "date");

  const { setError } = useFormDispatch();
  const value = useElementValue<string | null>(element.id, null);
  const errors = useElementErrors(element.id);
  const attr = element.attributes as DateElementAttributes;
  const validationSchema = buildDateValidationSchema({
    required: ruleState.required,
    minDate: attr.minDate,
    maxDate: attr.maxDate,
  });

  const handleChange = (value: Dayjs | null) => {
    onChange(value ? value.toISOString() : null);
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
      <DatePicker
        disabled={ruleState.readonly}
        value={value ? dayjs(value) : null}
        minDate={attr.minDate ? dayjs(attr.minDate) : undefined}
        maxDate={attr.maxDate ? dayjs(attr.maxDate) : undefined}
        onChange={handleChange}
        slotProps={{
          textField: {
            required: ruleState.required,
            id: element.id,
            error: errors.length > 0,
            helperText: errors[0],
            onBlur: handleBlur,
          },
        }}
      />
    </FieldElementContainer>
  );
};

function buildDateValidationSchema(options: {
  required: boolean;
  minDate?: string;
  maxDate?: string;
}): z.ZodTypeAny {
  const base = z.iso.datetime({ message: "Must be a valid date" });

  const refined = base
    .refine(
      (val) =>
        options.minDate == null || !dayjs(val).isBefore(dayjs(options.minDate)),
      options.minDate != null
        ? `Must be on or after ${options.minDate}`
        : undefined,
    )
    .refine(
      (val) =>
        options.maxDate == null || !dayjs(val).isAfter(dayjs(options.maxDate)),
      options.maxDate != null
        ? `Must be on or before ${options.maxDate}`
        : undefined,
    );

  if (!options.required) {
    return refined.optional().nullable();
  }

  return refined;
}
