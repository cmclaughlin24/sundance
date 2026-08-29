import { DragDropProvider, DragOverlay } from "@dnd-kit/react";
import { findPaletteItem } from "../panels/ToolboxPanel/palette";
import { PaletteItem } from "../panels/ToolboxPanel/PaletteItem";
import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
} from "@/store/formDesigner";
import {
  FormDragEventSource,
  type CanvasElementDragEventData,
  type CanvasSectionDragEventData,
  type FormDragEventData,
  type PaletteDragEventData,
} from "../types/formDragEvent";
import type {
  CanvasDropEventData,
  FormDropEventData,
  PaletteDropEventData,
} from "../types/formDropEvent";
import { createContext, useContext, useState } from "react";
import type {
  AddElementEvent,
  AddSectionEvent,
  FormDesignerEvent,
  MoveElementEvent,
  MoveSectionEvent,
} from "@/store/formDesigner/events";
import { generatedID } from "@/utils/id";
import { SectionItem } from "../panels/CanvasPanel/lists/SectionItem";
import { ElementItem } from "../panels/CanvasPanel/lists/ElementItem";

const FormDesignerDragContext = createContext<FormDragEventData | null>(null);

export function useFormDragEvent() {
  return useContext(FormDesignerDragContext);
}

export const FormDesignerDragProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const [activeDragData, setActiveDragData] =
      useState<FormDragEventData | null>(null);
    const { dispatch } = useFormDesignerDispatch();
    const { select } = useFormDesignerSelect();

    const handlePaletteDragEnd = (
      dragData: PaletteDragEventData,
      dropData: PaletteDropEventData,
    ) => {
      let event: FormDesignerEvent;
      const id = generatedID();

      switch (dragData.type) {
        case "section":
          event = {
            type: "AddSection",
            id,
            pageId: dropData.parentId,
            position: dropData.position,
          } satisfies AddSectionEvent;
          break;
        default:
          event = {
            type: "AddElement",
            id,
            elementType: dragData.type,
            sectionId: dropData.parentId,
            position: dropData.position,
          } satisfies AddElementEvent;
          break;
      }

      dispatch(event);
      select({ type: dragData.type, id });
    };

    const handleCanvasDragEnd = (
      dragData: CanvasElementDragEventData | CanvasSectionDragEventData,
      dropData: CanvasDropEventData,
    ) => {
      if (dragData.type === "section") {
        dispatch({
          type: "MoveSection",
          sectionId: dragData.section.id,
          targetPageId: dropData.parentId,
          position: dropData.position,
        } satisfies MoveSectionEvent);
      } else {
        dispatch({
          type: "MoveElement",
          elementId: dragData.element.id,
          targetSectionId: dropData.parentId,
          position: dropData.position,
        } satisfies MoveElementEvent);
      }
    };

    return (
      <FormDesignerDragContext value={activeDragData}>
        <DragDropProvider
          onDragStart={(event) =>
            setActiveDragData(event.operation.source?.data as FormDragEventData)
          }
          onDragEnd={(event) => {
            setActiveDragData(null);

            if (event.canceled) {
              return;
            }

            const { source, target } = event.operation;
            const dragData = source?.data as FormDragEventData;
            const dropData = target?.data as FormDropEventData;

            if (!dropData) {
              return;
            }

            switch (dragData.source) {
              case FormDragEventSource.Palette:
                handlePaletteDragEnd(
                  dragData,
                  dropData as PaletteDropEventData,
                );
                break;
              case FormDragEventSource.Canvas:
                handleCanvasDragEnd(dragData, dropData as CanvasDropEventData);
                break;
              default:
                throw new Error(
                  "failed to handle onDragEnd event; unknown FormDragEventData source",
                );
            }
          }}
        >
          {children}
          <DragOverlay>
            {(source) => {
              const data = source.data as FormDragEventData;

              switch (data.source) {
                case FormDragEventSource.Palette:
                  const paletteItem = findPaletteItem(source.data.type);
                  return <PaletteItem item={paletteItem!} draggable={false} />;
                case FormDragEventSource.Canvas:
                  if (data.type === "section") {
                    return (
                      <SectionItem
                        section={data.section}
                        parentId={data.fromPageId}
                        draggable={false}
                      />
                    );
                  }

                  return (
                    <ElementItem
                      element={data.element}
                      parentId={data.fromSectionId}
                      draggable={false}
                    />
                  );
                default:
                  throw new Error(
                    "cannot display draggable item; unknown FormDragEventData source",
                  );
              }
            }}
          </DragOverlay>
        </DragDropProvider>
      </FormDesignerDragContext>
    );
  };
