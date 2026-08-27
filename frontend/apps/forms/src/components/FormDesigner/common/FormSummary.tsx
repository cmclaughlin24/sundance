import type { IPage } from "@/types/page";
import Typography from "@mui/material/Typography";
import * as ArrayUtils from "@/utils/array";
import type { Styles } from "@/types/styles";

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

export const FormSummary: React.FC<{ pages: IPage[] }> = function ({ pages }) {
  const { fields, rules } = getCount(pages);

  return (
    <Typography sx={styles.title}>
      <Typography component="span" sx={styles.bold}>
        Form Layout:{" "}
      </Typography>
      {fields} fields · {rules} rules
    </Typography>
  );
};

function getCount(pages: IPage[]): { fields: number; rules: number } {
  let fields: number = 0;
  let rules: number = 0;

  if (!ArrayUtils.hasLengthGreaterThan(pages, 0)) {
    return { fields, rules };
  }

  for (const page of pages) {
    for (const section of page.sections) {
      for (const element of section.elements) {
        fields++;
        rules += element.rules.length;
      }
      rules += section.rules.length;
    }
    rules += page.rules.length;
  }

  return { fields, rules };
}
