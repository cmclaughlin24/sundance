import type { TextElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./BaseFieldElement";
import MuiTextField from "@mui/material/TextField";

export const TextField: ElementComponent = function ({ element }) {
  if (element.type !== "text") {
    return <>Incorrect Element Type: {element.type}</>;
  }

  const attr = element.attributes as TextElementAttributes;

  return (
    <FieldElementContainer element={element}>
      <MuiTextField id={element.id} placeholder={attr.placeholder} />
    </FieldElementContainer>
  );
};
