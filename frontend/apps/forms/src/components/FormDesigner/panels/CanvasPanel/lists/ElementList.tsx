import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import { ElementItem } from "./ElementItem";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
    display: "flex",
    flexDirection: "column",
    gap: 1.5,
  },
};

export const ElementList: React.FC<{ elements: IElement[] }> = function ({
  elements,
}) {
  return (
    <Box component="ul" sx={styles.list}>
      {elements.map((element) => (
        <ElementItem element={element} key={element.id} />
      ))}
    </Box>
  );
};
