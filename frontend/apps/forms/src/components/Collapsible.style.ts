import type { Styles } from "@/types/styles";

export const collapsibleStyles: Styles = {
  collapisble: {},
  summary: {
    marginBottom: 1.5,
    display: "flex",
    alignItems: "center",
    gap: 0.5,
    ":hover": {
      cursor: "pointer",
    },
  },
  content: {
    overflow: "hidden",
  },
  button: (theme) => ({
    padding: 0.5,
    borderRadius: 1.5,
    color: "#4B4444",
    ":hover": {
      color: theme.palette.primary.main,
      background: `${theme.palette.primary.main}25`,
    },
  }),
};
