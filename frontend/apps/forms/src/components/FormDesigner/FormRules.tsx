import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";

const styles: Styles = {
  container: {
    display: "grid",
    gridTemplateColumns: "minmax(296px, 18.5rem) auto",
    border: `1px solid ${Border.Primary}`,
  },
};

export const FormRules: React.FC = function () {
  return (
    <Box sx={styles.container}>
      <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
        <ToolboxPanel
          palette={[]}
          helpText="Drag the form rules onto the canvas."
        />
      </Box>
      <CanvasPanel></CanvasPanel>
    </Box>
  );
};
