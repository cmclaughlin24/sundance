import type { NumberElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./FieldElementContainer";
import MuiTextField from "@mui/material/TextField";
import type { ChangeEvent } from "react";
import { checkElementType } from "@/utils/error";

export const NumberFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "number");

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    onChange(event.target.value);
  };

  const attr = element.attributes as NumberElementAttributes;

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        type="number"
        required={ruleState.required}
        disabled={ruleState.readonly}
        onChange={handleChange}
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
