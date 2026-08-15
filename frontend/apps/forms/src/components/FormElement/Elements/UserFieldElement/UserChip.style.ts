import type { SxProps, Theme } from "@mui/material/styles";

export const userChipStyles: Readonly<Record<string, SxProps<Theme>>> = {
  chip: {
    px: 1,
    py: 0.75,
    borderRadius: 2,
    border: "1px solid #CAC5C3",
    opacity: 1,
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    gap: 1,
  },
  label: {
    flex: 1,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  removeBtn: {
    background: "none",
    border: "none",
    cursor: "pointer",
    padding: 0,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    fontWeight: 400,
    color: "black",
    lineHeight: 1,
    minWidth: 2,
    transition: "opacity 0.2s ease-in-out",
    ":hover": {
      opacity: 0.7,
    },
  },
};
