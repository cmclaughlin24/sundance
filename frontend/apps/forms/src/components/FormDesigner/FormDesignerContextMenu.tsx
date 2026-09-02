import {
  useFormDesignerDispatch,
  useFormDesignerUndo,
  useFormPagesSnapshot,
  type SelectedItem,
} from "@/store/formDesigner";
import type {
  CutElementEvent,
  CutSectionEvent,
  FormDesignerEvent,
  PasteElementEvent,
  PastePageEvent,
  PasteSectionEvent,
} from "@/store/formDesigner/events";
import { ContextMenu, useContextMenuDispatch } from "../ContextMenu";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import type { Styles } from "@/types/styles";
import { useState, useEffect } from "react";
import {
  ClipboardEventType,
  type ElementClipboardData,
  type SectionClipboardData,
  type PagesClipboardData,
  type ClipboardData,
} from "@/types/clipboard";

const styles: Styles = {
  btnWithShortcut: {
    display: "flex",
    justifyContent: "space-between",
  },
  shortcutText: {
    fontSize: "0.75rem",
    color: "#4B4444",
  },
};

export const FormDesignerContextMenu: React.FC<{ target: SelectedItem }> =
  function ({ target }) {
    const { undo, redo } = useFormDesignerUndo();
    const { dispatch } = useFormDesignerDispatch();
    const { close } = useContextMenuDispatch();
    const pages = useFormPagesSnapshot();
    const [clipboardData, setClipboardData] = useState<ClipboardData | null>(
      null,
    );

    useEffect(() => {
      navigator.clipboard
        .readText()
        .then((text) => {
          try {
            setClipboardData(JSON.parse(text) as ClipboardData);
          } catch {
            setClipboardData(null);
          }
        })
        .catch(() => {
          setClipboardData(null);
        });
    }, [target]);

    const handleCopy = () => {
      let data:
        ElementClipboardData | SectionClipboardData | PagesClipboardData;

      switch (target.type) {
        case "element":
          data = {
            type: ClipboardEventType.CopyElement,
            element: target.item,
          } satisfies ElementClipboardData;
          break;
        case "section":
          data = {
            type: ClipboardEventType.CopySection,
            section: target.item,
          } satisfies SectionClipboardData;
          break;
        case "page":
          data = {
            type: ClipboardEventType.CopyPage,
            page: target.item,
          } satisfies PagesClipboardData;
          break;
      }

      navigator.clipboard.writeText(JSON.stringify(data!));
      close();
    };

    const handleCut = () => {
      let event: FormDesignerEvent;
      let data: ClipboardData;

      switch (target.type) {
        case "element": {
          data = {
            type: ClipboardEventType.CutElement,
            element: target.item,
          } satisfies ElementClipboardData;
          event = {
            type: "CutElement",
            elementId: target.item.id,
          } satisfies CutElementEvent;
          break;
        }
        case "section": {
          data = {
            type: ClipboardEventType.CutSection,
            section: target.item,
          } satisfies SectionClipboardData;
          event = {
            type: "CutSection",
            sectionId: target.item.id,
          } satisfies CutSectionEvent;
          navigator.clipboard.writeText(JSON.stringify(data));
          break;
        }
        case "page":
          return;
      }

      dispatch(event!);
      navigator.clipboard.writeText(JSON.stringify(data!));
      close();
    };

    const handlePaste = async () => {
      try {
        const text = await navigator.clipboard.readText();
        const data: ClipboardData = JSON.parse(text);
        let event: FormDesignerEvent;

        switch (data.type) {
          case ClipboardEventType.CopyElement:
          case ClipboardEventType.CutElement:
            if (target.type !== "section") {
              return;
            }
            event = {
              type: "PasteElement",
              element: data.element,
              targetSectionId: target.item.id,
              clipboardOp: data.type,
            } satisfies PasteElementEvent;
            break;
          case ClipboardEventType.CopySection:
          case ClipboardEventType.CutSection:
            event = {
              type: "PasteSection",
              section: data.section,
              targetPageId: pages[0].id,
              clipboardOp: data.type,
            } satisfies PasteSectionEvent;
            break;
          case ClipboardEventType.CopyPage:
            event = {
              type: "PastePage",
              page: data.page,
            } satisfies PastePageEvent;
            break;
        }

        dispatch(event!);

        if (
          data.type === ClipboardEventType.CutElement ||
          data.type === ClipboardEventType.CutSection
        ) {
          navigator.clipboard.writeText("");
        }

        close();
      } catch {
        return;
      }
    };

    const canPasteItem = !canPaste(clipboardData, target);

    return (
      <>
        <ContextMenu.Button onClick={handleCopy}>Copy</ContextMenu.Button>
        <ContextMenu.Button sx={styles.btnWithShortcut} onClick={handleCut}>
          <Typography>Cut</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+x</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={handlePaste}
          disabled={canPasteItem}
        >
          <Typography>Paste</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+v</Typography>
        </ContextMenu.Button>
        <Divider sx={{ my: 1 }} />
        <ContextMenu.Button sx={styles.btnWithShortcut} onClick={undo}>
          <Typography>Undo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+z</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button sx={styles.btnWithShortcut} onClick={redo}>
          <Typography>Redo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+Shift+z</Typography>
        </ContextMenu.Button>
      </>
    );
  };

function canPaste(
  clipboardData: ClipboardData | null,
  target: SelectedItem,
): boolean {
  if (!clipboardData) {
    return false;
  }

  switch (clipboardData.type) {
    case ClipboardEventType.CopyElement:
    case ClipboardEventType.CutElement:
      return target.type === "section";
    case ClipboardEventType.CopySection:
    case ClipboardEventType.CutSection:
    case ClipboardEventType.CopyPage:
      return true;
  }
}
