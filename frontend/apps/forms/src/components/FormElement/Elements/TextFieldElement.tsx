import type { TextElementAttributes } from "@/types/elementAttributes";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import { FieldElementContainer } from "./FieldElementContainer";
import MuiTextField from "@mui/material/TextField";
import type { ChangeEvent } from "react";
import { checkElementType } from "@/utils/error";
import { useElementValue } from "@/store/useFormStoreContext";

export const TextFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "text");

  const value = useElementValue<string>(element.id, "");

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    onChange(event.target.value);
  };

  const attr = element.attributes as TextElementAttributes;

  return (
    <FieldElementContainer element={element}>
      <MuiTextField
        id={element.id}
        value={value}
        required={ruleState.required}
        disabled={ruleState.readonly}
        placeholder={attr.placeholder}
        onChange={handleChange}
      />
    </FieldElementContainer>
  );
};
