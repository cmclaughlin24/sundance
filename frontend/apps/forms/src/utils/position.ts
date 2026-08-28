import * as ArrayUtils from "@/utils/array";
import type { HasPosition } from "@/types/hasPosition";

export function sortPositioned<T extends HasPosition>(items: T[]): T[] {
  if (!items) {
    return [];
  }

  return items.sort((a, b) => {
    return a.position - b.position;
  });
}

export function getNextPosition<T extends HasPosition>(items: T[]): number {
  if (ArrayUtils.hasLengthGreaterThan(items, 1)) {
    return items[items.length - 1].position + 1;
  }

  return 0;
}
