import type { DataSource } from "@/types/dataSource";
import Box from "@mui/material/Box";
import { canvasPanelStyles } from "../panels/CanvasPanel/CanvasPanel.style";
import Typography from "@mui/material/Typography";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import { useState } from "react";

enum DataSourceCanvasTab {
  Attributes = "attributes",
  ResponseExtraction = "response-extraction",
  TestBench = "TestBench",
  FormBindings = "bindings",
}

export const DataSourceCanvas: React.FC<{
  dataSource: DataSource | undefined;
}> = function ({ dataSource }) {
  const [tab, setTab] = useState<DataSourceCanvasTab>(
    DataSourceCanvasTab.Attributes,
  );

  const handleTabChange = (
    _event: React.SyntheticEvent,
    tab: DataSourceCanvasTab,
  ) => setTab(tab);

  return (
    <Box sx={canvasPanelStyles.canvas}>
      <Typography>{dataSource?.name}</Typography>
      <Typography>{dataSource?.description}</Typography>
      <Box>
        <Tabs value={tab} onChange={handleTabChange}>
          <Tab
            label="1. Configuration & Attributes"
            value={DataSourceCanvasTab.Attributes}
          />
          <Tab
            label="2. Response Extraction"
            value={DataSourceCanvasTab.ResponseExtraction}
          />
          <Tab
            label="3. Live Test Bench"
            value={DataSourceCanvasTab.TestBench}
          />
          <Tab
            label="4. Form Field Bindings"
            value={DataSourceCanvasTab.FormBindings}
          />
        </Tabs>
      </Box>
    </Box>
  );
};
