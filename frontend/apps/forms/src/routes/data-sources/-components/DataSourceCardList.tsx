import { CardList } from "@/components/CardList";
import {
  DataSourceTypeCard,
  type DataSourceTypeCardProps,
} from "./DataSourceTypeCard";
import HttpIcon from "@mui/icons-material/Http";
import MenuIcon from "@mui/icons-material/Menu";
import StorageIcon from "@mui/icons-material/Storage";
import AccessAlarm from "@mui/icons-material/AccessAlarm";
import type { Styles } from "@/types/styles";
import type { DataSourceType } from "@/types/dataSource";

const dataSourceTypes: DataSourceTypeCardProps[] = [
  {
    type: "scheduled",
    description: "Periodic background sync with TTL cache.",
    icon: <AccessAlarm fontSize="large" />,
  },
  {
    type: "webhook",
    description: "On-demand HTTP fetch parameterized by form.",
    icon: <HttpIcon fontSize="large" />,
  },
  {
    type: "data-lake",
    description: "Direct BigQuery SQL query to lakehouse.",
    icon: <StorageIcon fontSize="large" />,
  },
  {
    type: "static",
    description: "In-memory key-value list stored in schema.",
    icon: <MenuIcon fontSize="large" />,
  },
];

const styles: Styles = {
  typesList: {
    gridTemplateColumns: {
      xs: "1fr",
      sm: "repeat(3, 1fr)",
      md: "repeat(auto-fit, minmax(5rem, 1fr))",
    },
    gap: 1,
    width: "100%",
  },
  selectedTypeCard: (theme) => ({
    borderColor: `${theme.palette.primary.main}`,
    borderStyle: "solid",
    background: `${theme.palette.primary.main}25`,
  }),
};

export const DataSourceCardList: React.FC<{
  selection: DataSourceType;
  onSelection: (type: DataSourceType) => void;
}> = function ({ selection, onSelection }) {
  return (
    <CardList
      data={dataSourceTypes}
      keyFn={({ type }) => type}
      onClick={(prop) => onSelection(prop.type)}
      sx={styles.typesList}
      disableAnimation
    >
      {(props) => (
        <DataSourceTypeCard
          {...props}
          sx={selection === props.type ? styles.selectedTypeCard : undefined}
        />
      )}
    </CardList>
  );
};
