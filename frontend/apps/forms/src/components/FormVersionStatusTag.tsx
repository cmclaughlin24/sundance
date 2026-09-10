import type { FormVersionStatus } from "@/types/formVersion";
import { Tag } from "./Tag";
import { mergeSx } from "merge-sx";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  base: {
    fontSize: "1rem",
    fontWeight: 600,
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
  retired: (theme) => ({
    color: theme.palette.secondary.main,
    background: `${theme.palette.secondary.main}20`,
  }),
};

const versionStatusLabels: Readonly<Record<FormVersionStatus, string>> = {
  draft: "Draft",
  active: "Published",
  retired: "Retired",
};

export const FormVersionTag: React.FC<{
  status: FormVersionStatus;
  version?: number;
}> = function ({ status, version }) {
  let text = versionStatusLabels[status];

  if (version != null) {
    text += ` v${version}`;
  }

  return <Tag sx={mergeSx(styles.base, styles[status])}>{text}</Tag>;
};
