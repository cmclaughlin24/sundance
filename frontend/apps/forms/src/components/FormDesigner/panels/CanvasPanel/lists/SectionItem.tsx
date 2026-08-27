import type { ISection } from "@/types/section";
import Box from "@mui/material/Box";
import { ElementList } from "./ElementList";
import { useFormDesignerSelect } from "@/store/formDesigner";
import type { MouseEventHandler } from "react";
import { getSectionItemStyles } from "./SectionItem.style";
import { findPaletteItem } from "../../ToolboxPanel/palette";
import { Tag } from "@/components/Tag";
import { ItemTools } from "../../../common/ItemTools";

export const SectionItem: React.FC<{ section: ISection }> = function ({
  section,
}) {
  const { select, isSelected } = useFormDesignerSelect(section.id);
  const styles = getSectionItemStyles(isSelected);
  const paletteItem = findPaletteItem("section");

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: "section", id: section.id } : null);
  };

  return (
    <Box component="li" sx={styles.item} onClick={handleClk} role="button">
      <Box sx={{ alignSelf: "end", display: "flex", alignItems: "center" }}>
        {paletteItem && <Tag>{paletteItem.label}</Tag>}
        <ItemTools
          onReorder={(_inc) => {}}
          onCopy={() => {}}
          onDelete={() => {}}
        />
      </Box>
      <ElementList elements={section.elements} />
    </Box>
  );
};
