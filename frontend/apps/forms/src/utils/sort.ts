import type { HasPosition } from "@/types/hasPosition";

export function sortFormElements<T extends HasPosition>(items: T[]): T[] {
  if (!items) {
    return [];
  }

  return items.sort((a, b) => {
    return a.position - b.position;
  });
}
