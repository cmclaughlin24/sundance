import {
  useFormDesignerDispatch,
  useFormDesignerSelect,
} from "@/store/formDesigner";
import type { IElement } from "@/types/element";
import Box from "@mui/material/Box";
import type { MouseEventHandler } from "react";
import { getElementItemStyles } from "./ElementItem.style";
import Typography from "@mui/material/Typography";
import { ItemTools } from "../../../common/ItemTools";
import { Tag } from "@/components/Tag";
import { findPaletteItem } from "../../ToolboxPanel/palette";
import { DraggableCard } from "@/components/DragDrop/DraggableCard";
import { motion, type Variants } from "motion/react";
import type { RemoveElementEvent } from "@/store/formDesigner/events";

const variants: Variants = {
  initial: { opacity: 0, height: 0, marginBottom: 0 },
  animate: { opacity: 1, height: "auto", marginBottom: "1.25rem" },
  exit: { opacity: 0, height: 0, marginBottom: 0 },
};

export const ElementItem: React.FC<{ element: IElement }> = function ({
  element,
}) {
  const { select, isSelected } = useFormDesignerSelect(element.id);
  const { dispatch } = useFormDesignerDispatch();
  const styles = getElementItemStyles(
    isSelected,
    element.attributes.isRequired,
  );

  const handleClk: MouseEventHandler<HTMLLIElement> = (event) => {
    event.stopPropagation();
    select(!isSelected ? { type: element.type, id: element.id } : null);
  };

  const handleDelete = () => {
    dispatch({
      type: "RemoveElement",
      elementId: element.id,
    } satisfies RemoveElementEvent);
    isSelected && select(null);
  };

  const paletteItem = findPaletteItem(element.type);

  return (
    <Box
      component={motion.li}
      layout="position"
      variants={variants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ type: "spring", bounce: 0, duration: 0.3 }}
      sx={styles.item}
      onClick={handleClk}
      role="button"
    >
      <DraggableCard sx={styles.card}>
        <Box sx={styles.content}>
          <Box>
            <Typography sx={styles.name}>{element.name}</Typography>
            <Typography sx={styles.key}>{element.key}</Typography>
          </Box>
          <Box sx={{ display: "flex", alignItems: "center" }}>
            {paletteItem && <Tag>{paletteItem.label}</Tag>}
            <ItemTools
              onReorder={(_inc) => {}}
              onCopy={() => {}}
              onDelete={handleDelete}
            />
          </Box>
        </Box>
      </DraggableCard>
    </Box>
  );
};
