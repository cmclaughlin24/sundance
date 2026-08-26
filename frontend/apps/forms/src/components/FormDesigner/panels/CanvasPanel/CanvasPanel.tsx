import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";

export const CanvasPanel: React.FC = function () {
  return <Panel sx={canvasPanelStyles.canvas}>Canvas Panel</Panel>;
};
