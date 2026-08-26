import { Border } from "@/constants/colors";
import type { SxProps, Theme } from "@mui/material/styles";

export const formDesignerStyles: Readonly<Record<string, SxProps<Theme>>> = {
  container: {
    display: "grid",
    gridTemplateColumns: "minmax(296px, 18.5rem) auto minmax(296px, 18.5rem)",
    border: `1px solid ${Border.Primary}`,
  },
};
