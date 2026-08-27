import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";
import type { IPaletteItem } from "./palette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  item: {
    listStyle: "none",
    display: "flex",
    alignItems: "center",
    gap: 1,
    px: 1,
    py: 1.5,
    border: `1px solid ${Border.Primary}`,
    borderRadius: 2,
  },
};

export const PaletteItem: React.FC<{ item: IPaletteItem }> = function ({
  item,
}) {
  return (
    <Box component="li" sx={styles.item}>
      {item.icon}
      <Typography>{item.label}</Typography>
    </Box>
  );
};
