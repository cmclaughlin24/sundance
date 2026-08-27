import { Panel } from "@/components/layout/Panel";
import { canvasPanelStyles } from "./CanvasPanel.style";
import { PageList } from "./lists/PageList";
import { useFormPagesSnapshot } from "@/store/formDesigner";
import { FormSummary } from "../../common/FormSummary";

export const CanvasPanel: React.FC = function () {
  const pages = useFormPagesSnapshot();

  return (
    <Panel sx={canvasPanelStyles.canvas}>
      <FormSummary pages={pages} />
      <PageList pages={pages} />
    </Panel>
  );
};
