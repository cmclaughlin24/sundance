import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";
import { PageList } from "./lists/PageList";
import { useFormPagesSnapshot } from "@/store/formDesigner";

export const CanvasPanel: React.FC = function () {
  const pages = useFormPagesSnapshot();

  return (
    <Panel sx={canvasPanelStyles.canvas}>
      <PageList pages={pages} />
    </Panel>
  );
};
