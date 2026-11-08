import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const dataSourceTypeCardStyles: Styles = {
  card: (theme) => ({
    p: 2,
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    textAlign: "center",
    gap: 0.5,
    borderRadius: 2.5,
    border: `1px dashed ${Border.Primary}`,
    overflow: "hidden",
    ":hover:not(:has(li:hover, button:hover))": {
      cursor: "pointer",
      borderColor: theme.palette.primary.main,
      borderStyle: "solid",
    },
  }),
  label: {
    fontWeight: 600,
  },
  description: {
    flex: 1,
    fontSize: "0.75rem",
    display: "flex",
    alignItems: "center",
  },
};
