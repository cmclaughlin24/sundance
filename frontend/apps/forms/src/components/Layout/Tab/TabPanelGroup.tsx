import Box from "@mui/material/Box";
import type { TabPanelProps } from "./TabPanel";
import {
  Children,
  isValidElement,
  useEffect,
  useRef,
  type ReactElement,
} from "react";
import { AnimatePresence, motion, type Variants } from "motion/react";
import type { SxProps, Theme } from "@mui/material/styles";
import { mergeSx } from "merge-sx";

export type TabPanelComponent<T> = ReactElement<TabPanelProps<T>>;

export interface TabPanelGroupProps<T> {
  active: T;
  order: T[];
  children: TabPanelComponent<T>[];
  sx?: SxProps<Theme>;
}

const variants: Variants = {
  enter: (dir: 1 | -1) => ({ x: dir > 0 ? "100%" : "-100%", opacity: 0 }),
  center: { x: 0, opacity: 1 },
  exit: (dir: 1 | -1) => ({ x: dir > 0 ? "-100%" : "100%", opacity: 0 }),
};

export const TabPanelGroup = function <T>({
  active,
  order,
  children,
  sx,
}: TabPanelGroupProps<T>) {
  const prevActiveRef = useRef<T>(active);
  const currentIndex = order.indexOf(active);
  const prevIndex = order.indexOf(prevActiveRef.current);
  const direction: 1 | -1 = currentIndex >= prevIndex ? 1 : -1;

  useEffect(() => {
    prevActiveRef.current = active;
  }, [active]);

  const tab = Children.toArray(children).find((child) => {
    return (
      isValidElement(child) &&
      (child as TabPanelComponent<T>).props.value === active
    );
  }) as TabPanelComponent<T> | undefined;

  return (
    <Box
      sx={mergeSx(
        { overflow: "hidden", position: "relative", width: "100%" },
        sx,
      )}
    >
      <AnimatePresence initial={false} custom={direction} mode="wait">
        <Box
          key={String(active)}
          component={motion.div}
          custom={direction}
          variants={variants}
          initial="enter"
          animate="center"
          exit="exit"
          transition={{
            duration: 0.35,
            ease: "easeInOut",
            type: "spring",
            bounce: 0.175,
          }}
          sx={{ py: 2.5, height: "100%" }}
          role="tabpanel"
          id={`tab-panel-${active}`}
          aria-labelledby={`tab-panel-${active}`}
        >
          {tab}
        </Box>
      </AnimatePresence>
    </Box>
  );
};
