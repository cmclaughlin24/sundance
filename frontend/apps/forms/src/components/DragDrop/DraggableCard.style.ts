import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const draggableCardStyles = (
  orientation: "horizontal" | "vertical",
): Styles => {
  const isHorizontal = orientation === "horizontal";

  return {
    draggableCard: {
      display: "flex",
      flexDirection: isHorizontal ? "row" : "column",
      border: `1px solid ${Border.Primary}`,
      padding: 1.5,
    },
    icon: {
      transform: isHorizontal ? "" : "rotate(-90deg)",
      alignSelf: isHorizontal ? "" : "center"
    },
    content: {
      flex: 1,
    },
  };
};
