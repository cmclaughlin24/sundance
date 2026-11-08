import type { DataSource } from "@/types/dataSource";
import type { Styles } from "@/types/styles";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  card: {
    textAlign: "left",
    display: "flex",
    flexDirection: "column",
    gap: 0.5,
  },
};

export const DataSourceCatalogItem: React.FC<{
  dataSource: DataSource;
  onClick: (dataSource: DataSource) => void;
}> = function ({ dataSource, onClick }) {
  return (
    <Card
      component="button"
      sx={styles.card}
      onClick={() => onClick(dataSource)}
    >
      <Typography>{dataSource.name}</Typography>
      <Typography>{dataSource.description}</Typography>
    </Card>
  );
};
