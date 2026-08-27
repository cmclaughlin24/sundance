import Box from "@mui/material/Box";
import { findPaletteItem } from "../ToolboxPanel/palette";
import Typography from "@mui/material/Typography";
import type { ElementType } from "@/types/element";
import type { Styles } from "@/types/styles";
import { Border } from "@/constants/colors";

const styles: Styles = {
  activeObject: {
    display: "flex",
    alignItems: "center",
    gap: 1,
    py: 2.5,
    borderTop: `1px solid ${Border.Primary}`,
    borderBottom: `1px solid ${Border.Primary}`,
  },
  label: {
    fontWeight: 600,
  },
};

export const ActiveObjectTitle: React.FC<{
  elementType: ElementType | "section" | "page";
}> = function ({ elementType }) {
  const paletteItem = findPaletteItem(elementType);

  return (
    <Box sx={styles.activeObject}>
      {paletteItem!.icon}
      <Typography component="h4" sx={styles.title}>
        {paletteItem!.label}
      </Typography>
    </Box>
  );
};
