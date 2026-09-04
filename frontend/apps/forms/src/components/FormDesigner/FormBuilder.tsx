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
import { useFormPagesSnapshot, type SelectedItem } from "@/store/formDesigner";
import { FormDesignerContextMenu } from "./FormDesignerContextMenu";
import { ClipboardEventType, type PagesClipboardData } from "@/types/clipboard";
import { PageList } from "./panels/CanvasPanel/lists/PageList";
import {
  FORM_OBJECT_PALETTE,
  type PaletteItemType,
} from "./panels/ToolboxPanel/constants/formObjectPalette";
import type { IPaletteCategory } from "./panels/ToolboxPanel/palette";

export interface FormDesignerProps {}

export const FormBuilder: React.FC<FormDesignerProps> = function () {
  const pages = useFormPagesSnapshot();

  const handleCopy = () => {
    // FIXME: When multi-page support is enabled, need to move into PageItem.
    const data: PagesClipboardData = {
      type: ClipboardEventType.CopyPage,
      page: pages[0],
    };

    navigator.clipboard.writeText(JSON.stringify(data));
  };

  return (
    <KeyboardShortcutProvider>
      <ContextMenuProvider>
        <FormDesignerDragProvider>
          <FormDesignerKeyboardShortcuts>
            <Box sx={formBuilderStyles.container}>
              <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
                <ToolboxPanel
                  palette={
                    FORM_OBJECT_PALETTE as IPaletteCategory<PaletteItemType>[]
                  }
                  helpText="Drag the form elements into the preferred section on the canvas."
                />
              </Box>
              <CanvasPanel onCopy={handleCopy}>
                <PageList pages={pages} />
              </CanvasPanel>
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
