import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import {
  useFormDesignerDispatch,
  useSelectedItem,
  type SelectedItem,
} from "@/store/formDesigner";
import type { MouseEventHandler } from "react";
import { getSectionItemStyles } from "./SectionItem.style";

export const SectionItem: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const item = useSelectedItem();
  const { select } = useFormDesignerDispatch();
  const styles = getSectionItemStyles(!!item && item.id === section.id);

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();

    let selection: SelectedItem | null = null;

    if (!item || item.id !== section.id) {
      selection = { type: "section", id: section.id };
    }

    select(selection);
  };

  return (
    <Box component="li" sx={styles.item} onClick={handleClk} role="button">
      <ElementList elements={section.elements} />
    </Box>
  );
};
