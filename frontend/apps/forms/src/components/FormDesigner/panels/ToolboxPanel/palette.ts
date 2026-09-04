import * as ArrayUtils from "@/utils/array";
import { PaletteItemDragType } from "../../types/formDragEvent";
import {
  FORM_OBJECT_PALETTE,
  type PaletteItemType,
} from "./constants/formObjectPalette";

export interface IPaletteCategory<T> {
  label: string;
  items: IPaletteItem<T>[];
}

export interface IPaletteItem<T> {
  icon: React.ReactNode;
  label: string;
  type: T;
  dragType: PaletteItemDragType;
}

/**
 * Filters the pallette based on the search term.
 * @param searchTerm The term to filter the pallette items by.
 * @returns The filtered pallette categories containing items that match the search term.
 */
export function filterPalette<T>(
  searchTerm: string,
  palette: IPaletteCategory<T>[],
): Readonly<IPaletteCategory<T>[]> {
  if (!searchTerm) {
    return palette;
  }

  const filtered: IPaletteCategory<T>[] = [];

  for (const category of palette) {
    const items = category.items.filter((item) =>
      item.label.toLowerCase().includes(searchTerm.toLowerCase()),
    );

    if (!ArrayUtils.hasLengthGreaterThan(items, 0)) {
      continue;
    }

    filtered.push({ ...category, items });
  }

  return filtered;
}

/**
 * Finds a `IPalletteItem` by its type.
 * @param type The type of pallette item to find.
 * @returns Teh pallette item if found, otherwise null.
 */
export function findFormObjectPaletteItem(
  type: PaletteItemType,
): IPaletteItem<PaletteItemType> | null {
  return findPaletteItem(type, FORM_OBJECT_PALETTE);
}

export function findPaletteItem<T>(
  type: T,
  palette: Readonly<IPaletteCategory<T>[]>,
): IPaletteItem<T> | null {
  for (const category of palette) {
    const item = category.items.find((i) => i.type === type);

    if (item) {
      return item;
    }
  }

  return null;
}
