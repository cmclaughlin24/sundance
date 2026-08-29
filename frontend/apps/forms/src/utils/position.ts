import * as ArrayUtils from "@/utils/array";
import type { HasPosition } from "@/types/hasPosition";

export function sortPositioned<T extends HasPosition>(items: T[]): T[] {
  if (!items) {
    return [];
  }

  return [...items].sort((a, b) => {
    return a.position - b.position;
  });
}

export function getNextPosition<T extends HasPosition>(items: T[]): number {
  if (ArrayUtils.hasLengthGreaterThan(items, 0)) {
    return items[items.length - 1].position + 1;
  }

  return 0;
}

export function getBetweenPosition<T extends HasPosition>(
  index: number,
  items: T[],
): number {
  if (index === 0) {
    return ArrayUtils.hasLengthGreaterThan(items, 0) ? items[0].position - 1 : 0;
  }

  if (index === items.length) {
    return items[items.length - 1].position + 1;
  }

  return (items[index - 1].position + items[index].position) / 2;
}

export function swapPositions<T extends HasPosition & { id: string }>(
  items: T[],
  id: string,
  inc: -1 | 1,
): T[] {
  const sorted = sortPositioned(items);
  const sortedIndex = sorted.findIndex((item) => item.id === id);

  if (sortedIndex === -1) {
    return items;
  }

  const siblingIndex = sortedIndex + inc;

  if (siblingIndex < 0 || siblingIndex >= sorted.length) {
    return items;
  }

  const current = sorted[sortedIndex];
  const sibling = sorted[siblingIndex];
  const currentPosition = current.position;
  const siblingPosition = sibling.position;

  return items.map((item) => {
    if (item.id === current.id) {
      return { ...item, position: siblingPosition };
    }

    if (item.id === sibling.id) {
      return { ...item, position: currentPosition };
    }

    return item;
  });
}
