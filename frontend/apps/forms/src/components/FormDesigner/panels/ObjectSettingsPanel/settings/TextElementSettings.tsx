import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";
import { settingsStyle } from "./Settings.style";
import Box from "@mui/material/Box";
import type { TextElementAttributes } from "@/types/elementAttributes";
import TextField from "@mui/material/TextField";

export const TextElementSettings: ElementSettingsComponent = function ({
  element,
}) {
  checkElementType(element.type, "text");

  const attr = element.attributes as TextElementAttributes;

  return (
    <Box sx={settingsStyle.container}>
      <TextField type="number" label="Min. Length" />
      <TextField type="number" label="Max Length" />
    </Box>
  );
};
