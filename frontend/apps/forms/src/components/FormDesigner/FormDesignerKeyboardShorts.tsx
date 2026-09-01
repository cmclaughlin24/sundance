import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
  useFormDesignerUndo,
  useFormPagesSnapshot,
} from "@/store/formDesigner";
import type {
  CutElementEvent,
  CutSectionEvent,
  FormDesignerEvent,
  PasteElementEvent,
  PasteSectionEvent,
  RemoveElementEvent,
  RemoveSectionEvent,
} from "@/store/formDesigner/events";
import { useKeyboardShortcut } from "@/store/keyboardShortcut/useKeyboardShortcut";
import {
  ClipboardEventType,
  type ClipboardData,
  type ElementClipboardData,
  type SectionClipboardData,
} from "@/types/clipboard";

export const FormDesignerKeyboardShortcuts: React.FC<
  React.PropsWithChildren<{}>
> = function ({ children }) {
  const pages = useFormPagesSnapshot();
  const { undo, redo } = useFormDesignerUndo();
  const { dispatch } = useFormDesignerDispatch();
  const { selected } = useFormDesignerSelect();

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
            case ClipboardEventType.CutElement:
              if (selected?.type !== "section") {
                return;
              }

              dispatch({
                type: "PasteElement",
                element: data.element,
                targetSectionId: selected.item.id,
                clipboardOp: data.type,
              } satisfies PasteElementEvent);
              break;
            case ClipboardEventType.CopySection:
            case ClipboardEventType.CutSection:
              // FIXME: When multi-page support is enabled, need to pull the page ID from the current selection.
              dispatch({
                type: "PasteSection",
                section: data.section,
                targetPageId: pages[0].id,
                clipboardOp: data.type,
              } satisfies PasteSectionEvent);
              break;
          }

          if (
            data.type === ClipboardEventType.CutElement ||
            data.type === ClipboardEventType.CutSection
          ) {
            navigator.clipboard.writeText("");
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

  useKeyboardShortcut(
    {
      name: "Cut",
      combination: { key: "x", ctrlOrMeta: true },
      action: () => {
        if (!selected) {
          return;
        }

        let data: ClipboardData;
        let event: FormDesignerEvent;

        switch (selected.type) {
          case "section":
            data = {
              type: ClipboardEventType.CutSection,
              section: selected.item,
            } satisfies SectionClipboardData;
            event = {
              type: "CutSection",
              sectionId: selected.item.id,
            } satisfies CutSectionEvent;
            break;
          case "element":
            data = {
              type: ClipboardEventType.CutElement,
              element: selected.item,
            } satisfies ElementClipboardData;
            event = {
              type: "CutElement",
              elementId: selected.item.id,
            } satisfies CutElementEvent;
            break;
        }

        dispatch(event!);
        navigator.clipboard.writeText(JSON.stringify(data!));
      },
    },
    [dispatch, selected],
  );

  return children;
};
