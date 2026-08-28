import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { SectionItem } from "./SectionItem";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
};

export const SectionList: React.FC<{ sections: ISection[] }> = function ({
  sections,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      {sections.map((section) => (
        <SectionItem section={section} key={section.id} />
      ))}
    </Box>
  );
};
