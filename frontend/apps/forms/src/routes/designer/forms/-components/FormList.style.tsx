import type { SxProps, Theme } from "@mui/material/styles";

export const formListStyles: Record<string, SxProps<Theme>> = {
  list: {
    margin: 0,
    padding: 0,
    display: "grid",
    gridTemplateColumns: {
      xs: "1fr",
      sm: "repeat(2, 1fr)",
      md: "repeat(3, 1fr)",
    },
    gap: 2,
  },
  item: {
    borderRadius: "10px",
    cursor: "pointer",
    listStyle: "none",
  },
};
