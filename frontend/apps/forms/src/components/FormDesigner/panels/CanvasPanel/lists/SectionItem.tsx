import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
  type SelectedSection,
} from "@/store/formDesigner";
import { useMemo, type MouseEventHandler } from "react";
import { getSectionItemStyles } from "./SectionItem.style";
import { findPaletteItem } from "../../ToolboxPanel/palette";
import { Tag } from "@/components/Tag";
import { ItemTools } from "../../../common/ItemTools";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import { useDroppable } from "@dnd-kit/react";
import type { PaletteDropEventData } from "@/components/FormDesigner/types/formDropEvent";
import { DropZoneIndicator } from "@/components/DragDrop/DropZoneIndicator";
import { useFormDragEvent } from "@/components/FormDesigner/providers/FormDesignerDragProvider";
import {
  CanvasDragType,
  FormDragEventSource,
  PaletteItemDragType,
  type CanvasSectionDragEventData,
  type FormDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import { getNextPosition, sortPositioned } from "@/utils/position";
import { motion, type Variants } from "motion/react";
import type {
  ReorderSectionEvent,
  RemoveSectionEvent,
} from "@/store/formDesigner/events";
import type { ItemComponentProps } from "@/components/FormDesigner/types/componentProps";
import { useSortable } from "@dnd-kit/react/sortable";
import {
  ClipboardEventType,
  type SectionClipboardData,
} from "@/types/clipboard";
import { useContextMenuDispatch } from "@/components/ContextMenu";
import { ElementList } from "./ElementList";

const variants: Variants = {
  initial: { opacity: 0, height: 0, marginBottom: 0 },
  animate: { opacity: 1, height: "auto", marginBottom: "1.25rem" },
  exit: { opacity: 0, height: 0, marginBottom: 0 },
};

export interface SectionItemProps extends ItemComponentProps {
  section: ISection;
}

export const SectionItem: React.FC<SectionItemProps> = function ({
  section,
  parentId,
  index,
  draggable = true,
}) {
  const elements = sortPositioned(section.elements);
  const { dispatch } = useFormDesignerDispatch();
  const { select, isSelected } = useFormDesignerSelect(section.id);
  const { open } = useContextMenuDispatch();

  const { ref: dragRef, handleRef } = useSortable({
    id: `canvas-section-${section.id}`,
    index,
    type: CanvasDragType.Section,
    group: parentId,
    accept: CanvasDragType.Section,
    data: {
      source: FormDragEventSource.Canvas,
      type: "section",
      fromPageId: parentId,
      section,
    } satisfies CanvasSectionDragEventData,
    disabled: !draggable,
  });

  const { ref: dropRef, isDropTarget } = useDroppable({
    id: `section-${section.id}`,
    accept: PaletteItemDragType.Element,
    data: {
      source: "palette",
      parentId: section.id,
      position: getNextPosition(elements),
    } satisfies PaletteDropEventData,
  });

  const dragData = useFormDragEvent();

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: "section", item: section } : null);
  };

  const handleContext: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    event.preventDefault();

    open({
      position: { x: event.clientX, y: event.clientY },
      data: { type: "section", item: section } satisfies SelectedSection,
    });
  };

  const handleCopy = () => {
    const data: SectionClipboardData = {
      type: ClipboardEventType.CopySection,
      section,
    };

    navigator.clipboard.writeText(JSON.stringify(data));
  };

  const handleReorder = (inc: -1 | 1) => {
    dispatch({
      type: "ReorderSection",
      sectionId: section.id,
      inc,
    } satisfies ReorderSectionEvent);
  };

  const handleDelete = () => {
    dispatch({
      type: "RemoveSection",
      sectionId: section.id,
    } satisfies RemoveSectionEvent);
    isSelected && select(null);
  };

  const paletteItem = useMemo(() => findPaletteItem("section"), []);
  const dragPaletteItem = useMemo(
    () =>
      dragData && dragData.source === "palette"
        ? findPaletteItem(dragData.objectType)
        : null,
    [dragData],
  );
  const canDrop = canDropItem(dragData);
  const styles = getSectionItemStyles(isSelected);

  return (
    <Box
      component={motion.li}
      layout="position"
      variants={variants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ type: "spring", bounce: 0, duration: 0.35 }}
      sx={styles.item}
      onClick={handleClk}
      onContextMenu={handleContext}
      ref={(el: Element) => {
        dragRef(el);
        dropRef(el);
      }}
      role="button"
    >
      <DraggableCard
        sx={styles.card}
        handleRef={handleRef}
        orientation="vertical"
        dragHandleOnly
      >
        <Box sx={{ alignSelf: "end", display: "flex", alignItems: "center" }}>
          {paletteItem && <Tag>{paletteItem.label}</Tag>}
          <ItemTools
            onReorder={handleReorder}
            onCopy={handleCopy}
            onDelete={handleDelete}
          />
        </Box>
        <Box sx={styles.elements}>
          <ElementList elements={elements} parentId={section.id} />
          <DropZoneIndicator
            text={`Drop ${dragPaletteItem?.label} here`}
            isVisible={canDrop}
            isDropTarget={isDropTarget}
          />
        </Box>
      </DraggableCard>
    </Box>
  );
};

function canDropItem(data: FormDragEventData | null): boolean {
  if (!data) {
    return false;
  }

  // TODO: Improve this conditional such that it will not introduce a bug if additional layout
  // elements are added.
  return (
    data.source === FormDragEventSource.Palette && data.objectType !== "section"
  );
}
