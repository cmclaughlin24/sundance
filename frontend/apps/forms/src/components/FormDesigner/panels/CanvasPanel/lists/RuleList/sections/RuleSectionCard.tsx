import type { Styles } from "@/types/styles";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  section: {
    display: "flex",
    flexDirection: "column",
    gap: 1.25,
    p: 2.5,
    borderRadius: "10px",
  },
  sectionHeader: {
    fontWeight: 600,
    color: "text.secondary",
  },
};

export type RuleSectionCardProps = React.PropsWithChildren<{
  title: string;
}>;

export const RuleSectionCard: React.FC<RuleSectionCardProps> = function ({
  title,
  children,
}) {
  return (
    <Card variant="outlined" sx={styles.section}>
      <Typography sx={styles.sectionHeader}>{title}</Typography>
      {children}
    </Card>
  );
};
