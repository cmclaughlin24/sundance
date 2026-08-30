import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import { ElementItem } from "./ElementItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";
import type { ListComponentProps } from "@/components/FormDesigner/types/componentProps";
import React from "react";

const styles: Styles = {
  list: {
    margin: 0,
    marginBottom: -1.5,
    padding: 0,
    display: "flex",
    flexDirection: "column",
  },
  dragZone: {
    marginBottom: "1.25rem",
  },
};

export interface ElementListProps extends ListComponentProps {
  elements: IElement[];
}

export const ElementList: React.FC<ElementListProps> = function ({
  elements,
  parentId,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {elements.map((element, index) => (
          <ElementItem
            element={element}
            parentId={parentId}
            index={index}
            key={element.id}
          />
        ))}
      </AnimatePresence>
    </Box>
  );
};
