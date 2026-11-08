import { Tag } from "@/components/Tag";
import type { DataSourceType } from "@/types/dataSource";
import type { Styles } from "@/types/styles";
import { mergeSx } from "merge-sx";
import { dataSourceTypeLabels } from "../-constants/dataSourceLabels";

const styles: Styles = {
  base: {
    fontSize: "0.75rem",
    px: 1,
    py: 0.5,
  },
  scheduled: {
    background: "rgba(192, 132, 252, 0.18)",
    color: "#d8b4fe",
    border: " 1px solid rgba(192, 132, 252, 0.4)",
  },
  webhook: {
    background: "rgba(6, 182, 212, 0.18)",
    color: "#67e8f9",
    border: "1px solid rgba(6, 182, 212, 0.4)",
  },
  "data-lake": {
    background: "rgba(251, 191, 36, 0.18)",
    color: "#fde68a",
    border: "1px solid rgba(251, 191, 36, 0.4)",
  },
  static: {
    background: "rgba(52, 211, 153, 0.18)",
    color: "#6ee7b7",
    border: "1px solid rgba(52, 211, 153, 0.4)",
  },
};

export const DataSourceTypeTag: React.FC<{ type: DataSourceType }> = function ({
  type,
}) {
  return (
    <Tag sx={mergeSx(styles.base, styles[type])}>
      {dataSourceTypeLabels[type]}
    </Tag>
  );
};
