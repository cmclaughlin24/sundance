import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";
import { useFormDesignerHistory, useFormSnapshot } from "@/store/formDesigner";
import { FormSummary } from "../../common/FormSummary";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import RedoIcon from "@mui/icons-material/Redo";
import UndoIcon from "@mui/icons-material/Undo";
import Tooltip from "@mui/material/Tooltip";
import ContentCopy from "@mui/icons-material/ContentCopy";
import type { MouseEventHandler } from "react";

export const CanvasPanel: React.FC<
  React.PropsWithChildren<{
    onCopy?: () => void;
    onContextMenu?: MouseEventHandler<HTMLDivElement>;
  }>
> = function ({ children, onCopy, onContextMenu }) {
  const { undo, canUndo, redo, canRedo } = useFormDesignerHistory();
  const { version, rules } = useFormSnapshot();

  return (
    <Panel sx={canvasPanelStyles.canvas} onContextMenu={onContextMenu}>
      <Box sx={canvasPanelStyles.toolbar}>
        <FormSummary pages={version.pages} rules={rules} />
        <Box sx={canvasPanelStyles.buttons}>
          <Tooltip title="Undo">
            <IconButton
              size="small"
              aria-label="undo"
              data-testid="undo-btn"
              onClick={undo}
              disabled={!canUndo}
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
              disabled={!canRedo}
            >
              <RedoIcon fontSize="inherit" />
            </IconButton>
          </Tooltip>
          {onCopy && (
            <Tooltip title="Copy">
              <IconButton
                size="small"
                aria-label="copy"
                data-testid="copy-btn"
                onClick={onCopy}
              >
                <ContentCopy fontSize="inherit" />
              </IconButton>
            </Tooltip>
          )}
        </Box>
      </Box>
      {children}
    </Panel>
  );
};
