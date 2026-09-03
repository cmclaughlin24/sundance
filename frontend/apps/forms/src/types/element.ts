import type { ElementAttributes } from "./elementAttributes";
import type { IFormObject } from "./formObject";
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

export interface IElement extends IFormObject, HasRules {
  description: string;
  type: ElementType;
  tags: any[];
  attributes: ElementAttributes;
}
