import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import { DataSourcesCatalogPanel } from "../panels/DataSourceCatalogPanel/DataSourcesCatalogPanel";
import { DataSourceCanvas } from "./DataSourceCanvas";
import { useState } from "react";
import type { DataSource } from "@/types/dataSource";
import { Border } from "@/constants/colors";

const styles: Styles = {
  layout: {
    display: "grid",
    gridTemplateColumns: "minmax(380px, 13rem) auto",
    border: `1px solid ${Border.Primary}`,
    height: "100%",
  },
};

export const ReferenceData: React.FC = function () {
  const [dataSource, setDataSource] = useState<DataSource | undefined>();

  return (
    <Box sx={styles.layout}>
      <Box sx={{ borderRight: `1px solid ${Border.Primary}` }}>
        <DataSourcesCatalogPanel
          dataSource={dataSource}
          onChange={(ds) => setDataSource(ds)}
        />
      </Box>
      <DataSourceCanvas dataSource={dataSource} />
    </Box>
  );
};
