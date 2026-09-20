import { Background, Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const tagsPanelCardStyles = (
  isSelected: boolean,
  isClickable: boolean = false,
): Styles => ({
  card: (theme) => ({
    mb: 2.5,
    p: 1.5,
    display: 'flex',
    borderRadius: 2.5,
    gap: 1,
    borderColor: isSelected ? `${theme.palette.primary.main}` : Border.Primary,
    borderStyle: isSelected ? "solid" : "dashed",
    background: isSelected
      ? `${theme.palette.primary.main}25`
      : Background.Primary,
    overflow: "hidden",
    ":hover:not(:has(li:hover, button:hover))": isClickable
      ? {
          cursor: "pointer",
          border: `1px solid ${theme.palette.primary.main}`,
        }
      : {},
  }),
  title: {
    mb: 0.5,
  },
  description: {
    fontSize: "0.75rem",
    color: "#4B4444",
  },
  content: {
    flex: 1,
  },
});
