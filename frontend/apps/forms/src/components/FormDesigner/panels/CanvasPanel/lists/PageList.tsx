import type { IPage } from "@/types/page";
import Box from "@mui/material/Box";
import { PageItem } from "./PageItem";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  list: {
    margin: 0,
    padding: 0,
  },
};

export const PageList: React.FC<{ pages: IPage[] }> = function ({ pages }) {
  return (
    <Box component="ul" sx={styles.list}>
      {pages?.map((page) => (
        <PageItem page={page} key={page.id} />
      ))}
    </Box>
  );
};
