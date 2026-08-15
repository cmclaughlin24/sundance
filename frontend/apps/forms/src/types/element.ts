import type { ElementAttributes } from "./elementAttributes";
import type { HasKey } from "./hasKey";
import type { HasName } from "./hasName";
import type { HasPosition } from "./hasPosition";
import type { HasRules } from "./rule";

export type ElementType =
  | "checkbox"
  | "date"
  | "number"
  | "radio"
  | "select"
  | "segmented"
  | "text"
  | "toggle"
  | "user";

export interface IElement extends HasKey, HasPosition, HasName, HasRules {
  id: string;
  type: ElementType;
  attributes: ElementAttributes;
}
