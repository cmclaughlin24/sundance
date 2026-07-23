import type { NumberElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./BaseFieldElement";
import MuiTextField from "@mui/material/TextField";
import type { ChangeEvent } from "react";

export const NumberField: ElementComponent = function ({ element, onChange }) {
  if (element.type !== "number") {
    return <>Incorrect Element Type: {element.type}</>;
  }

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
