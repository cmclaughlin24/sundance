import type { IElement } from "@/types/element";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";

const styles: Styles = {
  item: {
    listStyle: "none",
  },
};

export const ElementItem: React.FC<{ element: IElement }> = function ({
  element,
}) {
  return (
    <Box component="li" sx={styles.item}>
      {element.name}
    </Box>
  );
};
