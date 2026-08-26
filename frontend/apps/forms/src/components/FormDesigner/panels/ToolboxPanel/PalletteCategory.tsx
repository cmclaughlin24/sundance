import type { Styles } from "@/types/styles";
import type { IPalletteCategory } from "./pallette";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { PalletteItem } from "./PalletteItem";

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

export const PalletteCategory: React.FC<{ category: IPalletteCategory }> =
  function ({ category }) {
    return (
      <Box component="section" sx={styles.category}>
        <Typography component="h4" sx={styles.label}>
          {category.label}
        </Typography>
        <Box component="ul" sx={styles.list}>
          {category.items.map((item) => (
            <PalletteItem key={item.type} item={item} />
          ))}
        </Box>
      </Box>
    );
  };
