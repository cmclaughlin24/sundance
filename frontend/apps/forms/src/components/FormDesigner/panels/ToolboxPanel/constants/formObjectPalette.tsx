import { PaletteItemDragType } from "@/components/FormDesigner/types/formDragEvent";
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

export type PaletteItemType = ElementType | "section";

export const FORM_OBJECT_PALETTE: Readonly<
  IPaletteCategory<PaletteItemType>[]
> = [
  {
    label: "Basic",
    items: [
      {
        icon: <TextFields />,
        label: "Text",
        type: "text",
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <Numbers />,
        label: "Number",
        type: "number",
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <CalendarToday />,
        label: "Date",
        type: "date",
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <ToggleOn />,
        label: "Toggle",
        type: "toggle",
        dragType: PaletteItemDragType.Element,
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
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <RadioButtonChecked />,
        label: "Radio",
        type: "radio",
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <ArrowDropDownCircle />,
        label: "Select",
        type: "select",
        dragType: PaletteItemDragType.Element,
      },
      {
        icon: <ViewStream />,
        label: "Segmented",
        type: "segmented",
        dragType: PaletteItemDragType.Element,
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
        dragType: PaletteItemDragType.Element,
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
        dragType: PaletteItemDragType.Section,
      },
    ],
  },
];
