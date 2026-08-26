import type { Styles } from "@/types/styles";

export const pageTitleStyles: Styles = {
  container: {
    display: "flex",
    flexDirection: "column",
    gap: 1,
    paddingBottom: 2.5,
  },
  name: {
    fontSize: "2.5rem",
    fontWeight: 300,
    letterSpacing: 0,
    lineHeight: 1.3,
  },
  description: {
    fontSize: "1rem",
    fontWeight: 400,
    letterSpacing: 0,
    lineHeight: 1.25,
  },
};
