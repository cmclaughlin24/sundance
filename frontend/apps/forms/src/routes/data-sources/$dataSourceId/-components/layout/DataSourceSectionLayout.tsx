import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

export type DataSourceAttributeFormLayoutProps = React.PropsWithChildren<{
  title: string;
}>;

export const DataSourceSectionLayout: React.FC<DataSourceAttributeFormLayoutProps> =
  function ({ title, children }) {
    return (
      <Box sx={{ width: "100%" }}>
        <Typography>{title}</Typography>
        {children}
      </Box>
    );
  };
