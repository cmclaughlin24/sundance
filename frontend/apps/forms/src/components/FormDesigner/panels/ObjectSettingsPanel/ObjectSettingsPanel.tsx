import { useFormDesignerDispatch } from "@/store/formDesigner";
import { ElementType, type IElement } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";
import Box from "@mui/material/Box";

export type ElementAttributesComponent = React.FC<{
  element: IElement;
  onChange: (attr: ElementAttributes) => void;
}>;

const registry = new Map<ElementType, ElementAttributesComponent>([]);

export const ObjectSettingsPanel: React.FC = function () {
  const { dispatch } = useFormDesignerDispatch();

  const handleChanges = (elementId: string, attr: ElementAttributes) => {
    console.log(elementId, attr);
  };

  return <Box>Object Settings Panel</Box>;
};
