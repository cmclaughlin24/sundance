import type { DataSourceType } from "@/types/dataSource";
import { dataSourceTypeLabels } from "../-constants/dataSourceLabels";
import Typography from "@mui/material/Typography";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";
import Box from "@mui/material/Box";
import { dataSourceTypeCardStyles as styles } from "./DataSourceTypeCard.style";

export interface DataSourceTypeCardProps {
  type: DataSourceType;
  icon: React.ReactNode;
  description: string;
  sx?: SxProps<Theme>;
}

export const DataSourceTypeCard: React.FC<DataSourceTypeCardProps> = function ({
  type,
  icon,
  description,
  sx,
}) {
  return (
    <Box sx={mergeSx(styles.card, sx)}>
      {icon}
      <Typography component="span" sx={styles.label}>
        {dataSourceTypeLabels[type]}
      </Typography>
      <Typography component="span" sx={styles.description}>
        {description}
      </Typography>
    </Box>
  );
};
