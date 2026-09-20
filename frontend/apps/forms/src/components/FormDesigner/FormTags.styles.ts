import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const formTagsStyles: Styles = {
  workspace: {
    display: "grid",
    gridTemplateColumns: "1fr minmax(208px, 13rem) 1fr",
    border: `1px solid ${Border.Primary}`,
    height: "100%",
  },
  badge: {
    fontWeight: 600,
    color: "#7E7472",
  },
  fields: {
    borderRight: `1px solid ${Border.Primary}`,
  },
  tags: {
    borderLeft: `1px solid ${Border.Primary}`,
  },
};
