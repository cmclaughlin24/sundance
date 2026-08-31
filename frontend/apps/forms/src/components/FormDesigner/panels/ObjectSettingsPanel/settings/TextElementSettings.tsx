import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";

export const TextElementSettings: ElementSettingsComponent = function ({
  element,
}) {
  checkElementType(element.type, "text");
  return <></>;
};
