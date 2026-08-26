import { FormDesignerProvider } from "@/store/formDesigner";
import type { IForm } from "@/types/form";
import type { IFormVersion } from "@/types/formVersion";
import Box from "@mui/material/Box";
import { formDesignerStyles } from "./FormDesigner.styles";
import { ToolboxPanel } from "./panels/ToolboxPanel/ToolboxPanel";
import { CanvasPanel } from "./panels/CanvasPanel/CanvasPanel";
import { ObjectSettingsPanel } from "./panels/ObjectSettingsPanel/ObjectSettingsPanel";
import { Border } from "@/constants/colors";

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
        <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
          <ToolboxPanel />
        </Box>
        <CanvasPanel />
        <Box sx={{ borderLeft: `1px solid ${Border.Primary}` }}>
          <ObjectSettingsPanel />
        </Box>
      </Box>
    </FormDesignerProvider>
  );
};
