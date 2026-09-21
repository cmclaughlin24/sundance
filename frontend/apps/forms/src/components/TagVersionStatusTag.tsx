import { Tag } from "./Tag";
import { mergeSx } from "merge-sx";
import type { Styles } from "@/types/styles";
import type { TagVersionStatus } from "@/types/tag";

const styles: Styles = {
  base: {
    fontSize: "0.75rem",
    px: 1,
    py: 0.5,
  },
  draft: (theme) => ({
    color: theme.palette.primary.main,
    background: `${theme.palette.primary.main}20}`,
  }),
  active: {
    color: "#026F43",
    background: "#CCFFEE30",
  },
  deprecated: {
    color: "#865603",
    background: "#FFDE80",
  },
  retired: (theme) => ({
    color: theme.palette.secondary.main,
    background: `${theme.palette.secondary.main}20`,
  }),
};

const versionStatusLabels: Readonly<Record<TagVersionStatus, string>> = {
  draft: "Draft",
  active: "Active",
  deprecated: "Deprecated",
  retired: "Retired",
};

export const TagVersionStatusTag: React.FC<{
  status: TagVersionStatus;
  version?: number;
}> = function ({ status, version }) {
  let text = versionStatusLabels[status];

  if (version != null) {
    text += ` v${version}`;
  }

  return <Tag sx={mergeSx(styles.base, styles[status])}>{text}</Tag>;
};
