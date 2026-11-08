import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import type { DataSource } from "@/types/dataSource";
import Divider from "@mui/material/Divider";
import { DataSourceTypeTag } from "./DataSourceTypeTag";
import { dataSourceCardStyles as styles } from "./DataSourceCard.style";


export const DataSourceCard: React.FC<{ dataSource: DataSource }> = function ({
  dataSource,
}) {
  return (
    <Card variant="outlined" sx={styles["card"]}>
      <Typography sx={styles["name"]}>{dataSource.name}</Typography>
      <Typography sx={styles["description"]}>
        {dataSource.description}
      </Typography>
      <Divider sx={styles["divider"]} />
      <Box sx={styles["information"]}>
        <DataSourceTypeTag type={dataSource.type} />
        <Typography sx={styles["usage"]}>Used in ## forms</Typography>
      </Box>
    </Card>
  );
};
