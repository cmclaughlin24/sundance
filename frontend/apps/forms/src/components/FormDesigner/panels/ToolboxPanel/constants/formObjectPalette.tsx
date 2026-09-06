import { BuilderItemDragType } from "@/components/FormDesigner/types/formDragEvent";
import type { ElementType } from "@/types/element";
import ArrowDropDownCircle from "@mui/icons-material/ArrowDropDownCircle";
import CalendarToday from "@mui/icons-material/CalendarToday";
import CheckBox from "@mui/icons-material/CheckBox";
import Numbers from "@mui/icons-material/Numbers";
import Person from "@mui/icons-material/Person";
import RadioButtonChecked from "@mui/icons-material/RadioButtonChecked";
import TextFields from "@mui/icons-material/TextFields";
import ToggleOn from "@mui/icons-material/ToggleOn";
import ViewStream from "@mui/icons-material/ViewStream";
import WebAsset from "@mui/icons-material/WebAsset";
import type { IPaletteCategory } from "../palette";

export type FormObjectItemType = ElementType | "section";

export const FORM_OBJECT_PALETTE: Readonly<
  IPaletteCategory<FormObjectItemType, BuilderItemDragType>[]
> = [
  {
    label: "Basic",
    items: [
      {
        icon: <TextFields />,
        label: "Text",
        type: "text",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <Numbers />,
        label: "Number",
        type: "number",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <CalendarToday />,
        label: "Date",
        type: "date",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <ToggleOn />,
        label: "Toggle",
        type: "toggle",
        dragType: BuilderItemDragType.Element,
      },
    ],
  },
  {
    label: "Choice",
    items: [
      {
        icon: <CheckBox />,
        label: "Checkbox",
        type: "checkbox",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <RadioButtonChecked />,
        label: "Radio",
        type: "radio",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <ArrowDropDownCircle />,
        label: "Select",
        type: "select",
        dragType: BuilderItemDragType.Element,
      },
      {
        icon: <ViewStream />,
        label: "Segmented",
        type: "segmented",
        dragType: BuilderItemDragType.Element,
      },
    ],
  },
  {
    label: "Directory",
    items: [
      {
        icon: <Person />,
        label: "User",
        type: "user",
        dragType: BuilderItemDragType.Element,
      },
    ],
  },
  {
    label: "Layout",
    items: [
      {
        icon: <WebAsset />,
        label: "Section",
        type: "section",
        dragType: BuilderItemDragType.Section,
      },
    ],
  },
];
