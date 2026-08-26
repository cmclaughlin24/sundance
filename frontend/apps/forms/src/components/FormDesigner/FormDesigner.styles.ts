import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const formDesignerStyles: Styles = {
  container: {
    display: "grid",
    gridTemplateColumns: "minmax(296px, 18.5rem) auto minmax(296px, 18.5rem)",
    border: `1px solid ${Border.Primary}`,
  },
};
