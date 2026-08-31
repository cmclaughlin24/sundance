import { Panel } from "@/components/layout/Panel";
import {
  selectedToPaletteType,
  useFormDesignerSelect,
  type SelectedItem,
} from "@/store/formDesigner";
import type { ElementType, IElement } from "@/types/element";
import type { ElementAttributes } from "@/types/elementAttributes";
import Typography from "@mui/material/Typography";
import { ActiveObjectTitle } from "./ActiveObjectTitle";
import { Collapisble } from "@/components/Collapisble";
import Box from "@mui/material/Box";
import { Border } from "@/constants/colors";
import { TextElementSettings } from "./settings/TextElementSettings";
import { NumberElementSettings } from "./settings/NumberElementSettings";
import {
  IdentitySettings,
  type IdentitySettingsProps,
} from "./settings/IdentitySettings";
import { BehaviorSettings } from "./settings/BehaviorSettings";

export type ElementSettingsComponent = React.FC<{
  element: IElement;
  onChange: (attr: ElementAttributes) => void;
}>;

const registry = new Map<ElementType, ElementSettingsComponent>([
  ["text", TextElementSettings],
  ["number", NumberElementSettings],
]);

export const ObjectSettingsPanel: React.FC = function () {
  const { selected } = useFormDesignerSelect();

  const handleChanges = (elementId: string, attr: ElementAttributes) => {
    console.log(elementId, attr);
  };

  let content = (
    <Typography>
      Select an element in the canvas to view and edit it's properties.
    </Typography>
  );

  if (selected) {
    const elementType = selectedToPaletteType(selected);

    content = (
      <>
        <ActiveObjectTitle elementType={elementType} />
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            gap: 2.5,
            "& > *": {
              borderBottom: `1px solid ${Border.Primary}`,
            },
          }}
        >
          <Collapisble summary="Identity">
            <IdentitySettings {...selectedToIdentityProps(selected)} />
          </Collapisble>
          <Collapisble summary="Properties">Identity Content</Collapisble>
          <Collapisble summary="Data Sources">Identity Content</Collapisble>
          <Collapisble summary="Behavior">
            <BehaviorSettings />
          </Collapisble>
        </Box>
      </>
    );
  }

  return (
    <Panel sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
      <Panel.Title title="Object Settings" />
      {content}
    </Panel>
  );
};

function selectedToIdentityProps(
  selected: SelectedItem,
): IdentitySettingsProps {
  switch (selected.type) {
    case "element":
      return { type: "element", object: selected.item };
    case "page":
      return { type: "page", object: selected.item };
    case "section":
      return { type: "section", object: selected.item };
  }
}
