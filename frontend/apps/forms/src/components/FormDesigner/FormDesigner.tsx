import { FormDesignerProvider } from "@/store/formDesigner";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import Box from "@mui/material/Box";
import { formDesignerStyles } from "./FormDesigner.styles";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { ObjectSettingsPanel } from "./panels/ObjectSettingsPanel/ObjectSettingsPanel";
import { Border } from "@/constants/colors";
import { FormDesignerDragProvider } from "./providers/FormDesignerDragProvider";
import { KeyboardShortcutProvider } from "@/store/keyboardShortcut/KeyboardShortcutProvider";
import { FormDesignerKeyboardShortcuts } from "./FormDesignerKeyboardShorts";
import { ContextMenuProvider } from "../ContextMenu";

export interface FormDesignerProps {
  form: IForm;
  version: IFormVersion;
}

export const FormDesigner: React.FC<FormDesignerProps> = function ({
  form,
  version,
}) {
  return (
    <KeyboardShortcutProvider>
      <ContextMenuProvider>
        <FormDesignerProvider form={form} version={version}>
          <FormDesignerDragProvider>
            <FormDesignerKeyboardShortcuts>
              <Box sx={formDesignerStyles.container}>
                <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
                  <ToolboxPanel />
                </Box>
                <CanvasPanel />
                <Box sx={{ borderLeft: `1px solid ${Border.Primary}` }}>
                  <ObjectSettingsPanel />
                </Box>
              </Box>
            </FormDesignerKeyboardShortcuts>
          </FormDesignerDragProvider>
        </FormDesignerProvider>
      </ContextMenuProvider>
    </KeyboardShortcutProvider>
  );
};
