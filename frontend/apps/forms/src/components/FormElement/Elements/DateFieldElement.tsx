import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./FieldElementContainer";
import type { DateElementAttributes } from "@/types/elementAttributes";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";

export const DateFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "date");

  const attr = element.attributes as DateElementAttributes;

  const handleChange = (value: Dayjs | null) => {
    onChange(value ? value.toISOString() : null);
  };

  return (
    <FieldElementContainer element={element}>
      <DatePicker
        disabled={ruleState.readonly}
        minDate={attr.minDate ? dayjs(attr.minDate) : undefined}
        maxDate={attr.maxDate ? dayjs(attr.maxDate) : undefined}
        onChange={handleChange}
        slotProps={{
          textField: { required: ruleState.required, id: element.id },
        }}
      />
    </FieldElementContainer>
  );
};
