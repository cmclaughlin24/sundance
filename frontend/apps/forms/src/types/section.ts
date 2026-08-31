import type { IElement } from "./element";
import type { IFormObject } from "./formObject";
import type { HasRules } from "./rule";

export interface ISection extends IFormObject, HasRules {
  elements: IElement[];
}
