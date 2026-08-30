import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { SectionItem } from "./SectionItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";
import type { ListComponentProps } from "@/components/FormDesigner/types/componentProps";
import React from "react";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
    marginBottom: -1.5,
    display: "flex",
    flexDirection: "column",
  },
  dragZone: {
    marginBottom: "1.25rem",
  },
};

export interface SectionListProps extends ListComponentProps {
  sections: ISection[];
}

export const SectionList: React.FC<SectionListProps> = function ({
  sections,
  parentId,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {sections.map((section, index) => (
          <SectionItem
            section={section}
            parentId={parentId}
            index={index}
            key={`${parentId}-${section.id}`}
          />
        ))}
      </AnimatePresence>
    </Box>
  );
};
