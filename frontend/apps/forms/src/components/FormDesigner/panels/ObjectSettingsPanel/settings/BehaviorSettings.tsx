import Box from "@mui/material/Box";
import { settingsStyle } from "./Settings.style";
import FormControlLabel from "@mui/material/FormControlLabel";
import Switch from "@mui/material/Switch";
import type { ElementAttributes } from "@/types/elementAttributes";
import type { IElement } from "@/types/element";
import type { ChangeEvent } from "react";

export type BehaviorSettingsEvent = Partial<
  Pick<ElementAttributes, "isRequired" | "isReadOnly">
>;

export const BehaviorSettings: React.FC<{
  element: IElement;
  onChange: (event: BehaviorSettingsEvent) => void;
}> = function ({ element, onChange }) {
  const attr = element.attributes;

  const handleChange = (field: "isRequired" | "isReadOnly") => {
    return (_event: ChangeEvent<HTMLInputElement>, checked: boolean) => {
      onChange({ [field]: checked });
    };
  };

  return (
    <Box sx={settingsStyle.container}>
      <Box sx={{ flex: 1 }}>
        <FormControlLabel
          label="Required"
          control={
            <Switch
              checked={attr.isRequired}
              onChange={handleChange("isRequired")}
              color="success"
            />
          }
        />
      </Box>
      <Box sx={{ flex: 1 }}>
        <FormControlLabel
          label="Readonly"
          control={
            <Switch
              checked={attr.isReadOnly}
              onChange={handleChange("isReadOnly")}
              color="success"
            />
          }
        />
      </Box>
    </Box>
  );
};
