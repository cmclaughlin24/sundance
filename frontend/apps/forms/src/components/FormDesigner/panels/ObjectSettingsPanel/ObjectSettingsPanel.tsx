import { Panel } from "@/components/layout/Panel";
import {
  selectedToPaletteType,
  useFormDesignerDispatch,
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
  type IdentitySettingsEvent,
  type IdentitySettingsProps,
} from "./settings/IdentitySettings";
import { BehaviorSettings } from "./settings/BehaviorSettings";
import type {
  FormDesignerEvent,
  UpdateElementEvent,
  UpdateSectionEvent,
} from "@/store/formDesigner/events";

export type ElementSettingsComponent = React.FC<{
  element: IElement;
  onChange: (attr: ElementAttributes) => void;
}>;

const registry = new Map<ElementType, ElementSettingsComponent>([
  ["text", TextElementSettings],
  ["number", NumberElementSettings],
]);

export const ObjectSettingsPanel: React.FC = function () {
  const { dispatch } = useFormDesignerDispatch();
  const { selected } = useFormDesignerSelect();

  const handleIdentityChanges = (e: IdentitySettingsEvent) => {
    let event: FormDesignerEvent;

    switch (e.type) {
      case "section":
        event = {
          type: "UpdateSection",
          sectionId: selected!.item.id,
          changes: e.changes,
        } satisfies UpdateSectionEvent;
        break;
      case "element":
        event = {
          type: "UpdateElement",
          elementId: selected!.item.id,
          changes: e.changes,
        } satisfies UpdateElementEvent;
        break;
    }

    dispatch(event!);
  };

  let content = (
    <Typography>
      Select an element in the canvas to view and edit it's properties.
    </Typography>
  );

  if (selected) {
    const elementType = selectedToPaletteType(selected);
    let PropertyComponent =
      selected.type === "element"
        ? registry.get(selected.item.type)
        : undefined;

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
            <IdentitySettings
              onChange={handleIdentityChanges}
              {...selectedToIdentityProps(selected)}
            />
          </Collapisble>
          {selected.type === "element" && (
            <>
              {PropertyComponent && (
                <Collapisble summary="Properties">
                  <PropertyComponent
                    element={selected.item}
                    onChange={() => {}}
                  />
                </Collapisble>
              )}
              <Collapisble summary="Behavior">
                <BehaviorSettings />
              </Collapisble>
            </>
          )}
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
