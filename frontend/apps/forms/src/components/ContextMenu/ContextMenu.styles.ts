import { Background } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const contextMenuStyles: Styles = {
  backdrop: {
    position: "fixed",
    width: "100vw",
    height: "100vh",
    left: 0,
    right: 0,
    zIndex: 10_000,
  },
  contextMenu: {
    position: "relative",
    background: Background.Primary,
    borderRadius: 2.5,
    boxShadow: "0 4px 12px rgba(0, 0, 0, 0.12), 0 1px 3px rgba(0, 0, 0, 0.08)",
    p: 0.5
  },
};
