import type { SxProps, Theme } from "@mui/material/styles";
import type { Styles } from "@/types/styles";
import { Border } from "@/constants/colors";

const baseStyles: SxProps<Theme> = {
  display: "flex",
  flexDirection: "column",
  alignItems: "start",
  gap: 5,
};

export const rendererStyles: Styles = {
  form: {
    ...baseStyles,
    maxWidth: "901px",
    position: "relative",
    flex: 1,
    padding: 6.5,
    paddingBottom: 0,
    minHeight: 0,
  },
  page: { ...baseStyles, width: "100%" },
  section: { ...baseStyles, width: "100%" },
  repeatableSection: {
    borderLeft: `3px solid ${Border.Primary}`,
    paddingLeft: 2.5,
    width: "unset",
  },
  repeatableSectionActionTitle: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    width: "100%",
  },
};
