import MuiSelectField, { type SelectChangeEvent } from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import type { ILookup } from "@/types/data";
import type { SelectElementAttributes } from "@/types/elementAttributes";
import { checkElementType } from "@/utils/error";

export const SelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "select");

  const attr = element.attributes as SelectElementAttributes;
  const data: ILookup[] = attr.data;

  const handleChange = (event: SelectChangeEvent) => {
    onChange(event.target.value);
  };

  let content = data.map((lookup) => (
    <MenuItem value={lookup.value} key={`${lookup.value}=${lookup.label}`}>
      {lookup.label}
    </MenuItem>
  ));

  return (
    <MuiSelectField
      id={element.id}
      required={ruleState.required}
      disabled={ruleState.readonly}
      onChange={handleChange}
    >
      {content}
    </MuiSelectField>
  );
};
