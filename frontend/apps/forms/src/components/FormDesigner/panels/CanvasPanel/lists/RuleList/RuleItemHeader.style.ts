import type { Styles } from "@/types/styles";

export const ruleItemHeaderStyles: Styles = {
  ruleItemHeader: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  toggle: {
    display: "flex",
    alignItems: "center",
    gap: 0.5,
    ":hover": {
      cursor: "pointer",
      textDecoration: "underline",
    },
  },
  titleContainer: {
    display: "flex",
    alignItems: "center",
    gap: 1.25,
    flexWrap: "wrap",
  },
  titleText: {
    fontWeight: 600,
    lineHeight: 1.2,
  },
  tag: {
    alignSelf: "center",
    fontWeight: 600,
    fontSize: "0.75rem",
    letterSpacing: "0.02em",
    px: 1,
    py: 0.25,
    color: "white",
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
