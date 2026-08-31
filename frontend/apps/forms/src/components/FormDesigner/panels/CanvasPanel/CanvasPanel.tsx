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
import type {
  PasteElementEvent,
  PasteSectionEvent,
  RemoveElementEvent,
  RemoveSectionEvent,
} from "@/store/formDesigner/events";
import { type ClipboardData, ClipboardEventType } from "@/types/clipboard";
import { useKeyboardShortcut } from "@/store/keyboardShortcut/useKeyboardShortcut";

export const CanvasPanel: React.FC = function () {
  const { undo, redo } = useFormDesignerUndo();
  const { dispatch } = useFormDesignerDispatch();
  const { selected } = useFormDesignerSelect();
  const pages = useFormPagesSnapshot();

  useKeyboardShortcut(
    {
      name: "Paste",
      combination: { key: "v", ctrlOrMeta: true },
      action: async () => {
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
                targetSectionId: selected.item.id,
              } satisfies PasteElementEvent);
              break;
            case ClipboardEventType.CopySection:
              // FIXME: When multi-page support is enabled, need to pull the page ID from the current selection.
              dispatch({
                type: "PasteSection",
                section: data.section,
                targetPageId: pages[0].id,
              } satisfies PasteSectionEvent);
              break;
          }
        } catch {
          return;
        }
      },
    },
    [selected, pages, dispatch],
  );

  useKeyboardShortcut(
    { name: "Undo", combination: { key: "z", ctrlOrMeta: true }, action: undo },
    [undo],
  );

  useKeyboardShortcut(
    {
      name: "Redo",
      combination: { key: "z", ctrlOrMeta: true, shift: true },
      action: redo,
    },
    [redo],
  );

  useKeyboardShortcut(
    {
      name: "Delete",
      combination: { key: "Delete" },
      action: () => {
        if (!selected) {
          return;
        }

        switch (selected.type) {
          case "section":
            dispatch({
              type: "RemoveSection",
              sectionId: selected.item.id,
            } satisfies RemoveSectionEvent);
            break;
          default:
            dispatch({
              type: "RemoveElement",
              elementId: selected.item.id,
            } satisfies RemoveElementEvent);
            break;
        }
      },
    },
    [dispatch, selected],
  );

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
