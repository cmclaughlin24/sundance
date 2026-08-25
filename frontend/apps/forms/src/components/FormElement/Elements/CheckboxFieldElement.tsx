import { checkElementType } from "@/utils/error";
import type { ElementComponent } from "../renderer/ElementRenderer";
import type { CheckboxElementAttributes } from "@/types/elementAttributes";
import type { ILookup } from "@/types/data";
import MuiCheckbox from "@mui/material/Checkbox";
import MuiFormControlLabel from "@mui/material/FormControlLabel";
import Box from "@mui/material/Box";
import type { ChangeEvent } from "react";
import { FieldElementContainer } from "../layout/FieldElementContainer";
import { useDataSource } from "@/hooks/useDataSource";

interface BaseCheckboxFieldElementProps {
  data: ILookup[];
  required?: boolean;
  readonly?: boolean;
  isCheckedByDefault?: boolean;
  onChange: (event: any) => void;
}

export const BaseCheckboxFieldElement: React.FC<BaseCheckboxFieldElementProps> =
  function ({ data, required, readonly, isCheckedByDefault, onChange }) {
    const handleChange =
      (lookup: ILookup) => (event: ChangeEvent<HTMLInputElement>) => {
        onChange({ value: lookup.value, checked: event.target.checked });
      };

    const content = data.map((lookup) => (
      <MuiFormControlLabel
        label={lookup.label}
        required={required}
        disabled={readonly}
        key={`${lookup.value}=${lookup.label}`}
        control={
          <MuiCheckbox
            defaultChecked={isCheckedByDefault}
            onChange={handleChange(lookup)}
          />
        }
      />
    ));

    return <Box>{content}</Box>;
  };

const StaticCheckboxElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as CheckboxElementAttributes;

  return (
    <BaseCheckboxFieldElement
      data={attr.data}
      isCheckedByDefault={attr.isCheckedByDefault}
      required={ruleState.required}
      readonly={ruleState.readonly}
      onChange={onChange}
    />
  );
};

const DynamicCheckboxElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as CheckboxElementAttributes;
  const { data, isLoading } = useDataSource(attr.dataSourceRef!);

  return (
    <BaseCheckboxFieldElement
      data={data || []}
      isCheckedByDefault={attr.isCheckedByDefault}
      required={ruleState.required}
      readonly={ruleState.readonly || isLoading}
      onChange={onChange}
    />
  );
};

export const CheckboxFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "checkbox");

  const attr = element.attributes as CheckboxElementAttributes;
  let Component = StaticCheckboxElement;

  if (attr.dataSourceRef) {
    Component = DynamicCheckboxElement;
  }

  return (
    <FieldElementContainer element={element}>
      <Component element={element} ruleState={ruleState} onChange={onChange} />
    </FieldElementContainer>
  );
};
