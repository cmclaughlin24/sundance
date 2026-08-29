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
  dragHandleOnly?: boolean;
}>;

export const DraggableCard: React.FC<DraggableCardProps> = function ({
  sx = {},
  orientation = "horizontal",
  handleRef,
  dragHandleOnly = false,
  children,
}) {
  const { draggableCard, icon } = draggableCardStyles(orientation);

  return (
    <Box ref={!dragHandleOnly ? handleRef : null} sx={mergeSx(draggableCard, sx)}>
      <Box
        ref={dragHandleOnly ? handleRef : null}
        sx={{
          display: "flex",
          justifyContent: "center",
          alignContent: "center",
        }}
      >
        <DragIndicator sx={icon} />
      </Box>
      {children}
    </Box>
  );
};
