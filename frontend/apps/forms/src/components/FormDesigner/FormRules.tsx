import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { FORM_RULES_PALETTE } from "./panels/ToolboxPanel/constants/formRulesPalette";
import type { IPaletteCategory } from "./panels/ToolboxPanel/palette";
import { RuleList } from "./panels/CanvasPanel/lists/RuleList";
import { FormRuleDragProvider } from "./providers/FormRuleDragProvider";
import type { RuleItemDragType } from "./types/formDragEvent";
import type { RuleType } from "@/types/rule";
import { useFormSnapshot } from "@/store/formDesigner";

const styles: Styles = {
  container: {
    display: "grid",
    gridTemplateColumns: "minmax(296px, 18.5rem) auto",
    border: `1px solid ${Border.Primary}`,
    height: "100%",
  },
  rules: {
    width: "100%",
    maxWidth: "800px",
    alignSelf: "center",
    display: "flex",
    justifyContent: "center",
  },
};

export const FormRules: React.FC = function () {
  const { rules } = useFormSnapshot();

  return (
    <FormRuleDragProvider>
      <Box sx={styles.container}>
        <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
          <ToolboxPanel
            palette={
              FORM_RULES_PALETTE as IPaletteCategory<
                RuleType,
                RuleItemDragType
              >[]
            }
            helpText="Drag the form rules onto the canvas."
          />
        </Box>
        <CanvasPanel>
          <Box sx={styles.rules}>
            <RuleList rules={rules} />
          </Box>
        </CanvasPanel>
      </Box>
    </FormRuleDragProvider>
  );
};
