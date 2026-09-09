import type { Styles } from "@/types/styles";

export const ruleConditionsSectionStyles: Styles = {
  conditionsList: {
    display: "flex",
    flexDirection: "column",
    gap: 1.25,
  },
  emptyState: {
    dropZone: (theme) => ({
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      border: `1px dashed ${theme.palette.primary.main}`,
      borderRadius: 2.5,
      background: `${theme.palette.primary.main}20`,
      px: 3,
      height: "4.25rem",
    }),
  },
  addButton: {
    textTransform: "none",
    fontWeight: 600,
  },
};
