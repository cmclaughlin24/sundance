import type { Styles } from "@/types/styles";

export const ruleActionSectionStyles: Styles = {
  section: {
    display: "flex",
    flexDirection: "column",
    gap: 1.25,
    p: 2.5,
    borderRadius: "10px",
  },
  sectionHeader: {
    fontSize: "0.75rem",
    fontWeight: 700,
    letterSpacing: "0.08em",
    textTransform: "uppercase",
    color: "text.secondary",
  },
  actionBox: {
    display: "flex",
    alignItems: "center",
    flexWrap: "wrap",
    gap: 1.5,
  },
  actionText: {
    fontSize: "0.9rem",
    color: "text.secondary",
  },
};
