import type { Styles } from "@/types/styles";

export const dataSourceCardStyles: Styles = {
  card: {
    flex: 1,
    borderRadius: "10px",
    padding: 2,
    boxShadow: "0 4px 8px 0 rgba(156, 145, 145, .25)",
    display: "flex",
    flexDirection: "column",
  },
  name: {
    fontSize: "1rem",
    fontWeight: 600,
    letterSpacing: 0,
    lineHeight: 1.5,
    marginBottom: 1,
  },
  description: {
    fontSize: "1rem",
    fontWeight: 300,
    letterSpacing: 0,
    lineHeight: 1.25,
    flex: 1,
  },
  divider: {
    my: 1,
  },
  information: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  usage: {
    fontSize: "0.75rem",
  },
};
