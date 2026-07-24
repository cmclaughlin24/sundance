import { type ElementType, type IElement } from "@/types/element";
import { TextField } from "../Elements/TextFieldElement";
import { useElementRuleState, useFormDispatch } from "@/store/useFormContext";
import { NumberField } from "../Elements/NumberFieldElement";
import type { IRuleState } from "@/types/rule";

export type ElementComponent = React.FC<{
  element: IElement;
  ruleState: IRuleState;
  onChange: (value: any) => void;
}>;

const registry = new Map<ElementType, ElementComponent>([
  ["text", TextField],
  ["number", NumberField],
]);

export const ElementRenderer: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const dispatch = useFormDispatch();
  const ruleState = useElementRuleState(element);

  const handleChange = (value: any) => {
    dispatch({ type: "SET_VALUE", elementId: element.id, value });
  };

  const Component = registry.get(element.type);

  if (!Component) {
    return <>Element not defined!</>;
  }

  return (
    <Component
      element={element}
      ruleState={ruleState}
      onChange={handleChange}
    />
  );
};
