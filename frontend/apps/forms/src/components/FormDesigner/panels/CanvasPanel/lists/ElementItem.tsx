import {
  useFormDesignerDispatch,
  useSelectedItem,
  type SelectedItem,
} from "@/store/formDesigner";
import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import type { MouseEventHandler } from "react";
import { getElementItemStyles } from "./ElementItem.style";

export const ElementItem: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const item = useSelectedItem();
  const { select } = useFormDesignerDispatch();
  const styles = getElementItemStyles(!!item && item.id === element.id);

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();

    let selection: SelectedItem | null = null;

    if (!item || item.id !== element.id) {
      selection = { type: "element", id: element.id };
    }

    select(selection);
  };

  return (
    <Box component="li" sx={styles.item} onClick={handleClk} role="button">
      {element.name}
    </Box>
  );
};
