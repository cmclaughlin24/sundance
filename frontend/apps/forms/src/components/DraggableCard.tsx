import Box from "@mui/material/Box";
import { draggableCardStyles } from "./DraggableCard.style";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";
import DragIndicator from "@mui/icons-material/DragIndicator";
import type { Ref } from "react";

export type DraggableCardProps = React.PropsWithChildren<{
  sx?: SxProps<Theme>;
  orientation?: "horizontal" | "vertical";
  handleRef?: Ref<Element>;
}>;

export const DraggableCard: React.FC<DraggableCardProps> = function ({
  sx = {},
  orientation = "horizontal",
  handleRef,
  children,
}) {
  const { draggableCard, icon } = draggableCardStyles(orientation);

  return (
    <Box ref={handleRef} sx={mergeSx(draggableCard, sx)}>
      <DragIndicator sx={icon} />
      {children}
    </Box>
  );
};
