import { DragDropProvider } from "@dnd-kit/react";
import { createContext, useContext, useState } from "react";
import type { PaletteDragEventData } from "../types/formDragEvent";
import type { RuleType } from "@/types/rule";

const FormRuleDragContext =
  createContext<PaletteDragEventData<RuleType> | null>(null);

export function useFormRuleDragData() {
  return useContext(FormRuleDragContext);
}

export type FormRuleDragProviderProps = React.PropsWithChildren<{
  onPaletteDrop: (itemType: RuleType) => void;
}>;

export const FormRuleDragProvider: React.FC<FormRuleDragProviderProps> =
  function ({ children, onPaletteDrop }) {
    const [activeDragData, setActiveDragData] =
      useState<PaletteDragEventData<RuleType> | null>(null);

    const handleDragEnd = (
      dragData: PaletteDragEventData<RuleType>,
      _dropData: any,
    ) => {
      onPaletteDrop(dragData.itemType);
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

            if (event.canceled) {
              return;
            }

            const { source, target } = event.operation;
            const dragData = source?.data as PaletteDragEventData<RuleType>;
            const dropData = target?.data as any;

            handleDragEnd(dragData, dropData);
          }}
        >
          {children}
        </DragDropProvider>
      </FormRuleDragContext>
    );
  };
