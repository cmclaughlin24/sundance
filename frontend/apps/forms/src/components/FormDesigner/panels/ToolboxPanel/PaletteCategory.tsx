import type { Styles } from "@/types/styles";
import type { IPaletteCategory } from "./palette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { PaletteItem } from "./PaletteItem";

const styles: Styles = {
  category: {
    display: "flex",
    flexDirection: "column",
    gap: 1,
  },
  label: {
    fontWeight: 400,
  },
  list: {
    margin: 0,
    padding: 0,
    display: "flex",
    flexDirection: "column",
    gap: 0.5,
  },
};

export const PaletteCategory: React.FC<{ category: IPaletteCategory }> =
  function ({ category }) {
    return (
      <Box component="section" sx={styles.category}>
        <Typography component="h4" sx={styles.label}>
          {category.label}
        </Typography>
        <Box component="ul" sx={styles.list}>
          {category.items.map((item) => (
            <PaletteItem key={item.type} item={item} />
          ))}
        </Box>
      </Box>
    );
  };
