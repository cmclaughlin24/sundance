import type { Styles } from "@/types/styles";
import { useDroppable } from "@dnd-kit/react";
import Box from "@mui/material/Box";
import type { SxProps, Theme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import { mergeSx } from "merge-sx";
import type { CanvasDropEventData } from "../FormDesigner/types/formDropEvent";

const styles: Styles = {
  dragZone: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    gap: 1.5,
    py: 1.5,
    transition: "background 0.25s ease",
  },
  hr: (theme) => ({
    border: 0,
    borderTop: `1px dashed ${theme.palette.primary.main}`,
    flex: 1,
  }),
  text: (theme) => ({
    fontWeight: 600,
    color: theme.palette.primary.main,
  }),
};

export interface BetweenDragZone {
  id: string;
  text: string;
  accept?: string | string[];
  parentId: string;
  position: number;
  sx?: SxProps<Theme>;
}

export const BetweenDragZone: React.FC<BetweenDragZone> = function ({
  id,
  accept,
  text,
  parentId,
  position,
  sx,
}) {
  const { ref } = useDroppable({
    id,
    accept,
    data: {
      source: "canvas",
      parentId,
      position,
    } satisfies CanvasDropEventData,
  });

  return (
    <Box ref={ref} sx={mergeSx(styles.dragZone, sx)}>
      <Box component="hr" sx={styles.hr} />
      <Typography variant="body2" sx={styles.text}>
        {text}
      </Typography>
      <Box component="hr" sx={styles.hr} />
    </Box>
  );
};
