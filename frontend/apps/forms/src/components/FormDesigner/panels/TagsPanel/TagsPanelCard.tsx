import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import type { KeyboardEventHandler, MouseEventHandler } from "react";
import { tagsPanelCardStyles } from "./TagsPanelCard.style";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";

export type TagsPanelCardProps = React.PropsWithChildren<{
  title: string;
  description: string;
  isSelected?: boolean;
  onClick?: MouseEventHandler<HTMLDivElement>;
  onKeyDown?: KeyboardEventHandler<HTMLDivElement>;
  slotProps?: Partial<{
    title: { sx: SxProps<Theme> };
    content: { sx: SxProps<Theme> };
  }>;
}>;

export const TagsPanelCard: React.FC<TagsPanelCardProps> = function ({
  title,
  description,
  isSelected = false,
  children,
  slotProps = {},
  onClick,
  onKeyDown,
}) {
  const styles = tagsPanelCardStyles(isSelected, !!onClick || !!onKeyDown);

  const handleKeyDown: KeyboardEventHandler<HTMLDivElement> = (event) => {
    if (!onKeyDown) {
      return;
    }

    if (event.target !== event.currentTarget) {
      return;
    }

    if (event.key !== "Enter" && event.key !== " ") {
      return;
    }

    event.preventDefault();
    event.stopPropagation();
    onKeyDown(event);
  };

  return (
    <Card
      variant="outlined"
      sx={styles.card}
      onClick={onClick}
      onKeyDown={handleKeyDown}
      role="button"
      tabIndex={0}
      aria-pressed={isSelected}
    >
      <Box>
        <Typography sx={mergeSx(styles.title, slotProps.title?.sx)}>
          {title}
        </Typography>
        <Typography sx={styles.description}>{description}</Typography>
      </Box>
      <Box sx={mergeSx(styles.content, slotProps.content?.sx)}>{children}</Box>
    </Card>
  );
};
