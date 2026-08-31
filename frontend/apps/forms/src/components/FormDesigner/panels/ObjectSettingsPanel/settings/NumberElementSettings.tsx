import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";

export const NumberElementSettings: ElementSettingsComponent = function ({
  element,
}) {
  checkElementType(element.type, "number");
  return <></>;
};
