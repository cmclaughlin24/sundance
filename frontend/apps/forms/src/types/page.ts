import type { HasKey } from "./hasKey";
import type { HasName } from "./hasName";
import type { HasPosition } from "./hasPosition";
import type { HasRules } from "./rule";
import type { ISection } from "./section";

export interface IPage extends HasKey, HasPosition, HasName, HasRules {
    id: string;
    sections: ISection[];
}
