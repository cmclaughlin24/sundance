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

export interface PaletteCategoryProps<T> {
  category: IPaletteCategory<T>;
}

export function PaletteCategory<T>({ category }: PaletteCategoryProps<T>) {
  return (
    <Box component="section" sx={styles.category}>
      <Typography component="h4" sx={styles.label}>
        {category.label}
      </Typography>
      <Box component="ul" sx={styles.list}>
        {category.items.map((item) => (
          <PaletteItem<T> key={item.label} item={item} />
        ))}
      </Box>
    </Box>
  );
}
