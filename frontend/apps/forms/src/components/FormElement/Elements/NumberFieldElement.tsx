import type { NumberElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./BaseFieldElement";
import MuiTextField from "@mui/material/TextField";

export const NumberField: ElementComponent = function ({ element }) {
  if (element.type !== "number") {
    return <>Incorrect Element Type: {element.type}</>;
  }

  const attr = element.attributes as NumberElementAttributes;

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        type="number"
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
