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

export interface IPalletteCategory {
  label: string;
  items: IPalletteItem[];
}

export interface IPalletteItem {
  icon: React.ReactNode;
  label: string;
  type: ElementType | "section";
}

/**
 * Filters the pallette based on the search term.
 * @param searchTerm The term to filter the pallette items by.
 * @returns The filtered pallette categories containing items that match the search term.
 */
export function filterPallette(
  searchTerm: string,
): Readonly<IPalletteCategory[]> {
  if (!searchTerm) {
    return PALLETTE;
  }

  const pallette: IPalletteCategory[] = [];

  for (const category of PALLETTE) {
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

const PALLETTE: Readonly<IPalletteCategory[]> = [
  {
    label: "Basic",
    items: [
      { icon: <TextFields />, label: "Text", type: "text" },
      { icon: <Numbers />, label: "Number", type: "number" },
      { icon: <CalendarToday />, label: "Date", type: "date" },
      { icon: <ToggleOn />, label: "Toggle", type: "toggle" },
    ],
  },
  {
    label: "Choice",
    items: [
      { icon: <CheckBox />, label: "Checkbox", type: "checkbox" },
      { icon: <RadioButtonChecked />, label: "Radio", type: "radio" },
      { icon: <ArrowDropDownCircle />, label: "Select", type: "select" },
      { icon: <ViewStream />, label: "Segmented", type: "segmented" },
    ],
  },
  {
    label: "Directory",
    items: [{ icon: <Person />, label: "User", type: "user" }],
  },
  {
    label: "Layout",
    items: [{ icon: <WebAsset />, label: "Section", type: "section" }],
  },
];
