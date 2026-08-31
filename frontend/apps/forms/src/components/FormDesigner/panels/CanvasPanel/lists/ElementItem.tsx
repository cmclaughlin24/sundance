import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
} from "@/store/formDesigner";
import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import { type MouseEventHandler } from "react";
import { getElementItemStyles } from "./ElementItem.style";
import Typography from "@mui/material/Typography";
import { ItemTools } from "../../../common/ItemTools";
import { Tag } from "@/components/Tag";
import { findPaletteItem } from "../../ToolboxPanel/palette";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import type {
  ReorderElementEvent,
  RemoveElementEvent,
} from "@/store/formDesigner/events";
import {
  CanvasDragType,
  FormDragEventSource,
  type CanvasElementDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import type { ItemComponentProps } from "@/components/FormDesigner/types/componentProps";
import { useSortable } from "@dnd-kit/react/sortable";
import { motion, type Variants } from "motion/react";
import {
  ClipboardEventType,
  type CopyElementClipboardData,
} from "@/types/clipboard";

export interface ElementItemProps extends ItemComponentProps {
  element: IElement;
}

const variants: Variants = {
  initial: { opacity: 0, height: 0, marginBottom: 0 },
  animate: { opacity: 1, height: "auto", marginBottom: "1.25rem" },
  exit: { opacity: 0, height: 0, marginBottom: 0 },
};

export const ElementItem: React.FC<ElementItemProps> = function ({
  element,
  parentId,
  index,
  draggable = true,
}) {
  const { select, isSelected } = useFormDesignerSelect(element.id);
  const { dispatch } = useFormDesignerDispatch();

  const { ref, handleRef } = useSortable({
    id: `canvas-element-${element.id}`,
    index,
    type: CanvasDragType.Element,
    group: parentId,
    accept: CanvasDragType.Element,
    data: {
      source: FormDragEventSource.Canvas,
      type: "element",
      fromSectionId: parentId,
      element,
    } satisfies CanvasElementDragEventData,
    disabled: !draggable,
  });

  const styles = getElementItemStyles(
    isSelected,
    element.attributes.isRequired,
  );

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: "element", item: element } : null);
  };

  const handleCopy = () => {
    const data: CopyElementClipboardData = {
      type: ClipboardEventType.CopyElement,
      element,
    };

    navigator.clipboard.writeText(JSON.stringify(data));
  };

  const handleReorder = (inc: -1 | 1) => {
    dispatch({
      type: "ReorderElement",
      elementId: element.id,
      inc,
    } satisfies ReorderElementEvent);
  };

  const handleDelete = () => {
    dispatch({
      type: "RemoveElement",
      elementId: element.id,
    } satisfies RemoveElementEvent);
    isSelected && select(null);
  };

  const paletteItem = findPaletteItem(element.type);

  return (
    <>
      <Box
        component={motion.li}
        layout="position"
        variants={variants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ type: "spring", bounce: 0, duration: 0.3 }}
        sx={styles.item}
        onClick={handleClk}
        role="button"
        ref={ref}
      >
        <DraggableCard sx={styles.card} handleRef={handleRef} dragHandleOnly>
          <Box sx={styles.content}>
            <Box>
              <Typography sx={styles.name}>{element.name}</Typography>
              <Typography sx={styles.key}>{element.key}</Typography>
            </Box>
            <Box sx={{ display: "flex", alignItems: "center" }}>
              {paletteItem && <Tag>{paletteItem.label}</Tag>}
              <ItemTools
                onReorder={handleReorder}
                onCopy={handleCopy}
                onDelete={handleDelete}
              />
            </Box>
          </Box>
        </DraggableCard>
      </Box>
    </>
  );
};
