import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  dropZone: (theme) => ({
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    border: `1px dashed ${theme.palette.primary.main}`,
    borderRadius: 2.5,
    background: `${theme.palette.primary.main}20`,
    px: 3,
    height: "4.25rem",
  }),
  text: (theme) => ({
    fontWeight: 600,
    color: theme.palette.primary.main,
  }),
};

export const DropZoneIndicator: React.FC<{ text: string }> = function ({
  text,
}) {
  return (
    <Box sx={styles.dropZone}>
      <Typography sx={styles.text}>{text}</Typography>
    </Box>
  );
};
