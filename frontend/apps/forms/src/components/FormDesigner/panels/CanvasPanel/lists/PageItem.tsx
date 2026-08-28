import type { IPage } from "@/types/page";
import Box from "@mui/material/Box";
import { SectionList } from "./SectionList";
import type { Styles } from "@/types/styles";
import { useDroppable } from "@dnd-kit/react";
import { PaletteItemDragType } from "../../ToolboxPanel/palette";
import type { PaletteDropEventData } from "@/components/FormDesigner/types/formDropEvent";
import { DropZoneIndicator } from "@/components/DragDrop/DropZoneIndicator";
import {
  FormDragEventSource,
  type FormDragEventData,
} from "@/components/FormDesigner/types/formDragEvent";
import { useFormDragEvent } from "@/components/FormDesigner/providers/FormDesignerDragProvider";

const styles: Styles = {
  item: {
    listStyle: "none",
  },
};

export const PageItem: React.FC<{ page: IPage }> = function ({ page }) {
  const { ref } = useDroppable({
    id: `page-${page.id}`,
    accept: PaletteItemDragType.Section,
    data: { source: "palette", parentId: page.id } as PaletteDropEventData,
  });
  const dragData = useFormDragEvent();
  const canDrop = canDropItem(dragData);

  return (
    <Box component="li" sx={styles.item} ref={ref}>
      <SectionList sections={page.sections} />
      {canDrop && <DropZoneIndicator text="Drop you section here" />}
    </Box>
  );
};

function canDropItem(data: FormDragEventData | null): boolean {
  if (!data) {
    return false;
  }

  // TODO: Improve this conditional such that it will not introduce a bug if additional layout
  // elements are added.
  return data.source === FormDragEventSource.Palette && data.type === "section";
}
