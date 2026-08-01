import type { Variants } from "motion/react";

export const pageVariants: Variants = {
  initial: { opacity: 0, x: -50 },
  animate: {
    opacity: 1,
    x: 0,
    transition: { duration: 0.3, type: "spring", bounce: 0.4 },
  },
};

export const sectionVariants: Variants = {
  initial: { opacity: 0, y: -25 },
  animate: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.3, type: "spring", bounce: 0.4 },
  },
};

export const elementVariants: Variants = {
  initial: { opacity: 0, y: -25 },
  animate: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.3, type: "spring", bounce: 0.4 },
  },
};
