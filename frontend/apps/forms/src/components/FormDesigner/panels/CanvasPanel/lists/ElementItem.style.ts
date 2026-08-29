import { Background, Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const getElementItemStyles = (
  isSelected: boolean,
  isRequired: boolean,
): Styles => {
  return {
    item: {
      listStyle: "none",
    },
    card: (theme) => ({
      padding: 1.5,
      alignItems: "center",
      gap: 1.5,
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
    content: {
      flex: 1,
      display: "flex",
      justifyContent: "space-between",
      alignItems: "center",
    },
    name: {
      marginBottom: 0.5,
      "::after": isRequired
        ? {
            content: '"*"',
            color: "#971E28",
            marginLeft: 0.25,
          }
        : {},
    },
    key: {
      fontSize: "0.75rem",
      color: "#4B4444",
    },
    onDrag: (theme) => ({
      background: `${theme.palette.primary.main}70`,
      border: `1px solid ${theme.palette.primary.main}`,
    }),
  };
};
