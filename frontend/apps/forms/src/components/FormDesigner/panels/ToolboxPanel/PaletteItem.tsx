import type { Styles } from "@/types/styles";
import type { IPaletteItem } from "./palette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import { useDraggable } from "@dnd-kit/react";
import {
  FormDragEventSource,
  type PaletteDragEventData,
} from "../../types/formDragEvent";
import { mergeSx } from "merge-sx";

const styles: Styles = {
  item: {
    listStyle: "none",
  },
  card: (theme) => ({
    display: "flex",
    alignItems: "center",
    gap: 1,
    px: 1,
    py: 1.5,
    borderRadius: 2,
    ":hover": {
      cursor: "pointer",
      border: `1px solid ${theme.palette.primary.main}`,
    },
  }),
  onDrag: (theme) => ({
    background: `${theme.palette.primary.main}70`,
    border: `1px solid ${theme.palette.primary.main}`,
  }),
};

export interface PaletteItemProps<IType, DType> {
  item: IPaletteItem<IType, DType>;
  draggable?: boolean;
}

export function PaletteItem<IType, DType>({
  item,
  draggable = true,
}: PaletteItemProps<IType, DType>) {
  const { ref, handleRef, isDragging } = useDraggable({
    id: `palette-${item.type}`,
    type: item.dragType as any,
    data: {
      source: FormDragEventSource.Palette,
      itemType: item.type as any,
    } satisfies PaletteDragEventData<IType>,
    disabled: !draggable,
  });

  return (
    <Box component="li" sx={styles.item} ref={ref}>
      <DraggableCard
        sx={mergeSx(styles.card, isDragging && !draggable ? styles.onDrag : {})}
        handleRef={draggable ? handleRef : null}
      >
        {item.icon}
        <Typography>{item.label}</Typography>
      </DraggableCard>
    </Box>
  );
}
