import * as ArrayUtils from "@/utils/array";
import {
  BuilderItemDragType,
  RuleItemDragType,
} from "../../types/formDragEvent";
import {
  FORM_OBJECT_PALETTE,
  type FormObjectItemType,
} from "./constants/formObjectPalette";
import { FORM_RULES_PALETTE } from "./constants/formRulesPalette";
import type { RuleType } from "@/types/rule";

export interface IPaletteCategory<IType, DType> {
  label: string;
  items: IPaletteItem<IType, DType>[];
}

export interface IPaletteItem<IType, DType> {
  icon: React.ReactNode;
  label: string;
  type: IType;
  dragType: DType;
}

/**
 * Filters the pallette based on the search term.
 * @param searchTerm The term to filter the pallette items by.
 * @returns The filtered pallette categories containing items that match the search term.
 */
export function filterPalette<IType, DType>(
  searchTerm: string,
  palette: IPaletteCategory<IType, DType>[],
): Readonly<IPaletteCategory<IType, DType>[]> {
  if (!searchTerm) {
    return palette;
  }

  const filtered: IPaletteCategory<IType, DType>[] = [];

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
  type: FormObjectItemType,
): IPaletteItem<FormObjectItemType, BuilderItemDragType> | null {
  return findPaletteItem(type, FORM_OBJECT_PALETTE);
}

export function findFormRulePaletteItem(
  type: RuleType,
): IPaletteItem<RuleType, RuleItemDragType> | null {
  return findPaletteItem(type, FORM_RULES_PALETTE);
}

export function findPaletteItem<IType, DType>(
  type: IType,
  palette: Readonly<IPaletteCategory<IType, DType>[]>,
): IPaletteItem<IType, DType> | null {
  for (const category of palette) {
    const item = category.items.find((i) => i.type === type);

    if (item) {
      return item;
    }
  }

  return null;
}
