import { checkElementType } from "@/utils/error";
import type { ElementAttributesComponent } from "../ObjectSettingsPanel";

export const TextElementAttributes: ElementAttributesComponent = function ({
  element,
}) {
  checkElementType(element.type, "text");
};
