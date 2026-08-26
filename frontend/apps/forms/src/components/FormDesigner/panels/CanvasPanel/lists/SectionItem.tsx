import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  item: {
    listStyle: "none",
  },
};

export const SectionItem: React.FC<{ section: ISection }> = function ({
  section,
}) {
  return (
    <Box component="li" sx={styles.item}>
      <ElementList elements={section.elements} />
    </Box>
  );
};
