import { DragDropProvider, DragOverlay } from "@dnd-kit/react";
import { findPaletteItem } from "../panels/ToolboxPanel/palette";
import { PaletteItem } from "../panels/ToolboxPanel/PaletteItem";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import {
  FormDragEventSource,
  type FormDragEventData,
  type PaletteDragEventData,
} from "../types/formDragEvent";
import type {
  FormDropEventData,
  PaletteDropEventData,
} from "../types/formDropEvent";
import { createContext, useContext, useState } from "react";

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
                handlePaletteDragEnd(dragData, dropData);
                break;
              default:
                throw new Error(
                  `failed to handle onDragEnd event; unknown FormDragEventData source ${dragData.source}`,
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
                default:
                  throw new Error(
                    `cannot display draggable item; unknown FormDragEventData source ${data.source}`,
                  );
              }
            }}
          </DragOverlay>
        </DragDropProvider>
      </FormDesignerDragContext>
    );
  };
