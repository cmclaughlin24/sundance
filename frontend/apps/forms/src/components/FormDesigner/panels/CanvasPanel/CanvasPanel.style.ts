import type { Styles } from "@/types/styles";

export const canvasPanelStyles: Styles = {
  canvas: {
    padding: 5,
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
  toolbar: {
    display: "flex",
    justifyContent: "space-between",
  },
  buttons: {
    "> button": {
      color: "#4B4444",
    },
    "> button:not(:last-of-type)": {
      marginRight: 1,
    },
  },
};
