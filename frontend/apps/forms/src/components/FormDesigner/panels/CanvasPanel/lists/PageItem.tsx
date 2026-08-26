import type { IPage } from "@/types/page";
import Box from "@mui/material/Box";
import { SectionList } from "./SectionList";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  item: {
    listStyle: "none",
  },
};

export const PageItem: React.FC<{ page: IPage }> = function ({ page }) {
  return (
    <Box component="li" sx={styles.item}>
      <SectionList sections={page.sections} />
    </Box>
  );
};
