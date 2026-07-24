import type { ElementType } from "@/types/element";

export function checkElementType(
  elementType: ElementType,
  expected: ElementType,
) {
  if (elementType !== expected) {
    throw new Error(
      `invalid element type ${elementType}; expected ${expected}`,
    );
  }
}
