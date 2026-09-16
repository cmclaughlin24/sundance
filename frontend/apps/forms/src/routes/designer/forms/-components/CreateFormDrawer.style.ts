import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";

export const createFormDrawerStyles: Styles = {
  form: {
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
    height: "100%",
  },
  title: {
    p: 4,
    pr: 9,
    display: "flex",
    flexDirection: "column",
    gap: 1.5,
  },
  h3: {
    fontSize: "2rem",
    weight: 300,
    letterSpacing: 0,
  },
  fields: {
    p: 4,
    pt: 0,
    pr: 9,
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
    height: "100%",
  },
  buttons: {
    p: 4,
    mt: "auto",
    borderTop: `1px solid ${Border.Primary}`,
    display: "flex",
    justifyContent: "end",
    gap: 2.5,
  },
};
