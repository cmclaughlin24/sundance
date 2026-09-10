import type { Styles } from "@/types/styles";

export const ruleConditionsSectionStyles: Styles = {
  content: {
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
  conditionsList: {
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
  emptyState: (theme) => ({
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    border: `1px dashed ${theme.palette.primary.main}`,
    borderRadius: 2.5,
    px: 3,
    height: "4.25rem",
  }),
  addButton: {
    alignSelf: "center",
    textTransform: "none",
  },
};
