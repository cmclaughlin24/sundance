import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { SectionItem } from "./SectionItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
    marginBottom: -1.5,
    display: "flex",
    flexDirection: "column",
  },
};

export const SectionList: React.FC<{ sections: ISection[] }> = function ({
  sections,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {sections.map((section) => (
          <SectionItem section={section} key={section.id} />
        ))}
      </AnimatePresence>
    </Box>
  );
};
