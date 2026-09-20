import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  header: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  h3: {
    fontWeight: 600,
    fontSize: "1.25rem",
  },
};

export interface TagsPanelHeaderProps {
  title: string;
  slot?: React.ReactNode;
}

export const TagsPanelHeader: React.FC<TagsPanelHeaderProps> = function ({
  title,
  slot,
}) {
  return (
    <Box sx={styles.header}>
      <Typography variant="h3" sx={styles.h3}>
        {title}
      </Typography>
      {slot}
    </Box>
  );
};
