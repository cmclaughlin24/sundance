import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
} from "@/store/formDesigner";
import { useMemo, type MouseEventHandler } from "react";
import { getSectionItemStyles } from "./SectionItem.style";
import {
  findPaletteItem,
  PaletteItemDragType,
} from "../../ToolboxPanel/palette";
import { Tag } from "@/components/Tag";
import { ItemTools } from "../../../common/ItemTools";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import { useDroppable } from "@dnd-kit/react";
import type { PaletteDropEventData } from "@/components/FormDesigner/types/formDropEvent";
import { DropZoneIndicator } from "@/components/DragDrop/DropZoneIndicator";
import { useFormDragEvent } from "@/components/FormDesigner/providers/FormDesignerDragProvider";
import {
  FormDragEventSource,
  type FormDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import { getNextPosition, sortPositioned } from "@/utils/position";
import { motion, type Variants } from "motion/react";
import type { RemoveSectionEvent } from "@/store/formDesigner/events";

const variants: Variants = {
  initial: { opacity: 0, height: 0, marginBottom: 0 },
  animate: { opacity: 1, height: "auto", marginBottom: "1.25rem" },
  exit: { opacity: 0, height: 0, marginBottom: 0 },
};

export const SectionItem: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const elements = sortPositioned(section.elements);
  const { dispatch } = useFormDesignerDispatch();
  const { select, isSelected } = useFormDesignerSelect(section.id);
  const { ref, isDropTarget } = useDroppable({
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
    select(!isSelected ? { type: "section", id: section.id } : null);
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
    () => (dragData ? findPaletteItem(dragData.type) : null),
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
      ref={ref}
      role="button"
    >
      <DraggableCard sx={styles.card} orientation="vertical">
        <Box sx={{ alignSelf: "end", display: "flex", alignItems: "center" }}>
          {paletteItem && <Tag>{paletteItem.label}</Tag>}
          <ItemTools
            onReorder={(_inc) => {}}
            onCopy={() => {}}
            onDelete={handleDelete}
          />
        </Box>
        <Box sx={styles.elements}>
          <ElementList elements={elements} />
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
    data.source === FormDragEventSource.Palette &&
    data.type !== PaletteItemDragType.Section
  );
}
