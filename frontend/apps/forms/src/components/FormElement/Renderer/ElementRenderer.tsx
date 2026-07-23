import { type ElementType, type IElement } from "@/types/element";
import { TextField } from "../Elements/TextFieldElement";

export type ElementComponent = React.FC<{ element: IElement }>;

const registry = new Map<ElementType, ElementComponent>([["text", TextField]]);

export const ElementRenderer: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const Component = registry.get(element.type);

  if (!Component) {
    return <>Element not defined!</>;
  }

  return <Component element={element} />;
};
