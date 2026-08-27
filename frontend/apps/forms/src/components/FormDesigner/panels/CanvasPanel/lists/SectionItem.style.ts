import { Background, Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const getSectionItemStyles = (isSelected: boolean): Styles => {
  return {
    item: {
      listStyle: "none",
    },
    card: (theme) => ({
      padding: 2.5,
      gap: 2.5,
      borderRadius: 2.5,
      borderColor: isSelected
        ? `${theme.palette.primary.main}`
        : Border.Primary,
      borderStyle: isSelected ? "solid" : "dashed",
      background: isSelected
        ? `${theme.palette.primary.main}25`
        : Background.Primary,
      overflow: "hidden",
      ":hover:not(:has(li:hover, button:hover))": {
        cursor: "pointer",
        border: `1px solid ${theme.palette.primary.main}`,
      },
    }),
  };
};
