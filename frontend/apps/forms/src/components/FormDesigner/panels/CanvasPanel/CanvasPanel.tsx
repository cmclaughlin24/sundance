import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";
import { PageList } from "./lists/PageList";
import {
  useFormDesignerUndo,
  useFormPagesSnapshot,
} from "@/store/formDesigner";
import { FormSummary } from "../../common/FormSummary";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import RedoIcon from "@mui/icons-material/Redo";
import UndoIcon from "@mui/icons-material/Undo";

export const CanvasPanel: React.FC = function () {
  const { undo, redo } = useFormDesignerUndo();
  const pages = useFormPagesSnapshot();

  return (
    <Panel sx={canvasPanelStyles.canvas}>
      <Box sx={canvasPanelStyles.toolbar}>
        <FormSummary pages={pages} />
        <Box sx={canvasPanelStyles.buttons}>
          <IconButton
            size="small"
            aria-label="undo"
            data-testid="undo-btn"
            onClick={undo}
          >
            <UndoIcon fontSize="inherit" />
          </IconButton>
          <IconButton
            size="small"
            aria-label="redo"
            data-testid="redo-btn"
            onClick={redo}
          >
            <RedoIcon fontSize="inherit" />
          </IconButton>
        </Box>
      </Box>
      <PageList pages={pages} />
    </Panel>
  );
};
