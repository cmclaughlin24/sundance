import type { TagGroup } from "@/utils/tag";
import Box from "@mui/material/Box";
import { TagsPanel } from "../panels/TagsPanel";
import type { Styles } from "@/types/styles";
import Typography from "@mui/material/Typography";

const styles: Styles = {
  list: {
    p: 0,
    m: 0,
  },
  item: {
    listStyle: "none",
    mb: 0.5,
  },
};

export interface TagGroupListProps {
  group: TagGroup;
}

export const TagGroupList: React.FC<TagGroupListProps> = function ({ group }) {
  return (
    <>
      <Typography component="h4" sx={{ mb: 1 }}>
        {group.title} {group.keyPath ?? `(${group.keyPath})`}
      </Typography>
      <Box component="ul" sx={styles.list}>
        {group.items?.map((t) => (
          <Box component="li" sx={styles.item} key={t.id}>
            <TagsPanel.Card title={t.keyPath} description={t.displayName} />
          </Box>
        ))}
      </Box>
    </>
  );
};
