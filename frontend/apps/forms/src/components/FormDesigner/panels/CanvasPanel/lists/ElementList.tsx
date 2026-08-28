import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import { ElementItem } from "./ElementItem";
import type { Styles } from "@/types/styles";
import { AnimatePresence } from "motion/react";

const styles: Styles = {
  list: {
    margin: 0,
    marginBottom: -1.5,
    padding: 0,
    display: "flex",
    flexDirection: "column",
  },
};

export const ElementList: React.FC<{ elements: IElement[] }> = function ({
  elements,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      <AnimatePresence initial={false}>
        {elements.map((element) => (
          <ElementItem element={element} key={element.id} />
        ))}
      </AnimatePresence>
    </Box>
  );
};
