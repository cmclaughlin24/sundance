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
import * as ArrayUtils from "@/utils/array";
import { PaletteItemDragType } from "../../types/formDragEvent";

export type PaletteItemType = ElementType | "section";

export interface IPaletteCategory {
  label: string;
  items: IPaletteItem[];
}

export interface IPaletteItem {
  icon: React.ReactNode;
  label: string;
  type: PaletteItemType;
  dragType: PaletteItemDragType;
}

/**
 * Filters the pallette based on the search term.
 * @param searchTerm The term to filter the pallette items by.
 * @returns The filtered pallette categories containing items that match the search term.
 */
export function filterPalette(
  searchTerm: string,
): Readonly<IPaletteCategory[]> {
  if (!searchTerm) {
    return PALETTE;
  }

  const pallette: IPaletteCategory[] = [];

  for (const category of PALETTE) {
    const items = category.items.filter((item) =>
      item.label.toLowerCase().includes(searchTerm.toLowerCase()),
    );

    if (!ArrayUtils.hasLengthGreaterThan(items, 0)) {
      continue;
    }

    pallette.push({ ...category, items });
  }

  return pallette;
}

/**
 * Finds a `IPalletteItem` by its type.
 * @param type The type of pallette item to find.
 * @returns Teh pallette item if found, otherwise null.
 */
export function findPaletteItem(
  type: ElementType | "section" | "page",
): IPaletteItem | null {
  for (const category of PALETTE) {
    const item = category.items.find((i) => i.type === type);

    if (item) {
      return item;
    }
  }

  return null;
}

const PALETTE: Readonly<IPaletteCategory[]> = [
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
