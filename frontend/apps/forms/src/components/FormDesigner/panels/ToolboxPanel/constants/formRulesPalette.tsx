import EmergencyIcon from "@mui/icons-material/Emergency";
import EditOffIcon from "@mui/icons-material/EditOff";
import VisibilityIcon from "@mui/icons-material/Visibility";
import type { IPaletteCategory } from "../palette";
import { RuleItemDragType } from "@/components/FormDesigner/types/formDragEvent";
import type { RuleType } from "@/types/rule";


export const FORM_RULES_PALETTE: Readonly<
  IPaletteCategory<RuleType, RuleItemDragType>[]
> = [
  {
    label: "Behavior",
    items: [
      {
        label: "Required",
        type: "required",
        icon: <EmergencyIcon />,
        dragType: RuleItemDragType.Rule,
      },
      {
        label: "Read Only",
        type: "readonly",
        icon: <EditOffIcon />,
        dragType: RuleItemDragType.Rule,
      },
      {
        label: "Visbility",
        type: "visible",
        icon: <VisibilityIcon />,
        dragType: RuleItemDragType.Rule,
      },
    ],
  },
];
