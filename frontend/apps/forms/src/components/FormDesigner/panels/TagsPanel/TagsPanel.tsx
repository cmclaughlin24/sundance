import Box from "@mui/material/Box";
import { TagsPanelHeader } from "./TagsPanelHeader";
import type { Styles } from "@/types/styles";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";
import { TagsPanelContent } from "./TagPanelContent";
import { TagsPanelCard } from "./TagsPanelCard";

const styles: Styles = {
  panel: {
    p: 5,
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
};

export interface TagsPanelProps extends React.PropsWithChildren {
  sx?: SxProps<Theme>;
}

interface TagsPanelComponent extends React.FC<TagsPanelProps> {
  Header: typeof TagsPanelHeader;
  Content: typeof TagsPanelContent;
  Card: typeof TagsPanelCard;
}

const TagsPanel: TagsPanelComponent = function ({ children, sx }) {
  return (
    <Box component="section" sx={mergeSx(styles.panel, sx)}>
      {children}
    </Box>
  );
};

TagsPanel.Header = TagsPanelHeader;
TagsPanel.Content = TagsPanelContent;
TagsPanel.Card = TagsPanelCard;

export default TagsPanel;
