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
import Tooltip from "@mui/material/Tooltip";
import ContentCopy from "@mui/icons-material/ContentCopy";
import { ClipboardEventType, type PagesClipboardData } from "@/types/clipboard";

export const CanvasPanel: React.FC = function () {
  const { undo, redo } = useFormDesignerUndo();
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
    <Panel sx={canvasPanelStyles.canvas}>
      <Box sx={canvasPanelStyles.toolbar}>
        <FormSummary pages={pages} />
        <Box sx={canvasPanelStyles.buttons}>
          <Tooltip title="Undo">
            <IconButton
              size="small"
              aria-label="undo"
              data-testid="undo-btn"
              onClick={undo}
            >
              <UndoIcon fontSize="inherit" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Redo">
            <IconButton
              size="small"
              aria-label="redo"
              data-testid="redo-btn"
              onClick={redo}
            >
              <RedoIcon fontSize="inherit" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Copy">
            <IconButton
              size="small"
              aria-label="copy"
              data-testid="copy-btn"
              onClick={handleCopy}
            >
              <ContentCopy fontSize="inherit" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      <PageList pages={pages} />
    </Panel>
  );
};
