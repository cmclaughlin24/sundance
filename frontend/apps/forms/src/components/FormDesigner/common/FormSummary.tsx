import type { IPage } from "@/types/page";
import Typography from "@mui/material/Typography";
import * as ArrayUtils from "@/utils/array";
import type { Styles } from "@/types/styles";
import type { IFlatRule } from "@/types/rule";

const styles: Styles = {
  title: {
    display: "inline",
    alignItems: "center",
    color: "#9C9191",
  },
  bold: {
    fontWeight: 600,
    color: "#7E7472",
  },
};

export const FormSummary: React.FC<{ pages: IPage[]; rules: IFlatRule[] }> =
  function ({ pages, rules }) {
    const fields = getFieldCount(pages);

    return (
      <Typography sx={styles.title}>
        <Typography component="span" sx={styles.bold}>
          Form Layout:{" "}
        </Typography>
        {fields} fields · {rules?.length ?? 0} rules
      </Typography>
    );
  };

function getFieldCount(pages: IPage[]): number {
  let fields: number = 0;

  if (!ArrayUtils.hasLengthGreaterThan(pages, 0)) {
    return fields;
  }

  for (const page of pages) {
    for (const section of page.sections) {
      fields += section.elements?.length ?? 0;
    }
  }

  return fields;
}
