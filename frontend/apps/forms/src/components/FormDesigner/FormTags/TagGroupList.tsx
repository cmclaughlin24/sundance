import type { TagGroup } from "@/utils/tag";
import Box from "@mui/material/Box";
import type { Styles } from "@/types/styles";
import Typography from "@mui/material/Typography";
import { TagGroupListItem } from "./TagGroupListItem";

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
    <Box sx={{ mb: 2.5 }}>
      <Typography component="h4" sx={{ mb: 1 }}>
        {group.title} {group.keyPath ?? `(${group.keyPath})`}
      </Typography>
      <Box component="ul" sx={styles.list}>
        {group.items?.map((t) => (
          <Box component="li" sx={styles.item} key={t.id}>
            <TagGroupListItem tag={t} />
          </Box>
        ))}
      </Box>
    </Box>
  );
};
