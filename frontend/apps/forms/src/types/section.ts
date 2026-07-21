import type { IElement } from "./element";
import type { HasKey } from "./hasKey";
import type { HasName } from "./hasName";
import type { HasPosition } from "./hasPosition";
import type { HasRules } from "./rule";

export interface ISection extends HasKey, HasPosition, HasName, HasRules {
  id: string;
  elements: IElement[];
}
