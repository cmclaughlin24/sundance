import Box from "@mui/material/Box";
import { cardListStyles as styles } from "./CardList.styles";
import { motion } from "motion/react";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";

export interface CardListProps<T> {
  data: T[] | null;
  keyFn: (item: T) => string;
  onClick: (item: T) => void;
  children: (item: T) => React.ReactNode | undefined;
  sx?: SxProps<Theme>;
  disableAnimation?: boolean;
}

export function CardList<T>({
  data,
  keyFn,
  onClick,
  children,
  sx,
  disableAnimation = false,
}: CardListProps<T>) {
  return (
    <Box component="ul" sx={mergeSx(styles["list"], sx)}>
      {data?.map((item) => (
        <Box
          component={motion.li}
          whileHover={
            !disableAnimation ? { border: "1px solid white", y: -5 } : {}
          }
          transition={
            !disableAnimation ? { duration: 0.4, type: "spring" } : {}
          }
          sx={styles["item"]}
          onClick={() => onClick(item)}
          key={keyFn(item)}
        >
          {children(item)}
        </Box>
      ))}
    </Box>
  );
}
