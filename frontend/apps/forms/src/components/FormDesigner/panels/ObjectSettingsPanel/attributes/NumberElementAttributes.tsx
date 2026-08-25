import { checkElementType } from "@/utils/error";
import type { ElementAttributesComponent } from "../ObjectSettingsPanel";

export const NumberElementAttributes: ElementAttributesComponent = function ({
  element,
}) {
  checkElementType(element.type, "number");
};
