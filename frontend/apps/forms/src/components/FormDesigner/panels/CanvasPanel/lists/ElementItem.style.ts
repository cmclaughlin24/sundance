import { Background, Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const getElementItemStyles = (isSelected: boolean): Styles => {
  return {
    item: (theme) => ({
      listStyle: "none",
      padding: 1.5,
      border: `1px dashed ${Border.Primary}`,
      borderRadius: 2.5,
      borderColor: isSelected
        ? `${theme.palette.primary.main}`
        : Border.Primary,
      borderStyle: isSelected ? "solid" : "dashed",
      background: isSelected
        ? `${theme.palette.primary.main}25`
        : Background.Primary,
      overflow: "hidden",
      ":hover:not(:has(li:hover))": {
        cursor: "pointer",
        border: `1px solid ${theme.palette.primary.main}`,
        background: `${theme.palette.primary.main}25`,
      },
    }),
  };
};
