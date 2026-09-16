import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const drawerStyles: Styles = {
  header: {
    px: 1.5,
    py: 2.5,
    borderBottom: `1px solid ${Border.Primary}`,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
  title: {
    fontSize: "1.5rem",
    fontWeight: 600,
    flex: 1,
    textAlign: "center",
  },
  closeIcon: {
    width: 24,
  },
};


