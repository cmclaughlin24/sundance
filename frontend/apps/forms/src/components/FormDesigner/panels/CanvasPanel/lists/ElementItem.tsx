import { Border } from "@/constants/colors";
import type { IElement } from "@/types/element";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";

const styles: Styles = {
  item: {
    listStyle: "none",
    padding: 1.5,
    border: `1px dashed ${Border.Primary}`,
    borderRadius: 2.5,
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
