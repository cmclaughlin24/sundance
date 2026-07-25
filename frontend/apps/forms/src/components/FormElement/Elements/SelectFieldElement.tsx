import MuiSelectField, { type SelectChangeEvent } from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import type { ILookup } from "@/types/data";
import type { SelectElementAttributes } from "@/types/elementAttributes";
import { checkElementType } from "@/utils/error";
import { FieldElementContainer } from "./FieldElementContainer";
import { useDataSource } from "@/hooks/useDataSource";

interface BaseSelectFieldElementProps {
  id: string;
  data: ILookup[];
  required?: boolean;
  readonly?: boolean;
  multiple?: boolean;
  onChange: (event: any) => void;
}

const BaseSelectFieldElement: React.FC<BaseSelectFieldElementProps> =
  function ({ data, required, readonly, multiple, id, onChange }) {
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
        id={id}
        required={required}
        disabled={readonly}
        multiple={multiple}
        onChange={handleChange}
      >
        {content}
      </MuiSelectField>
    );
  };

const StaticSelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as SelectElementAttributes;
  return (
    <BaseSelectFieldElement
      id={element.id}
      multiple={attr.multiple}
      readonly={ruleState.readonly}
      required={ruleState.required}
      data={attr.data}
      onChange={onChange}
    />
  );
};

const DynamicSelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as SelectElementAttributes;
  const { data, isLoading } = useDataSource(attr.dataSourceRef!);

  return (
    <BaseSelectFieldElement
      id={element.id}
      multiple={attr.multiple}
      readonly={ruleState.readonly || isLoading}
      required={ruleState.required}
      data={data || []}
      onChange={onChange}
    />
  );
};

export const SelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "select");

  const attr = element.attributes as SelectElementAttributes;
  let Component = StaticSelectFieldElement;

  if (attr.dataSourceRef) {
    Component = DynamicSelectFieldElement;
  }

  return (
    <FieldElementContainer element={element}>
      <Component element={element} ruleState={ruleState} onChange={onChange} />
    </FieldElementContainer>
  );
};
