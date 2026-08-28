import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import { useFormDesignerSelect } from "@/store/formDesigner";
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

export const SectionItem: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const { select, isSelected } = useFormDesignerSelect(section.id);
  const { ref, isDropTarget } = useDroppable({
    id: `section-${section.id}`,
    accept: PaletteItemDragType.Element,
    data: { source: "palette", parentId: section.id } as PaletteDropEventData,
  });
  const dragData = useFormDragEvent();

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: "section", id: section.id } : null);
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
      component="li"
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
            onDelete={() => {}}
          />
        </Box>
        <ElementList elements={section.elements} />
        {canDrop && (
          <DropZoneIndicator
            text={`Drop ${dragPaletteItem!.label} here`}
            isDropTarget={isDropTarget}
          />
        )}
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
  return data.source === FormDragEventSource.Palette && data.type !== "section";
}
