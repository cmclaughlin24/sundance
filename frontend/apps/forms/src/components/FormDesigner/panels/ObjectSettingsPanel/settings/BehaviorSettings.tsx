import Box from "@mui/material/Box";
import { settingsStyle } from "./Settings.style";
import FormControlLabel from "@mui/material/FormControlLabel";
import Switch from "@mui/material/Switch";

export const BehaviorSettings: React.FC = function () {
  return (
    <Box sx={settingsStyle.container}>
      <Box sx={{ flex: 1 }}>
        <FormControlLabel label="Required" control={<Switch />} />
      </Box>
      <Box sx={{ flex: 1 }}>
        <FormControlLabel label="Readonly" control={<Switch />} />
      </Box>
    </Box>
  );
};
