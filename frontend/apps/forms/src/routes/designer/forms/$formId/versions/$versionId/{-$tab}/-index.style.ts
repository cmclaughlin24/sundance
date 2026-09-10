import type { Styles } from "@/types/styles";

export const formDesignerPageStyles: Styles = {
  page: {
    marginTop: 2.5,
    padding: 5,
    width: "100%",
    maxWidth: 1840,
    display: "flex",
    flexDirection: "column",
    alignItems: "start",
    gap: 2.5,
  },
  header: {
    display: "flex",
    justifyContent: "space-between",
    width: "100%",
  },
  headerTitle: {
    display: "flex",
    gap: 2,
  },
  headerIcons: {
    mt: 2,
    "> *:not(:last-child)": {
      mr: 2,
    },
  },
  headerActions: {
    "> button": {
      marginLeft: 2.5,
    },
  },
};
