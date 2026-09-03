import Box from "@mui/material/Box";
import { formBuilderStyles } from "./FormBuilder.styles";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { ObjectSettingsPanel } from "./panels/ObjectSettingsPanel/ObjectSettingsPanel";
import { Border } from "@/constants/colors";
import { FormDesignerDragProvider } from "./providers/FormDesignerDragProvider";
import { KeyboardShortcutProvider } from "@/store/keyboardShortcut/KeyboardShortcutProvider";
import { FormDesignerKeyboardShortcuts } from "./FormDesignerKeyboardShorts";
import { ContextMenu, ContextMenuProvider } from "../ContextMenu";
import { FORMS_HUB_PORTAL_REF } from "@/constants/portalRef";
import type { SelectedItem } from "@/store/formDesigner";
import { FormDesignerContextMenu } from "./FormDesignerContextMenu";

export interface FormDesignerProps {}

export const FormBuilder: React.FC<FormDesignerProps> = function () {
  return (
    <KeyboardShortcutProvider>
      <ContextMenuProvider>
        <FormDesignerDragProvider>
          <FormDesignerKeyboardShortcuts>
            <Box sx={formBuilderStyles.container}>
              <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
                <ToolboxPanel />
              </Box>
              <CanvasPanel />
              <Box sx={{ borderLeft: `1px solid ${Border.Primary}` }}>
                <ObjectSettingsPanel />
              </Box>
            </Box>
            <ContextMenu
              container={document.getElementById(FORMS_HUB_PORTAL_REF)!}
            >
              {(data: unknown) => {
                const target = data as SelectedItem;
                return <FormDesignerContextMenu target={target} />;
              }}
            </ContextMenu>
          </FormDesignerKeyboardShortcuts>
        </FormDesignerDragProvider>
      </ContextMenuProvider>
    </KeyboardShortcutProvider>
  );
};
