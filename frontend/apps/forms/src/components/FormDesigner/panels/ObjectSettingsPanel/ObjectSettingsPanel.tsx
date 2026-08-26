import { Panel } from "@/components/layout/Panel";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import type { ElementType, IElement } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";
import Typography from "@mui/material/Typography";

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

  return (
    <Panel sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
      <Panel.Title title="Object Settings" />
      <Typography>Select an element in the canvas to view and edit it's properties.</Typography>
    </Panel>
  );
};
