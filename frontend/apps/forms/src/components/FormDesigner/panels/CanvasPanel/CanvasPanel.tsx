import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";
import { PageList } from "./lists/PageList";
import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
  useFormDesignerUndo,
  useFormPagesSnapshot,
} from "@/store/formDesigner";
import { FormSummary } from "../../common/FormSummary";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import RedoIcon from "@mui/icons-material/Redo";
import UndoIcon from "@mui/icons-material/Undo";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, } from "react";
import type {
  PasteElementEvent,
  PasteSectionEvent,
} from "@/store/formDesigner/events";
import { type ClipboardData, ClipboardEventType } from "@/types/clipboard";

export const CanvasPanel: React.FC = function () {
  const { undo, redo } = useFormDesignerUndo();
  const { dispatch } = useFormDesignerDispatch();
  const { selected } = useFormDesignerSelect();
  const pages = useFormPagesSnapshot();

  useEffect(() => {
    const handleKeyDown = async (event: KeyboardEvent) => {
      if (!(event.ctrlKey || event.metaKey) || event.key !== "v") {
        return;
      }

      try {
        const text = await navigator.clipboard.readText();
        const data: ClipboardData = JSON.parse(text);

        switch (data.type) {
          case ClipboardEventType.CopyElement:
            if (selected?.type !== "section") {
              return;
            }

            dispatch({
              type: "PasteElement",
              element: data.element,
              targetSectionId: selected.id,
            } satisfies PasteElementEvent);
            break;
          case ClipboardEventType.CopySection:
            dispatch({
              type: "PasteSection",
              section: data.section,
              targetPageId: data.pageId,
            } satisfies PasteSectionEvent);
            break;
        }
      } catch {
        return;
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [selected, pages, dispatch]);

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
        </Box>
      </Box>
      <PageList pages={pages} />
    </Panel>
  );
};
