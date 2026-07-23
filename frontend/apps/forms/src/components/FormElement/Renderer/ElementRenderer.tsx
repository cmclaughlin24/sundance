import { type ElementType, type IElement } from "@/types/element";
import { TextField } from "../Elements/TextFieldElement";
import { useFormDispatch } from "@/store/useFormContext";
import { NumberField } from "../Elements/NumberFieldElement";

export type ElementComponent = React.FC<{
  element: IElement;
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

  const handleChange = (value: any) => {
    dispatch({ type: "SET_VALUE", elementId: element.id, value });
  };

  const Component = registry.get(element.type);

  if (!Component) {
    return <>Element not defined!</>;
  }

  return <Component element={element} onChange={handleChange} />;
};
