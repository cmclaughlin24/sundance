import EmergencyIcon from "@mui/icons-material/Emergency";
import EditOffIcon from "@mui/icons-material/EditOff";
import VisibilityIcon from "@mui/icons-material/Visibility";
import type { RuleType } from "@/types/rule";
import type { IPaletteCategory } from "../palette";
import { PaletteItemDragType } from "@/components/FormDesigner/types/formDragEvent";

export type FormRUlesPaletteItemType = RuleType;

export const FORM_RULES_PALETTE: Readonly<
  IPaletteCategory<FormRUlesPaletteItemType>[]
> = [
  {
    label: "Behavior",
    items: [
      {
        label: "Required",
        type: "required",
        icon: <EmergencyIcon />,
        dragType: PaletteItemDragType.Rule,
      },
      {
        label: "Read Only",
        type: "readonly",
        icon: <EditOffIcon />,
        dragType: PaletteItemDragType.Rule,
      },
      {
        label: "Visbility",
        type: "visible",
        icon: <VisibilityIcon />,
        dragType: PaletteItemDragType.Rule,
      },
    ],
  },
];
