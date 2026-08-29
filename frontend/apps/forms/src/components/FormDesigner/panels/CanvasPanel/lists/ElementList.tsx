import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import { ElementItem } from "./ElementItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";
import type { ListComponentProps } from "@/components/FormDesigner/types/componentProps";
import {
  CanvasDragType,
  FormDragEventSource,
  type CanvasElementDragEventData,
  type FormDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import { useFormDragEvent } from "@/components/FormDesigner/providers/FormDesignerDragProvider";
import { BetweenDragZone } from "@/components/DragDrop/BetweenDragZone";
import { getBetweenPosition } from "@/utils/position";
import React from "react";

const styles: Styles = {
  list: {
    margin: 0,
    marginBottom: -1.5,
    padding: 0,
    display: "flex",
    flexDirection: "column",
  },
  dragZone: {
    marginBottom: "1.25rem",
  },
};

export interface ElementListProps extends ListComponentProps {
  elements: IElement[];
}

export const ElementList: React.FC<ElementListProps> = function ({
  elements,
  parentId,
}) {
  const dragData = useFormDragEvent();
  const isElementDrag = isCanvasElementDrag(dragData);

  const createDragZone = (index: number) => (
    <BetweenDragZone
      key={`zone-${parentId}-${index}`}
      id={`zone-${parentId}-${index}`}
      accept={CanvasDragType.Element}
      parentId={parentId}
      position={getBetweenPosition(index, elements)}
      text="Drop Here"
      sx={styles.dragZone}
    />
  );

  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {isElementDrag && createDragZone(0)}
        {elements.map((element, index) => {
          return (
            <React.Fragment key={element.id}>
              <ElementItem
                element={element}
                parentId={parentId}
                key={`${parentId}-${element.id}`}
              />
              {isElementDrag &&
                (dragData as CanvasElementDragEventData).element.id !==
                  element.id &&
                createDragZone(index + 1)}
            </React.Fragment>
          );
        })}
      </AnimatePresence>
    </Box>
  );
};

function isCanvasElementDrag(dragData: FormDragEventData | null): boolean {
  return (
    dragData?.source === FormDragEventSource.Canvas &&
    dragData.type === "element"
  );
}
