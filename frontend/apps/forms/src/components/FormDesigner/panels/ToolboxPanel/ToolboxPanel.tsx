import { Panel } from "@/components/layout/Panel";
import Typography from "@mui/material/Typography";
import { PALLETTE } from "./pallette";
import { PalletteCategory } from "./PalletteCategory";

export const ToolboxPanel: React.FC = function () {
  return (
    <Panel sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
      <Panel.Title title="Library" />
      <Typography>
        Drag the form elements into the preferred section on the canvas.
      </Typography>
      {PALLETTE.map((category) => (
        <PalletteCategory category={category} key={category.label} />
      ))}
    </Panel>
  );
};
