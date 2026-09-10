import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const ruleConditionRowStyles: Styles = {
  conditionRow: {
    display: "flex",
    alignItems: "center",
    flexWrap: "wrap",
    gap: 1.25,
    borderRadius: 2.5,
    border: `1px dashed ${Border.Primary}`,
    px: 1,
    py: 1.25,
  },
  joinOpSelect: {
    minWidth: "85px",
  },
  fieldSelect: {
    minWidth: "220px",
    flex: 1,
  },
  operatorSelect: {
    minWidth: "160px",
  },
  valueInput: {
    minWidth: "160px",
    flex: 1,
  },
  deleteButton: (theme) => ({
    color: "#6b7280",
    ":hover": {
      color: theme.palette.error.main,
      background: `${theme.palette.error.main}15`,
    },
  }),
};
