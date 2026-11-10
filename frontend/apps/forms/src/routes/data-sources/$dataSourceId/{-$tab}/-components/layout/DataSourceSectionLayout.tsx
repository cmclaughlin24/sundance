import { Border } from "@/constants/colors";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  section: {
    p: 2.5,
    border: `1px solid ${Border.Primary}`,
  },
};

export type DataSourceAttributeFormLayoutProps = React.PropsWithChildren<{
  title: string;
}>;

export const DataSourceSectionLayout: React.FC<DataSourceAttributeFormLayoutProps> =
  function ({ title, children }) {
    return (
      <Box component="section" sx={styles.section}>
        <Typography>{title}</Typography>
        {children}
      </Box>
    );
  };
