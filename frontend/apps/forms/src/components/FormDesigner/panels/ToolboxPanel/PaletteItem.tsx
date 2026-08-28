import type { Styles } from "@/types/styles";
import type { IPaletteItem } from "./palette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import { useDraggable } from "@dnd-kit/react";
import { FormDragEventSource, type FormDragEventData } from "../../types/formDragEvent";

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
};

export const PaletteItem: React.FC<{
  item: IPaletteItem;
  draggable?: boolean;
}> = function ({ item, draggable = true }) {
  const { ref, handleRef } = useDraggable({
    id: `palette-${item.type}`,
    type: item.dragType,
    data: {
      source: FormDragEventSource.Palette,
      type: item.type,
    } satisfies FormDragEventData,
    disabled: !draggable,
  });

  return (
    <Box component="li" sx={styles.item} ref={ref}>
      <DraggableCard sx={styles.card} handleRef={draggable ? handleRef : null}>
        {item.icon}
        <Typography>{item.label}</Typography>
      </DraggableCard>
    </Box>
  );
};
