import type { IElement } from "@/types/element";
import type { IPage } from "@/types/page";
import type { ISection } from "@/types/section";

export interface SelectedElement {
  type: "element";
  item: IElement;
}

export interface SelectedSection {
  type: "section";
  item: ISection;
}

export interface SelectedPage {
  type: "page";
  item: IPage;
}

export type SelectedItem = SelectedElement | SelectedSection | SelectedPage;

export function selectedToPaletteType(selected: SelectedItem) {
  switch (selected.type) {
    case "element":
      return selected.item.type;
    default:
      return selected.type;
  }
}
