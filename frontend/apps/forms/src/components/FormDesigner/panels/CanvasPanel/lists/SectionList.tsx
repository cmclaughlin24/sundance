import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { SectionItem } from "./SectionItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";
import type { ListComponentProps } from "@/components/FormDesigner/types/componentProps";
import { BetweenDragZone } from "@/components/DragDrop/BetweenDragZone";
import { useFormDragEvent } from "@/components/FormDesigner/providers/FormDesignerDragProvider";
import {
  CanvasDragType,
  FormDragEventSource,
  type CanvasSectionDragEventData,
  type FormDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import { getBetweenPosition } from "@/utils/position";
import React from "react";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
    marginBottom: -1.5,
    display: "flex",
    flexDirection: "column",
  },
  dragZone: {
    marginBottom: "1.25rem",
  },
};

export interface SectionListProps extends ListComponentProps {
  sections: ISection[];
}

export const SectionList: React.FC<SectionListProps> = function ({
  sections,
  parentId,
}) {
  const dragData = useFormDragEvent();
  const isSectionDrag = isCanvasSectionDrag(dragData);

  const createDragZone = (index: number) => (
    <BetweenDragZone
      key={`zone-${parentId}-${index}`}
      id={`zone-${parentId}-${index}`}
      accept={CanvasDragType.Section}
      parentId={parentId}
      position={getBetweenPosition(index, sections)}
      text="Drop Here"
      sx={styles.dragZone}
    />
  );

  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {isSectionDrag && createDragZone(0)}
        {sections.map((section, index) => {
          return (
            <React.Fragment key={section.id}>
              <SectionItem
                section={section}
                parentId={parentId}
                key={`${parentId}-${section.id}`}
              />
              {isSectionDrag &&
                section.id !==
                  (dragData as CanvasSectionDragEventData).section.id &&
                createDragZone(index + 1)}
            </React.Fragment>
          );
        })}
      </AnimatePresence>
    </Box>
  );
};

function isCanvasSectionDrag(dragData: FormDragEventData | null): boolean {
  return (
    dragData?.source === FormDragEventSource.Canvas &&
    dragData.type === "section"
  );
}
