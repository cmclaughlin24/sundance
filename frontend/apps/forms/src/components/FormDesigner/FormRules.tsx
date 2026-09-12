import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { FORM_RULES_PALETTE } from "./panels/ToolboxPanel/constants/formRulesPalette";
import type { IPaletteCategory } from "./panels/ToolboxPanel/palette";
import { RuleList } from "./panels/CanvasPanel/lists/RuleList/RuleList";
import { FormRuleDragProvider } from "./providers/FormRuleDragProvider";
import type { RuleItemDragType } from "./types/formDragEvent";
import type { IFlatRule, RuleType } from "@/types/rule";
import { useFormSnapshot } from "@/store/formDesigner";
import {
  ContextMenu,
  ContextMenuProvider,
  useContextMenuDispatch,
} from "../ContextMenu";
import { FORMS_HUB_PORTAL_REF } from "@/constants/portalRef";
import { RuleContextMenu } from "./menus/RuleContextMenu";
import type { MouseEventHandler } from "react";
import { RuleKeyboardShortcuts } from "./providers/RuleKeyboardShortcuts";

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
  return (
    <FormRuleDragProvider>
      <ContextMenuProvider>
        <RuleKeyboardShortcuts>
          <Component />
        </RuleKeyboardShortcuts>
      </ContextMenuProvider>
    </FormRuleDragProvider>
  );
};

const Component: React.FC = function () {
  const { rules } = useFormSnapshot();
  const { open } = useContextMenuDispatch();

  const handleContext: MouseEventHandler<HTMLDivElement> = (event) => {
    event.preventDefault();
    event.stopPropagation();
    open({ position: { x: event.clientX, y: event.clientY }, data: null });
  };

  return (
    <>
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
        <CanvasPanel onContextMenu={handleContext}>
          <Box sx={styles.rules}>
            <RuleList rules={rules} />
          </Box>
        </CanvasPanel>
        <ContextMenu container={document.getElementById(FORMS_HUB_PORTAL_REF)!}>
          {(data: unknown) => {
            const target = data as IFlatRule | undefined;
            return <RuleContextMenu target={target} />;
          }}
        </ContextMenu>
      </Box>
    </>
  );
};
