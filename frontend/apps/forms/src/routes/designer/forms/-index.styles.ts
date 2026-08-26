import type { Styles } from "@/types/styles";

export const designerStyles: Styles = {
  page: {
    maxWidth: 1440,
    display: "flex",
    flexDirection: "column",
    gap: 5,
  },
  toolbar: {
    marginBottom: 3,
    display: "grid",
    gridTemplateColumns: "repeat(3, 1fr)",
    gap: 2,
  },
  toolbarLeft: {
    gridColumn: "span 2",
    display: "flex",
    gap: 2,
  },
  searchInput: {
    width: "100%",
  },
};
