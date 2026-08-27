import { DragDropProvider, DragOverlay } from "@dnd-kit/react";
import { findPaletteItem } from "../panels/ToolboxPanel/palette";
import { PaletteItem } from "../panels/ToolboxPanel/PaletteItem";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import type { FormDragEventData } from "../types/drag-event";

export const FormDesignerDragProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const {} = useFormDesignerDispatch();

    return (
      <DragDropProvider
        onDragEnd={(event) => {
          const { source, target } = event.operation;
          const data = source?.data as FormDragEventData;
        }}
      >
        {children}
        <DragOverlay>
          {(source) => {
            const data = source.data as FormDragEventData;

            if (data?.source === "palette") {
              const paletteItem = findPaletteItem(source.data.type);
              return <PaletteItem item={paletteItem!} draggable={false} />;
            }
          }}
        </DragOverlay>
      </DragDropProvider>
    );
  };
