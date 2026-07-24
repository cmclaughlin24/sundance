import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import type { CheckboxElementAttributes } from "@/types/elementAttributes";
import type { ILookup } from "@/types/data";
import MuiCheckbox from "@mui/material/Checkbox";
import MuiFormControlLabel from "@mui/material/FormControlLabel";
import Box from "@mui/material/Box";
import type { ChangeEvent } from "react";
import { FieldElementContainer } from "./FieldElementContainer";

export const CheckboxFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "checkbox");

  const attr = element.attributes as CheckboxElementAttributes;
  const data: ILookup[] = attr.data;

  const handleChange =
    (lookup: ILookup) => (event: ChangeEvent<HTMLInputElement>) => {
      onChange({ value: lookup.value, checked: event.target.checked });
    };

  const content = data.map((lookup) => (
    <MuiFormControlLabel
      label={lookup.label}
      required={ruleState.required}
      disabled={ruleState.readonly}
      key={`${lookup.value}=${lookup.label}`}
      control={
        <MuiCheckbox
          defaultChecked={attr.isCheckedByDefault}
          onChange={handleChange(lookup)}
        />
      }
    />
  ));

  return (
    <FieldElementContainer element={element}>
      <Box>{content}</Box>
    </FieldElementContainer>
  );
};
