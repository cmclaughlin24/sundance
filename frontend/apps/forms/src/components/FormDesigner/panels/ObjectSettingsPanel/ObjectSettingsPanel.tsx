import { Panel } from "@/components/layout/Panel";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import type { ElementType, IElement } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";

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
    <Panel>
      <Panel.Title title="Object Settings Panel" />
    </Panel>
  );
};
