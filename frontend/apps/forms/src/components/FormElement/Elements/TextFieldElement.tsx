import type { TextElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./BaseFieldElement";
import MuiTextField from "@mui/material/TextField";
import type { ChangeEvent } from "react";

export const TextField: ElementComponent = function ({ element, onChange }) {
  if (element.type !== "text") {
    return <>Incorrect Element Type: {element.type}</>;
  }

  const handleChange = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    onChange(event.target.value);
  };

  const attr = element.attributes as TextElementAttributes;

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        placeholder={attr.placeholder}
        onChange={handleChange}
      />
    </FieldElementContainer>
  );
};
