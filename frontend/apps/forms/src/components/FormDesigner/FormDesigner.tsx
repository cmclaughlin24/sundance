import { FormDesignerProvider } from "@/store/formDesigner";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import Box from "@mui/material/Box";
import { formDesignerStyles } from "./FormDesigner.styles";
import { LibraryPanel } from "./panels/LibraryPanel/LibraryPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { ObjectSettingsPanel } from "./panels/ObjectSettingsPanel/ObjectSettingsPanel";

export interface FormDesignerProps {
  form: IForm;
  version: IFormVersion;
}

export const FormDesigner: React.FC<FormDesignerProps> = function ({
  form,
  version,
}) {
  return (
    <FormDesignerProvider form={form} version={version}>
      <Box sx={formDesignerStyles.container}>
        <LibraryPanel />
        <CanvasPanel />
        <ObjectSettingsPanel />
      </Box>
    </FormDesignerProvider>
  );
};
