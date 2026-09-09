import type { Styles } from "@/types/styles";

export const ruleItemHeaderStyles: Styles = {
  ruleItemHeader: {
    display: "flex",
    justifyContent: "space-between",
  },
  toggle: {
    mb: 1.5,
    display: "flex",
    alignItems: "center",
    gap: 0.5,
    ":hover": {
      cursor: "pointer",
      textDecoration: "underline",
    },
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
