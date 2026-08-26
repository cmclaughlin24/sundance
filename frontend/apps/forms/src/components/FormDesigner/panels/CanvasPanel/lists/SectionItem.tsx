import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import type { Styles } from "@/types/styles";
import { Border } from "@/constants/colors";

const styles: Styles = {
  item: {
    listStyle: "none",
    border: `1px dashed ${Border.Primary}`,
    borderRadius: 2.5,
    padding: 2.5,
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
