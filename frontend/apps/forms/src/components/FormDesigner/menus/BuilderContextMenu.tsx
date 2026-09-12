import {
  useFormDesignerDispatch,
  useFormDesignerHistory,
  useFormPagesSnapshot,
  useFormSnapshot,
  type SelectedItem,
} from "@/store/formDesigner";
import type {
  CutElementEvent,
  CutSectionEvent,
  FormDesignerEvent,
  PasteElementEvent,
  PastePageEvent,
  PasteSectionEvent,
  RemoveElementEvent,
  RemoveSectionEvent,
} from "@/store/formDesigner/events";
import { ContextMenu, useContextMenuDispatch } from "../../ContextMenu";
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
import { isDraftVersion, versionToRequest } from "@/utils/form";
import { useFormsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";

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

export const BuilderContextMenu: React.FC<{ target: SelectedItem }> =
  function ({ target }) {
    const { undo, canUndo, redo, canRedo, commit } = useFormDesignerHistory();
    const { dispatch } = useFormDesignerDispatch();
    const { close } = useContextMenuDispatch();
    const { form, version, rules } = useFormSnapshot();
    const formsService = useFormsService();
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
      if (!clipboardData) {
        return;
      }

      let event: FormDesignerEvent;

      switch (clipboardData.type) {
        case ClipboardEventType.CopyElement:
        case ClipboardEventType.CutElement:
          if (target.type !== "section") {
            return;
          }
          event = {
            type: "PasteElement",
            element: clipboardData.element,
            targetSectionId: target.item.id,
            clipboardOp: clipboardData.type,
          } satisfies PasteElementEvent;
          break;
        case ClipboardEventType.CopySection:
        case ClipboardEventType.CutSection:
          event = {
            type: "PasteSection",
            section: clipboardData.section,
            targetPageId: pages[0].id,
            clipboardOp: clipboardData.type,
          } satisfies PasteSectionEvent;
          break;
        case ClipboardEventType.CopyPage:
          event = {
            type: "PastePage",
            page: clipboardData.page,
          } satisfies PastePageEvent;
          break;
      }

      dispatch(event!);

      if (
        clipboardData.type === ClipboardEventType.CutElement ||
        clipboardData.type === ClipboardEventType.CutSection
      ) {
        navigator.clipboard.writeText("");
      }

      close();
    };

    const handleSaveDraft = async () => {
      if (!isDraftVersion(version.status)) {
        throw new Error("cannot update a non-draft-version");
      }

      const request = versionToRequest(version, rules);
      const updated = await formsService.updateFormVersion(
        form.id,
        version.id,
        request,
        {
          tenantId: TENANT_ID,
          token: "placeholder",
        },
      );
      commit(updated);
      close();
    };

    const handleDelete = () => {
      switch (target.type) {
        case "section":
          dispatch({
            type: "RemoveSection",
            id: target.item.id,
          } satisfies RemoveSectionEvent);
          break;
        default:
          dispatch({
            type: "RemoveElement",
            id: target.item.id,
          } satisfies RemoveElementEvent);
          break;
      }
      close();
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
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={undo}
          disabled={!canUndo}
        >
          <Typography>Undo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+z</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={redo}
          disabled={!canRedo}
        >
          <Typography>Redo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+Shift+z</Typography>
        </ContextMenu.Button>
        <Divider sx={{ my: 1 }} />
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={handleSaveDraft}
        >
          <Typography>Save Draft</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button sx={styles.btnWithShortcut} onClick={handleDelete}>
          <Typography>Delete</Typography>
          <Typography sx={styles.shortcutText}>Del</Typography>
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

  return false;
}
