import type { IFormObject } from "./formObject";
import type { HasRules } from "./rule";
import type { ISection } from "./section";

export interface IPage extends IFormObject, HasRules {
    sections: ISection[];
}
