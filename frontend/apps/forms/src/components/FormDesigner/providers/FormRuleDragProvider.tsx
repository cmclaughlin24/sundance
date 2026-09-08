import { DragDropProvider, DragOverlay } from "@dnd-kit/react";
import { createContext, useContext, useState } from "react";
import type { PaletteDragEventData } from "../types/formDragEvent";
import type { RuleType } from "@/types/rule";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import { generatedID } from "@/utils/id";
import { findFormRulePaletteItem } from "../panels/ToolboxPanel/palette";
import { PaletteItem } from "../panels/ToolboxPanel/PaletteItem";

const FormRuleDragContext =
  createContext<PaletteDragEventData<RuleType> | null>(null);

export function useFormRuleDragData() {
  return useContext(FormRuleDragContext);
}

export const FormRuleDragProvider: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    const [activeDragData, setActiveDragData] =
      useState<PaletteDragEventData<RuleType> | null>(null);
    const { dispatch } = useFormDesignerDispatch();

    const handleDragEnd = (dragData: PaletteDragEventData<RuleType>) => {
      dispatch({
        type: "AddRule",
        id: generatedID(),
        ruleType: dragData.itemType,
      });
    };

    return (
      <FormRuleDragContext value={activeDragData}>
        <DragDropProvider
          onDragStart={(event) => {
            setActiveDragData(
              event.operation.source?.data as PaletteDragEventData<RuleType>,
            );
          }}
          onDragEnd={(event) => {
            setActiveDragData(null);

            const { source, target } = event.operation;

            if (event.canceled || !target) {
              return;
            }

            const dragData = source?.data as PaletteDragEventData<RuleType>;

            handleDragEnd(dragData);
          }}
        >
          {children}
          <DragOverlay>
            {(source) => {
              const data = source.data as PaletteDragEventData<RuleType>;
              const paletteItem = findFormRulePaletteItem(data.itemType);

              if (!paletteItem) {
                return null;
              }

              return <PaletteItem item={paletteItem} draggable={false} />;
            }}
          </DragOverlay>
        </DragDropProvider>
      </FormRuleDragContext>
    );
  };
