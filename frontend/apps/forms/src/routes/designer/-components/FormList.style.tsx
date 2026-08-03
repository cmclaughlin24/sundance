import type { SxProps, Theme } from "@mui/material/styles";

export const formListStyles: Record<string, SxProps<Theme>> = {
  list: {
    margin: 0,
    padding: 0,
    display: "flex",
    gap: 2,
  },
  item: {
    listStyle: "none",
  },
};
