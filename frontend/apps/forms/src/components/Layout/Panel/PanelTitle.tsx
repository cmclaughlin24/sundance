import type { SxProps, Theme } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import { mergeSx } from "merge-sx";

const styles: Readonly<SxProps<Theme>> = {
  fontWeight: 600,
  fontSize: "1.25rem",
};

export const PanelTitle: React.FC<{ title: string; sx?: SxProps<Theme> }> =
  function ({ title, sx }) {
    return (
      <Typography variant="h3" sx={mergeSx(styles, sx)}>
        {title}
      </Typography>
    );
  };
