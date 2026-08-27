import type { Styles } from "@/types/styles";
import type { IPaletteItem } from "./palette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { DraggableCard } from "@/components/DraggableCard";

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

export const PaletteItem: React.FC<{ item: IPaletteItem }> = function ({
  item,
}) {
  return (
    <Box component="li" sx={styles.item}>
      <DraggableCard sx={styles.card}>
        {item.icon}
        <Typography>{item.label}</Typography>
      </DraggableCard>
    </Box>
  );
};
