import { DragDropProvider, DragOverlay } from "@dnd-kit/react";
import { findFormObjectPaletteItem } from "../panels/ToolboxPanel/palette";
import { PaletteItem } from "../panels/ToolboxPanel/PaletteItem";
import { useFormDesignerDispatch } from "@/store/formDesigner";
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
} from "@/store/formDesigner/events";
import { generatedID } from "@/utils/id";

const FormDesignerDragContext = createContext<FormDragEventData | null>(null);

export function useFormDragEvent() {
  return useContext(FormDesignerDragContext);
}

export const FormDesignerDragProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const [activeDragData, setActiveDragData] =
      useState<FormDragEventData | null>(null);
    const { dispatch } = useFormDesignerDispatch();

    const handlePaletteDragEnd = (
      dragData: PaletteDragEventData,
      dropData: PaletteDropEventData,
    ) => {
      let event: FormDesignerEvent;
      const id = generatedID();

      switch (dragData.objectType) {
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
            elementType: dragData.objectType,
            sectionId: dropData.parentId,
            position: dropData.position,
          } satisfies AddElementEvent;
          break;
      }

      dispatch(event);
    };

    const handleCanvasDragEnd = (
      dragData: CanvasElementDragEventData | CanvasSectionDragEventData,
      dropData: CanvasDropEventData,
    ) => {
      console.log(dragData, dropData);
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
                  const paletteItem = findFormObjectPaletteItem(
                    source.data.objectType,
                  );
                  return <PaletteItem item={paletteItem!} draggable={false} />;
                case FormDragEventSource.Canvas:
                  return;
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
