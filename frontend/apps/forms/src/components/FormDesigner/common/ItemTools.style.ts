import type { Styles } from "@/types/styles";

export const itemToolsStyles: Styles = {
  tools: {
    display: "inline-flex",
    alignItems: "center",
    gap: 0.5,
    padding: 0.5,
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
  deleteButton: (theme) => ({
    padding: 0.5,
    borderRadius: 1.5,
    color: "#4B4444",
    ":hover": {
      color: theme.palette.error.main,
      background: `${theme.palette.error.main}25`,
    },
  }),
};
