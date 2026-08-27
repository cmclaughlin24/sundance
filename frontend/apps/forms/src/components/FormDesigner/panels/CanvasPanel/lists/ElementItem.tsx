import { useFormDesignerSelect } from "@/store/formDesigner";
import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import type { MouseEventHandler } from "react";
import { getElementItemStyles } from "./ElementItem.style";
import Typography from "@mui/material/Typography";
import { ItemTools } from "../../../common/ItemTools";
import { Tag } from "@/components/Tag";
import { findPaletteItem } from "../../ToolboxPanel/palette";

export const ElementItem: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const { select, isSelected } = useFormDesignerSelect(element.id);
  const styles = getElementItemStyles(
    isSelected,
    element.attributes.isRequired,
  );
  const paletteItem = findPaletteItem(element.type);

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: element.type, id: element.id } : null);
  };

  return (
    <Box component="li" sx={styles.item} onClick={handleClk} role="button">
      <Box>
        <Typography sx={styles.name}>{element.name}</Typography>
        <Typography sx={styles.key}>{element.key}</Typography>
      </Box>
      <Box sx={{ display: "flex", alignItems: "center" }}>
        {paletteItem && <Tag>{paletteItem.label}</Tag>}
        <ItemTools
          onReorder={(_inc) => {}}
          onCopy={() => {}}
          onDelete={() => {}}
        />
      </Box>
    </Box>
  );
};
